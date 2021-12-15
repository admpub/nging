package handler

import (
	"github.com/admpub/godownloader/service"
	"github.com/webx-top/echo"
	mw "github.com/webx-top/echo/middleware"

	"github.com/admpub/nging/v4/application/handler"
	"github.com/admpub/nging/v4/application/library/config"
)

var downloadDir = func() string {
	if len(config.DefaultConfig.Download.SavePath) == 0 {
		return service.GetDownloadPath()
	}
	return config.DefaultConfig.Download.SavePath
}

func init() {
	server := &service.DServ{}
	server.SetTmpl(`download/index`)
	server.SetSavePath(downloadDir)
	handler.RegisterToGroup(`/download`, func(g echo.RouteRegister) {
		server.Register(g, true)
		g.Route(`GET,POST`, `/file`, File, mw.CORS())
	})
}
