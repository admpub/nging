/*
   Nging is a toolbox for webmasters
   Copyright (C) 2018-present Wenhui Shen <swh@admpub.com>

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
	"fmt"
	"os"
	"reflect"

	"github.com/admpub/color"
	"github.com/admpub/confl"
	"github.com/admpub/log"
	"github.com/admpub/securecookie"
	"github.com/webx-top/codec"
	"github.com/webx-top/com"
	"github.com/webx-top/db"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/middleware/language"

	"github.com/admpub/nging/v5/application/library/common"
	"github.com/admpub/nging/v5/application/library/config/extend"
	"github.com/admpub/nging/v5/application/library/config/subconfig/scookie"
	"github.com/admpub/nging/v5/application/library/config/subconfig/scron"
	"github.com/admpub/nging/v5/application/library/config/subconfig/sdb"
	"github.com/admpub/nging/v5/application/library/config/subconfig/ssystem"
)

func NewConfig() *Config {
	c := &Config{
		Validations: Validations{},
	}
	c.InitExtend()
	c.settings = NewSettings(c)
	c.DB.MaxIdleConns = db.DefaultSettings.MaxIdleConns()
	c.DB.MaxOpenConns = db.DefaultSettings.MaxOpenConns()
	return c
}

type Config struct {
	DB       sdb.DB          `json:"db"`
	Sys      ssystem.System  `json:"sys"`
	Cron     scron.Cron      `json:"cron"`
	Cookie   scookie.Config  `json:"cookie"`
	Language language.Config `json:"language"`
	Extend   echo.H          `json:"extend,omitempty"`
	settings *Settings       `json:"-"`

	// 自定义validator 验证规则。map 的 key 为规则标识，值为规则
	Validations Validations `json:"validations"`

	connectedDB bool
	reloaders   []func(newConfig *Config)
}

func (c *Config) IsEnv(env string) bool {
	return c.Sys.IsEnv(env)
}

func (c *Config) IsEnvProd() bool {
	return c.Sys.IsEnv(`prod`)
}

func (c *Config) IsEnvDev() bool {
	return c.Sys.IsEnv(`dev`)
}

func (c *Config) Settings() *Settings {
	return c.settings
}

func (c *Config) Debug() bool {
	return c.settings.Debug
}

func (c *Config) GetMaxRequestBodySize() int {
	if c.settings.MaxRequestBodySizeBytes() > 0 {
		return c.settings.MaxRequestBodySizeBytes()
	}
	return c.Sys.MaxRequestBodySizeBytes()
}

// ConnectedDB 数据库是否已连接，如果没有连接则自动连接
func (c *Config) ConnectedDB(autoConn ...bool) bool {
	if c.connectedDB {
		return c.connectedDB
	}
	n := len(autoConn)
	if n == 0 || (n > 0 && autoConn[0]) {
		err := c.connectDB()
		if err != nil {
			log.Error(err)
		}
	}
	return c.connectedDB
}

func (c *Config) connectDB() error {
	err := ConnectDB(c.DB, 0, `default`)
	if err != nil {
		return err
	}
	c.connectedDB = true
	return nil
}

func (c *Config) APIKey() string {
	return c.settings.APIKey
}

func (c *Config) GetExtend() echo.H {
	return c.Extend
}

func (c *Config) CookieConfig() scookie.Config {
	return c.Cookie
}

func (c *Config) InitExtend() *Config {
	c.Extend = echo.H{}
	extend.Range(c.registerExtend)
	return c
}

func (c *Config) PrintExtend() {
	for k, v := range c.Extend {
		fmt.Printf("[Extend] %s : %T\n", k, v)
	}
}

func (c *Config) registerExtend(key string, recv interface{}) {
	fmt.Printf(color.YellowString(`[Register Extend Config]`)+` `+color.MagentaString(`P%d`, FromCLI().Pid())+` %s: %T`+"\n", key, recv)
	c.Extend[key] = recv
}

func (c *Config) UnregisterExtend(key string) {
	if recv, ok := c.Extend[key]; ok {
		fmt.Printf(color.YellowString(`[Unregister Extend Config]`)+` `+color.MagentaString(`P%d`, FromCLI().Pid())+` %s: %T`+"\n", key, recv)
		delete(c.Extend, key)
	}
}

func (c *Config) ConfigFromDB() echo.H {
	return c.settings.GetConfig()
}

func (c *Config) SetDebug(on bool) *Config {
	c.settings.SetDebug(on)
	return c
}

func (c *Config) Codec(lengths ...int) codec.Codec {
	length := Bit128
	if len(lengths) > 0 {
		length = lengths[0]
	}
	if length == Bit256 {
		return default256Codec
	}
	return defaultCodec
}

const (
	Bit128 = 128
	Bit256 = 256
)

var (
	defaultCodec    = codec.NewAES(`AES-128-CBC`)
	default256Codec = codec.NewAES(`AES-256-CBC`)
)

func (c *Config) encode(bit int, raw string, keys ...string) string {
	var key string
	if len(keys) > 0 && len(keys[0]) > 0 {
		key = com.Md5(keys[0])
	} else {
		key = c.Cookie.HashKey
	}
	return c.Codec(bit).Encode(raw, key)
}

func (c *Config) decode(bit int, encrypted string, keys ...string) string {
	if len(encrypted) == 0 {
		return ``
	}
	var key string
	if len(keys) > 0 && len(keys[0]) > 0 {
		key = com.Md5(keys[0])
	} else {
		key = c.Cookie.HashKey
	}
	return c.Codec(bit).Decode(encrypted, key)
}

func (c *Config) Encode(raw string, keys ...string) string {
	return c.encode(Bit128, raw, keys...)
}

func (c *Config) Decode(encrypted string, keys ...string) string {
	return c.decode(Bit128, encrypted, keys...)
}

func (c *Config) Encode256(raw string, keys ...string) string {
	return c.encode(Bit256, raw, keys...)
}

func (c *Config) Decode256(encrypted string, keys ...string) string {
	return c.decode(Bit256, encrypted, keys...)
}

func (c *Config) InitSecretKey() *Config {
	c.Cookie.BlockKey = c.GenerateRandomKey()
	c.Cookie.HashKey = c.GenerateRandomKey()
	return c
}

func (c *Config) GenerateRandomKey() string {
	return com.ByteMd5(securecookie.GenerateRandomKey(32))
}

func (c *Config) AddReloader(rd ...func(newConfig *Config)) *Config {
	c.reloaders = append(c.reloaders, rd...)
	return c
}

func (c *Config) Reload(newConfig *Config) error {
	var engines []string
	for name, newExtConfig := range newConfig.Extend {
		oldExtConfig := c.Extend.Get(name)
		if !reflect.DeepEqual(newExtConfig, oldExtConfig) {
			log.Debugf(`reloading extend config: %v`, name)
			if newRd, ok := newExtConfig.(extend.Reloader); ok {
				log.Debugf(`reload extend config(Reload): %v`, name)
				if err := newRd.Reload(); err != nil {
					log.Errorf(`failed to %s.Reload: %v`, name, err)
				}
				continue
			}
			engines = append(engines, name)
		}
	}
	if !reflect.DeepEqual(c.Validations, newConfig.Validations) {
		log.Debug(`reloading validation config`)
		if err := newConfig.Validations.Register(); err != nil {
			return err
		}
	}
	for _, rd := range c.reloaders {
		rd(newConfig)
	}
	return FromCLI().Reload(newConfig, engines...)
}

var configInitors = []func(*Config){}

func AddConfigInitor(initors ...func(*Config)) {
	configInitors = append(configInitors, initors...)
}

func (c *Config) AsDefault() *Config {
	c.Validations.Register()
	echo.Set(common.ConfigName, c)
	for _, initor := range configInitors {
		initor(c)
	}
	defaultConfigMu.Lock()
	defaultConfig = c
	defaultConfigMu.Unlock()
	err := c.settings.Init(nil)
	if err != nil {
		log.Errorf(`failed to config.settings.Init: %v`, err)
	}
	return c
}

func (c *Config) SaveToFile() error {
	b, err := confl.Marshal(c)
	if err != nil {
		return err
	}
	dir := FromCLI().ConfDir()
	err = com.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return err
	}
	err = os.WriteFile(FromCLI().Conf, b, os.ModePerm)
	return err
}

func (c *Config) GenerateSample() error {
	_, err := os.Stat(FromCLI().Conf + `.sample`)
	if err == nil || !os.IsNotExist(err) {
		return err
	}
	var old []byte
	old, err = os.ReadFile(FromCLI().Conf)
	if err == nil {
		err = os.WriteFile(FromCLI().Conf+`.sample`, old, os.ModePerm)
	}
	return err
}

func (c *Config) SetDefaults() *Config {
	c.Sys.Init()
	if len(c.Cookie.Path) == 0 {
		c.Cookie.Path = `/`
	}
	for _, value := range c.Extend {
		if sd, ok := value.(extend.SetDefaults); ok {
			sd.SetDefaults()
		}
	}
	return c
}
