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
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/admpub/nging/application/registry/perm"

	"github.com/admpub/nging/application/handler"
	"github.com/admpub/nging/application/library/common"
	"github.com/admpub/nging/application/library/config"
	"github.com/admpub/nging/application/middleware"
	"github.com/admpub/nging/application/registry/navigate"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/handler/captcha"
	"github.com/webx-top/echo/handler/mvc/static/resource"
	"github.com/webx-top/echo/middleware/render"
)

func Initialize(e *echo.Echo) {
	e.Use(middleware.FuncMap(), render.Auto())
	addRouter(e)
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
			return config.ParseConfig()
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
	}, mustOk)
}

func addRouter(e *echo.Echo) {
	opt := captcha.Options{EnableImage: true}
	opt.Wrapper(e)
	e.Get(`/icon`, func(c echo.Context) error {
		return c.Render(`icon`, nil)
	}, middleware.AuthCheck)
	handler.Use(`*`, middleware.AuthCheck)
	handler.Apply(e)
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
	e.Get(`/routeList`, func(ctx echo.Context) error {
		return ctx.JSON(e.Routes())
	}, middleware.AuthCheck)
	e.Get(`/routeNotin`, func(ctx echo.Context) error {
		var unuse []string
		for _, route := range e.Routes() {
			if strings.HasPrefix(route.Path, `/term/client/`) {
				continue
			}
			if strings.HasPrefix(route.Path, `/frp/dashboard/`) {
				continue
			}
			var exists bool
			for _, navGroup := range navigate.TopNavigate {
				for _, navItem := range navGroup.Children {
					var navRoute string
					if len(navItem.Action) > 0 {
						navRoute = `/` + navGroup.Action + `/` + navItem.Action
					} else {
						navRoute = `/` + navGroup.Action
					}
					if navRoute == route.Path {
						exists = true
						break
					}
				}
			}
			if exists {
				continue
			}
			for _, navGroup := range navigate.LeftNavigate {
				for _, navItem := range navGroup.Children {
					var navRoute string
					if len(navItem.Action) > 0 {
						navRoute = `/` + navGroup.Action + `/` + navItem.Action
					} else {
						navRoute = `/` + navGroup.Action
					}
					if navRoute == route.Path {
						exists = true
						break
					}
				}
			}
			if exists {
				continue
			}
			for _, v := range unuse {
				if v == route.Path {
					exists = true
					break
				}
			}
			if exists {
				continue
			}
			_, exists = perm.SpecialAuths[route.Path]
			if exists {
				continue
			}
			unuse = append(unuse, route.Path)
		}

		return ctx.JSON(unuse)
	}, middleware.AuthCheck)
	e.Route(`GET,POST`, `/ping`, func(ctx echo.Context) error {
		header := ctx.Request().Header()
		body := ctx.Request().Body()
		b, _ := ioutil.ReadAll(body)
		body.Close()
		r := echo.H{
			`header`: header.Object(),
			`form`:   echo.NewMapx(ctx.Request().Form().All()).AsStore(),
			`body`:   string(b),
		}
		data := ctx.Data()
		data.SetData(r)
		callback := ctx.Form(`callback`)
		if len(callback) > 0 {
			return ctx.JSONP(callback, data)
		}
		return ctx.JSON(data)
	})
}
