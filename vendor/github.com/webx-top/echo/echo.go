package echo

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"regexp"
	"runtime"
	"strings"
	"sync"

	"github.com/admpub/log"
	"github.com/webx-top/echo/engine"
	"github.com/webx-top/echo/logger"
)

type (
	Echo struct {
		engine            engine.Engine
		prefix            string
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
		meta              map[string]H
		handlerWrapper    []func(interface{}) Handler
		middlewareWrapper []func(interface{}) Middleware
		FuncMap           map[string]interface{}
	}

	Route struct {
		Method      string
		Path        string
		Handler     Handler `json:"-" xml:"-"`
		HandlerName string
		Format      string
		Params      []string
		Prefix      string
		Meta        H
	}

	HTTPError struct {
		Code    int
		Message string
	}

	Middleware interface {
		Handle(Handler) Handler
	}

	MiddlewareFunc func(Handler) Handler

	MiddlewareFuncd func(Handler) HandlerFunc

	Handler interface {
		Handle(Context) error
	}

	HandleNamer interface {
		HandleName() string
	}

	HandlerFunc func(Context) error

	// HTTPErrorHandler is a centralized HTTP error handler.
	HTTPErrorHandler func(error, Context)

	// Validator is the interface that wraps the Validate method.
	Validator interface {
		Validate() error
	}

	// Renderer is the interface that wraps the Render method.
	Renderer interface {
		Render(w io.Writer, name string, data interface{}, c Context) error
	}
)

const (
	// CONNECT HTTP method
	CONNECT = "CONNECT"
	// DELETE HTTP method
	DELETE = "DELETE"
	// GET HTTP method
	GET = "GET"
	// HEAD HTTP method
	HEAD = "HEAD"
	// OPTIONS HTTP method
	OPTIONS = "OPTIONS"
	// PATCH HTTP method
	PATCH = "PATCH"
	// POST HTTP method
	POST = "POST"
	// PUT HTTP method
	PUT = "PUT"
	// TRACE HTTP method
	TRACE = "TRACE"

	//-------------
	// Media types
	//-------------

	MIMEApplicationJSON                  = "application/json"
	MIMEApplicationJSONCharsetUTF8       = MIMEApplicationJSON + "; " + CharsetUTF8
	MIMEApplicationJavaScript            = "application/javascript"
	MIMEApplicationJavaScriptCharsetUTF8 = MIMEApplicationJavaScript + "; " + CharsetUTF8
	MIMEApplicationXML                   = "application/xml"
	MIMEApplicationXMLCharsetUTF8        = MIMEApplicationXML + "; " + CharsetUTF8
	MIMEApplicationForm                  = "application/x-www-form-urlencoded"
	MIMEApplicationProtobuf              = "application/protobuf"
	MIMEApplicationMsgpack               = "application/msgpack"
	MIMETextHTML                         = "text/html"
	MIMETextHTMLCharsetUTF8              = MIMETextHTML + "; " + CharsetUTF8
	MIMETextPlain                        = "text/plain"
	MIMETextPlainCharsetUTF8             = MIMETextPlain + "; " + CharsetUTF8
	MIMEMultipartForm                    = "multipart/form-data"
	MIMEOctetStream                      = "application/octet-stream"
	MIMEEventStream                      = "text/event-stream"

	//---------
	// Charset
	//---------

	CharsetUTF8 = "charset=utf-8"

	//---------
	// Headers
	//---------

	HeaderAcceptEncoding                = "Accept-Encoding"
	HeaderAuthorization                 = "Authorization"
	HeaderContentDisposition            = "Content-Disposition"
	HeaderContentEncoding               = "Content-Encoding"
	HeaderContentLength                 = "Content-Length"
	HeaderContentType                   = "Content-Type"
	HeaderIfModifiedSince               = "If-Modified-Since"
	HeaderCookie                        = "Cookie"
	HeaderSetCookie                     = "Set-Cookie"
	HeaderLastModified                  = "Last-Modified"
	HeaderLocation                      = "Location"
	HeaderUpgrade                       = "Upgrade"
	HeaderVary                          = "Vary"
	HeaderWWWAuthenticate               = "WWW-Authenticate"
	HeaderXForwardedProto               = "X-Forwarded-Proto"
	HeaderXHTTPMethodOverride           = "X-HTTP-Method-Override"
	HeaderXForwardedFor                 = "X-Forwarded-For"
	HeaderXRealIP                       = "X-Real-IP"
	HeaderServer                        = "Server"
	HeaderOrigin                        = "Origin"
	HeaderAccessControlRequestMethod    = "Access-Control-Request-Method"
	HeaderAccessControlRequestHeaders   = "Access-Control-Request-Headers"
	HeaderAccessControlAllowOrigin      = "Access-Control-Allow-Origin"
	HeaderAccessControlAllowMethods     = "Access-Control-Allow-Methods"
	HeaderAccessControlAllowHeaders     = "Access-Control-Allow-Headers"
	HeaderAccessControlAllowCredentials = "Access-Control-Allow-Credentials"
	HeaderAccessControlExposeHeaders    = "Access-Control-Expose-Headers"
	HeaderAccessControlMaxAge           = "Access-Control-Max-Age"

	// Security
	HeaderStrictTransportSecurity = "Strict-Transport-Security"
	HeaderXContentTypeOptions     = "X-Content-Type-Options"
	HeaderXXSSProtection          = "X-XSS-Protection"
	HeaderXFrameOptions           = "X-Frame-Options"
	HeaderContentSecurityPolicy   = "Content-Security-Policy"
	HeaderXCSRFToken              = "X-CSRF-Token"
)

