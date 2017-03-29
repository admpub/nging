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
package mvc

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/admpub/confl"
	codec "github.com/gorilla/securecookie"
	"github.com/webx-top/com"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/engine"
	"github.com/webx-top/echo/engine/fasthttp"
	"github.com/webx-top/echo/engine/standard"
	"github.com/webx-top/echo/handler/mvc/events"
	"github.com/webx-top/echo/handler/mvc/static/resource"
	"github.com/webx-top/echo/handler/pprof"
	mw "github.com/webx-top/echo/middleware"
	"github.com/webx-top/echo/middleware/render"
	"github.com/webx-top/echo/middleware/render/driver"
	"github.com/webx-top/echo/middleware/tplfunc"
)

// New 创建Application实例
func New(name string, middlewares ...interface{}) (s *Application) {
	s = NewWithContext(name, nil, middlewares...)
	globalApp[name] = s
	return
}

func HandlerWrapper(h interface{}) echo.Handler {
	if handle, y := h.(func(*Context) error); y {
		return echo.HandlerFunc(func(c echo.Context) error {
			return handle(c.(*Context))
		})
	}
	return nil
}

// NewWithContext 创建Application实例
func NewWithContext(name string, newContext func(*echo.Echo) echo.Context, middlewares ...interface{}) (s *Application) {
	s = &Application{
		Name:           name,
		moduleHosts:    make(map[string]*Module),
		moduleNames:    make(map[string]*Module),
		TemplateDir:    `template`,
		URL:            `/`,
		MaxUploadSize:  10 * 1024 * 1024,
		StaticDir:      `assets`,
		RootModuleName: `base`,
		FuncMap:        tplfunc.New(),
		RouteTagName:   `webx`,
		URLConvert:     LowerCaseFirst,
		URLRecovery:    UpperCaseFirst,
		MapperCheck:    DefaultMapperCheck,
	}
	mwNum := len(middlewares)
	if mwNum == 1 && middlewares[0] == nil {
		s.DefaultMiddlewares = []interface{}{}
	} else {
		s.DefaultMiddlewares = []interface{}{
			mw.Log(),
			mw.Recover(),
			mw.FuncMap(s.FuncMap, func(ctx echo.Context) bool {
				return ctx.Format() != `html`
			}),
		}
		if mwNum > 0 {
			s.DefaultMiddlewares = append(s.DefaultMiddlewares, middlewares...)
		}
	}
	s.SessionOptions = &echo.SessionOptions{
		Engine: `cookie`,
		Name:   `GOSID`,
		CookieOptions: &echo.CookieOptions{
			Prefix:   name + `_`,
			HttpOnly: true,
			Path:     `/`,
		},
	}
	s.FuncMap = s.DefaultFuncMap()
	s.ContextInitial = func(ctx echo.Context, wrp *Wrapper, controller interface{}, actionName string) (err error, exit bool) {
		return ctx.(*Context).Init(wrp, controller, actionName)
	} //[1]
	if newContext == nil {
		newContext = func(e *echo.Echo) echo.Context {
			return NewContext(s, echo.NewContext(nil, nil, e))
		}
	}
	s.ContextCreator = newContext //[1]
	s.Core = echo.NewWithContext(s.ContextCreator)
	s.Core.Use(s.DefaultMiddlewares...)
	s.Core.SetHandlerWrapper(HandlerWrapper) //[1]
	s.SetErrorPages(nil)

	s.URLs = NewURLs(name, s)
	return
}

type (
	//URLConvert 网址转换
	URLConvert func(string) string

	//URLRecovery 网址还原
	URLRecovery func(string) string
)

var (
	//SnakeCase 单词全部小写并用下划线连接
	SnakeCase URLConvert = com.SnakeCase

	//LowerCaseFirst 小写首字母
	LowerCaseFirst URLConvert = com.LowerCaseFirst

	//PascalCase 帕斯卡命名法
	PascalCase URLRecovery = com.PascalCase

	//UpperCaseFirst 大写首字母
	UpperCaseFirst URLRecovery = strings.Title
)

