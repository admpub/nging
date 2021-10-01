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
	"reflect"
	"strings"

	"github.com/webx-top/echo"
	"github.com/webx-top/echo/subdomains"

	"github.com/admpub/log"
	"github.com/admpub/nging/v3/application/library/common"
	"github.com/admpub/nging/v3/application/library/notice"
	"github.com/admpub/nging/v3/application/registry/settings"
)

var Setting = common.Setting

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
		r.Set(`base`, defaults.GetStore(`base`))
	}
	c.Base = r.GetStore(`base`)
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
	actGroups          = []string{`base`, `smtp`, `log`}
	onKeySetSettings   = map[string][]func(Diff) error{}
	onGroupSetSettings = map[string][]func(Diffs) error{}
)

func OnGroupSetSettings(groupAndKey string, fn func(Diffs) error) {
	if _, ok := onGroupSetSettings[groupAndKey]; !ok {
		onGroupSetSettings[groupAndKey] = []func(Diffs) error{}
	}
	onGroupSetSettings[groupAndKey] = append(onGroupSetSettings[groupAndKey], fn)
}

func OnKeySetSettings(groupAndKey string, fn func(Diff) error) {
	if _, ok := onKeySetSettings[groupAndKey]; !ok {
		onKeySetSettings[groupAndKey] = []func(Diff) error{}
	}
	onKeySetSettings[groupAndKey] = append(onKeySetSettings[groupAndKey], fn)
}

func FireInitSettings(configs echo.H) error {
	for group, fnList := range onGroupSetSettings {
		values := configs.GetStore(group)
		diffs := Diffs{}
		for k, v := range values {
			diffs[k] = &Diff{
				Old: v,
				New: v,
			}
		}
		for _, fn := range fnList {
			if err := fn(diffs); err != nil {
				return err
			}
		}
	}
	for groupAndKey, fnList := range onKeySetSettings {
		args := strings.SplitN(groupAndKey, `.`, 2)
		values := configs.GetStore(args[0])
		var val interface{}
		if len(args) == 2 {
			val = values.Get(args[1])
		} else {
			val = values
		}
		for _, fn := range fnList {
			if err := fn(Diff{
				Old: val,
				New: val,
			}); err != nil {
				return err
			}
		}
	}
	return nil
}

func FireSetSettings(group string, diffs Diffs) error {
	if fnList, ok := onGroupSetSettings[group]; ok {
		for _, fn := range fnList {
			if err := fn(diffs); err != nil {
				return err
			}
		}
	}
	for key, diff := range diffs {
		k := group + `.` + key
		if fnList, ok := onKeySetSettings[k]; ok {
			for _, fn := range fnList {
				if err := fn(*diff); err != nil {
					return err
				}
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

type Diff struct {
	Old    interface{}
	New    interface{}
	IsDiff bool
}

func (d Diff) String() string {
	s, _ := d.New.(string)
	return s
}

type Diffs map[string]*Diff

func (d Diffs) Get(key string) interface{} {
	return d[key]
}

func (c *Settings) SetConfigs(groups ...string) {
	newConfigs := settings.ConfigAsStore(groups...)
	oldConfigs := c.GetConfig()
	c.setConfigs(newConfigs, oldConfigs)
}

func (c *Settings) setConfigs(newConfigs echo.H, oldConfigs echo.H) {
	for group, conf := range newConfigs {
		keyCfg := conf.(echo.H)
		keyOldCfg := oldConfigs.GetStore(group)
		diffs := Diffs{}
		for k, v := range keyCfg {
			if !reflect.DeepEqual(keyOldCfg.Get(k), v) {
				diffs[k] = &Diff{
					Old:    keyOldCfg.Get(k),
					New:    v,
					IsDiff: true,
				}
			}
		}
		if len(diffs) == 0 {
			continue
		}
		oldConfigs.Set(group, keyCfg)
		FireSetSettings(group, diffs)
		//log.Debug(`Change configuration:`, group, `:`, echo.Dump(conf, false))
		c.SetConfig(group, oldConfigs, nil)
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
