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
	"crypto/tls"
	"flag"
	"fmt"
	"strings"

	"github.com/admpub/caddyui/application/library/caddy"
	"github.com/admpub/log"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/engine"
	"github.com/webx-top/echo/engine/standard"
	"github.com/webx-top/echo/middleware"
	"github.com/webx-top/echo/middleware/render"
	"github.com/webx-top/echo/middleware/session"
	"github.com/webx-top/echo/middleware/tplfunc"

	"github.com/admpub/caddyui/application"
	"github.com/admpub/caddyui/application/library/config"
	"github.com/admpub/letsencrypt"
)

var Version = `v0.1.0 beta1`

func main() {
	config.DefaultCLIConfig.InitFlag()
	flag.Parse()

	config.SetVersion(Version)

	config.MustOK(config.ParseConfig())

	switch config.DefaultCLIConfig.Type {
	case `webserver`:
		caddy.TrapSignals()
		config.ParseConfig()
		config.DefaultConfig.Caddy.Init().Start()
		return
	case `ftpserver`:
		config.ParseConfig()
		config.DefaultConfig.FTP.Init().Start()
		return
	}

	if err := config.DefaultCLIConfig.CaddyRestart(); err != nil {
		log.Error(err)
	}
	if err := config.DefaultCLIConfig.FTPRestart(); err != nil {
		log.Error(err)
	}

	e := echo.New()

	e.Use(middleware.Log(), middleware.Recover())

	// 注册静态资源文件
	e.Use(middleware.Static(&middleware.StaticOptions{
		Root: "./public",
		Path: "/public",
	}))

	// 启用session
	e.Use(session.Middleware(config.SessionOptions))

	// 为模板注册常用函数
	e.Use(middleware.FuncMap(tplfunc.TplFuncMap, func(c echo.Context) bool {
		return c.Format() != `html`
	}))

	// 注册模板引擎
	d := render.New(`standard`, `./template`)
	d.Init(true)
	d.SetContentProcessor(func(b []byte) []byte {
		s := string(b)
		s = strings.Replace(s, `__PUBLIC__`, `/public`, -1)
		s = strings.Replace(s, `__ASSETS__`, `/public/assets`, -1)
		s = strings.Replace(s, `__TMPL__`, `./template`, -1)
		return []byte(s)
	})
	e.Use(render.Middleware(d))

	application.Initialize(e)

	c := &engine.Config{
		Address: fmt.Sprintf(`:%v`, config.DefaultCLIConfig.Port),

		TLSCertFile: config.DefaultConfig.Sys.SSLCertFile,
		TLSKeyFile:  config.DefaultConfig.Sys.SSLKeyFile,
	}
	if len(config.DefaultConfig.Sys.SSLHosts) > 0 {
		var tlsManager letsencrypt.Manager
		tlsManager.SetHosts(config.DefaultConfig.Sys.SSLHosts)
		if err := tlsManager.CacheFile(config.DefaultConfig.Sys.SSLCacheFile); err != nil {
			panic(err.Error())
		}
		c.TLSConfig = &tls.Config{
			GetCertificate: tlsManager.GetCertificate,
		}
	}
	e.Run(standard.NewWithConfig(c))

}
