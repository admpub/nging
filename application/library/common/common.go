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
	"fmt"

	"github.com/webx-top/captcha"
	"github.com/webx-top/echo"
	stdCode "github.com/webx-top/echo/code"
	"github.com/webx-top/echo/middleware/tplfunc"
	"github.com/webx-top/echo/subdomains"
)

// Ok 操作成功
func Ok(v string) Successor {
	return NewOk(v)
}

// Err 获取错误信息
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

// SendOk 记录成功信息
func SendOk(ctx echo.Context, msg string) {
	if ctx.IsAjax() || ctx.Format() != `html` {
		ctx.Data().SetInfo(msg, 1)
		return
	}
	ctx.Session().AddFlash(Ok(msg))
}

// SendFail 记录失败信息
func SendFail(ctx echo.Context, msg string) {
	if ctx.IsAjax() || ctx.Format() != `html` {
		ctx.Data().SetInfo(msg, 0)
		return
	}
	ctx.Session().AddFlash(msg)
}

// SendErr 记录错误信息 (SendFail的别名)
func SendErr(ctx echo.Context, err error) {
	SendFail(ctx, err.Error())
}

// VerifyCaptcha 验证码验证
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
		data.SetError(ErrCaptcha)
	}
	return data
}

// VerifyAndSetCaptcha 验证码验证并设置新验证码信息
func VerifyAndSetCaptcha(ctx echo.Context, hostAlias string, captchaName string, args ...string) echo.Data {
	data := VerifyCaptcha(ctx, hostAlias, captchaName, args...)
	if data.GetCode() != stdCode.CaptchaError {
		data.SetData(CaptchaInfo(hostAlias, captchaName, args...))
	}
	return data
}

// CaptchaInfo 新验证码信息
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

type ConfigFromDB interface {
	ConfigFromDB() echo.H
}

func MakeMap(values ...interface{}) echo.H {
	h := echo.H{}
	if len(values) == 0 {
		return h
	}
	var k string
	for i, j := 0, len(values); i < j; i++ {
		if i%2 == 0 {
			k = fmt.Sprint(values[i])
			continue
		}
		h.Set(k, values[i])
		k = ``
	}
	if len(k) > 0 {
		h.Set(k, nil)
	}
	return h
}
