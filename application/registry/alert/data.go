package alert

import (
	"github.com/admpub/nging/v5/application/dbschema"
	"github.com/webx-top/echo/param"
)

func NewData(title string, ct ContentType) *AlertData {
	return &AlertData{
		Title:   title,
		Content: ct,
		Data:    param.Store{},
	}
}

type AlertData struct {
	Title   string
	Content ContentType
	Data    param.Store
}

type AlertTopicExt struct {
	*dbschema.NgingAlertTopic
	Recipient *dbschema.NgingAlertRecipient `db:"-,relation=id:recipient_id"`
}
