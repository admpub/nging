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
	"fmt"
	htmlTpl "html/template"
	"io"
	"io/ioutil"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/admpub/log"
	"golang.org/x/sync/singleflight"

	"github.com/webx-top/com"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/logger"
	"github.com/webx-top/echo/middleware/render/driver"
	"github.com/webx-top/echo/middleware/render/manager"
	"github.com/webx-top/poolx/bufferpool"
)

var Debug = false

func New(templateDir string, args ...logger.Logger) driver.Driver {
	var err error
	templateDir, err = filepath.Abs(templateDir)
	if err != nil {
		panic(err.Error())
	}
	t := &Standard{
		CachedRelation:    NewCache(),
		TemplateDir:       templateDir,
		DelimLeft:         "{{",
		DelimRight:        "}}",
		IncludeTag:        "Include",
		FunctionTag:       "Function",
		ExtendTag:         "Extend",
		BlockTag:          "Block",
		SuperTag:          "Super",
		StripTag:          "Strip",
		Ext:               ".html",
		debug:             Debug,
		fileEvents:        make([]func(string), 0),
		contentProcessors: make([]func([]byte) []byte, 0),
	}
	if len(args) > 0 {
		t.logger = args[0]
	} else {
		t.logger = log.New("render-standard")
	}
	t.InitRegexp()
	t.SetManager(manager.Default)
	return t
}

type Standard struct {
	CachedRelation     *CacheData
	TemplateDir        string
	TemplateMgr        driver.Manager
	contentProcessors  []func([]byte) []byte
	DelimLeft          string
	DelimRight         string
	incTagRegex        *regexp.Regexp
	funcTagRegex       *regexp.Regexp
	extTagRegex        *regexp.Regexp
	blkTagRegex        *regexp.Regexp
	rplTagRegex        *regexp.Regexp
	innerTagBlankRegex *regexp.Regexp
	stripTagRegex      *regexp.Regexp
	IncludeTag         string
	FunctionTag        string
	ExtendTag          string
	BlockTag           string
	SuperTag           string
	StripTag           string
	Ext                string
	tmplPathFixer      func(echo.Context, string) string
	debug              bool
	getFuncs           func() map[string]interface{}
	logger             logger.Logger
	fileEvents         []func(string)
	mutex              sync.RWMutex
	quotedLeft         string
	quotedRight        string
	quotedRfirst       string
	sg                 singleflight.Group
}

func (a *Standard) Debug() bool {
	return a.debug
}

func (a *Standard) SetDebug(on bool) {
	a.debug = on
}

func (a *Standard) SetLogger(l logger.Logger) {
	a.logger = l
	if a.TemplateMgr != nil {
		a.TemplateMgr.SetLogger(a.logger)
	}
}
func (a *Standard) Logger() logger.Logger {
	return a.logger
}

func (a *Standard) TmplDir() string {
	return a.TemplateDir
}

func (a *Standard) MonitorEvent(fn func(string)) {
	if fn == nil {
		return
	}
	a.fileEvents = append(a.fileEvents, fn)
}

func (a *Standard) SetContentProcessor(fn func([]byte) []byte) {
	if fn == nil {
		return
	}
	a.contentProcessors = append(a.contentProcessors, fn)
}

func (a *Standard) SetFuncMap(fn func() map[string]interface{}) {
	a.getFuncs = fn
}

func (a *Standard) deleteCachedRelation(name string) {
	if cs, ok := a.CachedRelation.GetOk(name); ok {
		_ = cs
		a.CachedRelation.Reset()
		a.logger.Info("remove cached template object")
		/*
			for key, _ := range cs.Rel {
				if key == name {
					continue
				}
				a.deleteCachedRelation(key)
			}
			a.Logger.Info("remove cached template object:", name)
			delete(a.CachedRelation, name)
		*/
	}
}

func (a *Standard) Init() {
	a.InitRegexp()
	callback := func(name, typ, event string) {
		switch event {
		case "create":
		case "delete", "modify", "rename":
			if typ == "dir" {
				return
			}
			a.deleteCachedRelation(name)
			for _, fn := range a.fileEvents {
				fn(name)
			}
		}
	}
	a.TemplateMgr.AddAllow("*" + a.Ext)
	a.TemplateMgr.AddWatchDir(a.TemplateDir)
	a.TemplateMgr.AddCallback(a.TemplateDir, callback)
	a.TemplateMgr.Start()
}

