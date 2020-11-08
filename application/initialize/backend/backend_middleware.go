package backend

import (
	"strings"

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
	c.SetFunc(`AssetsURL`, func(paths ...string) (r string) {
		r = AssetsURLPath
		for _, ppath := range paths {
			r += ppath
		}
		return r
	})
	c.SetFunc(`BackendURL`, func(paths ...string) (r string) {
		r = handler.BackendPrefix
		for _, ppath := range paths {
			r += ppath
		}
		if !strings.HasPrefix(r, `/`) {
			r = `/` + r
		}
		return r
		//return subdomains.Default.URL(r, `backend`)
	})
	c.SetFunc(`FrontendURL`, func(paths ...string) (r string) {
		r = handler.FrontendPrefix
		for _, ppath := range paths {
			r += ppath
		}
		return subdomains.Default.URL(r, `frontend`)
	})
	return nil
}
