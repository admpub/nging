package cloud

import (
	"github.com/admpub/nging/v5/application/library/module"
	"github.com/admpub/nging/v5/application/registry/navigate"
)

const ID = `cloud`

var Module = module.Module{
	Navigate: func(nc *navigate.Collection) {
		nc.Backend.AddLeftItems(-1, LeftNavigate)
	},
}