func (a *Standard) SetManager(mgr driver.Manager) {
	if a.TemplateMgr != nil {
		a.TemplateMgr.Close()
	}
	a.TemplateMgr = mgr
}

func (a *Standard) Manager() driver.Manager {
	return a.TemplateMgr
}

func (a *Standard) SetTmplPathFixer(fn func(echo.Context, string) string) {
	a.tmplPathFixer = fn
}

func (a *Standard) TmplPath(c echo.Context, p string) string {
	if a.tmplPathFixer != nil {
		return a.tmplPathFixer(c, p)
	}
	p = filepath.Join(a.TemplateDir, p)
	return p
}

func (a *Standard) InitRegexp() {
	a.quotedLeft = regexp.QuoteMeta(a.DelimLeft)
	a.quotedRight = regexp.QuoteMeta(a.DelimRight)
	a.quotedRfirst = regexp.QuoteMeta(a.DelimRight[0:1])

	//{{Include "tmpl"}} or {{Include "tmpl" .}}
	a.incTagRegex = regexp.MustCompile(a.quotedLeft + a.IncludeTag + `[\s]+"([^"]+)"(?:[\s]+([^` + a.quotedRfirst + `]+))?[\s]*\/?` + a.quotedRight)

	//{{Function "funcName"}} or {{Function "funcName" .}}
	a.funcTagRegex = regexp.MustCompile(a.quotedLeft + a.FunctionTag + `[\s]+"([^"]+)"(?:[\s]+([^` + a.quotedRfirst + `]+))?[\s]*\/?` + a.quotedRight)

	//{{Extend "name"}}
	a.extTagRegex = regexp.MustCompile(`^[\s]*` + a.quotedLeft + a.ExtendTag + `[\s]+"([^"]+)"(?:[\s]+([^` + a.quotedRfirst + `]+))?[\s]*\/?` + a.quotedRight)

	//{{Block "name"}}content{{/Block}}
	a.blkTagRegex = regexp.MustCompile(`(?s)` + a.quotedLeft + a.BlockTag + `[\s]+"([^"]+)"[\s]*` + a.quotedRight + `(.*?)` + a.quotedLeft + `\/` + a.BlockTag + a.quotedRight)

	//{{Block "name"/}}
	a.rplTagRegex = regexp.MustCompile(a.quotedLeft + a.BlockTag + `[\s]+"([^"]+)"[\s]*\/` + a.quotedRight)

	//}}...{{ or >...<
	a.innerTagBlankRegex = regexp.MustCompile(`(?s)(` + a.quotedRight + `|>)[\s]{2,}(` + a.quotedLeft + `|<)`)

	//{{Strip}}...{{/Strip}}
	a.stripTagRegex = regexp.MustCompile(`(?s)` + a.quotedLeft + a.StripTag + a.quotedRight + `(.*?)` + a.quotedLeft + `\/` + a.StripTag + a.quotedRight)
}

// Render HTML
func (a *Standard) Render(w io.Writer, tmplName string, values interface{}, c echo.Context) error {
	// if c.Get(`webx:render.locked`) == nil {
	// 	c.Set(`webx:render.locked`, true)
	// 	a.mutex.Lock()
	// 	defer func() {
	// 		a.mutex.Unlock()
	// 		c.Delete(`webx:render.locked`)
	// 	}()
	// }
	tmpl, err := a.parse(c, tmplName)
	if err != nil {
		return err
	}
	return tmpl.ExecuteTemplate(w, tmpl.Name(), values)
}

