/*
   Nging is a toolbox for webmasters
   Copyright (C) 2018-present Wenhui Shen <swh@admpub.com>

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
	"strings"

	"github.com/webx-top/com"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/subdomains"

	"github.com/admpub/log"
	"github.com/admpub/nging/v5/application/dbschema"
	"github.com/admpub/nging/v5/application/library/common"
	"github.com/admpub/nging/v5/application/library/notice"
	"github.com/admpub/nging/v5/application/library/sessionguard"
	"github.com/admpub/nging/v5/application/registry/route"
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
	GetRoleList          = func(echo.Context) []*dbschema.NgingUserRole {
		return nil
	}
)
var (
	WebSocketLogger = log.GetLogger(`websocket`)
	GlobalPrefix    string //路由前缀（全局）
	FrontendPrefix  string //路由前缀（前台）
	BackendPrefix   string //路由前缀（后台）
	//=============================
	// 后台路由注册函数
	//=============================

	IRegister          = route.IRegister
	WithMeta           = route.MetaHandler
	WithMetaAndRequest = route.MetaHandlerWithRequest
	WithRequest        = route.HandlerWithRequest
	Apply              = route.Apply
	SetRootGroup       = route.SetRootGroup
	PublicHandler      = route.PublicHandler
	GuestHandler       = route.GuestHandler
	Pre                = route.Pre
	PreToGroup         = route.PreToGroup
	Use                = route.Use
	// UseToGroup “@”符号代表后台网址前缀
	UseToGroup = func(groupName string, middlewares ...interface{}) {
		if groupName != `*` {
			groupName = `@` + groupName
		}
		route.UseToGroup(groupName, middlewares...)
	}
	Register = func(fn func(echo.RouteRegister)) {
		route.RegisterToGroup(`@`, fn)
	}
	RegisterToGroup = func(groupName string, fn func(echo.RouteRegister), middlewares ...interface{}) {
		route.RegisterToGroup(`@`+groupName, fn, middlewares...)
	}
	Host = route.Host
)

func init() {
	WebSocketLogger.SetLevel(`Info`)
	route.AddGroupNamer(func(group string) string {
		if len(group) == 0 {
			return group
		}
		if group == `@` {
			return BackendPrefix
		}
		if strings.HasPrefix(group, `@`) {
			return BackendPrefix + group[1:]
		}
		return group
	})
}

func User(ctx echo.Context) *dbschema.NgingUser {
	user, ok := ctx.Internal().Get(`user`).(*dbschema.NgingUser)
	if ok && user != nil {
		return user
	}
	user, ok = ctx.Session().Get(`user`).(*dbschema.NgingUser)
	if ok {
		if !sessionguard.Validate(ctx, user.LastIp, `user`, uint64(user.Id)) {
			log.Warn(ctx.T(`用户“%s”的会话环境发生改变，需要重新登录`, user.Username))
			ctx.Session().Delete(`user`)
			return nil
		}
		ctx.Internal().Set(`user`, user)
	}
	return user
}

func Prefix() string {
	return IRegister().Prefix() + BackendPrefix
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
