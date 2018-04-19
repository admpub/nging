package echo

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"sync"

	"github.com/admpub/log"
	"github.com/webx-top/echo/engine"
	"github.com/webx-top/echo/logger"
)

type (
	Echo struct {
		engine            engine.Engine
		middleware        []interface{}
		head              Handler
		maxParam          *int
		notFoundHandler   HandlerFunc
		httpErrorHandler  HTTPErrorHandler
		binder            Binder
		renderer          Renderer
		pool              sync.Pool
		debug             bool
		router            *Router
		logger            logger.Logger
		groups            map[string]*Group
		handlerWrapper    []func(interface{}) Handler
		middlewareWrapper []func(interface{}) Middleware
		acceptFormats     map[string]string //mime=>format
		FuncMap           map[string]interface{}
		RouteDebug        bool
		MiddlewareDebug   bool
	}

	Middleware interface {
		Handle(Handler) Handler
	}

	MiddlewareFunc func(Handler) Handler

	MiddlewareFuncd func(Handler) HandlerFunc

	Handler interface {
		Handle(Context) error
	}

	Name interface {
		Name() string
	}

	Meta interface {
		Meta() H
	}

	HandlerFunc func(Context) error

	// HTTPErrorHandler is a centralized HTTP error handler.
	HTTPErrorHandler func(error, Context)

	// Renderer is the interface that wraps the Render method.
	Renderer interface {
		Render(w io.Writer, name string, data interface{}, c Context) error
	}
)

// New creates an instance of Echo.
func New() (e *Echo) {
	return NewWithContext(func(e *Echo) Context {
		return NewContext(nil, nil, e)
	})
}

func NewWithContext(fn func(*Echo) Context) (e *Echo) {
	e = &Echo{maxParam: new(int)}
	e.pool.New = func() interface{} {
		return fn(e)
	}
	e.router = NewRouter(e)
	e.groups = make(map[string]*Group)

	//----------
	// Defaults
	//----------
	e.SetHTTPErrorHandler(e.DefaultHTTPErrorHandler)
	e.SetBinder(NewBinder(e))

	// Logger
	e.logger = log.GetLogger("echo")
	e.acceptFormats = map[string]string{
		//json
		`application/json`:       `json`,
		`text/javascript`:        `json`,
		`application/javascript`: `json`,

		//xml
		`application/xml`: `xml`,
		`text/xml`:        `xml`,

		//text
		`text/plain`: `text`,

		//html
		`*/*`:               `html`,
		`application/xhtml`: `html`,
		`text/html`:         `html`,

		//default
		`*`: `html`,
	}
	return
}

func (m MiddlewareFunc) Handle(h Handler) Handler {
	return m(h)
}

func (m MiddlewareFuncd) Handle(h Handler) Handler {
	return m(h)
}

func (h HandlerFunc) Handle(c Context) error {
	return h(c)
}

func (e *Echo) SetAcceptFormats(acceptFormats map[string]string) *Echo {
	e.acceptFormats = acceptFormats
	return e
}

func (e *Echo) AddAcceptFormat(mime, format string) *Echo {
	e.acceptFormats[mime] = format
	return e
}

// Router returns router.
func (e *Echo) Router() *Router {
	return e.router
}

// SetLogger sets the logger instance.
func (e *Echo) SetLogger(l logger.Logger) {
	e.logger = l
}

// Logger returns the logger instance.
func (e *Echo) Logger() logger.Logger {
	return e.logger
}

// DefaultHTTPErrorHandler invokes the default HTTP error handler.
func (e *Echo) DefaultHTTPErrorHandler(err error, c Context) {
	code := http.StatusInternalServerError
	msg := http.StatusText(code)
	if he, ok := err.(*HTTPError); ok {
		code = he.Code
		msg = he.Message
	}
	if e.debug {
		msg = err.Error()
	}
	if !c.Response().Committed() {
		if c.Request().Method() == HEAD {
			c.NoContent(code)
		} else {
			if code > 0 {
				c.String(msg, code)
			} else {
				c.String(msg)
			}
		}
	}
	e.logger.Debug(err)
}

// SetHTTPErrorHandler registers a custom Echo.HTTPErrorHandler.
func (e *Echo) SetHTTPErrorHandler(h HTTPErrorHandler) {
	e.httpErrorHandler = h
}

