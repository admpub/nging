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

package servermanager

import (
	"strings"

	"github.com/webx-top/db"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/param"

	"github.com/admpub/nging/v4/application/handler"
	"github.com/admpub/nging/v4/application/library/role"
	"github.com/admpub/nging/v4/application/library/role/roleutils"
	"github.com/admpub/nging/v4/application/registry/dashboard"

	"github.com/nging-plugins/servermanager/application/library/system"
	"github.com/nging-plugins/servermanager/application/model"
)

func init() {
	dashboard.BlockRegister((&dashboard.Block{
		Tmpl:   `server/chart/cpu`,
		Footer: `server/chart/cpu.js`,
	}).SetContentGenerator(func(ctx echo.Context) error {
		ctx.Set(`systemRealtimeStatusIsListening`, system.RealTimeStatusIsListening())
		return nil
	}))
	dashboard.BlockRegister((&dashboard.Block{
		Tmpl: `server/dashbord/cmd_list`,
	}).SetContentGenerator(func(ctx echo.Context) error {
		user := handler.User(ctx)
		//指令集
		cmdMdl := model.NewCommand(ctx)
		if user.Id == 1 {
			cmdMdl.ListByOffset(nil, nil, 0, -1)
		} else {
			roleList := roleutils.UserRoles(ctx)
			cmdIds := []string{}
			for _, row := range roleList {
				for _, p := range row.Permissions {
					if p.Type == role.RolePermissionTypeCommand {
						cmdIds = append(cmdIds, strings.Split(p.Permission, `,`)...)
					}
				}
			}
			if len(cmdIds) > 0 {
				cmdIds = param.StringSlice(cmdIds).Unique().String()
				cmdMdl.ListByOffset(nil, nil, 0, -1, db.Cond{`id`: db.In(cmdIds)})
			}
		}
		ctx.Set(`cmdList`, cmdMdl.Objects())
		return nil
	}))
}
