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
	"strings"

	"github.com/nging-plugins/firewallmanager/application/library/driver"
	"github.com/nging-plugins/firewallmanager/application/library/enums"
	"github.com/webx-top/com"
)

func appendArgs(to *[]string, from []string) {
	if len(from) == 0 {
		return
	}
	*to = append(*to, from...)
}

func (a *IPTables) buildCommonRule(rule *driver.Rule) (args []string, err error) {
	args, err = a.buildProtoRule(rule)
	if err != nil {
		return
	}
	if com.InSlice(rule.Direction, enums.InputIfaceChainList) && !enums.IsEmptyIface(rule.Interface) {
		args = append(args, `-i`, rule.Interface)
	}

	if com.InSlice(`localIp`, enums.ChainParams[rule.Direction]) {
		_args, _err := a.buildLocalIPRule(rule)
		if _err != nil {
			err = _err
			return
		}
		appendArgs(&args, _args)
	}

	if com.InSlice(`localPort`, enums.ChainParams[rule.Direction]) {
		_args, _err := a.buildLocalPortRule(rule)
		if _err != nil {
			err = _err
			return
		}
		appendArgs(&args, _args)
	}

	if com.InSlice(rule.Direction, enums.OutputIfaceChainList) && !enums.IsEmptyIface(rule.Outerface) {
		args = append(args, `-o`, rule.Outerface)
	}

	if com.InSlice(`remoteIp`, enums.ChainParams[rule.Direction]) {
		_args, _err := a.buildRemoteIPRule(rule)
		if _err != nil {
			err = _err
			return
		}
		appendArgs(&args, _args)
	}

	if com.InSlice(`remotePort`, enums.ChainParams[rule.Direction]) {
		_args, _err := a.buildRemotePortRule(rule)
		if _err != nil {
			err = _err
			return
		}
		appendArgs(&args, _args)
	}

	return
}

func (a *IPTables) buildProtoRule(rule *driver.Rule) (args []string, err error) {
	args = []string{
		`-p`, rule.Protocol,
	}
	return
}

func (a *IPTables) buildLocalIPRule(rule *driver.Rule) (args []string, err error) {
	if enums.IsEmptyIP(rule.LocalIP) {
		return
	}
	if strings.Contains(rule.LocalIP, `-`) {
		args = append(args, `-m`, `iprange`)
		args = append(args, `--dst-range`, rule.LocalIP)
	} else {
		args = append(args, `-d`, rule.LocalIP)
	}
	return
}

func (a *IPTables) buildRemoteIPRule(rule *driver.Rule) (args []string, err error) {
	if enums.IsEmptyIP(rule.RemoteIP) {
		return
	}
	if strings.Contains(rule.RemoteIP, `-`) {
		args = append(args, `-m`, `iprange`)
		args = append(args, `--src-range`, rule.RemoteIP)
	} else {
		args = append(args, `-s`, rule.RemoteIP)
	}
	return
}

func (a *IPTables) buildLocalPortRule(rule *driver.Rule) (args []string, err error) {
	if len(rule.LocalPort) == 0 {
		return
	}
	if strings.Contains(rule.LocalPort, `,`) {
		args = append(args, `-m`, `multiport`)
		args = append(args, `--dports`, rule.LocalPort)
	} else {
		rule.LocalPort = strings.ReplaceAll(rule.LocalPort, `-`, `:`)
		args = append(args, `--dport`, rule.LocalPort)
	}
	return
}

func (a *IPTables) buildRemotePortRule(rule *driver.Rule) (args []string, err error) {
	if len(rule.RemotePort) == 0 {
		return
	}
	if strings.Contains(rule.RemotePort, `,`) {
		args = append(args, `-m`, `multiport`)
		args = append(args, `--sports`, rule.RemotePort)
	} else {
		rule.RemotePort = strings.ReplaceAll(rule.RemotePort, `-`, `:`)
		args = append(args, `--sport`, rule.RemotePort) // 支持用“:”指定端口范围，例如 “22:25” 指端口 22-25，或者 “:22” 指端口 0-22 或者 “22:” 指端口 22-65535
	}
	return
}

func (a *IPTables) buildStateRule(rule *driver.Rule) (args []string, err error) {
	if len(rule.State) == 0 {
		return
	}
	args = append(args, `-m`, `state`)
	args = append(args, `--state`, rule.State)
	return
}