func (a *Standard) parse(c echo.Context, tmplName string) (tmpl *htmlTpl.Template, err error) {
	funcs := c.Funcs()
	tmplOriginalName := tmplName
	tmplName = tmplName + a.Ext
	tmplName = a.TmplPath(c, tmplName)
	cachedKey := tmplName
	var funcMap htmlTpl.FuncMap
	if a.getFuncs != nil {
		funcMap = htmlTpl.FuncMap(a.getFuncs())
	}
	if funcMap == nil {
		funcMap = htmlTpl.FuncMap{}
	}
	for k, v := range funcs {
		funcMap[k] = v
	}
	rel, ok := a.CachedRelation.GetOk(cachedKey)
	if ok && rel.Tpl[0].Template != nil {
		tmpl = rel.Tpl[0].Template
		if a.debug {
			a.logger.Debug(` `+tmplName, tmpl.DefinedTemplates())
		}
		funcMap = setFunc(rel.Tpl[0], funcMap)
		tmpl.Funcs(funcMap)
		return
	}
	var v interface{}
	var shared bool
	v, err, shared = a.sg.Do(cachedKey, func() (interface{}, error) {
		return a.find(c, rel, tmplOriginalName, tmplName, cachedKey, funcMap)
	})
	if err != nil {
		return
	}
	if !shared {
		tmpl = v.(*htmlTpl.Template)
		return
	}
	rel, ok = a.CachedRelation.GetOk(cachedKey)
	if ok && rel.Tpl[0].Template != nil {
		tmpl = rel.Tpl[0].Template
		funcMap = setFunc(rel.Tpl[0], funcMap)
		tmpl.Funcs(funcMap)
		return
	}
	return
	//return a.find(c, rel, tmplOriginalName, tmplName, cachedKey, funcMap)
}

var bytesBOM = []byte("\xEF\xBB\xBF")

func (a *Standard) find(c echo.Context, rel *CcRel, tmplOriginalName string, tmplName string, cachedKey string, funcMap htmlTpl.FuncMap) (tmpl *htmlTpl.Template, err error) {
	if a.debug {
		start := time.Now()
		a.logger.Debug(` ◐ compile template: `, tmplName)
		defer func() {
			a.logger.Debug(` ◑ finished compile: `+tmplName, ` (elapsed: `+time.Since(start).String()+`)`)
		}()
	}
	t := htmlTpl.New(driver.CleanTemplateName(tmplName))
	t.Delims(a.DelimLeft, a.DelimRight)
	if rel == nil {
		rel = &CcRel{
			Rel: map[string]uint8{cachedKey: 0},
			Tpl: [2]*tplInfo{NewTplInfo(nil), NewTplInfo(nil)},
		}
	}
	funcMap = setFunc(rel.Tpl[0], funcMap)
	t.Funcs(funcMap)
	b, err := a.RawContent(tmplName)
	if err != nil {
		tmpl, _ = t.Parse(err.Error())
		return
	}
	content := string(b)
	subcs := map[string]string{} //子模板内容
	extcs := map[string]string{} //母板内容
	m := a.extTagRegex.FindAllStringSubmatch(content, 1)
	content = a.rplTagRegex.ReplaceAllString(content, ``)
	for i := 0; i < 10 && len(m) > 0; i++ {
		a.ParseBlock(c, content, subcs, extcs)
		extFile := m[0][1] + a.Ext
		passObject := m[0][2]
		extFile = a.TmplPath(c, extFile)
		b, err = a.RawContent(extFile)
		if err != nil {
			tmpl, _ = t.Parse(err.Error())
			return
		}
		content = string(b)
		content, m = a.ParseExtend(c, content, extcs, passObject, subcs)

		if v, ok := a.CachedRelation.GetOk(extFile); !ok {
			a.CachedRelation.Set(extFile, NewRel(cachedKey))
		} else if _, ok := v.GetOk(cachedKey); !ok {
			v.Set(cachedKey, 0)
		}
	}
	content = a.ContainsSubTpl(c, content, subcs)
	clips := map[string]string{}
	content = a.ContainsFunctionResult(c, tmplOriginalName, content, clips)
	tmpl, err = t.Parse(content)
	if err != nil {
		content = fmt.Sprintf("Parse %v err: %v", tmplName, err)
		tmpl, _ = t.Parse(content)
		return
	}
	for name, subc := range subcs {
		v, ok := a.CachedRelation.GetOk(name)
		if ok && v.Tpl[1].Template != nil {
			v.Set(cachedKey, 0)
			tmpl.AddParseTree(name, v.Tpl[1].Template.Tree)
			continue
		}
		subc = a.ContainsFunctionResult(c, tmplOriginalName, subc, clips)
		t = tmpl.New(name)
		subc = a.Tag(`define "`+driver.CleanTemplateName(name)+`"`) + subc + a.Tag(`end`)
		_, err = t.Parse(subc)
		if err != nil {
			t.Parse(fmt.Sprintf("Parse File %v err: %v", name, err))
			return
		}

		if ok {
			v.Set(cachedKey, 0)
			v.Tpl[1].Template = t
		} else {
			a.CachedRelation.Set(name, &CcRel{
				Rel: map[string]uint8{cachedKey: 0},
				Tpl: [2]*tplInfo{NewTplInfo(nil), NewTplInfo(t)},
			})
		}

	}

	for name, extc := range extcs {
		t = tmpl.New(name)
		extc = a.ContainsFunctionResult(c, tmplOriginalName, extc, clips)
		extc = a.Tag(`define "`+driver.CleanTemplateName(name)+`"`) + extc + a.Tag(`end`)
		_, err = t.Parse(extc)
		if err != nil {
			t.Parse(fmt.Sprintf("Parse Block %v err: %v", name, err))
			return
		}
		rel.Tpl[0].Blocks[name] = struct{}{}
	}

	rel.Tpl[0].Template = tmpl
	a.CachedRelation.Set(cachedKey, rel)
	return
}

