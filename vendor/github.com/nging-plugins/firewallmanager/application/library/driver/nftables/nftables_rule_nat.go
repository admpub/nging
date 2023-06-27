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

package nftables

import (
	"fmt"
	"net"
	"strings"

	"github.com/admpub/nftablesutils"
	"github.com/google/nftables"
	"github.com/nging-plugins/firewallmanager/application/library/driver"
	"github.com/nging-plugins/firewallmanager/application/library/enums"
	"github.com/webx-top/echo/param"
)

func parseIPRange(ipStr string, isIPv4 bool) (ipStart net.IP, ipEnd net.IP, err error) {
	parts := strings.SplitN(ipStr, `-`, 2)
	ip1 := parts[0]
	ipStart = net.ParseIP(ip1)
	if len(parts) == 2 {
		ipEnd = net.ParseIP(parts[1])
	}
	if isIPv4 {
		if ipStart != nil {
			ipStart = ipStart.To4()
		}
		if ipStart == nil {
			err = fmt.Errorf(`%w: %s`, driver.ErrInvalidIPv4, ip1)
			return
		}
		if ipEnd != nil {
			ipEnd = ipEnd.To4()
		}
	} else {
		if ipStart != nil {
			ipStart = ipStart.To16()
		}
		if ipStart == nil {
			err = fmt.Errorf(`%w: %s`, driver.ErrInvalidIPv6, ip1)
			return
		}
		if ipEnd != nil {
			ipEnd = ipEnd.To16()
		}
	}
	return
}

func parsePortRange(portStr string) (portMin uint16, portMax uint16, err error) {
	parts := strings.SplitN(portStr, `-`, 2)
	portMin = param.AsUint16(parts[0])
	err = nftablesutils.ValidatePort(portMin)
	if err != nil {
		return
	}
	if len(parts) == 2 {
		portMax = param.AsUint16(parts[1])
		err = nftablesutils.ValidatePort(portMax)
		if err != nil {
			return
		}
		err = nftablesutils.ValidatePortRange(portMin, portMax)
		if err != nil {
			return
		}
	}
	return
}

func (a *NFTables) ruleNATFrom(c *nftables.Conn, rule *driver.Rule) (args nftablesutils.Exprs, err error) {
	args, err = a.buildCommonRule(c, rule)
	if err != nil {
		return
	}
	switch rule.Direction {
	case enums.ChainPreRouting:
		var portMin, portMax uint16
		portMin, portMax, err = parsePortRange(rule.NatPort)
		if err != nil {
			return
		}
		if len(rule.NatIP) > 0 {
			var ipStart, ipEnd net.IP
			ipStart, ipEnd, err = parseIPRange(rule.NatIP, a.base.isIPv4())
			if err != nil {
				return
			}
			if a.base.isIPv4() {
				args = args.Add(nftablesutils.SetDNATRange(ipStart, ipEnd, portMin, portMax)...)
			} else {
				args = args.Add(nftablesutils.SetDNATv6Range(ipStart, ipEnd, portMin, portMax)...)
			}
		} else if portMin > 0 {
			args = args.Add(nftablesutils.SetRedirect(portMin, portMax)...)
		} else {
			err = driver.ErrNatIPOrNatPortRequired
		}
	case enums.ChainPostRouting:
		var portMin, portMax uint16
		portMin, portMax, err = parsePortRange(rule.NatPort)
		if err != nil {
			return
		}
		if len(rule.NatIP) > 0 { // 发送给访客
			var ipStart, ipEnd net.IP
			ipStart, ipEnd, err = parseIPRange(rule.NatIP, a.base.isIPv4())
			if err != nil {
				return
			}
			if a.base.isIPv4() {
				args = args.Add(nftablesutils.SetSNATRange(ipStart, ipEnd, portMin, portMax)...)
			} else {
				args = args.Add(nftablesutils.SetSNATv6Range(ipStart, ipEnd, portMin, portMax)...)
			}
		} else {
			args = args.Add(nftablesutils.ExprMasquerade(1, 0))
		}
	default:
		err = fmt.Errorf(`%w: %s (table=%v)`, driver.ErrUnsupportedChain, rule.Direction, rule.Type)
	}
	return
}
