package dnspod

import (
	"github.com/admpub/nging/v3/application/library/ddnsmanager"
	"github.com/admpub/nging/v3/application/library/ddnsmanager/interfaces"
)

func New() interfaces.Updater {
	return &Dnspod{}
}

func init() {
	ddnsmanager.Register(`dnspod`, New)
}
