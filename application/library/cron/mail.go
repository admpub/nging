package cron

import (
	"strings"

	"github.com/admpub/nging/v5/application/dbschema"
	"github.com/admpub/nging/v5/application/library/cron/send"
	"github.com/admpub/nging/v5/application/library/cron/writer"
	"github.com/admpub/nging/v5/application/registry/alert"
	"github.com/webx-top/echo/defaults"
)

var (
	// SendMail 发送Email
	// @param toEmail 收信邮箱
	// @param toUsername 收信人名称
	// @param title 邮件标题
	// @param content 邮件内容
	// @param ccList 抄送地址
	SendMail = send.Mail

	// SendMailWithID 发送Email(带ID参数)
	SendMailWithID = send.MailWithID

	// SendMailWithNoticer 发送Email(带Noticer参数)
	SendMailWithNoticer = send.MailWithNoticer

	// SendMailWithIDAndNoticer 发送Email(带ID和Noticer参数)
	SendMailWithIDAndNoticer = send.MailWithIDAndNoticer

	NewCmdRec = writer.New
)

type OutputWriter = writer.OutputWriter

func EmailSender(alertData *alert.AlertData) error {
	params := alertData.Data
	task, ok := params.Get(`task`).(dbschema.NgingTask)
	if !ok {
		return nil
	}
	var ccList []string
	if len(task.NotifyEmail) > 0 {
		ccList = strings.Split(task.NotifyEmail, "\n")
		for index, email := range ccList {
			email = strings.TrimSpace(email)
			if len(email) == 0 {
				continue
			}
			ccList[index] = email
		}
	}
	if len(ccList) == 0 {
		return nil
	}
	toEmail := ccList[0]
	toUsername := strings.SplitN(toEmail, "@", 2)[0]
	if len(ccList) > 1 {
		ccList = append([]string{}, ccList[1:]...)
	} else {
		ccList = []string{}
	}
	ct := alertData.Content
	content := ct.EmailContent(params)
	return SendMail(toEmail, toUsername, alertData.Title, content, ccList...)
}

func OtherSender(alertData *alert.AlertData) error {
	ctx := defaults.NewMockContext()
	return alert.SendTopic(ctx, `cron`, alertData)
}

func init() {
	AddSender(EmailSender)
	AddSender(OtherSender)
	alert.Topics.Add(`cron`, `定时任务`)
}
