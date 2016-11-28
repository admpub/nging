/*

   Copyright 2016 Wenhui Shen <www.webx.top>

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

/**
 * 模板扩展
 * @author swh <swh@admpub.com>
 */
package standard

import (
	"bytes"
	"errors"
	"fmt"
	htmlTpl "html/template"
	"io"
	"io/ioutil"
	"path/filepath"
	"regexp"
	"strings"
	"sync"

	"github.com/admpub/log"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/logger"
	. "github.com/webx-top/echo/middleware/render/driver"
	. "github.com/webx-top/echo/middleware/render/manager"
)

var Debug = false

func New(templateDir string, args ...logger.Logger) Driver {
	var err error
	templateDir, err = filepath.Abs(templateDir)
	if err != nil {
		panic(err.Error())
	}
	t := &Standard{
		CachedRelation:    make(map[string]*CcRel),
		TemplateDir:       templateDir,
		DelimLeft:         "{{",
		DelimRight:        "}}",
		IncludeTag:        "Include",
		ExtendTag:         "Extend",
		BlockTag:          "Block",
		SuperTag:          "Super",
		Ext:               ".html",
		Debug:             Debug,
		fileEvents:        make([]func(string), 0),
		contentProcessors: make([]func([]byte) []byte, 0),
	}
	if len(args) > 0 {
		t.logger = args[0]
	} else {
		t.logger = log.New("render-standard")
	}
	t.InitRegexp()
	return t
}

type tplInfo struct {
	Template *htmlTpl.Template
	Blocks   map[string]struct{}
}

func NewTplInfo(t *htmlTpl.Template) *tplInfo {
	return &tplInfo{
		Template: t,
		Blocks:   map[string]struct{}{},
	}
}

type CcRel struct {
	Rel map[string]uint8
	Tpl [2]*tplInfo //0是独立模板；1是子模板
}

type Standard struct {
	CachedRelation     map[string]*CcRel
	TemplateDir        string
	TemplateMgr        *Manager
	contentProcessors  []func([]byte) []byte
	DelimLeft          string
	DelimRight         string
	incTagRegex        *regexp.Regexp
	extTagRegex        *regexp.Regexp
	blkTagRegex        *regexp.Regexp
	cachedRegexIdent   string
	IncludeTag         string
	ExtendTag          string
	BlockTag           string
	SuperTag           string
	Ext                string
	TemplatePathParser func(string) string
	Debug              bool
	FuncMapFn          func() map[string]interface{}
	logger             logger.Logger
	fileEvents         []func(string)
	mutex              *sync.RWMutex
}

func (self *Standard) SetLogger(l logger.Logger) {
	self.logger = l
	if self.TemplateMgr != nil {
		self.TemplateMgr.Logger = self.logger
	}
}
func (self *Standard) Logger() logger.Logger {
	return self.logger
}

func (self *Standard) TmplDir() string {
	return self.TemplateDir
}

func (self *Standard) MonitorEvent(fn func(string)) {
	if fn == nil {
		return
	}
	self.fileEvents = append(self.fileEvents, fn)
}

func (self *Standard) SetContentProcessor(fn func([]byte) []byte) {
	if fn == nil {
		return
	}
	self.contentProcessors = append(self.contentProcessors, fn)
}

func (self *Standard) SetFuncMapFn(fn func() map[string]interface{}) {
	self.FuncMapFn = fn
}

func (self *Standard) deleteCachedRelation(name string) {
	if cs, ok := self.CachedRelation[name]; ok {
		_ = cs
		self.CachedRelation = make(map[string]*CcRel)
		self.logger.Info("remove cached template object")
		/*
			for key, _ := range cs.Rel {
				if key == name {
					continue
				}
				self.deleteCachedRelation(key)
			}
			self.Logger.Info("remove cached template object:", name)
			delete(self.CachedRelation, name)
		*/
	}
}

func (self *Standard) Init(cached ...bool) {
	self.TemplateMgr = new(Manager)
	self.mutex = &sync.RWMutex{}

	ln := len(cached)
	if ln < 1 || !cached[0] {
		return
	}
	reloadTemplates := true
	if ln > 1 {
		reloadTemplates = cached[1]
	}
	self.TemplateMgr.OnChangeCallback = func(name, typ, event string) {
		switch event {
		case "create":
		case "delete", "modify", "rename":
			if typ == "dir" {
				return
			}
			self.deleteCachedRelation(name)
			for _, fn := range self.fileEvents {
				fn(name)
			}
		}
	}
	self.TemplateMgr.Init(self.logger, self.TemplateDir, reloadTemplates, "*"+self.Ext)
}

