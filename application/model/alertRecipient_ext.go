package model

import (
	"encoding/json"

	"github.com/admpub/nging/v3/application/dbschema"
	"github.com/admpub/nging/v3/application/registry/alert"
	"github.com/webx-top/com"
	"github.com/webx-top/echo"
)

type AlertRecipientExt struct {
	*dbschema.NgingAlertRecipient
	Extra echo.H
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

func (a *AlertRecipientExt) Send(alertData *alert.AlertData) (err error) {
	a.Parse()
	return alertSend(a.NgingAlertRecipient, a.Extra, alertData)
}
