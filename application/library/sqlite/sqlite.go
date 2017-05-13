// +build sqlite

package sqlite

import (
	"os"
	"strings"

	"github.com/admpub/nging/application/library/config"
	"github.com/webx-top/db/lib/factory"
	"github.com/webx-top/db/sqlite"
)

func register() {
	config.DBCreaters[`sqlite`] = CreaterSQLite
	config.DBConnecters[`sqlite`] = ConnectSQLite
	config.DBInstallers[`sqlite`] = Exec
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
	factory.SetDebug(c.DB.Debug)
	cluster := factory.NewCluster().AddMaster(database)
	factory.SetCluster(0, cluster).Cluster(0).SetPrefix(c.DB.Prefix)
	return nil
}

func CreaterSQLite(err error, c *config.Config) error {
	if strings.Contains(err.Error(), `unable to open database file`) {
		var f *os.File
		f, err = os.Create(c.DB.Database)
		if err == nil {
			f.Close()
			err = config.ConnectDB()
		}
	}
	return err
}
