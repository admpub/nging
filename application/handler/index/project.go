/*
   Nging is a toolbox for webmasters
   Copyright (C) 2018-present Wenhui Shen <swh@admpub.com>

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

package index

import (
	"github.com/admpub/nging/v5/application/handler"
	"github.com/admpub/nging/v5/application/library/role"
	"github.com/admpub/nging/v5/application/middleware"
	"github.com/admpub/nging/v5/application/registry/navigate"

	"github.com/webx-top/echo"
)

func Project(ctx echo.Context) error {
	ident := ctx.Param(`ident`)
	partial := ctx.Formx(`partial`).Bool()
	var list navigate.List
	proj := navigate.ProjectGet(ident)
	if proj != nil {
		user := handler.User(ctx)
		if user == nil || !role.IsFounder(user) {
			permission := middleware.UserPermission(ctx)
			list = permission.FilterNavigate(ctx, proj.NavList)
		} else {
			list = *proj.NavList
		}
	}
	//echo.Dump(navigate.ProjectURLsIdent())
	data := ctx.Data()
	result := echo.H{}
	if !partial {
		result.Set(`list`, list)
	} else {
		ctx.Set(`leftNavigate`, list)
		b, err := ctx.Fetch(`sidebar_nav`, nil)
		if err != nil {
			return ctx.JSON(data.SetError(err))
		}
		result.Set(`list`, string(b))
	}
	data.SetData(result)
	return ctx.JSON(data)
}
