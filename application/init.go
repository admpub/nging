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
package application

import (
	"path/filepath"
	"strings"

	"github.com/admpub/log"
	"github.com/webx-top/com"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/handler/captcha"

	. "github.com/admpub/nging/application/handler"
	_ "github.com/admpub/nging/application/handler/caddy"
	_ "github.com/admpub/nging/application/handler/database"
	_ "github.com/admpub/nging/application/handler/ftp"
	_ "github.com/admpub/nging/application/handler/index"
	_ "github.com/admpub/nging/application/handler/server"
	_ "github.com/admpub/nging/application/handler/setup"
	_ "github.com/admpub/nging/application/handler/user"
	"github.com/admpub/nging/application/library/config"
	"github.com/admpub/nging/application/middleware"
)

func Initialize(e *echo.Echo) {
	e.Use(middleware.FuncMap())
	addRouter(e)
	me := com.MonitorEvent{
		Modify: func(file string) {
			if !strings.HasSuffix(file, `.yaml`) {
				return
			}
			log.Info(`reload config from ` + file)
			err := config.ParseConfig()
			if err == nil {
				return
			}
			if config.IsInstalled() {
				config.MustOK(err)
			} else {
				log.Error(err)
			}
		},
	}
	me.Watch(filepath.Dir(config.DefaultCLIConfig.Conf))
}

func addRouter(e *echo.Echo) {
	opt := captcha.Options{EnableImage: true}
	opt.Wrapper(e)
	for _, register := range Handlers {
		register(e)
	}
	for group, handlers := range GroupHandlers {
		g := e.Group(group, middleware.AuthCheck)
		for _, register := range handlers {
			register(g)
		}
	}
}
