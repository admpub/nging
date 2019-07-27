package test

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/webx-top/echo"
	testings "github.com/webx-top/echo/testing"
)

type HandlerTest func(*bytes.Buffer) echo.HandlerFunc
type MiddlewareTest func(*bytes.Buffer) echo.MiddlewareFuncd

func Hit(t *testing.T, configs []*Config, middelwares ...MiddlewareTest) {
	e := echo.New()
	hit(e, t, configs, middelwares...)
}

func HitBy(custom func(e *echo.Echo) echo.Context, t *testing.T, configs []*Config, middelwares ...MiddlewareTest) {
	e := echo.NewWithContext(custom)
	hit(e, t, configs, middelwares...)
}

func hit(e *echo.Echo, t *testing.T, configs []*Config, middelwares ...MiddlewareTest) {
	buf := new(bytes.Buffer)
	for _, h := range middelwares {
		e.Use(h(buf))
	}
	for _, cfg := range configs {
		ms := make([]interface{}, len(cfg.Middlewares))
		for k, m := range cfg.Middlewares {
			ms[k] = m(buf)
		}
		e.Match([]string{cfg.Method}, cfg.Path, cfg.Handler(buf), ms...)
		r := testings.Request(cfg.Method, cfg.Path, e, cfg.ReqRewrite...)
		cfg.Checker(t, r, buf)
	}
}

var (
	Eq          = assert.Equal
	NotEq       = assert.NotEqual
	True        = assert.True
	False       = assert.False
	NotNil      = assert.NotNil
	Empty       = assert.Empty
	NotEmpty    = assert.NotEmpty
	Len         = assert.Len
	Contains    = assert.Contains
	NotContains = assert.NotContains
	Subset      = assert.Subset
	NotSubset   = assert.NotSubset
)
