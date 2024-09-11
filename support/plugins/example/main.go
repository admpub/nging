package main

import (
	"github.com/coscms/webcore/registry/route"
	"github.com/webx-top/echo"
)

func main() {
	route.RegisterToGroup(`/plugins`, func(g echo.RouteRegister) {
		g.Route(`GET,POST`, `/example`, echo.HandlerFunc(func(ctx echo.Context) error {
			return ctx.String(`plugins.example`)
		}))
	})
	println(`------------------PLUGINS:example-------------`)
}
