package user

import (
	"errors"

	"encoding/gob"

	QR "github.com/RaymondChou/goqr/pkg"
	GAuth "github.com/admpub/dgoogauth"
	"github.com/admpub/nging/application/dbschema"
	"github.com/admpub/nging/application/handler"
	"github.com/admpub/nging/application/middleware"
	"github.com/admpub/nging/application/model"
	"github.com/webx-top/echo"
)

func init() {
	GAuth.Issuer = `nging`
	handler.Register(func(e *echo.Echo) {
		e.Route("GET,POST", `/gauth_bind`, GAuthBind, middleware.AuthCheck)
		e.Route("GET,POST", `/gauth_check`, GAuthCheck)
		e.Route("GET", `/qrcode`, QrCode)
	})
}

func QrCode(ctx echo.Context) error {
	data := ctx.Form("data")
	size := QR.L
	switch ctx.Form("size") {
	case `M`:
		size = QR.M
	case `Q`:
		size = QR.Q
	case `H`:
		size = QR.H
	}
	c, err := QR.Encode(data, size)
	if err != nil {
		return err
	}
	pngBytes := c.PNG()
	ctx.Response().Header().Set("Content-Type", "image/png")
	return ctx.Blob(pngBytes)
}

func init() {
	gob.Register(&GAuth.KeyData{})
}

func GAuthBind(ctx echo.Context) error {
	var err error
	user := handler.User(ctx)
	if user == nil {
		return errors.New(ctx.T(`登录信息获取失败，请重新登录`))
	}
	var (
		binded bool
		u2f    *dbschema.UserU2f
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
			keyData, qrCodeUrl = GAuth.GenQrCode(user.Username, "/qrcode?size=%s&data=%s")
			ctx.Session().Set(`GAuthKeyData`, keyData)
		} else {
			qrCodeUrl = GAuth.QrCode(user.Username, keyData.Encoded, "/qrcode?size=%s&data=%s")
		}
		ctx.Set(`keyData`, keyData)
		ctx.Set(`qrCodeUrl`, qrCodeUrl)
	}
	ctx.Set(`binded`, binded)
	return ctx.Render(`gauth/bind`, handler.Err(ctx, err))
}

func GAuthCheck(ctx echo.Context) error {
	//直接从session中读取
	user, _ := ctx.Session().Get(`user`).(*dbschema.User)
	if user == nil {
		return ctx.Redirect(`/login`)
	}
	ctx.Set(`user`, user)
	var err error
	if ctx.IsPost() {
		err = GAuthVerify(ctx, ``)
		if err == nil {
			ctx.Session().Delete(`auth2ndURL`)
			returnTo := ctx.Form(`return_to`)
			if len(returnTo) == 0 {
				returnTo = `/`
			}
			return ctx.Redirect(returnTo)
		}
	}
	return ctx.Render(`gauth/check`, handler.Err(ctx, err))
}

func GAuthVerify(ctx echo.Context, fieldName string, test ...bool) error {
	var keyData *GAuth.KeyData
	user := handler.User(ctx)
	if user == nil {
		return errors.New(ctx.T(`登录信息获取失败，请重新登录`))
	}
	testAndBind := len(test) > 0 && test[0]
	if testAndBind {
		var ok bool
		keyData, ok = ctx.Session().Get(`GAuthKeyData`).(*GAuth.KeyData)
		if !ok {
			return errors.New(ctx.T(`从session获取GAuthKeyData失败`))
		}
	} else {
		m := model.NewUser(ctx)
		u2f, err := m.U2F(user.Id, `google`)
		if err != nil && u2f.Id < 1 {
			return errors.New(ctx.T(`从用户资料中获取token失败`))
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
		return errors.New(ctx.T(`验证码不正确`))
	}
	if err != nil {
		return err
	}
	if testAndBind {
		u2f := &dbschema.UserU2f{}
		u2f.Uid = user.Id
		u2f.Token = keyData.Original
		u2f.Extra = keyData.Encoded
		u2f.Type = `google`
		_, err = u2f.Add()
	}
	return err
}
