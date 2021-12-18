//go:build !bindata
// +build !bindata

package module

import (
	"github.com/admpub/nging/v4/application/library/bindata"
	"github.com/admpub/nging/v4/application/library/common"
	"github.com/admpub/nging/v4/application/library/config"
	"github.com/admpub/nging/v4/application/library/route"
	"github.com/admpub/nging/v4/application/registry/dashboard"
	"github.com/admpub/nging/v4/application/registry/navigate"
	"github.com/webx-top/echo"
)

func Register(modules ...IModule) {
	schemaVer := echo.Float64(`SCHEMA_VER`)
	versionNumbers := []float64{schemaVer}
	for _, module := range modules {
		module.SetNavigate(navigate.Default)
		module.SetConfig(config.DefaultConfig)
		module.SetCmder(config.DefaultCLIConfig)
		module.SetTemplate(bindata.PathAliases)
		module.SetAssets(bindata.StaticOptions)
		module.SetSQL(config.GetSQLCollection())
		module.SetDashboard(dashboard.Default)
		module.SetRoute(route.Default)
		module.SetLogParser(common.LogParsers)
		module.SetSettings()
		versionNumbers = append(versionNumbers, module.DBSchemaVersion())
	}
	echo.Set(`SCHEMA_VER`, common.Float64Sum(versionNumbers...))
}
