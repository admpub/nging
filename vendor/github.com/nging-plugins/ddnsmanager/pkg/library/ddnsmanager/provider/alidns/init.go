package alidns

import (
	"github.com/nging-plugins/ddnsmanager/pkg/library/ddnsmanager"
	"github.com/nging-plugins/ddnsmanager/pkg/library/ddnsmanager/interfaces"
)

func New() interfaces.Updater {
	return &Alidns{}
}

func init() {
	ddnsmanager.Register(Name, New)
}
