//go:build mongo
// +build mongo

package sdb

import (
	mongo "github.com/webx-top/db/mongoq"
)

func (d *DB) ToMongoDB() mongo.ConnectionURL {
	settings := mongo.ConnectionURL{
		Host:     d.Host,
		Database: d.Database,
		User:     d.User,
		Password: d.Password,
		Options:  d.Options,
	}
	if d.ConnMaxDuration() > 0 {
		mongoq.ConnTimeout = d.ConnMaxDuration()
	}
	return settings
}
