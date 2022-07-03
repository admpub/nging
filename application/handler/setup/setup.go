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

package setup

import (
	"fmt"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/admpub/color"
	"github.com/webx-top/com"
	"github.com/webx-top/echo"
	stdCode "github.com/webx-top/echo/code"

	"github.com/admpub/errors"
	"github.com/admpub/log"
	"github.com/admpub/nging/v4/application/handler"
	"github.com/admpub/nging/v4/application/library/common"
	"github.com/admpub/nging/v4/application/library/config"
	"github.com/admpub/nging/v4/application/library/config/subconfig/sdb"
	"github.com/admpub/nging/v4/application/model"
	"github.com/admpub/nging/v4/application/registry/settings"
)

type ProgressInfo struct {
	Finished  int64
	TotalSize int64
	Summary   string
	Timestamp int64
}

var (
	lockProgress      sync.RWMutex
	installProgress   *ProgressInfo
	installedProgress = &ProgressInfo{
		Finished:  1,
		TotalSize: 1,
	}
	uninstallProgress = &ProgressInfo{
		Finished:  0,
		TotalSize: 1,
	}

	onInstalled        []func(ctx echo.Context) error
	RegisterInstallSQL = config.RegisterInstallSQL
)

func getInstallProgress() *ProgressInfo {
	lockProgress.RLock()
	v := installProgress
	lockProgress.RUnlock()
	return v
}

func setInstallProgress(v *ProgressInfo) {
	lockProgress.Lock()
	installProgress = v
	lockProgress.Unlock()
}

func OnInstalled(cb func(ctx echo.Context) error) {
	if cb == nil {
		return
	}
	onInstalled = append(onInstalled, cb)
}

