package user

import (
	"errors"

	"encoding/gob"

	QR "github.com/RaymondChou/goqr/pkg"
	GAuth "github.com/admpub/dgoogauth"
	"github.com/admpub/nging/application/library/config"
	"github.com/webx-top/echo"
)

func init() {
	GAuth.Issuer = `nging`
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
	user, ok := ctx.Get(`user`).(string)
	if !ok {
		return errors.New(ctx.T(`登录信息获取失败，请重新登录`))
	}
	var binded bool
	if profile, ok := config.DefaultConfig.Sys.Accounts[user]; ok && profile.GAuthKey != nil {
		binded = true
	}
	if !binded {
		if ctx.IsPost() {
			err = GAuthVerify(ctx, true)
			if err == nil {
				binded = true
			}
		}
		var qrCodeUrl string
		keyData, ok := ctx.Session().Get(`GAuthKeyData`).(*GAuth.KeyData)
		if !ok {
			keyData, qrCodeUrl = GAuth.GenQrCode(user, "/qrcode?size=%s&data=%s")
			ctx.Session().Set(`GAuthKeyData`, keyData)
		} else {
			qrCodeUrl = GAuth.QrCode(user, keyData.Encoded, "/qrcode?size=%s&data=%s")
		}
		ctx.Set(`keyData`, keyData)
		ctx.Set(`qrCodeUrl`, qrCodeUrl)
	}
	ctx.Set(`binded`, binded)
	return ctx.Render(`gauth/bind`, Err(ctx, err))
}

func GAuthCheck(ctx echo.Context) error {
	user, ok := ctx.Session().Get(`user`).(string)
	if !ok {
		return ctx.Redirect(`/login`)
	}
	ctx.Set(`user`, user)
	var err error
	if ctx.IsPost() {
		err = GAuthVerify(ctx)
		if err == nil {
			ctx.Session().Delete(`auth2ndURL`)
			returnTo := ctx.Form(`return_to`)
			if len(returnTo) == 0 {
				returnTo = `/`
			}
			return ctx.Redirect(returnTo)
		}
	}
	return ctx.Render(`gauth/check`, Err(ctx, err))
}

func GAuthVerify(ctx echo.Context, test ...bool) error {
	var keyData *GAuth.KeyData
	user, ok := ctx.Get(`user`).(string)
	if !ok {
		return errors.New(ctx.T(`登录信息获取失败，请重新登录`))
	}
	testAndBind := len(test) > 0 && test[0]
	if testAndBind {
		keyData, ok = ctx.Session().Get(`GAuthKeyData`).(*GAuth.KeyData)
		if !ok {
			return errors.New(ctx.T(`从session获取GAuthKeyData失败`))
		}
	} else {
		if profile, ok := config.DefaultConfig.Sys.Accounts[user]; ok && profile.GAuthKey != nil {
			keyData = profile.GAuthKey
		} else {
			return errors.New(ctx.T(`从用户资料中获取GAuthKey失败`))
		}
	}
	ok, err := GAuth.VerifyFrom(keyData, ctx.Form("code"))
	if err != nil {
		return err
	}
	if !ok {
		return errors.New(ctx.T(`验证码不正确`))
	}
	if testAndBind {
		if profile, ok := config.DefaultConfig.Sys.Accounts[user]; ok {
			profile.GAuthKey = keyData
			config.DefaultConfig.Sys.Accounts[user] = profile
			err := config.DefaultConfig.SaveToFile()
			if err != nil {
				return err
			}
		}
	}
	return nil
}
