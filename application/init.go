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
	sockjsHandler "github.com/webx-top/echo/handler/sockjs"
	ws "github.com/webx-top/echo/handler/websocket"
	"github.com/webx-top/echo/middleware/render"

	. "github.com/admpub/nging/application/handler"
	"github.com/admpub/nging/application/library/config"
	"github.com/admpub/nging/application/middleware"
)

var DefaultRequestMethods = []string{echo.GET}
var DefaultFormPageMethods = []string{echo.GET, echo.POST}

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
	addHandler(e, `/`, Index)
	addHandler(e, `/logout`, Logout)
	addHandler(e, `/addon_form`, AddonForm)

	addFormHandler(e, `/login`, Login)
	addFormHandler(e, `/setup`, Setup)

	g := e.Group(`/manage`, middleware.AuthCheck)
	{
		addHandler(g, ``, ManageIndex)
		addFormHandler(g, `/vhost_add`, ManageVhostAdd)
		addFormHandler(g, `/vhost_edit`, ManageVhostEdit)
		addFormHandler(g, `/vhost_delete`, ManageVhostDelete)
		addFormHandler(g, `/vhost_file`, ManageVhostFile)
		addFormHandler(g, `/service`, ManageService)
		addFormHandler(g, `/web_restart`, ManageWebRestart)
		addFormHandler(g, `/web_stop`, ManageWebStop)
		addFormHandler(g, `/ftp_restart`, ManageFTPRestart)
		addFormHandler(g, `/web_log`, ManageWebLog)
		addFormHandler(g, `/ftp_log`, ManageFTPLog)
		addFormHandler(g, `/ftp_stop`, ManageFTPStop)
		addFormHandler(g, `/clear_cache`, ManageClearCache)
		addFormHandler(g, `/execmd`, ManageSysCmd)

		sockjsOpts := sockjsHandler.Options{
			Handle: ManageSockJSSendCmd,
			Prefix: "/execmd_send",
		}
		sockjsOpts.Wrapper(g)
		_ = sockjsOpts
		wsOpts := ws.Options{
			Handle: ManageWSSendCmd,
			Prefix: "/execmd_send_ws",
		}
		wsOpts.Wrapper(g)

		addHandler(g, `/sysinfo`, ManageSysInfo, render.AutoOutput(nil))

		wsOpts = ws.Options{
			Handle: ManageNotice,
			Prefix: "/notice",
		}
		wsOpts.Wrapper(g)
	}

	g = e.Group(`/ftp`, middleware.AuthCheck)
	{
		addHandler(g, `/account`, FTPAccountIndex)
		addFormHandler(g, `/account_add`, FTPAccountAdd)
		addFormHandler(g, `/account_edit`, FTPAccountEdit)
		addFormHandler(g, `/account_delete`, FTPAccountDelete)

		addHandler(g, `/group`, FTPGroupIndex)
		addFormHandler(g, `/group_add`, FTPGroupAdd)
		addFormHandler(g, `/group_edit`, FTPGroupEdit)
		addFormHandler(g, `/group_delete`, FTPGroupDelete)
	}

	addFormHandler(e, `/db`, DbManager, middleware.AuthCheck)

	opt := captcha.Options{EnableImage: true}
	opt.Wrapper(e)
}

func addFormHandler(rr echo.RouteRegister, urlPath string, handler interface{}, middlewares ...interface{}) {
	rr.Match(DefaultFormPageMethods, urlPath, handler, middlewares...)
}

func addHandler(rr echo.RouteRegister, urlPath string, handler interface{}, middlewares ...interface{}) {
	rr.Match(DefaultRequestMethods, urlPath, handler, middlewares...)
}