var (
	httpMethodRegexp = regexp.MustCompile(`[^A-Z]+`)

	methods = []string{
		CONNECT,
		DELETE,
		GET,
		HEAD,
		OPTIONS,
		PATCH,
		POST,
		PUT,
		TRACE,
	}

	//--------
	// Errors
	//--------

	ErrUnsupportedMediaType        error = NewHTTPError(http.StatusUnsupportedMediaType)
	ErrNotFound                    error = NewHTTPError(http.StatusNotFound)
	ErrUnauthorized                error = NewHTTPError(http.StatusUnauthorized)
	ErrStatusRequestEntityTooLarge error = NewHTTPError(http.StatusRequestEntityTooLarge)
	ErrMethodNotAllowed            error = NewHTTPError(http.StatusMethodNotAllowed)
	ErrRendererNotRegistered             = errors.New("renderer not registered")
	ErrInvalidRedirectCode               = errors.New("invalid redirect status code")

	//----------------
	// Error handlers
	//----------------

	NotFoundHandler = HandlerFunc(func(c Context) error {
		return ErrNotFound
	})

	MethodNotAllowedHandler = HandlerFunc(func(c Context) error {
		return ErrMethodNotAllowed
	})

	_ MiddlewareFuncd = func(h Handler) HandlerFunc {
		return func(c Context) error {
			return h.Handle(c)
		}
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
	e.meta = map[string]H{}

	//----------
	// Defaults
	//----------
	e.SetHTTPErrorHandler(e.DefaultHTTPErrorHandler)
	e.SetBinder(&binder{Echo: e})

	// Logger
	e.logger = log.GetLogger("echo")

	return
}

func (m MiddlewareFunc) Handle(h Handler) Handler {
	return m(h)
}

func (m MiddlewareFunc) SetMeta(e *Echo, meta H) MiddlewareFunc {
	e.meta[HandlerName(m)] = meta
	return m
}

func (m MiddlewareFuncd) Handle(h Handler) Handler {
	return m(h)
}

func (m MiddlewareFuncd) SetMeta(e *Echo, meta H) MiddlewareFuncd {
	e.meta[HandlerName(m)] = meta
	return m
}

func (h HandlerFunc) Handle(c Context) error {
	return h(c)
}

func (h HandlerFunc) SetMeta(e *Echo, meta H) HandlerFunc {
	e.meta[HandlerName(h)] = meta
	return h
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
	}
}

