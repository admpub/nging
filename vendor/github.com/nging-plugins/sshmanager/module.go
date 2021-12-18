package sshmanager

import (
	"github.com/admpub/nging/v4/application/library/module"

	"github.com/nging-plugins/sshmanager/pkg/handler"
	"github.com/nging-plugins/sshmanager/pkg/library/setup"
)

const ID = `term`

var Module = module.Module{
	TemplatePath: map[string]string{
		ID: `sshmanager/template/backend`,
	},
	AssetsPath:    []string{},
	SQLCollection: setup.RegisterSQL,
	Navigate:      RegisterNavigate,
	Route:         handler.RegisterRoute,
	DBSchemaVer:   0.0000,
}
