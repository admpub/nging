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
	"encoding/gob"
	"strings"

	GAuth "github.com/admpub/dgoogauth"
	"github.com/admpub/nging/v4/application/dbschema"
	"github.com/admpub/nging/v4/application/handler"
	"github.com/admpub/nging/v4/application/model"
	"github.com/admpub/qrcode"
	"github.com/webx-top/echo"
)

func init() {
	GAuth.Issuer = `nging`
	GAuth.Size = `300x300`
	handler.Register(func(e echo.RouteRegister) {
		e.Route("GET,POST", `/gauth_check`, GAuthCheck)
		e.Route("GET", `/qrcode`, QrCode)
	})
}

func QrCode(ctx echo.Context) error {
	data := ctx.Form("data")
	size := ctx.Form("size")
	var (
		width  = 300
		height = 300
	)
	siz := strings.SplitN(size, `x`, 2)
	switch len(siz) {
	case 2:
		if i := ctx.Atop(siz[1]).Int(); i > 0 {
			height = i
		}
		fallthrough
	case 1:
		if i := ctx.Atop(siz[0]).Int(); i > 0 {
			width = i
		}
	}
	ctx.Response().Header().Set("Content-Type", "image/png")
	return qrcode.EncodeToWriter(data, width, height, ctx.Response())
}

func init() {
	gob.Register(&GAuth.KeyData{})
}

func GAuthBind(ctx echo.Context) error {
	var err error
	user := handler.User(ctx)
	if user == nil {
		return ctx.E(`登录信息获取失败，请重新登录`)
	}
	var (
		binded bool
		u2f    *dbschema.NgingUserU2f
		typ    = `google`
	)
	m := model.NewUser(ctx)
	u2f, _ = m.U2F(user.Id, typ)
	if u2f.Id > 0 {
		binded = true
	}
	if !binded {
		if ctx.IsPost() {
			err = GAuthVerify(ctx, ``, true)
			if err == nil {
				binded = true
			}
		}
		var qrCodeUrl string
		keyData, ok := ctx.Session().Get(`GAuthKeyData`).(*GAuth.KeyData)
		if !ok {
			keyData, qrCodeUrl = GAuth.GenQrCode(user.Username, handler.URLFor("/qrcode")+"?size=%s&data=%s")
			ctx.Session().Set(`GAuthKeyData`, keyData)
		} else {
			qrCodeUrl = GAuth.QrCode(user.Username, keyData.Encoded, handler.URLFor("/qrcode")+"?size=%s&data=%s")
		}
		ctx.Set(`keyData`, keyData)
		ctx.Set(`qrCodeUrl`, qrCodeUrl)
	}
	ctx.Set(`binded`, binded)
	return ctx.Render(`gauth/bind`, handler.Err(ctx, err))
}

func GAuthCheck(ctx echo.Context) error {
	//直接从session中读取
	user, _ := ctx.Session().Get(`user`).(*dbschema.NgingUser)
	if user == nil {
		return ctx.Redirect(handler.URLFor(`/login`))
	}
	ctx.Set(`user`, user)
	var err error
	if ctx.IsPost() {
		err = GAuthVerify(ctx, ``)
		if err == nil {
			ctx.Session().Delete(`auth2ndURL`)
			next := ctx.Form(`next`)
			if len(next) == 0 {
				next = handler.URLFor(`/`)
			}
			return ctx.Redirect(next)
		}
	}
	return ctx.Render(`gauth/check`, handler.Err(ctx, err))
}

func GAuthVerify(ctx echo.Context, fieldName string, test ...bool) error {
	var keyData *GAuth.KeyData
	user := handler.User(ctx)
	if user == nil {
		return ctx.E(`登录信息获取失败，请重新登录`)
	}
	testAndBind := len(test) > 0 && test[0]
	if testAndBind {
		var ok bool
		keyData, ok = ctx.Session().Get(`GAuthKeyData`).(*GAuth.KeyData)
		if !ok {
			return ctx.E(`从session获取GAuthKeyData失败`)
		}
	} else {
		m := model.NewUser(ctx)
		u2f, err := m.U2F(user.Id, `google`)
		if err != nil && u2f.Id < 1 {
			return ctx.E(`从用户资料中获取token失败`)
		}
		keyData = &GAuth.KeyData{
			Original: u2f.Token,
			Encoded:  u2f.Extra,
		}
	}
	if len(fieldName) == 0 {
		fieldName = `code`
	}
	ok, err := GAuth.VerifyFrom(keyData, ctx.Form(fieldName))
	if !ok {
		return ctx.E(`验证码不正确`)
	}
	if err != nil {
		return err
	}
	if testAndBind {
		u2f := dbschema.NewNgingUserU2f(ctx)
		u2f.Uid = user.Id
		u2f.Token = keyData.Original
		u2f.Extra = keyData.Encoded
		u2f.Type = `google`
		_, err = u2f.Add()
	}
	return err
}
