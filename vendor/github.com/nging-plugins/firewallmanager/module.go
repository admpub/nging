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

package firewallmanager

import (
	"github.com/admpub/nging/v5/application/library/config/cmder"
	"github.com/admpub/nging/v5/application/library/config/extend"
	"github.com/admpub/nging/v5/application/library/module"

	"github.com/nging-plugins/firewallmanager/application/handler"
	pluginCmder "github.com/nging-plugins/firewallmanager/application/library/cmder"
	"github.com/nging-plugins/firewallmanager/application/library/setup"
)

const ID = `firewall`

var Module = module.Module{
	Startup: `firewall`,
	Extend: map[string]extend.Initer{
		`firewall`: pluginCmder.Initer,
	},
	Cmder: map[string]cmder.Cmder{
		`firewall`: pluginCmder.New(),
	},
	TemplatePath: map[string]string{
		ID: `firewallmanager/template/backend`,
	},
	AssetsPath:    []string{},
	SQLCollection: setup.RegisterSQL,
	Navigate:      RegisterNavigate,
	Route:         handler.RegisterRoute,
	DBSchemaVer:   0.0004,
}