// HTTPErrorHandler returns the HTTPErrorHandler
func (e *Echo) HTTPErrorHandler() HTTPErrorHandler {
	return e.httpErrorHandler
}

// SetBinder registers a custom binder. It's invoked by Context.Bind().
func (e *Echo) SetBinder(b Binder) {
	e.binder = b
}

// Binder returns the binder instance.
func (e *Echo) Binder() Binder {
	return e.binder
}

// SetRenderer registers an HTML template renderer. It's invoked by Context.Render().
func (e *Echo) SetRenderer(r Renderer) {
	e.renderer = r
}

// Renderer returns the renderer instance.
func (e *Echo) Renderer() Renderer {
	return e.renderer
}

// SetDebug enable/disable debug mode.
func (e *Echo) SetDebug(on bool) {
	e.debug = on
	if logger, ok := e.logger.(logger.LevelSetter); ok {
		if on {
			logger.SetLevel(`Debug`)
		} else {
			logger.SetLevel(`Info`)
		}
	}
}

// Debug returns debug mode (enabled or disabled).
func (e *Echo) Debug() bool {
	return e.debug
}

// Use adds handler to the middleware chain.
func (e *Echo) Use(middleware ...interface{}) {
	for _, m := range middleware {
		e.ValidMiddleware(m)
		e.middleware = append(e.middleware, m)
		if e.MiddlewareDebug {
			e.logger.Debugf(`Middleware[Use](%p): [] -> %s `, m, HandlerName(m))
		}
	}
}

// Pre is an alias for `PreUse` function.
func (e *Echo) Pre(middleware ...interface{}) {
	e.PreUse(middleware...)
}

// PreUse adds handler to the middleware chain.
func (e *Echo) PreUse(middleware ...interface{}) {
	var middlewares []interface{}
	for _, m := range middleware {
		e.ValidMiddleware(m)
		middlewares = append(middlewares, m)
		if e.MiddlewareDebug {
			e.logger.Debugf(`Middleware[Pre](%p): [] -> %s`, m, HandlerName(m))
		}
	}
	e.middleware = append(middlewares, e.middleware...)
}

// Clear middleware
func (e *Echo) Clear(middleware ...interface{}) {
	if len(middleware) > 0 {
		for _, dm := range middleware {
			var decr int
			for i, m := range e.middleware {
				if m != dm {
					continue
				}
				i -= decr
				start := i + 1
				if start < len(e.middleware) {
					e.middleware = append(e.middleware[0:i], e.middleware[start:]...)
				} else {
					e.middleware = e.middleware[0:i]
				}
				decr++
			}
		}
	} else {
		e.middleware = []interface{}{}
	}
	e.head = nil
}

// Connect adds a CONNECT route > handler to the router.
func (e *Echo) Connect(path string, h interface{}, m ...interface{}) {
	e.add(CONNECT, "", path, h, m...)
}

// Delete adds a DELETE route > handler to the router.
func (e *Echo) Delete(path string, h interface{}, m ...interface{}) {
	e.add(DELETE, "", path, h, m...)
}

// Get adds a GET route > handler to the router.
func (e *Echo) Get(path string, h interface{}, m ...interface{}) {
	e.add(GET, "", path, h, m...)
}

// Head adds a HEAD route > handler to the router.
func (e *Echo) Head(path string, h interface{}, m ...interface{}) {
	e.add(HEAD, "", path, h, m...)
}

// Options adds an OPTIONS route > handler to the router.
func (e *Echo) Options(path string, h interface{}, m ...interface{}) {
	e.add(OPTIONS, "", path, h, m...)
}

// Patch adds a PATCH route > handler to the router.
func (e *Echo) Patch(path string, h interface{}, m ...interface{}) {
	e.add(PATCH, "", path, h, m...)
}

// Post adds a POST route > handler to the router.
func (e *Echo) Post(path string, h interface{}, m ...interface{}) {
	e.add(POST, "", path, h, m...)
}

// Put adds a PUT route > handler to the router.
func (e *Echo) Put(path string, h interface{}, m ...interface{}) {
	e.add(PUT, "", path, h, m...)
}

// Trace adds a TRACE route > handler to the router.
func (e *Echo) Trace(path string, h interface{}, m ...interface{}) {
	e.add(TRACE, "", path, h, m...)
}

