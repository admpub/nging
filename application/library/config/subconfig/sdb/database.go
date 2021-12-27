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

package sdb

import (
	"time"

	"github.com/admpub/nging/v4/application/library/common"
	"github.com/webx-top/db/lib/factory"
	"github.com/webx-top/db/lib/sqlbuilder"
	"github.com/webx-top/db/mongo"
	"github.com/webx-top/db/mysql"
)

var (
	MySQLSupportCharsetList = []string{
		MySQLDefaultCharset,
		`utf8`,
	}
)

const MySQLDefaultCharset = `utf8mb4`

type DB struct {
	Type              string            `json:"type"`
	User              string            `json:"user"`
	Password          string            `json:"password"`
	Host              string            `json:"host"`
	Database          string            `json:"database"`
	Prefix            string            `json:"prefix"`
	Options           map[string]string `json:"options"`
	Debug             bool              `json:"debug"`
	ConnMaxLifetime   string            `json:"connMaxLifetime"` //example: 10s
	MaxIdleConns      int               `json:"maxIdleConns"`
	MaxOpenConns      int               `json:"maxOpenConns"`
	connMaxDuration   time.Duration
	parsedMaxDuration bool
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
	if d.parsedMaxDuration {
		return d.connMaxDuration
	}
	if len(d.ConnMaxLifetime) > 0 {
		d.connMaxDuration, _ = time.ParseDuration(d.ConnMaxLifetime)
		d.parsedMaxDuration = true
	}
	return d.connMaxDuration
}

func (d *DB) SetConn(setter DBConnSetter) error {
	setter.SetMaxIdleConns(d.MaxIdleConns)
	setter.SetMaxOpenConns(d.MaxOpenConns)
	connMaxLifetime := d.ConnMaxDuration()
	if connMaxLifetime > 0 {
		setter.SetConnMaxLifetime(connMaxLifetime)
	}
	return nil
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
		settings.Options["charset"] = MySQLDefaultCharset
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
