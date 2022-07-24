package standard

import (
	"net/http"
	"sync"
)

func NewHeader(header http.Header) *Header {
	return &Header{header: header}
}

type Header struct {
	header http.Header
	lock   sync.RWMutex
}

func (h *Header) Add(key string, value string) {
	h.lock.Lock()
	h.header.Add(key, value)
	h.lock.Unlock()
}

func (h *Header) Del(key string) {
	h.lock.Lock()
	h.header.Del(key)
	h.lock.Unlock()
}

func (h *Header) Get(key string) string {
	h.lock.RLock()
	v := h.header.Get(key)
	h.lock.RUnlock()
	return v
}

func (h *Header) Values(key string) []string {
	h.lock.RLock()
	v := h.header.Values(key)
	h.lock.RUnlock()
	return v
}

func (h *Header) Set(key string, value string) {
	h.lock.Lock()
	h.header.Set(key, value)
	h.lock.Unlock()
}

func (h *Header) reset(hdr http.Header) {
	h.lock.Lock()
	h.header = hdr
	h.lock.Unlock()
}

func (h *Header) Object() interface{} {
	h.lock.RLock()
	v := h.header
	h.lock.RUnlock()
	return v
}

func (h *Header) Std() http.Header {
	h.lock.RLock()
	v := h.header
	h.lock.RUnlock()
	return v
}
