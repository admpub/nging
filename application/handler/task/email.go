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
		data := ctx.NewData()
		if err != nil {
			data.SetError(err)
		}
		return ctx.JSON(data)
	}
	return ctx.Render(`/task/email_test`, nil)
}
