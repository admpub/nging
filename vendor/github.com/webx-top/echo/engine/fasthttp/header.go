// +build !appengine

package fasthttp

import (
	"net/http"

	"github.com/admpub/fasthttp"
	"github.com/webx-top/echo/engine"
)

type (
	RequestHeader struct {
		header *fasthttp.RequestHeader
		stdhdr *http.Header
	}

	ResponseHeader struct {
		header *fasthttp.ResponseHeader
		stdhdr *http.Header
	}
)

func (h *RequestHeader) Add(key, val string) {
	h.header.Set(key, val)
}

func (h *RequestHeader) Del(key string) {
	h.header.Del(key)
}

func (h *RequestHeader) Get(key string) string {
	return engine.Bytes2str(h.header.Peek(key))
}

func (h *RequestHeader) Set(key, val string) {
	h.header.Set(key, val)
}

func (h *RequestHeader) Object() interface{} {
	return h.header
}

func (h *ResponseHeader) Add(key, val string) {
	h.header.Set(key, val)
}

func (h *RequestHeader) reset(hdr *fasthttp.RequestHeader) {
	h.header = hdr
}

func (h *RequestHeader) Std() http.Header {
	if h.stdhdr != nil {
		return *h.stdhdr
	}
	hdr := http.Header{}
	h.header.VisitAll(func(key, value []byte) {
		hdr.Add(string(key), string(value))
	})
	h.stdhdr = &hdr
	return hdr
}

func (h *ResponseHeader) Del(key string) {
	h.header.Del(key)
}

func (h *ResponseHeader) Get(key string) string {
	return engine.Bytes2str(h.header.Peek(key))
}

func (h *ResponseHeader) Set(key, val string) {
	h.header.Set(key, val)
}

func (h *ResponseHeader) Object() interface{} {
	return h.header
}

func (h *ResponseHeader) reset(hdr *fasthttp.ResponseHeader) {
	h.header = hdr
}

func (h *ResponseHeader) Std() http.Header {
	if h.stdhdr != nil {
		return *h.stdhdr
	}
	hdr := http.Header{}
	h.header.VisitAll(func(key, value []byte) {
		hdr.Add(string(key), string(value))
	})
	h.stdhdr = &hdr
	return hdr
}
