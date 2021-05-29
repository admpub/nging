package mock

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/admpub/log"

	"github.com/webx-top/echo/engine"
	"github.com/webx-top/echo/engine/standard"
	"github.com/webx-top/echo/logger"
)

type Response struct {
	*standard.Response
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
		case *standard.Request:
			r = a.StdRequest()
		case *Request:
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
		Response: standard.NewResponse(w, r, l),
	}
}