func (self *Standard) SetMgr(mgr *Manager) {
	self.TemplateMgr = mgr
}

func (self *Standard) TemplatePath(p string) string {
	if self.TemplatePathParser == nil {
		return p
	}
	return self.TemplatePathParser(p)
}

func (self *Standard) echo(messages ...string) {
	if self.Debug {
		var message string
		for _, v := range messages {
			message += v + ` `
		}
		fmt.Println(`[tplex]`, message)
	}
}

func (self *Standard) InitRegexp() {
	left := regexp.QuoteMeta(self.DelimLeft)
	right := regexp.QuoteMeta(self.DelimRight)
	rfirst := regexp.QuoteMeta(self.DelimRight[0:1])
	self.incTagRegex = regexp.MustCompile(left + self.IncludeTag + `[\s]+"([^"]+)"(?:[\s]+([^` + rfirst + `]+))?[\s]*` + right)
	self.extTagRegex = regexp.MustCompile(left + self.ExtendTag + `[\s]+"([^"]+)"(?:[\s]+([^` + rfirst + `]+))?[\s]*` + right)
	self.blkTagRegex = regexp.MustCompile(`(?s)` + left + self.BlockTag + `[\s]+"([^"]+)"[\s]*` + right + `(.*?)` + left + `\/` + self.BlockTag + right)
}

// Render HTML
func (self *Standard) Render(w io.Writer, tmplName string, values interface{}, c echo.Context) error {
	var funcMap htmlTpl.FuncMap
	funcs := c.Funcs()
	if self.FuncMapFn != nil {
		funcMap = self.FuncMapFn()
		if funcs != nil {
			for k, v := range funcs {
				funcMap[k] = v
			}
		}
	} else {
		if funcs != nil {
			funcMap = funcs
		}
	}
	tmpl := self.parse(tmplName, funcMap)
	buf := new(bytes.Buffer)
	err := tmpl.ExecuteTemplate(buf, tmpl.Name(), values)
	if err != nil {
		return errors.New(fmt.Sprintf("Parse %v err: %v", tmpl.Name(), err))
	}
	_, err = io.Copy(w, buf)
	if err != nil {
		return errors.New(fmt.Sprintf("Parse %v err: %v", tmpl.Name(), err))
	}
	return err
}

