package proxy

import (
	"fmt"
	"net/url"

	"github.com/webx-top/com"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/code"
	mw "github.com/webx-top/echo/middleware"

	"github.com/admpub/nging/v4/application/library/config"
	"github.com/nging-plugins/frpmanager/pkg/model"
)

func presetProxyHTTP(c echo.Context) {
	user := c.Internal().String(`frp.admin.user`)
	password := c.Internal().String(`frp.admin.password`)
	auth := com.Base64Encode(user + `:` + password)
	c.Request().Header().Set(echo.HeaderAuthorization, "Basic "+auth)
}

func ProxyClient(next echo.Handler) echo.HandlerFunc {
	return func(c echo.Context) error {
		id := c.Paramx("id").Uint()
		m := model.NewFrpClient(c)
		err := m.Get(nil, `id`, id)
		if err != nil {
			return err
		}
		if m.AdminPort < 1 {
			return echo.ErrNotFound
		}
		if m.Disabled == `Y` {
			return c.NewError(code.DataStatusIncorrect, `未启用`)
		}
		if !config.FromCLI().IsRunning(`frpclient.` + fmt.Sprint(m.Id)) {
			return c.NewError(code.DataStatusIncorrect, `未启动`)
		}
		c.Internal().Set(`frp.admin.user`, m.AdminUser)
		c.Internal().Set(`frp.admin.password`, m.AdminPwd)
		address := fmt.Sprintf("%s:%d", m.AdminAddr, m.AdminPort)
		c.Internal().Set(`frp.admin.address`, address)
		return next.Handle(c)
	}
}

func ProxyServer(next echo.Handler) echo.HandlerFunc {
	return func(c echo.Context) error {
		id := c.Paramx("id").Uint()
		m := model.NewFrpServer(c)
		err := m.Get(nil, `id`, id)
		if err != nil {
			return err
		}
		if m.DashboardPort < 1 {
			return echo.ErrNotFound
		}
		if m.Disabled == `Y` {
			return c.NewError(code.DataStatusIncorrect, `未启用`)
		}
		if !config.FromCLI().IsRunning(`frpserver.` + fmt.Sprint(m.Id)) {
			return c.NewError(code.DataStatusIncorrect, `未启动`)
		}
		c.Internal().Set(`frp.admin.user`, m.DashboardUser)
		c.Internal().Set(`frp.admin.password`, m.DashboardPwd)
		address := fmt.Sprintf("%s:%d", m.DashboardAddr, m.DashboardPort)
		c.Internal().Set(`frp.admin.address`, address)
		return next.Handle(c)
	}
}

func Proxy() echo.MiddlewareFuncd {
	target := NewTarget(&mw.ProxyTarget{
		Name: `default`,
		URL: &url.URL{
			Scheme: "http",
			Host:   "",
		},
	})
	balancer := mw.NewRandomBalancer([]mw.ProxyTargeter{target})
	proxyMWCfg := mw.ProxyConfig{
		Skipper: func(c echo.Context) bool {
			presetProxyHTTP(c)
			return false
		},
		Handler: mw.DefaultProxyHandler,
		Rewrite: mw.RewriteConfig{
			Skipper: echo.DefaultSkipper,
			Rules: map[string]string{
				`/frp/dashboard/server/:id/*`: `/$2`,
				`/frp/dashboard/client/:id/*`: `/$2`,
			},
		},
		Balancer:   balancer,
		ContextKey: "target",
	}
	return mw.ProxyWithConfig(proxyMWCfg)
}
