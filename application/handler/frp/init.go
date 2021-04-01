/*
   Nging is a toolbox for webmasters
   Copyright (C) 2018-present  Wenhui Shen <swh@admpub.com>

   This program is free software: you can redistribute it and/or modify
   it under the terms of the GNU Affero General Public License as published
   by the Free Software Foundation, either version 3 of the License, or
   (at your option) any later version.

   This program is distributed in the hope that it will be useful,
   but WITHOUT ANY WARRANTY; without even the implied warranty of
   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
   GNU Affero General Public License for more details.

   You should have received a copy of the GNU Affero General Public License
   along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/

package frp

import (
	"fmt"
	"net/http"
	"regexp"

	"github.com/admpub/nging/application/dbschema"
	"github.com/admpub/nging/application/handler"
	"github.com/admpub/nging/application/handler/caddy"
	"github.com/webx-top/db"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/code"
)

var regexNumEnd = regexp.MustCompile(`_[\d]+$`)

type Section struct {
	Section string
	Addon   string
}

func init() {
	handler.RegisterToGroup(`/frp`, func(g echo.RouteRegister) {
		e := handler.Echo()
		g.Route(`GET`, `/server_index`, ServerIndex)
		g.Route(`GET,POST`, `/server_add`, ServerAdd)
		g.Route(`GET,POST`, `/server_edit`, ServerEdit)
		g.Route(`GET,POST`, `/server_delete`, ServerDelete)
		g.Route(`GET,POST`, `/server_log`, ServerLog)
		g.Route(`GET`, `/client_index`, ClientIndex)
		g.Route(`GET,POST`, `/client_add`, ClientAdd)
		g.Route(`GET,POST`, `/client_edit`, ClientEdit)
		g.Route(`GET,POST`, `/client_delete`, ClientDelete)
		g.Route(`GET,POST`, `/client_log`, ClientLog)

		g.Route(`GET`, `/group_index`, GroupIndex)
		g.Route(`GET,POST`, `/group_add`, GroupAdd)
		g.Route(`GET,POST`, `/group_edit`, GroupEdit)
		g.Route(`GET,POST`, `/group_delete`, GroupDelete)
		g.Route(`GET,POST`, `/server_restart`, ServerRestart)
		g.Route(`GET,POST`, `/server_stop`, ServerStop)
		g.Route(`GET,POST`, `/client_restart`, ClientRestart)
		g.Route(`GET,POST`, `/client_stop`, ClientStop)
		g.Route(`GET`, `/addon_form`, e.MetaHandler(echo.H{`name`: `FRP客户端配置表单`}, AddonForm))
	})
	handler.RegisterToGroup(`/frp/dashboard`, func(g echo.RouteRegister) {
		g.Get(``, func(c echo.Context) error {
			m := &dbschema.NgingFrpServer{}
			err := m.Get(nil, `disabled`, `N`)
			if err != nil {
				if err != db.ErrNoMoreRows {
					return err
				}
				return c.NewError(code.DataNotFound, c.T(`没有找到启用的配置信息`))
			}
			if m.DashboardPort > 0 {
				dashboardHost := m.DashboardAddr
				if m.DashboardAddr == `0.0.0.0` || len(m.DashboardAddr) == 0 {
					dashboardHost = `127.0.0.1`
				}
				return c.Redirect(fmt.Sprintf(`http://%s:%d/`, dashboardHost, m.DashboardPort))
			}
			return c.NewError(code.Unsupported, c.T(`配置信息中未启用管理面板访问功能。如要启用，请将面板端口设为一个大于0的数值`))
		})
	})
}

func AddonForm(ctx echo.Context) error {
	addon := ctx.Query(`addon`)
	if len(addon) == 0 {
		return echo.NewHTTPError(http.StatusBadRequest, ctx.T("参数 addon 的值不能为空"))
	}
	if !caddy.ValidAddonName(addon) {
		return echo.NewHTTPError(http.StatusBadRequest, ctx.T("参数 addon 的值包含非法字符"))
	}
	section := ctx.Query(`section`, addon)
	setAddonFunc(ctx)
	return ctx.Render(`frp/client/form/`+addon, section)
}

func setAddonFunc(ctx echo.Context) {
	prefix := `extra`
	formKey := func(key string, keys ...string) string {
		key = prefix + `[` + key + `]`
		for _, k := range keys {
			key += `[` + k + `]`
		}
		return key
	}
	ctx.SetFunc(`Val`, func(key string, keys ...string) string {
		return ctx.Form(formKey(key, keys...))
	})
	ctx.SetFunc(`Vals`, func(key string, keys ...string) []string {
		return ctx.FormValues(formKey(key, keys...))
	})
	ctx.SetFunc(`Key`, formKey)
}
