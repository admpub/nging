package servermanager

import (
	"github.com/admpub/nging/v4/application/registry/navigate"

	"github.com/nging-plugins/servermanager/application/handler"
)

var LeftNavigate = handler.LeftNavigate

func RegisterNavigate(nc *navigate.Collection) {
	nc.Backend.AddLeftItems(-1, LeftNavigate)
}
