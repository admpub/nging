package multiuser

import (
	"strconv"

	plugin "github.com/admpub/frp/pkg/plugin/server"
	"github.com/admpub/log"
	"github.com/webx-top/echo"

	"github.com/admpub/nging/v4/application/cmd/bootconfig"
	"github.com/admpub/nging/v4/application/handler"
	"github.com/admpub/nging/v4/application/library/common"
	"github.com/admpub/nging/v4/application/library/config"
	"github.com/nging-plugins/frpmanager/application/library/frp"
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
		} else {
			if config.FromCLI().Port != config.DefaultPort {
				p.Addr = `127.0.0.1:` + strconv.Itoa(config.FromCLI().Port)
			}
		}
		return p
	})
	config.OnKeySetSettings(`base.backendURL`, func(diff config.Diff) error {
		if !bootconfig.IsWeb() || !config.IsInstalled() || !diff.IsDiff {
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
