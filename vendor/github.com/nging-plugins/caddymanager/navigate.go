package caddymanager

import (
	"github.com/admpub/nging/v4/application/registry/navigate"

	"github.com/nging-plugins/caddymanager/pkg/handler"
	_ "github.com/nging-plugins/caddymanager/pkg/library/cmder"
	_ "github.com/nging-plugins/caddymanager/pkg/library/setup"
)

var LeftNavigate = handler.LeftNavigate

func RegisterNavigate(nc *navigate.Collection) {
	nc.Backend.AddLeftItems(-1, LeftNavigate)
}