type Application struct {
	Core               *echo.Echo
	Name               string
	TemplateDir        string
	StaticDir          string
	StaticRes          *resource.Static
	RouteTagName       string
	URLConvert         URLConvert  `json:"-" xml:"-"`
	URLRecovery        URLRecovery `json:"-" xml:"-"`
	MaxUploadSize      int64
	RootModuleName     string
	URL                string
	URLs               *URLs
	DefaultMiddlewares []interface{} `json:"-" xml:"-"`
	SessionOptions     *echo.SessionOptions
	Renderer           driver.Driver                                                   `json:"-" xml:"-"`
	FuncMap            map[string]interface{}                                          `json:"-" xml:"-"`
	ContextCreator     func(*echo.Echo) echo.Context                                   `json:"-" xml:"-"`
	ContextInitial     func(echo.Context, *Wrapper, interface{}, string) (error, bool) `json:"-" xml:"-"`
	Codec              codec.Codec                                                     `json:"-" xml:"-"`
	MapperCheck        func(t reflect.Type) bool                                       `json:"-" xml:"-"`
	moduleHosts        map[string]*Module                                              //域名关联
	moduleNames        map[string]*Module                                              //名称关联
	rootDir            string
	theme              string
}

// ServeHTTP HTTP服务执行入口
func (s *Application) ServeHTTP(r engine.Request, w engine.Response) {
	var h *echo.Echo
	host := r.Host()
	module, ok := s.moduleHosts[host]
	if !ok {
		if p := strings.LastIndexByte(host, ':'); p > -1 {
			module, ok = s.moduleHosts[host[0:p]]
		}
	}
	if !ok || module.Handler == nil {
		h = s.Core
	} else {
		h = module.Handler
	}

	if h != nil {
		h.ServeHTTP(r, w)
	} else {
		w.NotFound()
	}
}

func (s *Application) SetErrorPages(templates map[int]string, options ...*render.Options) *Application {
	s.Core.SetHTTPErrorHandler(render.HTTPErrorHandler(templates, options...))
	return s
}

// SetDomain 为模块设置域名
func (s *Application) SetDomain(name string, domain string) *Application {
	a, ok := s.moduleNames[name]
	if !ok {
		s.Core.Logger().Warn(`Module does not exist: `, name)
		return s
	}
	if len(domain) == 0 { // 取消域名，加入到Core的Group中
		domain = a.Domain
		if _, ok := s.moduleHosts[domain]; !ok {
			return s
		}
		delete(s.moduleHosts, domain)
		var prefix string
		if name != s.RootModuleName {
			prefix = `/` + name
			a.Dir = prefix + `/`
		} else {
			a.Dir = `/`
		}
		routes := a.Handler.Routes()
		for _, r := range routes {
			if r.Path == `/` {
				if len(prefix) > 0 {
					r.Path = prefix
					r.Format = prefix
				}
			} else {
				r.Path = prefix + r.Path
				r.Format = prefix + r.Format
			}
			r.Prefix = prefix
		}
		a.URL = a.Dir
		if s.URL != `/` {
			a.URL = strings.TrimSuffix(s.URL, `/`) + a.URL
		}
		a.Domain = ``
		a.Group = s.Core.Group(prefix)
		a.Group.Use(a.Middlewares...)
		s.Core.AppendRouter(routes)
		a.Handler = nil
		return s
	}
	if len(a.Domain) > 0 { // 从一个域名换为另一个域名
		if a.Domain == domain {
			return s
		}
		if _, ok := s.moduleHosts[a.Domain]; ok {
			delete(s.moduleHosts, a.Domain)
		}
		s.moduleHosts[domain] = a
		a.Domain = domain
		return s
	}
	// 从Group移到域名
	s.moduleHosts[domain] = a
	routes := []*echo.Route{}
	coreRoutes := []*echo.Route{}
	for _, r := range s.Core.Routes() {
		if r.Prefix == `/`+name {
			if r.Path == `/`+name {
				r.Path = `/`
				r.Format = `/`
			} else {
				r.Path = `/` + strings.TrimPrefix(r.Path, `/`+name+`/`)
				r.Format = `/` + strings.TrimPrefix(r.Format, `/`+name+`/`)
			}
			r.Prefix = ``
			routes = append(routes, r)
		} else {
			coreRoutes = append(coreRoutes, r)
		}
	}
	a.Domain = domain
	a.Group = nil
	e := echo.NewWithContext(s.ContextCreator)
	e.Use(s.DefaultMiddlewares...)
	e.Use(a.Middlewares...)
	s.Core.RebuildRouter(coreRoutes)
	e.RebuildRouter(routes)
	a.Handler = e
	scheme := `http`
	if s.SessionOptions.Secure {
		scheme = `https`
	}
	a.URL = scheme + `://` + a.Domain + `/`
	a.Dir = `/`
	return s
}

