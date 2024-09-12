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
	"encoding/gob"
	"strings"

	GAuth "github.com/admpub/dgoogauth"
	"github.com/coscms/webcore/dbschema"
	"github.com/coscms/webcore/library/backend"
	"github.com/coscms/webcore/library/common"
	"github.com/coscms/webcore/model"
	"github.com/coscms/webcore/registry/route"
	"github.com/webx-top/db"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/code"
)

func init() {
	GAuth.Issuer = `nging`
	GAuth.Size = `300x300`
	route.Register(func(e echo.RouteRegister) {
		e.Route("GET,POST", `/gauth_check`, GAuthCheck)
		e.Route("GET", `/qrcode`, backend.QrCode)
	})
}

func init() {
	gob.Register(&GAuth.KeyData{})
}

func GAuthBind(ctx echo.Context) error {
	var err error
	user := backend.User(ctx)
	if user == nil {
		return ctx.NewError(code.Unauthenticated, `登录信息获取失败，请重新登录`)
	}
	var (
		binded bool
		u2f    *dbschema.NgingUserU2f
		typ         = `google`
		step   uint = 2
	)
	m := model.NewUser(ctx)
	u2f, err = m.U2F(user.Id, typ, step)
	if err != nil {
		if err != db.ErrNoMoreRows {
			return err
		}
		err = nil
	}
	if u2f.Id > 0 {
		binded = true
	}
	// binded = true // for test only
	if !binded {
		if ctx.IsPost() {
			err = gAuthBind(ctx)
			if err == nil {
				return ctx.Redirect(backend.URLFor(`/user/gauth_bind`))
			}
		}
		keyData, ok := ctx.Session().Get(`GAuthKeyData`).(*GAuth.KeyData)
		if !ok {
			keyData = GAuth.GenKeyData()
			ctx.Session().Set(`GAuthKeyData`, keyData)
		}
		qrCodeUrl := GAuth.QrCode(user.Username, keyData.Encoded, backend.URLFor("/qrcode")+"?size=%s&data=%s")
		ctx.Set(`keyData`, keyData)
		ctx.Set(`qrCodeUrl`, qrCodeUrl)
	} else {
		if ctx.IsPost() {
			operation := ctx.Form(`operation`)
			switch operation {
			case `unbind`:
				err = gAuthUnbind(ctx, user.Id, typ, step)
			case `modify`:
				precondition := strings.Join(ctx.FormValues(`precondition`), `,`)
				err = gAuthUpdatePrecondition(ctx, user.Id, typ, step, precondition)
			}
			if err == nil {
				return ctx.Redirect(backend.URLFor(`/user/gauth_bind`))
			}
		} else {
			ctx.Request().Form().Set(`precondition`, u2f.Precondition)
		}
	}
	ctx.Set(`binded`, binded)
	ctx.Set(`activeSafeItem`, `gauth_bind`)
	ctx.Set(`safeItems`, model.SafeItems.Slice())
	ctx.SetFunc(`getSafeItemName`, model.SafeItems.Get)
	ctx.Set(`step1SafeItems`, model.ListSafeItemsByStep(1, model.AuthTypePassword))
	return ctx.Render(`gauth/bind`, common.Err(ctx, err))
}

func gAuthBind(ctx echo.Context) error {
	return GAuthVerify(ctx, ``, true)
}

func gAuthUnbind(ctx echo.Context, uid uint, typ string, step uint) error {
	err := GAuthVerify(ctx, ``)
	if err == nil {
		u2f := model.NewUserU2F(ctx)
		err = u2f.Unbind(uid, typ, step)
	}
	return err
}

func gAuthUpdatePrecondition(ctx echo.Context, uid uint, typ string, step uint, precondition string) error {
	err := GAuthVerify(ctx, ``)
	if err == nil {
		u2f := model.NewUserU2F(ctx)
		err = u2f.UpdatePrecondition(uid, typ, step, precondition)
	}
	return err
}

func GAuthCheck(ctx echo.Context) error {
	//直接从session中读取
	user, _ := ctx.Session().Get(`user`).(*dbschema.NgingUser)
	if user == nil {
		return ctx.Redirect(backend.URLFor(`/login`))
	}
	ctx.Set(`user`, user)
	var err error
	if ctx.IsPost() {
		err = GAuthVerify(ctx, ``)
		if err == nil {
			ctx.Session().Delete(`auth2ndURL`)
			next := ctx.Form(`next`)
			if len(next) == 0 {
				next = backend.URLFor(`/`)
			}
			return ctx.Redirect(next)
		}
	}
	return ctx.Render(`gauth/check`, common.Err(ctx, err))
}

func GAuthVerify(ctx echo.Context, fieldName string, test ...bool) error {
	var keyData *GAuth.KeyData
	user := backend.User(ctx)
	if user == nil {
		return ctx.NewError(code.Unauthenticated, `登录信息获取失败，请重新登录`)
	}
	testAndBind := len(test) > 0 && test[0]
	if testAndBind {
		var ok bool
		keyData, ok = ctx.Session().Get(`GAuthKeyData`).(*GAuth.KeyData)
		if !ok {
			return ctx.NewError(code.Failure, `从session获取GAuthKeyData失败`)
		}
	} else {
		m := model.NewUser(ctx)
		u2f, err := m.U2F(user.Id, `google`, 2)
		if err != nil && u2f.Id < 1 {
			return ctx.NewError(code.Failure, `从用户资料中获取token失败`)
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
		return ctx.NewError(code.InvalidParameter, `验证码不正确`)
	}
	if err != nil {
		return err
	}
	if testAndBind {
		u2f := model.NewUserU2F(ctx)
		u2f.Uid = user.Id
		u2f.Name = `Two-factor authentication`
		u2f.Token = keyData.Original
		u2f.Extra = keyData.Encoded
		u2f.Type = `google`
		u2f.Step = 2
		u2f.Precondition = strings.Join(ctx.FormValues(`precondition`), `,`)
		_, err = u2f.Add()
	}
	return err
}
