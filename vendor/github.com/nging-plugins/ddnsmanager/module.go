package ddnsmanager

import (
	"github.com/admpub/nging/v4/application/library/module"

	"github.com/nging-plugins/ddnsmanager/pkg/handler"
)

const ID = `ddns`

var Module = module.Module{
	TemplatePath: map[string]string{
		ID: `ddnsmanager/template/backend`,
	},
	AssetsPath:  []string{},
	Navigate:    RegisterNavigate,
	Route:       handler.RegisterRoute,
	DBSchemaVer: 0.0000,
}
