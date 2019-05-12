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
package user

import (
	"strings"

	"github.com/admpub/nging/application/handler"
	"github.com/admpub/nging/application/handler/term"
	"github.com/admpub/nging/application/library/filemanager"
	"github.com/admpub/nging/application/model"
	"github.com/webx-top/com"
	"github.com/webx-top/echo"
)

func Edit(ctx echo.Context) error {
	var err error
	user := handler.User(ctx)
	if user == nil {
		return ctx.E(`登录信息获取失败，请重新登录`)
	}
	m := model.NewUser(ctx)
	err = m.Get(nil, `id`, user.Id)
	if err != nil {
		return err
	}
	needCheckU2F := m.NeedCheckU2F(user.Id)
	if ctx.IsPost() {
		email := strings.TrimSpace(ctx.Form(`email`))
		mobile := strings.TrimSpace(ctx.Form(`mobile`))
		modifyPass := ctx.Form(`modifyPass`) == `1`

		//新密码
		newPass := strings.TrimSpace(ctx.Form(`newPass`))
		confirmPass := strings.TrimSpace(ctx.Form(`confirmPass`))

		//旧密码
		passwd := strings.TrimSpace(ctx.Form(`pass`))

		gender := strings.TrimSpace(ctx.Form(`gender`))

		if len(email) == 0 {
			err = ctx.E(`Email不能为空`)
		} else if modifyPass && len(newPass) < 8 {
			err = ctx.E(`新密码不能少于8个字符`)
		} else if modifyPass && newPass != confirmPass {
			err = ctx.E(`新密码与确认新密码不一致`)
		} else if !ctx.Validate(`email`, email, `email`).Ok() {
			err = ctx.E(`Email地址"%s"格式不正确`, email)
		} else if len(mobile) > 0 && !ctx.Validate(`mobile`, mobile, `mobile`).Ok() {
			err = ctx.E(`手机号"%s"格式不正确`, mobile)
		} else if m.User.Password != com.MakePassword(passwd, m.User.Salt) {
			err = ctx.E(`旧密码输入不正确`)
		} else if needCheckU2F {
			//两步验证码
			err = GAuthVerify(ctx, `u2fCode`)
		}
		if err == nil && !ctx.Validate(`email`, email, `email`).Ok() {
			err = ctx.E(`Email地址格式不正确`)
		}
		if err == nil {
			set := map[string]interface{}{
				`email`:  email,
				`mobile`: mobile,
				`avatar`: ctx.Form(`avatar`),
				`gender`: gender,
			}
			if modifyPass {
				set[`password`] = com.MakePassword(newPass, m.User.Salt)
			}
			err = m.Param().SetSend(set).SetArgs(`id`, user.Id).Update()
		}
		if err == nil {
			handler.SendOk(ctx, ctx.T(`修改成功`))
			m.Get(nil, `id`, user.Id)
			m.SetSession()
			return ctx.Redirect(handler.URLFor(`/user/edit`))
		}
	}
	ctx.Set(`needCheckU2F`, needCheckU2F)
	return ctx.Render(`user/edit`, handler.Err(ctx, err))
}

func AutoCompletePath(ctx echo.Context) error {
	sshAccountID := ctx.Formx(`sshAccountId`).Uint()
	if sshAccountID > 0 {
		check, _ := ctx.Funcs()[`CheckPerm`].(func(string) error)
		data := ctx.Data()
		if check == nil {
			data.SetData([]string{})
			return ctx.JSON(data)
		}
		if err := check(`manager/command_add`); err != nil {
			return err
		}
		if err := check(`manager/command_edit`); err != nil {
			return err
		}
		return term.SftpSearch(ctx, sshAccountID)
	}
	data := ctx.Data()
	prefix := ctx.Form(`query`)
	size := ctx.Formx(`size`, `10`).Int()
	var paths []string
	switch ctx.Form(`type`) {
	case `dir`:
		paths = filemanager.SearchDir(prefix, size)
	case `file`:
		paths = filemanager.SearchFile(prefix, size)
	default:
		paths = filemanager.Search(prefix, size)
	}
	data.SetData(paths)
	return ctx.JSON(data)
}
