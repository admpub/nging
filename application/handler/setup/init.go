package setup

import (
	"github.com/admpub/nging/v5/application/handler"
	"github.com/admpub/nging/v5/application/request"
	"github.com/webx-top/echo"
)

func init() {
	handler.Register(func(e echo.RouteRegister) {
		e.Route("GET,POST", `/setup`, handler.WithRequest(Setup, request.Setup{}, `POST`))
		e.Route("GET", `/progress`, Progress)
		e.Route("GET,POST", `/license`, License)
	})
}
