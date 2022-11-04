/*

   Copyright 2017-present Wenhui Shen <www.webx.top>

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
	"runtime"
	"strconv"
	"strings"

	pkgCode "github.com/webx-top/echo/code"
)

const (
	SnippetLineNumbers = 13
	StackSize          = 4 << 10 // 4 KB
)

// ==========================================
// Error
// ==========================================

func NewError(msg string, code ...pkgCode.Code) *Error {
	e := &Error{Code: pkgCode.Failure, Message: msg, Extra: H{}}
	if len(code) > 0 {
		e.Code = code[0]
	}
	if len(msg) == 0 {
		e.Message = e.Code.String()
	}
	return e
}

func NewErrorWith(err error, msg string, code ...pkgCode.Code) *Error {
	e := &Error{Code: pkgCode.Failure, Message: msg, Extra: H{}, cause: err}
	if len(code) > 0 {
		e.Code = code[0]
	}
	if len(msg) == 0 {
		if err != nil {
			e.Message = err.Error()
		} else {
			e.Message = e.Code.String()
		}
	}
	return e
}

func IsErrorCode(err error, code pkgCode.Code) bool {
	val, ok := err.(*Error)
	if !ok {
		return false
	}
	return val.Code.Is(code)
}

func InErrorCode(err error, codes ...pkgCode.Code) bool {
	val, ok := err.(*Error)
	if !ok {
		return false
	}
	return val.Code.In(codes...)
}

type Error struct {
	Code    pkgCode.Code
	Message string
	Zone    string
	Extra   H
	cause   error
}

// Error returns message.
func (e *Error) Error() string {
	return e.Message
}

func (e *Error) Set(key string, value interface{}) *Error {
	e.Extra.Set(key, value)
	return e
}

func (e *Error) SetMessage(message string) *Error {
	e.Message = message
	return e
}

func (e *Error) SetZone(zone string) *Error {
	e.Zone = zone
	return e
}

func (e *Error) SetError(err error) *Error {
	e.cause = err
	return e
}

func (e *Error) Delete(keys ...string) *Error {
	e.Extra.Delete(keys...)
	return e
}

func (e *Error) ErrorCode() int {
	return e.Code.Int()
}

func (e *Error) Cause() error {
	return e.cause
}

func (e *Error) Unwrap() error {
	return e.cause
}

// ==========================================
// HTTPError
// ==========================================

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
	raw     error
}

// Error returns message.
func (e *HTTPError) Error() string {
	if e.raw != nil {
		return e.Message + `: ` + e.raw.Error()
	}
	return e.Message
}

// SetRaw sets the raw error
func (e *HTTPError) SetRaw(err error) *HTTPError {
	e.raw = err
	return e
}

// Raw gets the raw error
func (e *HTTPError) Raw() error {
	return e.raw
}

func (e *HTTPError) Unwrap() error {
	return e.raw
}

// ==========================================
// PanicError
// ==========================================

func NewPanicError(recovered interface{}, err error, debugAndDisableStackAll ...bool) *PanicError {
	var debug, disableStackAll bool
	switch len(debugAndDisableStackAll) {
	case 2:
		disableStackAll = debugAndDisableStackAll[1]
		fallthrough
	case 1:
		debug = debugAndDisableStackAll[0]
	}

	return &PanicError{
		error:           err,
		Raw:             recovered,
		Traces:          make([]*Trace, 0),
		Snippets:        make([]*SnippetGroup, 0),
		debug:           debug,
		disableStackAll: disableStackAll,
	}
}

type PanicError struct {
	error
	Raw             interface{}
	Traces          []*Trace
	Snippets        []*SnippetGroup
	debug           bool
	disableStackAll bool
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

func (p *PanicError) AddTrace(trace *Trace, content ...string) *PanicError {
	if len(p.Snippets) == 0 {
		if !trace.HasErr && strings.Contains(trace.File, workDir) {
			trace.HasErr = true
		}
		if trace.HasErr {
			index := len(p.Traces)
			if len(content) > 0 {
				p.ExtractSnippets(content[0], trace.File, trace.Line, index)
			} else {
				p.ExtractFileSnippets(trace.File, trace.Line, index)
			}
		}
	}
	p.Traces = append(p.Traces, trace)
	return p
}

func (p *PanicError) ExtractFileSnippets(file string, curLineNum int, index int) error {
	content, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}
	return p.ExtractSnippets(string(content), file, curLineNum, index)
}

func (p *PanicError) ExtractSnippets(content string, file string, curLineNum int, index int) error {
	half := SnippetLineNumbers / 2
	lines := strings.Split(string(content), "\n")
	group := &SnippetGroup{
		Path:    file,
		Index:   index,
		Snippet: []*Snippet{},
	}
	for lineNum := curLineNum - half; lineNum <= curLineNum+half; lineNum++ {
		if len(lines) >= lineNum && lineNum > 0 {
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

func (p *PanicError) Parse(stackSizes ...int) *PanicError {
	stackSize := StackSize
	if len(stackSizes) > 0 {
		stackSize = stackSizes[0]
	}
	var err error
	switch r := p.Raw.(type) {
	case error:
		err = r
	default:
		err = fmt.Errorf("%v", r)
	}
	if p.disableStackAll {
		p.SetError(err)
		return p
	}
	content := "[PANIC RECOVER] " + err.Error()
	for i := 2; len(content) < stackSize; i++ {
		pc, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}
		t := &Trace{
			File: file,
			Line: line,
			Func: runtime.FuncForPC(pc).Name(),
		}
		p.AddTrace(t)
		content += "\n" + fmt.Sprintf(`%v:%v`, file, line)
	}
	p.SetErrorString(content)
	return p
}
