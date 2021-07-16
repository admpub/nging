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

package sdb

import (
	"time"

	"github.com/webx-top/db/lib/factory"
	"github.com/webx-top/db/mongo"
	"github.com/webx-top/db/mysql"

	"github.com/admpub/log"
	"github.com/admpub/nging/application/library/common"
	"github.com/admpub/null"
	"github.com/webx-top/db/lib/sqlbuilder"
)

var dbConnMaxDuration = 10 * time.Second

type DB struct {
	Type            string            `json:"type"`
	User            string            `json:"user"`
	Password        string            `json:"password"`
	Host            string            `json:"host"`
	Database        string            `json:"database"`
	Prefix          string            `json:"prefix"`
	Options         map[string]string `json:"options"`
	Debug           bool              `json:"debug"`
	ConnMaxLifetime string            `json:"connMaxLifetime"` //example: 10s
	MaxIdleConns    int               `json:"maxIdleConns"`
	MaxOpenConns    int               `json:"maxOpenConns"`
	connMaxDuration time.Duration
}

func (d *DB) SetKV(key string, value string) *DB {
	if d.Options == nil {
		d.Options = map[string]string{}
	}
	d.Options[key] = value
	return d
}

func (d *DB) GetByKey(key string) (string, bool) {
	if d.Options == nil {
		return ``, false
	}
	value, ok := d.Options[key]
	return value, ok
}

func (d *DB) Charset() string {
	charset, _ := d.GetByKey(`charset`)
	return charset
}

type DBConnSetter interface {
	SetConnMaxLifetime(time.Duration)
	SetMaxIdleConns(int)
	SetMaxOpenConns(int)
}

func (d *DB) SetDebug(on bool) {
	d.Debug = on
	factory.SetDebug(on)
}

func (d *DB) Table(table string) string {
	return d.Prefix + table
}

func (d *DB) ToTable(m sqlbuilder.Name_) string {
	return d.Table(m.Name_())
}

func (d *DB) ConnMaxDuration() time.Duration {
	if d.connMaxDuration > 0 {
		return d.connMaxDuration
	}
	if len(d.ConnMaxLifetime) > 0 {
		d.connMaxDuration, _ = time.ParseDuration(d.ConnMaxLifetime)
		if d.connMaxDuration <= 0 {
			d.connMaxDuration = dbConnMaxDuration
		}
	} else {
		d.connMaxDuration = dbConnMaxDuration
	}
	return d.connMaxDuration
}

func (d *DB) SetConn(setter DBConnSetter) error {
	setter.SetMaxIdleConns(d.MaxIdleConns)
	setter.SetMaxOpenConns(d.MaxOpenConns)
	database, ok := setter.(sqlbuilder.Database)
	if !ok {
		setter.SetConnMaxLifetime(d.ConnMaxDuration())
		return nil
	}
	var retErr error
	switch d.Type {
	case `mysql`:
		rows, err := database.Query(`show variables where Variable_name = 'wait_timeout'`)
		if err != nil {
			log.Error(err)
		} else {
			if rows.Next() {
				name := null.String{}
				timeout := null.String{}
				err = rows.Scan(&name, &timeout)
				if err != nil {
					log.Error(err)
				} else {
					d.connMaxDuration = time.Duration(timeout.Int64()) * time.Second
					d.connMaxDuration /= 2
				}
			}
			rows.Close()
		}
		retErr = err
	default:
	}
	database.SetConnMaxLifetime(d.ConnMaxDuration())
	return retErr
}

func (d *DB) ToMySQL() mysql.ConnectionURL {
	settings := mysql.ConnectionURL{
		Host:     d.Host,
		Database: d.Database,
		User:     d.User,
		Password: d.Password,
		Options:  d.Options,
	}
	common.ParseMysqlConnectionURL(&settings)
	if settings.Options == nil {
		settings.Options = map[string]string{}
	}
	// Default options.
	if _, ok := settings.Options["charset"]; !ok {
		settings.Options["charset"] = "utf8mb4"
	}
	return settings
}

func (d *DB) ToMongoDB() mongo.ConnectionURL {
	settings := mongo.ConnectionURL{
		Host:     d.Host,
		Database: d.Database,
		User:     d.User,
		Password: d.Password,
		Options:  d.Options,
	}
	if d.ConnMaxDuration() > 0 {
		mongo.ConnTimeout = d.ConnMaxDuration()
	}
	return settings
}
