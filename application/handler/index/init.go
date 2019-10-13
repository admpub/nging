package index

import (
	"github.com/webx-top/echo"

	"github.com/admpub/nging/application/handler"
	"github.com/admpub/nging/application/middleware"
)

func init() {
	handler.Register(func(e echo.RouteRegister) {
		e.Route("GET", ``, Index)
		e.Route("GET", `/`, Index)
		//e.Route("GET", `/ts2m3u8`, TS2M3U8, mw.CORS())
		e.Route("GET", `/project/:ident`, Project)
		e.Route("GET", `/index`, Index)
		e.Route("GET,POST", `/login`, Login)
		e.Route("GET,POST", `/register`, Register)
		e.Route("GET", `/logout`, Logout)
		e.Route("GET", `/donation`, Donation)
		//e.Route(`GET,POST`, `/ping`, Ping)
		e.Get(`/icon`, Icon, middleware.AuthCheck)
		e.Get(`/routeList`, RouteList, middleware.AuthCheck)
		e.Get(`/routeNotin`, RouteNotin, middleware.AuthCheck)
		e.Get(`/navTree`, NavTree, middleware.AuthCheck)
	})
}