// Any adds a route > handler to the router for all HTTP methods.
func (e *Echo) Any(path string, h interface{}, middleware ...interface{}) {
	for _, m := range methods {
		e.add(m, "", path, h, middleware...)
	}
}

func (e *Echo) Route(methods string, path string, h interface{}, middleware ...interface{}) {
	e.Match(splitHTTPMethod.Split(methods, -1), path, h, middleware...)
}

// Match adds a route > handler to the router for multiple HTTP methods provided.
func (e *Echo) Match(methods []string, path string, h interface{}, middleware ...interface{}) {
	for _, m := range methods {
		e.add(m, "", path, h, middleware...)
	}
}

// Static registers a new route with path prefix to serve static files from the
// provided root directory.
func (e *Echo) Static(prefix, root string) {
	if root == "" {
		root = "." // For security we want to restrict to CWD.
	}
	static(e, prefix, root)
}

// File registers a new route with path to serve a static file.
func (e *Echo) File(path, file string) {
	e.Get(path, func(c Context) error {
		return c.File(file)
	})
}

func (e *Echo) ValidHandler(v interface{}) (h Handler) {
	if e.handlerWrapper != nil {
		for _, wrapper := range e.handlerWrapper {
			h = wrapper(v)
			if h != nil {
				return
			}
		}
	}
	return WrapHandler(v)
}

func (e *Echo) ValidMiddleware(v interface{}) (m Middleware) {
	if e.middlewareWrapper != nil {
		for _, wrapper := range e.middlewareWrapper {
			m = wrapper(v)
			if m != nil {
				return
			}
		}
	}
	return WrapMiddleware(v)
}

func (e *Echo) SetHandlerWrapper(funcs ...func(interface{}) Handler) {
	e.handlerWrapper = funcs
}

func (e *Echo) SetMiddlewareWrapper(funcs ...func(interface{}) Middleware) {
	e.middlewareWrapper = funcs
}

func (e *Echo) AddHandlerWrapper(funcs ...func(interface{}) Handler) {
	e.handlerWrapper = append(e.handlerWrapper, funcs...)
}

func (e *Echo) AddMiddlewareWrapper(funcs ...func(interface{}) Middleware) {
	e.middlewareWrapper = append(e.middlewareWrapper, funcs...)
}

func (e *Echo) add(method, prefix string, path string, h interface{}, middleware ...interface{}) {
	r := &Route{
		Method:     method,
		Path:       path,
		Prefix:     prefix,
		handler:    h,
		middleware: middleware,
	}
	r.apply(e)
	rid := len(e.router.routes)
	e.router.Add(r, rid)
	if e.RouteDebug {
		e.logger.Debugf(`Route: %7v %-30v -> %v`, method, r.Format, r.HandlerName)
	}
	if _, ok := e.router.nroute[r.HandlerName]; !ok {
		e.router.nroute[r.HandlerName] = []int{rid}
	} else {
		e.router.nroute[r.HandlerName] = append(e.router.nroute[r.HandlerName], rid)
	}
	e.router.routes = append(e.router.routes, r)
}

// MetaHandler Add meta information about endpoint
func (e *Echo) MetaHandler(m H, handler interface{}) Handler {
	return &MetaHandler{m, e.ValidHandler(handler)}
}

// RebuildRouter rebuild router
func (e *Echo) RebuildRouter(args ...[]*Route) {
	routes := e.router.routes
	if len(args) > 0 {
		routes = args[0]
	}
	e.router = NewRouter(e)
	for i, r := range routes {
		//e.logger.Debugf(`%p rebuild: %#v`, e, *r)
		r.apply(e)
		e.router.Add(r, i)

		if _, ok := e.router.nroute[r.HandlerName]; !ok {
			e.router.nroute[r.HandlerName] = []int{i}
		} else {
			e.router.nroute[r.HandlerName] = append(e.router.nroute[r.HandlerName], i)
		}
	}
	e.router.routes = routes
	e.head = nil
}

// AppendRouter append router
func (e *Echo) AppendRouter(routes []*Route) {
	for i, r := range routes {
		i = len(e.router.routes)
		r.apply(e)
		e.router.Add(r, i)
		if _, ok := e.router.nroute[r.HandlerName]; !ok {
			e.router.nroute[r.HandlerName] = []int{i}
		} else {
			e.router.nroute[r.HandlerName] = append(e.router.nroute[r.HandlerName], i)
		}
		e.router.routes = append(e.router.routes, r)
	}
	e.head = nil
}

