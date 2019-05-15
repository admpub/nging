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

package handler

import (
	"io"
	"os"

	"github.com/admpub/log"
	"github.com/admpub/nging/application/dbschema"
	"github.com/admpub/nging/application/library/common"
	"github.com/admpub/nging/application/library/notice"
	"github.com/admpub/nging/application/registry/route"
	"github.com/webx-top/com"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/subdomains"
)

var (
	NewLister            = common.NewLister
	Paging               = common.Paging
	PagingWithPagination = common.PagingWithPagination
	PagingWithLister     = common.PagingWithLister
	PagingWithListerCond = common.PagingWithListerCond
	PagingWithSelectList = common.PagingWithSelectList
	Ok                   = common.Ok
	Err                  = common.Err
	SendErr              = common.SendErr
	SendFail             = common.SendFail
	SendOk               = common.SendOk
)
var (
	WebSocketLogger = log.GetLogger(`websocket`)
	OfficialSQL     string
	GlobalPrefix    string //路由前缀（全局）
	FrontendPrefix  string //路由前缀（前台）
	BackendPrefix   string //路由前缀（后台）
	//=============================
	// 后台路由注册函数
	//=============================

	Echo         = route.Echo
	Apply        = route.Apply
	SetRootGroup = route.SetRootGroup
	Register     = func(fn func(echo.RouteRegister)) {
		route.RegisterToGroup(BackendPrefix, fn)
	}
	Use = func(groupName string, middlewares ...interface{}) {
		if groupName != `*` {
			groupName = BackendPrefix + groupName
		}
		route.Use(groupName, middlewares...)
	}
	RegisterToGroup = func(groupName string, fn func(echo.RouteRegister), middlewares ...interface{}) {
		route.RegisterToGroup(BackendPrefix+groupName, fn, middlewares...)
	}
)

func init() {
	WebSocketLogger.SetLevel(`Info`)
}

func User(ctx echo.Context) *dbschema.User {
	user, _ := ctx.Get(`user`).(*dbschema.User)
	return user
}

func RoleList(ctx echo.Context) []*dbschema.UserRole {
	roleList, _ := ctx.Get(`roleList`).([]*dbschema.UserRole)
	return roleList
}

func Prefix() string {
	return Echo().Prefix() + BackendPrefix
}

func NoticeWriter(ctx echo.Context, noticeType string) (wOut io.Writer, wErr io.Writer, err error) {
	user := User(ctx)
	if user == nil {
		return nil, nil, ctx.Redirect(URLFor(`/login`))
	}
	typ := `service:` + noticeType
	notice.OpenMessage(user.Username, typ)

	wOut = &com.CmdResultCapturer{Do: func(b []byte) error {
		os.Stdout.Write(b)
		notice.Send(user.Username, notice.NewMessageWithValue(typ, noticeType, string(b), notice.Succeed))
		return nil
	}}
	wErr = &com.CmdResultCapturer{Do: func(b []byte) error {
		os.Stderr.Write(b)
		notice.Send(user.Username, notice.NewMessageWithValue(typ, noticeType, string(b), notice.Failed))
		return nil
	}}
	return
}

func URLFor(purl string) string {
	return subdomains.Default.URL(BackendPrefix+purl, `backend`)
}

func FrontendURLFor(purl string) string {
	return subdomains.Default.URL(FrontendPrefix+purl, `frontend`)
}
