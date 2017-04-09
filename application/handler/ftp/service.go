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
package ftp

import (
	"os"

	"github.com/admpub/nging/application/handler"
	"github.com/admpub/nging/application/library/config"
	"github.com/admpub/nging/application/library/notice"
	"github.com/webx-top/echo"
)

func init() {
	handler.RegisterToGroup(`/manage`, func(g *echo.Group) {
		g.Route(`GET,POST`, `/ftp_restart`, Restart)
		g.Route(`GET,POST`, `/ftp_stop`, Stop)
		g.Route(`GET,POST`, `/ftp_log`, Log)
	})
}

func Restart(ctx echo.Context) error {
	wOut, wErr, err := handler.NoticeWriter(ctx, ctx.T(`FTP服务`))
	if err != nil {
		return ctx.String(err.Error())
	}
	if err := config.DefaultCLIConfig.FTPRestart(wOut, wErr); err != nil {
		return ctx.String(err.Error())
	}
	return ctx.String(ctx.T(`已经重启FTP服务`))
}

func Stop(ctx echo.Context) error {
	if err := config.DefaultCLIConfig.FTPStop(); err != nil {
		return ctx.String(err.Error())
	}
	return ctx.String(ctx.T(`已经关闭FTP服务`))
}

func Log(ctx echo.Context) error {
	on := ctx.Formx(`on`).Bool()
	if on {
		wOut, wErr, err := handler.NoticeWriter(ctx, ctx.T(`FTP服务`))
		if err != nil {
			return ctx.String(err.Error())
		}
		err = config.DefaultCLIConfig.SetLogWriter(`ftp`, wOut, wErr)
		if err != nil {
			return ctx.String(err.Error())
		}
		return ctx.String(ctx.T(`已经开始直播FTP服务状态`))
	}
	err := config.DefaultCLIConfig.SetLogWriter(`ftp`, os.Stdout, os.Stderr)
	if err != nil {
		return ctx.String(err.Error())
	}
	user := handler.User(ctx)
	if user == nil {
		return ctx.String(ctx.T(`请先登录`))
	}
	typ := `service:` + ctx.T(`FTP服务`)
	notice.CloseMessage(user.Username, typ)
	return ctx.String(ctx.T(`已经停止直播FTP服务状态`))
}