func (a *Standard) Fetch(tmplName string, data interface{}, c echo.Context) string {
	content, _ := a.parse(c, tmplName)
	return a.execute(content, data)
}

func (a *Standard) execute(tmpl *htmlTpl.Template, data interface{}) string {
	buf := bufferpool.Get()
	defer bufferpool.Release(buf)
	err := tmpl.ExecuteTemplate(buf, tmpl.Name(), data)
	if err != nil {
		return fmt.Sprintf("Parse %v err: %v", tmpl.Name(), err)
	}
	return com.Bytes2str(buf.Bytes())
}

func (a *Standard) ParseBlock(c echo.Context, content string, subcs map[string]string, extcs map[string]string) {
	matches := a.blkTagRegex.FindAllStringSubmatch(content, -1)
	for _, v := range matches {
		blockName := v[1]
		content := v[2]
		extcs[blockName] = a.ContainsSubTpl(c, content, subcs)
	}
}

func (a *Standard) ParseExtend(c echo.Context, content string, extcs map[string]string, passObject string, subcs map[string]string) (string, [][]string) {
	m := a.extTagRegex.FindAllStringSubmatch(content, 1)
	hasParent := len(m) > 0
	if len(passObject) == 0 {
		passObject = "."
	}
	content = a.rplTagRegex.ReplaceAllStringFunc(content, func(match string) string {
		match = match[strings.Index(match, `"`)+1:]
		match = match[0:strings.Index(match, `"`)]
		if v, ok := extcs[match]; ok {
			return v
		}
		return ``
	})
	matches := a.blkTagRegex.FindAllStringSubmatch(content, -1)
	var superTag string
	if len(a.SuperTag) > 0 {
		superTag = a.Tag(a.SuperTag)
	}
	rec := make(map[string]uint8)
	sup := make(map[string]string)
	for _, v := range matches {
		matched := v[0]
		blockName := v[1]
		innerStr := v[2]
		if v, ok := extcs[blockName]; ok {
			var suffix string
			if idx, ok := rec[blockName]; ok {
				idx++
				rec[blockName] = idx
				suffix = fmt.Sprintf(`.%v`, idx)
			} else {
				rec[blockName] = 0
			}
			if len(superTag) > 0 {
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
					innerStr = a.ContainsSubTpl(c, innerStr, subcs)
					v = strings.Replace(v, superTag, innerStr, 1)
					if suffix == `` {
						extcs[blockName] = v
					}
				}
			}
			if len(suffix) > 0 {
				extcs[blockName+suffix] = v
				rec[blockName+suffix] = 0
			}
			if hasParent {
				content = strings.Replace(content, matched, a.DelimLeft+a.BlockTag+` "`+blockName+`"`+a.DelimRight+v+a.DelimLeft+`/`+a.BlockTag+a.DelimRight, 1)
			} else {
				content = strings.Replace(content, matched, a.Tag(`template "`+blockName+suffix+`" `+passObject), 1)
			}
		} else {
			if !hasParent {
				content = strings.Replace(content, matched, innerStr, 1)
			}
		}
	}
	//只保留layout中存在的Block
	for k := range extcs {
		if _, ok := rec[k]; !ok {
			delete(extcs, k)
		}
	}
	return content, m
}

