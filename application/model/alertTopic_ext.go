package model

import (
	"github.com/admpub/nging/v5/application/dbschema"
	"github.com/admpub/nging/v5/application/model/alert"
	alertRegistry "github.com/admpub/nging/v5/application/registry/alert"
	"github.com/webx-top/echo"
)

type AlertTopicExt struct {
	*dbschema.NgingAlertTopic
	Recipient *dbschema.NgingAlertRecipient `db:"-,relation=id:recipient_id"`
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

func (a *AlertTopicExt) Send(alertData *alertRegistry.AlertData) (err error) {
	if a.Recipient == nil || a.Recipient.Disabled == `Y` {
		return
	}
	return alertSend(a.Recipient, alertData)
}

func alertSend(a *dbschema.NgingAlertRecipient, alertData *alertRegistry.AlertData) (err error) {
	return alert.Send(a, alertData)
}
