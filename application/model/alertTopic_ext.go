package model

import (
	"encoding/json"

	"github.com/admpub/nging/application/dbschema"
	"github.com/admpub/nging/application/model/alert"
	alertRegistry "github.com/admpub/nging/application/registry/alert"
	"github.com/webx-top/com"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/param"
)

type AlertTopicExt struct {
	*dbschema.NgingAlertTopic
	Recipient *dbschema.NgingAlertRecipient `db:"-,relation=id:recipient_id"`
	Extra echo.H
}

func init() {
	alertRegistry.SendTopic = func(ctx echo.Context, topic string, data param.Store) error {
		return AlertSend(ctx, topic, data.String(`email-content`), data.String(`title`))
	}
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
	return alert.Send(a, extra, title, message)
}
