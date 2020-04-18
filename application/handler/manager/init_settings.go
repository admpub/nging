package manager

import (
	"github.com/webx-top/com"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/param"

	"github.com/admpub/nging/application/dbschema"
	"github.com/admpub/nging/application/library/config"
	"github.com/admpub/nging/application/model"
	"github.com/admpub/nging/application/model/file/storer"
	"github.com/admpub/nging/application/registry/settings"
	"github.com/admpub/nging/application/registry/upload/driver"
	"github.com/admpub/nging/application/registry/upload/driver/local"
	"github.com/admpub/nging/application/registry/upload/driver/s3"

	"github.com/webx-top/image"
)

var configDefaults = map[string]map[string]*dbschema.NgingConfig{
	`base`: {
		`storer`: {
			Key:         `storer`,
			Label:       `存储引擎`,
			Description: ``,
			Value:       `{"name": "local", "id": ""}`,
			Group:       `base`,
			Type:        `json`,
			Sort:        30,
			Disabled:    `N`,
		},
		`watermark`: {
			Key:         `watermark`,
			Label:       `图片水印`,
			Description: ``,
			Value:       `{"watermark": "/public/assets/backend/images/nging-gear.png", "type": "image", "position": 0, "padding": 0, "on": true}`,
			Group:       `base`,
			Type:        `json`,
			Sort:        30,
			Disabled:    `N`,
		},
	},
}

var defaultStorer = storer.Info{
	Name: local.Name,
}

func init() {
	// 添加默认配置数据
	for group, configs := range configDefaults {
		settings.AddDefaultConfig(group, configs)
	}
	// 注册配置模板和逻辑
	if index, setting := settings.Get(`base`); index != -1 && setting != nil {
		setting.AddHookGet(func(ctx echo.Context) error {
			ctx.Set(`storerNames`, driver.AllNames())
			m := model.NewCloudStorage(ctx)
			m.ListByOffset(nil, nil, 0, -1)
			ctx.Set(`cloudStorageAccounts`, m.Objects())
			return nil
		})
	}
	settings.RegisterDecoder(`base.storer`, func(v *dbschema.NgingConfig, r echo.H) error {
		jsonData := storer.NewInfo()
		if len(v.Value) > 0 {
			com.JSONDecode(com.Str2bytes(v.Value), jsonData)
		}
		r[`ValueObject`] = jsonData
		return nil
	})
	settings.RegisterDecoder(`base.watermark`, func(v *dbschema.NgingConfig, r echo.H) error {
		jsonData := image.NewWatermarkOptions()
		if len(v.Value) > 0 {
			com.JSONDecode(com.Str2bytes(v.Value), jsonData)
		}
		r[`ValueObject`] = jsonData
		return nil
	})
	settings.RegisterEncoder(`base.storer`, func(v *dbschema.NgingConfig, r echo.H) ([]byte, error) {
		cfg := storer.NewInfo().FromStore(r)
		if cfg.Name == local.Name {
			cfg.ID = ``
		} else {
			id := param.AsUint(cfg.ID)
			if id > 0 {
				cfg.Name = s3.Name
			} else {
				cfg.Name = local.Name
			}
		}
		return com.JSONEncode(cfg)
	})
	settings.RegisterEncoder(`base.watermark`, func(v *dbschema.NgingConfig, r echo.H) ([]byte, error) {
		echo.Dump(r)
		cfg := image.NewWatermarkOptions().FromStore(r)
		return com.JSONEncode(cfg)
	})
	var updateStorer = func(cfg echo.H) error {
		settings, ok := cfg.Store(`base`).Get(`storer`).(*storer.Info)
		if !ok || settings == nil {
			settings = &defaultStorer
		}
		echo.Set(storer.StorerInfoKey, settings)
		return nil
	}
	config.OnInitSettings(updateStorer)
	config.OnSetSettings(func(group string, cfg echo.H) error {
		if group != `base` {
			return nil
		}
		return updateStorer(cfg)
	})
}
