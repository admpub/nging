package cloudflare

import (
	"github.com/admpub/nging/v3/application/library/ddnsmanager"
	"github.com/admpub/nging/v3/application/library/ddnsmanager/interfaces"
)

func New() interfaces.Updater {
	return &Cloudflare{}
}

func init() {
	ddnsmanager.Register(Name, New)
}
