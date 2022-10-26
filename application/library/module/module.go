package module

import (
	"strings"

	"github.com/admpub/nging/v5/application/library/common"
	"github.com/admpub/nging/v5/application/library/config"
	"github.com/admpub/nging/v5/application/library/config/cmder"
	"github.com/admpub/nging/v5/application/library/config/extend"
	"github.com/admpub/nging/v5/application/library/cron"
	"github.com/admpub/nging/v5/application/library/ntemplate"
	"github.com/admpub/nging/v5/application/library/route"
	"github.com/admpub/nging/v5/application/registry/dashboard"
	"github.com/admpub/nging/v5/application/registry/navigate"
	"github.com/admpub/nging/v5/application/registry/settings"
	"github.com/webx-top/echo/middleware"
)

type IModule interface {
	Apply()
	Version() float64
}

var _ IModule = &Module{}

type Module struct {
	Startup       string                         // 默认启动项(多个用半角逗号“,”隔开)
	Navigate      func(nc *navigate.Collection)  // 注册导航菜单
	Extend        map[string]extend.Initer       // 注册扩展配置项
	Cmder         map[string]cmder.Cmder         // 注册命令
	TemplatePath  map[string]string              // 注册模板路径
	AssetsPath    []string                       // 注册素材路径
	SQLCollection func(sc *config.SQLCollection) // 注册SQL语句
	Dashboard     func(dd *dashboard.Dashboards) // 注册控制面板首页区块
	Route         func(r *route.Collection)      // 注册网址路由
	LogParser     map[string]common.LogParser    // 注册日志解析器
	Settings      []*settings.SettingForm        // 注册配置选项
	CronJobs      []*cron.Jobx                   // 注册定时任务
	DBSchemaVer   float64                        // 设置数据库结构版本号
}

func (m *Module) setNavigate(nc *navigate.Collection) {
	if m.Navigate == nil {
		return
	}
	m.Navigate(nc)
}

func (m *Module) setConfig(*config.Config) {
	if m.Extend == nil {
		return
	}
	for k, v := range m.Extend {
		extend.Register(k, v)
	}
}

func (m *Module) setCmder(*config.CLIConfig) {
	if m.Cmder == nil {
		return
	}
	for k, v := range m.Cmder {
		cmder.Register(k, v)
	}
}

func (m *Module) setTemplate(pa *ntemplate.PathAliases) {
	if m.TemplatePath == nil {
		return
	}
	for k, v := range m.TemplatePath {
		SetTemplate(pa, k, v)
	}
}

func (m *Module) setAssets(so *middleware.StaticOptions) {
	for _, v := range m.AssetsPath {
		SetAssets(so, v)
	}
}

func (m *Module) setSQL(sc *config.SQLCollection) {
	if m.SQLCollection == nil {
		return
	}
	m.SQLCollection(sc)
}

func (m *Module) setDashboard(dd *dashboard.Dashboards) {
	if m.Dashboard == nil {
		return
	}
	m.Dashboard(dd)
}

func (m *Module) setRoute(r *route.Collection) {
	if m.Route == nil {
		return
	}
	m.Route(r)
}

func (m *Module) setLogParser(parsers map[string]common.LogParser) {
	if m.LogParser == nil {
		return
	}
	for k, p := range m.LogParser {
		parsers[k] = p
	}
}

func (m *Module) setSettings() {
	settings.Register(m.Settings...)
}

func (m *Module) setCronJob() {
	for _, jobx := range m.CronJobs {
		jobx.Register()
	}
}

func (m *Module) setDefaultStartup() {
	if len(m.Startup) > 0 {
		if len(config.DefaultStartup) > 0 && !strings.HasPrefix(m.Startup, `,`) {
			config.DefaultStartup += `,` + m.Startup
		} else {
			config.DefaultStartup += m.Startup
		}
	}
}

func (m *Module) Version() float64 {
	return m.DBSchemaVer
}

func (m *Module) Apply() {
	m.setNavigate(navigate.Default)
	m.setConfig(config.FromFile())
	m.setCmder(config.FromCLI())
	m.applyTemplateAndAssets()
	//m.setTemplate(bindata.PathAliases)
	//m.setAssets(bindata.StaticOptions)
	m.setSQL(config.GetSQLCollection())
	m.setDashboard(dashboard.Default)
	m.setRoute(route.Default)
	m.setLogParser(common.LogParsers)
	m.setSettings()
	m.setDefaultStartup()
	m.setCronJob()
}

func SetTemplate(pa *ntemplate.PathAliases, key string, templatePath string) {
	if len(templatePath) == 0 {
		return
	}
	if templatePath[0] != '.' && templatePath[0] != '/' && !strings.HasPrefix(templatePath, `vendor/`) {
		templatePath = NgingPluginDir + `/` + templatePath
	}
	pa.Add(key, templatePath)
}

func SetAssets(so *middleware.StaticOptions, assetsPath string) {
	if len(assetsPath) == 0 {
		return
	}
	if assetsPath[0] != '.' && assetsPath[0] != '/' && !strings.HasPrefix(assetsPath, `vendor/`) {
		assetsPath = NgingPluginDir + `/` + assetsPath
	}
	so.AddFallback(assetsPath)
}
