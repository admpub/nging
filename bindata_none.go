//go:build !bindata
// +build !bindata

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

package main

import (
	"github.com/admpub/nging/v4/application/library/bindata"
)

func init() {
	bindata.PathAliases.Add(`caddy`, `../../nging-plugins/caddymanager/template/backend`)
	bindata.PathAliases.Add(`collector`, `../../nging-plugins/collector/template/backend`)
	bindata.PathAliases.Add(`db`, `../../nging-plugins/dbmanager/template/backend`)
	bindata.PathAliases.Add(`ddns`, `../../nging-plugins/ddnsmanager/template/backend`)
	bindata.PathAliases.Add(`download`, `../../nging-plugins/dlmanager/template/backend`)
	bindata.PathAliases.Add(`frp`, `../../nging-plugins/frpmanager/template/backend`)
	bindata.PathAliases.Add(`ftp`, `../../nging-plugins/ftpmanager/template/backend`)
	bindata.PathAliases.Add(`server`, `../../nging-plugins/servermanager/template/backend`)
	bindata.PathAliases.Add(`term`, `../../nging-plugins/sshmanager/template/backend`)

	bindata.Initialize()
}
