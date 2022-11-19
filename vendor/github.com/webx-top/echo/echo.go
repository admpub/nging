package echo

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"strings"
	"sync"

	"github.com/admpub/log"

	"github.com/webx-top/echo/engine"
	"github.com/webx-top/echo/logger"
)

type (
	Echo struct {
		engine             engine.Engine
		prefix             string
		premiddleware      []interface{}
		middleware         []interface{}
		hosts              map[string]*Host
		hostAlias          map[string]string
		onHostFound        func(Context) (bool, error)
		maxParam           *int
		notFoundHandler    HandlerFunc
		httpErrorHandler   HTTPErrorHandler
		binder             Binder
		renderer           Renderer
		renderDataWrapper  DataWrapper
		pool               sync.Pool
		debug              bool
		router             *Router
		logger             logger.Logger
		groups             map[string]*Group
		handlerWrapper     []func(interface{}) Handler
		middlewareWrapper  []func(interface{}) Middleware
		acceptFormats      map[string]string //mime=>format
		formatRenderers    map[string]FormatRender
		FuncMap            map[string]interface{}
		RouteDebug         bool
		MiddlewareDebug    bool
		JSONPVarName       string
		Validator          Validator
		FormSliceMaxIndex  int
		parseHeaderAccept  bool
		defaultExtension   string
		rewriter           Rewriter
		maxRequestBodySize int
	}

	Middleware interface {
		Handle(Handler) Handler
	}

	AsMiddleware func(Context) error

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
		RenderBy(w io.Writer, name string, content func(string) ([]byte, error), data interface{}, c Context) error
	}
)

func (m MiddlewareFunc) Handle(h Handler) Handler {
	return m(h)
}

func (m MiddlewareFuncd) Handle(h Handler) Handler {
	return m(h)
}

func (m AsMiddleware) Handle(h Handler) Handler {
	return HandlerFunc(func(c Context) error {
		if err := m(c); err != nil {
			return err
		}
		return h.Handle(c)
	})
}

func (h HandlerFunc) Handle(c Context) error {
	return h(c)
}

// New creates an instance of Echo.
func New() (e *Echo) {
	return NewWithContext(func(e *Echo) Context {
		return NewContext(nil, nil, e)
	})
}

func NewWithContext(fn func(*Echo) Context) (e *Echo) {
	e = &Echo{}
	e.pool.New = func() interface{} {
		return fn(e)
	}
	return e.Reset()
}

func (e *Echo) Reset() *Echo {
	e.engine = nil
	e.prefix = ``
	e.middleware = []interface{}{}
	e.premiddleware = []interface{}{}
	e.hosts = make(map[string]*Host)
	e.hostAlias = make(map[string]string)
	e.maxParam = new(int)
	e.SetHTTPErrorHandler(e.DefaultHTTPErrorHandler)
	e.SetBinder(NewBinder(e))
	e.notFoundHandler = nil
	e.renderer = nil
	e.debug = false
	e.router = NewRouter(e)
	e.logger = log.GetLogger("echo")
	e.groups = make(map[string]*Group)
	e.handlerWrapper = []func(interface{}) Handler{}
	e.middlewareWrapper = []func(interface{}) Middleware{}
	e.acceptFormats = DefaultAcceptFormats
	e.formatRenderers = DefaultFormatRenderers
	e.FuncMap = make(map[string]interface{})
	e.RouteDebug = false
	e.MiddlewareDebug = false
	e.JSONPVarName = `callback`
	e.Validator = DefaultNopValidate
	e.FormSliceMaxIndex = 100
	e.parseHeaderAccept = false
	e.defaultExtension = ``
	e.maxRequestBodySize = 0
	e.renderDataWrapper = nil
	e.rewriter = nil
	return e
}

func (e *Echo) ParseHeaderAccept(on bool) *Echo {
	e.parseHeaderAccept = on
	return e
}

