package echo

import (
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/webx-top/echo/engine"
	"github.com/webx-top/echo/param"
)

// Request returns *http.Request.
func (c *xContext) Request() engine.Request {
	return c.request
}

// Path returns the registered path for the handler.
func (c *xContext) Path() string {
	return c.path
}

// P returns path parameter by index.
func (c *xContext) P(i int, defaults ...string) (value string) {
	l := len(c.pvalues)
	if i < l {
		value = c.pvalues[i]
	}
	if len(value) == 0 && len(defaults) > 0 {
		return defaults[0]
	}
	return
}

func (c *xContext) Px(n int, defaults ...string) param.String {
	return param.String(c.P(n, defaults...))
}

// Param returns path parameter by name.
func (c *xContext) Param(name string, defaults ...string) (value string) {
	l := len(c.pvalues)
	for i, n := range c.pnames {
		if i < l && n == name {
			value = c.pvalues[i]
			break
		}
	}

	if len(value) == 0 && len(defaults) > 0 {
		return defaults[0]
	}
	return
}

func (c *xContext) Paramx(name string, defaults ...string) param.String {
	return param.String(c.Param(name, defaults...))
}

func (c *xContext) ParamNames() []string {
	return c.pnames
}

func (c *xContext) ParamValues() []string {
	return c.pvalues[:len(c.pvalues)]
}

func (c *xContext) SetParamNames(names ...string) {
	c.pnames = names

	l := len(names)
	if *c.echo.maxParam < l {
		*c.echo.maxParam = l
	}

	if len(c.pvalues) < l {
		// Keeping the old pvalues just for backward compatibility, but it sounds that doesn't make sense to keep them,
		// probably those values will be overriden in a Context#SetParamValues
		newPvalues := make([]string, l)
		copy(newPvalues, c.pvalues)
		c.pvalues = newPvalues
	}
}

func (c *xContext) SetParamValues(values ...string) {
	// NOTE: Don't just set c.pvalues = values, because it has to have length c.echo.maxParam at all times
	// It will brake the Router#Find code
	limit := len(values)
	if limit > *c.echo.maxParam {
		limit = *c.echo.maxParam
	}
	for i := 0; i < limit; i++ {
		c.pvalues[i] = values[i]
	}
}

func (c *xContext) AddHostParam(name string, value string) {
	c.hnames = append(c.hnames, name)
	c.hvalues = append(c.hvalues, value)
}

func (c *xContext) SetHostParamNames(names ...string) {
	c.hnames = names
}

func (c *xContext) SetHostParamValues(values ...string) {
	c.hvalues = values
}

func (c *xContext) HostParamNames() []string {
	return c.hnames
}

func (c *xContext) HostParamValues() []string {
	return c.hvalues[:len(c.hvalues)]
}

// HostP returns host parameter by index.
func (c *xContext) HostP(i int, defaults ...string) (value string) {
	l := len(c.hvalues)
	if i < l {
		value = c.hvalues[i]
	}
	if len(value) == 0 && len(defaults) > 0 {
		return defaults[0]
	}
	return
}

// HostParam returns host parameter by name.
func (c *xContext) HostParam(name string, defaults ...string) (value string) {
	l := len(c.hvalues)
	for i, n := range c.hnames {
		if i < l && n == name {
			value = c.hvalues[i]
			break
		}
	}

	if len(value) == 0 && len(defaults) > 0 {
		return defaults[0]
	}
	return
}

// Query returns query parameter by name.
func (c *xContext) Query(name string, defaults ...string) (value string) {
	value = c.request.URL().QueryValue(name)
	if len(value) == 0 && len(defaults) > 0 {
		return defaults[0]
	}
	return
}

