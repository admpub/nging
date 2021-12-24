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
	"io"

	"github.com/webx-top/echo"
)

func Ping(ctx echo.Context) error {
	header := ctx.Request().Header()
	body := ctx.Request().Body()
	b, _ := io.ReadAll(body)
	body.Close()
	r := echo.H{
		`header`: header.Object(),
		`form`:   echo.NewMapx(ctx.Request().Form().All()).AsStore(),
		`body`:   string(b),
	}
	data := ctx.Data()
	data.SetData(r)
	callback := ctx.Form(`callback`)
	if len(callback) > 0 {
		return ctx.JSONP(callback, data)
	}
	return ctx.JSON(data)
}