// Group creates a new sub-router with prefix.
func (e *Echo) Group(prefix string, m ...interface{}) (g *Group) {
	g = &Group{prefix: prefix, echo: e}
	g.Use(m...)
	e.groups[prefix] = g
	return
}

func (e *Echo) GetGroup(prefix string) (g *Group) {
	g, _ = e.groups[prefix]
	return
}

// URI generates a URI from handler.
func (e *Echo) URI(handler interface{}, params ...interface{}) string {
	var uri, name string
	if h, ok := handler.(Handler); ok {
		if hn, ok := h.(Name); ok {
			name = hn.Name()
		} else {
			name = HandlerName(h)
		}
	} else if h, ok := handler.(string); ok {
		name = h
	} else {
		return uri
	}
	if indexes, ok := e.router.nroute[name]; ok && len(indexes) > 0 {
		r := e.router.routes[indexes[0]]
		length := len(params)
		if length == 1 {
			switch val := params[0].(type) {
			case url.Values:
				uri = r.Path
				for _, name := range r.Params {
					tag := `:` + name
					v := val.Get(name)
					uri = strings.Replace(uri, tag+`/`, v+`/`, -1)
					if strings.HasSuffix(uri, tag) {
						uri = strings.TrimSuffix(uri, tag) + v
					}
					val.Del(name)
				}
				q := val.Encode()
				if len(q) > 0 {
					uri += `?` + q
				}
			case map[string]string:
				uri = r.Path
				for _, name := range r.Params {
					tag := `:` + name
					v, y := val[name]
					if y {
						delete(val, name)
					}
					uri = strings.Replace(uri, tag+`/`, v+`/`, -1)
					if strings.HasSuffix(uri, tag) {
						uri = strings.TrimSuffix(uri, tag) + v
					}
				}
				sep := `?`
				keys := make([]string, 0, len(val))
				for k := range val {
					keys = append(keys, k)
				}
				sort.Strings(keys)
				for _, k := range keys {
					uri += sep + url.QueryEscape(k) + `=` + url.QueryEscape(val[k])
					sep = `&`
				}
			case []interface{}:
				uri = fmt.Sprintf(r.Format, val...)
			default:
				uri = fmt.Sprintf(r.Format, val)
			}
		} else {
			uri = fmt.Sprintf(r.Format, params...)
		}
	}
	return uri
}

// URL is an alias for `URI` function.
func (e *Echo) URL(h interface{}, params ...interface{}) string {
	return e.URI(h, params...)
}

// Routes returns the registered routes.
func (e *Echo) Routes() []*Route {
	return e.router.routes
}

// NamedRoutes returns the registered handler name.
func (e *Echo) NamedRoutes() map[string][]int {
	return e.router.nroute
}

// Chain middleware
func (e *Echo) chainMiddleware() {
	if e.head != nil {
		return
	}
	e.head = e.router.Handle(nil)
	for i := len(e.middleware) - 1; i >= 0; i-- {
		e.head = e.ValidMiddleware(e.middleware[i]).Handle(e.head)
	}
}

func (e *Echo) ServeHTTP(req engine.Request, res engine.Response) {
	c := e.pool.Get().(Context)
	c.Reset(req, res)

	e.chainMiddleware()

	if err := e.head.Handle(c); err != nil {
		c.Error(err)
	}

	e.pool.Put(c)
}

// Run starts the HTTP engine.
func (e *Echo) Run(eng engine.Engine, handler ...engine.Handler) {
	e.setEngine(eng).start(handler...)
}

func (e *Echo) start(handler ...engine.Handler) {
	if len(handler) > 0 {
		e.engine.SetHandler(handler[0])
	} else {
		e.engine.SetHandler(e)
	}
	e.engine.SetLogger(e.logger)
	if e.Debug() {
		e.logger.Debug("running in debug mode")
	}
	e.engine.Start()
}

func (e *Echo) setEngine(eng engine.Engine) *Echo {
	e.engine = eng
	return e
}

func (e *Echo) Engine() engine.Engine {
	return e.engine
}

// Stop stops the HTTP server.
func (e *Echo) Stop() error {
	if e.engine == nil {
		return nil
	}
	return e.engine.Stop()
}
