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

//const NgingPluginDir = `../../nging-plugins`
const NgingPluginDir = `vendor/github.com/nging-plugins`

func init() {
	bindata.PathAliases.Add(`caddy`, NgingPluginDir+`/caddymanager/template/backend`)
	bindata.PathAliases.Add(`collector`, NgingPluginDir+`/collector/template/backend`)
	bindata.PathAliases.Add(`db`, NgingPluginDir+`/dbmanager/template/backend`)
	bindata.PathAliases.Add(`ddns`, NgingPluginDir+`/ddnsmanager/template/backend`)
	bindata.PathAliases.Add(`download`, NgingPluginDir+`/dlmanager/template/backend`)
	bindata.PathAliases.Add(`frp`, NgingPluginDir+`/frpmanager/template/backend`)
	bindata.PathAliases.Add(`ftp`, NgingPluginDir+`/ftpmanager/template/backend`)
	bindata.PathAliases.Add(`server`, NgingPluginDir+`/servermanager/template/backend`)
	bindata.PathAliases.Add(`term`, NgingPluginDir+`/sshmanager/template/backend`)

	bindata.Initialize()
}
