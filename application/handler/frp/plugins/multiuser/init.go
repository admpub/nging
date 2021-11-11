package multiuser

import (
	"strconv"

	plugin "github.com/admpub/frp/pkg/plugin/server"
	"github.com/admpub/log"
	"github.com/admpub/nging/v3/application/cmd/event"
	"github.com/admpub/nging/v3/application/handler"
	"github.com/admpub/nging/v3/application/library/common"
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
	register(`multiuser_login`, `多用户登录`, func() plugin.HTTPPluginOptions {
		p := definePlugin
		backendURL := config.Setting(`base`).String(`backendURL`)
		if len(backendURL) > 0 {
			p.Addr = backendURL
		}
		return p
	})
	config.OnKeySetSettings(`base.backendURL`, func(config.Diff) error {
		if !event.IsWeb() {
			return nil
		}
		go func() {
			ctx := common.NewMockContext()
			if err := OnChangeBackendURL(ctx); err != nil {
				log.Error(err)
			}
		}()
		return nil
	})
	handler.Register(func(g echo.RouteRegister) {
		g.Post(`/frp_login`, Login)
	})
}
