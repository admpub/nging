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
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/admpub/log"
	"github.com/admpub/nging/application/handler"
	"github.com/admpub/nging/application/library/config"
	"github.com/admpub/nging/application/library/cron"
	"github.com/admpub/nging/application/model"
	"github.com/webx-top/com"
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
		adminUser := ctx.Form(`adminUser`)
		adminPass := ctx.Form(`adminPass`)
		adminEmail := ctx.Form(`adminEmail`)
		if len(adminUser) == 0 {
			err = errors.New(ctx.T(`管理员用户名不能为空`))
			goto DIE
		}
		if !com.IsUsername(adminUser) {
			err = errors.New(ctx.T(`管理员名不能包含特殊字符(只能由字母、数字、下划线和汉字组成)`))
			goto DIE
		}
		if len(adminPass) < 8 {
			err = errors.New(ctx.T(`管理员密码不能少于8个字符`))
			goto DIE
		}
		if len(adminEmail) == 0 {
			err = errors.New(ctx.T(`管理员邮箱不能为空`))
			goto DIE
		}
		if !ctx.Validate(`adminEmail`, adminEmail, `email`).Ok() {
			err = errors.New(ctx.T(`管理员邮箱格式不正确`))
			goto DIE
		}
		err = ctx.MustBind(&config.DefaultConfig.DB)
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
		if err != nil {
			goto DIE
		}
		//连接数据库
		err = config.ConnectDB()
		if err != nil {
			err = createDatabase(err)
		}
		if err != nil {
			goto DIE
		}
		//创建数据库数据
		var sqlStr string
		installer, ok := config.DBInstallers[config.DefaultConfig.DB.Type]
		if !ok {
			err = errors.New(ctx.T(`不支持安装到%s`, config.DefaultConfig.DB.Type))
			goto DIE
		}
		err = com.SeekFileLines(sqlFile, func(line string) error {
			if strings.HasPrefix(line, `--`) {
				return nil
			}
			line = strings.TrimSpace(line)
			sqlStr += line
			if strings.HasSuffix(line, `;`) && len(sqlStr) > 0 {
				defer func() {
					sqlStr = ``
				}()
				return installer(sqlStr)
			}
			return nil
		})
		if err != nil {
			goto DIE
		}

		if err2 := cron.InitJobs(); err2 != nil {
			log.Error(err2)
		}

		m := model.NewUser(ctx)
		err = m.Register(adminUser, adminPass, adminEmail)
		if err != nil {
			goto DIE
		}

		//保存数据库账号到配置文件
		err = config.DefaultConfig.SaveToFile()
		if err != nil {
			goto DIE
		}

		//生成锁文件
		err = ioutil.WriteFile(lockFile, []byte(time.Now().Format(`2006-01-02 15:04:05`)), os.ModePerm)
		if err != nil {
			goto DIE
		}
		handler.SendOk(ctx, ctx.T(`安装成功`))
		return ctx.Redirect(`/`)
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
