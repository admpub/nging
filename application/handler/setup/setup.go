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
	"github.com/admpub/copier"
	"github.com/webx-top/com"
	"github.com/webx-top/echo"
	stdCode "github.com/webx-top/echo/code"

	"github.com/admpub/errors"
	"github.com/admpub/log"
	"github.com/admpub/nging/v5/application/handler"
	"github.com/admpub/nging/v5/application/library/common"
	"github.com/admpub/nging/v5/application/library/config"
	"github.com/admpub/nging/v5/application/model"
	"github.com/admpub/nging/v5/application/registry/settings"
	"github.com/admpub/nging/v5/application/request"
)

type ProgressInfo struct {
	Finished  int64
	TotalSize int64
	Summary   string
	Timestamp int64
	mu        *sync.RWMutex
}

func (p *ProgressInfo) Clone() ProgressInfo {
	p.mu.RLock()
	r := *p
	p.mu.RUnlock()
	return r
}

func (p *ProgressInfo) GetTs() int64 {
	p.mu.RLock()
	r := p.Timestamp
	p.mu.RUnlock()
	return r
}

func (p *ProgressInfo) SetTs(ts int64) {
	p.mu.Lock()
	p.Timestamp = ts
	p.mu.Unlock()
}

func (p *ProgressInfo) Done(n int64) {
	p.mu.Lock()
	newVal := p.Finished + n
	if newVal > p.TotalSize {
		p.Finished = p.TotalSize
	} else {
		p.Finished = newVal
	}
	p.mu.Unlock()
}

func (p *ProgressInfo) Add(n int64) {
	p.mu.Lock()
	p.TotalSize += n
	p.mu.Unlock()
}

var (
	installProgress   = &ProgressInfo{mu: &sync.RWMutex{}}
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
		prog := installProgress.Clone()
		if prog.GetTs() <= 0 {
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
		strLen := len(sqlStr)
		sqlStr = common.ReplaceCharset(sqlStr, charset, true)
		err := installer(sqlStr)
		installProgress.Done(int64(strLen))
		return err
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
	lockFile = filepath.Join(config.FromCLI().ConfDir(), config.LockFileName)
	sqlFiles, err := config.GetSQLInstallFiles()
	if err != nil && len(config.GetInstallSQLs()[`nging`]) == 0 {
		err = ctx.NewError(stdCode.DataNotFound, ctx.T(`找不到文件%s，无法安装`, `config/install.sql`))
		return err
	}
	insertSQLFiles := config.GetSQLInsertFiles()
	var requestData *request.Setup
	if ctx.IsPost() && installProgress.GetTs() <= 0 {
		requestData = echo.GetValidated(ctx).(*request.Setup)
		err = copier.Copy(&config.FromFile().DB, requestData)
		if err != nil {
			return ctx.NewError(stdCode.Failure, err.Error())
		}
		installProgress.SetTs(time.Now().Unix())
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
		installProgress.Add(totalSize)
		config.FromFile().DB.SetKV(`charset`, requestData.Charset)
		//连接数据库
		err = config.ConnectDB(config.FromFile().DB, 0, `default`)
		if err != nil {
			err = createDatabase(err)
			if err != nil {
				return ctx.NewError(stdCode.Failure, err.Error())
			}
		}
		//创建数据库数据
		installer, ok := config.DBInstallers[config.FromFile().DB.Type]
		if !ok {
			err = ctx.NewError(stdCode.Unsupported, ctx.T(`不支持安装到%s`, config.FromFile().DB.Type))
			return err
		}
		data := ctx.Data()
		// 先执行 sql struct (建表)
		for _, sqlFile := range sqlFiles {
			log.Info(color.GreenString(`[installer]`), `Execute SQL file: `, sqlFile)
			err = install(ctx, sqlFile, true, requestData.Charset, installer)
			if err != nil {
				return ctx.NewError(stdCode.Failure, err.Error())
			}
		}
		for _, sqlContents := range installSQLs {
			for _, sqlContent := range sqlContents {
				log.Info(color.GreenString(`[installer]`), `Execute SQL: `, sqlContent)
				err = install(ctx, sqlContent, false, requestData.Charset, installer)
				if err != nil {
					return ctx.NewError(stdCode.Failure, err.Error())
				}
			}
		}
		// 再执行 insert (插入数据)
		for _, sqlFile := range insertSQLFiles {
			log.Info(color.GreenString(`[installer]`), `Execute SQL file: `, sqlFile)
			err = install(ctx, sqlFile, true, requestData.Charset, installer)
			if err != nil {
				return ctx.NewError(stdCode.Failure, err.Error())
			}
		}
		for _, sqlContents := range insertSQLs {
			for _, sqlContent := range sqlContents {
				log.Info(color.GreenString(`[installer]`), `Execute SQL: `, sqlContent)
				err = install(ctx, sqlContent, false, requestData.Charset, installer)
				if err != nil {
					return ctx.NewError(stdCode.Failure, err.Error())
				}
			}
		}

		// 重新连接数据库
		log.Info(color.GreenString(`[installer]`), `Reconnect the database`)
		err = config.ConnectDB(config.FromFile().DB, 0, `default`)
		if err != nil {
			return ctx.NewError(stdCode.Failure, err.Error())
		}

		// 添加创始人
		m := model.NewUser(ctx)
		log.Info(color.GreenString(`[installer]`), `Create Administrator`)
		err = m.Register(requestData.AdminUser, requestData.AdminPass, requestData.AdminEmail, ``)
		if err != nil {
			err = errors.WithMessage(err, `Create Administrator`)
			return ctx.NewError(stdCode.Failure, err.Error())
		}

		// 生成安全密钥
		log.Info(color.GreenString(`[installer]`), `Generate a security key`)
		config.FromFile().InitSecretKey()

		// 保存数据库账号到配置文件
		log.Info(color.GreenString(`[installer]`), `Save the configuration file`)
		err = config.FromFile().SaveToFile()
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
		config.FromCLI().RunStartup()

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
	if requestData == nil {
		requestData = &request.Setup{}
	}
	ctx.Set(`data`, requestData)
	ctx.Set(`dbEngines`, config.DBEngines.Slice())
	return ctx.Render(`setup`, handler.Err(ctx, err))
}

func createDatabase(err error) error {
	if fn, ok := config.DBCreaters[config.FromFile().DB.Type]; ok {
		return fn(err, config.FromFile().DB)
	}
	return err
}
