package setup

import (
	"github.com/admpub/nging/v4/application/handler"
	"github.com/webx-top/echo"
)

func init() {
	handler.Register(func(e echo.RouteRegister) {
		e.Route("GET,POST", `/setup`, Setup)
		e.Route("GET", `/progress`, Progress)
		e.Route("GET,POST", `/license`, License)
	})
}
