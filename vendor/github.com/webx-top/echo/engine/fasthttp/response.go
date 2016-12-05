// +build !appengine

package fasthttp

import (
	"errors"
	"io"
	"net"
	"net/http"
	"strings"

	"github.com/admpub/fasthttp"
	"github.com/admpub/log"
	"github.com/webx-top/echo/engine"
	"github.com/webx-top/echo/logger"
)

var ErrAlreadyCommitted = errors.New(`response already committed`)

type (
	Response struct {
		context           *fasthttp.RequestCtx
		header            engine.Header
		status            int
		size              int64
		committed         bool
		writer            io.Writer
		logger            logger.Logger
		stdResponseWriter http.ResponseWriter
	}
)

func NewResponse(c *fasthttp.RequestCtx) *Response {
	return &Response{
		context: c,
		header:  &ResponseHeader{header: &c.Response.Header, stdhdr: nil},
		writer:  c,
		logger:  log.New("echo"),
	}
}

func (r *Response) Object() interface{} {
	return r.context
}

func (r *Response) Header() engine.Header {
	return r.header
}

func (r *Response) WriteHeader(code int) {
	if r.committed {
		r.logger.Warn(ErrAlreadyCommitted.Error())
		return
	}
	r.status = code
	r.context.SetStatusCode(code)
	r.committed = true
}

func (r *Response) Write(b []byte) (n int, err error) {
	n, err = r.writer.Write(b)
	r.size += int64(n)
	return
}

func (r *Response) Status() int {
	return r.status
}

func (r *Response) Size() int64 {
	return r.size
}

func (r *Response) Committed() bool {
	return r.committed
}

func (r *Response) SetWriter(w io.Writer) {
	r.writer = w
}

func (r *Response) Writer() io.Writer {
	return r.writer
}

func (r *Response) Hijack(fn func(net.Conn)) {
	r.context.Hijack(fasthttp.HijackHandler(fn))
	r.committed = true
}

func (r *Response) Body() []byte {
	switch strings.ToLower(r.header.Get(`Content-Encoding`)) {
	case `gzip`:
		body, err := r.context.Response.BodyGunzip()
		if err != nil {
			r.logger.Error(err)
		}
		return body
	case `deflate`:
		body, err := r.context.Response.BodyInflate()
		if err != nil {
			r.logger.Error(err)
		}
		return body
	default:
		return r.context.Response.Body()
	}
}

func (r *Response) Redirect(url string, code int) {
	//r.context.Redirect(url, code)  bug: missing port number
	r.header.Set(`Location`, url)
	r.WriteHeader(code)
}

func (r *Response) NotFound() {
	r.context.NotFound()
	r.committed = true
}

func (r *Response) SetCookie(cookie *http.Cookie) {
	r.header.Set("Set-Cookie", cookie.String())
}

func (r *Response) ServeFile(file string) {
	fasthttp.ServeFile(r.context, file)
	r.committed = true
}

func (r *Response) Error(errMsg string, args ...int) {
	if len(args) > 0 {
		r.status = args[0]
	} else {
		r.status = fasthttp.StatusInternalServerError
	}
	r.Write(engine.Str2bytes(errMsg))
	r.WriteHeader(r.status)
}

func (r *Response) reset(c *fasthttp.RequestCtx, h engine.Header) {
	r.context = c
	r.header = h
	r.status = http.StatusOK
	r.size = 0
	r.committed = false
	r.writer = c
	r.stdResponseWriter = nil
}

func (r *Response) StdResponseWriter() http.ResponseWriter {
	if r.stdResponseWriter != nil {
		return r.stdResponseWriter
	}
	w := &netHTTPResponseWriter{
		response: r,
	}
	r.stdResponseWriter = w
	return w
}

type netHTTPResponseWriter struct {
	h        http.Header
	response *Response
}

func (w *netHTTPResponseWriter) StatusCode() int {
	if w.response.Status() == 0 {
		return http.StatusOK
	}
	return w.response.Status()
}

func (w *netHTTPResponseWriter) Header() http.Header {
	if w.h == nil {
		w.h = make(http.Header)
	}
	return w.h
}

func (w *netHTTPResponseWriter) WriteHeader(statusCode int) {
	if w.response.committed {
		return
	}
	w.response.WriteHeader(statusCode)
	h := w.response.Header()
	for k, vv := range w.Header() {
		for _, v := range vv {
			h.Set(k, v)
		}
	}
}

func (w *netHTTPResponseWriter) Write(b []byte) (int, error) {
	if w.response.committed {
		return 0, ErrAlreadyCommitted
	}
	return w.response.Write(b)
}
