package standard

import (
	"io"
	"net"
	"net/http"

	"github.com/webx-top/echo/engine"
	"github.com/webx-top/echo/logger"
)

type (
	Response struct {
		config    *engine.Config
		response  http.ResponseWriter
		request   *http.Request
		header    engine.Header
		status    int
		size      int64
		committed bool
		writer    io.Writer
		logger    logger.Logger
		body      []byte
	}
)

func NewResponse(w http.ResponseWriter, r *http.Request, l logger.Logger) *Response {
	return &Response{
		response: w,
		request:  r,
		header:   &Header{header: w.Header()},
		writer:   w,
		logger:   l,
	}
}

func (r *Response) Header() engine.Header {
	return r.header
}

func (r *Response) WriteHeader(code int) {
	if r.committed {
		r.logger.Warn("response already committed")
		return
	}
	r.status = code
	r.response.WriteHeader(code)
	r.committed = true
}

func (r *Response) Write(b []byte) (n int, err error) {
	r.body = b
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

func (r *Response) Object() interface{} {
	return r.response
}

func (r *Response) Error(errMsg string, args ...int) {
	if len(args) > 0 {
		r.status = args[0]
	} else {
		r.status = http.StatusInternalServerError
	}
	r.Write(engine.Str2bytes(errMsg))
	r.WriteHeader(r.status)
}

func (r *Response) reset(w http.ResponseWriter, req *http.Request, h engine.Header) {
	r.response = w
	r.request = req
	r.header = h
	r.status = http.StatusOK
	r.size = 0
	r.committed = false
	r.writer = w
	r.body = nil
}

func (r *Response) Hijack(fn func(net.Conn)) {
	conn, bufrw, err := r.response.(http.Hijacker).Hijack()
	if err != nil {
		r.logger.Error(err)
	}
	_ = bufrw
	fn(conn)
	conn.Close()
	r.committed = true
}

func (r *Response) Body() []byte {
	return r.body
}

func (r *Response) Redirect(url string, code int) {
	http.Redirect(r.response, r.request, url, code)
	r.committed = true
}

func (r *Response) NotFound() {
	http.Error(r.response, http.StatusText(http.StatusNotFound), http.StatusNotFound)
	r.committed = true
}

func (r *Response) SetCookie(cookie *http.Cookie) {
	r.header.Set("Set-Cookie", cookie.String())
}

func (r *Response) ServeFile(file string) {
	http.ServeFile(r.response, r.request, file)
	r.committed = true
}

func (r *Response) StdResponseWriter() http.ResponseWriter {
	return r.response
}
