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
	"strings"

	"github.com/webx-top/echo"
	"github.com/webx-top/echo/subdomains"

	"github.com/admpub/log"
	"github.com/admpub/nging/application/library/notice"
	"github.com/admpub/nging/application/registry/settings"
)

func Setting(group ...string) echo.H {
	st, ok := echo.Get(`NgingConfig`).(echo.H)
	if !ok {
		if st == nil {
			st = echo.H{}
		}
		return st
	}
	sz := len(group)
	if sz <= 0 {
		return st
	}
	cfg := st.Store(group[0])
	if sz == 1 {
		return cfg
	}
	for _, key := range group[1:] {
		cfg = cfg.Store(key)
	}
	return cfg
}

func NewSettings(config *Config) *Settings {
	c := &Settings{
		config: config,
	}
	return c
}

type Settings struct {
	Email              Email  `json:"email"`
	Log                Log    `json:"log"`
	APIKey             string `json:"-"` //API密钥
	Debug              bool   `json:"debug"`
	MaxRequestBodySize int    `json:"maxRequestBodySize"`
	Base               echo.H `json:"base"`
	config             *Config
}

func (c *Settings) SetBy(r echo.H, defaults echo.H) *Settings {
	if !r.Has(`base`) && defaults != nil {
		r.Set(`base`, defaults.Store(`base`))
	}
	c.Base = r.Store(`base`)
	c.APIKey = c.Base.String(`apiKey`)
	c.Debug = c.Base.Bool(`debug`)
	c.MaxRequestBodySize = c.Base.Int(`maxRequestBodySize`)
	return c
}

func (c *Settings) SetDebug(on bool) {
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

var (
	actGroups     = []string{`base`, `smtp`, `log`}
	onSetSettings = map[string][]func(globalCfg echo.H) error{}
)

func OnSetSettings(groupAndKey string, fn func(echo.H) error) {
	if _, ok := onSetSettings[groupAndKey]; !ok {
		onSetSettings[groupAndKey] = []func(globalCfg echo.H) error{}
	}
	onSetSettings[groupAndKey] = append(onSetSettings[groupAndKey], fn)
}

func FireInitSettings(cfg echo.H) error {
	for _, fnList := range onSetSettings {
		for _, fn := range fnList {
			if err := fn(cfg); err != nil {
				return err
			}
		}
	}
	return nil
}

func FireSetSettings(group string, globalCfg echo.H) error {
	for groupAndKey, fnList := range onSetSettings {
		if !strings.HasPrefix(group+`.`, groupAndKey) {
			continue
		}
		for _, fn := range fnList {
			if err := fn(globalCfg); err != nil {
				return err
			}
		}
	}
	return nil
}

func (c *Settings) Init() {
	defaults := settings.ConfigDefaultsAsStore()
	var configs = defaults
	if IsInstalled() {
		if c.config.ConnectedDB() {
			configs = settings.ConfigAsStore()
		}
	}
	echo.Set(`NgingConfig`, configs)
	for _, group := range actGroups {
		c.SetConfig(group, configs, defaults)
	}
	err := FireInitSettings(configs)
	if err != nil {
		log.Error(err)
	}
}

func (c *Settings) GetConfig() echo.H {
	r, _ := echo.Get(`NgingConfig`).(echo.H)
	return r
}

func (c *Settings) SetConfigs(groups ...string) {
	ngingConfig := c.GetConfig()
	configs := settings.ConfigAsStore(groups...)
	for group, conf := range configs {
		ngingConfig.Set(group, conf)
		FireSetSettings(group, ngingConfig)
		//log.Debug(`Change configuration:`, group, `:`, echo.Dump(conf, false))
		c.SetConfig(group, ngingConfig, nil)
	}
}

func (c *Settings) SetConfig(group string, ngingConfig echo.H, defaults echo.H) {
	switch group {
	case `base`:
		c.SetBy(ngingConfig, defaults)
		c.SetDebug(c.Debug)
	case `smtp`:
		c.Email.SetBy(ngingConfig, defaults).Init()
	case `log`:
		c.Log.SetBy(ngingConfig, defaults).Init()
	}
}
