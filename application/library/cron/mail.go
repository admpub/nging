package cron

import (
	"strings"
	"bytes"

	"github.com/admpub/nging/application/library/cron/send"
	alertRegistry "github.com/admpub/nging/application/registry/alert"
	"github.com/admpub/nging/application/library/cron/writer"
	"github.com/admpub/nging/application/dbschema"
	"github.com/webx-top/echo/param"
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

func OtherSender(params param.Store) error {
	if alertRegistry.SendTopic == nil {
		return nil
	}
	ctx := defaults.NewMockContext()
	return alertRegistry.SendTopic(ctx, `cron`, params)
}

func EmailSender(params param.Store) error {
	user := new(dbschema.NgingUser)
	task, ok := params.Get(`task`).(dbschema.NgingTask)
	if !ok {
		return nil
	}
	err := user.Get(nil, `id`, task.Uid)
	if err != nil || len(user.Email) == 0 {
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
	params["username"] = user.Username
	var content []byte
	if !params.Has(`email-content`) {
		b := new(bytes.Buffer)
		send.MailTpl().Execute(b, params)
		content = b.Bytes()
		params["email-content"] = b.Bytes()
	} else {
		content, _ = params.Get(`email-content`).([]byte)
	}
	return SendMail(user.Email, user.Username, params.String(`title`), content, ccList...)
}

func init() {
	AddSender(EmailSender)
	AddSender(OtherSender)
	alertRegistry.Topics.Add(`cron`, `定时任务`)
}