func (a *Standard) ContainsSubTpl(c echo.Context, content string, subcs map[string]string) string {
	matches := a.incTagRegex.FindAllStringSubmatch(content, -1)
	for _, v := range matches {
		matched := v[0]
		tmplFile := v[1]
		passObject := v[2]
		tmplFile += a.Ext
		tmplFile = a.TmplPath(c, tmplFile)
		if _, ok := subcs[tmplFile]; !ok {
			// if v, ok := a.CachedRelation[tmplFile]; ok && v.Tpl[1] != nil {
			// 	subcs[tmplFile] = ""
			// } else {
			b, err := a.RawContent(tmplFile)
			if err != nil {
				return fmt.Sprintf("RenderTemplate %v read err: %s", tmplFile, err)
			}
			str := string(b)
			subcs[tmplFile] = "" //先登记，避免死循环
			str = a.ContainsSubTpl(c, str, subcs)
			subcs[tmplFile] = str
			//}
		}
		if len(passObject) == 0 {
			passObject = "."
		}
		content = strings.Replace(content, matched, a.Tag(`template "`+driver.CleanTemplateName(tmplFile)+`" `+passObject), -1)
	}
	return content
}

func (a *Standard) ContainsFunctionResult(c echo.Context, tmplOriginalName string, content string, clips map[string]string) string {
	matches := a.funcTagRegex.FindAllStringSubmatch(content, -1)
	for _, v := range matches {
		matched := v[0]
		funcName := v[1]
		passArg := v[2]
		key := funcName + `:` + passArg
		if _, ok := clips[key]; !ok {
			if fn, ok := c.GetFunc(funcName).(func(string, string) string); ok {
				clips[key] = fn(tmplOriginalName, passArg)
			} else {
				clips[key] = ``
			}
		}

		content = strings.Replace(content, matched, clips[key], -1)
	}
	return content
}

func (a *Standard) Tag(content string) string {
	return a.DelimLeft + content + a.DelimRight
}

func (a *Standard) preprocess(b []byte) []byte {
	if b == nil {
		return nil
	}
	if a.contentProcessors != nil {
		for _, fn := range a.contentProcessors {
			b = fn(b)
		}
	}
	return a.strip(b)
}

func (a *Standard) RawContent(tmpl string) (b []byte, e error) {
	if a.TemplateMgr != nil {
		b, e = a.TemplateMgr.GetTemplate(tmpl)
	} else {
		b, e = ioutil.ReadFile(tmpl)
	}
	if e != nil {
		return
	}
	b = bytes.TrimPrefix(b, bytesBOM)
	b = a.preprocess(b)
	return
}

func (a *Standard) strip(src []byte) []byte {
	if a.debug {
		src = bytes.ReplaceAll(src, []byte(a.DelimLeft+a.StripTag+a.DelimRight), []byte{})
		return bytes.ReplaceAll(src, []byte(a.DelimLeft+`/`+a.StripTag+a.DelimRight), []byte{})
	}
	src = a.stripTagRegex.ReplaceAllFunc(src, func(b []byte) []byte {
		b = bytes.TrimPrefix(b, []byte(a.DelimLeft+a.StripTag+a.DelimRight))
		b = bytes.TrimSuffix(b, []byte(a.DelimLeft+`/`+a.StripTag+a.DelimRight))
		var pres [][]byte
		b, pres = driver.ReplacePRE(b)
		b = a.innerTagBlankRegex.ReplaceAll(b, driver.FE)
		b = driver.RemoveMultiCRLF(b)
		b = bytes.TrimSpace(b)
		b = driver.RecoveryPRE(b, pres)
		return b
	})
	return src
}

func (a *Standard) stripSpace(b []byte) []byte {
	var pres [][]byte
	b, pres = driver.ReplacePRE(b)
	b = a.innerTagBlankRegex.ReplaceAll(b, driver.FE)
	b = bytes.TrimSpace(b)
	b = driver.RecoveryPRE(b, pres)
	return b
}

func (a *Standard) ClearCache() {
	if a.TemplateMgr != nil {
		a.TemplateMgr.ClearCache()
	}
	a.CachedRelation.Reset()
}

func (a *Standard) Close() {
	a.ClearCache()
	if a.TemplateMgr != nil {
		if a.TemplateMgr == manager.Default {
			a.TemplateMgr.CancelWatchDir(a.TemplateDir)
			a.TemplateMgr.DelCallback(a.TemplateDir)
		} else {
			a.TemplateMgr.Close()
		}
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
