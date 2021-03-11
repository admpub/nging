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
	Extra     echo.H
}

func init() {
	alertRegistry.SendTopic = func(ctx echo.Context, topic string, params param.Store) error {
		return AlertSend(ctx, topic, params)
	}
}

func AlertSend(ctx echo.Context, topic string, params param.Store) error {
	m := NewAlertTopic(ctx)
	return m.Send(topic, params)
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

func (a *AlertTopicExt) Send(params param.Store) (err error) {
	if a.Recipient == nil || a.Recipient.Disabled == `Y` {
		return
	}
	a.Parse()
	return alertSend(a.Recipient, a.Extra, params)
}

func alertSend(a *dbschema.NgingAlertRecipient, extra echo.H, params param.Store) (err error) {
	return alert.Send(a, extra, params)
}