// NewModule 创建新模块
func (s *Application) NewModule(name string, middlewares ...interface{}) *Module {
	r := strings.Split(name, `@`) //blog@www.blog.com
	var domain string
	if len(r) > 1 {
		name = r[0]
		domain = r[1]
	}
	a := NewModule(name, domain, s, middlewares...)
	if len(domain) > 0 {
		s.moduleHosts[domain] = a
	}
	s.moduleNames[name] = a
	return a
}

// NewRenderer 为特殊module(比如admin)单独新建渲染接口
func (s *Application) NewRenderer(conf *render.Config, a *Module, funcMap map[string]interface{}) driver.Driver {
	themeAbsPath := s.ThemeDir(conf.Theme)
	staticURLPath := `/assets`
	if a != nil && len(a.Name) > 0 {
		staticURLPath = `/` + a.Name + staticURLPath
	}
	staticAbsPath := themeAbsPath + `/assets`
	te := s.NewTemplateEngine(themeAbsPath, conf)
	static := s.NewStatic(staticURLPath, staticAbsPath, funcMap)
	te.SetFuncMap(func() map[string]interface{} {
		return funcMap
	})
	te.MonitorEvent(static.OnUpdate(themeAbsPath))
	te.SetContentProcessor(conf.Parser())
	return te
}

// NewTemplateEngine 新建模板引擎实例
func (s *Application) NewTemplateEngine(tmplPath string, conf *render.Config) driver.Driver {
	if tmplPath == `` {
		tmplPath = s.ThemeDir()
	}
	eng := render.New(conf.Engine, tmplPath, s.Core.Logger())
	eng.Init(true, conf.Reload)
	eng.SetDebug(conf.Debug)
	return eng
}

// 重置模板引擎
func (s *Application) resetRenderer(conf *render.Config) *Application {
	if s.Renderer != nil {
		s.Renderer.Close()
	}
	s.Renderer = s.NewTemplateEngine(s.ThemeDir(conf.Theme), conf)
	s.Core.SetRenderer(s.Renderer)
	s.TemplateMonitor()
	s.Renderer.SetFuncMap(func() map[string]interface{} {
		return s.FuncMap
	})
	s.Renderer.SetContentProcessor(conf.Parser())
	return s
}

// Module 获取模块实例
func (s *Application) Module(args ...string) *Module {
	name := s.RootModuleName
	if len(args) > 0 {
		name = args[0]
		if ap, ok := s.moduleNames[name]; ok {
			return ap
		}
	}
	return s.NewModule(name)
}

func (s *Application) SetSessionOptions(sessionOptions *echo.SessionOptions) *Application {
	if sessionOptions.CookieOptions == nil {
		sessionOptions.CookieOptions = &echo.CookieOptions{
			Path:     `/`,
			HttpOnly: true,
		}
	}
	if sessionOptions.Name == `` {
		sessionOptions.Name = `GOSID`
	}
	if sessionOptions.Engine == `` {
		sessionOptions.Engine = `cookie`
	}
	s.SessionOptions = sessionOptions
	return s
}

