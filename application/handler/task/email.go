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
package task

import (
	"time"

	"github.com/admpub/nging/v4/application/handler"
	"github.com/admpub/nging/v4/application/library/common"
	"github.com/admpub/nging/v4/application/library/cron"
	"github.com/admpub/nging/v4/application/library/notice"
	"github.com/webx-top/echo"
)

//EmailTest 测试邮件发送是否正常
func EmailTest(ctx echo.Context) error {
	user := handler.User(ctx)
	if user == nil {
		return common.ErrUserNotLoggedIn
	}
	if ctx.IsPost() {
		clientID := ctx.Formx(`clientID`).String()
		if len(clientID) == 0 {
			return ctx.E(`clientID值不正确`)
		}
		toEmail := ctx.Form(`email`)
		toUsername := `test`
		title := ctx.T(`恭喜！邮件发送功能正常`)
		content := []byte(ctx.T(`如果您收到这封邮件，说明邮件发送功能正常。<br /><br /> 来自：%s<br />时间：%s`, ctx.Site(), time.Now().Format(time.RFC3339)))
		noticerConfig := &notice.HTTPNoticerConfig{
			User:     user.Username,
			Type:     `emailTest`,
			ClientID: clientID,
		}
		err := cron.SendMailWithNoticer(notice.NewNoticer(ctx, noticerConfig), toEmail, toUsername, title, content)
		data := ctx.Data()
		if err != nil {
			data.SetError(err)
		}
		return ctx.JSON(data)
	}
	return ctx.Render(`/task/email_test`, nil)
}
