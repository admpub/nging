package backend

import (
	"github.com/admpub/nging/v5/application/handler"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/subdomains"
)

func addGlobalFuncMap(fm map[string]interface{}) map[string]interface{} {
	fm[`AssetsURL`] = getAssetsURL
	fm[`BackendURL`] = getBackendURL
	fm[`FrontendURL`] = getFrontendURL
	return fm
}

func getAssetsURL(paths ...string) (r string) {
	r = AssetsURLPath
	for _, ppath := range paths {
		r += ppath
	}
	return r
}

func BackendURLFuncMW() echo.MiddlewareFunc {
	return func(h echo.Handler) echo.Handler {
		return echo.HandlerFunc(func(c echo.Context) error {
			BackendURLFunc(c)
			return h.Handle(c)
		})
	}
}

func getBackendURL(paths ...string) (r string) {
	r = handler.BackendPrefix
	for _, ppath := range paths {
		r += ppath
	}
	return r
	//return subdomains.Default.URL(r, `backend`)
}

func getFrontendURL(paths ...string) (r string) {
	r = handler.FrontendPrefix
	for _, ppath := range paths {
		r += ppath
	}
	return subdomains.Default.URL(r, `frontend`)
}

func BackendURLFunc(c echo.Context) error {
	return nil
}
