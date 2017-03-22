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
	"reflect"
	"regexp"
	"strings"

	"github.com/webx-top/echo"
)

var (
	mapperType         = reflect.TypeOf(Mapper{})
	methodSuffixRegex  = regexp.MustCompile(`(?:_(?:` + strings.Join(echo.Methods(), `|`) + `))+$`)
	routeTagRegex      = regexp.MustCompile(`^[A-Z.]+(\|[A-Z]+)*$`)
	DefaultMapperCheck = func(t reflect.Type) bool {
		return t == mapperType
	}
	DefaultContextInitial = func(ctx echo.Context, wrp *Wrapper, controller interface{}, actionName string) (err error, exit bool) {
		return
	}
)

//Mapper 结构体中定义路由的字段类型
type Mapper struct{}

//BeforeHandler 静态实例中的前置行为
type BeforeHandler interface {
	Before(echo.Context) error
}

//AfterHandler 静态实例中的后置行为
type AfterHandler interface {
	After(echo.Context) error
}

//Initer 动态实例中的初始化行为
type Initer interface {
	Init(echo.Context) error
}

//Before 动态实例中的前置行为
type Before interface {
	Before() error
}

//Main 动态实例中的行为入口
type Main interface {
	Main() error
}

//After 动态实例中的后置行为
type After interface {
	After() error
}

type StaticIniter interface {
	Init(wrp *Wrapper, act string) error
}

type ExitChecker interface {
	IsExit() bool
}

type AllowFormat interface {
	AllowFormat(urlPath string, extension string) bool
}

func NewHandler(h func(echo.Context) error, name string) *Handler {
	return &Handler{
		name:   name,
		handle: h,
	}
}

type Handler struct {
	name   string
	handle func(echo.Context) error
}

func (h *Handler) Handle(c echo.Context) error {
	return h.handle(c)
}

func (h *Handler) HandleName() string {
	return h.name
}

type Wrapper struct {
	// 静态实例中的行为
	beforeHandler echo.HandlerFunc
	afterHandler  echo.HandlerFunc
	// 动态实例中的行为状态
	hasBefore bool
	hasMain   bool
	hasAfter  bool
	// 实例对象
	Controller     interface{}        `json:"-" xml:"-"`
	RouteRegister  echo.RouteRegister `json:"-" xml:"-"`
	Module         *Module            `json:"-" xml:"-"`
	ControllerName string
}

func (a *Wrapper) wrapHandler(v interface{}, ctl string, act string) func(echo.Context) error {
	h := a.Module.Core.ValidHandler(v)
	a.ControllerName = ctl
	if a.beforeHandler != nil && a.afterHandler != nil {
		return func(ctx echo.Context) error {
			if a.Module.ContextInitial != nil {
				if err, exit := a.Module.ContextInitial(ctx, a, nil, act); err != nil {
					return err
				} else if exit {
					return nil
				}
			}
			ex, ok := ctx.(ExitChecker)
			if err := a.beforeHandler(ctx); err != nil {
				return err
			}
			if ok && ex.IsExit() {
				return nil
			}
			if err := h.Handle(ctx); err != nil {
				return err
			}
			if ok && ex.IsExit() {
				return nil
			}
			return a.afterHandler(ctx)
		}
	}
	if a.beforeHandler != nil {
		return func(ctx echo.Context) error {
			if a.Module.ContextInitial != nil {
				if err, exit := a.Module.ContextInitial(ctx, a, nil, act); err != nil {
					return err
				} else if exit {
					return nil
				}
			}
			ex, ok := ctx.(ExitChecker)
			if err := a.beforeHandler(ctx); err != nil {
				return err
			}
			if ok && ex.IsExit() {
				return nil
			}
			return h.Handle(ctx)
		}
	}
	if a.afterHandler != nil {
		return func(ctx echo.Context) error {
			if a.Module.ContextInitial != nil {
				if err, exit := a.Module.ContextInitial(ctx, a, nil, act); err != nil {
					return err
				} else if exit {
					return nil
				}
			}
			ex, ok := ctx.(ExitChecker)
			if err := h.Handle(ctx); err != nil {
				return err
			}
			if ok && ex.IsExit() {
				return nil
			}
			return a.afterHandler(ctx)
		}
	}
	return func(ctx echo.Context) error {
		if a.Module.ContextInitial != nil {
			if err, exit := a.Module.ContextInitial(ctx, a, nil, act); err != nil {
				return err
			} else if exit {
				return nil
			}
		}
		return h.Handle(ctx)
	}
}

func (a *Wrapper) HandleName(h interface{}) string {
	return echo.HandlerName(h)
}

