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

package driver

import (
	"encoding/gob"
	"fmt"
	"net/url"
)

type DbAuth struct {
	Driver       string
	Username     string
	Password     string
	Host         string
	Db           string
	Charset      string
	AccountID    uint
	AccountTitle string
}

func (d *DbAuth) GenKey() string {
	if d.AccountID > 0 {
		return GenKey(``, ``, ``, ``, d.AccountID)
	}
	return GenKey(d.Driver, d.Username, d.Host, d.Db, d.AccountID)
}

func GenKey(driver string, username string, host string, database string, accountID uint) string {
	key := fmt.Sprintf(
		"%s://%s:@%s/%s?%d",
		driver,
		url.QueryEscape(username),
		host,
		database,
		accountID,
	)
	return key
}

func (d *DbAuth) CopyFrom(auth *DbAuth) *DbAuth {
	d.Driver = auth.Driver
	d.Username = auth.Username
	d.Password = auth.Password
	d.Host = auth.Host
	d.Db = auth.Db
	d.Charset = auth.Charset
	d.AccountID = auth.AccountID
	d.AccountTitle = auth.AccountTitle
	return d
}

type AuthAccounts map[string]*DbAuth

func (a *AuthAccounts) Add(account *DbAuth) *AuthAccounts {
	key := account.GenKey()
	(*a)[key] = account
	return a
}

func (a AuthAccounts) Get(key string) *DbAuth {
	if v, y := a[key]; y {
		return v
	}
	return nil
}

func (a *AuthAccounts) Delete(account *DbAuth) {
	key := account.GenKey()
	a.DeleteByKey(key)
}

func (a *AuthAccounts) DeleteByKey(key string) {
	if _, y := (*a)[key]; y {
		delete(*a, key)
	}
}

func init() {
	gob.Register(&DbAuth{})
	gob.Register(AuthAccounts{})
}
