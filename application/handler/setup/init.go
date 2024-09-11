package setup

import (
	"github.com/coscms/webcore/cmd/bootconfig"
	"github.com/coscms/webcore/registry/route"
	"github.com/coscms/webcore/request"
	"github.com/webx-top/echo"
)

func init() {
	route.Register(func(e echo.RouteRegister) {
		e.Route("GET,POST", `/setup`, route.HandlerWithRequest(Setup, request.Setup{}, `POST`))
		e.Route("GET", `/progress`, Progress)
		e.Route("GET,POST", `/license`, License)
	})
	bootconfig.Setup = Setup
	bootconfig.Upgrade = Upgrade
}
