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

package setup

import (
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/webx-top/com"
	"github.com/webx-top/echo"
	stdCode "github.com/webx-top/echo/code"

	"github.com/admpub/errors"
	"github.com/admpub/log"
	"github.com/admpub/nging/application/handler"
	"github.com/admpub/nging/application/library/config"
	"github.com/admpub/nging/application/model"
	"github.com/admpub/nging/application/registry/settings"
)

type ProgressInfo struct {
	Finished  int64
	TotalSize int64
	Summary   string
	Timestamp int64
}

var (
	installProgress   *ProgressInfo
	installedProgress = &ProgressInfo{
		Finished:  1,
		TotalSize: 1,
	}
	uninstallProgress = &ProgressInfo{
		Finished:  0,
		TotalSize: 1,
	}

	onInstalled []func(ctx echo.Context) error
)

func OnInstalled(cb func(ctx echo.Context) error) {
	if cb == nil {
		return
	}
	onInstalled = append(onInstalled, cb)
}

func init() {
	handler.Register(func(e echo.RouteRegister) {
		e.Route("GET,POST", `/setup`, Setup)
		e.Route("GET", `/progress`, Progress)
		e.Route("GET,POST", `/license`, License)
	})
}

func Progress(ctx echo.Context) error {
	data := ctx.Data()
	if config.IsInstalled() {
		data.SetInfo(ctx.T(`已经安装过了`), 0)
		data.SetData(installedProgress)
	} else {
		if installProgress == nil {
			data.SetInfo(ctx.T(`尚未开始`), 1)
			data.SetData(uninstallProgress)
		} else {
			data.SetInfo(ctx.T(`安装中`), 1)
			data.SetData(installProgress)
		}
	}
	return ctx.JSON(data)
}

func install(ctx echo.Context, sqlFile string, isFile bool, installer func(string) error) (err error) {
	var sqlStr string
	installFunction := func(line string) (rErr error) {
		installProgress.Finished += int64(len(line)) + 1
		if strings.HasPrefix(line, `--`) {
			return nil
		}
		if strings.HasPrefix(line, `/*`) && strings.HasSuffix(line, `*/;`) {
			return nil
		}
		line = strings.TrimSpace(line)
		sqlStr += line
		if strings.HasSuffix(line, `;`) && len(sqlStr) > 0 {
			//installProgress.Summary = sqlStr
			defer func() {
				sqlStr = ``
			}()
			return installer(sqlStr)
		}
		return nil
	}
	if isFile {
		return com.SeekFileLines(sqlFile, installFunction)
	}
	sqlContent := sqlFile
	for _, line := range strings.Split(sqlContent, "\n") {
		err = installFunction(line)
		if err != nil {
			return err
		}
	}
	return err
}

