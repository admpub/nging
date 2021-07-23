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

package backend

import (
	"path/filepath"
	"strings"

	"github.com/webx-top/echo/handler/captcha"
	"github.com/webx-top/echo/handler/mvc/static/resource"
	"github.com/webx-top/echo/middleware/render"

	"github.com/admpub/nging/v3/application/handler"
	"github.com/admpub/nging/v3/application/library/common"
	"github.com/admpub/nging/v3/application/library/config"
	"github.com/admpub/nging/v3/application/middleware"
)

func Initialize() {
	handler.Echo().Use(BackendURLFuncMW(), middleware.FuncMap(), middleware.BackendFuncMap(), render.Auto())
	handler.Echo().Use(middleware.Middlewares...)
	addRouter()
	DefaultConfigWatcher(true)
	//config.RunDaemon()
}

func DefaultConfigWatcher(mustOk bool) {
	if config.DefaultCLIConfig.Type != `manager` {
		return
	}
	conf := filepath.Base(config.DefaultCLIConfig.Conf)
	config.WatchConfig(func(file string) error {
		name := filepath.Base(file)
		switch name {
		case conf:
			err := config.ParseConfig()
			if err != nil {
				if mustOk && config.IsInstalled() {
					config.MustOK(err)
				}
			}
			return err
		default:
			if !config.IsInstalled() {
				return nil
			}
			filePath := filepath.ToSlash(file)
			if strings.Contains(filePath, `/frp/server/`) {
				id := config.DefaultCLIConfig.GenerateIDFromConfigFileName(file, true)
				if len(id) == 0 {
					return common.ErrIgnoreConfigChange
				}
				if !config.DefaultCLIConfig.IsRunning(`frpserver.` + id) {
					return common.ErrIgnoreConfigChange
				}
				return config.DefaultCLIConfig.FRPRestartID(id)
			}
			if strings.Contains(filePath, `/frp/client/`) {
				id := config.DefaultCLIConfig.GenerateIDFromConfigFileName(file, true)
				if len(id) == 0 {
					return common.ErrIgnoreConfigChange
				}
				if !config.DefaultCLIConfig.IsRunning(`frpclient.` + id) {
					return common.ErrIgnoreConfigChange
				}
				return config.DefaultCLIConfig.FRPClientRestartID(id)
			}
			return common.ErrIgnoreConfigChange
		}
	})
}

func addRouter() {
	opt := captcha.Options{EnableImage: true}
	opt.Wrapper(handler.Echo())
	handler.Use(`*`, middleware.AuthCheck) //应用中间件到所有子组
	handler.Apply()
	/*
		res := resource.NewStatic(`/public/assets`, filepath.Join(echo.Wd(), `public/assets`))
		resPath := func(rpath string) string {
			return filepath.Join(echo.Wd(), `public/assets`, rpath)
		}
		e.Get(`/minify/*`, func(ctx echo.Context) error {
			return res.HandleMinify(ctx, resPath)
		})
	*/
	_ = resource.Static{}
}
