package manager

import (
	"time"

	"github.com/coscms/webcore/cmd"
	"github.com/coscms/webcore/library/config"
	"github.com/coscms/webcore/library/license"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/code"
	"github.com/webx-top/echo/param"
)

func selfExit() {
	cmd.SelfRestart()
}

func selfUpgrade(ctx echo.Context) error {
	if config.FromFile().Settings().Debug && ctx.Formx(`restart`).Bool() { // url: /manager/upgrade?restart=true
		selfExit()
		return nil
	}
	data := ctx.Data()
	if ctx.Formx(`local`).Bool() { // url: /manager/upgrade?local=true
		return ctx.JSON(data.SetData(echo.H{
			`local`: config.Version,
		}))
	}
	download := ctx.Formx(`download`).Bool()
	exit := ctx.Formx(`exit`).Bool()
	if !download && exit { // url: /manager/upgrade?exit=true
		nonce := ctx.Formx(`nonce`).String()
		expected, ok := ctx.Session().Get(`nging.exit.nonce`).(string)
		if !ok {
			return ctx.JSON(data.SetError(ctx.NewError(code.InvalidParameter, `无效参数: %s`, `nonce`).SetZone(`nonce`)))
		}
		if nonce != expected {
			return ctx.JSON(data.SetError(ctx.NewError(code.InvalidParameter, `无效参数: %s`, `nonce`).SetZone(`nonce`)))
		}
		ctx.Session().Delete(`nging.exit.nonce`).Save()
		selfExit()
		return ctx.JSON(data.SetInfo(ctx.T(`升级成功`), code.Success.Int()))
	}
	version := ctx.Formx(`version`).String()
	info, err := license.LatestVersion(ctx, version, download)
	if err != nil {
		return ctx.JSON(data.SetError(err))
	}
	if download { // url: /manager/upgrade?download=true
		nonce := ctx.Formx(`nonce`).String()
		expected, ok := ctx.Session().Get(`nging.upgrade.nonce`).(string)
		if !ok {
			return ctx.JSON(data.SetError(ctx.NewError(code.InvalidParameter, `无效参数: %s`, `nonce`).SetZone(`nonce`)))
		}
		if nonce != expected {
			return ctx.JSON(data.SetError(ctx.NewError(code.InvalidParameter, `无效参数: %s`, `nonce`).SetZone(`nonce`)))
		}
		ctx.Session().Delete(`nging.upgrade.nonce`)
		err = info.Upgrade(ctx, echo.Wd(), `none`)
		if err != nil {
			return ctx.JSON(data.SetError(err))
		}
		if !exit {
			nonce = param.AsString(time.Now().UnixMilli())
			ctx.Session().Set(`nging.exit.nonce`, nonce)
			data.SetData(echo.H{`nonce`: nonce})
		} else {
			ctx.Session().Save()
			selfExit()
		}
		return ctx.JSON(data.SetInfo(ctx.T(`升级成功`), code.Success.Int()))
	}
	nonce := time.Now().UnixMilli()
	ctx.Session().Set(`nging.upgrade.nonce`, param.AsString(nonce))
	return ctx.JSON(data.SetData(echo.H{
		`local`:  config.Version,
		`remote`: info,
		`isNew`:  info.IsNew(),
		`nonce`:  nonce,
	}))
}