// Pre is alias
func (e *Echo) Pre(middleware ...interface{}) {
	e.PreUse(middleware...)
}

// PreUse adds handler to the middleware chain.
func (e *Echo) PreUse(middleware ...interface{}) {
	var middlewares []interface{}
	for _, m := range middleware {
		e.ValidMiddleware(m)
		middlewares = append(middlewares, m)
	}
	e.middleware = append(middlewares, e.middleware...)
}

// Connect adds a CONNECT route > handler to the router.
func (e *Echo) Connect(path string, h interface{}, m ...interface{}) {
	e.add(CONNECT, path, h, m...)
}

// Delete adds a DELETE route > handler to the router.
func (e *Echo) Delete(path string, h interface{}, m ...interface{}) {
	e.add(DELETE, path, h, m...)
}

// Get adds a GET route > handler to the router.
func (e *Echo) Get(path string, h interface{}, m ...interface{}) {
	e.add(GET, path, h, m...)
}

// Head adds a HEAD route > handler to the router.
func (e *Echo) Head(path string, h interface{}, m ...interface{}) {
	e.add(HEAD, path, h, m...)
}

// Options adds an OPTIONS route > handler to the router.
func (e *Echo) Options(path string, h interface{}, m ...interface{}) {
	e.add(OPTIONS, path, h, m...)
}

// Patch adds a PATCH route > handler to the router.
func (e *Echo) Patch(path string, h interface{}, m ...interface{}) {
	e.add(PATCH, path, h, m...)
}

// Post adds a POST route > handler to the router.
func (e *Echo) Post(path string, h interface{}, m ...interface{}) {
	e.add(POST, path, h, m...)
}

// Put adds a PUT route > handler to the router.
func (e *Echo) Put(path string, h interface{}, m ...interface{}) {
	e.add(PUT, path, h, m...)
}

// Trace adds a TRACE route > handler to the router.
func (e *Echo) Trace(path string, h interface{}, m ...interface{}) {
	e.add(TRACE, path, h, m...)
}

// Any adds a route > handler to the router for all HTTP methods.
func (e *Echo) Any(path string, h interface{}, middleware ...interface{}) {
	for _, m := range methods {
		e.add(m, path, h, middleware...)
	}
}

func (e *Echo) Route(methods string, path string, h interface{}, middleware ...interface{}) {
	e.Match(httpMethodRegexp.Split(methods, -1), path, h, middleware...)
}

