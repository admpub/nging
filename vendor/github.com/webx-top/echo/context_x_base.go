package echo

import (
	"bytes"
	"context"
	"fmt"
	"time"

	"github.com/admpub/events"
	"github.com/admpub/events/emitter"

	pkgCode "github.com/webx-top/echo/code"
	"github.com/webx-top/echo/engine"
	"github.com/webx-top/echo/logger"
	"github.com/webx-top/echo/param"
)

type xContext struct {
	Validator
	Translator
	events.Emitter
	transaction         *BaseTransaction
	sessioner           Sessioner
	cookier             Cookier
	context             context.Context
	request             engine.Request
	response            engine.Response
	path                string
	pnames              []string
	pvalues             []string
	hnames              []string
	hvalues             []string
	store               Store
	internal            *param.SafeMap
	handler             Handler
	route               *Route
	rid                 int
	echo                *Echo
	funcs               map[string]interface{}
	renderer            Renderer
	sessionOptions      *SessionOptions
	withFormatExtension bool
	format              string
	code                int
	preResponseHook     []func() error
	dataEngine          Data
	accept              *Accepts
	auto                bool
}

// NewContext creates a Context object.
func NewContext(req engine.Request, res engine.Response, e *Echo) Context {
	c := &xContext{
		Validator:   e.Validator,
		Translator:  DefaultNopTranslate,
		Emitter:     emitter.DefaultCondEmitter,
		transaction: DefaultNopTransaction,
		context:     context.Background(),
		request:     req,
		response:    res,
		echo:        e,
		pvalues:     make([]string, *e.maxParam),
		internal:    param.NewMap(),
		store:       make(Store),
		handler:     NotFoundHandler,
		funcs:       make(map[string]interface{}),
		sessioner:   DefaultSession,
	}
	c.cookier = NewCookier(c)
	c.dataEngine = NewData(c)
	return c
}

func (c *xContext) StdContext() context.Context {
	return c.context
}

func (c *xContext) Internal() *param.SafeMap {
	return c.internal
}

func (c *xContext) SetStdContext(ctx context.Context) {
	c.context = ctx
}

func (c *xContext) SetEmitter(emitter events.Emitter) {
	c.Emitter = emitter
}

func (c *xContext) Handler() Handler {
	return c.handler
}

func (c *xContext) Deadline() (deadline time.Time, ok bool) {
	return c.context.Deadline()
}

func (c *xContext) Done() <-chan struct{} {
	return c.context.Done()
}

func (c *xContext) Err() error {
	return c.context.Err()
}

func (c *xContext) Value(key interface{}) interface{} {
	return c.context.Value(key)
}

func (c *xContext) Handle(ctx Context) error {
	return c.handler.Handle(ctx)
}

func (c *xContext) Route() *Route {
	if c.route == nil {
		if c.rid < 0 || c.rid >= len(c.echo.router.routes) {
			c.route = defaultRoute
		} else {
			c.route = c.echo.router.routes[c.rid]
		}
	}
	return c.route
}

func (c *xContext) SetAuto(on bool) Context {
	c.auto = on
	return c
}

// Error invokes the registered HTTP error handler. Generally used by middleware.
func (c *xContext) Error(err error) {
	c.echo.httpErrorHandler(err, c)
}

func (c *xContext) NewError(code pkgCode.Code, msg string, args ...interface{}) *Error {
	return NewError(c.T(msg, args...), code)
}

func (c *xContext) NewErrorWith(err error, code pkgCode.Code, args ...interface{}) *Error {
	var msg string
	if len(args) > 0 {
		msg = param.AsString(args[0])
		if len(args) > 1 {
			msg = c.T(msg, args[1:]...)
		} else {
			msg = c.T(msg)
		}
	}
	return NewErrorWith(err, msg, code)
}

// Logger returns the `Logger` instance.
func (c *xContext) Logger() logger.Logger {
	return c.echo.logger
}

// Object returns the `context` object.
func (c *xContext) Object() *xContext {
	return c
}

// Echo returns the `Echo` instance.
func (c *xContext) Echo() *Echo {
	return c.echo
}

