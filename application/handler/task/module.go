package task

import (
	"github.com/admpub/nging/v4/application/library/module"
	"github.com/admpub/nging/v4/application/registry/navigate"
)

const ID = `task`

var Module = module.Module{
	Navigate: func(nc *navigate.Collection) {
		nc.Backend.AddLeftItems(-1, LeftNavigate)
	},
}
