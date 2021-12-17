//go:build !bindata
// +build !bindata

package module

import (
	"github.com/admpub/nging/v4/application/library/bindata"
	"github.com/admpub/nging/v4/application/library/common"
	"github.com/admpub/nging/v4/application/library/config"
	"github.com/admpub/nging/v4/application/registry/dashboard"
	"github.com/admpub/nging/v4/application/registry/navigate"
	"github.com/admpub/nging/v4/application/registry/route"
)

func Register(modules ...IModule) {
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
	}
}
