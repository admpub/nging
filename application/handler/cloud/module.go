package cloud

import (
	"github.com/coscms/webcore/library/module"
	"github.com/coscms/webcore/registry/navigate"
)

const ID = `cloud`

var Module = module.Module{
	Navigate: func(nc *navigate.Collection) {
		nc.Backend.AddLeftItems(-1, LeftNavigate)
	},
}
