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

package index

import (
	"github.com/admpub/nging/application/registry/navigate"

	"github.com/webx-top/com"
	"github.com/webx-top/echo"
)

func Project(ctx echo.Context) error {
	ident := ctx.Param(`ident`)
	partial := ctx.Formx(`partial`).Bool()
	proj := navigate.ProjectGet(ident)
	var list navigate.List
	if proj != nil {
		list = *proj.NavList
	}
	//echo.Dump(navigate.ProjectURLsIdent())
	data := ctx.Data()
	result := echo.H{}
	if !partial {
		result.Set(`list`, list)
	} else {
		b, err := ctx.Fetch(`sidebar_nav`, list)
		if err != nil {
			return ctx.JSON(data.SetError(err))
		}
		result.Set(`list`, com.Bytes2str(b))
	}
	data.SetData(result)
	return ctx.JSON(data)
}
