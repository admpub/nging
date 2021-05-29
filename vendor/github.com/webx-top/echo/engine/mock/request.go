package mock

import (
	"net/http"
	"net/url"

	"github.com/webx-top/echo/engine"
	"github.com/webx-top/echo/engine/standard"
)

func NewRequest(reqs ...*http.Request) engine.Request {
	var req *http.Request
	if len(reqs) > 0 {
		req = reqs[0]
	}
	if req == nil {
		req = &http.Request{
			URL:    &url.URL{},
			Header: http.Header{},
		}
	}
	return &Request{
		Request: standard.NewRequest(req),
	}
}

type Request struct {
	*standard.Request
}
