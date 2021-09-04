package alert

import (
	"strings"

	"github.com/admpub/nging/v3/application/dbschema"
	"github.com/admpub/nging/v3/application/library/cron/send"
	"github.com/admpub/nging/v3/application/library/imbot"
	"github.com/admpub/nging/v3/application/registry/alert"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/param"
)

// Send 发送警报
func Send(a *dbschema.NgingAlertRecipient, extra echo.H, alertData *alert.AlertData) (err error) {
	title := alertData.Title
	ct := alertData.Content
	if ct == nil {
		return
	}
	params := alertData.Data
	switch a.Type {
	case `email`:
		content := ct.EmailContent(params)
		if len(content) == 0 {
			return
		}
		err = send.Mail(a.Account, strings.SplitN(a.Account, `@`, 2)[0], title, content)
	case `webhook`:
		mess := imbot.Open(a.Platform)
		if mess == nil || mess.Messager == nil {
			return
		}

		content := ct.MarkdownContent(params)
		if len(content) == 0 {
			return
		}
		message := string(content)

		var apiURL string
		if len(a.Account) > 7 {
			switch a.Account[0:7] {
			case `https:/`, `http://`:
				apiURL = a.Account
			}
		}
		if len(apiURL) == 0 {
			apiURL = mess.Messager.BuildURL(a.Account)
		}
		var atMobiles []string
		if extra.Has(`at`) {
			switch v := extra.Get(`at`).(type) {
			case []interface{}:
				atMobiles = make([]string, len(v))
				for k, m := range v {
					atMobiles[k] = param.AsString(m)
				}
			case []string:
				atMobiles = v
			}
		}
		go func(apiURL string, title string, message string, atMobiles ...string) {
			err = mess.Messager.SendMarkdown(apiURL, title, message, atMobiles...)
		}(apiURL, title, message, atMobiles...)
	}
	return
}
