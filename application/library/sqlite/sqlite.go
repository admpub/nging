//go:build sqlite
// +build sqlite

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

package sqlite

import (
	"os"
	"strings"

	syncSQLite "github.com/admpub/mysql-schema-sync/sqlite"
	"github.com/admpub/mysql-schema-sync/sync"
	"github.com/admpub/nging/v5/application/library/config"
	"github.com/webx-top/db/lib/factory"
	"github.com/webx-top/db/sqlite"
)

func register() {
	config.DBCreaters[`sqlite`] = CreaterSQLite
	config.DBConnecters[`sqlite`] = ConnectSQLite
	config.DBInstallers[`sqlite`] = ExecSQL
	config.DBUpgraders[`sqlite`] = UpgradeSQLite
	config.DBEngines.Add(`sqlite`, `SQLite`)
}

func ConnectSQLite(c *config.Config) error {
	settings := sqlite.ConnectionURL{
		Database: c.DB.Database,
		Options:  c.DB.Options,
	}
	database, err := sqlite.Open(settings)
	if err != nil {
		return err
	}
	cluster := factory.NewCluster().AddMaster(database)
	factory.SetCluster(0, cluster)
	factory.SetDebug(c.DB.Debug)
	return nil
}

func CreaterSQLite(err error, c *config.Config) error {
	if strings.Contains(err.Error(), `unable to open database file`) {
		var f *os.File
		f, err = os.Create(c.DB.Database)
		if err == nil {
			f.Close()
			err = config.ConnectDB(c)
		}
	}
	return err
}

func UpgradeSQLite(schema string, syncConfig *sync.Config, cfg *config.Config) (config.DBOperators, error) {
	syncConfig.DestDSN = cfg.DB.Database
	syncConfig.Comparer = syncSQLite.NewCompare()
	var err error
	schema, err = ConvertMySQL(schema)
	return config.DBOperators{
		Source:      syncSQLite.NewSchemaData(schema, `source`),
		Destination: syncSQLite.New(cfg.DB.Database, `dest`),
	}, err
}
