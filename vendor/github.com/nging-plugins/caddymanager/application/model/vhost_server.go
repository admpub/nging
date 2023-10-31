package model

import (
	"strings"

	"github.com/admpub/nging/v5/application/library/common"
	"github.com/nging-plugins/caddymanager/application/dbschema"
	"github.com/nging-plugins/caddymanager/application/library/engine"
	"github.com/webx-top/com"
	"github.com/webx-top/db"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/code"
)

func NewVhostServer(ctx echo.Context) *VhostServer {
	return &VhostServer{
		NgingVhostServer: dbschema.NewNgingVhostServer(ctx),
	}
}

type VhostServer struct {
	*dbschema.NgingVhostServer
}

func (f *VhostServer) check() error {
	ctx := f.Context()
	f.Name = strings.TrimSpace(f.Name)
	if len(f.Name) == 0 {
		return ctx.NewError(code.InvalidParameter, `名称不能为空`).SetZone(`name`)
	}
	f.Ident = strings.TrimSpace(f.Ident)
	if len(f.Ident) == 0 {
		return ctx.NewError(code.InvalidParameter, `唯一标识不能为空`).SetZone(`ident`)
	}
	if strings.ToLower(f.Ident) == `default` {
		return ctx.NewError(code.InvalidParameter, `唯一标识“default”是系统内保留标识，不可使用，请修改`).SetZone(`ident`)
	}
	if !com.IsAlphaNumericUnderscoreHyphen(f.Ident) {
		return ctx.NewError(code.InvalidParameter, `唯一标识只能包含字母、数字、下划线或短横`).SetZone(`ident`)
	}
	f.Engine = strings.TrimSpace(f.Engine)
	if len(f.Engine) == 0 {
		return ctx.NewError(code.InvalidParameter, `引擎必选`).SetZone(`engine`)
	}
	if f.Engine == `default` || !engine.Engines.Has(f.Engine) {
		return ctx.NewError(code.InvalidParameter, `引擎无效`).SetZone(`engine`)
	}
	if !engine.Environs.Has(f.Environ) {
		return ctx.NewError(code.InvalidParameter, `环境选项值无效`).SetZone(`environ`)
	}
	f.ExecutableFile = strings.TrimSpace(f.ExecutableFile)
	f.Endpoint = strings.TrimSpace(f.Endpoint)
	f.Env = strings.TrimSpace(f.Env)
	if len(f.ExecutableFile) == 0 {
		return ctx.NewError(code.InvalidParameter, `执行文件路径不能为空`).SetZone(`executableFile`)
	}
	if f.Environ == engine.EnvironContainer {
		if len(f.Endpoint) > 0 && (strings.HasPrefix(f.Endpoint, `ftp`) || !com.IsURL(f.Endpoint)) {
			return ctx.NewError(code.InvalidParameter, `API接口网址无效`).SetZone(`endpoint`)
		}
	} else {
		f.Endpoint = ``
	}
	var exists bool
	var err error
	if f.Id > 0 {
		exists, err = f.Exists(f.Ident, f.Id)
	} else {
		exists, err = f.Exists(f.Ident)
	}
	if err != nil {
		return err
	}
	if exists {
		return ctx.NewError(code.DataAlreadyExists, `唯一标识“%s”已经存在`, f.Ident).SetZone(`ident`)
	}
	f.ConfigLocalFile = strings.TrimSpace(f.ConfigLocalFile)
	f.ConfigContainerFile = strings.TrimSpace(f.ConfigContainerFile)
	f.VhostConfigLocalDir = strings.TrimSpace(f.VhostConfigLocalDir)
	f.VhostConfigContainerDir = strings.TrimSpace(f.VhostConfigContainerDir)
	f.CertLocalDir = strings.TrimSpace(f.CertLocalDir)
	f.CertContainerDir = strings.TrimSpace(f.CertContainerDir)
	f.CertPathFormatCert = strings.TrimSpace(f.CertPathFormatCert)
	f.CertPathFormatKey = strings.TrimSpace(f.CertPathFormatKey)
	f.CertPathFormatTrust = strings.TrimSpace(f.CertPathFormatTrust)
	f.WorkDir = strings.TrimSpace(f.WorkDir)
	f.CmdWithConfig = common.GetBoolFlag(f.CmdWithConfig, common.BoolN)
	f.Disabled = common.GetBoolFlag(f.Disabled, common.BoolN)
	f.AutoModifyConfig = common.GetBoolFlag(f.AutoModifyConfig, common.BoolN)
	f.CertAutoRenew = common.GetBoolFlag(f.CertAutoRenew, common.BoolN)
	return nil
}

func (f *VhostServer) Exists(ident string, exclude ...uint) (bool, error) {
	cond := db.NewCompounds()
	cond.AddKV(`ident`, ident)
	if len(exclude) > 0 {
		cond.AddKV(`id`, db.NotEq(exclude[0]))
	}
	return f.NgingVhostServer.Exists(nil, cond)
}

func (f *VhostServer) Add() (interface{}, error) {
	if err := f.check(); err != nil {
		return nil, err
	}
	return f.NgingVhostServer.Insert()
}

func (f *VhostServer) Edit(mw func(db.Result) db.Result, args ...interface{}) (err error) {
	if err = f.check(); err != nil {
		return err
	}
	return f.NgingVhostServer.Update(mw, args...)
}

func (f *VhostServer) Delete(mw func(db.Result) db.Result, args ...interface{}) (err error) {
	return f.NgingVhostServer.Delete(mw, args...)
}

func (f *VhostServer) SetConfigFileUpdated(ts uint) error {
	return f.NgingVhostServer.UpdateField(nil, `config_file_updated`, ts, `id`, f.Id)
}
