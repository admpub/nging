package model

import (
	"github.com/admpub/nging/v5/application/dbschema"
	"github.com/admpub/nging/v5/application/registry/alert"
)

type AlertRecipientExt struct {
	*dbschema.NgingAlertRecipient
}

func (a *AlertRecipientExt) Send(alertData *alert.AlertData) (err error) {
	return alertSend(a.NgingAlertRecipient, alertData)
}
