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

	"github.com/admpub/nging/application/dbschema"
	"github.com/admpub/nging/application/model/base"
	"github.com/webx-top/db"
	"github.com/webx-top/echo"
)

func NewDbSync(ctx echo.Context) *DbSync {
	return &DbSync{
		DbSync: &dbschema.DbSync{},
		Base:   base.New(ctx),
	}
}

type DbSync struct {
	*dbschema.DbSync
	*base.Base
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
	idx := strings.Index(dsn, `:`)
	user, _ = url.QueryUnescape(dsn[0:idx])

	dsn = dsn[idx+1:]
	idx = strings.Index(dsn, `@`)
	pwd, _ = url.QueryUnescape(dsn[0:idx])

	dsn = dsn[idx+1:]
	idx = strings.Index(dsn, `/`)
	host = dsn[0:idx]

	host = strings.TrimPrefix(host, `(`)
	host = strings.TrimSuffix(host, `)`)

	dbName = dsn[idx+1:]
	return
}

func (a *DbSync) HidePassword(dsn string) string {
	idx := strings.Index(dsn, `:`)
	user, _ := url.QueryUnescape(dsn[0:idx])

	dsn = dsn[idx+1:]
	idx = strings.Index(dsn, `@`)
	return user + `:***@` + dsn[idx+1:]
}

func (a *DbSync) Add() (interface{}, error) {
	return a.DbSync.Add()
}

func (a *DbSync) Edit(mw func(db.Result) db.Result, args ...interface{}) error {
	return a.DbSync.Edit(mw, args...)
}
