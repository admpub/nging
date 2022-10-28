package setup

import (
	"github.com/admpub/nging/v5/application/handler"
	"github.com/admpub/nging/v5/application/request"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/defaults"
)

func init() {
	handler.Register(func(e echo.RouteRegister) {
		e.Route("GET,POST", `/setup`, defaults.MetaHandlerWithRequest(nil, Setup, request.Setup{}, `POST`))
		e.Route("GET", `/progress`, Progress)
		e.Route("GET,POST", `/license`, License)
	})
}
