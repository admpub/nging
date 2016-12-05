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
	"github.com/admpub/caddyui/application/library/caddy"
	"github.com/admpub/caddyui/application/library/ftp"
)

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
		VhostsfileDir string            `json:"vhostsfileDir"`
		AllowIP       []string          `json:"allowIP"`
		Accounts      map[string]string `json:"accounts"`
		SSLHosts      []string          `json:"sslHosts"`
		SSLCacheFile  string            `json:"sslCacheFile"`
		SSLKeyFile    string            `json:"sslKeyFile"`
		SSLCertFile   string            `json:"sslCertFile"`
		Debug         bool              `json:"debug"`
	} `json:"sys"`

	Cookie struct {
		Domain   string `json:"domain"`
		MaxAge   int    `json:"maxAge"`
		Path     string `json:"path"`
		HttpOnly bool   `json:"httpOnly"`
		HashKey  string `json:"hashKey"`
		BlockKey string `json:"blockKey"`
	} `json:"cookie"`

	Caddy caddy.Config `json:"caddy"`
	FTP   ftp.Config   `json:"ftp"`
}
