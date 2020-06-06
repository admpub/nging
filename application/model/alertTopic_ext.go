package model

import (
	"encoding/json"
	"strings"

	"github.com/admpub/nging/application/dbschema"
	"github.com/admpub/nging/application/library/cron"
	"github.com/admpub/nging/application/library/imbot"
	"github.com/webx-top/com"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/param"
)

type AlertTopicExt struct {
	*dbschema.NgingAlertTopic
	Recipient *dbschema.NgingAlertRecipient `db:"-,relation=id:recipient_id"`
	Extra echo.H
}

func AlertSend(ctx echo.Context, topic string, message string, extra ...string) error {
	m := NewAlertTopic(ctx)
	var title string
	if len(extra) > 0 {
		title = extra[0]
	}
	return m.Send(topic, title, message)
}

func (a *AlertTopicExt) Parse() *AlertTopicExt {
	if a.Extra != nil {
		return a
	}
	a.Extra = echo.H{}
	if a.Recipient == nil {
		return a
	}
	if len(a.Recipient.Extra) > 0 {
		json.Unmarshal(com.Str2bytes(a.Recipient.Extra), &a.Extra)
	}
	return a
}

func (a *AlertTopicExt) Send(title string, message string) (err error) {
	if a.Recipient == nil || a.Recipient.Disabled == `Y` {
		return
	}
	a.Parse()
	return alertSend(a.Recipient, a.Extra, title, message)
}

func alertSend(a *dbschema.NgingAlertRecipient, extra echo.H, title string, message string) (err error) {
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
