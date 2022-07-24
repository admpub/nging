package mock

import (
	"bufio"
	"bytes"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/admpub/log"

	"github.com/webx-top/echo/engine"
	"github.com/webx-top/echo/engine/standard"
	"github.com/webx-top/echo/logger"
)

type Response struct {
	http.ResponseWriter
	request   *http.Request
	status    int
	size      int64
	committed bool
	logger    logger.Logger
	mutex     sync.RWMutex
}

func (r *Response) Header() engine.Header {
	r.mutex.RLock()
	h := standard.NewHeader(r.ResponseWriter.Header())
	r.mutex.RUnlock()
	return h
}

func (r *Response) WriteHeader(code int) {
	if r.Committed() {
		r.logger.Warn("response already committed")
		return
	}

	r.mutex.Lock()
	r.status = code
	r.ResponseWriter.WriteHeader(code)
	r.committed = true
	r.mutex.Unlock()
}

func (r *Response) writeHeaderNoLock(code int) {
	if r.committed {
		r.logger.Warn("response already committed")
		return
	}

	r.status = code
	r.ResponseWriter.WriteHeader(code)
	r.committed = true
}

func (r *Response) KeepBody(on bool) {
}

func (r *Response) Write(b []byte) (n int, err error) {
	r.mutex.Lock()
	if !r.committed {
		if r.status == 0 {
			r.status = http.StatusOK
		}
		r.writeHeaderNoLock(r.status)
	}
	n, err = r.ResponseWriter.Write(b)
	r.size += int64(n)
	r.mutex.Unlock()
	return
}

func (r *Response) writeNoLock(b []byte) (n int, err error) {
	if !r.committed {
		if r.status == 0 {
			r.status = http.StatusOK
		}
		r.writeHeaderNoLock(r.status)
	}
	n, err = r.ResponseWriter.Write(b)
	r.size += int64(n)
	return
}

func (r *Response) Status() int {
	r.mutex.RLock()
	status := r.status
	r.mutex.RUnlock()
	return status
}

func (r *Response) Size() int64 {
	r.mutex.RLock()
	size := r.size
	r.mutex.RUnlock()
	return size
}

func (r *Response) Committed() bool {
	r.mutex.RLock()
	committed := r.committed
	r.mutex.RUnlock()
	return committed
}

func (r *Response) SetWriter(w io.Writer) {
}

func (r *Response) Writer() io.Writer {
	return r.ResponseWriter
}

func (r *Response) Object() interface{} {
	return r.ResponseWriter
}

func (r *Response) Error(errMsg string, args ...int) {
	r.mutex.Lock()
	if len(args) > 0 {
		r.status = args[0]
	} else {
		r.status = http.StatusInternalServerError
	}
	r.writeNoLock(engine.Str2bytes(errMsg))
	r.writeHeaderNoLock(r.status)
	r.mutex.Unlock()
}

func (r *Response) Reset(w http.ResponseWriter, req *http.Request) {
	r.mutex.Lock()
	r.ResponseWriter = w
	r.request = req
	r.status = http.StatusOK
	r.size = 0
	r.committed = false
	r.mutex.Unlock()
}

func (r *Response) Hijacker(fn func(net.Conn)) error {
	conn, bufrw, err := r.ResponseWriter.(http.Hijacker).Hijack()
	if err != nil {
		return err
	}
	_ = bufrw
	fn(conn)
	conn.Close()
	r.setCommitted(true)
	return nil
}

func (r *Response) setCommitted(committed bool) {
	r.mutex.Lock()
	r.committed = committed
	r.mutex.Unlock()
}

func (r *Response) Body() []byte {
	w, y := r.ResponseWriter.(*ResponseWriter)
	if y {
		r.mutex.Lock()
		b := w.Bytes()
		r.mutex.Unlock()
		return b
	}
	return nil
}

func (r *Response) Redirect(url string, code int) {
	r.mutex.Lock()
	http.Redirect(r.ResponseWriter, r.request, url, code)
	r.committed = true
	r.mutex.Unlock()
}

func (r *Response) NotFound() {
	r.mutex.Lock()
	http.Error(r.ResponseWriter, http.StatusText(http.StatusNotFound), http.StatusNotFound)
	r.committed = true
	r.mutex.Unlock()
}

func (r *Response) SetCookie(cookie *http.Cookie) {
	r.mutex.Lock()
	r.ResponseWriter.Header().Add(engine.HeaderSetCookie, cookie.String())
	r.mutex.Unlock()
}

func (r *Response) ServeFile(file string) {
	r.mutex.Lock()
	http.ServeFile(r.ResponseWriter, r.request, file)
	r.committed = true
	r.mutex.Unlock()
}

func (r *Response) ServeContent(content io.ReadSeeker, name string, modtime time.Time) {
	r.mutex.Lock()
	http.ServeContent(r.ResponseWriter, r.request, name, modtime, content)
	r.committed = true
	r.mutex.Unlock()
}

func (r *Response) Stream(step func(io.Writer) bool) (err error) {
	r.mutex.RLock()
	ctx := r.request.Context()
	r.mutex.RUnlock()
	for {
		select {
		case <-ctx.Done():
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
	return r.ResponseWriter
}

func NewResponseWriter(w *http.Response) http.ResponseWriter {
	b := bytes.NewBuffer(nil)
	w.Body = ioutil.NopCloser(b)
	return &ResponseWriter{
		Response: w,
		bytes:    b,
	}
}

type ResponseWriter struct {
	*http.Response
	bytes *bytes.Buffer
}

func (w *ResponseWriter) Header() http.Header {
	return w.Response.Header
}

func (w *ResponseWriter) Write(b []byte) (int, error) {
	return w.bytes.Write(b)
}

func (w *ResponseWriter) WriteHeader(statusCode int) {
	w.Response.StatusCode = statusCode
	w.Response.Status = http.StatusText(statusCode)
}

func (w *ResponseWriter) Bytes() []byte {
	r := w.bytes.Bytes()
	w.bytes.Reset()
	return r
}

func (w *ResponseWriter) String() string {
	r := w.bytes.String()
	w.bytes.Reset()
	return r
}

func NewResponse(args ...interface{}) *Response {
	var w http.ResponseWriter
	var r *http.Request
	var l logger.Logger
	for _, arg := range args {
		switch a := arg.(type) {
		case http.ResponseWriter:
			w = a
		case *http.Request:
			r = a
		case logger.Logger:
			l = a
		case engine.Request:
			r = a.StdRequest()
		}
	}
	if r == nil {
		r = &http.Request{
			URL:    &url.URL{},
			Header: http.Header{},
		}
	}
	if w == nil {
		w = NewResponseWriter(&http.Response{
			Request: r,
			Header:  http.Header{},
		})
	}
	if l == nil {
		l = log.GetLogger(`mock`)
	}
	return &Response{
		ResponseWriter: w,
		request:        r,
		logger:         l,
	}
}
