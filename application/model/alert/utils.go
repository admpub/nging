package alert

import (
	"strings"

	"github.com/admpub/nging/application/dbschema"
	"github.com/admpub/nging/application/library/cron/send"
	"github.com/admpub/nging/application/library/imbot"
	"github.com/webx-top/com"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/param"
)

// Send 发送警报
func Send(a *dbschema.NgingAlertRecipient, extra echo.H, title string, message string) (err error) {
	switch a.Type {
	case `email`:
		err = send.Mail(a.Account, strings.SplitN(a.Account, `@`, 2)[0], title, com.Str2bytes(message))
	case `webhook`:
		mess := imbot.Open(a.Platform)
		if mess == nil || mess.Messager == nil {
			return
		}
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
