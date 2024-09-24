package task

import (
	"github.com/coscms/webcore/library/config/cmder"
	cronCmder "github.com/coscms/webcore/library/cron/cmder"
	"github.com/coscms/webcore/library/module"
)

const ID = `task`

var Module = module.Module{
	Startup: ID,
	Cmder: map[string]cmder.Cmder{
		ID: cronCmder.New(),
	},
	Navigate: func(nc module.Navigate) {
		nc.Backend().AddLeftItems(-1, LeftNavigate)
	},
}
