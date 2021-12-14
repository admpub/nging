package model

import (
	"github.com/admpub/nging/v4/application/dbschema"
	"github.com/admpub/nging/v4/application/registry/alert"
)

type AlertRecipientExt struct {
	*dbschema.NgingAlertRecipient
}

func (a *AlertRecipientExt) Send(alertData *alert.AlertData) (err error) {
	return alertSend(a.NgingAlertRecipient, alertData)
}
