package echo

import (
	"bytes"
	"context"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/admpub/events"
	"github.com/admpub/events/emitter"
	"github.com/webx-top/echo/engine"
	"github.com/webx-top/echo/logger"
	"github.com/webx-top/echo/param"
)

type (
	// Context represents context for the current request. It holds request and
	// response objects, path parameters, data and registered handler.
	Context interface {
		context.Context
		events.Emitter
		SetEmitter(events.Emitter)
		StdContext() context.Context
		SetStdContext(context.Context)
		Validator
		SetValidator(Validator)
		Translator
		SetTranslator(Translator)
		Request() engine.Request
		Response() engine.Response
		Handle(Context) error
		Logger() logger.Logger
		Object() *xContext
		Echo() *Echo
		Route() *Route
		Reset(engine.Request, engine.Response)

		//----------------
		// Param
		//----------------

		Path() string
		P(int) string
		Param(string) string
		// ParamNames returns path parameter names.
		ParamNames() []string
		ParamValues() []string
		SetParamValues(values ...string)

		// Queries returns the query parameters as map. It is an alias for `engine.URL#Query()`.
		Queries() map[string][]string
		QueryValues(string) []string
		QueryxValues(string) param.StringSlice
		Query(string) string

		//----------------
		// Form data
		//----------------

		Form(string) string
		FormValues(string) []string
		FormxValues(string) param.StringSlice
		// Forms returns the form parameters as map. It is an alias for `engine.Request#Form().All()`.
		Forms() map[string][]string

		// Param+
		Px(int) param.String
		Paramx(string) param.String
		Queryx(string) param.String
		Formx(string) param.String
		// string to param.String
		Atop(string) param.String
		ToParamString(string) param.String
		ToStringSlice([]string) param.StringSlice

		//----------------
		// Context data
		//----------------

		Set(string, interface{})
		Get(string) interface{}
		Delete(...string)
		Stored() store

		//----------------
		// Bind
		//----------------

		Bind(interface{}, ...FormDataFilter) error
		MustBind(interface{}, ...FormDataFilter) error

		//----------------
		// Response data
		//----------------

		Render(string, interface{}, ...int) error
		HTML(string, ...int) error
		String(string, ...int) error
		Blob([]byte, ...int) error
		JSON(interface{}, ...int) error
		JSONBlob([]byte, ...int) error
		JSONP(string, interface{}, ...int) error
		XML(interface{}, ...int) error
		XMLBlob([]byte, ...int) error
		Stream(func(io.Writer) bool)
		SSEvent(string, chan interface{}) error
		File(string, ...http.FileSystem) error
		Attachment(io.ReadSeeker, string) error
		NoContent(...int) error
		Redirect(string, ...int) error
		Error(err error)
		SetCode(int)
		Code() int
		SetData(Data)
		Data() Data

		// ServeContent sends static content from `io.Reader` and handles caching
		// via `If-Modified-Since` request header. It automatically sets `Content-Type`
		// and `Last-Modified` response headers.
		ServeContent(io.ReadSeeker, string, time.Time) error

		//----------------
		// FuncMap
		//----------------

		SetFunc(string, interface{})
		GetFunc(string) interface{}
		ResetFuncs(map[string]interface{})
		Funcs() map[string]interface{}
		PrintFuncs()

		//----------------
		// Render
		//----------------

		Fetch(string, interface{}) ([]byte, error)
		SetRenderer(Renderer)

		//----------------
		// Cookie
		//----------------

		SetCookieOptions(*CookieOptions)
		CookieOptions() *CookieOptions
		NewCookie(string, string) *Cookie
		Cookie() Cookier
		GetCookie(string) string
		SetCookie(string, string, ...interface{})

		//----------------
		// Session
		//----------------

		SetSessionOptions(*SessionOptions)
		SessionOptions() *SessionOptions
		SetSessioner(Sessioner)
		Session() Sessioner
		Flash(...string) interface{}

		//----------------
		// Request data
		//----------------

		Header(string) string
		IsAjax() bool
		IsPjax() bool
		PjaxContainer() string
		Method() string
		Format() string
		SetFormat(string)
		IsPost() bool
		IsGet() bool
		IsPut() bool
		IsDel() bool
		IsHead() bool
		IsPatch() bool
		IsOptions() bool
		IsSecure() bool
		IsWebsocket() bool
		IsUpload() bool
		ResolveContentType() string
		WithFormatExtension(bool)
		ResolveFormat() string
		Protocol() string
		Site() string
		Scheme() string
		Domain() string
		Host() string
		Proxy() []string
		Referer() string
		Port() int
		RealIP() string
		HasAnyRequest() bool

		MapForm(i interface{}, names ...string) error
		MapData(i interface{}, data map[string][]string, names ...string) error
		SaveUploadedFile(fieldName string, saveAbsPath string, saveFileName ...string) (*multipart.FileHeader, error)
		SaveUploadedFileToWriter(string, io.Writer) (*multipart.FileHeader, error)
		//Multiple file upload
		SaveUploadedFiles(fieldName string, savePath func(*multipart.FileHeader) string) error
		SaveUploadedFilesToWriter(fieldName string, writer func(*multipart.FileHeader) io.Writer) error

		//----------------
		// Hook
		//----------------

		AddPreResponseHook(func() error) Context
		SetPreResponseHook(...func() error) Context
	}

	xContext struct {
		Validator
		Translator
		events.Emitter
		sessioner           Sessioner
		cookier             Cookier
		context             context.Context
		request             engine.Request
		response            engine.Response
		path                string
		pnames              []string
		pvalues             []string
		store               store
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
	}
)