func (e *Echo) SetValidator(validator Validator) *Echo {
	e.Validator = validator
	return e
}

func (e *Echo) SetFormSliceMaxIndex(max int) *Echo {
	e.FormSliceMaxIndex = max
	return e
}

func (e *Echo) SetAcceptFormats(acceptFormats map[string]string) *Echo {
	e.acceptFormats = acceptFormats
	return e
}

func (e *Echo) SetMaxRequestBodySize(maxRequestBodySize int) *Echo {
	e.maxRequestBodySize = maxRequestBodySize
	return e
}

func (e *Echo) MaxRequestBodySize() int {
	if e.engine != nil && e.engine.Config().MaxRequestBodySize > 0 {
		return e.engine.Config().MaxRequestBodySize
	}
	return e.maxRequestBodySize
}

func (e *Echo) AddAcceptFormat(mime, format string) *Echo {
	e.acceptFormats[mime] = format
	return e
}

func (e *Echo) SetFormatRenderers(formatRenderers map[string]FormatRender) *Echo {
	e.formatRenderers = formatRenderers
	return e
}

func (e *Echo) AddFormatRenderer(format string, renderer func(c Context, data interface{}) error) *Echo {
	e.formatRenderers[format] = renderer
	return e
}

func (e *Echo) RemoveFormatRenderer(formats ...string) *Echo {
	for _, format := range formats {
		delete(e.formatRenderers, format)
	}
	return e
}

func (e *Echo) SetDefaultExtension(ext string) {
	e.defaultExtension = ext
}

func (e *Echo) DefaultExtension() string {
	return e.defaultExtension
}

// Router returns router.
func (e *Echo) Router() *Router {
	return e.router
}