func (c *xContext) Queryx(name string, defaults ...string) param.String {
	return param.String(c.Query(name, defaults...))
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
func (c *xContext) Form(name string, defaults ...string) (value string) {
	value = c.request.FormValue(name)
	if len(value) == 0 && len(defaults) > 0 {
		return defaults[0]
	}
	return
}

func (c *xContext) Formx(name string, defaults ...string) param.String {
	return param.String(c.Form(name, defaults...))
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

// Bind binds the request body into specified type `i`. The default binder does
// it based on Content-Type header.
func (c *xContext) Bind(i interface{}, filter ...FormDataFilter) error {
	return c.echo.binder.Bind(i, c, filter...)
}

func (c *xContext) BindAndValidate(i interface{}, filter ...FormDataFilter) error {
	return c.echo.binder.BindAndValidate(i, c, filter...)
}

func (c *xContext) MustBind(i interface{}, filter ...FormDataFilter) error {
	return c.echo.binder.MustBind(i, c, filter...)
}

func (c *xContext) MustBindAndValidate(i interface{}, filter ...FormDataFilter) error {
	return c.echo.binder.MustBindAndValidate(i, c, filter...)
}

func (c *xContext) Header(name string) string {
	return c.Request().Header().Get(name)
}

func (c *xContext) IsAjax() bool {
	return c.Header(`X-Requested-With`) == `XMLHttpRequest`
}

func (c *xContext) IsPjax() bool {
	return len(c.Header(`X-PJAX`)) > 0 || len(c.PjaxContainer()) > 0
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

	info := c.Accept()
	for _, accepts := range info.Accepts {
		for _, accept := range accepts.Type {
			if format, ok := c.echo.acceptFormats[accept.Mime]; ok {
				return format
			}
		}
	}
	if format, ok := c.echo.acceptFormats[`*`]; ok {
		return format
	}
	return `html`
}

func (c *xContext) Accept() *Accepts {
	if c.accept != nil {
		return c.accept
	}
	c.accept = NewAccepts(c.Header(HeaderAccept))
	if c.echo.parseHeaderAccept {
		return c.accept.Advance()
	}
	return c.accept.Simple(3)
}

// Protocol returns request protocol name, such as HTTP/1.1 .
func (c *xContext) Protocol() string {
	return c.Request().Proto()
}

// Site returns base site url as scheme://domain/ type.
func (c *xContext) Site() string {
	return c.Scheme() + `://` + c.Request().Host() + `/`
}

func (c *xContext) RequestURI() string {
	return c.Scheme() + `://` + c.Request().Host() + c.Request().URL().String()
}

// Scheme returns request scheme as `http` or `https`.
func (c *xContext) Scheme() string {
	scheme := c.Header(HeaderXForwardedProto)
	if len(scheme) > 0 {
		return scheme
	}
	return c.Request().Scheme()
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
		var isIPv6 bool
		if host[0] == '[' {
			isIPv6 = true
			delim = `]:`
		}
		hostParts := strings.SplitN(host, delim, 2)
		if len(hostParts) == 2 {
			host = hostParts[0]
			if isIPv6 {
				host += `]`
			}
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

func (c *xContext) SaveUploadedFile(fieldName string, saveAbsPath string, saveFileName ...func(*multipart.FileHeader) (string, error)) (*multipart.FileHeader, error) {
	fileSrc, fileHdr, err := c.Request().FormFile(fieldName)
	if err != nil {
		return fileHdr, err
	}
	defer fileSrc.Close()

	// Destination
	fileName := fileHdr.Filename
	if len(saveFileName) > 0 && saveFileName[0] != nil {
		fileName, err = saveFileName[0](fileHdr)
		if err != nil {
			return fileHdr, err
		}
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

func (c *xContext) SaveUploadedFiles(fieldName string, savePath func(*multipart.FileHeader) (string, error)) error {
	m := c.Request().MultipartForm()
	files, ok := m.File[fieldName]
	if !ok {
		return ErrNotFoundFileInput
	}
	var dstFile string
	for _, fileHdr := range files {
		//for each fileheader, get a handle to the actual file
		file, err := fileHdr.Open()
		if err != nil {
			file.Close()
			return err
		}
		dstFile, err = savePath(fileHdr)
		if err != nil {
			file.Close()
			return err
		}
		if len(dstFile) == 0 {
			file.Close()
			continue
		}
		//create destination file making sure the path is writeable.
		dst, err := os.Create(dstFile)
		if err != nil {
			file.Close()
			return err
		}
		//copy the uploaded file to the destination file
		if _, err := io.Copy(dst, file); err != nil {
			file.Close()
			dst.Close()
			return err
		}
		file.Close()
		dst.Close()
	}
	return nil
}

func (c *xContext) SaveUploadedFilesToWriter(fieldName string, writer func(*multipart.FileHeader) (io.Writer, error)) error {
	m := c.Request().MultipartForm()
	files, ok := m.File[fieldName]
	if !ok {
		return ErrNotFoundFileInput
	}
	var w io.Writer
	for _, fileHdr := range files {
		//for each fileheader, get a handle to the actual file
		file, err := fileHdr.Open()
		if err != nil {
			file.Close()
			return err
		}
		w, err = writer(fileHdr)
		if err != nil {
			file.Close()
			return err
		}
		if w == nil {
			continue
		}
		//copy the uploaded file to the destination file
		if _, err := io.Copy(w, file); err != nil {
			file.Close()
			if v, ok := w.(Closer); ok {
				v.Close()
			}
			return err
		}
		file.Close()
		if v, ok := w.(Closer); ok {
			v.Close()
		}
	}
	return nil
}

// HasAnyRequest 是否提交了参数
func (c *xContext) HasAnyRequest() bool {
	return len(c.Request().Form().All()) > 0
}
