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

package model

import (
	"net/url"
	"strings"

	"github.com/webx-top/db"
	"github.com/webx-top/echo"

	"github.com/nging-plugins/dbmanager/pkg/dbschema"
)

type DbSyncWithAccount struct {
	*dbschema.NgingDbSync
	SrcAccount *dbschema.NgingDbAccount `db:"-,relation=id:source_account_id"`
	DstAccount *dbschema.NgingDbAccount `db:"-,relation=id:destination_account_id"`
}

func NewDbSync(ctx echo.Context) *DbSync {
	return &DbSync{
		NgingDbSync: dbschema.NewNgingDbSync(ctx),
	}
}

type DbSync struct {
	*dbschema.NgingDbSync
}

func (a *DbSync) ToDSNFromAccount(acc *dbschema.NgingDbAccount) string {
	return a.ToDSN(acc.User, acc.Password, acc.Host, acc.Name)
}

func (a *DbSync) ToDSN(user, pwd, host, db string) string {
	//test:test@(127.0.0.1:3306)/test_0
	if len(user) == 0 {
		user = `root`
	}
	if len(pwd) == 0 {
		pwd = `root`
	}
	if len(host) == 0 {
		host = `127.0.0.1:3306`
	}
	return url.QueryEscape(user) + `:` + url.QueryEscape(pwd) + `@(` + host + `)/` + db
}

func (a *DbSync) ParseDSN(dsn string) (user string, pwd string, host string, dbName string) {
	if len(dsn) == 0 {
		return
	}
	idx := strings.Index(dsn, `:`)
	if idx > 0 {
		user, _ = url.QueryUnescape(dsn[0:idx])
	}

	dsn = dsn[idx+1:]
	idx = strings.Index(dsn, `@`)
	if idx > 0 {
		pwd, _ = url.QueryUnescape(dsn[0:idx])
	}

	dsn = dsn[idx+1:]
	idx = strings.Index(dsn, `/`)
	if idx > 0 {
		host = dsn[0:idx]
	}

	host = strings.TrimPrefix(host, `(`)
	host = strings.TrimSuffix(host, `)`)
	if len(dsn) > idx+1 {
		dbName = dsn[idx+1:]
	}
	return
}

func (a *DbSync) HidePassword(dsn string) string {
	if len(dsn) == 0 {
		return dsn
	}
	idx := strings.Index(dsn, `:`)
	var user string
	if idx > 0 {
		user, _ = url.QueryUnescape(dsn[0:idx])
	}

	dsn = dsn[idx+1:]
	idx = strings.Index(dsn, `@`)
	return user + `:***@` + dsn[idx+1:]
}

func (a *DbSync) Add() (interface{}, error) {
	return a.NgingDbSync.Add()
}

func (a *DbSync) Edit(mw func(db.Result) db.Result, args ...interface{}) error {
	return a.NgingDbSync.Edit(mw, args...)
}
