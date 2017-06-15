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
		content := []byte(ctx.T(`如果您收到这封邮件，说明邮件发送功能正常。<br /><br /> 来自：%s<br />时间：%s`, ctx.Site(), time.Now().String()))
		err := cron.SendMail(toEmail, toUsername, title, content)
		data := ctx.Data()
		if err != nil {
			data.SetError(err)
		}
		return ctx.JSON(data)
	}
	return ctx.Render(`/task/email_test`, nil)
}
