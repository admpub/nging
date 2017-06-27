/*

   Copyright 2016 Wenhui Shen <www.webx.top>

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.

*/
package index

import (
	"errors"

	"strings"

	"github.com/admpub/nging/application/handler"
	"github.com/admpub/nging/application/library/config"
	"github.com/admpub/nging/application/middleware"
	"github.com/admpub/nging/application/model"
	"github.com/webx-top/com"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/middleware/tplfunc"
)

func init() {
	handler.Register(func(e *echo.Echo) {
		e.Route("GET", `/`, Index)
		e.Route("GET,POST", `/login`, Login)
		e.Route("GET,POST", `/register`, Register)
		e.Route("GET", `/logout`, Logout)
		e.Route("GET", `/donation`, Donation)
	})
}

func Index(ctx echo.Context) error {
	return ctx.Redirect(`/manage`)
}

func Login(ctx echo.Context) error {
	//检查是否已安装
	if !config.IsInstalled() {
		return ctx.Redirect(`/setup`)
	}

	returnTo := ctx.Form(`return_to`)
	if len(returnTo) == 0 {
		returnTo = `/manage`
	}

	user := handler.User(ctx)
	if user != nil {
		return ctx.Redirect(returnTo)
	}
	var err error
	if ctx.IsPost() {
		if !tplfunc.CaptchaVerify(ctx.Form(`code`), ctx.Form) {
			err = errors.New(ctx.T(`验证码不正确`))
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
			err = errors.New(ctx.T(`邀请码不能为空`))
		} else if len(user) == 0 {
			err = errors.New(ctx.T(`用户名不能为空`))
		} else if len(email) == 0 {
			err = errors.New(ctx.T(`Email不能为空`))
		} else if len(passwd) < 8 {
			err = errors.New(ctx.T(`密码不能少于8个字符`))
		} else if repwd != passwd {
			err = errors.New(ctx.T(`密码与确认密码不一致`))
		} else if !com.IsUsername(user) {
			err = errors.New(ctx.T(`用户名不能包含特殊字符(只能由字母、数字、下划线和汉字组成)`))
		} else {
			var exists bool
			exists, err = m.Exists(user)
			if exists {
				err = errors.New(ctx.T(`用户名已经存在`))
			}
			if err == nil {
				err = c.VerfyInvitationCode(code)
			}
			if err == nil && !ctx.Validate(`email`, email, `email`).Ok() {
				err = errors.New(ctx.T(`Email地址格式不正确`))
			}
		}
		if err == nil {
			err = m.Register(user, passwd, email)
			if err == nil {
				c.UseInvitationCode(c.Invitation, m.User.Id)
			}
		}
		if err == nil {
			ctx.Session().Set(`user`, m.User)
			returnTo := ctx.Query(`return_to`)
			if len(returnTo) == 0 {
				returnTo = `/manage`
			}
			return ctx.Redirect(returnTo)
		}
	}
	return ctx.Render(`register`, handler.Err(ctx, err))
}

func Logout(ctx echo.Context) error {
	ctx.Session().Delete(`user`)
	return ctx.Redirect(`/login`)
}

func Donation(ctx echo.Context) error {
	var langSuffix string
	lang := strings.ToLower(ctx.Lang())
	if strings.HasPrefix(lang, `zh`) {
		langSuffix = `_zh-CN`
	}
	return ctx.Redirect(`/public/donation` + langSuffix + `.html`)
}
