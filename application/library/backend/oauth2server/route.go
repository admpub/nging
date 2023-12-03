package oauth2server

import (
	"github.com/admpub/nging/v5/application/registry/route"
	"github.com/webx-top/echo"
)

var (
	// RoutePrefix 路由前缀
	RoutePrefix string
)

func Route(g echo.RouteRegister) {
	if len(RoutePrefix) > 0 {
		g = g.Group(RoutePrefix)
	}
	g.Route(`GET,POST`, `/authorize`, authorizeHandler).SetMetaKV(route.PermGuestKV()) // (step.1)(step.4)
	g.Route(`GET,POST`, `/login`, loginHandler).SetMetaKV(route.PermGuestKV())         // (step.2) 登录页面
	g.Route(`GET,POST`, `/auth`, authHandler).SetMetaKV(route.PermGuestKV())           // (step.3) 授权页面
	g.Route(`GET,POST`, `/logout`, logoutHandler).SetMetaKV(route.PermGuestKV())       // 退出登录
	g.Route(`GET,POST`, `/token`, tokenHandler).SetMetaKV(route.PermGuestKV())
	g.Route(`GET,POST`, `/profile`, profileHandler).SetMetaKV(route.PermGuestKV()) // 获取用户个人资料
}
