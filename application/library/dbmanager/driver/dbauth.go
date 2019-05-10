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

import "encoding/gob"

type DbAuth struct {
	Driver   string
	Username string
	Password string
	Host     string
	Db       string
}

func (d *DbAuth) CopyFrom(auth *DbAuth) *DbAuth {
	d.Driver = auth.Driver
	d.Username = auth.Username
	d.Password = auth.Password
	d.Host = auth.Host
	d.Db = auth.Db
	return d
}

func init() {
	gob.Register(&DbAuth{})
}
