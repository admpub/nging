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

package user

import (
	"strings"

	"github.com/webx-top/com"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/code"

	"github.com/admpub/nging/v5/application/handler"
	"github.com/admpub/nging/v5/application/library/codec"
	"github.com/admpub/nging/v5/application/library/filemanager"
	"github.com/admpub/nging/v5/application/model"
)

func Edit(ctx echo.Context) error {
	var err error
	user := handler.User(ctx)
	if user == nil {
		return ctx.NewError(code.Unauthenticated, `登录信息获取失败，请重新登录`)
	}
	m := model.NewUser(ctx)
	err = m.Get(nil, `id`, user.Id)
	if err != nil {
		return err
	}
	needCheckU2F, err := m.NeedCheckU2F(model.AuthTypePassword, user.Id, 2)
	if err != nil {
		return err
	}
	if ctx.IsPost() {
		email := strings.TrimSpace(ctx.Form(`email`))
		mobile := strings.TrimSpace(ctx.Form(`mobile`))
		modifyPass := ctx.Form(`modifyPass`) == `1`

		//新密码
		newPass := strings.TrimSpace(ctx.Form(`newPass`))
		confirmPass := strings.TrimSpace(ctx.Form(`confirmPass`))

		//旧密码
		passwd := strings.TrimSpace(ctx.Form(`pass`))

		passwd, err = codec.DefaultSM2DecryptHex(passwd)
		if err != nil {
			return ctx.NewError(code.InvalidParameter, `旧密码解密失败: %v`, err).SetZone(`pass`)
		}
		if modifyPass {
			newPass, err = codec.DefaultSM2DecryptHex(newPass)
			if err != nil {
				return ctx.NewError(code.InvalidParameter, `新密码解密失败: %v`, err).SetZone(`newPass`)
			}
			confirmPass, err = codec.DefaultSM2DecryptHex(confirmPass)
			if err != nil {
				return ctx.NewError(code.InvalidParameter, `您输入的确认密码解密失败: %v`, err).SetZone(`confirmPass`)
			}
		}

		gender := strings.TrimSpace(ctx.Form(`gender`))

		if len(email) == 0 {
			err = ctx.NewError(code.InvalidParameter, `Email不能为空`).SetZone(`email`)
		} else if modifyPass && len(newPass) < 8 {
			err = ctx.NewError(code.InvalidParameter, `新密码不能少于8个字符`).SetZone(`newPass`)
		} else if modifyPass && newPass != confirmPass {
			err = ctx.NewError(code.InvalidParameter, `新密码与确认新密码不一致`).SetZone(`confirmPass`)
		} else if ctx.Validate(`email`, email, `email`) != nil {
			err = ctx.NewError(code.InvalidParameter, `Email地址"%s"格式不正确`, email).SetZone(`email`)
		} else if len(mobile) > 0 && ctx.Validate(`mobile`, mobile, `mobile`) != nil {
			err = ctx.NewError(code.InvalidParameter, `手机号"%s"格式不正确`, mobile).SetZone(`mobile`)
		} else if m.NgingUser.Password != com.MakePassword(passwd, m.NgingUser.Salt) {
			err = ctx.NewError(code.InvalidParameter, `旧密码输入不正确`).SetZone(`email`)
		} else if needCheckU2F {
			//两步验证码
			err = GAuthVerify(ctx, `u2fCode`)
		}
		if err == nil && ctx.Validate(`email`, email, `email`) != nil {
			err = ctx.NewError(code.InvalidParameter, `Email地址格式不正确`).SetZone(`email`)
		}
		if err == nil {
			set := map[string]interface{}{
				`email`:  email,
				`mobile`: mobile,
				`avatar`: ctx.Form(`avatar`),
				`gender`: gender,
			}
			if modifyPass {
				set[`password`] = com.MakePassword(newPass, m.NgingUser.Salt)
			}
			err = m.UpdateFields(nil, set, `id`, user.Id)
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

var onAutoCompletePath = []func(echo.Context) (bool, error){}

func OnAutoCompletePath(fn func(echo.Context) (bool, error)) {
	onAutoCompletePath = append(onAutoCompletePath, fn)
}

func FireAutoCompletePath(c echo.Context) (bool, error) {
	for _, fn := range onAutoCompletePath {
		ok, err := fn(c)
		if ok || err != nil {
			return true, err
		}
	}
	return false, nil
}

func AutoCompletePath(ctx echo.Context) error {
	user := handler.User(ctx)
	if user == nil {
		return ctx.NewError(code.Unauthenticated, `登录信息获取失败，请重新登录`)
	}
	if ok, err := FireAutoCompletePath(ctx); ok || err != nil {
		return err
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
