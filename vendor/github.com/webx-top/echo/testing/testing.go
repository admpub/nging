package testing

import (
	"net/http"
	"net/http/httptest"

	"github.com/admpub/log"

	"github.com/webx-top/echo/engine"
	"github.com/webx-top/echo/engine/standard"
)

// Request testing
func Request(method, path string, handler engine.Handler, reqRewrite ...func(*http.Request)) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, path, nil)
	if len(reqRewrite) > 0 && reqRewrite[0] != nil {
		reqRewrite[0](req)
	}
	rec := httptest.NewRecorder()

	handler.ServeHTTP(WrapRequest(req), WrapResponse(req, rec))
	//rec.Code, rec.Body.String(),rec.Header
	return rec
}

func WrapRequest(req *http.Request) engine.Request {
	return standard.NewRequest(req)
}

func WrapResponse(req *http.Request, rw http.ResponseWriter) engine.Response {
	return standard.NewResponse(rw, req, log.New().Sync())
}
