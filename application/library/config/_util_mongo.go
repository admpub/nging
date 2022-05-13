//go:build mongo
// +build mongo

package config

import (
	mongo "github.com/webx-top/db/mongoq"
)

func init() {
	DBConnecters[`mongo`] = ConnectMongoDB
}

func ConnectMongoDB(c *Config) error {
	settings := c.DB.ToMongoDB()
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
