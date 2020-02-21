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
	"strings"

	"github.com/webx-top/echo"
	"github.com/webx-top/echo/middleware/tplfunc"

	"github.com/admpub/nging/application/handler"
	"github.com/admpub/nging/application/library/common"
	"github.com/admpub/nging/application/library/config"
	"github.com/admpub/nging/application/library/license"
	"github.com/admpub/nging/application/middleware"
	"github.com/admpub/nging/application/model"
	"github.com/admpub/nging/application/registry/dashboard"
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

	common.SetRandomSecret(ctx, `loginPassword`, `passwordSecrect`)
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
		common.DecryptedByRandomSecret(ctx, `registerPassword`, &passwd, &repwd)
		if len(code) == 0 {
			err = ctx.E(`邀请码不能为空`)
			goto END
		}
		if repwd != passwd {
			err = ctx.E(`密码与确认密码不一致`)
			goto END
		}
		err = c.VerfyInvitationCode(code)
		if err != nil {
			goto END
		}
		err = m.Register(user, passwd, email, c.Invitation.RoleIds)
		if err != nil {
			goto END
		}
		c.UseInvitationCode(c.Invitation, m.NgingUser.Id)
		common.DeleteRandomSecret(ctx, `registerPassword`)
		m.SetSession()
		returnTo := ctx.Query(`return_to`)
		if len(returnTo) == 0 {
			returnTo = handler.URLFor(`/index`)
		}
		return ctx.Redirect(returnTo)
	}

END:
	common.SetRandomSecret(ctx, `registerPassword`, `passwordSecrect`)
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
