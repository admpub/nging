// Copyright 2016 fatedier, fatedier@gmail.com
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/admpub/ini"
)

// ClientCommonConf client common config
type ClientCommonConf struct {
	ServerAddr        string              `json:"server_addr"`
	ServerPort        int                 `json:"server_port"`
	HttpProxy         string              `json:"http_proxy"`
	LogFile           string              `json:"log_file"`
	LogWay            string              `json:"log_way"`
	LogLevel          string              `json:"log_level"`
	LogMaxDays        int64               `json:"log_max_days"`
	Token             string              `json:"token"`
	AdminAddr         string              `json:"admin_addr"`
	AdminPort         int                 `json:"admin_port"`
	AdminUser         string              `json:"admin_user"`
	AdminPwd          string              `json:"admin_pwd"`
	PoolCount         int                 `json:"pool_count"`
	TcpMux            bool                `json:"tcp_mux"`
	User              string              `json:"user"`
	DnsServer         string              `json:"dns_server"`
	LoginFailExit     bool                `json:"login_fail_exit"`
	Start             map[string]struct{} `json:"start"`
	Protocol          string              `json:"protocol"`
	TLSEnable         bool                `json:"tls_enable"`
	HeartBeatInterval int64               `json:"heartbeat_interval"`
	HeartBeatTimeout  int64               `json:"heartbeat_timeout"`
}

func GetDefaultClientConf() *ClientCommonConf {
	return &ClientCommonConf{
		ServerAddr:        "0.0.0.0",
		ServerPort:        7000,
		HttpProxy:         os.Getenv("http_proxy"),
		LogFile:           "console",
		LogWay:            "console",
		LogLevel:          "info",
		LogMaxDays:        3,
		Token:             "",
		AdminAddr:         "127.0.0.1",
		AdminPort:         0,
		AdminUser:         "",
		AdminPwd:          "",
		PoolCount:         1,
		TcpMux:            true,
		User:              "",
		DnsServer:         "",
		LoginFailExit:     true,
		Start:             make(map[string]struct{}),
		Protocol:          "tcp",
		HeartBeatInterval: 30,
		HeartBeatTimeout:  90,
	}
}

func UnmarshalClientConfFromIni(defaultCfg *ClientCommonConf, content string) (cfg *ClientCommonConf, err error) {
	cfg = defaultCfg
	if cfg == nil {
		cfg = GetDefaultClientConf()
	}

	conf, err := ini.LoadContent(content)
	if err != nil {
		err = fmt.Errorf("parse ini conf file error: %v", err)
		return nil, err
	}

	var (
		v int64
	)
	commonSection := conf.Section("common")
	if tmpStr := commonSection.Key("server_addr").String(); len(tmpStr) > 0 {
		cfg.ServerAddr = tmpStr
	}

	if tmpStr := commonSection.Key("server_port").String(); len(tmpStr) > 0 {
		v, err = strconv.ParseInt(tmpStr, 10, 64)
		if err != nil {
			err = fmt.Errorf("Parse conf error: invalid server_port")
			return
		}
		cfg.ServerPort = int(v)
	}

	if tmpStr := commonSection.Key("http_proxy").String(); len(tmpStr) > 0 {
		cfg.HttpProxy = tmpStr
	}

	if tmpStr := commonSection.Key("log_file").String(); len(tmpStr) > 0 {
		cfg.LogFile = tmpStr
		if cfg.LogFile == "console" {
			cfg.LogWay = "console"
		} else {
			cfg.LogWay = "file"
		}
	}

	if tmpStr := commonSection.Key("log_level").String(); len(tmpStr) > 0 {
		cfg.LogLevel = tmpStr
	}

	if tmpStr := commonSection.Key("log_max_days").String(); len(tmpStr) > 0 {
		if v, err = strconv.ParseInt(tmpStr, 10, 64); err == nil {
			cfg.LogMaxDays = v
		}
	}

	if tmpStr := commonSection.Key("token").String(); len(tmpStr) > 0 {
		cfg.Token = tmpStr
	}

	if tmpStr := commonSection.Key("admin_addr").String(); len(tmpStr) > 0 {
		cfg.AdminAddr = tmpStr
	}

	if tmpStr := commonSection.Key("admin_port").String(); len(tmpStr) > 0 {
		if v, err = strconv.ParseInt(tmpStr, 10, 64); err == nil {
			cfg.AdminPort = int(v)
		} else {
			err = fmt.Errorf("Parse conf error: invalid admin_port")
			return
		}
	}

	if tmpStr := commonSection.Key("admin_user").String(); len(tmpStr) > 0 {
		cfg.AdminUser = tmpStr
	}

	if tmpStr := commonSection.Key("admin_pwd").String(); len(tmpStr) > 0 {
		cfg.AdminPwd = tmpStr
	}

	if tmpStr := commonSection.Key("pool_count").String(); len(tmpStr) > 0 {
		if v, err = strconv.ParseInt(tmpStr, 10, 64); err == nil {
			cfg.PoolCount = int(v)
		}
	}

	if tmpStr := commonSection.Key("tcp_mux").String(); tmpStr == "false" {
		cfg.TcpMux = false
	} else {
		cfg.TcpMux = true
	}

	if tmpStr := commonSection.Key("user").String(); len(tmpStr) > 0 {
		cfg.User = tmpStr
	}

	if tmpStr := commonSection.Key("dns_server").String(); len(tmpStr) > 0 {
		cfg.DnsServer = tmpStr
	}

	if tmpStr := commonSection.Key("start").String(); len(tmpStr) > 0 {
		proxyNames := strings.Split(tmpStr, ",")
		for _, name := range proxyNames {
			cfg.Start[strings.TrimSpace(name)] = struct{}{}
		}
	}

	if tmpStr := commonSection.Key("login_fail_exit").String(); tmpStr == "false" {
		cfg.LoginFailExit = false
	} else {
		cfg.LoginFailExit = true
	}

	if tmpStr := commonSection.Key("protocol").String(); len(tmpStr) > 0 {
		// Now it only support tcp and kcp and websocket.
		if tmpStr != "tcp" && tmpStr != "kcp" && tmpStr != "websocket" {
			err = fmt.Errorf("Parse conf error: invalid protocol")
			return
		}
		cfg.Protocol = tmpStr
	}

	if tmpStr := commonSection.Key("heartbeat_timeout").String(); len(tmpStr) > 0 {
		if v, err = strconv.ParseInt(tmpStr, 10, 64); err != nil {
			err = fmt.Errorf("Parse conf error: invalid heartbeat_timeout")
			return
		} else {
			cfg.HeartBeatTimeout = v
		}
	}

	if tmpStr := commonSection.Key("heartbeat_interval").String(); len(tmpStr) > 0 {
		if v, err = strconv.ParseInt(tmpStr, 10, 64); err != nil {
			err = fmt.Errorf("Parse conf error: invalid heartbeat_interval")
			return
		} else {
			cfg.HeartBeatInterval = v
		}
	}
	return
}

func (cfg *ClientCommonConf) Check() (err error) {
	if cfg.HeartBeatInterval <= 0 {
		err = fmt.Errorf("Parse conf error: invalid heartbeat_interval")
		return
	}

	if cfg.HeartBeatTimeout < cfg.HeartBeatInterval {
		err = fmt.Errorf("Parse conf error: invalid heartbeat_timeout, heartbeat_timeout is less than heartbeat_interval")
		return
	}
	return
}
