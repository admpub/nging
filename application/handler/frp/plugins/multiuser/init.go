package multiuser

import (
	"strconv"

	frpConfig "github.com/admpub/frp/pkg/config"
	plugin "github.com/admpub/frp/pkg/plugin/server"
	frpLog "github.com/admpub/frp/pkg/util/log"
	"github.com/admpub/nging/v3/application/handler"
	"github.com/admpub/nging/v3/application/library/config"
	"github.com/admpub/nging/v3/application/library/frp"
	"github.com/webx-top/echo"
)

var (
	register     = frp.ServerPluginRegister
	definePlugin = plugin.HTTPPluginOptions{
		Name:      `multiuser_login`,
		Addr:      `127.0.0.1:` + strconv.Itoa(config.DefaultPort),
		Path:      `/frp_login`,
		Ops:       []string{`Login`}, // Login / NewProxy / NewWorkConn / NewUserConn / Ping
		TLSVerify: false,
	}
)

func init() {
	register(`multiuser_login`, `多用户登录`, func(_ *frpConfig.ServerCommonConf) plugin.HTTPPluginOptions {
		p := definePlugin
		backendURL := config.Setting(`base`).String(`backendURL`)
		if len(backendURL) > 0 {
			p.Addr = backendURL
		}
		frpLog.Info(`[frp] 注册多用户登录插件，插件接口地址: %s`, p.Addr)
		return p
	})
	handler.Register(func(g echo.RouteRegister) {
		g.Post(`/frp_login`, Login)
	})
}
