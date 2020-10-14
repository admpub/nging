package main

import (
	"github.com/admpub/nging/application/handler"
	"github.com/webx-top/echo"
)

func main() {
	handler.RegisterToGroup(`/plugins`, func(g echo.RouteRegister) {
		g.Route(`GET,POST`, `/example`, echo.HandlerFunc(func(ctx echo.Context) error {
			return ctx.String(`plugins.example`)
		}))
	})
	println(`------------------PLUGINS:example-------------`)
}
