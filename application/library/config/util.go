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
	"io/ioutil"
	stdLog "log"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/admpub/confl"
	"github.com/admpub/log"
	"github.com/admpub/nging/application/cmd/event"
	"github.com/admpub/nging/application/library/caddy"
	"github.com/admpub/nging/application/library/common"
	"github.com/admpub/nging/application/library/cron"
	cronSend "github.com/admpub/nging/application/library/cron/send"
	"github.com/admpub/nging/application/library/ftp"
	"github.com/webx-top/com"
	"github.com/webx-top/db/lib/factory"
	"github.com/webx-top/db/lib/sqlbuilder"
	"github.com/webx-top/db/mongo"
	"github.com/webx-top/db/mysql"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/middleware/bytes"
)

var (
	reNumeric                  = regexp.MustCompile(`^[0-9]+$`)
	defaultMaxRequestBodyBytes = 2 << 20 // 2M
)

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

func InitConfig() (*Config, error) {
	configFiles := []string{
		DefaultCLIConfig.Conf,
		filepath.Join(echo.Wd(), `config/config.yaml.sample`),
	}
	var (
		configFile      string
		err             error
		temporaryConfig = NewConfig()
	)
	temporaryConfig.Debug = event.Develop
	for key, conf := range configFiles {
		if !filepath.IsAbs(conf) {
			conf = filepath.Join(echo.Wd(), conf)
			configFiles[key] = conf
			if key == 0 {
				DefaultCLIConfig.Conf = conf
			}
		}
		_, err = os.Stat(conf)
		if err == nil {
			configFile = conf
			break
		}
		if !os.IsNotExist(err) {
			return temporaryConfig, err
		}
	}
	if err != nil {
		return temporaryConfig, err
	}
	_, err = confl.DecodeFile(configFile, temporaryConfig)
	if err != nil {
		return temporaryConfig, err
	}
	confDir := filepath.Dir(configFile)
	if len(temporaryConfig.Caddy.Caddyfile) == 0 {
		temporaryConfig.Caddy.Caddyfile = `./Caddyfile`
	} else if strings.HasSuffix(temporaryConfig.Caddy.Caddyfile, `/`) || strings.HasSuffix(temporaryConfig.Caddy.Caddyfile, `\`) {
		temporaryConfig.Caddy.Caddyfile = path.Join(temporaryConfig.Caddy.Caddyfile, `Caddyfile`)
	}
	if len(temporaryConfig.Sys.VhostsfileDir) == 0 {
		temporaryConfig.Sys.VhostsfileDir = path.Join(confDir, `vhosts`)
	}
	if temporaryConfig.Sys.MaxRequestBodySize <= 0 {
		temporaryConfig.Sys.MaxRequestBodySize = defaultMaxRequestBodyBytes
	}
	if temporaryConfig.Sys.EditableFileMaxBytes < 1 && len(temporaryConfig.Sys.EditableFileMaxSize) > 0 {
		temporaryConfig.Sys.EditableFileMaxBytes, err = bytes.Parse(temporaryConfig.Sys.EditableFileMaxSize)
		if err != nil {
			log.Error(err.Error())
		}
	}
	temporaryConfig.Sys.CmdTimeoutDuration = ParseTimeDuration(temporaryConfig.Sys.CmdTimeout)
	if temporaryConfig.Sys.CmdTimeoutDuration <= 0 {
		temporaryConfig.Sys.CmdTimeoutDuration = time.Second * 30
	}
	if len(temporaryConfig.Cookie.Path) == 0 {
		temporaryConfig.Cookie.Path = `/`
	}
	if len(temporaryConfig.Sys.SSLCacheDir) == 0 {
		temporaryConfig.Sys.SSLCacheDir = filepath.Join(echo.Wd(), `data`, `cache`, `autocert`)
	}
	caddy.Fixed(&temporaryConfig.Caddy)
	ftp.Fixed(&temporaryConfig.FTP)

	return temporaryConfig, nil
}

func ParseConfig() error {
	if false {
		b, err := confl.Marshal(DefaultConfig)
		err = ioutil.WriteFile(DefaultCLIConfig.Conf, b, os.ModePerm)
		if err != nil {
			return err
		}
	}
	conf, err := InitConfig()
	if err != nil {
		return err
	}
	InitSessionOptions(conf)
	if conf.Cron.PoolSize > 0 {
		cron.PoolSize = conf.Cron.PoolSize
	}
	cronSend.DefaultEmailConfig.Template = conf.Cron.Template
	if IsInstalled() {
		err = conf.connectDB()
		if err != nil {
			return err
		}
		if DefaultConfig != nil {
			err = DefaultConfig.Reload(conf)
			if err != nil {
				return err
			}
		}
	}
	conf.AsDefault()
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
		err2 := ConnectDB(c)
		if err2 != nil {
			return err
		}
		sqlStr := "CREATE DATABASE `" + dbName + "`"
		_, err = factory.NewParam().SetCollection(sqlStr).Exec()
		if err != nil {
			return err
		}
		c.DB.Database = dbName
		err = ConnectDB(c)
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
	common.ParseMysqlConnectionURL(&settings)
	if settings.Options == nil {
		settings.Options = map[string]string{}
	}
	// Default options.
	if _, ok := settings.Options["charset"]; !ok {
		settings.Options["charset"] = "utf8mb4"
	}
	database, err := mysql.Open(settings)
	if err != nil {
		return err
	}
	c.DB.SetConn(database)
	cluster := factory.NewCluster().AddMaster(database)
	factory.SetCluster(0, cluster)
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
	if c.DB.ConnMaxDuration() > 0 {
		mongo.ConnTimeout = c.DB.ConnMaxDuration()
	}
	database, err := mongo.Open(settings)
	if err != nil {
		return err
	}
	c.DB.SetConn(database)
	cluster := factory.NewCluster().AddMaster(database)
	factory.SetCluster(0, cluster)
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

func QueryTo(sqlStr string, result interface{}) (sqlbuilder.Iterator, error) {
	return factory.NewParam().SetRecv(result).SetCollection(sqlStr).QueryTo()
}

func ConnectDB(c *Config) error {
	factory.CloseAll()
	if fn, ok := DBConnecters[c.DB.Type]; ok {
		return fn(c)
	}
	return ErrUnknowDatabaseType
}

func MustOK(err error) {
	if err != nil {
		log.Fatal(err.Error())
	}
}

var CmdIsRunning = com.CmdIsRunning

func Table(table string) string {
	return DefaultConfig.DB.Table(table)
}

func ToTable(m sqlbuilder.Name_) string {
	return DefaultConfig.DB.ToTable(m)
}
