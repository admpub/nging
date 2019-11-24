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
	"context"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/webx-top/com"
	"github.com/webx-top/echo"

	"github.com/admpub/log"
	"github.com/admpub/nging/application/handler"
	"github.com/admpub/nging/application/library/config"
	"github.com/admpub/nging/application/model"
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

	OnInstalled func(ctx context.Context) error
)

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

func install(ctx echo.Context, sqlFile string, installer func(string) error) (err error) {
	installProgress = &ProgressInfo{
		Timestamp: time.Now().Local().Unix(),
	}
	defer func() {
		if err != nil {
			installProgress = nil
		}
	}()
	var sqlStr string
	installProgress.TotalSize, err = com.FileSize(sqlFile)
	if err != nil {
		return
	}
	installProgress.TotalSize += int64(len(handler.OfficialSQL))
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
	err = com.SeekFileLines(sqlFile, installFunction)
	if err != nil {
		return
	}
	for _, line := range strings.Split(handler.OfficialSQL, "\n") {
		err = installFunction(line)
		if err != nil {
			return
		}
	}
	return
}

func Setup(ctx echo.Context) error {
	var err error
	lockFile := filepath.Join(echo.Wd(), `installed.lock`)
	if info, err := os.Stat(lockFile); err == nil && info.IsDir() == false {
		msg := ctx.T(`已经安装过了。如要重新安装，请先删除%s`, lockFile)
		if ctx.IsAjax() {
			return ctx.JSON(ctx.Data().SetInfo(msg, 0))
		}
		return ctx.String(msg)
	}
	sqlFiles, err := config.GetSQLInstallFiles()
	if err != nil {
		msg := ctx.T(`找不到文件%s，无法安装`, `config/install.sql`)
		if ctx.IsAjax() {
			return ctx.JSON(ctx.Data().SetInfo(msg, 0))
		}
		return ctx.String(msg)
	}
	if ctx.IsPost() && installProgress == nil {

		err = ctx.MustBind(&config.DefaultConfig.DB)
		if err != nil {
			return err
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
		}
		if err != nil {
			return err
		}
		//创建数据库数据
		installer, ok := config.DBInstallers[config.DefaultConfig.DB.Type]
		if !ok {
			err = ctx.E(`不支持安装到%s`, config.DefaultConfig.DB.Type)
			return err
		}

		adminUser := ctx.Form(`adminUser`)
		adminPass := ctx.Form(`adminPass`)
		adminEmail := ctx.Form(`adminEmail`)
		if len(adminUser) == 0 {
			err = ctx.E(`管理员用户名不能为空`)
			return err
		}
		if !com.IsUsername(adminUser) {
			err = ctx.E(`管理员名不能包含特殊字符(只能由字母、数字、下划线和汉字组成)`)
			return err
		}
		if len(adminPass) < 8 {
			err = ctx.E(`管理员密码不能少于8个字符`)
			return err
		}
		if len(adminEmail) == 0 {
			err = ctx.E(`管理员邮箱不能为空`)
			return err
		}
		if !ctx.Validate(`adminEmail`, adminEmail, `email`).Ok() {
			err = ctx.E(`管理员邮箱格式不正确`)
			return err
		}
		data := ctx.Data()
		for _, sqlFile := range sqlFiles {
			err = install(ctx, sqlFile, installer)
			if err != nil {
				break
			}
		}
		err = config.ConnectDB(config.DefaultConfig)
		if err != nil {
			return err
		}
		m := model.NewUser(ctx)
		err = m.Register(adminUser, adminPass, adminEmail)
		if err != nil {
			return err
		}
		defer func(cfg *config.Config) {
			if err == nil {
				time.Sleep(1 * time.Second)
				config.DefaultCLIConfig.RunStartup()
				if err := Upgrade(); err != nil {
					log.Error(err)
				}

				// 保存配置
				cfg.InitSecretKey()

				//保存数据库账号到配置文件
				err = cfg.SaveToFile()
				if err != nil {
					return
				}

				//生成锁文件
				err = config.SetInstalled(lockFile)
				if err == nil && OnInstalled != nil {
					err = OnInstalled(ctx)
				}
			}
		}(config.DefaultConfig)
		if ctx.IsAjax() {
			if err != nil {
				data.SetError(err)
			} else {
				data.SetInfo(ctx.T(`安装成功`)).SetData(installProgress)
			}
			return ctx.JSON(data)
		}
		if err != nil {
			goto DIE
		}
		handler.SendOk(ctx, ctx.T(`安装成功`))
		return ctx.Redirect(handler.URLFor(`/`))
	}

DIE:
	ctx.Set(`dbEngines`, config.DBEngines.Slice())
	return ctx.Render(`setup`, handler.Err(ctx, err))
}

func createDatabase(err error) error {
	if fn, ok := config.DBCreaters[config.DefaultConfig.DB.Type]; ok {
		return fn(err, config.DefaultConfig)
	}
	return err
}
