package test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

type Config struct {
	Method, Path string
	Handler      HandlerTest
	ReqRewrite   []func(*http.Request)
	Checker      func(*testing.T, *httptest.ResponseRecorder, *bytes.Buffer)
	Middlewares  []MiddlewareTest
}

func DefaultChecker(value string) func(t *testing.T, r *httptest.ResponseRecorder, buf *bytes.Buffer) {
	return func(t *testing.T, r *httptest.ResponseRecorder, buf *bytes.Buffer) {
		Eq(t, value, buf.String())
		Eq(t, `OK`, r.Body.String())
		Eq(t, http.StatusOK, r.Code)
	}
}
