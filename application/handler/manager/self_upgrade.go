package manager

import (
	"context"
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/coscms/webcore/cmd"
	"github.com/coscms/webcore/library/backend"
	"github.com/coscms/webcore/library/config"
	"github.com/coscms/webcore/library/license"
	"github.com/coscms/webcore/library/notice"
	"github.com/webx-top/com"
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
	if ctx.Formx(`upload`).Bool() {
		if ctx.IsPost() {
			return selfUpgradeUpload(ctx)
		}
		return ctx.Render(`manager/upgrade`, nil)
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

func selfUpgradeUpload(ctx echo.Context) error {
	var username string
	user := backend.User(ctx)
	if user != nil {
		username = user.Username
	}
	pv := &license.ProductVersion{}
	np := notice.NewP(ctx, `ngingUploadNewVersion`, username, context.Background()).AutoComplete(true)
	pv.SetProgressor(np)
	saveTo := filepath.Join(echo.Wd(), `data/cache/nging-new-version`, `upload`)
	err := com.MkdirAll(saveTo, os.ModePerm)
	if err != nil {
		return err
	}
	var newFileName string
	_, err = ctx.SaveUploadedFile(`file`, saveTo, func(fh *multipart.FileHeader) (string, error) {
		newFileName = filepath.Base(fh.Filename)
		extension := strings.ToLower(filepath.Ext(fh.Filename))
		if extension != `.gz` && extension != `.tar` && extension != `.tar.gz` {
			return newFileName, fmt.Errorf(`sorry, uploading this file is not supported. only files with the extension “%s” can be uploaded`, `.tar.gz`)
		}
		return newFileName, nil
	})
	if err != nil {
		return err
	}
	saveTo += echo.FilePathSeparator + newFileName
	pv.SetDownloadedPath(saveTo)
	err = pv.Extract()
	if err != nil {
		return err
	}
	err = pv.Upgrade(ctx, echo.Wd(), `none`)
	if err != nil {
		return err
	}
	go func() {
		time.Sleep(time.Millisecond * 500)
		selfExit()
	}()
	data := ctx.Data()
	data.SetData(echo.H{`newVersion`: pv.Version})
	data.SetInfo(ctx.T(`上传成功。正在升级，请稍候...`), code.Success.Int())
	return ctx.JSON(data)
}
