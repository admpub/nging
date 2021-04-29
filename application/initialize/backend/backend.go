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
	"fmt"
	"net/http"
	"strings"

	"github.com/webx-top/com"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/handler/pprof"
	"github.com/webx-top/echo/middleware"
	"github.com/webx-top/echo/middleware/language"
	"github.com/webx-top/echo/middleware/render"
	"github.com/webx-top/echo/middleware/render/driver"
	"github.com/webx-top/echo/middleware/session"
	"github.com/webx-top/echo/subdomains"

	"github.com/admpub/log"
	"github.com/admpub/nging/application/cmd/event"
	"github.com/admpub/nging/application/handler"
	"github.com/admpub/nging/application/library/common"
	"github.com/admpub/nging/application/library/config"
	ngingMW "github.com/admpub/nging/application/middleware"
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
	SkippedGzipPaths = map[string]bool{}
	GzipSkipper      = func(skippedPaths map[string]bool) func(c echo.Context) bool {
		return func(c echo.Context) bool {
			skipped, _ := skippedPaths[c.Request().URL().Path()]
			if !skipped {
				skipped, _ = skippedPaths[c.Path()]
			}
			return skipped
		}
	}
	DefaultLocalHostNames = []string{
		`127.0.0.1`, `localhost`,
	}
)

func MakeSubdomains(domain string, appends []string) string {
	domainList := strings.Split(domain, `,`)
	domain = domainList[0]
	if pos := strings.Index(domain, `://`); pos > 0 {
		pos += 3
		if pos < len(domain) {
			domain = domain[pos:]
		} else {
			domain = ``
		}
	}
	var myPort string
	domain, myPort = com.SplitHost(domain)
	if len(myPort) == 0 && len(domainList) > 1 {
		_, myPort = com.SplitHost(domainList[1])
	}
	port := fmt.Sprintf("%d", config.DefaultCLIConfig.Port)
	newDomainList := []string{}
	if !com.InSlice(domain+`:`+port, domainList) {
		newDomainList = append(newDomainList, domain+`:`+port)
	}
	if myPort == port {
		myPort = ``
	}
	if len(myPort) > 0 {
		if !com.InSlice(domain+`:`+myPort, domainList) {
			newDomainList = append(newDomainList, domain+`:`+myPort)
		}
	}
	for _, hostName := range appends {
		if hostName == domain {
			continue
		}
		if !com.InSlice(hostName+`:`+port, domainList) {
			newDomainList = append(newDomainList, hostName+`:`+port)
		}
		if len(myPort) > 0 {
			if !com.InSlice(hostName+`:`+myPort, domainList) {
				newDomainList = append(newDomainList, hostName+`:`+myPort)
			}
		}
	}
	if len(newDomainList) > 0 {
		domainList = append(domainList, newDomainList...)
	}
	return strings.Join(domainList, `,`)
}

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
		subdomains.Default.Boot = `backend`
		domainName := subdomains.Default.Default
		backendDomain := config.DefaultCLIConfig.BackendDomain
		if len(backendDomain) > 0 {
			domainName += `@` + MakeSubdomains(backendDomain, DefaultLocalHostNames)
		}
		subdomains.Default.Add(domainName, e)

		e.Use(middleware.Log(), middleware.Recover())
		e.Use(middleware.Gzip(&middleware.GzipConfig{
			Skipper: GzipSkipper(SkippedGzipPaths),
		}))
		e.Use(func(h echo.Handler) echo.HandlerFunc {
			return func(c echo.Context) error {
				c.Response().Header().Set(`Server`, event.SoftwareName+`/`+config.Version.Number)
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
		e.Use(ngingMW.Transaction())
		// 注册模板引擎
		renderOptions := &render.Config{
			TmplDir: TemplateDir,
			Engine:  `standard`,
			ParseStrings: map[string]string{
				`__TMPL__`: TemplateDir,
			},
			DefaultHTTPErrorCode: http.StatusOK,
			Reload:               true,
			ErrorPages:           config.DefaultConfig.Sys.ErrorPages,
			ErrorProcessors:      common.ErrorProcessors,
		}
		for key, val := range ParseStrings {
			renderOptions.ParseStrings[key] = val
		}
		for key, val := range ParseStringFuncs {
			renderOptions.ParseStringFuncs[key] = val
		}
		if RendererDo != nil {
			renderOptions.AddRendererDo(RendererDo)
		}
		renderOptions.AddFuncSetter(BackendURLFunc)
		renderOptions.AddFuncSetter(ngingMW.ErrorPageFunc)
		renderOptions.ApplyTo(e, event.BackendTmplMgr)
		//RendererDo(renderOptions.Renderer())
		echo.On(`clearCache`, func(_ echo.H) error {
			log.Debug(`clear: Backend Template Object Cache`)
			renderOptions.Renderer().ClearCache()
			return nil
		})
		e.Get(`/favicon.ico`, event.FaviconHandler)
		if event.Develop {
			pprof.Wrap(e)
		}
		Initialize()
	})
}
