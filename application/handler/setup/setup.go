/*

   Copyright 2016 Wenhui Shen <www.webx.top>

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.

*/
package setup

import (
	"os"
	"path/filepath"
	"time"

	"strings"

	"io/ioutil"

	"github.com/admpub/nging/application/handler"
	"github.com/admpub/nging/application/library/config"
	"github.com/webx-top/com"
	"github.com/webx-top/db/lib/factory"
	"github.com/webx-top/echo"
)

func init() {
	handler.Register(func(e *echo.Echo) {
		e.Route("GET,POST", `/setup`, Setup)
	})
}

func Setup(ctx echo.Context) error {
	var err error
	lockFile := filepath.Join(com.SelfDir(), `installed.lock`)
	if info, err := os.Stat(lockFile); err == nil && info.IsDir() == false {
		return ctx.String(ctx.T(`已经安装过了。如要重新安装，请先删除%s`, lockFile))
	}
	sqlFile := filepath.Join(filepath.Dir(config.DefaultCLIConfig.Conf), `install.sql`)
	if !com.FileExists(sqlFile) {
		return ctx.String(ctx.T(`找不到文件%s，无法安装`, sqlFile))
	}
	if ctx.IsPost() {
		err = ctx.MustBind(&config.DefaultConfig.DB)
		config.DefaultConfig.DB.Database = strings.Replace(config.DefaultConfig.DB.Database, "'", "", -1)
		config.DefaultConfig.DB.Database = strings.Replace(config.DefaultConfig.DB.Database, "`", "", -1)
		if err == nil {
			//连接数据库
			err = config.ConnectDB()
			if err != nil {
				err = createDatabase(err)
			}
		}
		if err == nil {
			//创建数据库数据
			var sqlStr string
			err = com.SeekFileLines(sqlFile, func(line string) error {
				if strings.HasPrefix(line, `--`) {
					return nil
				}
				line = strings.TrimSpace(line)
				sqlStr += line
				if strings.HasSuffix(line, `;`) && len(sqlStr) > 0 {
					_, err := factory.NewParam().SetCollection(sqlStr).Exec()
					sqlStr = ``
					if err != nil {
						return err
					}
				}
				return nil
			})
		}
		if err == nil {
			//保存数据库账号到配置文件
			err = config.DefaultConfig.SaveToFile()
		}
		if err == nil {
			//生成锁文件
			err = ioutil.WriteFile(lockFile, []byte(time.Now().Format(`2006-01-02 15:04:05`)), os.ModePerm)
		}
		if err == nil {
			handler.SendOk(ctx, ctx.T(`安装成功`))
			return ctx.Redirect(`/`)
		}
	}
	return ctx.Render(`setup`, handler.Err(ctx, err))
}

func createDatabase(err error) error {
	switch config.DefaultConfig.DB.Type {
	case `mysql`:
		if strings.Contains(err.Error(), `Unknown database`) {
			dbName := config.DefaultConfig.DB.Database
			config.DefaultConfig.DB.Database = ``
			err2 := config.ConnectDB()
			if err2 == nil {
				sqlStr := "CREATE DATABASE `" + dbName + "`"
				_, err2 = factory.NewParam().SetCollection(sqlStr).Exec()
				if err2 == nil {
					config.DefaultConfig.DB.Database = dbName
					err = config.ConnectDB()
				}
			}
		}
	}
	return err
}
