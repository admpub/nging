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
	"errors"
	"flag"
	"strings"

	"github.com/admpub/confl"
	"github.com/admpub/log"
	"github.com/webx-top/db/lib/factory"
	"github.com/webx-top/db/mongo"
	"github.com/webx-top/db/mysql"
)

type CLIConfig struct {
	Port int
	Conf string
}

func (c *CLIConfig) InitFlag() {
	flag.IntVar(&c.Port, `p`, 9999, `port`)
	flag.StringVar(&c.Conf, `c`, `config/config.yaml`, `config`)
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
		Accounts map[string]string `json:"accounts"`
	} `json:"sys"`
}

var (
	DefaultConfig         = &Config{}
	DefaultCLIConfig      = &CLIConfig{}
	ErrUnknowDatabaseType = errors.New(`unkown database type`)
)

func ParseConfig() error {
	_, err := confl.DecodeFile(DefaultCLIConfig.Conf, DefaultConfig)
	if err != nil {
		return err
	}
	InitLog()
	return ConnectDB()
}

func ConnectDB() error {
	factory.CloseAll()
	switch DefaultConfig.DB.Type {
	case `mysql`:
		settings := mysql.ConnectionURL{
			Host:     DefaultConfig.DB.Host,
			Database: DefaultConfig.DB.Database,
			User:     DefaultConfig.DB.User,
			Password: DefaultConfig.DB.Password,
			Options:  DefaultConfig.DB.Options,
		}
		database, err := mysql.Open(settings)
		if err != nil {
			return err
		}
		factory.SetDebug(DefaultConfig.DB.Debug)
		cluster := factory.NewCluster().AddW(database)
		factory.SetCluster(0, cluster).Cluster(0).SetPrefix(DefaultConfig.DB.Prefix)
	case `mongo`:
		settings := mongo.ConnectionURL{
			Host:     DefaultConfig.DB.Host,
			Database: DefaultConfig.DB.Database,
			User:     DefaultConfig.DB.User,
			Password: DefaultConfig.DB.Password,
			Options:  DefaultConfig.DB.Options,
		}
		database, err := mongo.Open(settings)
		if err != nil {
			return err
		}
		factory.SetDebug(DefaultConfig.DB.Debug)
		cluster := factory.NewCluster().AddW(database)
		factory.SetCluster(0, cluster).Cluster(0).SetPrefix(DefaultConfig.DB.Prefix)
	default:
		return ErrUnknowDatabaseType
	}
	return nil
}

func InitLog() {

	//======================================================
	// 配置日志
	//======================================================
	if DefaultConfig.Log.Debug {
		log.DefaultLog.MaxLevel = log.LevelDebug
	} else {
		log.DefaultLog.MaxLevel = log.LevelInfo
	}
	targets := []log.Target{}

	for _, targetName := range strings.Split(DefaultConfig.Log.Targets, `,`) {
		switch targetName {
		case "console":
			//输出到命令行
			consoleTarget := log.NewConsoleTarget()
			consoleTarget.ColorMode = DefaultConfig.Log.Colorable
			targets = append(targets, consoleTarget)

		case "file":
			//输出到文件
			fileTarget := log.NewFileTarget()
			fileTarget.FileName = DefaultConfig.Log.SaveFile
			fileTarget.Filter.MaxLevel = log.LevelInfo
			fileTarget.MaxBytes = DefaultConfig.Log.FileMaxBytes
			targets = append(targets, fileTarget)
		}
	}

	log.SetTarget(targets...)
	log.SetFatalAction(log.ActionExit)
}

func MustOK(err error) {
	if err != nil {
		log.Fatal(err.Error())
	}
}