// NewContext creates a Context object.
func NewContext(req engine.Request, res engine.Response, e *Echo) Context {
	c := &xContext{
		Validator:  DefaultNopValidate,
		Translator: DefaultNopTranslate,
		Emitter:    emitter.DefaultCondEmitter,
		context:    context.Background(),
		request:    req,
		response:   res,
		echo:       e,
		pvalues:    make([]string, *e.maxParam),
		store:      make(store),
		handler:    NotFoundHandler,
		funcs:      make(map[string]interface{}),
		sessioner:  DefaultSession,
	}
	c.cookier = NewCookier(c)
	c.dataEngine = NewData(c)
	return c
}

func (c *xContext) StdContext() context.Context {
	return c.context
}

func (c *xContext) SetStdContext(ctx context.Context) {
	c.context = ctx
}

func (c *xContext) SetEmitter(emitter events.Emitter) {
	c.Emitter = emitter
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

// Request returns *http.Request.
func (c *xContext) Request() engine.Request {
	return c.request
}

// Response returns *Response.
func (c *xContext) Response() engine.Response {
	return c.response
}

// Path returns the registered path for the handler.
func (c *xContext) Path() string {
	return c.path
}

// P returns path parameter by index.
func (c *xContext) P(i int) (value string) {
	l := len(c.pnames)
	if i < l {
		value = c.pvalues[i]
	}
	return
}

// Param returns path parameter by name.
func (c *xContext) Param(name string) (value string) {
	l := len(c.pnames)
	for i, n := range c.pnames {
		if n == name && i < l {
			value = c.pvalues[i]
			break
		}
	}
	return
}

func (c *xContext) ParamNames() []string {
	return c.pnames
}

func (c *xContext) ParamValues() []string {
	return c.pvalues
}

func (c *xContext) SetParamValues(values ...string) {
	c.pvalues = values
}

// Query returns query parameter by name.
func (c *xContext) Query(name string) string {
	return c.request.URL().QueryValue(name)
}

func (c *xContext) QueryValues(name string) []string {
	return c.request.URL().QueryValues(name)
}

func (c *xContext) QueryxValues(name string) param.StringSlice {
	return param.StringSlice(c.request.URL().QueryValues(name))
}

func (c *xContext) Queries() map[string][]string {
	return c.request.URL().Query()
}

// Form returns form parameter by name.
func (c *xContext) Form(name string) string {
	return c.request.FormValue(name)
}

func (c *xContext) FormValues(name string) []string {
	return c.request.Form().Gets(name)
}

func (c *xContext) FormxValues(name string) param.StringSlice {
	return param.StringSlice(c.request.Form().Gets(name))
}

func (c *xContext) Forms() map[string][]string {
	return c.request.Form().All()
}

// Get retrieves data from the context.
func (c *xContext) Get(key string) interface{} {
	return c.store.Get(key)
}

// Set saves data in the context.
func (c *xContext) Set(key string, val interface{}) {
	c.store.Set(key, val)
}

// Delete saves data in the context.
func (c *xContext) Delete(keys ...string) {
	c.store.Delete(keys...)
}

func (c *xContext) Stored() store {
	return c.store
}

// Bind binds the request body into specified type `i`. The default binder does
// it based on Content-Type header.
func (c *xContext) Bind(i interface{}, filter ...FormDataFilter) error {
	return c.echo.binder.Bind(i, c, filter...)
}

func (c *xContext) MustBind(i interface{}, filter ...FormDataFilter) error {
	return c.echo.binder.MustBind(i, c, filter...)
}

// Render renders a template with data and sends a text/html response with status
// code. Templates can be registered using `Echo.SetRenderer()`.
func (c *xContext) Render(name string, data interface{}, codes ...int) (err error) {
	b, err := c.Fetch(name, data)
	if err != nil {
		return
	}
	b = bytes.TrimLeftFunc(b, unicode.IsSpace)
	c.response.Header().Set(HeaderContentType, MIMETextHTMLCharsetUTF8)
	err = c.Blob(b, codes...)
	return
}

// HTML sends an HTTP response with status code.
func (c *xContext) HTML(html string, codes ...int) (err error) {
	c.response.Header().Set(HeaderContentType, MIMETextHTMLCharsetUTF8)
	err = c.Blob([]byte(html), codes...)
	return
}

// String sends a string response with status code.
func (c *xContext) String(s string, codes ...int) (err error) {
	c.response.Header().Set(HeaderContentType, MIMETextPlainCharsetUTF8)
	err = c.Blob([]byte(s), codes...)
	return
}

func (c *xContext) Blob(b []byte, codes ...int) (err error) {
	if len(codes) > 0 {
		c.code = codes[0]
	}
	if c.code == 0 {
		c.code = http.StatusOK
	}
	err = c.preResponse()
	if err != nil {
		return
	}
	c.response.WriteHeader(c.code)
	_, err = c.response.Write(b)
	return
}

// JSON sends a JSON response with status code.
func (c *xContext) JSON(i interface{}, codes ...int) (err error) {
	var b []byte
	if c.echo.Debug() {
		b, err = json.MarshalIndent(i, "", "  ")
	} else {
		b, err = json.Marshal(i)
	}
	if err != nil {
		return err
	}
	return c.JSONBlob(b, codes...)
}

// JSONBlob sends a JSON blob response with status code.
func (c *xContext) JSONBlob(b []byte, codes ...int) (err error) {
	c.response.Header().Set(HeaderContentType, MIMEApplicationJSONCharsetUTF8)
	err = c.Blob(b, codes...)
	return
}

// JSONP sends a JSONP response with status code. It uses `callback` to construct
// the JSONP payload.
func (c *xContext) JSONP(callback string, i interface{}, codes ...int) (err error) {
	b, err := json.Marshal(i)
	if err != nil {
		return err
	}
	c.response.Header().Set(HeaderContentType, MIMEApplicationJavaScriptCharsetUTF8)
	b = []byte(callback + "(" + string(b) + ");")
	err = c.Blob(b, codes...)
	return
}

// XML sends an XML response with status code.
func (c *xContext) XML(i interface{}, codes ...int) (err error) {
	var b []byte
	if c.echo.Debug() {
		b, err = xml.MarshalIndent(i, "", "  ")
	} else {
		b, err = xml.Marshal(i)
	}
	if err != nil {
		return err
	}
	return c.XMLBlob(b, codes...)
}

// XMLBlob sends a XML blob response with status code.
func (c *xContext) XMLBlob(b []byte, codes ...int) (err error) {
	c.response.Header().Set(HeaderContentType, MIMEApplicationXMLCharsetUTF8)
	b = []byte(xml.Header + string(b))
	err = c.Blob(b, codes...)
	return
}

func (c *xContext) Stream(step func(w io.Writer) bool) {
	c.response.Stream(step)
}

func (c *xContext) SSEvent(event string, data chan interface{}) (err error) {
	hdr := c.response.Header()
	hdr.Set(HeaderContentType, MIMEEventStream)
	hdr.Set(`Cache-Control`, `no-cache`)
	hdr.Set(`Connection`, `keep-alive`)
	c.Stream(func(w io.Writer) bool {
		b, e := c.Fetch(event, <-data)
		if e != nil {
			err = e
			return false
		}
		_, e = w.Write(b)
		if e != nil {
			err = e
			return false
		}
		return true
	})
	return
}

func (c *xContext) File(file string, fs ...http.FileSystem) (err error) {
	var f http.File
	customFS := len(fs) > 0 && fs[0] != nil
	if customFS {
		f, err = fs[0].Open(file)
	} else {
		f, err = os.Open(file)
	}
	if err != nil {
		return ErrNotFound
	}
	defer f.Close()

	fi, _ := f.Stat()
	if fi.IsDir() {
		file = filepath.Join(file, "index.html")
		if customFS {
			f, err = fs[0].Open(file)
		} else {
			f, err = os.Open(file)
		}
		if err != nil {
			return ErrNotFound
		}
		fi, _ = f.Stat()
	}
	return c.ServeContent(f, fi.Name(), fi.ModTime())
}

func (c *xContext) Attachment(r io.ReadSeeker, name string) (err error) {
	c.response.Header().Set(HeaderContentType, ContentTypeByExtension(name))
	c.response.Header().Set(HeaderContentDisposition, "attachment; filename="+name)
	c.response.WriteHeader(http.StatusOK)
	c.response.KeepBody(false)
	_, err = io.Copy(c.response, r)
	return
}

// NoContent sends a response with no body and a status code.
func (c *xContext) NoContent(codes ...int) error {
	if len(codes) > 0 {
		c.code = codes[0]
	}
	if c.code == 0 {
		c.code = http.StatusOK
	}
	c.response.WriteHeader(c.code)
	return nil
}

// Redirect redirects the request with status code.
func (c *xContext) Redirect(url string, codes ...int) error {
	code := http.StatusFound
	if len(codes) > 0 {
		code = codes[0]
	}
	if code < http.StatusMultipleChoices || code > http.StatusTemporaryRedirect {
		return ErrInvalidRedirectCode
	}
	err := c.preResponse()
	if err != nil {
		return err
	}
	c.response.Redirect(url, code)
	return nil
}

// Error invokes the registered HTTP error handler. Generally used by middleware.
func (c *xContext) Error(err error) {
	c.echo.httpErrorHandler(err, c)
}

// Logger returns the `Logger` instance.
func (c *xContext) Logger() logger.Logger {
	return c.echo.logger
}

// Object returns the `context` object.
func (c *xContext) Object() *xContext {
	return c
}

func (c *xContext) ServeContent(content io.ReadSeeker, name string, modtime time.Time) error {
	rq := c.Request()
	rs := c.Response()

	if t, err := time.Parse(http.TimeFormat, rq.Header().Get(HeaderIfModifiedSince)); err == nil && modtime.Before(t.Add(1*time.Second)) {
		rs.Header().Del(HeaderContentType)
		rs.Header().Del(HeaderContentLength)
		return c.NoContent(http.StatusNotModified)
	}

	rs.Header().Set(HeaderContentType, ContentTypeByExtension(name))
	rs.Header().Set(HeaderLastModified, modtime.UTC().Format(http.TimeFormat))
	rs.WriteHeader(http.StatusOK)
	rs.KeepBody(false)
	_, err := io.Copy(rs, content)
	return err
}

// Echo returns the `Echo` instance.
func (c *xContext) Echo() *Echo {
	return c.echo
}

func (c *xContext) SetTranslator(t Translator) {
	c.Translator = t
}

func (c *xContext) Reset(req engine.Request, res engine.Response) {
	c.Validator = DefaultNopValidate
	c.Emitter = emitter.DefaultCondEmitter
	c.Translator = DefaultNopTranslate
	c.sessioner = DefaultSession
	c.cookier = NewCookier(c)
	c.context = context.Background()
	c.request = req
	c.response = res
	c.store = make(store)
	c.path = ""
	c.pnames = nil
	c.funcs = make(map[string]interface{})
	c.renderer = nil
	c.handler = NotFoundHandler
	c.route = nil
	c.rid = -1
	c.sessionOptions = nil
	c.withFormatExtension = false
	c.format = ""
	c.code = 0
	c.preResponseHook = nil
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

func (c *xContext) Session() Sessioner {
	return c.sessioner
}

func (c *xContext) Flash(name ...string) (r interface{}) {
	if v := c.sessioner.Flashes(name...); len(v) > 0 {
		r = v[0]
	}
	return r
}

func (c *xContext) SetCookieOptions(opts *CookieOptions) {
	c.SessionOptions().CookieOptions = opts
}

func (c *xContext) CookieOptions() *CookieOptions {
	return c.SessionOptions().CookieOptions
}

func (c *xContext) SetSessionOptions(opts *SessionOptions) {
	c.sessionOptions = opts
}

func (c *xContext) SessionOptions() *SessionOptions {
	if c.sessionOptions == nil {
		c.sessionOptions = DefaultSessionOptions
	}
	return c.sessionOptions
}

func (c *xContext) NewCookie(key string, value string) *Cookie {
	return NewCookie(key, value, c.CookieOptions())
}

func (c *xContext) Cookie() Cookier {
	return c.cookier
}

func (c *xContext) GetCookie(key string) string {
	return c.cookier.Get(key)
}

func (c *xContext) SetCookie(key string, val string, args ...interface{}) {
	c.cookier.Set(key, val, args...)
}

func (c *xContext) Px(n int) param.String {
	return param.String(c.P(n))
}

func (c *xContext) Paramx(name string) param.String {
	return param.String(c.Param(name))
}

func (c *xContext) Queryx(name string) param.String {
	return param.String(c.Query(name))
}

func (c *xContext) Formx(name string) param.String {
	return param.String(c.Form(name))
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

func (c *xContext) Header(name string) string {
	return c.Request().Header().Get(name)
}

func (c *xContext) IsAjax() bool {
	return c.Header(`X-Requested-With`) == `XMLHttpRequest`
}

func (c *xContext) IsPjax() bool {
	return len(c.Header(`X-PJAX`)) > 0
}

func (c *xContext) PjaxContainer() string {
	container := c.Header(`X-PJAX-Container`)
	if len(container) > 0 {
		return container
	}
	return c.Query(`_pjax`)
}

func (c *xContext) Method() string {
	return c.Request().Method()
}

func (c *xContext) Format() string {
	if len(c.format) == 0 {
		c.format = c.ResolveFormat()
	}
	return c.format
}

func (c *xContext) SetFormat(format string) {
	c.format = format
}

// IsPost CREATE：在服务器新建一个资源
func (c *xContext) IsPost() bool {
	return c.Method() == POST
}

// IsGet SELECT：从服务器取出资源（一项或多项）
func (c *xContext) IsGet() bool {
	return c.Method() == GET
}

// IsPut UPDATE：在服务器更新资源（客户端提供改变后的完整资源）
func (c *xContext) IsPut() bool {
	return c.Method() == PUT
}

// IsDel DELETE：从服务器删除资源
func (c *xContext) IsDel() bool {
	return c.Method() == DELETE
}

// IsHead 获取资源的元数据
func (c *xContext) IsHead() bool {
	return c.Method() == HEAD
}

//IsPatch UPDATE：在服务器更新资源（客户端提供改变的属性）
func (c *xContext) IsPatch() bool {
	return c.Method() == PATCH
}

// IsOptions 获取信息，关于资源的哪些属性是客户端可以改变的
func (c *xContext) IsOptions() bool {
	return c.Method() == OPTIONS
}

func (c *xContext) IsSecure() bool {
	return c.Scheme() == `https`
}

// IsWebsocket returns boolean of this request is in webSocket.
func (c *xContext) IsWebsocket() bool {
	upgrade := c.Header(`Upgrade`)
	return upgrade == `websocket` || upgrade == `Websocket`
}

// IsUpload returns boolean of whether file uploads in this request or not..
func (c *xContext) IsUpload() bool {
	return c.ResolveContentType() == MIMEMultipartForm
}

// ResolveContentType Get the content type.
// e.g. From `multipart/form-data; boundary=--` to `multipart/form-data`
// If none is specified, returns `text/html` by default.
func (c *xContext) ResolveContentType() string {
	contentType := c.Header(HeaderContentType)
	if len(contentType) == 0 {
		return `text/html`
	}
	return strings.ToLower(strings.TrimSpace(strings.SplitN(contentType, `;`, 2)[0]))
}

func (c *xContext) WithFormatExtension(on bool) {
	c.withFormatExtension = on
}

// ResolveFormat maps the request's Accept MIME type declaration to
// a Request.Format attribute, specifically `html`, `xml`, `json`, or `txt`,
// returning a default of `html` when Accept header cannot be mapped to a
// value above.
func (c *xContext) ResolveFormat() string {
	if format := c.Query(`format`); len(format) > 0 {
		return format
	}
	if c.withFormatExtension {
		urlPath := c.Request().URL().Path()
		if pos := strings.LastIndex(urlPath, `.`); pos > -1 {
			return strings.ToLower(urlPath[pos+1:])
		}
	}

	accept := c.Header(HeaderAccept)
	for _, mimeType := range strings.Split(strings.SplitN(accept, `;`, 2)[0], `,`) {
		mimeType = strings.TrimSpace(mimeType)
		if format, ok := c.echo.acceptFormats[mimeType]; ok {
			return format
		}
	}
	if format, ok := c.echo.acceptFormats[`*`]; ok {
		return format
	}
	return `html`
}

// Protocol returns request protocol name, such as HTTP/1.1 .
func (c *xContext) Protocol() string {
	return c.Request().Proto()
}

// Site returns base site url as scheme://domain/ type.
func (c *xContext) Site() string {
	return c.Scheme() + `://` + c.Request().Host() + `/`
}

// Scheme returns request scheme as `http` or `https`.
func (c *xContext) Scheme() string {
	scheme := c.Request().Scheme()
	if len(scheme) > 0 {
		return scheme
	}
	if c.Request().IsTLS() == false {
		return `http`
	}
	return `https`
}

// Domain returns host name.
// Alias of Host method.
func (c *xContext) Domain() string {
	return c.Host()
}

// Host returns host name.
// if no host info in request, return localhost.
func (c *xContext) Host() string {
	host := c.Request().Host()
	if len(host) > 0 {
		delim := `:`
		if host[0] == '[' {
			host = strings.TrimPrefix(host, `[`)
			delim = `]:`
		}
		hostParts := strings.SplitN(host, delim, 2)
		if len(hostParts) > 0 {
			return hostParts[0]
		}
		return host
	}
	return `localhost`
}

// Proxy returns proxy client ips slice.
func (c *xContext) Proxy() []string {
	if ips := c.Header(`X-Forwarded-For`); len(ips) > 0 {
		return strings.Split(ips, `,`)
	}
	return []string{}
}

// Referer returns http referer header.
func (c *xContext) Referer() string {
	return c.Header(`Referer`)
}

func (c *xContext) RealIP() string {
	return c.Request().RealIP()
}

// Port returns request client port.
// when error or empty, return 80.
func (c *xContext) Port() int {
	host := c.Request().Host()
	delim := `:`
	if len(host) > 0 && host[0] == '[' {
		delim = `]:`
	}
	parts := strings.SplitN(host, delim, 2)
	if len(parts) > 1 {
		port, _ := strconv.Atoi(parts[1])
		return port
	}
	return 80
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

// MapForm 映射表单数据到结构体
// ParseStruct mapping forms' name and values to struct's field
// For example:
//		<form>
//			<input name=`user.id`/>
//			<input name=`user.name`/>
//			<input name=`user.age`/>
//		</form>
//
//		type User struct {
//			Id int64
//			Name string
//			Age string
//		}
//
//		var user User
//		err := c.MapForm(&user,`user`)
//
func (c *xContext) MapForm(i interface{}, names ...string) error {
	return c.MapData(i, c.Request().Form().All(), names...)
}

// MapData 映射数据到结构体
func (c *xContext) MapData(i interface{}, data map[string][]string, names ...string) error {
	var name string
	if len(names) > 0 {
		name = names[0]
	}
	return NamedStructMap(c.echo, i, data, name)
}

func (c *xContext) SaveUploadedFile(fieldName string, saveAbsPath string, saveFileName ...string) (*multipart.FileHeader, error) {
	fileSrc, fileHdr, err := c.Request().FormFile(fieldName)
	if err != nil {
		return fileHdr, err
	}
	defer fileSrc.Close()

	// Destination
	fileName := fileHdr.Filename
	if len(saveFileName) > 0 {
		fileName = saveFileName[0]
	}
	fileDst, err := os.Create(filepath.Join(saveAbsPath, fileName))
	if err != nil {
		return fileHdr, err
	}
	defer fileDst.Close()

	// Copy
	if _, err = io.Copy(fileDst, fileSrc); err != nil {
		return fileHdr, err
	}
	return fileHdr, nil
}

func (c *xContext) SaveUploadedFileToWriter(fieldName string, writer io.Writer) (*multipart.FileHeader, error) {
	fileSrc, fileHdr, err := c.Request().FormFile(fieldName)
	if err != nil {
		return fileHdr, err
	}
	defer fileSrc.Close()
	if _, err = io.Copy(writer, fileSrc); err != nil {
		return fileHdr, err
	}
	return fileHdr, nil
}

func (c *xContext) SaveUploadedFiles(fieldName string, savePath func(*multipart.FileHeader) string) error {
	m := c.Request().MultipartForm()
	files := m.File[fieldName]
	for _, fileHdr := range files {
		//for each fileheader, get a handle to the actual file
		file, err := fileHdr.Open()
		defer file.Close()
		if err != nil {
			return err
		}

		//create destination file making sure the path is writeable.
		dst, err := os.Create(savePath(fileHdr))
		defer dst.Close()
		if err != nil {
			return err
		}
		//copy the uploaded file to the destination file
		if _, err := io.Copy(dst, file); err != nil {
			return err
		}
	}
	return nil
}

func (c *xContext) SaveUploadedFilesToWriter(fieldName string, writer func(*multipart.FileHeader) io.Writer) error {
	m := c.Request().MultipartForm()
	files := m.File[fieldName]
	for _, fileHdr := range files {
		//for each fileheader, get a handle to the actual file
		file, err := fileHdr.Open()
		defer file.Close()
		if err != nil {
			return err
		}
		w := writer(fileHdr)
		if v, ok := w.(Closer); ok {
			defer v.Close()
		}
		//copy the uploaded file to the destination file
		if _, err := io.Copy(w, file); err != nil {
			return err
		}
	}
	return nil
}

// HasAnyRequest 是否提交了参数
func (c *xContext) HasAnyRequest() bool {
	return len(c.Request().Form().All()) > 0
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
		return nil
	}
	for _, hook := range c.preResponseHook {
		if err := hook(); err != nil {
			return err
		}
	}
	return nil
}

func (c *xContext) PrintFuncs() {
	for key, fn := range c.Funcs() {
		fmt.Printf("[Template Func](%p) %-15s -> %s \n", fn, key, HandlerName(fn))
	}
}
