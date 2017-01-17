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
package handler

import (
	"errors"

	"strings"

	"github.com/admpub/nging/application/library/config"
	"github.com/admpub/nging/application/middleware"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/middleware/tplfunc"
)

func Index(ctx echo.Context) error {
	return ctx.Redirect(`/manage`)
}

func Login(ctx echo.Context) error {
	if user, _ := ctx.Get(`user`).(string); len(user) > 0 {
		returnTo := ctx.Query(`return_to`)
		if len(returnTo) == 0 {
			returnTo = `/manage`
		}
		return ctx.Redirect(returnTo)
	}
	var err error
	if ctx.IsPost() {
		if !tplfunc.CaptchaVerify(ctx.Form(`code`), ctx.Form) {
			err = errors.New(ctx.T(`验证码不正确`))
		} else {
			err = middleware.Auth(ctx, true)
			if err == nil {
				returnTo := ctx.Form(`return_to`)
				if len(returnTo) == 0 {
					returnTo = `/manage`
				}
				return ctx.Redirect(returnTo)
			}
		}
	}

	//检查是否已安装
	if !config.IsInstalled() {
		return ctx.Redirect(`/setup`)
	}

	return ctx.Render(`login`, Err(ctx, err))
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
