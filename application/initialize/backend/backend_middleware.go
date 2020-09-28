package backend

import (
	"github.com/admpub/nging/application/handler"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/subdomains"
)

func BackendURLFuncMW() echo.MiddlewareFunc {
	return func(h echo.Handler) echo.Handler {
		return echo.HandlerFunc(func(c echo.Context) error {
			BackendURLFunc(c)
			return h.Handle(c)
		})
	}
}

func BackendURLFunc(c echo.Context) error {
	c.SetFunc(`AssetsURL`, func() string {
		return AssetsURLPath
	})
	c.SetFunc(`BackendURL`, func() string {
		return subdomains.Default.URL(handler.BackendPrefix, `backend`)
	})
	c.SetFunc(`FrontendURL`, func() string {
		return subdomains.Default.URL(handler.FrontendPrefix, `frontend`)
	})
	return nil
}
