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

	"github.com/admpub/nging/application/library/cron"
	"github.com/webx-top/echo"
)

//EmailTest 测试邮件发送是否正常
func EmailTest(ctx echo.Context) error {
	if ctx.IsPost() {
		toEmail := ctx.Form(`email`)
		toUsername := `test`
		title := ctx.T(`恭喜！邮件发送功能正常`)
		content := []byte(ctx.T(`如果您收到这封邮件，说明邮件发送功能正常。<br /><br /> 来自：%s<br />时间：%s`, ctx.Site(), time.Now().Format(time.RFC3339)))
		err := cron.SendMail(toEmail, toUsername, title, content)
		data := ctx.Data()
		if err != nil {
			data.SetError(err)
		}
		return ctx.JSON(data)
	}
	return ctx.Render(`/task/email_test`, nil)
}
