// Copyright 2017 fatedier, fatedier@gmail.com
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package client

import (
	"net/http"
	"time"

	"github.com/admpub/frp/assets"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/engine"
	"github.com/webx-top/echo/engine/standard"
	"github.com/webx-top/echo/middleware"
)

var (
	httpServerReadTimeout  = 10 * time.Second
	httpServerWriteTimeout = 10 * time.Second
)

func (svr *Service) RunAdminServer(address string) (err error) {
	e := echo.New()
	e.Use(middleware.Log(), middleware.Recover(), middleware.Gzip())
	if len(svr.cfg.AdminUser) > 0 && len(svr.cfg.AdminPwd) > 0 {
		e.Use(middleware.BasicAuth(func(user string, passwd string) bool {
			return user == svr.cfg.AdminUser && passwd == svr.cfg.AdminPwd
		}))
	}
	e.Get("/api/reload", svr.APIReload)
	e.Get("/api/status", svr.APIStatus)
	e.Get("/api/config", svr.APIGetConfig)
	e.Put("/api/config", svr.APIPutConfig)

	// view
	err = assets.Load(svr.cfg.AssetsDir, `client`)
	if err != nil {
		panic(err)
	}
	fs := assets.FS(`client`)
	e.Get("/favicon.ico", http.FileServer(fs))
	e.Get("/static*", func(c echo.Context) error {
		file := c.Param("*")
		if len(file) == 0 || file == `/` {
			file = `/index.html`
		}
		return c.File(file, fs)
	})
	e.Get("/", func(c echo.Context) error {
		return c.Redirect("/static/")
	})
	//address := fmt.Sprintf("%s:%d", addr, port)
	cfg := &engine.Config{
		Address:      address,
		ReadTimeout:  httpServerReadTimeout,
		WriteTimeout: httpServerWriteTimeout,
	}
	go e.Run(standard.NewWithConfig(cfg))
	return
}
