/*

   Copyright 2017 Wenhui Shen <www.webx.top>

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.

*/

package echo

import (
	"errors"
	"fmt"
	"html"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const SnippetLineNumbers = 13

var workDir string

func init() {
	workDir, _ = os.Getwd()
	workDir = filepath.ToSlash(workDir) + "/"
}

func Wd() string {
	return workDir
}

func NewHTTPError(code int, msg ...string) *HTTPError {
	he := &HTTPError{Code: code, Message: http.StatusText(code)}
	if len(msg) > 0 {
		he.Message = msg[0]
	}
	return he
}

type HTTPError struct {
	Code    int
	Message string
}

// Error returns message.
func (e *HTTPError) Error() string {
	return e.Message
}

func NewPanicError(recovered interface{}, err error) *PanicError {
	return &PanicError{
		error:    err,
		Raw:      recovered,
		Traces:   make([]*Trace, 0),
		Snippets: make([]*SnippetGroup, 0),
	}
}

type PanicError struct {
	error
	Raw      interface{}
	Traces   []*Trace
	Snippets []*SnippetGroup
	debug    bool
}

type SnippetGroup struct {
	Path    string
	Index   int
	Snippet []*Snippet
}

func (sg *SnippetGroup) String() string {
	var s string
	var l int
	for _, snippet := range sg.Snippet {
		ns := strconv.Itoa(snippet.Number)
		if l == 0 {
			l = len(ns) + 5
		}
		if snippet.Current {
			ns = "[" + ns + "]"
		} else {
			ns = " " + ns
		}
		s += fmt.Sprintf("\n\t%*s%v", -l, ns, snippet.Code)
	}
	return s
}

func (sg *SnippetGroup) TableRow() string {
	var s string
	for _, snippet := range sg.Snippet {
		ns := strconv.Itoa(snippet.Number)
		cd := snippet.Code
		cd = html.EscapeString(cd)
		cd = strings.Replace(cd, "\t", "&nbsp;&nbsp;&nbsp;&nbsp;", -1)
		cd = strings.Replace(cd, " ", "&nbsp;", -1)
		if snippet.Current {
			ns = "<strong>" + ns + "</strong>"
			cd = "<strong>" + cd + "</strong>"
		}
		s += "<tr><td class='left'>" + ns + "</td><td class='right'>" + cd + "</td></tr>"
	}
	return s
}

type Snippet struct {
	Number  int
	Code    string
	Current bool
}

type Trace struct {
	Line   int
	File   string
	Func   string
	HasErr bool
}

func (p *PanicError) JSONString() string {
	return Dump(p, false)
}

func (p *PanicError) Error() string {
	return p.String()
}

func (p *PanicError) String() string {
	if len(p.Snippets) == 0 {
		return p.error.Error()
	}
	e := p.error.Error()
	for _, sg := range p.Snippets {
		f := sg.Path + ":" + strconv.Itoa(p.Traces[sg.Index].Line)
		if pos := strings.Index(e, f); pos > -1 {
			start := pos + len(f) + 1
			var s string
			if start < len(e) {
				s = e[start:]
			}
			e = e[0:pos+len(f)] + "\n" + sg.String() + "\n\n" + s
		}
	}
	return e
}

func (p *PanicError) HTML() template.HTML {
	if len(p.Snippets) == 0 {
		return template.HTML(`<pre>` + p.error.Error() + `</pre>`)
	}
	table := "<style>.panic-table-snippet td.left{width:100px;text-align:right}.panic-table-trace td.left{width:50%}</style>"
	for _, sg := range p.Snippets {
		table += `<table class="table table-bordered panic-table panic-table-snippet">`
		table += `<thead><tr><th colspan="2">` + sg.Path + `</th></tr></thead>`
		table += `<tbody>`
		table += sg.TableRow()
		table += `</tbody>`
		table += `</table>`
	}
	table += `<table class="table table-bordered panic-table panic-table-trace">`
	table += `<thead><tr><th colspan="2">Trace</th></tr><tr><th>File</th><th>Func</th></tr></thead>`
	table += `<tbody>`
	for _, ts := range p.Traces {
		f := html.EscapeString(ts.File) + `:` + strconv.Itoa(ts.Line)
		table += `<tr><td class='left'>`
		if ts.HasErr {
			table += `<strong>` + f + `</strong>`
			table += `</td><td class='right'><strong>` + ts.Func + `</strong>`
		} else {
			table += f + `</td><td class='right'>` + ts.Func
		}
		table += `</td></tr>`
	}
	table += `</tbody>`
	table += `</table>`
	return template.HTML(table)
}

func (p *PanicError) AddTrace(trace *Trace) *PanicError {
	if len(p.Snippets) == 0 {
		var index int
		if strings.Index(trace.File, workDir) != -1 {
			trace.HasErr = true
			index = len(p.Traces)
		}
		if trace.HasErr {
			p.ExtractSnippets(trace.File, trace.Line, index)
		}
	}
	p.Traces = append(p.Traces, trace)
	return p
}

func (p *PanicError) ExtractSnippets(file string, curLineNum int, index int) error {
	content, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}
	half := SnippetLineNumbers / 2
	lines := strings.Split(string(content), "\n")
	group := &SnippetGroup{
		Path:    file,
		Index:   index,
		Snippet: []*Snippet{},
	}
	for lineNum := curLineNum - half; lineNum <= curLineNum+half; lineNum++ {
		if len(lines) >= lineNum {
			group.Snippet = append(group.Snippet, &Snippet{
				Number:  lineNum,
				Code:    lines[lineNum-1],
				Current: lineNum == curLineNum,
			})
		}
	}
	if len(group.Snippet) > 0 {
		p.Snippets = append(p.Snippets, group)
	}
	return nil
}

func (p *PanicError) SetError(err error) *PanicError {
	p.error = err
	return p
}

func (p *PanicError) SetErrorString(errStr string) *PanicError {
	p.error = errors.New(errStr)
	return p
}

func (p *PanicError) SetDebug(on bool) *PanicError {
	p.debug = on
	return p
}

func (p *PanicError) Debug() bool {
	return p.debug
}
