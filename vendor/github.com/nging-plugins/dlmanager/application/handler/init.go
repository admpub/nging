package handler

import (
	"github.com/admpub/godownloader/service"
	"github.com/admpub/log"
	"github.com/webx-top/echo"
	mw "github.com/webx-top/echo/middleware"

	"github.com/admpub/nging/v5/application/library/config/startup"
	"github.com/admpub/nging/v5/application/library/route"
	dlconfig "github.com/nging-plugins/dlmanager/application/library/config"
)

var Server = &service.DServ{}

func RegisterRoute(r *route.Collection) {
	r.Backend.RegisterToGroup(`/download`, registerRoute)
}

func registerRoute(g echo.RouteRegister) {
	Server.Register(g, true)
	g.Route(`GET,POST`, `/file`, File, mw.CORS())
}

var downloadDir = func() string {
	if len(dlconfig.Get().SavePath) == 0 {
		return service.GetDownloadPath()
	}
	return dlconfig.Get().SavePath
}

func init() {
	startup.OnAfter(`web.installed`, func() {
		go Server.LoadSettings()
	})
	startup.OnAfter(`web`, func() {
		if err := Server.SaveSettings(); err != nil {
			log.Error(err)
		}
	})
	Server.SetTmpl(`download/index`)
	Server.SetSavePath(downloadDir)
}
