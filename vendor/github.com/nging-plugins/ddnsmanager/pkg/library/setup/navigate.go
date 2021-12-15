package setup

import (
	"github.com/admpub/nging/v4/application/handler/manager"
	"github.com/admpub/nging/v4/application/registry/navigate"
)

func init() {
	manager.TopNavigate.AddChild(`tool`, -1, &navigate.Item{
		Display: true,
		Name:    `DDNS`,
		Action:  `ddns`,
	})
}
