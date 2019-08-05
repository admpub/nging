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

package config

import (
	"github.com/admpub/log"
	"github.com/admpub/nging/application/library/notice"
	"github.com/admpub/nging/application/registry/settings"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/subdomains"
)

func InDB(group ...string) echo.H {
	sz := len(group)
	if sz > 0 {
		cfg := DefaultConfig.ConfigInDB.GetConfig().Store(group[0])
		if sz == 1 {
			return cfg
		}
		for _, key := range group[1:] {
			cfg = cfg.Store(key)
		}
		return cfg
	}
	return DefaultConfig.ConfigInDB.GetConfig()
}

func NewConfigInDB(config *Config) *ConfigInDB {
	c := &ConfigInDB{
		config: config,
	}
	return c
}

type ConfigInDB struct {
	Email  Email  `json:"email"`
	Log    Log    `json:"log"`
	APIKey string `json:"-"` //API密钥
	Debug  bool   `json:"debug"`
	Base   echo.H `json:"base"`
	config *Config
}

func (c *ConfigInDB) SetBy(r echo.H, defaults echo.H) *ConfigInDB {
	if !r.Has(`base`) && defaults != nil {
		r.Set(`base`, defaults.Store(`base`))
	}
	c.Base = r.Store(`base`)
	c.APIKey = c.Base.String(`apiKey`)
	c.Debug = c.Base.Bool(`debug`)
	return c
}

func (c *ConfigInDB) SetDebug(on bool) {
	c.Log.Debug = on
	c.Debug = on

	c.config.DB.SetDebug(on)
	notice.SetDebug(on)
	if on {
		log.Info(`Currently running in debug mode`)
		log.SetLevel(`Debug`)
	} else {
		log.Info(`Currently running in normal mode`)
		log.SetLevel(`Info`)
	}
	subdomains.Default.SetDebug(on)
}

func (c *ConfigInDB) Init() {
	defaults := settings.ConfigDefaultsAsStore()
	var configs = defaults
	if IsInstalled() {
		if c.config.ConnectedDB() {
			configs = settings.ConfigAsStore()
		}
	}
	echo.Set(`NgingConfig`, configs)
	for _, group := range []string{`base`, `email`, `log`} {
		c.SetConfig(group, configs, defaults)
	}
}

func (c *ConfigInDB) GetConfig() echo.H {
	r, _ := echo.Get(`NgingConfig`).(echo.H)
	return r
}

func (c *ConfigInDB) SetConfigs(groups ...string) {
	ngingConfig := c.GetConfig()
	configs := settings.ConfigAsStore(groups...)
	for _, group := range []string{`base`, `email`, `log`} {
		conf, ok := configs[group]
		if !ok {
			continue
		}
		ngingConfig.Set(group, conf)
		c.SetConfig(group, ngingConfig, nil)
	}
	echo.Set(`NgingConfig`, ngingConfig)
}

func (c *ConfigInDB) SetConfig(group string, ngingConfig echo.H, defaults echo.H) {
	switch group {
	case `base`:
		c.SetBy(ngingConfig, defaults)
		c.SetDebug(c.Debug)
	case `email`:
		c.Email.SetBy(ngingConfig, defaults).Init()
	case `log`:
		c.Log.SetBy(ngingConfig, defaults).Init()
	}
}
