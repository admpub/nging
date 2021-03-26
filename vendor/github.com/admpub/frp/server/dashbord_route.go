package server

import (
	"github.com/admpub/frp/assets"
	"github.com/webx-top/echo"
)

// RegisterTo 为echo框架创建路由
func (svr *Service) RegisterTo(router echo.RouteRegister) {
	// api, see dashboard_api.go
	router.Get("/api/serverinfo", svr.ApiServerInfo)
	router.Get("/api/proxy/:type", svr.ApiProxyByType)
	router.Get("/api/proxy/:type/:name", svr.ApiProxyByTypeAndName)
	router.Get("/api/traffic/:name", svr.ApiProxyTraffic)
	// view
	router.Get("/", func(c echo.Context) error {
		return c.Redirect("./static/")
	})
	cfg := &svr.cfg
	//cfg.AssetsDir = `/Users/hank/go/src/github.com/admpub/frp/assets/static`
	err := assets.Load(cfg.AssetsDir)
	if err != nil {
		panic(err)
	}
	fs := assets.FS(`server`)
	router.Get("/static*", func(c echo.Context) error {
		file := c.Param(`*`)
		if len(file) == 0 || file == `/` {
			file = `/index.html`
		}
		return c.File(file, fs)
	})
}