func Setup(ctx echo.Context) error {
	var err error
	lockFile := filepath.Join(echo.Wd(), `installed.lock`)
	if info, err := os.Stat(lockFile); err == nil && info.IsDir() == false {
		err = ctx.NewError(stdCode.RepeatOperation, `已经安装过了。如要重新安装，请先删除%s`, lockFile)
		return err
	}
	sqlFiles, err := config.GetSQLInstallFiles()
	if err != nil {
		err = ctx.NewError(stdCode.DataNotFound, `找不到文件%s，无法安装`, `config/install.sql`)
		return err
	}
	insertSQLFiles := config.GetSQLInsertFiles()
	if len(insertSQLFiles) > 0 {
		sqlFiles = append(sqlFiles, insertSQLFiles...)
	}

	if ctx.IsPost() && installProgress == nil {
		installProgress = &ProgressInfo{
			Timestamp: time.Now().Local().Unix(),
		}
		defer func() {
			if err != nil {
				installProgress = nil
			}
		}()
		var totalSize int64
		for _, sqlFile := range sqlFiles {
			var fileSize int64
			fileSize, err = com.FileSize(sqlFile)
			if err != nil {
				err = errors.WithMessage(err, sqlFile)
				return ctx.NewError(stdCode.Failure, err.Error())
			}
			totalSize += fileSize
		}
		installProgress.TotalSize = totalSize
		installProgress.TotalSize += int64(len(handler.OfficialSQL))
		err = ctx.MustBind(&config.DefaultConfig.DB)
		if err != nil {
			return ctx.NewError(stdCode.Failure, err.Error())
		}
		config.DefaultConfig.DB.Database = strings.Replace(config.DefaultConfig.DB.Database, "'", "", -1)
		config.DefaultConfig.DB.Database = strings.Replace(config.DefaultConfig.DB.Database, "`", "", -1)
		if config.DefaultConfig.DB.Type == `sqlite` {
			config.DefaultConfig.DB.User = ``
			config.DefaultConfig.DB.Password = ``
			config.DefaultConfig.DB.Host = ``
			if strings.HasSuffix(config.DefaultConfig.DB.Database, `.db`) == false {
				config.DefaultConfig.DB.Database += `.db`
			}
		}
		//连接数据库
		err = config.ConnectDB(config.DefaultConfig)
		if err != nil {
			err = createDatabase(err)
			if err != nil {
				return ctx.NewError(stdCode.Failure, err.Error())
			}
		}
		//创建数据库数据
		installer, ok := config.DBInstallers[config.DefaultConfig.DB.Type]
		if !ok {
			err = ctx.NewError(stdCode.Unsupported, `不支持安装到%s`, config.DefaultConfig.DB.Type)
			return err
		}

		adminUser := ctx.Form(`adminUser`)
		adminPass := ctx.Form(`adminPass`)
		adminEmail := ctx.Form(`adminEmail`)
		if len(adminUser) == 0 {
			err = ctx.NewError(stdCode.InvalidParameter, `管理员用户名不能为空`)
			return err
		}
		if !com.IsUsername(adminUser) {
			err = ctx.NewError(stdCode.InvalidParameter, `管理员名不能包含特殊字符(只能由字母、数字、下划线和汉字组成)`)
			return err
		}
		if len(adminPass) < 8 {
			err = ctx.NewError(stdCode.InvalidParameter, `管理员密码不能少于8个字符`)
			return err
		}
		if len(adminEmail) == 0 {
			err = ctx.NewError(stdCode.InvalidParameter, `管理员邮箱不能为空`)
			return err
		}
		if !ctx.Validate(`adminEmail`, adminEmail, `email`).Ok() {
			err = ctx.NewError(stdCode.InvalidParameter, `管理员邮箱格式不正确`)
			return err
		}
		data := ctx.Data()
		for _, sqlFile := range sqlFiles {
			log.Info(color.GreenString(`[installer]`), `Execute SQL file: `, sqlFile)
			err = install(ctx, sqlFile, true, installer)
			if err != nil {
				return ctx.NewError(stdCode.Failure, err.Error())
			}
		}
		if len(handler.OfficialSQL) > 0 {
			err = install(ctx, handler.OfficialSQL, false, installer)
			if err != nil {
				return ctx.NewError(stdCode.Failure, err.Error())
			}
		}

		// 重新连接数据库
		log.Info(color.GreenString(`[installer]`), `Reconnect the database`)
		err = config.ConnectDB(config.DefaultConfig)
		if err != nil {
			return ctx.NewError(stdCode.Failure, err.Error())
		}

		// 添加创始人
		m := model.NewUser(ctx)
		log.Info(color.GreenString(`[installer]`), `Create Administrator`)
		err = m.Register(adminUser, adminPass, adminEmail, ``)
		if err != nil {
			err = errors.WithMessage(err, `Create Administrator`)
			return ctx.NewError(stdCode.Failure, err.Error())
		}

		// 生成安全密钥
		log.Info(color.GreenString(`[installer]`), `Generate a security key`)
		config.DefaultConfig.InitSecretKey()

		// 保存数据库账号到配置文件
		log.Info(color.GreenString(`[installer]`), `Save the configuration file`)
		err = config.DefaultConfig.SaveToFile()
		if err != nil {
			return ctx.NewError(stdCode.Failure, err.Error())
		}

		for _, cb := range onInstalled {
			log.Info(color.GreenString(`[installer]`), `Execute Hook: `, com.FuncName(cb))
			if err = cb(ctx); err != nil {
				return ctx.NewError(stdCode.Failure, err.Error())
			}
		}

		// 生成锁文件
		log.Info(color.GreenString(`[installer]`), `Generated file: `, lockFile)
		err = config.SetInstalled(lockFile)
		if err != nil {
			return ctx.NewError(stdCode.Failure, err.Error())
		}

		time.Sleep(1 * time.Second) // 等1秒
		settings.Init()

		// 启动
		log.Info(color.GreenString(`[installer]`), `Start up`)
		config.DefaultCLIConfig.RunStartup()

		// 升级
		if err := Upgrade(); err != nil {
			log.Error(err)
		}

		if ctx.IsAjax() {
			data.SetInfo(ctx.T(`安装成功`)).SetData(installProgress)
			return ctx.JSON(data)
		}
		handler.SendOk(ctx, ctx.T(`安装成功`))
		return ctx.Redirect(handler.URLFor(`/`))
	}

	ctx.Set(`dbEngines`, config.DBEngines.Slice())
	return ctx.Render(`setup`, handler.Err(ctx, err))
}

func createDatabase(err error) error {
	if fn, ok := config.DBCreaters[config.DefaultConfig.DB.Type]; ok {
		return fn(err, config.DefaultConfig)
	}
	return err
}
