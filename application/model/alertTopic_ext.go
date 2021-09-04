package model

import (
	"encoding/json"

	"github.com/admpub/nging/v3/application/dbschema"
	"github.com/admpub/nging/v3/application/model/alert"
	alertRegistry "github.com/admpub/nging/v3/application/registry/alert"
	"github.com/webx-top/com"
	"github.com/webx-top/echo"
)

type AlertTopicExt struct {
	*dbschema.NgingAlertTopic
	Recipient *dbschema.NgingAlertRecipient `db:"-,relation=id:recipient_id"`
	Extra     echo.H
}

func init() {
	alertRegistry.SendTopic = func(ctx echo.Context, topic string, alertData *alertRegistry.AlertData) error {
		return AlertSend(ctx, topic, alertData)
	}
}

func AlertSend(ctx echo.Context, topic string, alertData *alertRegistry.AlertData) error {
	m := NewAlertTopic(ctx)
	return m.Send(topic, alertData)
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

func (a *AlertTopicExt) Send(alertData *alertRegistry.AlertData) (err error) {
	if a.Recipient == nil || a.Recipient.Disabled == `Y` {
		return
	}
	a.Parse()
	return alertSend(a.Recipient, a.Extra, alertData)
}

func alertSend(a *dbschema.NgingAlertRecipient, extra echo.H, alertData *alertRegistry.AlertData) (err error) {
	return alert.Send(a, extra, alertData)
}
