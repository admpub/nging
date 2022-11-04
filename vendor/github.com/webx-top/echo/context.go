package echo

import (
	"context"
	"io"
	"mime/multipart"
	"net/http"
	"time"

	"github.com/admpub/events"

	pkgCode "github.com/webx-top/echo/code"
	"github.com/webx-top/echo/engine"
	"github.com/webx-top/echo/logger"
	"github.com/webx-top/echo/param"
)

// Context represents context for the current request. It holds request and
// response objects, path parameters, data and registered handler.
type Context interface {
	context.Context
	events.Emitterer
	SetEmitterer(events.Emitterer)
	Handler() Handler

	//Transaction
	SetTransaction(t Transaction)
	Transaction() Transaction
	Begin() error
	Rollback() error
	Commit() error
	End(succeed bool) error

	//Standard Context
	StdContext() context.Context
	WithContext(ctx context.Context) *http.Request
	SetValue(key string, value interface{})

	SetValidator(Validator)
	Validator() Validator
	Validate(item interface{}, args ...interface{}) error
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
	Dispatch(route string) Handler

	//----------------
	// Param
	//----------------

	Path() string
	P(int, ...string) string
	Param(string, ...string) string
	// ParamNames returns path parameter names.
	ParamNames() []string
	ParamValues() []string
	SetParamNames(names ...string)
	SetParamValues(values ...string)
	// Host
	HostParamNames() []string
	HostParamValues() []string
	HostParam(string, ...string) string
	HostP(int, ...string) string
	SetHostParamNames(names ...string)
	SetHostParamValues(values ...string)

	// Queries returns the query parameters as map. It is an alias for `engine.URL#Query()`.
	Queries() map[string][]string
	QueryValues(string) []string
	QueryxValues(string) param.StringSlice
	Query(string, ...string) string

	//----------------
	// Form data
	//----------------

	Form(string, ...string) string
	FormValues(string) []string
	FormxValues(string) param.StringSlice
	// Forms returns the form parameters as map. It is an alias for `engine.Request#Form().All()`.
	Forms() map[string][]string

	// Param+
	Px(int, ...string) param.String
	Paramx(string, ...string) param.String
	Queryx(string, ...string) param.String
	Formx(string, ...string) param.String
	// string to param.String
	Atop(string) param.String
	ToParamString(string) param.String
	ToStringSlice([]string) param.StringSlice

	//----------------
	// Context data
	//----------------

	Set(string, interface{})
	Get(string, ...interface{}) interface{}
	Delete(...string)
	Stored() Store
	Internal() *param.SafeMap

	//----------------
	// Bind
	//----------------

	Bind(interface{}, ...FormDataFilter) error
	BindAndValidate(interface{}, ...FormDataFilter) error
	MustBind(interface{}, ...FormDataFilter) error
	MustBindAndValidate(interface{}, ...FormDataFilter) error

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
	Stream(func(io.Writer) bool) error
	SSEvent(string, chan interface{}) error
	File(string, ...http.FileSystem) error
	Attachment(io.Reader, string, time.Time, ...bool) error
	NoContent(...int) error
	Redirect(string, ...int) error
	Error(err error)
	NewError(code pkgCode.Code, msg string, args ...interface{}) *Error
	NewErrorWith(err error, code pkgCode.Code, args ...interface{}) *Error
	SetCode(int)
	Code() int
	SetData(Data)
	Data() Data

	// ServeContent sends static content from `io.Reader` and handles caching
	// via `If-Modified-Since` request header. It automatically sets `Content-Type`
	// and `Last-Modified` response headers.
	ServeContent(io.Reader, string, time.Time) error
	ServeCallbackContent(func(Context) (io.Reader, error), string, time.Time) error

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
	SetAuto(on bool) Context
	Fetch(string, interface{}) ([]byte, error)
	SetRenderer(Renderer)
	SetRenderDataWrapper(DataWrapper)
	Renderer() Renderer
	RenderDataWrapper() DataWrapper

	//----------------
	// Cookie
	//----------------

	SetCookieOptions(*CookieOptions)
	CookieOptions() *CookieOptions
	NewCookie(string, string) *http.Cookie
	Cookie() Cookier
	GetCookie(string) string
	// SetCookie @param:key,value,maxAge(seconds),path(/),domain,secure,httpOnly,sameSite(lax/strict/default)
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
	SetDefaultExtension(string)
	DefaultExtension() string
	ResolveFormat() string
	Accept() *Accepts
	Protocol() string
	Site() string
	RequestURI() string
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
	SaveUploadedFile(fieldName string, saveAbsPath string, saveFileName ...func(*multipart.FileHeader) (string, error)) (*multipart.FileHeader, error)
	SaveUploadedFileToWriter(string, io.Writer) (*multipart.FileHeader, error)
	//Multiple file upload
	SaveUploadedFiles(fieldName string, savePath func(*multipart.FileHeader) (string, error)) error
	SaveUploadedFilesToWriter(fieldName string, writer func(*multipart.FileHeader) (io.Writer, error)) error

	//----------------
	// Hook
	//----------------

	AddPreResponseHook(func() error) Context
	SetPreResponseHook(...func() error) Context
	OnHostFound(func(Context) (bool, error)) Context
	FireHostFound() (bool, error)
}
