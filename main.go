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
	"io/ioutil"
	"strconv"
	"strings"

	assetfs "github.com/elazarl/go-bindata-assetfs"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/engine"
	"github.com/webx-top/echo/engine/standard"
	"github.com/webx-top/echo/handler/mvc/events"
	"github.com/webx-top/echo/logger"
	"github.com/webx-top/echo/middleware"
	"github.com/webx-top/echo/middleware/render"
	"github.com/webx-top/echo/middleware/session"
	"github.com/webx-top/echo/middleware/tplfunc"

	"github.com/admpub/letsencrypt"
	"github.com/admpub/nging/application"
	"github.com/admpub/nging/application/library/config"
)

var Version = `v0.1.0 beta1`
var BindData = `1`

type assetManager struct {
	*assetfs.AssetFS
}

func (a *assetManager) Close()                                            {}
func (a *assetManager) SetOnChangeCallback(func(name, typ, event string)) {}
func (a *assetManager) SetLogger(logger.Logger)                           {}
func (a *assetManager) ClearCache()                                       {}
func (a *assetManager) GetTemplate(fileName string) ([]byte, error) {
	file, err := a.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	b, err := ioutil.ReadAll(file)
	return b, err
}
func (a *assetManager) Init(logger logger.Logger, rootDir string, reload bool, allows ...string) {}

func main() {
	config.DefaultCLIConfig.InitFlag()
	flag.Parse()

	config.SetVersion(Version)

	config.MustOK(config.ParseConfig())

	if config.DefaultCLIConfig.OnlyRunServer() {
		return
	}
	config.DefaultCLIConfig.RunStartup()

	e := echo.New()
	e.SetHTTPErrorHandler(render.HTTPErrorHandler(config.DefaultConfig.Sys.ErrorPages))
	e.Use(middleware.Log(), middleware.Recover())

	bindData, _ := strconv.ParseBool(BindData)

	// 注册静态资源文件
	if bindData {
		asset := &assetfs.AssetFS{Asset: Asset, AssetDir: AssetDir, Prefix: ""}
		e.Get("/public/*", func(c echo.Context) error {
			fileName := c.Request().URL().Path()
			file, err := asset.Open(fileName)
			if err != nil {
				return echo.ErrNotFound
			}
			defer file.Close()
			info, err := file.Stat()
			if err != nil {
				return echo.ErrNotFound
			}
			return c.ServeContent(file, info.Name(), info.ModTime())
		})
	} else {
		e.Use(middleware.Static(&middleware.StaticOptions{
			Root: "./public",
			Path: "/public",
		}))
	}

	// 启用session
	e.Use(session.Middleware(config.SessionOptions))

	// 为模板注册常用函数
	e.Use(middleware.FuncMap(tplfunc.TplFuncMap, func(c echo.Context) bool {
		return c.Format() != `html`
	}))

	// 注册模板引擎
	d := render.New(`standard`, `./template`)
	d.Init(true)
	if bindData {
		d.SetManager(&assetManager{AssetFS: &assetfs.AssetFS{Asset: Asset, AssetDir: AssetDir, Prefix: "template"}})
	}
	d.SetContentProcessor(func(b []byte) []byte {
		s := string(b)
		s = strings.Replace(s, `__PUBLIC__`, `/public`, -1)
		s = strings.Replace(s, `__ASSETS__`, `/public/assets`, -1)
		s = strings.Replace(s, `__TMPL__`, `./template`, -1)
		return []byte(s)
	})
	events.AddEvent(`clearCache`, func(next func(r bool), args ...interface{}) {
		d.ClearCache()
		next(true)
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
