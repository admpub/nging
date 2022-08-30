package namecheap

import (
	"github.com/nging-plugins/ddnsmanager/application/library/ddnsmanager"
	"github.com/nging-plugins/ddnsmanager/application/library/ddnsmanager/interfaces"
)

func New() interfaces.Updater {
	return &Namecheap{}
}

func init() {
	ddnsmanager.Register(Name, New)
}
