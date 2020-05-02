package model

import (
	"strings"
	"encoding/json"

	"github.com/webx-top/echo"
	"github.com/webx-top/echo/param"
	"github.com/webx-top/com"
	"github.com/admpub/nging/application/dbschema"
	"github.com/admpub/nging/application/library/imbot"
	"github.com/admpub/nging/application/library/cron"
)

type AlertRecipientExt struct {
	*dbschema.NgingAlertRecipient
	Extra echo.H
}

func SendAlert(ctx echo.Context, message string, extra ...string) error {
	m := NewAlertRecipient(ctx)
	var title string
	if len(extra)>0 {
		title = extra[0]
	}
	return m.Send(title,message)
}

func (a *AlertRecipientExt) Parse() *AlertRecipientExt {
	if a.Extra != nil {
		return a
	}
	a.Extra = echo.H{}
	if len(a.NgingAlertRecipient.Extra) > 0 {
		json.Unmarshal(com.Str2bytes(a.NgingAlertRecipient.Extra), &a.Extra)
	}
	return a
}

func (a *AlertRecipientExt) Send(title string, message string) (err error) {
	a.Parse()
	switch a.Type {
	case `email`:
		err = cron.SendMail(a.Account, strings.SplitN(a.Account, `@`, 2)[0], title, com.Str2bytes(message))
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
		if a.Extra.Has(`at`) {
			switch v := a.Extra.Get(`at`).(type) {
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