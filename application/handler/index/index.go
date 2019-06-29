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
	"errors"
	"strings"

	"github.com/admpub/nging/application/handler"
	"github.com/admpub/nging/application/library/config"
	"github.com/admpub/nging/application/library/license"
	"github.com/admpub/nging/application/middleware"
	"github.com/admpub/nging/application/model"
	"github.com/admpub/nging/application/registry/dashboard"
	"github.com/webx-top/com"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/middleware/tplfunc"
)

func Index(ctx echo.Context) error {
	user := handler.User(ctx)
	if user == nil {
		return ctx.Redirect(handler.URLFor(`/login`))
	}
	var err error
	ctx.Set(`cards`, dashboard.CardAll().Build(ctx))
	blocks := dashboard.BlockAll()
	if err = blocks.Ready(ctx); err != nil {
		return err
	}
	ctx.Set(`blocks`, blocks)
	ctx.Set(`license`, license.License())
	ctx.Set(`showExpirationTime`, config.DefaultConfig.Sys.ShowExpirationTime)
	machineID, _ := license.MachineID()
	ctx.Set(`productURL`, license.ProductURL()+config.Version.Number+`/`+machineID)
	return ctx.Render(`index`, handler.Err(ctx, err))
}

func Login(ctx echo.Context) error {
	//检查是否已安装
	if !config.IsInstalled() {
		return ctx.Redirect(handler.URLFor(`/setup`))
	}

	returnTo := ctx.Form(`return_to`)
	if len(returnTo) == 0 {
		returnTo = handler.URLFor(`/index`)
	}

	user := handler.User(ctx)
	if user != nil {
		return ctx.Redirect(returnTo)
	}
	var err error
	if ctx.IsPost() {
		if !tplfunc.CaptchaVerify(ctx.Form(`code`), ctx.Form) {
			err = ctx.E(`验证码不正确`)
		} else {
			err = middleware.Auth(ctx, true)
			if err == nil {
				return ctx.Redirect(returnTo)
			}
		}
	}

	return ctx.Render(`login`, handler.Err(ctx, err))
}

func Register(ctx echo.Context) error {
	var err error
	if ctx.IsPost() {
		c := model.NewCode(ctx)
		m := model.NewUser(ctx)
		code := ctx.Form(`invitationCode`)
		user := ctx.Form(`username`)
		email := ctx.Form(`email`)
		passwd := ctx.Form(`password`)
		repwd := ctx.Form(`confirmationPassword`)
		if len(code) == 0 {
			err = ctx.E(`邀请码不能为空`)
		} else if len(user) == 0 {
			err = ctx.E(`用户名不能为空`)
		} else if len(email) == 0 {
			err = ctx.E(`Email不能为空`)
		} else if len(passwd) < 8 {
			err = ctx.E(`密码不能少于8个字符`)
		} else if repwd != passwd {
			err = ctx.E(`密码与确认密码不一致`)
		} else if !com.IsUsername(user) {
			err = errors.New(ctx.T(`用户名不能包含特殊字符(只能由字母、数字、下划线和汉字组成)`))
		} else if !ctx.Validate(`email`, email, `email`).Ok() {
			err = ctx.E(`Email地址格式不正确`)
		} else {
			var exists bool
			exists, err = m.Exists(user)
			if exists {
				err = ctx.E(`用户名已经存在`)
			}
			if err == nil {
				err = c.VerfyInvitationCode(code)
			}
			if err == nil && !ctx.Validate(`email`, email, `email`).Ok() {
				err = ctx.E(`Email地址格式不正确`)
			}
		}
		if err == nil {
			err = m.Register(user, passwd, email)
			if err == nil {
				c.UseInvitationCode(c.Invitation, m.User.Id)
			}
		}
		if err == nil {
			m.SetSession()
			returnTo := ctx.Query(`return_to`)
			if len(returnTo) == 0 {
				returnTo = handler.URLFor(`/index`)
			}
			return ctx.Redirect(returnTo)
		}
	}
	return ctx.Render(`register`, handler.Err(ctx, err))
}

func Logout(ctx echo.Context) error {
	ctx.Session().Delete(`user`)
	return ctx.Redirect(handler.URLFor(`/login`))
}

func Donation(ctx echo.Context) error {
	var langSuffix string
	lang := strings.ToLower(ctx.Lang())
	if strings.HasPrefix(lang, `zh`) {
		langSuffix = `_zh-CN`
	}
	return ctx.Render(`index/donation/donation`+langSuffix, nil)
}
