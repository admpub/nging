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
	stdLog "log"
	"net/url"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/admpub/confl"
	"github.com/admpub/log"
	"github.com/admpub/mysql-schema-sync/sync"
	"github.com/admpub/nging/v5/application/cmd/bootconfig"
	"github.com/admpub/nging/v5/application/library/common"
	"github.com/admpub/nging/v5/application/library/config/subconfig/sdb"
	"github.com/admpub/nging/v5/application/library/config/subconfig/ssystem"
	"github.com/admpub/nging/v5/application/library/cron"
	cronSend "github.com/admpub/nging/v5/application/library/cron/send"
	"github.com/webx-top/com"
	"github.com/webx-top/db/lib/factory"
	"github.com/webx-top/db/lib/sqlbuilder"
	"github.com/webx-top/db/mysql"
	"github.com/webx-top/echo"
)

func MustGetConfig() *Config {
	if FromFile() == nil {
		FromCLI().ParseConfig()
	}
	return FromFile()
}

func InitConfig() (*Config, error) {
	configFiles := []string{
		FromCLI().Conf,
		filepath.Join(FromCLI().Confd, `config.yaml.sample`),
	}
	var (
		configFile      string
		err             error
		temporaryConfig = NewConfig()
	)
	temporaryConfig.Debug = bootconfig.Develop
	for key, conf := range configFiles {
		if !filepath.IsAbs(conf) {
			conf = filepath.Join(echo.Wd(), conf)
			configFiles[key] = conf
			if key == 0 {
				FromCLI().Conf = conf
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
	temporaryConfig.SetDefaults(configFile)

	return temporaryConfig, nil
}

func ParseConfig() error {
	if false {
		b, err := confl.Marshal(FromFile())
		if err != nil {
			return err
		}
		err = os.WriteFile(FromCLI().Conf, b, os.ModePerm)
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
		if FromFile() != nil {
			if !FromFile().connectedDB || !reflect.DeepEqual(conf.DB, FromFile().DB) {
				if err = conf.connectDB(); err != nil {
					return err
				}
			}
			err = FromFile().Reload(conf)
		} else {
			err = conf.connectDB()
		}
		if err != nil {
			return err
		}
	}
	conf.AsDefault()
	return err
}

var (
	DBConnecters = map[string]func(*Config) error{
		`mysql`: ConnectMySQL,
	}
	DBInstallers = map[string]func(string) error{
		`mysql`: ExecMySQL,
	}
	DBCreaters = map[string]func(error, *Config) error{
		`mysql`: CreaterMySQL,
	}
	DBUpgraders = map[string]func(string, *sync.Config, *Config) (DBOperators, error){
		`mysql`: UpgradeMySQL,
	}
	DBEngines         = echo.NewKVData().Add(`mysql`, `MySQL`)
	ParseTimeDuration = ssystem.ParseTimeDuration
	ParseBytes        = ssystem.ParseBytes
)

type DBOperators struct {
	Source      sync.DBOperator
	Destination sync.DBOperator
}

func CreaterMySQL(err error, c *Config) error {
	if strings.Contains(err.Error(), `Unknown database`) {
		dbName := c.DB.Database
		c.DB.Database = ``
		err2 := ConnectDB(c)
		if err2 != nil {
			return err
		}
		charset := c.DB.Charset()
		if len(charset) == 0 {
			charset = sdb.MySQLDefaultCharset
		}
		sqlStr := "CREATE DATABASE `" + dbName + "` COLLATE '" + charset + "_general_ci'"
		_, err = factory.NewParam().SetCollection(sqlStr).Exec()
		if err != nil {
			return err
		}
		c.DB.Database = dbName
		err = ConnectDB(c)
	}
	return err
}

func UpgradeMySQL(schema string, syncConfig *sync.Config, cfg *Config) (DBOperators, error) {
	syncConfig.DestDSN = cfg.DB.User + `:` + cfg.DB.Password + `@(` + cfg.DB.Host + `)/` + cfg.DB.Database
	t := `?`
	for key, value := range cfg.DB.Options {
		syncConfig.DestDSN += t + fmt.Sprintf("%s=%s", key, url.QueryEscape(value))
		t = `&`
	}
	syncConfig.SQLPreprocessor = func() func(string) string {
		charset := cfg.DB.Charset()
		if len(charset) == 0 {
			charset = sdb.MySQLDefaultCharset
		}
		return func(sqlStr string) string {
			return common.ReplaceCharset(sqlStr, charset)
		}
	}()
	return DBOperators{Source: sync.NewMySchemaData(schema, `source`)}, nil
}

func ConnectMySQL(c *Config) error {
	settings := c.DB.ToMySQL()
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
	return FromFile().DB.Table(table)
}

func ToTable(m sqlbuilder.Name_) string {
	return FromFile().DB.ToTable(m)
}
