package manager

import (
	"os"
	"time"

	"github.com/admpub/nging/v5/application/cmd"
	"github.com/admpub/nging/v5/application/library/config"
	"github.com/admpub/nging/v5/application/library/license"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/code"
	"github.com/webx-top/echo/param"
)

func selfExit() {
	cmd.SendSignal(os.Interrupt, cmd.ExitCodeSelfRestart)
	time.Sleep(time.Second)
	os.Exit(cmd.ExitCodeSelfRestart)
}

func selfUpgrade(ctx echo.Context) error {
	data := ctx.Data()
	if ctx.Formx(`local`).Bool() {
		return ctx.JSON(data.SetData(echo.H{
			`local`: config.Version,
		}))
	}
	download := ctx.Formx(`download`).Bool()
	exit := ctx.Formx(`exit`).Bool()
	if !download && exit {
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
	if download {
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
