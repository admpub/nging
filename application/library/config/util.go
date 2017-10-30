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
	stdLog "log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/admpub/confl"
	"github.com/admpub/log"
	"github.com/admpub/nging/application/library/caddy"
	"github.com/admpub/nging/application/library/ftp"
	"github.com/webx-top/db/lib/factory"
	"github.com/webx-top/db/mongo"
	"github.com/webx-top/db/mysql"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/middleware/bytes"
)

var reNumeric = regexp.MustCompile(`^[0-9]+$`)

func ParseTimeDuration(timeout string) time.Duration {
	var timeoutDuration time.Duration
	if len(timeout) > 0 {
		if reNumeric.MatchString(timeout) {
			if val, err := strconv.ParseUint(timeout, 10, 64); err != nil {
				log.Error(err)
			} else {
				timeoutDuration = time.Second * time.Duration(val)
			}
		} else {
			timeoutDuration, _ = time.ParseDuration(timeout)
		}
	}
	return timeoutDuration
}

func InitConfig() error {
	_, err := confl.DecodeFile(DefaultCLIConfig.Conf, DefaultConfig)
	if err != nil {
		return err
	}
	confDir := filepath.Dir(DefaultCLIConfig.Conf)
	if len(DefaultConfig.Caddy.Caddyfile) == 0 {
		DefaultConfig.Caddy.Caddyfile = `./Caddyfile`
	} else if strings.HasSuffix(DefaultConfig.Caddy.Caddyfile, `/`) || strings.HasSuffix(DefaultConfig.Caddy.Caddyfile, `\`) {
		DefaultConfig.Caddy.Caddyfile = path.Join(DefaultConfig.Caddy.Caddyfile, `Caddyfile`)
	}
	if len(DefaultConfig.Sys.VhostsfileDir) == 0 {
		DefaultConfig.Sys.VhostsfileDir = path.Join(confDir, `vhosts`)
	}
	if DefaultConfig.Sys.EditableFileMaxBytes < 1 && len(DefaultConfig.Sys.EditableFileMaxSize) > 0 {
		DefaultConfig.Sys.EditableFileMaxBytes, err = bytes.Parse(DefaultConfig.Sys.EditableFileMaxSize)
		if err != nil {
			log.Error(err.Error())
		}
	}
	DefaultConfig.Sys.CmdTimeoutDuration = ParseTimeDuration(DefaultConfig.Sys.CmdTimeout)
	if DefaultConfig.Sys.CmdTimeoutDuration <= 0 {
		DefaultConfig.Sys.CmdTimeoutDuration = time.Second * 30
	}
	caddy.Fixed(&DefaultConfig.Caddy)
	ftp.Fixed(&DefaultConfig.FTP)
	return nil
}

func ParseConfig() error {
	if false {
		b, err := confl.Marshal(DefaultConfig)
		err = ioutil.WriteFile(DefaultCLIConfig.Conf, b, os.ModePerm)
		if err != nil {
			return err
		}
	}
	err := InitConfig()
	if err != nil {
		return err
	}
	InitLog()
	InitSessionOptions()
	if DefaultConfig.Sys.Debug {
		log.SetLevel(`Debug`)
	} else {
		log.SetLevel(`Info`)
	}

	err = ConnectDB()
	if err != nil {
		return err
	}

	err = DefaultCLIConfig.Reload()
	if err != nil {
		return err
	}
	return err
}

var (
	DBConnecters = map[string]func(*Config) error{
		`mysql`: ConnectMySQL,
		`mongo`: ConnectMongoDB,
	}
	DBInstallers = map[string]func(string) error{
		`mysql`: ExecMySQL,
	}
	DBCreaters = map[string]func(error, *Config) error{
		`mysql`: CreaterMySQL,
	}
	DBEngines = echo.NewKVData().Add(`mysql`, `MySQL`)
)

func CreaterMySQL(err error, c *Config) error {
	if strings.Contains(err.Error(), `Unknown database`) {
		dbName := c.DB.Database
		c.DB.Database = ``
		err2 := ConnectDB()
		if err2 != nil {
			return err
		}
		sqlStr := "CREATE DATABASE `" + dbName + "`"
		_, err = factory.NewParam().SetCollection(sqlStr).Exec()
		if err != nil {
			return err
		}
		c.DB.Database = dbName
		err = ConnectDB()
	}
	return err
}

func ConnectMySQL(c *Config) error {
	settings := mysql.ConnectionURL{
		Host:     c.DB.Host,
		Database: c.DB.Database,
		User:     c.DB.User,
		Password: c.DB.Password,
		Options:  c.DB.Options,
	}
	database, err := mysql.Open(settings)
	if err != nil {
		return err
	}
	cluster := factory.NewCluster().AddMaster(database)
	factory.SetCluster(0, cluster).Cluster(0).SetPrefix(c.DB.Prefix)
	factory.SetDebug(c.DB.Debug)
	return nil
}

func ConnectMongoDB(c *Config) error {
	settings := mongo.ConnectionURL{
		Host:     c.DB.Host,
		Database: c.DB.Database,
		User:     c.DB.User,
		Password: c.DB.Password,
		Options:  c.DB.Options,
	}
	database, err := mongo.Open(settings)
	if err != nil {
		return err
	}
	cluster := factory.NewCluster().AddMaster(database)
	factory.SetCluster(0, cluster).Cluster(0).SetPrefix(c.DB.Prefix)
	factory.SetDebug(c.DB.Debug)
	return nil
}

func ExecMySQL(sqlStr string) error {
	_, err := factory.NewParam().SetCollection(sqlStr).Exec()
	if err != nil {
		stdLog.Println(err.Error(), `->SQL:`, sqlStr)
	}
	return err
}

func ConnectDB() error {
	factory.CloseAll()
	if fn, ok := DBConnecters[DefaultConfig.DB.Type]; ok {
		return fn(DefaultConfig)
	}
	return ErrUnknowDatabaseType
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

func CmdIsRunning(cmd *exec.Cmd) bool {
	return cmd != nil && cmd.ProcessState == nil
}
