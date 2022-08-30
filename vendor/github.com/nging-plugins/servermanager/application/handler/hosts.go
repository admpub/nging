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

package handler

import (
	"os"

	"github.com/webx-top/com"
	"github.com/webx-top/echo"

	"github.com/admpub/nging/v4/application/handler"

	"github.com/nging-plugins/servermanager/application/library/hosts"
)

func Hosts(ctx echo.Context) error {
	var err error
	if ctx.IsPost() {
		content := ctx.Form(`hosts`)
		err = hosts.WriteFile(com.Str2bytes(content))
		data := ctx.Data()
		if err != nil {
			data.SetError(err)
		} else {
			data.SetInfo(ctx.T(`操作成功`))
		}
		return ctx.JSON(data)
	}
	ctx.Set(`title`, ctx.T(`hosts文件编辑`))
	hostsPath := hosts.Path()
	ctx.Set(`path`, hostsPath)
	var perm string
	if fi, err := os.Stat(hostsPath); err == nil {
		perm = fi.Mode().String()
	}
	ctx.Set(`perm`, perm)
	var b []byte
	b, err = hosts.ReadFile()
	if err == nil {
		ctx.Request().Form().Set(`hosts`, com.Bytes2str(b))
	}
	return ctx.Render(`server/hosts`, handler.Err(ctx, err))
}
