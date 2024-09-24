package cloud

import (
	"github.com/coscms/webcore/library/module"
)

const ID = `cloud`

var Module = module.Module{
	Navigate: func(nc module.Navigate) {
		nc.Backend().AddLeftItems(-1, LeftNavigate)
	},
}
