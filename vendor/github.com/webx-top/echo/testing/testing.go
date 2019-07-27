package testing

import (
	"net/http"
	"net/http/httptest"

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

	handler.ServeHTTP(standard.NewRequest(req), standard.NewResponse(rec, req, nil))
	//rec.Code, rec.Body.String(),rec.Header
	return rec
}