// Register 路由注册方案1：注册函数(可匿名)或静态实例的成员函数
// 例如：Register(`/index`,Index.Index,"GET","POST")
func (a *Wrapper) Register(p string, h interface{}, methods ...string) *Wrapper {
	if len(methods) < 1 {
		methods = append(methods, "GET")
	}
	_, ctl, act := a.Module.Application.URLs.Set(h)
	a.Module.Router().Match(methods, p, NewHandler(a.wrapHandler(h, ctl, act), a.HandleName(h)))
	return a
}

// RouteTags 路由注册方案2：从动态实例内Mapper类型字段标签中获取路由信息
func (a *Wrapper) RouteTags() {
	if _, y := a.Controller.(Initer); !y {
		a.Module.Core.Logger().Infof(`%T is no method Init(echo.Context)error, skiped.`, a.Controller)
		return
	}
	t := reflect.TypeOf(a.Controller)
	e := t.Elem()
	v := reflect.ValueOf(a.Controller)
	ctlPath := e.PkgPath() + `.(*` + e.Name() + `).`
	//github.com/webx-top/{Project}/app/{Module}/controller.(*Index).

	var ctl string
	if a.Module.URLConvert != nil {
		ctl = a.Module.URLConvert(e.Name())
	} else {
		ctl = e.Name()
	}
	a.ControllerName = ctl
	for i := 0; i < e.NumField(); i++ {
		f := e.Field(i)
		if !a.Module.Application.MapperCheck(f.Type) {
			continue
		}
		name := strings.Title(f.Name)
		m := v.MethodByName(name)
		if !m.IsValid() && !a.hasMain {
			m = v.MethodByName(`Main`)
			if !m.IsValid() {
				continue
			}
			a.hasMain = true
		}
		/*
			支持的tag:
			1. webx - 路由规则
			2. memo - 注释说明
			webx标签内容支持以下格式：
			1、只指定http请求方式，如`webx:"POST|GET"`
			2、只指定路由规则，如`webx:"index"`
			3、只指定扩展名规则，如`webx:".JSON|XML"`
			4、指定以上全部规则，如`webx:"GET|POST.JSON|XML index"`
			注: 当路径规则以波浪线"~"开头时，表示该规则加在app下，否则加在controller下
		*/
		tag := e.Field(i).Tag
		tagv := tag.Get(a.Module.RouteTagName)
		methods := []string{}
		extends := []string{}
		var p, w string
		if len(tagv) > 0 {
			tags := strings.Split(tagv, ` `)
			length := len(tags)
			if length >= 2 { //`webx:"GET|POST /index"`
				w = tags[0]
				p = tags[1]
			} else if length == 1 {
				if !routeTagRegex.MatchString(tags[0]) {
					//非全大写字母时，判断为网址规则
					p = tags[0]
				} else { //`webx:"GET|POST"`
					w = tags[0]
				}
			}
		}
		if len(p) == 0 {
			if a.Module.URLConvert != nil {
				p = `/` + a.Module.URLConvert(name)
			} else {
				p = `/` + name
			}
		}
		var ppath string
		if p[0] == '~' {
			p = p[1:]
			if p[0] != '/' {
				p = `/` + p
			}
			ppath = p
		} else {
			if p[0] != '/' {
				p = `/` + p
			}
			ppath = `/` + ctl + p
		}
		met := ``
		ext := ``
		if w != `` {
			me := strings.Split(w, `.`)
			met = me[0]
			if len(me) > 1 {
				ext = me[1]
			}
		}
		if met != `` {
			methods = strings.Split(met, `|`)
		}
		if ext != `` {
			ext = strings.ToLower(ext)
			extends = strings.Split(ext, `|`)
		}
		k := ctlPath + name + `-fm`
		u := a.Module.Application.URLs.SetExtensions(k, extends)
		h := NewHandler(func(ctx echo.Context) error {
			return a.execute(ctx, k, e, u, name)
		}, k)
		switch len(methods) {
		case 0:
			methods = append(methods, `GET`)
			methods = append(methods, `POST`)
		case 1:
			if methods[0] == `ANY` {
				a.addRouter(ctl, ppath, h)
				continue
			}
		}
		a.addRouter(ctl, ppath, h, methods...)
	}
}

