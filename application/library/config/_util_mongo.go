//go:build mongo
// +build mongo

package config

import (
	"github.com/webx-top/db/lib/sqlbuilder"
	mongo "github.com/webx-top/db/mongoq"
)

func init() {
	DBConnecters[`mongo`] = ConnectMongoDB
}

func ConnectMongoDB(c sdb.DB) (sqlbuilder.Database, error) {
	settings := c.ToMongoDB()
	return mongo.Open(settings)
}
