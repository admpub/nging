package module

import (
	"github.com/webx-top/echo"

	"github.com/admpub/nging/v5/application/library/common"
)

var NgingPluginDir = `vendor/github.com/nging-plugins` // `../../nging-plugins`

func Register(modules ...IModule) {
	schemaVer := echo.Float64(`SCHEMA_VER`)
	versionNumbers := []float64{schemaVer}
	for _, module := range modules {
		module.Apply()
		versionNumbers = append(versionNumbers, module.Version())
	}
	echo.Set(`SCHEMA_VER`, common.Float64Sum(versionNumbers...))
}
