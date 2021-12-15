package standard

import (
	"bufio"
	"io"
	"net"
	"net/http"
	"time"

	"github.com/webx-top/echo/engine"
	"github.com/webx-top/echo/logger"
)

type Response struct {
	http.ResponseWriter
	config         *engine.Config
	request        *http.Request
	header         *Header
	status         int
	size           int64
	committed      bool
	writer         io.Writer
	logger         logger.Logger
	body           []byte
	keepBody       bool
	responseWriter *responseWriter
}

func NewResponse(w http.ResponseWriter, r *http.Request, l logger.Logger) *Response {
	return &Response{
		ResponseWriter: w,
		request:        r,
		header:         &Header{header: w.Header()},
		writer:         w,
		logger:         l,
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
	r.header.lock.Lock()
	r.ResponseWriter.WriteHeader(code)
	r.header.lock.Unlock()
	r.committed = true
}

func (r *Response) KeepBody(on bool) {
	r.keepBody = on
}

func (r *Response) Write(b []byte) (n int, err error) {
	if !r.committed {
		if r.status == 0 {
			r.status = http.StatusOK
		}
		r.WriteHeader(r.status)
	}
	if r.keepBody {
		r.body = append(r.body, b...)
	}
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
	return r.ResponseWriter
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

func (r *Response) reset(w http.ResponseWriter, req *http.Request, h *Header) {
	r.ResponseWriter = w
	r.request = req
	r.header = h
	r.status = http.StatusOK
	r.size = 0
	r.committed = false
	r.writer = w
	r.body = nil
	r.keepBody = false
	r.responseWriter = &responseWriter{r}
}

func (r *Response) Hijacker(fn func(net.Conn)) error {
	conn, bufrw, err := r.ResponseWriter.(http.Hijacker).Hijack()
	if err != nil {
		return err
	}
	_ = bufrw
	fn(conn)
	conn.Close()
	r.committed = true
	return nil
}

func (r *Response) Body() []byte {
	return r.body
}

func (r *Response) Redirect(url string, code int) {
	http.Redirect(r.ResponseWriter, r.request, url, code)
	r.committed = true
}

func (r *Response) NotFound() {
	http.Error(r.ResponseWriter, http.StatusText(http.StatusNotFound), http.StatusNotFound)
	r.committed = true
}

func (r *Response) SetCookie(cookie *http.Cookie) {
	r.header.Add(engine.HeaderSetCookie, cookie.String())
}

func (r *Response) ServeFile(file string) {
	r.keepBody = false
	http.ServeFile(r.ResponseWriter, r.request, file)
	r.committed = true
}

func (r *Response) ServeContent(content io.ReadSeeker, name string, modtime time.Time) {
	r.keepBody = false
	http.ServeContent(r.ResponseWriter, r.request, name, modtime, content)
	r.committed = true
}

func (r *Response) Stream(step func(io.Writer) bool) (err error) {
	for {
		select {
		case <-r.request.Context().Done():
			return
		default:
			keepOpen := step(r)
			r.Flush()
			if !keepOpen {
				return
			}
		}
	}
}

func (r *Response) Flush() {
	if flusher, ok := r.ResponseWriter.(http.Flusher); ok {
		flusher.Flush()
	}
}

func (r *Response) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	return r.ResponseWriter.(http.Hijacker).Hijack()
}

func (r *Response) StdResponseWriter() http.ResponseWriter {
	if r.responseWriter == nil {
		r.responseWriter = &responseWriter{r}
	}
	return r.responseWriter
}

type responseWriter struct {
	*Response
}

func (r *responseWriter) StatusCode() int {
	if r.Response.Status() == 0 {
		return http.StatusOK
	}
	return r.Response.Status()
}

func (r *responseWriter) Header() http.Header {
	return r.Response.header.Std()
}

func (r *responseWriter) Write(b []byte) (n int, err error) {
	return r.Response.Write(b)
}

func (r *responseWriter) WriteHeader(code int) {
	r.Response.WriteHeader(code)
}
