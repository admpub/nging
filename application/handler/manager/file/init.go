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

package file

import (
	"fmt"
	"time"

	"github.com/webx-top/echo"

	"github.com/admpub/nging/application/handler"
	"github.com/admpub/nging/application/middleware"
	_ "github.com/admpub/nging/application/model/file/initialize"
	"github.com/admpub/nging/application/registry/upload"
	"github.com/admpub/nging/application/registry/upload/table"
)

func init() {
	handler.RegisterToGroup(`/manager`, func(g echo.RouteRegister) {
		r := g.Group(`/file`)
		r.Route(`GET,POST`, `/list`, FileList)
		r.Route(`GET,POST`, `/delete/:id`, FileDelete)
	})
	handler.Register(func(r echo.RouteRegister) {
		r.Route(`GET,POST`, `/finder`, Finder, middleware.AuthCheck)
	})

	// 用户上传个人文件时的文件命名方式
	upload.CheckerRegister(`user`, func(ctx echo.Context, tis table.TableInfoStorer) (subdir string, name string, err error) {
		user := handler.User(ctx)
		if user == nil {
			err = ctx.E(`登录信息获取失败，请重新登录`)
			return
		}
		userID := uint64(user.Id)
		timestamp := ctx.Formx(`time`).Int64()
		// 验证签名（避免上传接口被滥用）
		if ctx.Form(`token`) != upload.Token(ctx.Queries()) {
			err = ctx.E(`令牌错误`)
			return
		}
		if time.Now().Local().Unix()-timestamp > upload.UploadLinkLifeTime {
			err = ctx.E(`上传网址已过期`)
			return
		}
		uid := fmt.Sprint(userID)
		subdir = uid + `/`
		subdir += time.Now().Format(`2006/01/02/`)
		tis.SetTableID(uid)
		tis.SetTableName(`user`)
		tis.SetFieldName(``)
		return
	})
}