func Progress(ctx echo.Context) error {
	data := ctx.Data()
	if config.IsInstalled() {
		data.SetInfo(ctx.T(`已经安装过了`), 0)
		data.SetData(installedProgress)
	} else {
		installProgress := getInstallProgress()
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

func install(ctx echo.Context, sqlFile string, isFile bool, charset string, installer func(string) error) (err error) {
	installFunction := common.SQLLineParser(func(sqlStr string) error {
		sqlStr = common.ReplaceCharset(sqlStr, charset, true)
		return installer(sqlStr)
	})
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
	lockFile := config.InstalledLockFile()
	if len(lockFile) > 0 {
		err := ctx.NewError(stdCode.RepeatOperation, ctx.T(`已经安装过了。如要重新安装，请先删除%s`, filepath.Base(lockFile)))
		return err
	}
	lockFile = filepath.Join(config.DefaultCLIConfig.ConfDir(), config.LockFileName)
	sqlFiles, err := config.GetSQLInstallFiles()
	if err != nil && len(config.GetInstallSQLs()[`nging`]) == 0 {
		err = ctx.NewError(stdCode.DataNotFound, ctx.T(`找不到文件%s，无法安装`, `config/install.sql`))
		return err
	}
	insertSQLFiles := config.GetSQLInsertFiles()
	if ctx.IsPost() && getInstallProgress() == nil {
		installProgress := &ProgressInfo{
			Timestamp: time.Now().Unix(),
		}
		setInstallProgress(installProgress)
		defer func() {
			installProgress = nil
			setInstallProgress(installProgress)
		}()
		installSQLs := config.GetInstallSQLs()
		insertSQLs := config.GetInsertSQLs()
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
		for _, sqlContents := range installSQLs {
			for _, sqlContent := range sqlContents {
				totalSize += int64(len(sqlContent))
			}
		}
		for _, sqlFile := range insertSQLFiles {
			var fileSize int64
			fileSize, err = com.FileSize(sqlFile)
			if err != nil {
				err = errors.WithMessage(err, sqlFile)
				return ctx.NewError(stdCode.Failure, err.Error())
			}
			totalSize += fileSize
		}
		for _, sqlContents := range insertSQLs {
			for _, sqlContent := range sqlContents {
				totalSize += int64(len(sqlContent))
			}
		}
		installProgress.TotalSize = totalSize
		err = ctx.MustBind(&config.DefaultConfig.DB)
		if err != nil {
			return ctx.NewError(stdCode.Failure, err.Error())
		}
		charset := sdb.MySQLDefaultCharset
		config.DefaultConfig.DB.Database = strings.Replace(config.DefaultConfig.DB.Database, "'", "", -1)
		config.DefaultConfig.DB.Database = strings.Replace(config.DefaultConfig.DB.Database, "`", "", -1)
		if config.DefaultConfig.DB.Type == `sqlite` {
			config.DefaultConfig.DB.User = ``
			config.DefaultConfig.DB.Password = ``
			config.DefaultConfig.DB.Host = ``
			if !strings.HasSuffix(config.DefaultConfig.DB.Database, `.db`) {
				config.DefaultConfig.DB.Database += `.db`
			}
		} else {
			charset = ctx.Form(`charset`, sdb.MySQLDefaultCharset)
			if !com.InSlice(charset, sdb.MySQLSupportCharsetList) {
				return ctx.NewError(stdCode.InvalidParameter, ctx.T(`字符集参数无效`)).SetZone(`charset`)
			}
		}
		config.DefaultConfig.DB.SetKV(`charset`, charset)
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
			err = ctx.NewError(stdCode.Unsupported, ctx.T(`不支持安装到%s`, config.DefaultConfig.DB.Type))
			return err
		}

		adminUser := ctx.Form(`adminUser`)
		adminPass := ctx.Form(`adminPass`)
		adminEmail := ctx.Form(`adminEmail`)
		if len(adminUser) == 0 {
			err = ctx.NewError(stdCode.InvalidParameter, ctx.T(`管理员用户名不能为空`))
			return err
		}
		if !com.IsUsername(adminUser) {
			err = ctx.NewError(stdCode.InvalidParameter, ctx.T(`管理员名不能包含特殊字符(只能由字母、数字、下划线和汉字组成)`))
			return err
		}
		if len(adminPass) < 8 {
			err = ctx.NewError(stdCode.InvalidParameter, ctx.T(`管理员密码不能少于8个字符`))
			return err
		}
		if len(adminEmail) == 0 {
			err = ctx.NewError(stdCode.InvalidParameter, ctx.T(`管理员邮箱不能为空`))
			return err
		}
		if !ctx.Validate(`adminEmail`, adminEmail, `email`).Ok() {
			err = ctx.NewError(stdCode.InvalidParameter, ctx.T(`管理员邮箱格式不正确`))
			return err
		}
		data := ctx.Data()
		// 先执行 sql struct (建表)
		for _, sqlFile := range sqlFiles {
			log.Info(color.GreenString(`[installer]`), `Execute SQL file: `, sqlFile)
			err = install(ctx, sqlFile, true, charset, installer)
			if err != nil {
				return ctx.NewError(stdCode.Failure, err.Error())
			}
		}
		for _, sqlContents := range installSQLs {
			for _, sqlContent := range sqlContents {
				log.Info(color.GreenString(`[installer]`), `Execute SQL: `, sqlContent)
				err = install(ctx, sqlContent, false, charset, installer)
				if err != nil {
					return ctx.NewError(stdCode.Failure, err.Error())
				}
			}
		}
		// 再执行 insert (插入数据)
		for _, sqlFile := range insertSQLFiles {
			log.Info(color.GreenString(`[installer]`), `Execute SQL file: `, sqlFile)
			err = install(ctx, sqlFile, true, charset, installer)
			if err != nil {
				return ctx.NewError(stdCode.Failure, err.Error())
			}
		}
		for _, sqlContents := range insertSQLs {
			for _, sqlContent := range sqlContents {
				log.Info(color.GreenString(`[installer]`), `Execute SQL: `, sqlContent)
				err = install(ctx, sqlContent, false, charset, installer)
				if err != nil {
					return ctx.NewError(stdCode.Failure, err.Error())
				}
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

		initConfigDB := func(ctx echo.Context) (ierr error) {
			defer func() {
				if panicErr := recover(); panicErr != nil {
					ierr = fmt.Errorf(`%v`, panicErr)
				}
			}()
			ierr = settings.Init(ctx)
			return
		}
		for i := 1; i <= 5; i++ {
			time.Sleep(time.Duration(i) * time.Second) // 等1秒
			err = initConfigDB(ctx)
			if err == nil {
				break
			}
		}
		if err != nil {
			return err
		}

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
