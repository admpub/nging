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
)

func (m *Module) Apply() {
	m.setNavigate(navigate.Default)
	m.setConfig(config.DefaultConfig)
	m.setCmder(config.DefaultCLIConfig)
	m.setTemplate(bindata.PathAliases)
	m.setAssets(bindata.StaticOptions)
	m.setSQL(config.GetSQLCollection())
	m.setDashboard(dashboard.Default)
	m.setRoute(route.Default)
	m.setLogParser(common.LogParsers)
	m.setSettings()
	m.setDefaultStartup()
	m.setCronJob()
}