func (self *Standard) parse(tmplName string, funcMap htmlTpl.FuncMap) (tmpl *htmlTpl.Template) {
	self.mutex.Lock()
	defer self.mutex.Unlock()
	tmplName = tmplName + self.Ext
	tmplName = self.TemplatePath(tmplName)
	cachedKey := tmplName
	if tmplName[0] == '/' {
		cachedKey = tmplName[1:]
	}
	rel, ok := self.CachedRelation[cachedKey]
	if ok && rel.Tpl[0].Template != nil {
		tmpl = rel.Tpl[0].Template
		funcMap = setFunc(rel.Tpl[0], funcMap)
		tmpl.Funcs(funcMap)
		if self.Debug {
			fmt.Println(`Using the template object to be cached:`, tmplName)
			fmt.Println("_________________________________________")
			fmt.Println("")
			for k, v := range tmpl.Templates() {
				fmt.Printf("%v. %#v\n", k, v.Name())
			}
			fmt.Println("_________________________________________")
			fmt.Println("")
		}
		return
	}
	t := htmlTpl.New(tmplName)
	t.Delims(self.DelimLeft, self.DelimRight)
	if rel == nil {
		rel = &CcRel{
			Rel: map[string]uint8{cachedKey: 0},
			Tpl: [2]*tplInfo{NewTplInfo(nil), NewTplInfo(nil)},
		}
	}
	funcMap = setFunc(rel.Tpl[0], funcMap)
	t.Funcs(funcMap)
	self.echo(`Read not cached template content:`, tmplName)
	b, err := self.RawContent(tmplName)
	if err != nil {
		tmpl, _ = t.Parse(err.Error())
		return
	}

	content := string(b)
	subcs := make(map[string]string, 0) //子模板内容
	extcs := make(map[string]string, 0) //母板内容

	ident := self.DelimLeft + self.IncludeTag + self.DelimRight
	if self.cachedRegexIdent != ident || self.incTagRegex == nil {
		self.InitRegexp()
	}
	m := self.extTagRegex.FindAllStringSubmatch(content, 1)
	if len(m) > 0 {
		self.ParseBlock(content, &subcs, &extcs)
		extFile := m[0][1] + self.Ext
		passObject := m[0][2]
		extFile = self.TemplatePath(extFile)
		self.echo(`Read layout template content:`, extFile)
		b, err = self.RawContent(extFile)
		if err != nil {
			tmpl, _ = t.Parse(err.Error())
			return
		}
		content = string(b)
		content = self.ParseExtend(content, &extcs, passObject, &subcs)

		if v, ok := self.CachedRelation[extFile]; !ok {
			self.CachedRelation[extFile] = &CcRel{
				Rel: map[string]uint8{cachedKey: 0},
				Tpl: [2]*tplInfo{NewTplInfo(nil), NewTplInfo(nil)},
			}
		} else if _, ok := v.Rel[cachedKey]; !ok {
			self.CachedRelation[extFile].Rel[cachedKey] = 0
		}
	}
	content = self.ContainsSubTpl(content, &subcs)
	self.echo(`The template content:`, content)
	tmpl, err = t.Parse(content)
	if err != nil {
		content = fmt.Sprintf("Parse %v err: %v", tmplName, err)
		tmpl, _ = t.Parse(content)
		return
	}
	for name, subc := range subcs {
		v, ok := self.CachedRelation[name]
		if ok && v.Tpl[1].Template != nil {
			self.CachedRelation[name].Rel[cachedKey] = 0
			tmpl.AddParseTree(name, self.CachedRelation[name].Tpl[1].Template.Tree)
			continue
		}
		var t *htmlTpl.Template
		if name == tmpl.Name() {
			t = tmpl
		} else {
			t = tmpl.New(name)
			_, err = t.Parse(subc)
			if err != nil {
				t.Parse(fmt.Sprintf("Parse File %v err: %v", name, err))
			}
		}

		if ok {
			self.CachedRelation[name].Rel[cachedKey] = 0
			self.CachedRelation[name].Tpl[1].Template = t
		} else {
			self.CachedRelation[name] = &CcRel{
				Rel: map[string]uint8{cachedKey: 0},
				Tpl: [2]*tplInfo{NewTplInfo(nil), NewTplInfo(t)},
			}
		}

	}
	for name, extc := range extcs {
		var t *htmlTpl.Template
		if name == tmpl.Name() {
			t = tmpl
		} else {
			t = tmpl.New(name)
			_, err = t.Parse(extc)
			if err != nil {
				t.Parse(fmt.Sprintf("Parse Block %v err: %v", name, err))
			}
		}
		rel.Tpl[0].Blocks[name] = struct{}{}
	}

	rel.Tpl[0].Template = tmpl
	self.CachedRelation[cachedKey] = rel
	return
}

func (self *Standard) Fetch(tmplName string, data interface{}, funcMap map[string]interface{}) string {
	return self.execute(self.parse(tmplName, funcMap), data)
}

func (self *Standard) execute(tmpl *htmlTpl.Template, data interface{}) string {
	buf := new(bytes.Buffer)
	err := tmpl.ExecuteTemplate(buf, tmpl.Name(), data)
	if err != nil {
		return fmt.Sprintf("Parse %v err: %v", tmpl.Name(), err)
	}
	b, err := ioutil.ReadAll(buf)
	if err != nil {
		return fmt.Sprintf("Parse %v err: %v", tmpl.Name(), err)
	}
	return string(b)
}

func (self *Standard) ParseBlock(content string, subcs *map[string]string, extcs *map[string]string) {
	matches := self.blkTagRegex.FindAllStringSubmatch(content, -1)
	for _, v := range matches {
		blockName := v[1]
		content := v[2]
		(*extcs)[blockName] = self.Tag(`define "`+blockName+`"`) + self.ContainsSubTpl(content, subcs) + self.Tag(`end`)
	}
}

