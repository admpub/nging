/*

   Copyright 2016 Wenhui Shen <www.webx.top>

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.

*/
package config

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"time"

	"github.com/admpub/confl"
	"github.com/admpub/mail"
	"github.com/admpub/nging/application/library/caddy"
	"github.com/admpub/nging/application/library/ftp"
	"github.com/webx-top/echo/middleware/language"
)

func SetVersion(version string) {
	caddy.DefaultVersion = version
}

type Config struct {
	DB struct {
		Type     string            `json:"type"`
		User     string            `json:"user"`
		Password string            `json:"password"`
		Host     string            `json:"host"`
		Database string            `json:"database"`
		Prefix   string            `json:"prefix"`
		Options  map[string]string `json:"options"`
		Debug    bool              `json:"debug"`
	} `json:"db"`

	Log struct {
		Debug        bool   `json:"debug"`
		Colorable    bool   `json:"colorable"`    // for console
		SaveFile     string `json:"saveFile"`     // for file
		FileMaxBytes int64  `json:"fileMaxBytes"` // for file
		Targets      string `json:"targets"`
	} `json:"log"`

	Sys struct {
		VhostsfileDir          string            `json:"vhostsfileDir"`
		AllowIP                []string          `json:"allowIP"`
		SSLHosts               []string          `json:"sslHosts"`
		SSLCacheFile           string            `json:"sslCacheFile"`
		SSLKeyFile             string            `json:"sslKeyFile"`
		SSLCertFile            string            `json:"sslCertFile"`
		Debug                  bool              `json:"debug"`
		EditableFileExtensions map[string]string `json:"editableFileExtensions"`
		EditableFileMaxSize    string            `json:"editableFileMaxSize"`
		EditableFileMaxBytes   int64             `json:"editableFileMaxBytes"`
		ErrorPages             map[int]string    `json:"errorPages"`
		CmdTimeout             string            `json:"cmdTimeout"`
		CmdTimeoutDuration     time.Duration     `json:"-"`
	} `json:"sys"`

	Email struct {
		*mail.SMTPConfig
		Timeout   int64  //超时时间(秒)，采用默认引擎发信时，此项无效
		Engine    string //值为email时采用github.com/jordan-wright/email包发送，否则采用默认的github.com/admpub/mail发送
		From      string //发信人Email地址
		QueueSize int    //允许同一时间发信的数量
	} `json:"email"`

	Cron struct {
		PoolSize int    `json:"poolSize"`
		Template string `json:"template"` //发信模板
	} `json:"cron"`

	Cookie struct {
		Domain   string `json:"domain"`
		MaxAge   int    `json:"maxAge"`
		Path     string `json:"path"`
		HttpOnly bool   `json:"httpOnly"`
		HashKey  string `json:"hashKey"`
		BlockKey string `json:"blockKey"`
	} `json:"cookie"`

	Caddy    caddy.Config    `json:"caddy"`
	FTP      ftp.Config      `json:"ftp"`
	Language language.Config `json:"language"`
}

func (c *Config) SaveToFile() error {
	b, err := confl.Marshal(c)
	if err == nil {
		err = os.MkdirAll(filepath.Dir(DefaultCLIConfig.Conf), os.ModePerm)
		if err == nil {
			_, e := os.Stat(DefaultCLIConfig.Conf + `.sample`)
			if e != nil {
				if os.IsNotExist(e) {
					old, err := ioutil.ReadFile(DefaultCLIConfig.Conf)
					if err == nil {
						err = ioutil.WriteFile(DefaultCLIConfig.Conf+`.sample`, old, os.ModePerm)
					}
					if err != nil {
						return err
					}
				}
			}
			err = ioutil.WriteFile(DefaultCLIConfig.Conf, b, os.ModePerm)
		}
	}
	return err
}