// Match adds a route > handler to the router for multiple HTTP methods provided.
func (e *Echo) Match(methods []string, path string, h interface{}, middleware ...interface{}) {
	for _, m := range methods {
		e.add(m, path, h, middleware...)
	}
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

func (e *Echo) add(method, path string, h interface{}, middleware ...interface{}) {
	handler := e.ValidHandler(h)
	if handler == nil {
		return
	}
	var name string
	if hn, ok := handler.(HandleNamer); ok {
		name = hn.HandleName()
	} else {
		name = HandlerName(handler)
	}
	meta := H{}
	for i := len(middleware) - 1; i >= 0; i-- {
		m := middleware[i]
		e.addMeta(meta, HandlerName(m))
		mw := e.ValidMiddleware(m)
		handler = mw.Handle(handler)
	}
	e.addMeta(meta, name)
	hdl := HandlerFunc(func(c Context) error {
		return handler.Handle(c)
	})
	fpath, pnames := e.router.Add(method, path, hdl, e)
	e.logger.Debugf(`Route: %7v %-30v -> %v`, method, fpath, name)
	r := &Route{
		Method:      method,
		Path:        path,
		Handler:     hdl,
		HandlerName: name,
		Format:      fpath,
		Params:      pnames,
		Meta:        meta,
	}
	if _, ok := e.router.nroute[name]; !ok {
		e.router.nroute[name] = []int{len(e.router.routes)}
	} else {
		e.router.nroute[name] = append(e.router.nroute[name], len(e.router.routes))
	}
	e.router.routes = append(e.router.routes, r)
}

func (e *Echo) addMeta(meta H, handler string) {
	if m, ok := e.meta[handler]; ok {
		meta.DeepMerge(m)
	}
}

// Add meta information about endpoint
func (e *Echo) MetaMiddleware(m H, middleware interface{}) interface{} {
	name := HandlerName(middleware)
	e.meta[name] = m
	return middleware
}

// Add meta information about endpoint
func (e *Echo) MetaHandler(m H, handler interface{}) interface{} {
	name := HandlerName(handler)
	e.meta[name] = m
	return handler
}

// RebuildRouter rebuild router
func (e *Echo) RebuildRouter(args ...[]*Route) {
	routes := e.router.routes
	if len(args) > 0 {
		routes = args[0]
	}
	e.router = NewRouter(e)
	for index, r := range routes {
		//e.logger.Debugf(`%p rebuild: %#v`, e, *r)
		e.router.Add(r.Method, r.Path, r.Handler, e)

		if _, ok := e.router.nroute[r.HandlerName]; !ok {
			e.router.nroute[r.HandlerName] = []int{index}
		} else {
			e.router.nroute[r.HandlerName] = append(e.router.nroute[r.HandlerName], index)
		}
	}
	e.router.routes = routes
	e.head = nil
}

// AppendRouter append router
func (e *Echo) AppendRouter(routes []*Route) {
	for index, r := range routes {
		e.router.Add(r.Method, r.Path, r.Handler, e)
		index = len(e.router.routes)
		if _, ok := e.router.nroute[r.HandlerName]; !ok {
			e.router.nroute[r.HandlerName] = []int{index}
		} else {
			e.router.nroute[r.HandlerName] = append(e.router.nroute[r.HandlerName], index)
		}
		e.router.routes = append(e.router.routes, r)
	}
	e.head = nil
}

// Group creates a new sub-router with prefix.
func (e *Echo) Group(prefix string, m ...interface{}) (g *Group) {
	g = &Group{prefix: prefix, echo: e}
	g.Use(m...)
	return
}

// URI generates a URI from handler.
func (e *Echo) URI(handler interface{}, params ...interface{}) string {
	var uri, name string
	if h, ok := handler.(Handler); ok {
		if hn, ok := h.(HandleNamer); ok {
			name = hn.HandleName()
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
			switch params[0].(type) {
			case url.Values:
				val := params[0].(url.Values)
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
				if q != `` {
					uri += `?` + q
				}
			case map[string]string:
				val := params[0].(map[string]string)
				uri = r.Path
				for _, name := range r.Params {
					tag := `:` + name
					v, _ := val[name]
					uri = strings.Replace(uri, tag+`/`, v+`/`, -1)
					if strings.HasSuffix(uri, tag) {
						uri = strings.TrimSuffix(uri, tag) + v
					}
				}
			case []interface{}:
				val := params[0].([]interface{})
				uri = fmt.Sprintf(r.Format, val...)
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
func (e *Echo) Run(eng engine.Engine) {
	e.engine = eng
	eng.SetHandler(e)
	eng.SetLogger(e.logger)
	if e.Debug() {
		e.logger.Debug("running in debug mode")
	}
	eng.Start()
}

// Stop stops the HTTP server.
func (e *Echo) Stop() error {
	if e.engine == nil {
		return nil
	}
	return e.engine.Stop()
}

func NewHTTPError(code int, msg ...string) *HTTPError {
	he := &HTTPError{Code: code, Message: http.StatusText(code)}
	if len(msg) > 0 {
		he.Message = msg[0]
	}
	return he
}

// Error returns message.
func (e *HTTPError) Error() string {
	return e.Message
}

func HandlerName(h interface{}) string {
	v := reflect.ValueOf(h)
	t := v.Type()
	if t.Kind() == reflect.Func {
		return runtime.FuncForPC(v.Pointer()).Name()
	}
	return t.String()
}

func Methods() []string {
	return methods
}
