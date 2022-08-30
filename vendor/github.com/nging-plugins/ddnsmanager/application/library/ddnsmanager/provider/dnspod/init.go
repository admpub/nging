package dnspod

import (
	"github.com/nging-plugins/ddnsmanager/application/library/ddnsmanager"
	"github.com/nging-plugins/ddnsmanager/application/library/ddnsmanager/interfaces"
)

func New() interfaces.Updater {
	return &Dnspod{}
}

func init() {
	ddnsmanager.Register(Name, New)
}
