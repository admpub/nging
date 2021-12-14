/*
   Nging is a toolbox for webmasters
   Copyright (C) 2018-present  Wenhui Shen <swh@admpub.com>

   This program is free software: you can redistribute it and/or modify
   it under the terms of the GNU Affero General Public License as published
   by the Free Software Foundation, either version 3 of the License, or
   (at your option) any later version.

   This program is distributed in the hope that it will be useful,
   but WITHOUT ANY WARRANTY; without even the implied warranty of
   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
   GNU Affero General Public License for more details.

   You should have received a copy of the GNU Affero General Public License
   along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/

package settings

import (
	"github.com/admpub/nging/v4/application/dbschema"
	"github.com/webx-top/echo"
)

func InsertDefaultConfig(ctx echo.Context, group, key string, values ...string) error {
	gs, ok := configDefaults[group]
	if !ok {
		return nil
	}
	return InsertBy(ctx, gs, key, values...)
}

func InsertBy(ctx echo.Context, configs map[string]*dbschema.NgingConfig, key string, values ...string) error {
	cfg, ok := configs[key]
	if !ok {
		return nil
	}
	cfgCopy := *cfg
	cfgCopy.SetContext(ctx)
	switch len(values) {
	case 2:
		cfgCopy.Disabled = values[1]
		fallthrough
	case 1:
		cfgCopy.Value = values[0]
	}
	if len(cfgCopy.Disabled) == 0 {
		cfgCopy.Disabled = `N`
	}
	if len(cfgCopy.Encrypted) == 0 {
		cfgCopy.Encrypted = `N`
	}
	_, err := cfgCopy.Add()
	return err
}

func InsertMissing(ctx echo.Context, gm *echo.Mapx, added map[string]int, configs map[string]*dbschema.NgingConfig, encoder Encoder) error {
	for key, cfg := range configs {
		_, ok := added[key]
		if ok {
			continue
		}
		cfgCopy := *cfg
		cfgCopy.SetContext(ctx)
		setting := gm.Get(key)
		if setting != nil {
			_v := setting.Get(`value`)
			if _v != nil {
				value, err := EncodeConfigValue(_v, &cfgCopy, encoder)
				if err != nil {
					return err
				}
				cfgCopy.Value = value
			}
			disabled := setting.Value(`disabled`)
			if len(disabled) > 0 {
				cfgCopy.Disabled = disabled
			}
		}
		if len(cfgCopy.Disabled) == 0 {
			cfgCopy.Disabled = `N`
		}
		if len(cfgCopy.Encrypted) == 0 {
			cfgCopy.Encrypted = `N`
		}
		_, err := cfgCopy.Add()
		if err != nil {
			return err
		}
	}
	return nil
}

func InsertMissingDefaultConfig(ctx echo.Context, added map[string]map[string]struct{}) error {
	for group, configs := range configDefaults {
		addedConfig, y := added[group]
		if !y { //整个组都没有的时候，添加整组
			for _, _cfg := range configs {
				cfg := *_cfg
				cfg.SetContext(ctx)
				if len(cfg.Disabled) == 0 {
					cfg.Disabled = `N`
				}
				if len(cfg.Encrypted) == 0 {
					cfg.Encrypted = `N`
				}
				_, err := cfg.Add()
				if err != nil {
					return err
				}
			}
			continue
		}
		for key, _cfg := range configs {
			if _, y := addedConfig[key]; y {
				continue
			}
			cfg := *_cfg
			cfg.SetContext(ctx)
			if len(cfg.Disabled) == 0 {
				cfg.Disabled = `N`
			}
			if len(cfg.Encrypted) == 0 {
				cfg.Encrypted = `N`
			}
			_, err := cfg.Add()
			if err != nil {
				return err
			}
		}
	}
	return nil
}
