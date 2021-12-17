package module

import (
	"github.com/admpub/nging/v4/application/library/common"
	"github.com/admpub/nging/v4/application/library/config"
	"github.com/admpub/nging/v4/application/library/ntemplate"
	"github.com/admpub/nging/v4/application/library/route"
	"github.com/admpub/nging/v4/application/registry/dashboard"
	"github.com/admpub/nging/v4/application/registry/navigate"
	"github.com/webx-top/echo/middleware"
)

type IModule interface {
	SetNavigate(*navigate.Collection)
	SetConfig(*config.Config)
	SetCmder(*config.CLIConfig)
	SetTemplate(ntemplate.PathAliases)
	SetAssets(*middleware.StaticOptions)
	SetSQL(*config.SQLCollection)
	SetDashboard(*dashboard.Dashboards)
	SetRoute(*route.Collection)
	SetLogParser(map[string]common.LogParser)
}

var _ IModule = &Module{}

type Module struct{}

func (m *Module) SetNavigate(*navigate.Collection)         {}
func (m *Module) SetConfig(*config.Config)                 {}
func (m *Module) SetCmder(*config.CLIConfig)               {}
func (m *Module) SetTemplate(ntemplate.PathAliases)        {}
func (m *Module) SetAssets(*middleware.StaticOptions)      {}
func (m *Module) SetSQL(*config.SQLCollection)             {}
func (m *Module) SetDashboard(*dashboard.Dashboards)       {}
func (m *Module) SetRoute(*route.Collection)               {}
func (m *Module) SetLogParser(map[string]common.LogParser) {}
