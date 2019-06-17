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
	"net/http"

	"github.com/admpub/events"
	"github.com/admpub/log"

	"github.com/admpub/events/emitter"
	"github.com/admpub/nging/application/cmd/event"
	"github.com/admpub/nging/application/handler"
	"github.com/admpub/nging/application/library/config"
	ngingMW "github.com/admpub/nging/application/middleware"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/handler/pprof"
	"github.com/webx-top/echo/middleware"
	"github.com/webx-top/echo/middleware/language"
	"github.com/webx-top/echo/middleware/render"
	"github.com/webx-top/echo/middleware/render/driver"
	"github.com/webx-top/echo/middleware/session"
	"github.com/webx-top/echo/subdomains"
)

const (
	DefaultTemplateDir   = `./template/backend`
	DefaultAssetsDir     = `./public/assets`
	DefaultAssetsURLPath = `/public/assets/backend`
)

var (
	TemplateDir      = DefaultTemplateDir //模板文件夹
	AssetsDir        = DefaultAssetsDir   //素材文件夹
	AssetsURLPath    = DefaultAssetsURLPath
	DefaultAvatarURL = AssetsURLPath + `/images/user_128.png`
	RendererDo       = func(driver.Driver) {}
	ParseStrings     = map[string]string{}
	ParseStringFuncs = map[string]func() string{}
)

func init() {
	echo.Set(`BackendPrefix`, handler.BackendPrefix)
	echo.Set(`GlobalPrefix`, handler.GlobalPrefix)
	event.OnStart(0, func() {
		handler.GlobalPrefix = echo.String(`GlobalPrefix`)
		handler.BackendPrefix = echo.String(`BackendPrefix`)
		handler.FrontendPrefix = echo.String(`FrontendPrefix`)
		ngingMW.DefaultAvatarURL = DefaultAssetsURLPath
		e := handler.Echo()
		e.SetPrefix(handler.GlobalPrefix)
		handler.SetRootGroup(handler.BackendPrefix)
		subdomains.Default.Default = `backend`
		domainName := subdomains.Default.Default
		if len(config.DefaultCLIConfig.BackendDomain) > 0 {
			domainName += `@` + config.DefaultCLIConfig.BackendDomain
		}
		subdomains.Default.Add(domainName, e)

		e.Use(middleware.Log(), middleware.Recover())
		skippedGzipPaths := map[string]bool{
			e.Prefix() + `/server/cmdSend/info`:      true,
			e.Prefix() + `/download/progress/info`:   true,
			e.Prefix() + `/debug/pprof/`:             true,
			e.Prefix() + `/debug/pprof/allocs`:       true,
			e.Prefix() + `/debug/pprof/block`:        true,
			e.Prefix() + `/debug/pprof/cmdline`:      true,
			e.Prefix() + `/debug/pprof/goroutine`:    true,
			e.Prefix() + `/debug/pprof/heap`:         true,
			e.Prefix() + `/debug/pprof/mutex`:        true,
			e.Prefix() + `/debug/pprof/profile`:      true,
			e.Prefix() + `/debug/pprof/threadcreate`: true,
			e.Prefix() + `/debug/pprof/trace`:        true,
		}
		e.Use(middleware.Gzip(&middleware.GzipConfig{
			Skipper: func(c echo.Context) bool {
				upath := c.Request().URL().Path()
				skipped, _ := skippedGzipPaths[upath]
				return skipped
			},
		}))
		e.Use(func(h echo.Handler) echo.HandlerFunc {
			return func(c echo.Context) error {
				c.Response().Header().Set(`Server`, `nging/`+config.Version.Number)
				return h.Handle(c)
			}
		})

		// 注册静态资源文件(网站素材文件)
		e.Use(event.StaticMW) //打包的静态资源
		// 上传文件资源(改到manager中用File函数实现)
		// e.Use(middleware.Static(&middleware.StaticOptions{
		// 	Root: helper.UploadDir,
		// 	Path: helper.UploadURLPath,
		// }))

		// 启用session
		e.Use(session.Middleware(config.SessionOptions))
		// 启用多语言支持
		config.DefaultConfig.Language.SetFSFunc(event.LangFSFunc)
		e.Use(language.New(&config.DefaultConfig.Language).Middleware())

		// 启用Validation
		e.Use(middleware.Validate(echo.NewValidation))

		// 事物支持
		e.Use(ngingMW.Tansaction())
		// 注册模板引擎
		renderOptions := &render.Config{
			TmplDir: TemplateDir,
			Engine:  `standard`,
			ParseStrings: map[string]string{
				`__ASSETS__`: AssetsURLPath,
				`__TMPL__`:   TemplateDir,
			},
			ParseStringFuncs: map[string]func() string{
				`__BACKEND__`:  func() string { return subdomains.Default.URL(handler.BackendPrefix, `backend`) },
				`__FRONTEND__`: func() string { return subdomains.Default.URL(handler.FrontendPrefix, `frontend`) },
			},
			DefaultHTTPErrorCode: http.StatusOK,
			Reload:               true,
			ErrorPages:           config.DefaultConfig.Sys.ErrorPages,
		}
		if ParseStrings != nil {
			for key, val := range ParseStrings {
				renderOptions.ParseStrings[key] = val
			}
		}
		if ParseStringFuncs != nil {
			for key, val := range ParseStringFuncs {
				renderOptions.ParseStringFuncs[key] = val
			}
		}
		if RendererDo != nil {
			renderOptions.AddRendererDo(RendererDo)
		}
		renderOptions.AddFuncSetter(ngingMW.ErrorPageFunc)
		renderOptions.ApplyTo(e, event.BackendTmplMgr)
		RendererDo(renderOptions.Renderer())
		emitter.DefaultCondEmitter.On(`clearCache`, events.Callback(func(_ events.Event) error {
			log.Debug(`clear: Backend Template Object Cache`)
			renderOptions.Renderer().ClearCache()
			return nil
		}))
		e.Get(`/favicon.ico`, event.FaviconHandler)
		if event.Develop {
			pprof.Wrap(e)
		}
		Initialize()
	})
}
