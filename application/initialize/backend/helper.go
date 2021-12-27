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

package backend

import (
	"path/filepath"

	"github.com/webx-top/echo/handler/captcha"
	"github.com/webx-top/echo/middleware/render"

	"github.com/admpub/nging/v4/application/handler"
	"github.com/admpub/nging/v4/application/library/common"
	"github.com/admpub/nging/v4/application/library/config"
	"github.com/admpub/nging/v4/application/middleware"
)

func Initialize() {
	handler.Use(BackendURLFuncMW(), middleware.FuncMap(), middleware.BackendFuncMap(), render.Auto())
	handler.Use(middleware.Middlewares...)
	addRouter()
	DefaultConfigWatcher(true)
}

var onConfigChange = []func(file string) error{}

func OnConfigChange(fn func(file string) error) {
	onConfigChange = append(onConfigChange, fn)
}

func FireConfigChange(file string) error {
	for _, fn := range onConfigChange {
		if err := fn(file); err != nil {
			return err
		}
	}
	return common.ErrIgnoreConfigChange
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
			return FireConfigChange(filePath)
		}
	})
}

func addRouter() {
	opt := captcha.Options{EnableImage: true}
	opt.Wrapper(handler.IRegister().Echo())
	handler.UseToGroup(`*`, middleware.AuthCheck) //应用中间件到所有子组
	handler.Apply()
}
