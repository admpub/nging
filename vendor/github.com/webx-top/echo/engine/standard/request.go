package standard

import (
	"context"
	"io"
	"mime/multipart"
	"net/http"
	"sync"

	"github.com/admpub/realip"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/engine"
)

type Request struct {
	config    *engine.Config
	requestMu sync.RWMutex
	request   *http.Request
	url       *URL
	header    *Header
	value     *Value
	realIP    string
	maxSize   int
}

func NewRequest(r *http.Request) *Request {
	req := &Request{
		request: r,
		url:     &URL{url: r.URL},
		header:  &Header{header: r.Header},
	}
	req.value = NewValue(req)
	return req
}

func (r *Request) Context() context.Context {
	return r.request.Context()
}

func (r *Request) WithContext(ctx context.Context) *http.Request {
	return r.request.WithContext(ctx)
}

func (r *Request) SetValue(key string, value interface{}) {
	r.requestMu.Lock()
	*r.request = *r.WithContext(context.WithValue(r.request.Context(), key, value))
	r.requestMu.Unlock()
}

func (r *Request) SetMaxSize(maxSize int) {
	r.maxSize = maxSize
}

func (r *Request) MaxSize() int {
	if r.maxSize <= 0 {
		maxMemory := engine.DefaultMaxRequestBodySize
		if r.config != nil && r.config.MaxRequestBodySize != 0 {
			maxMemory = r.config.MaxRequestBodySize
		}
		return maxMemory
	}
	return r.maxSize
}

func (r *Request) Host() string {
	return r.request.Host
}

func (r *Request) URL() engine.URL {
	return r.url
}

func (r *Request) Header() engine.Header {
	return r.header
}

func (r *Request) Proto() string {
	return r.request.Proto
}

func (r *Request) RemoteAddress() string {
	return r.request.RemoteAddr
}

// RealIP implements `engine.Request#RealIP` function.
func (r *Request) RealIP() string {
	if len(r.realIP) > 0 {
		return r.realIP
	}

	r.realIP = realip.XRealIP(r.header.Get(echo.HeaderXRealIP), r.header.Get(echo.HeaderXForwardedFor), r.RemoteAddress())
	return r.realIP
}

func (r *Request) Method() string {
	return r.request.Method
}

func (r *Request) SetMethod(method string) {
	r.request.Method = method
}

func (r *Request) URI() string {
	return r.request.RequestURI
}

// SetURI implements `engine.Request#SetURI` function.
func (r *Request) SetURI(uri string) {
	r.request.RequestURI = uri
}

func (r *Request) Body() io.ReadCloser {
	return r.request.Body
}

// SetBody implements `engine.Request#SetBody` function.
func (r *Request) SetBody(reader io.Reader) {
	if readCloser, ok := reader.(io.ReadCloser); ok {
		r.request.Body = readCloser
	} else {
		r.request.Body = io.NopCloser(reader)
	}
}

func (r *Request) FormValue(name string) string {
	return r.value.Get(name)
}

func (r *Request) Form() engine.URLValuer {
	return r.value
}

func (r *Request) PostForm() engine.URLValuer {
	return r.value.postArgs
}

func (r *Request) MultipartForm() *multipart.Form {
	if r.request.MultipartForm == nil {
		r.request.ParseMultipartForm(int64(r.MaxSize()))
	}
	return r.request.MultipartForm
}

func (r *Request) IsTLS() bool {
	return r.request.TLS != nil
}

func (r *Request) Cookie(key string) string {
	if cookie, err := r.request.Cookie(key); err == nil {
		return cookie.Value
	}
	return ``
}

func (r *Request) Referer() string {
	return r.request.Referer()
}

func (r *Request) UserAgent() string {
	return r.request.UserAgent()
}

func (r *Request) Object() interface{} {
	return r.request
}

func (r *Request) reset(req *http.Request, h *Header, u *URL) {
	r.requestMu = sync.RWMutex{}
	r.request = req
	r.header = h
	r.url = u
	r.value = NewValue(r)
	r.realIP = ``
	r.maxSize = 0
}

func (r *Request) FormFile(key string) (multipart.File, *multipart.FileHeader, error) {
	r.MultipartForm()
	file, fileHeader, err := r.request.FormFile(key)
	if err != nil {
		return nil, nil, err
	}
	return file, fileHeader, err
}

// Size implements `engine.Request#ContentLength` function.
func (r *Request) Size() int64 {
	return r.request.ContentLength
}

func (r *Request) Scheme() string {
	if r.IsTLS() {
		return echo.SchemeHTTPS
	}
	if len(r.request.URL.Scheme) > 0 {
		return r.request.URL.Scheme
	}
	return echo.SchemeHTTP
}

func (r *Request) BasicAuth() (username, password string, ok bool) {
	return r.request.BasicAuth()
}

// SetHost implements `engine.Request#SetHost` function.
func (r *Request) SetHost(host string) {
	r.request.Host = host
}

func (r *Request) StdRequest() *http.Request {
	return r.request
}
