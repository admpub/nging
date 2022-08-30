package frpmanager

import (
	"github.com/admpub/nging/v4/application/library/config/cmder"
	"github.com/admpub/nging/v4/application/library/module"

	"github.com/nging-plugins/frpmanager/application/handler"
	pluginCmder "github.com/nging-plugins/frpmanager/application/library/cmder"
	"github.com/nging-plugins/frpmanager/application/library/setup"
)

const ID = `frp`

var Module = module.Module{
	Startup: `frpserver,frpclient`,
	Cmder: map[string]cmder.Cmder{
		`frpclient`: pluginCmder.NewClient(),
		`frpserver`: pluginCmder.NewServer(),
	},
	TemplatePath: map[string]string{
		ID: `frpmanager/template/backend`,
	},
	AssetsPath:    []string{},
	SQLCollection: setup.RegisterSQL,
	Navigate:      RegisterNavigate,
	Route:         handler.RegisterRoute,
	DBSchemaVer:   0.0000,
}
