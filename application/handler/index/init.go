package index

import (
	"github.com/admpub/nging/application/handler"
	"github.com/webx-top/echo"
)

func init() {
	handler.Register(func(e *echo.Echo) {
		e.Route("GET", `/`, Index)
		e.Route("GET", `/index`, Index)
		e.Route("GET,POST", `/login`, Login)
		e.Route("GET,POST", `/register`, Register)
		e.Route("GET", `/logout`, Logout)
		e.Route("GET", `/donation`, Donation)
		//e.Route(`GET,POST`, `/ping`, Ping)
	})
}
