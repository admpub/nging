package send

import (
	"errors"
	"html/template"

	"github.com/admpub/nging/application/library/notice"

	"github.com/admpub/mail"
	"github.com/admpub/nging/application/library/email"
)

var (
	htmlTpl, markdownTpl *template.Template
	defaultHTMLTmpl      = `
	你好 {{.username}}，<br/>
<p>以下是任务执行结果：</p>
<p>
任务编号：{{.task.Id}}<br/>
任务名称：{{.task.Name}}<br/>
执行时间：{{.startTime}}<br />
执行耗时：{{.elapsed}}秒<br />
执行状态：<span style="color:{{if eq "success" .status}}green{{else if eq "failure" .status}}red{{else}}orange{{end}}">{{.statusText}}</span><br />
详细信息：<a href="{{.detailURL}}" target="_blank">查看</a>
</p>
<p>-------------以下是任务执行输出-------------</p>
<p><pre>{{.output}}</pre></p>
<p>
--------------------------------------------<br />
本邮件由系统自动发出，请勿回复<br />
如果要取消邮件通知，请登录到系统进行设置<br />
</p>
`
	defaultMarkdownTmpl = `
### 任务执行结果
**任务编号**：{{.task.Id}}
**任务名称**：{{.task.Name}}
**执行时间**：{{.startTime}}
**执行耗时**：{{.elapsed}}秒
**执行状态**：<font color="{{if eq "success" .status}}info{{else if eq "failure" .status}}warning{{else}}warning{{end}}">{{.statusText}}</font>
**详细信息**：[查看]({{.detailURL}})

### 以下是任务执行输出
{{.output}}
`
	DefaultSMTPConfig     = &mail.SMTPConfig{} //STMP配置
	DefaultEmailConfig    = &EmailConfig{}
	ErrIncorrectRecipient = errors.New(`The recipient's email address is incorrect`)
)

type EmailConfig struct {
	Template   string
	Sender     string
	Engine     string
	Timeout    int64
	QueueSize  int
	TemplateMd string
}

func InitialHTMLTmpl() {
	tmpl := DefaultEmailConfig.Template
	if len(tmpl) == 0 {
		tmpl = defaultHTMLTmpl
	}
	var err error
	htmlTpl, err = template.New("notifyHTMLTmpl").Parse(tmpl)
	if err != nil {
		panic(err)
	}
}

func InitialMarkdownTmpl() {
	tmpl := DefaultEmailConfig.TemplateMd
	if len(tmpl) == 0 {
		tmpl = defaultMarkdownTmpl
	}
	var err error
	markdownTpl, err = template.New("notifyMarkdownTmpl").Parse(tmpl)
	if err != nil {
		panic(err)
	}
}

func MarkdownTmpl() *template.Template {
	if markdownTpl == nil {
		InitialMarkdownTmpl()
	}
	return markdownTpl
}

func HTMLTmpl() *template.Template {
	if htmlTpl == nil {
		InitialHTMLTmpl()
	}
	return htmlTpl
}

func MailTpl() *template.Template {
	return HTMLTmpl()
}

// Mail 发送Email
// @param toEmail 收信邮箱
// @param toUsername 收信人名称
// @param title 邮件标题
// @param content 邮件内容
// @param ccList 抄送地址
func Mail(toEmail string, toUsername string, title string, content []byte, ccList ...string) error {
	return MailWithID(0, toEmail, toUsername, title, content, ccList...)
}

// MailWithID 发送Email(带ID参数)
func MailWithID(id uint64, toEmail string, toUsername string, title string, content []byte, ccList ...string) error {
	return MailWithIDAndNoticer(id, nil, toEmail, toUsername, title, content, ccList...)
}

// MailWithNoticer 发送Email(带Noticer参数)
func MailWithNoticer(noticer notice.Noticer, toEmail string, toUsername string, title string, content []byte, ccList ...string) error {
	return MailWithIDAndNoticer(0, noticer, toEmail, toUsername, title, content, ccList...)
}

// MailWithIDAndNoticer 发送Email(带ID和Noticer参数)
func MailWithIDAndNoticer(id uint64, noticer notice.Noticer, toEmail string, toUsername string, title string, content []byte, ccList ...string) error {
	if len(toEmail) < 1 { //收信人邮箱地址不正确
		return ErrIncorrectRecipient
	}
	conf := &email.Config{
		ID:         id,
		Engine:     DefaultEmailConfig.Engine,
		SMTP:       DefaultSMTPConfig,
		From:       DefaultEmailConfig.Sender,
		ToAddress:  toEmail,
		ToUsername: toUsername,
		Subject:    title,
		Content:    content,
		CcAddress:  ccList,
		Timeout:    DefaultEmailConfig.Timeout,
		Noticer:    noticer,
	}
	return email.SendMail(conf)
}
