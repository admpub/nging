package task

import (
	"github.com/admpub/nging/v4/application/library/config/cmder"
	cronCmder "github.com/admpub/nging/v4/application/library/cron/cmder"
	"github.com/admpub/nging/v4/application/library/module"
	"github.com/admpub/nging/v4/application/registry/navigate"
)

const ID = `task`

var Module = module.Module{
	Startup: ID,
	Cmder: map[string]cmder.Cmder{
		ID: cronCmder.New(),
	},
	Navigate: func(nc *navigate.Collection) {
		nc.Backend.AddLeftItems(-1, LeftNavigate)
	},
}
