package sender

import (
	"time"

	"github.com/admpub/log"
	"github.com/admpub/nging/v4/application/registry/alert"
	"github.com/nging-plugins/ddnsmanager/application/library/ddnsmanager/domain/dnsdomain"
	"github.com/webx-top/com"
	"github.com/webx-top/echo/defaults"
	"github.com/webx-top/echo/param"
)

type alarmContent struct {
	title string
	tmpl  map[string]string
	*dnsdomain.TagValues
}

func (a *alarmContent) EmailContent(_ param.Store) []byte {
	if a.tmpl != nil {
		if tp, ok := a.tmpl[`html`]; ok && len(tp) > 0 {
			return []byte(a.TagValues.Parse(tp))
		}
	}
	var content string
	if len(a.IPv4Addr) > 0 && len(a.IPv4Domains) > 0 {
		content += `<p>IPv4: ` + a.IPv4Addr + `<br /><pre>` + a.IPv4Result.String() + `</pre><br />`
	}
	if len(a.IPv6Addr) > 0 && len(a.IPv6Domains) > 0 {
		content += `<p>IPv6: ` + a.IPv6Addr + `<br /><pre>` + a.IPv6Result.String() + `</pre><br />`
	}
	if len(a.Error) > 0 {
		content += `<p><span style="color:red">错误</span>: <br /><pre>` + a.Error + `</pre><br />`
	}
	if len(content) == 0 {
		return nil
	}
	return com.Str2bytes(`<h1>` + a.title + `</h1>` + content + `时间: ` + time.Now().Format(time.RFC3339) + `<br /></p>`)
}

func (a *alarmContent) MarkdownContent(_ param.Store) []byte {
	if a.tmpl != nil {
		if tp, ok := a.tmpl[`markdown`]; ok && len(tp) > 0 {
			return []byte(a.TagValues.Parse(tp))
		}
	}
	var content string
	if len(a.IPv4Addr) > 0 {
		content += `**IPv4**: ` + a.IPv4Addr + "\n**结果**:\n```\n" + a.IPv4Result.String() + "\n```\n"
	}
	if len(a.IPv6Addr) > 0 {
		content += `**IPv6**: ` + a.IPv6Addr + "\n**结果**:\n```\n" + a.IPv6Result.String() + "\n```\n"
	}
	if len(a.Error) > 0 {
		content += "**<font color=\"warning\">错误</font>**:\n```\n" + a.Error + "\n```\n"
	}
	if len(content) == 0 {
		return nil
	}
	return com.Str2bytes(`### ` + a.title + "\n" + content + `**时间**: ` + time.Now().Format(time.RFC3339) + "\n")
}

func Send(v dnsdomain.TagValues, tmpl map[string]string) (err error) {
	ctx := defaults.NewMockContext()
	ct := &alarmContent{
		title:     `[DDNS]IP变更通知`,
		tmpl:      tmpl,
		TagValues: &v,
	}
	alertData := &alert.AlertData{
		Title:   ct.title,
		Content: ct,
		Data:    param.Store{},
	}
	if err = alert.SendTopic(ctx, `ddnsUpdate`, alertData); err != nil {
		log.Warn(`alert.SendTopic: `, err)
	}
	return
}