func (c *xContext) SetTranslator(t Translator) {
	c.Translator = t
}

func (c *xContext) Reset(req engine.Request, res engine.Response) {
	c.Validator = c.echo.Validator
	c.Emitter = emitter.DefaultCondEmitter
	c.Translator = DefaultNopTranslate
	c.transaction = DefaultNopTransaction
	c.sessioner = DefaultSession
	c.cookier = NewCookier(c)
	c.context = context.Background()
	c.request = req
	c.response = res
	c.internal = param.NewMap()
	c.store = make(Store)
	c.path = ""
	c.pnames = nil
	c.hnames = nil
	c.hvalues = nil
	c.funcs = make(map[string]interface{})
	c.renderer = nil
	c.handler = NotFoundHandler
	c.route = nil
	c.rid = -1
	c.sessionOptions = nil
	c.withFormatExtension = false
	c.format = ""
	c.code = 0
	c.auto = false
	c.preResponseHook = nil
	c.accept = nil
	c.dataEngine = NewData(c)
	// NOTE: Don't reset because it has to have length c.echo.maxParam at all times
	// c.pvalues = nil
}

func (c *xContext) GetFunc(key string) interface{} {
	return c.funcs[key]
}

func (c *xContext) SetFunc(key string, val interface{}) {
	c.funcs[key] = val
}

func (c *xContext) ResetFuncs(funcs map[string]interface{}) {
	c.funcs = funcs
}

func (c *xContext) Funcs() map[string]interface{} {
	return c.funcs
}

func (c *xContext) Fetch(name string, data interface{}) (b []byte, err error) {
	if c.renderer == nil {
		if c.echo.renderer == nil {
			return nil, ErrRendererNotRegistered
		}
		c.renderer = c.echo.renderer
	}
	buf := new(bytes.Buffer)
	err = c.renderer.Render(buf, name, data, c)
	if err != nil {
		return
	}
	b = buf.Bytes()
	return
}

func (c *xContext) SetValidator(v Validator) {
	c.Validator = v
}

// SetRenderer registers an HTML template renderer.
func (c *xContext) SetRenderer(r Renderer) {
	c.renderer = r
}

func (c *xContext) SetSessioner(s Sessioner) {
	c.sessioner = s
}

func (c *xContext) Atop(v string) param.String {
	return param.String(v)
}

func (c *xContext) ToParamString(v string) param.String {
	return param.String(v)
}

func (c *xContext) ToStringSlice(v []string) param.StringSlice {
	return param.StringSlice(v)
}

func (c *xContext) SetFormat(format string) {
	c.format = format
}

func (c *xContext) WithFormatExtension(on bool) {
	c.withFormatExtension = on
}

func (c *xContext) SetCode(code int) {
	c.code = code
}

func (c *xContext) Code() int {
	return c.code
}

func (c *xContext) SetData(data Data) {
	c.dataEngine = data
}

func (c *xContext) Data() Data {
	return c.dataEngine
}

// MapData 映射数据到结构体
func (c *xContext) MapData(i interface{}, data map[string][]string, names ...string) error {
	var name string
	if len(names) > 0 {
		name = names[0]
	}
	return NamedStructMap(c.echo, i, data, name)
}

func (c *xContext) AddPreResponseHook(hook func() error) Context {
	if c.preResponseHook == nil {
		c.preResponseHook = []func() error{hook}
	} else {
		c.preResponseHook = append(c.preResponseHook, hook)
	}
	return c
}

func (c *xContext) SetPreResponseHook(hook ...func() error) Context {
	c.preResponseHook = hook
	return c
}

func (c *xContext) preResponse() error {
	if c.preResponseHook == nil {
		c.cookier.Send()
		return nil
	}
	for _, hook := range c.preResponseHook {
		if err := hook(); err != nil {
			return err
		}
	}
	c.cookier.Send()
	return nil
}

func (c *xContext) PrintFuncs() {
	for key, fn := range c.Funcs() {
		fmt.Printf("[Template Func](%p) %-15s -> %s \n", fn, key, HandlerName(fn))
	}
}