func (self *Standard) ParseExtend(content string, extcs *map[string]string, passObject string, subcs *map[string]string) string {
	if passObject == "" {
		passObject = "."
	}
	matches := self.blkTagRegex.FindAllStringSubmatch(content, -1)
	var superTag string
	if self.SuperTag != "" {
		superTag = self.Tag(self.SuperTag)
	}
	rec := make(map[string]uint8)
	sup := make(map[string]string)
	for _, v := range matches {
		matched := v[0]
		blockName := v[1]
		innerStr := v[2]
		if v, ok := (*extcs)[blockName]; ok {
			var suffix string
			if idx, ok := rec[blockName]; ok {
				idx++
				rec[blockName] = idx
				suffix = fmt.Sprintf(`.%v`, idx)
			} else {
				rec[blockName] = 0
			}
			if superTag != "" {
				sv, hasSuper := sup[blockName]
				if !hasSuper {
					hasSuper = strings.Contains(v, superTag)
					if hasSuper {
						sup[blockName] = v
					}
				} else {
					v = sv
				}
				if hasSuper {
					innerStr = self.ContainsSubTpl(innerStr, subcs)
					v = strings.Replace(v, superTag, innerStr, 1)
					if suffix == `` {
						(*extcs)[blockName] = v
					}
				}
			}
			if suffix != `` {
				(*extcs)[blockName+suffix] = strings.Replace(v, self.Tag(`define "`+blockName+`"`), self.Tag(`define "`+blockName+suffix+`"`), 1)
				rec[blockName+suffix] = 0
			}
			content = strings.Replace(content, matched, self.Tag(`template "`+blockName+suffix+`" `+passObject), 1)
		} else {
			content = strings.Replace(content, matched, innerStr, 1)
		}
	}
	//只保留layout中存在的Block
	for k := range *extcs {
		if _, ok := rec[k]; !ok {
			delete(*extcs, k)
		}
	}
	return content
}

func (self *Standard) ContainsSubTpl(content string, subcs *map[string]string) string {
	matches := self.incTagRegex.FindAllStringSubmatch(content, -1)
	for _, v := range matches {
		matched := v[0]
		tmplFile := v[1]
		passObject := v[2]
		tmplFile += self.Ext
		tmplFile = self.TemplatePath(tmplFile)
		if _, ok := (*subcs)[tmplFile]; !ok {
			// if v, ok := self.CachedRelation[tmplFile]; ok && v.Tpl[1] != nil {
			// 	(*subcs)[tmplFile] = ""
			// } else {
			b, err := self.RawContent(tmplFile)
			if err != nil {
				return fmt.Sprintf("RenderTemplate %v read err: %s", tmplFile, err)
			}
			str := string(b)
			(*subcs)[tmplFile] = "" //先登记，避免死循环
			str = self.ContainsSubTpl(str, subcs)
			(*subcs)[tmplFile] = self.Tag(`define "`+tmplFile+`"`) + str + self.Tag(`end`)
			//}
		}
		if passObject == "" {
			passObject = "."
		}
		content = strings.Replace(content, matched, self.Tag(`template "`+tmplFile+`" `+passObject), -1)
	}
	return content
}

func (self *Standard) Tag(content string) string {
	return self.DelimLeft + content + self.DelimRight
}

func (self *Standard) Include(tmplName string, funcMap htmlTpl.FuncMap, values interface{}) interface{} {
	return htmlTpl.HTML(self.Fetch(tmplName, values, funcMap))
}

func (self *Standard) RawContent(tmpl string) (b []byte, e error) {
	defer func() {
		if b != nil && self.contentProcessors != nil {
			for _, fn := range self.contentProcessors {
				b = fn(b)
			}
		}
	}()
	if self.TemplateMgr != nil && self.TemplateMgr.Caches != nil {
		b, e = self.TemplateMgr.GetTemplate(tmpl)
		if e != nil {
			self.logger.Error(e)
		}
		return
	}
	return ioutil.ReadFile(filepath.Join(self.TemplateDir, tmpl))
}

func (self *Standard) ClearCache() {
	if self.TemplateMgr != nil {
		self.TemplateMgr.ClearCache()
	}
	self.CachedRelation = make(map[string]*CcRel)
}

func (self *Standard) Close() {
	self.ClearCache()
	if self.TemplateMgr != nil {
		self.TemplateMgr.Close()
	}
}

func setFunc(tplInf *tplInfo, funcMap htmlTpl.FuncMap) htmlTpl.FuncMap {
	if funcMap == nil {
		funcMap = htmlTpl.FuncMap{}
	}
	funcMap["hasBlock"] = func(blocks ...string) bool {
		for _, blockName := range blocks {
			if _, ok := tplInf.Blocks[blockName]; !ok {
				return false
			}
		}
		return true
	}
	funcMap["hasAnyBlock"] = func(blocks ...string) bool {
		for _, blockName := range blocks {
			if _, ok := tplInf.Blocks[blockName]; ok {
				return true
			}
		}
		return false
	}
	return funcMap
}
