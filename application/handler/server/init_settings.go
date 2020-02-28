package server

import (
	"github.com/webx-top/com"
	"github.com/webx-top/echo"

	"github.com/admpub/nging/application/dbschema"
	"github.com/admpub/nging/application/library/config"
	"github.com/admpub/nging/application/library/system"
	"github.com/admpub/nging/application/registry/settings"
)

var configDefaults = map[string]map[string]*dbschema.NgingConfig{
	`base`: {
		`systemStatus`: {
			Key:         `systemStatus`,
			Label:       `系统状态`,
			Description: ``,
			Value:       `{"MonitorOn":true,"AlarmOn":false,"AlarmThreshold":{"Memory":90,"CPU":10},"ReportEmail":""}`,
			Group:       `base`,
			Type:        `json`,
			Sort:        30,
			Disabled:    `N`,
		},
	},
}

func init() {
	// 添加默认配置数据
	for group, configs := range configDefaults {
		settings.AddDefaultConfig(group, configs)
	}
	// 注册配置模板和逻辑
	if index, setting := settings.Get(`base`); index != -1 && setting != nil {
		setting.AddTmpl(`server/settings/base`)
	}
	settings.RegisterDecoder(`base.systemStatus`, func(v *dbschema.NgingConfig, r echo.H) error {
		jsonData := system.NewSettings()
		if len(v.Value) > 0 {
			com.JSONDecode(com.Str2bytes(v.Value), jsonData)
		}
		r[`ValueObject`] = jsonData
		return nil
	})
	settings.RegisterEncoder(`base.systemStatus`, func(v *dbschema.NgingConfig, r echo.H) ([]byte, error) {
		cfg := system.NewSettings().FromStore(r)
		return com.JSONEncode(cfg)
	})
	config.OnInitSettings(func(cfg echo.H) error {
		settings, ok := cfg.Store(`base`).Get(`systemStatus`).(*system.Settings)
		if !ok || !settings.MonitorOn {
			return nil
		}
		system.ListenRealTimeStatus(settings)
		return nil
	})
	config.OnSetSettings(func(group string, cfg echo.H) error {
		if group != `base` {
			return nil
		}
		settings, ok := cfg.Store(`base`).Get(`systemStatus`).(*system.Settings)
		if !ok || !settings.MonitorOn {
			system.CancelRealTimeStatusCollection()
			return nil
		}
		system.ListenRealTimeStatus(settings)
		return nil
	})
}