func (a *Wrapper) addRouter(ctl string, ppath string, h *Handler, methods ...string) {
	isAnyMethods := true
	if len(methods) > 0 {
		isAnyMethods = false
	}
	router := a.Module.Router()
	if isAnyMethods {
		router.Any(ppath, h)
	} else {
		router.Match(methods, ppath, h)
	}

	if ctl == `index` && ppath != `/index/index` {
		if ppath == `/index` {
			if isAnyMethods {
				router.Any(`/`, h)
			} else {
				router.Match(methods, `/`, h)
			}
		} else if strings.HasPrefix(ppath, `/index/`) {
			if isAnyMethods {
				router.Any(strings.TrimPrefix(ppath, `/index`), h)
			} else {
				router.Match(methods, strings.TrimPrefix(ppath, `/index`), h)
			}
		}
	}

	for strings.HasSuffix(ppath, `/index`) {
		ppath = strings.TrimSuffix(ppath, `/index`)
		if isAnyMethods {
			router.Any(ppath, h)
		} else {
			router.Match(methods, ppath+`/`, h)
		}
	}
}

func (a *Wrapper) execute(c echo.Context, k string, e reflect.Type, u *URLs, action string) error {
	if !u.AllowFormat(k, c.Format()) {
		return c.HTML(`The contents can not be displayed in this format: `+c.Format(), 404)
	}
	return a.Exec(c, e, action)
}

func (a *Wrapper) Exec(ctx echo.Context, t reflect.Type, action string) error {
	v := reflect.New(t)
	ac := v.Interface()
	if a.Module.ContextInitial != nil {
		if err, exit := a.Module.ContextInitial(ctx, a, ac, action); err != nil {
			return err
		} else if exit {
			return nil
		}
	}
	if err := ac.(Initer).Init(ctx); err != nil {
		return err
	}
	ex, ok := ac.(ExitChecker)
	if ok && ex.IsExit() {
		return nil
	}
	if a.hasBefore {
		if err := ac.(Before).Before(); err != nil {
			return err
		}
		if ok && ex.IsExit() {
			return nil
		}
	}
	if a.hasMain {
		if err := ac.(Main).Main(); err != nil {
			return err
		}
	} else {
		if err := a.Module.ExecAction(action, t, v, ctx); err != nil {
			return err
		}
	}
	if a.hasAfter && (!ok || !ex.IsExit()) {
		return ac.(After).After()
	}
	return nil
}

// RouteMethods 路由注册方案3：自动注册动态实例内带HTTP方法名后缀的成员函数作为路由
func (a *Wrapper) RouteMethods() {
	if _, valid := a.Controller.(Initer); !valid {
		a.Module.Core.Logger().Infof(`%T is no method Init(echo.Context)error, skiped.`, a.Controller)
		return
	}
	t := reflect.TypeOf(a.Controller)
	e := t.Elem()
	ctlPath := e.PkgPath() + `.(*` + e.Name() + `).`
	//github.com/webx-top/{Project}/app/{Module}/controller.(*Index).
	var ctl string
	if a.Module.URLConvert != nil {
		ctl = a.Module.URLConvert(e.Name())
	} else {
		ctl = e.Name()
	}
	a.ControllerName = ctl
	for i := t.NumMethod() - 1; i >= 0; i-- {
		m := t.Method(i)
		name := m.Name
		h := func(k string, u *URLs) func(ctx echo.Context) error {
			return func(ctx echo.Context) error {
				return a.execute(ctx, k, e, u, name)
			}
		}
		if strings.HasSuffix(name, `_ANY`) {
			name = strings.TrimSuffix(name, `_ANY`)
			var pp string
			if a.Module.URLConvert != nil {
				pp = a.Module.URLConvert(name)
			} else {
				pp = name
			}
			ppath := `/` + ctl + `/` + pp
			k := ctlPath + name + `-fm`
			u := a.Module.Application.URLs
			handler := NewHandler(h(k, u), k)
			a.addRouter(ctl, ppath, handler)
			continue
		}
		matches := methodSuffixRegex.FindAllString(name, 1)
		if len(matches) < 1 {
			continue
		}
		methods := strings.Split(strings.TrimPrefix(matches[0], `_`), `_`)
		name = strings.TrimSuffix(name, matches[0])
		var pp string
		if a.Module.URLConvert != nil {
			pp = a.Module.URLConvert(name)
		} else {
			pp = name
		}
		ppath := `/` + ctl + `/` + pp
		k := ctlPath + name + `-fm`
		u := a.Module.Application.URLs
		handler := NewHandler(h(k, u), k)
		a.addRouter(ctl, ppath, handler, methods...)
	}
}

// Auto 自动注册动态实例的路由：a.Auto()
func (a *Wrapper) Auto(args ...interface{}) {
	if len(args) > 0 {
		a.RouteMethods()
		return
	}
	a.RouteTags()
}
