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
package handler

import (
	"io"
	"os"
	"runtime"

	"github.com/admpub/log"
	"github.com/admpub/nging/application/dbschema"
	"github.com/admpub/nging/application/library/common"
	"github.com/admpub/nging/application/library/errors"
	"github.com/admpub/nging/application/library/notice"
	"github.com/webx-top/com"
	"github.com/webx-top/echo"
	"github.com/webx-top/pagination"
)

type handle struct {
	Methods     []string
	Function    interface{}
	Middlewares []interface{}
}

var (
	Handlers      = []func(*echo.Echo){}
	GroupHandlers = map[string][]func(*echo.Group){}
)

func Register(fn func(*echo.Echo)) {
	Handlers = append(Handlers, fn)
}

func RegisterToGroup(groupName string, fn func(*echo.Group)) {
	_, ok := GroupHandlers[groupName]
	if !ok {
		GroupHandlers[groupName] = []func(*echo.Group){}
	}
	GroupHandlers[groupName] = append(GroupHandlers[groupName], fn)
}

var (
	WebSocketLogger = log.GetLogger(`websocket`)
	IsWindows       bool
)

func init() {
	WebSocketLogger.SetLevel(`Info`)
	IsWindows = runtime.GOOS == `windows`
}

func Paging(ctx echo.Context) (page int, size int) {
	return common.Paging(ctx)
}

func PagingWithPagination(ctx echo.Context, delKeys ...string) (page int, size int, rows int, p *pagination.Pagination) {
	return common.PagingWithPagination(ctx, delKeys...)
}

func Ok(v string) errors.Successor {
	return common.Ok(v)
}

func Err(ctx echo.Context, err error) (ret interface{}) {
	return common.Err(ctx, err)
}

func SendOk(ctx echo.Context, msg string) {
	ctx.Session().AddFlash(Ok(msg))
}

func SendFail(ctx echo.Context, msg string) {
	ctx.Session().AddFlash(msg)
}

func User(ctx echo.Context) *dbschema.User {
	user, _ := ctx.Get(`user`).(*dbschema.User)
	return user
}

func NoticeWriter(ctx echo.Context, noticeType string) (wOut io.Writer, wErr io.Writer, err error) {
	user := User(ctx)
	if user == nil {
		return nil, nil, ctx.Redirect(`/login`)
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
