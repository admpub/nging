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
	"fmt"

	"github.com/webx-top/echo"

	"github.com/admpub/events"
	"github.com/coscms/webcore/library/backend"
	"github.com/coscms/webcore/library/backend/oauth2client"
	"github.com/coscms/webcore/library/captcha/captchabiz"
	"github.com/coscms/webcore/library/common"
	"github.com/coscms/webcore/library/config"
	"github.com/coscms/webcore/library/httpserver"
	"github.com/coscms/webcore/library/license"
	"github.com/coscms/webcore/library/nerrors"
	"github.com/coscms/webcore/middleware"
	"github.com/coscms/webcore/model"
	"github.com/coscms/webcore/request"
	stdCode "github.com/webx-top/echo/code"
	"github.com/webx-top/echo/handler/oauth2"
	"github.com/webx-top/echo/param"
)

func Index(ctx echo.Context) error {
	user := backend.User(ctx)
	if user == nil {
		return ctx.Redirect(backend.URLFor(`/login`, true))
	}
	dashboard := httpserver.Backend.Dashboard
	var err error
	ctx.Set(`cards`, dashboard.Cards.All(ctx).Build(ctx))
	blocks := dashboard.Blocks.All(ctx)
	if err = blocks.Ready(ctx); err != nil {
		return err
	}
	ctx.Set(`blocks`, blocks)
	ctx.Set(`license`, license.License())
	ctx.Set(`showExpirationTime`, config.FromFile().Sys.ShowExpirationTime)
	return ctx.Render(`index`, common.Err(ctx, err))
}

func Login(ctx echo.Context) error {
	//检查是否已安装
	if !config.IsInstalled() {
		return ctx.Redirect(backend.URLFor(`/setup`, true))
	}

	next := ctx.Form(`next`)
	if len(next) == 0 {
		next = backend.URLFor(`/index`, true)
	}

	user := backend.User(ctx)
	if user != nil {
		return ctx.Redirect(next)
	}
	var err error
	if ctx.IsPost() {
		if data := captchabiz.VerifyCaptcha(ctx, `backend`, `code`); data.GetCode() == stdCode.CaptchaError {
			err = nerrors.ErrCaptcha.SetMessage(param.AsString(data.GetInfo()))
		} else if data.GetCode() != stdCode.Success {
			err = fmt.Errorf("%v", data.GetInfo())
		} else {
			err = middleware.Auth(ctx)
			if err == nil {
				return ctx.Redirect(next)
			}
		}
	}
	ctx.SetFunc(`oAuthAccounts`, func() []oauth2.Account {
		return oauth2client.GetOAuthAccounts(true)
	})
	return ctx.Render(`login`, common.Err(ctx, err))
}

func Register(ctx echo.Context) error {
	var err error
	var req *request.Register
	if ctx.IsPost() {
		c := model.NewCode(ctx)
		m := model.NewUser(ctx)
		req = echo.GetValidated(ctx).(*request.Register)
		code := req.InvitationCode
		user := req.Username
		email := req.Email
		passwd := req.Password
		err = ctx.Begin()
		if err != nil {
			goto END
		}
		err = c.VerfyInvitationCode(code)
		if err != nil {
			ctx.Rollback()
			goto END
		}
		err = m.Register(user, passwd, email, c.Invitation.RoleIds)
		if err != nil {
			ctx.Rollback()
			goto END
		}
		c.UseInvitationCode(c.Invitation, m.NgingUser.Id)
		m.SetSession()
		err = m.NgingUser.UpdateField(nil, `session_id`, ctx.Session().ID(), `id`, m.NgingUser.Id)
		if err != nil {
			ctx.Rollback()
			m.UnsetSession()
			goto END
		}

		loginLogM := model.NewLoginLog(ctx)
		loginLogM.OwnerType = `user`
		loginLogM.Username = user
		loginLogM.SessionId = ctx.Session().ID()
		loginLogM.Success = `Y`
		loginLogM.AddAndSaveSession()

		err = echo.FireByNameWithMap(`nging.user.register.success`, events.Map{`user`: m.NgingUser})
		if err != nil {
			ctx.Rollback()
			m.UnsetSession()
			goto END
		}

		ctx.Commit()

		next := ctx.Query(`next`)
		if len(next) == 0 {
			next = backend.URLFor(`/index`)
		}
		return ctx.Redirect(next)
	}

END:
	if req == nil {
		req = &request.Register{}
	}
	ctx.Set(`user`, req)
	return ctx.Render(`register`, common.Err(ctx, err))
}

func Logout(ctx echo.Context) error {
	ctx.Session().Delete(`user`)
	user := backend.User(ctx)
	if user != nil {
		err := echo.FireByNameWithMap(`nging.user.logout.success`, events.Map{
			`user`:     user,
			`username`: user.Username,
			`uid`:      user.Id,
		})
		if err != nil {
			return err
		}
	}
	return ctx.Redirect(backend.URLFor(`/login`))
}
