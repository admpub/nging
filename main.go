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

	"github.com/webx-top/echo"
	"github.com/webx-top/echo/engine/standard"
	"github.com/webx-top/echo/middleware"
	"github.com/webx-top/echo/middleware/render"
	"github.com/webx-top/echo/middleware/session"

	"strings"

	"github.com/admpub/caddyui/application"
	"github.com/admpub/caddyui/application/library/config"
)

func main() {
	config.DefaultCLIConfig.InitFlag()
	flag.Parse()

	config.MustOK(config.ParseConfig())

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
	e.Use(middleware.FuncMap(nil, func(c echo.Context) bool {
		return c.Format() != `html`
	}))

	// 注册模板引擎
	d := render.New(`standard`, `./template`)
	d.Init(true)
	d.SetContentProcessor(func(b []byte) []byte {
		s := string(b)
		s = strings.Replace(s, `__PUBLIC__`, `/public`, -1)
		s = strings.Replace(s, `__ASSETS__`, `/public/assets`, -1)
		return []byte(s)
	})
	e.Use(render.Middleware(d))

	application.Initialize(e)
	e.Run(standard.New(fmt.Sprintf(`:%v`, config.DefaultCLIConfig.Port)))
}
