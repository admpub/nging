package module

import (
	"strings"

	"github.com/admpub/nging/v4/application/library/common"
	"github.com/admpub/nging/v4/application/library/config"
	"github.com/admpub/nging/v4/application/library/config/cmder"
	"github.com/admpub/nging/v4/application/library/config/extend"
	"github.com/admpub/nging/v4/application/library/ntemplate"
	"github.com/admpub/nging/v4/application/library/route"
	"github.com/admpub/nging/v4/application/registry/dashboard"
	"github.com/admpub/nging/v4/application/registry/navigate"
	"github.com/admpub/nging/v4/application/registry/settings"
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
	SetSettings()
	SetDefaultStartup()
	DBSchemaVersion() float64
}

var _ IModule = &Module{}

type Module struct {
	Startup       string
	Navigate      func(nc *navigate.Collection)
	Extend        map[string]extend.Initer
	Cmder         map[string]cmder.Cmder
	TemplatePath  map[string]string
	AssetsPath    []string
	SQLCollection func(sc *config.SQLCollection)
	Dashboard     func(dd *dashboard.Dashboards)
	Route         func(r *route.Collection)
	LogParser     map[string]common.LogParser
	Settings      []*settings.SettingForm
	DBSchemaVer   float64
}

func (m *Module) SetNavigate(nc *navigate.Collection) {
	if m.Navigate == nil {
		return
	}
	m.Navigate(nc)
}

func (m *Module) SetConfig(*config.Config) {
	if m.Extend == nil {
		return
	}
	for k, v := range m.Extend {
		extend.Register(k, v)
	}
}

func (m *Module) SetCmder(*config.CLIConfig) {
	if m.Cmder == nil {
		return
	}
	for k, v := range m.Cmder {
		cmder.Register(k, v)
	}
}

func (m *Module) SetTemplate(pa ntemplate.PathAliases) {
	if m.TemplatePath == nil {
		return
	}
	for k, v := range m.TemplatePath {
		pa.Add(k, v)
	}
}

func (m *Module) SetAssets(so *middleware.StaticOptions) {
	for _, v := range m.AssetsPath {
		so.AddFallback(v)
	}
}

func (m *Module) SetSQL(sc *config.SQLCollection) {
	if m.SQLCollection == nil {
		return
	}
	m.SQLCollection(sc)
}

func (m *Module) SetDashboard(dd *dashboard.Dashboards) {
	if m.Dashboard == nil {
		return
	}
	m.Dashboard(dd)
}

func (m *Module) SetRoute(r *route.Collection) {
	if m.Route == nil {
		return
	}
	m.Route(r)
}

func (m *Module) SetLogParser(parsers map[string]common.LogParser) {
	if m.LogParser == nil {
		return
	}
	for k, p := range m.LogParser {
		parsers[k] = p
	}
}

func (m *Module) SetSettings() {
	settings.Register(m.Settings...)
}

func (m *Module) SetDefaultStartup() {
	if len(m.Startup) > 0 {
		if len(config.DefaultStartup) > 0 && !strings.HasPrefix(m.Startup, `,`) {
			config.DefaultStartup += `,` + m.Startup
		} else {
			config.DefaultStartup += m.Startup
		}
	}
}

func (m *Module) DBSchemaVersion() float64 {
	return m.DBSchemaVer
}