// ModuleOk 获取模块实例
func (s *Application) ModuleOk(args ...string) (app *Module, ok bool) {
	name := s.RootModuleName
	if len(args) > 0 {
		name = args[0]
	}
	app, ok = s.moduleNames[name]
	return
}

func (s *Application) Modules(args ...bool) map[string]*Module {
	if len(args) > 0 && args[0] {
		return s.moduleHosts
	}
	return s.moduleNames
}

func (s *Application) HasModule(name string) bool {
	_, ok := s.moduleNames[name]
	return ok
}

// NewStatic 新建静态资源实例
func (s *Application) NewStatic(urlPath string, absPath string, f ...map[string]interface{}) *resource.Static {
	st := resource.NewStatic(urlPath, absPath)
	if len(f) > 0 {
		f[0] = st.Register(f[0])
	}
	s.Core.Use(mw.Static(&mw.StaticOptions{Path: urlPath, Root: absPath}))
	return st
}

// ThemeDir 主题所在文件夹的路径
func (s *Application) ThemeDir(args ...string) string {
	if len(args) < 1 {
		return filepath.Join(s.TemplateDir, s.theme)
	}
	return filepath.Join(s.TemplateDir, args[0])
}

// InitStatic 初始化静态资源
func (s *Application) InitStatic() *Application {
	absPath := filepath.Join(s.ThemeDir(), s.StaticDir)
	s.StaticRes = s.NewStatic(s.StaticDir, absPath, s.FuncMap)
	if s.Renderer != nil {
		s.TemplateMonitor()
	}
	return s
}

// TemplateMonitor 模板监控事件
func (s *Application) TemplateMonitor() *Application {
	s.Renderer.MonitorEvent(s.StaticRes.OnUpdate(s.ThemeDir()))
	return s
}

// Theme 当前使用的主题名称
func (s *Application) Theme() string {
	return s.theme
}

// SetTheme 设置模板主题
func (s *Application) SetTheme(conf *render.Config) *Application {
	if conf.Theme == `admin` {
		return s
	}
	s.theme = conf.Theme
	s.resetRenderer(conf)
	return s
}

// LoadConfig 载入confl支持的配置文件
func (s *Application) LoadConfig(file string, config interface{}) error {
	content, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}
	return confl.Unmarshal(content, config)
}

// RootDir 网站根目录
func (s *Application) RootDir() string {
	if len(s.rootDir) == 0 {
		ppath := os.Getenv(strings.ToUpper(s.Name) + `PATH`)
		if len(ppath) == 0 {
			ppath, _ = filepath.Abs(os.Args[0])
			ppath = filepath.Dir(ppath)
		}
		s.rootDir = ppath
	}
	return s.rootDir
}

// Debug 开关debug模式
func (s *Application) Debug(on bool) *Application {
	s.Core.SetDebug(on)
	return s
}

// 运行之前准备数据
func (s *Application) ready() {
	s.Event(`mvc.serverReady`, func(_ bool) {})
}

func (s *Application) AddEvent(eventName string, handler interface{}) *Application {
	if h, ok := handler.(func(func(bool), ...interface{})); ok {
		events.AddEvent(eventName, h)
		return s
	}
	if h, ok := handler.(func(...interface{}) bool); ok {
		events.AddEvent(eventName, func(next func(bool), sessions ...interface{}) {
			next(h(sessions...))
		})
		return s
	}
	s.Core.Logger().Warnf(`Invalid event function: %T`, handler)
	return s
}

func (s *Application) Event(eventName string, next func(bool), sessions ...interface{}) *Application {
	events.Event(eventName, next, sessions...)
	return s
}

