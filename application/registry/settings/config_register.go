package settings

import (
	"github.com/admpub/nging/v4/application/dbschema"
)

func AddConfigs(configs map[string]map[string]*dbschema.NgingConfig) {
	for group, configs := range configs {
		AddDefaultConfig(group, configs)
	}
}

func AddTmpl(group string, tmpl string, opts ...FormSetter) {
	// 注册配置模板和逻辑
	index, setting := Get(group)
	if index == -1 || setting == nil {
		return
	}
	if len(tmpl) > 0 {
		setting.AddTmpl(tmpl)
	}
	for _, option := range opts {
		option(setting)
	}
}
