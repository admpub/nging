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
package main

import (
	"flag"
	"fmt"
	stdLog "log"
	"net/http"
	"os"
	"strings"

	_ "github.com/admpub/nging/application/library/sqlite"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/engine"
	"github.com/webx-top/echo/engine/standard"
	"github.com/webx-top/echo/handler/mvc/events"
	"github.com/webx-top/echo/middleware"
	"github.com/webx-top/echo/middleware/language"
	"github.com/webx-top/echo/middleware/render"
	"github.com/webx-top/echo/middleware/render/driver"
	"github.com/webx-top/echo/middleware/session"

	"github.com/admpub/log"
	"github.com/admpub/nging/application"
	"github.com/admpub/nging/application/library/config"
	"github.com/admpub/nging/application/library/cron"
	"github.com/admpub/nging/application/library/service"
)

var (
	BUILD_TIME string
	CLOUD_GOX  string
	COMMIT     string
	VERSION    = `1.2.0`

	version    string
	binData    bool
	staticMW   interface{}
	tmplMgr    driver.Manager
	langFSFunc func(dir string) http.FileSystem
)

func main() {
	config.DefaultCLIConfig.InitFlag()
	flag.Parse()
	if binData {
		version = VERSION + ` (bindata)`
	} else {
		version = VERSION
	}
	config.SetVersion(version)

	// Service
	if len(os.Args) > 1 && !strings.HasPrefix(os.Args[1], `-`) {
		config.MustOK(config.InitConfig())
		application.WatchConfig(config.InitConfig, false)
		if err := service.Run(os.Args[1]); err != nil {
			stdLog.Println(err)
		}
		return
	}

	err := config.ParseConfig()
	if err != nil {
		if config.IsInstalled() {
			config.MustOK(err)
		} else {
			log.Error(err)
		}
	}

	if config.DefaultCLIConfig.OnlyRunServer() {
		return
	}

	//Manager
	config.DefaultCLIConfig.RunStartup()

	if config.IsInstalled() {
		// 继续上次任务
		if err := cron.InitJobs(); err != nil {
			log.Error(err)
		}
	}

	e := echo.New()
	if binData {
		e.SetDebug(false)
		log.SetLevel(`Info`)
	} else {
		e.SetDebug(true)
		log.SetLevel(`Debug`)
	}
	e.Use(middleware.Log(), middleware.Recover())
	e.Use(middleware.Gzip(&middleware.GzipConfig{
		Skipper: func(c echo.Context) bool {
			switch c.Request().URL().Path() {
			case `/manage/cmdSend/info`, `/download/progress/info`:
				return true
			}
			return false
		},
	}))
	e.Use(func(h echo.Handler) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Response().Header().Set(`Server`, `nging/`+Version)
			return h.Handle(c)
		}
	})

	// 注册静态资源文件(网站素材文件)
	e.Use(staticMW)

	// 启用session
	e.Use(session.Middleware(config.SessionOptions))

	// 启用多语言支持
	config.DefaultConfig.Language.SetFSFunc(langFSFunc)
	e.Use(language.New(&config.DefaultConfig.Language).Middleware())

	// 启用Validation
	e.Use(middleware.Validate(echo.NewValidation))

	// 注册模板引擎
	renderOptions := &render.Config{
		TmplDir: `./template`,
		Engine:  `standard`,
		ParseStrings: map[string]string{
			`__PUBLIC__`: `/public`,
			`__ASSETS__`: `/public/assets`,
			`__TMPL__`:   `./template`,
		},
		Reload:     true,
		ErrorPages: config.DefaultConfig.Sys.ErrorPages,
	}
	renderOptions.ApplyTo(e)
	if tmplMgr != nil {
		renderOptions.Renderer().SetManager(tmplMgr)
	}
	events.AddEvent(`clearCache`, func(next func(r bool), args ...interface{}) {
		renderOptions.Renderer().ClearCache()
		next(true)
	})

	application.Initialize(e)
	c := &engine.Config{
		Address:     fmt.Sprintf(`:%v`, config.DefaultCLIConfig.Port),
		TLSAuto:     false,
		TLSCacheDir: config.DefaultConfig.Sys.SSLCacheDir,
		TLSCertFile: config.DefaultConfig.Sys.SSLCertFile,
		TLSKeyFile:  config.DefaultConfig.Sys.SSLKeyFile,
	}
	e.Run(standard.NewWithConfig(c))
}
