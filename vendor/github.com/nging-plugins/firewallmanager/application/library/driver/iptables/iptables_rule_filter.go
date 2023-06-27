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

package iptables

import (
	"github.com/webx-top/com"

	"github.com/nging-plugins/firewallmanager/application/library/driver"
	"github.com/nging-plugins/firewallmanager/application/library/enums"
)

func (a *IPTables) ruleFilterFrom(rule *driver.Rule) (args []string, err error) {
	args, err = a.buildCommonRule(rule)
	if err != nil {
		return
	}

	if com.InSlice(`state`, enums.ChainParams[rule.Direction]) {
		_args, _err := a.buildStateRule(rule)
		if _err != nil {
			err = _err
			return
		}
		appendArgs(&args, _args)
	}

	if com.InSlice(`connLimit`, enums.ChainParams[rule.Direction]) {
		_args, _err := a.buildConnLimitRule(rule)
		if _err != nil {
			err = _err
			return
		}
		appendArgs(&args, _args)
	}

	if com.InSlice(`rateLimit`, enums.ChainParams[rule.Direction]) {
		_args, _err := a.buildHashLimitRule(rule)
		if _err != nil {
			err = _err
			return
		}
		appendArgs(&args, _args)
	}

	args = append(args, `-j`, rule.Action)
	return
}
