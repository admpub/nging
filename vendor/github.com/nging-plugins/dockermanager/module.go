package dockermanager

import (
	"github.com/admpub/nging/v5/application/library/config/cmder"
	"github.com/admpub/nging/v5/application/library/module"

	"github.com/nging-plugins/dockermanager/application/handler"
	"github.com/nging-plugins/dockermanager/application/library/setup"
)

const ID = `docker`

var Module = module.Module{
	Startup: ``,
	Cmder:   map[string]cmder.Cmder{},
	TemplatePath: map[string]string{
		ID: `dockermanager/template/backend`,
	},
	AssetsPath: []string{
		`dockermanager/public/assets`,
	},
	SQLCollection: setup.RegisterSQL,
	Navigate:      RegisterNavigate,
	Route:         handler.RegisterRoute,
	DBSchemaVer:   0.0000,
}
