package ddnsmanager

import (
	"github.com/admpub/nging/v4/application/registry/navigate"

	"github.com/nging-plugins/ddnsmanager/pkg/handler"
)

var TopNavigate = handler.TopNavigate

func RegisterNavigate(nc *navigate.Collection) {
	nc.Backend.GetTop().AddChild(`tool`, -1, TopNavigate)
}
