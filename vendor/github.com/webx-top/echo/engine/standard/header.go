package standard

import "net/http"

type Header struct {
	http.Header
}

func (h *Header) reset(hdr http.Header) {
	h.Header = hdr
}

func (h *Header) Object() interface{} {
	return h.Header
}

func (h *Header) Std() http.Header {
	return h.Header
}
