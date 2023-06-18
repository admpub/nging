package firewallmanager

import (
	"github.com/admpub/nging/v5/application/library/config/cmder"
	"github.com/admpub/nging/v5/application/library/config/extend"
	"github.com/admpub/nging/v5/application/library/module"

	"github.com/nging-plugins/firewallmanager/application/handler"
	pluginCmder "github.com/nging-plugins/firewallmanager/application/library/cmder"
	"github.com/nging-plugins/firewallmanager/application/library/setup"
)

const ID = `firewall`

var Module = module.Module{
	Startup: `firewall`,
	Extend: map[string]extend.Initer{
		`firewall`: pluginCmder.Initer,
	},
	Cmder: map[string]cmder.Cmder{
		`firewall`: pluginCmder.New(),
	},
	TemplatePath: map[string]string{
		ID: `firewallmanager/template/backend`,
	},
	AssetsPath:    []string{},
	SQLCollection: setup.RegisterSQL,
	Navigate:      RegisterNavigate,
	Route:         handler.RegisterRoute,
	DBSchemaVer:   0.0001,
}