func (s *Application) GoEvent(eventName string, next func(bool), sessions ...interface{}) *Application {
	events.GoEvent(eventName, next, sessions...)
	return s
}

func (s *Application) DelEvent(eventName string) *Application {
	events.DelEvent(eventName)
	return s
}

// Run 运行服务
func (s *Application) Run(args ...interface{}) {
	s.ready()
	var eng engine.Engine
	var arg interface{}
	size := len(args)
	if size > 0 {
		arg = args[0]
	}
	if size > 1 {
		if conf, ok := arg.(*engine.Config); ok {
			if v, ok := args[1].(string); ok {
				if v == `fast` {
					eng = fasthttp.NewWithConfig(conf)
				} else {
					eng = standard.NewWithConfig(conf)
				}
			} else {
				eng = fasthttp.NewWithConfig(conf)
			}
		} else {
			addr := `:80`
			if v, ok := arg.(string); ok && v != `` {
				addr = v
			}
			if v, ok := args[1].(string); ok {
				if v == `fast` {
					eng = fasthttp.New(addr)
				} else {
					eng = standard.New(addr)
				}
			} else {
				eng = fasthttp.New(addr)
			}
		}
	} else {
		switch arg.(type) {
		case string:
			eng = fasthttp.New(arg.(string))
		case engine.Engine:
			eng = arg.(engine.Engine)
		default:
			eng = fasthttp.New(`:80`)
		}
	}
	s.Core.Logger().Infof(`Server "%v" has been launched.`, s.Name)
	s.Core.Run(eng, s)
	s.Core.Logger().Infof(`Server "%v" has been closed.`, s.Name)
	s.GoEvent(`mvc.serverExit`, func(_ bool) {})
}

// InitCodec 初始化 加密/解密 接口
func (s *Application) InitCodec(hashKey []byte, blockKey []byte) {
	s.Codec = codec.New(hashKey, blockKey)
}

// Pprof 启用pprof
func (s *Application) Pprof() *Application {
	pprof.Wrapper(s.Core)
	return s
}

func (s *Application) DefaultFuncMap() (r map[string]interface{}) {
	r = tplfunc.New()
	r["RootURL"] = func(p ...string) string {
		if len(p) > 0 {
			return s.URL + p[0]
		}
		return s.URL
	}
	return
}

// Tree  module -> controller -> action -> HTTP-METHODS
func (s *Application) Tree(args ...*echo.Echo) (r map[string]map[string]map[string]map[string]string) {
	core := s.Core
	if len(args) > 0 {
		core = args[0]
	}
	nrs := core.NamedRoutes()
	rs := core.Routes()
	r = map[string]map[string]map[string]map[string]string{}
	for name, indexes := range nrs {
		p := strings.LastIndex(name, `/`)
		s := strings.Split(name[p+1:], `.`)
		var appName, ctlName, actName string
		switch len(s) {
		case 3:
			if !strings.HasSuffix(s[2], `-fm`) {
				continue
			}
			actName = strings.TrimSuffix(s[2], `-fm`)
			ctlName = strings.TrimPrefix(s[1], `(`)
			ctlName = strings.TrimPrefix(ctlName, `*`)
			ctlName = strings.TrimSuffix(ctlName, `)`)
			p2 := strings.LastIndex(name[0:p], `/`)
			appName = name[p2+1 : p]
		default:
			continue
		}
		if _, ok := r[appName]; !ok {
			r[appName] = map[string]map[string]map[string]string{}
		}
		if _, ok := r[appName][ctlName]; !ok {
			r[appName][ctlName] = map[string]map[string]string{}
		}
		if _, ok := r[appName][ctlName][actName]; !ok {
			r[appName][ctlName][actName] = map[string]string{}
			for _, index := range indexes {
				route := rs[index]
				if _, ok := r[appName][ctlName][actName][route.Method]; !ok {
					r[appName][ctlName][actName][route.Method] = route.Method
				}
			}
		}
	}
	return
}
