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

package handler

import (
	"sync/atomic"
	"time"

	"github.com/webx-top/echo"

	"github.com/admpub/nging/v5/application/library/config/startup"
	"github.com/admpub/nging/v5/application/library/route"
)

func RegisterRoute(r *route.Collection) {
	r.Backend.RegisterToGroup(`/firewall`, registerRoute)
}

var routeRegisters route.Registers

func registerRoute(g echo.RouteRegister) {
	ruleG := g.Group(`/rule`)
	ruleG.Route(`GET,POST`, `/static`, ruleStaticIndex)
	ruleG.Route(`GET,POST`, `/static_add`, ruleStaticAdd)
	ruleG.Route(`GET,POST`, `/static_edit`, ruleStaticEdit)
	ruleG.Route(`GET,POST`, `/static_delete`, ruleStaticDelete)
	ruleG.Route(`GET,POST`, `/static_apply`, ruleStaticApply)
	ruleG.Route(`GET,POST`, `/dynamic`, ruleDynamicIndex)
	ruleG.Route(`GET,POST`, `/dynamic_add`, ruleDynamicAdd)
	ruleG.Route(`GET,POST`, `/dynamic_edit`, ruleDynamicEdit)
	ruleG.Route(`GET,POST`, `/dynamic_delete`, ruleDynamicDelete)

	serviceG := g.Group(`/service`)
	serviceG.Route(`GET,POST`, `/restart`, Restart)
	serviceG.Route(`GET,POST`, `/stop`, Stop)
	serviceG.Route(`GET,POST`, `/log`, Log)

	routeRegisters.Apply(g)
}

var staticRuleLastModifyTs uint64

func setStaticRuleLastModifyTime(t time.Time) {
	atomic.StoreUint64(&staticRuleLastModifyTs, uint64(t.Unix()))
}

func getStaticRuleLastModifyTs() uint64 {
	return atomic.LoadUint64(&staticRuleLastModifyTs)
}
func init() {
	startup.OnAfter(`web.installed`, func() {
	})
	startup.OnAfter(`web`, func() {
	})
}
