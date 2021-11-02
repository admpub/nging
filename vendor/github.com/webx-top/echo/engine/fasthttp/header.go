//go:build !appengine
// +build !appengine

package fasthttp

import (
	"net/http"
	"sync"

	"github.com/admpub/fasthttp"
	"github.com/webx-top/echo/engine"
)

type (
	RequestHeader struct {
		header *fasthttp.RequestHeader
		stdhdr *http.Header
		lock   sync.RWMutex
	}

	ResponseHeader struct {
		header *fasthttp.ResponseHeader
		stdhdr *http.Header
		lock   sync.RWMutex
	}
)

func (h *RequestHeader) Add(key, val string) {
	h.lock.Lock()
	h.header.Set(key, val)
	h.lock.Unlock()
}

func (h *RequestHeader) Del(key string) {
	h.lock.Lock()
	h.header.Del(key)
	h.lock.Unlock()
}

func (h *RequestHeader) Get(key string) string {
	h.lock.RLock()
	v := engine.Bytes2str(h.header.Peek(key))
	h.lock.RUnlock()
	return v
}

func (h *RequestHeader) Values(key string) []string {
	h.lock.RLock()
	v := engine.Bytes2str(h.header.Peek(key))
	h.lock.RUnlock()
	return []string{v}
}

func (h *RequestHeader) Set(key, val string) {
	h.lock.Lock()
	h.header.Set(key, val)
	h.lock.Unlock()
}

func (h *RequestHeader) Object() interface{} {
	h.lock.RLock()
	v := h.header
	h.lock.RUnlock()
	return v
}

func (h *ResponseHeader) Add(key, val string) {
	h.lock.Lock()
	h.header.Set(key, val)
	h.lock.Unlock()
}

func (h *RequestHeader) reset(hdr *fasthttp.RequestHeader) {
	h.lock.Lock()
	h.header = hdr
	h.lock.Unlock()
}

func (h *RequestHeader) Std() http.Header {
	h.lock.Lock()
	if h.stdhdr != nil {
		h.lock.Unlock()
		return *h.stdhdr
	}
	hdr := http.Header{}
	h.header.VisitAll(func(key, value []byte) {
		hdr.Add(string(key), string(value))
	})
	h.stdhdr = &hdr
	h.lock.Unlock()
	return hdr
}

func (h *ResponseHeader) Del(key string) {
	h.lock.Lock()
	h.header.Del(key)
	h.lock.Unlock()
}

func (h *ResponseHeader) Get(key string) string {
	h.lock.RLock()
	v := engine.Bytes2str(h.header.Peek(key))
	h.lock.RUnlock()
	return v
}

func (h *ResponseHeader) Values(key string) []string {
	h.lock.RLock()
	v := engine.Bytes2str(h.header.Peek(key))
	h.lock.RUnlock()
	return []string{v}
}

func (h *ResponseHeader) Set(key, val string) {
	h.lock.Lock()
	h.header.Set(key, val)
	h.lock.Unlock()
}

func (h *ResponseHeader) Object() interface{} {
	h.lock.RLock()
	v := h.header
	h.lock.RUnlock()
	return v
}

func (h *ResponseHeader) reset(hdr *fasthttp.ResponseHeader) {
	h.lock.Lock()
	h.header = hdr
	h.lock.Unlock()
}

func (h *ResponseHeader) Std() http.Header {
	h.lock.Lock()
	if h.stdhdr != nil {
		return *h.stdhdr
	}
	hdr := http.Header{}
	h.header.VisitAll(func(key, value []byte) {
		hdr.Add(string(key), string(value))
	})
	h.stdhdr = &hdr
	h.lock.Unlock()
	return hdr
}
