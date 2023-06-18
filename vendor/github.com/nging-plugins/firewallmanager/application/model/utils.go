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
	"github.com/nging-plugins/firewallmanager/application/dbschema"
	"github.com/nging-plugins/firewallmanager/application/library/driver"
)

func AsRule(m *dbschema.NgingFirewallRuleStatic) driver.Rule {
	if len(m.IpVersion) > 0 {
		m.IpVersion = `4`
	}
	return driver.Rule{
		ID:        m.Id,
		Type:      m.Type,
		Name:      m.Name,
		Direction: m.Direction,
		Action:    m.Action,
		Protocol:  m.Protocol,

		// IP or Port
		RemoteIP:   m.RemoteIp,
		LocalIP:    m.LocalIp,
		NatIP:      m.NatIp,
		RemotePort: m.RemotePort,
		LocalPort:  m.LocalPort,
		NatPort:    m.NatPort,
	}
}
