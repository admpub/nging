package handler

import (
	"errors"

	"encoding/gob"

	QR "github.com/RaymondChou/goqr/pkg"
	GAuth "github.com/admpub/dgoogauth"
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
	user, ok := ctx.Get(`user`).(string)
	if !ok {
		return errors.New(ctx.T(`登录信息获取失败，请重新登录`))
	}
	if ctx.IsPost() {
		validOk := GAuthVerify(ctx)
		ctx.Set(`validOk`, validOk)
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
	return ctx.Render(`gauthbind`, nil)
}

func GAuthVerify(ctx echo.Context) bool {
	key := ctx.Form("key")
	return GAuth.Verify(key, ctx.Form("code"))
}