// Hosts returns the map of host => Host.
func (e *Echo) Hosts() map[string]*Host {
	return e.hosts
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
	e.logger.Debug(err, `: `, c.Request().URL().String())
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

// SetRenderDataWrapper .
func (e *Echo) SetRenderDataWrapper(dataWrapper DataWrapper) {
	e.renderDataWrapper = dataWrapper
}

// RenderDataWrapper .
func (e *Echo) RenderDataWrapper() DataWrapper {
	return e.renderDataWrapper
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

// Pre adds handler to the middleware chain.
func (e *Echo) Pre(middleware ...interface{}) {
	var middlewares []interface{}
	for _, m := range middleware {
		e.ValidMiddleware(m)
		middlewares = append(middlewares, m)
		if e.MiddlewareDebug {
			e.logger.Debugf(`Middleware[Pre](%p): [] -> %s`, m, HandlerName(m))
		}
	}
	e.premiddleware = append(middlewares, e.premiddleware...)
}

// Clear middleware
func (e *Echo) Clear(middleware ...interface{}) {
	e.middleware = Clear(e.middleware, middleware...)
}

// ClearPre Clear premiddleware
func (e *Echo) ClearPre(middleware ...interface{}) {
	e.premiddleware = Clear(e.premiddleware, middleware...)
}

// Connect adds a CONNECT route > handler to the router.
func (e *Echo) Connect(path string, h interface{}, m ...interface{}) IRouter {
	return e.Add(CONNECT, path, h, m...)
}

// Delete adds a DELETE route > handler to the router.
func (e *Echo) Delete(path string, h interface{}, m ...interface{}) IRouter {
	return e.Add(DELETE, path, h, m...)
}

// Get adds a GET route > handler to the router.
func (e *Echo) Get(path string, h interface{}, m ...interface{}) IRouter {
	return e.Add(GET, path, h, m...)
}

// Head adds a HEAD route > handler to the router.
func (e *Echo) Head(path string, h interface{}, m ...interface{}) IRouter {
	return e.Add(HEAD, path, h, m...)
}

// Options adds an OPTIONS route > handler to the router.
func (e *Echo) Options(path string, h interface{}, m ...interface{}) IRouter {
	return e.Add(OPTIONS, path, h, m...)
}

// Patch adds a PATCH route > handler to the router.
func (e *Echo) Patch(path string, h interface{}, m ...interface{}) IRouter {
	return e.Add(PATCH, path, h, m...)
}

// Post adds a POST route > handler to the router.
func (e *Echo) Post(path string, h interface{}, m ...interface{}) IRouter {
	return e.Add(POST, path, h, m...)
}

// Put adds a PUT route > handler to the router.
func (e *Echo) Put(path string, h interface{}, m ...interface{}) IRouter {
	return e.Add(PUT, path, h, m...)
}

// Trace adds a TRACE route > handler to the router.
func (e *Echo) Trace(path string, h interface{}, m ...interface{}) IRouter {
	return e.Add(TRACE, path, h, m...)
}

// Any adds a route > handler to the router for all HTTP methods.
func (e *Echo) Any(path string, h interface{}, middleware ...interface{}) IRouter {
	routes := Routes{}
	for _, m := range methods {
		routes = append(routes, e.Add(m, path, h, middleware...))
	}
	return routes
}

func (e *Echo) Route(methods string, path string, h interface{}, middleware ...interface{}) IRouter {
	return e.Match(splitHTTPMethod.Split(methods, -1), path, h, middleware...)
}

// Match adds a route > handler to the router for multiple HTTP methods provided.
func (e *Echo) Match(methods []string, path string, h interface{}, middleware ...interface{}) IRouter {
	routes := Routes{}
	for _, m := range methods {
		routes = append(routes, e.Add(m, path, h, middleware...))
	}
	return routes
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

func (e *Echo) Prefix() string {
	return e.prefix
}

func (e *Echo) SetPrefix(prefix string) *Echo {
	if len(prefix) == 0 {
		return e
	}
	if !strings.HasPrefix(prefix, `/`) {
		prefix = `/` + prefix
	}
	prefix = strings.TrimSuffix(prefix, `/`)
	e.prefix = prefix
	return e
}

func (e *Echo) add(host, method, prefix string, path string, h interface{}, middleware ...interface{}) *Route {
	r := &Route{
		Host:       host,
		Method:     method,
		Path:       e.prefix + path,
		Prefix:     prefix,
		handler:    h,
		middleware: middleware,
	}
	e.router.routes = append(e.router.routes, r)
	return r
}

func (e *Echo) buildRouter() *Echo {
	return e.RebuildRouter()
}

// Add registers a new route for an HTTP method and path with matching handler
// in the router with optional route-level middleware.
func (e *Echo) Add(method, path string, handler interface{}, middleware ...interface{}) *Route {
	return e.add("", method, "", path, handler, middleware...)
}

// MetaHandler Add meta information about endpoint
func (e *Echo) MetaHandler(m H, handler interface{}, requests ...interface{}) Handler {
	var request interface{}
	if len(requests) > 0 {
		request = requests[0]
	}
	return e.MetaHandlerWithRequest(m, handler, request)
}

// MetaHandlerWithRequest Add meta information about endpoint
func (e *Echo) MetaHandlerWithRequest(m H, handler interface{}, request interface{}, methods ...string) Handler {
	h := &MetaHandler{
		meta:    m,
		Handler: e.ValidHandler(handler),
	}
	if request != nil {
		switch r := request.(type) {
		case RequestValidator:
			h.request = r
		case func() MetaValidator:
			h.request = r
		case func() interface{}:
			h.request = func() MetaValidator {
				return NewBaseRequestValidator(r())
			}
		default:
			t := reflect.Indirect(reflect.ValueOf(r)).Type()
			if t.Kind() != reflect.Struct {
				panic(fmt.Sprintf(`unsupported validate data: %T`, r))
			}
			var method []string
			for _, me := range methods {
				method = append(method, splitHTTPMethod.Split(me, -1)...)
			}
			h.request = func() MetaValidator {
				return NewBaseRequestValidator(reflect.New(t).Interface(), method...)
			}
		}
	}
	return h
}

// RebuildRouter rebuild router
func (e *Echo) RebuildRouter(args ...[]*Route) *Echo {
	routes := e.router.routes
	if len(args) > 0 {
		routes = args[0]
	}
	e.router = NewRouter(e)
	for i, r := range routes {
		router, _, _, _ := e.findRouter(r.Host)
		r.apply(e)
		router.Add(r, i)
		if e.RouteDebug {
			e.logger.Debugf(`Route: %7v %-30v -> %v`, r.Method, r.Host+r.Format, r.Name)
		}

		if _, ok := e.router.nroute[r.Name]; !ok {
			e.router.nroute[r.Name] = []int{i}
		} else {
			e.router.nroute[r.Name] = append(e.router.nroute[r.Name], i)
		}
	}
	e.router.routes = routes
	return e
}

// AppendRouter append router
func (e *Echo) AppendRouter(routes []*Route) *Echo {
	for i, r := range routes {
		router, _, _, _ := e.findRouter(r.Host)
		i = len(e.router.routes)
		r.apply(e)
		router.Add(r, i)
		if _, ok := e.router.nroute[r.Name]; !ok {
			e.router.nroute[r.Name] = []int{i}
		} else {
			e.router.nroute[r.Name] = append(e.router.nroute[r.Name], i)
		}
		e.router.routes = append(e.router.routes, r)
	}
	return e
}

func parseHostConfig(name string) *host {
	vhost := &host{
		name: name,
	}
	return vhost.Parse()
}

// Host creates a new router group for the provided host and optional host-level middleware.
func (e *Echo) Host(name string, m ...interface{}) *Group {
	h, y := e.hosts[name]
	if !y {
		h = &Host{
			group:  &Group{host: parseHostConfig(name), echo: e},
			groups: map[string]*Group{},
			Router: NewRouter(e),
		}
		e.hosts[name] = h
	}
	if len(m) > 0 {
		h.group.Use(m...)
	}
	return h.group
}

// TypeHost TypeHost(`blog`).URI(`login`)
func (e *Echo) TypeHost(alias string, args ...interface{}) (r TypeHost) {
	if name, ok := e.hostAlias[alias]; ok {
		hs, ok := e.hosts[name]
		if !ok || hs == nil {
			return
		}
		r = hs.Host(args...)
	}
	return
}

func (e *Echo) OnHostFound(onHostFound func(Context) (bool, error)) *Echo {
	e.onHostFound = onHostFound
	return e
}

// Group creates a new sub-router with prefix.
func (e *Echo) Group(prefix string, m ...interface{}) *Group {
	g, y := e.groups[prefix]
	if !y {
		g = &Group{prefix: prefix, echo: e}
		e.groups[prefix] = g
	}
	if len(m) > 0 {
		g.Use(m...)
	}
	return g
}

// URI generates a URI from handler.
func (e *Echo) URI(handler interface{}, params ...interface{}) string {
	var uri, name string
	switch h := handler.(type) {
	case Handler:
		if hn, ok := h.(Name); ok {
			name = hn.Name()
		} else {
			name = HandlerName(h)
		}
	case string:
		name = h
	default:
		return uri
	}
	if indexes, ok := e.router.nroute[name]; ok && len(indexes) > 0 {
		r := e.router.routes[indexes[0]]
		uri = r.MakeURI(e, params...)
	}
	return uri
}

// URL is an alias for `URI` function.
func (e *Echo) URL(h interface{}, params ...interface{}) string {
	return e.URI(h, params...)
}

func (e *Echo) SetRewriter(r Rewriter) {
	e.rewriter = r
}

func (e *Echo) Rewriter() Rewriter {
	return e.rewriter
}

func (e *Echo) wrapURI(uri string) string {
	if e.rewriter != nil {
		uri = e.rewriter.Rewrite(uri)
	}
	if len(e.defaultExtension) > 0 && !strings.HasSuffix(uri, e.defaultExtension) {
		uri += e.defaultExtension
	}
	return uri
}

// Routes returns the registered routes.
func (e *Echo) Routes() []*Route {
	return e.router.routes
}

// NamedRoutes returns the registered handler name.
func (e *Echo) NamedRoutes() map[string][]int {
	return e.router.nroute
}

func (e *Echo) applyMiddleware(h Handler, middleware ...interface{}) Handler {
	for i := len(middleware) - 1; i >= 0; i-- {
		h = e.ValidMiddleware(middleware[i]).Handle(h)
	}
	return h
}

func (e *Echo) buildHandler(c Context) Handler {
	if r, names, values, exist := e.findRouter(c.Host()); exist {
		if len(names) > 0 {
			c.SetHostParamNames(names...)
			c.SetHostParamValues(values...)
		}
		found, err := c.FireHostFound()
		if err != nil {
			return e.applyMiddleware(HandlerFunc(func(Context) error {
				return err
			}), e.middleware...)
		}
		if !found {
			return e.applyMiddleware(e.router.Handle(c), e.middleware...)
		}
		return e.applyMiddleware(r.Handle(c), e.middleware...)
	}
	return e.applyMiddleware(e.router.Handle(c), e.middleware...)
}

func (e *Echo) ServeHTTP(req engine.Request, res engine.Response) {
	c := e.pool.Get().(Context)
	c.Reset(req, res)

	var h Handler
	if len(e.premiddleware) > 0 {
		h = e.applyMiddleware(HandlerFunc(func(c Context) error {
			return e.buildHandler(c).Handle(c)
		}), e.premiddleware...)
	} else {
		h = e.buildHandler(c)
	}
	if err := h.Handle(c); err != nil {
		c.Error(err)
	}

	e.pool.Put(c)
}

// Run starts the HTTP engine.
func (e *Echo) Run(eng engine.Engine, handler ...engine.Handler) error {
	err := e.buildRouter().setEngine(eng).start(handler...)
	if err != nil {
		fmt.Println(err)
	}
	return err
}

func (e *Echo) Commit() *Echo {
	e.buildRouter()
	return e
}

func (e *Echo) start(handler ...engine.Handler) error {
	if len(handler) > 0 {
		e.engine.SetHandler(handler[0])
	} else {
		e.engine.SetHandler(e)
	}
	e.engine.SetLogger(e.logger)
	if e.Debug() {
		e.logger.Debug("running in debug mode")
	}
	return e.engine.Start()
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

func (e *Echo) Shutdown(ctx context.Context) error {
	if e.engine == nil {
		return nil
	}
	return e.engine.Shutdown(ctx)
}

func (e *Echo) findRouter(host string) (*Router, []string, []string, bool) {
	if len(e.hosts) == 0 {
		return e.router, nil, nil, false
	}
	if r, ok := e.hosts[host]; ok {
		return r.Router, nil, nil, true
	}
	l := len(host)
	for h, r := range e.hosts {
		if r.group != nil && r.group.host != nil {
			values, hasExpr := r.group.host.Match(host)
			if hasExpr {
				if len(values) > 0 {
					return r.Router, r.group.host.names, values, true
				}
				continue
			}
		}
		if l <= len(h) {
			continue
		}
		if h[0] == '.' && strings.HasSuffix(host, h) { //.host(xxx.host)
			return r.Router, nil, nil, true
		}
		if h[len(h)-1] == '.' && strings.HasPrefix(host, h) { //host.(host.xxx)
			return r.Router, nil, nil, true
		}
	}
	return e.router, nil, nil, false
}

func (e *Echo) NewContext(req engine.Request, resp engine.Response) Context {
	return NewContext(req, resp, e)
}
