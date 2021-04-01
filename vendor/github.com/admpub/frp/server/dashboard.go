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

package server

import (
	"time"

	"github.com/admpub/frp/assets"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/webx-top/echo"
	echoPrometheus "github.com/webx-top/echo-prometheus"
	"github.com/webx-top/echo/engine"
	"github.com/webx-top/echo/engine/standard"
	"github.com/webx-top/echo/middleware"
)

var (
	httpServerReadTimeout  = 10 * time.Second
	httpServerWriteTimeout = 10 * time.Second
)

func (svr *Service) RunDashboardServer(address string) (err error) {
	// url router
	e := echo.New()
	e.Use(middleware.Log(), middleware.Recover(), middleware.Gzip())
	if len(svr.cfg.DashboardUser) > 0 && len(svr.cfg.DashboardPwd) > 0 {
		e.Use(middleware.BasicAuth(func(user string, passwd string) bool {
			return user == svr.cfg.DashboardUser && passwd == svr.cfg.DashboardPwd
		}))
	}
	// metrics
	if svr.cfg.EnablePrometheus {
		// Prometheus
		e.Use(echoPrometheus.MetricsMiddleware())
		e.Get("/metrics", echo.WrapHandler(promhttp.Handler()))
	}
	svr.RegisterTo(e)
	// view
	fs := assets.FS(`server`)
	e.Get("/favicon.ico", func(c echo.Context) error {
		return c.File(c.Path(), fs)
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
