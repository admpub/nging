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

package common

import (
	stdErr "errors"

	"github.com/webx-top/captcha"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/middleware/tplfunc"
	"github.com/webx-top/echo/subdomains"
)

func Ok(v string) Successor {
	return NewOk(v)
}

func Err(ctx echo.Context, err error) (ret interface{}) {
	if err == nil {
		flash := ctx.Flash()
		if flash != nil {
			if errMsg, ok := flash.(string); ok {
				ret = stdErr.New(errMsg)
			} else {
				ret = flash
			}
		}
	} else {
		ret = err
	}
	return
}

func SendOk(ctx echo.Context, msg string) {
	ctx.Session().AddFlash(Ok(msg))
}

func SendFail(ctx echo.Context, msg string) {
	ctx.Session().AddFlash(msg)
}

func SendErr(ctx echo.Context, err error) {
	SendFail(ctx, err.Error())
}

func VerifyCaptcha(ctx echo.Context, hostAlias string, captchaName string, args ...string) echo.Data {
	data := ctx.Data()
	idGet := ctx.Form
	if len(args) > 0 {
		idGet = func(_ string, defaults ...string) string {
			return ctx.Form(args[0], defaults...)
		}
	}
	code := ctx.Form(captchaName)
	var err error
	if len(code) == 0 {
		err = ctx.E(`请输入验证码`)
	} else if !tplfunc.CaptchaVerify(code, idGet) {
		err = ctx.E(`验证码不正确`)
	}
	if err != nil {
		data.SetZone(captchaName)
		data.SetData(CaptchaInfo(hostAlias, captchaName, args...))
		data.SetError(err, StatusCaptchaError)
	}
	return data
}

func VerifyAndSetCaptcha(ctx echo.Context, hostAlias string, captchaName string, args ...string) echo.Data {
	data := VerifyCaptcha(ctx, hostAlias, captchaName, args...)
	if data.GetCode() != StatusCaptchaError {
		data.SetData(CaptchaInfo(hostAlias, captchaName, args...))
	}
	return data
}

func CaptchaInfo(hostAlias string, captchaName string, args ...string) echo.H {
	captchaID := captcha.New()
	captchaIdent := `captchaId`
	if len(args) > 0 {
		captchaIdent = args[0]
	}
	return echo.H{
		`captchaName`:  captchaName,
		`captchaIdent`: captchaIdent,
		`captchaID`:    captchaID,
		`captchaURL`:   subdomains.Default.URL(`/captcha/`+captchaID+`.png`, hostAlias),
	}
}
