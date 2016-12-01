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
package pongo2

import (
	"bytes"
	"io"
	"io/ioutil"
	"path/filepath"
	"strings"
	"sync"

	"github.com/admpub/log"
	. "github.com/admpub/pongo2"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/logger"
	"github.com/webx-top/echo/middleware/render"
	. "github.com/webx-top/echo/middleware/render/driver"
	. "github.com/webx-top/echo/middleware/render/manager"
)

func init() {
	render.Reg(`pongo2`, func(tmplDir string) Driver {
		return New(tmplDir)
	})
}

func New(templateDir string, args ...logger.Logger) Driver {
	var err error
	templateDir, err = filepath.Abs(templateDir)
	if err != nil {
		panic(err.Error())
	}
	a := &Pongo2{
		templateDir:       templateDir,
		ext:               `.html`,
		fileEvents:        make([]func(string), 0),
		contentProcessors: make([]func([]byte) []byte, 0),
	}
	if len(args) > 0 {
		a.logger = args[0]
	} else {
		a.logger = log.New("render-pongo2")
	}
	a.templateDir, _ = filepath.Abs(templateDir)
	return a
}

type Pongo2 struct {
	templates         map[string]*Template
	mutex             *sync.RWMutex
	loader            *templateLoader
	set               *TemplateSet
	ext               string
	templateDir       string
	Mgr               *Manager
	logger            logger.Logger
	getFuncs          func() map[string]interface{}
	fileEvents        []func(string)
	contentProcessors []func([]byte) []byte
}

type templateLoader struct {
	templateDir string
	ext         string
	logger      logger.Logger
	template    *Pongo2
}

func (a *templateLoader) Abs(base, name string) string {
	//a.logger.Info(base+" => %v\n", name)
	return filepath.Join(``, name)
}

// Get returns an io.Reader where the template's content can be read from.
func (a *templateLoader) Get(tmpl string) (io.Reader, error) {
	var b []byte
	var e error
	tmpl += a.ext
	tmpl = strings.TrimPrefix(tmpl, a.templateDir)
	b, e = a.template.RawContent(tmpl)
	if e != nil {
		a.logger.Error(e)
	}
	buf := new(bytes.Buffer)
	buf.WriteString(string(b))
	return buf, e
}

func (self *Pongo2) SetLogger(l logger.Logger) {
	self.logger = l
	self.loader.logger = l
	if self.Mgr != nil {
		self.Mgr.Logger = self.logger
	}
}
func (self *Pongo2) Logger() logger.Logger {
	return self.logger
}

func (self *Pongo2) TmplDir() string {
	return self.templateDir
}

func (a *Pongo2) MonitorEvent(fn func(string)) {
	if fn == nil {
		return
	}
	a.fileEvents = append(a.fileEvents, fn)
}

func (a *Pongo2) Init(cached ...bool) {
	a.Mgr = new(Manager)
	a.templates = map[string]*Template{}
	a.mutex = &sync.RWMutex{}
	loader := &templateLoader{
		templateDir: a.templateDir,
		ext:         a.ext,
		logger:      a.logger,
		template:    a,
	}
	a.loader = loader
	a.set = NewSet(a.templateDir, a.loader)

	ln := len(cached)
	if ln < 1 || !cached[0] {
		return
	}
	reloadTemplates := true
	if ln > 1 {
		reloadTemplates = cached[1]
	}

	a.Mgr.OnChangeCallback = func(name, typ, event string) {
		switch event {
		case "create":
		case "delete", "modify", "rename":
			if typ == "dir" || !strings.HasSuffix(name, a.ext) {
				return
			}
			key := strings.TrimSuffix(name, a.ext)
			//布局模板被修改时，清空缓存
			if strings.HasSuffix(key, `layout`) {
				a.templates = make(map[string]*Template)
				a.logger.Info(`remove all cached template object:`, name)
			} else if _, ok := a.templates[key]; ok {
				delete(a.templates, key)
				a.logger.Info(`remove cached template object:`, name)
			}
			for _, fn := range a.fileEvents {
				fn(name)
			}
		}
	}
	a.Mgr.Init(a.logger, a.templateDir, reloadTemplates, "*"+a.ext)
}

func (a *Pongo2) SetContentProcessor(fn func([]byte) []byte) {
	if fn == nil {
		return
	}
	a.contentProcessors = append(a.contentProcessors, fn)
}

func (a *Pongo2) SetFuncMapFn(fn func() map[string]interface{}) {
	a.getFuncs = fn
}

func (a *Pongo2) Render(w io.Writer, tmpl string, data interface{}, c echo.Context) error {
	t, context := a.parse(tmpl, data, c.Funcs())
	return t.ExecuteWriter(context, w)
}

func (a *Pongo2) parse(tmpl string, data interface{}, funcMap map[string]interface{}) (*Template, Context) {
	a.mutex.Lock()
	defer a.mutex.Unlock()
	k := tmpl
	if tmpl[0] == '/' {
		k = tmpl[1:]
	}
	t, ok := a.templates[k]
	if !ok {
		var err error
		t, err = a.set.FromFile(tmpl)
		if err != nil {
			t = Must(a.set.FromString(err.Error()))
			return t, Context{}
		}
		a.templates[k] = t
	}
	var context Context
	if a.getFuncs != nil {
		context = Context(a.getFuncs())
	}
	if v, ok := data.(Context); ok {
		if context == nil {
			context = v
		} else {
			for n, f := range v {
				context[n] = f
			}
		}
	} else if v, ok := data.(map[string]interface{}); ok {
		if context == nil {
			context = v
		} else {
			for n, f := range v {
				context[n] = f
			}
		}
	} else {
		if context == nil {
			context = Context{
				`value`: data,
			}
		} else {
			context[`value`] = data
		}
	}
	if funcMap != nil {
		for name, function := range funcMap {
			context[name] = function
			//a.Logger.Info("added func: %v => %#v", name, function)
		}
	}
	return t, context
}

func (a *Pongo2) Fetch(tmpl string, data interface{}, funcMap map[string]interface{}) string {
	t, context := a.parse(tmpl, data, funcMap)
	r, err := t.Execute(context)
	if err != nil {
		r = err.Error()
	}
	return r
}

func (a *Pongo2) RawContent(tmpl string) (b []byte, e error) {
	defer func() {
		if b != nil && a.contentProcessors != nil {
			for _, fn := range a.contentProcessors {
				b = fn(b)
			}
		}
	}()
	if a.Mgr != nil && a.Mgr.Caches != nil {
		b, e = a.Mgr.GetTemplate(tmpl)
	}
	if b == nil || e != nil {
		b, e = ioutil.ReadFile(filepath.Join(a.templateDir, tmpl))
	}
	return
}

func (a *Pongo2) ClearCache() {
	if a.Mgr != nil {
		a.Mgr.ClearCache()
	}
	a.templates = make(map[string]*Template)
}

func (a *Pongo2) Close() {
	a.ClearCache()
	if a.Mgr != nil {
		a.Mgr.Close()
	}
}
