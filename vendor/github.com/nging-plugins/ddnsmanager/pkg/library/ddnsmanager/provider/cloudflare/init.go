package cloudflare

import (
	"github.com/nging-plugins/ddnsmanager/pkg/library/ddnsmanager"
	"github.com/nging-plugins/ddnsmanager/pkg/library/ddnsmanager/interfaces"
)

func New() interfaces.Updater {
	return &Cloudflare{}
}

func init() {
	ddnsmanager.Register(Name, New)
}
