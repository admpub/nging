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

func (a *NFTables) ruleNATFrom(c *nftables.Conn, rule *driver.Rule) (args nftablesutils.Exprs, err error) {
	args, err = a.buildCommonRule(c, rule)
	if err != nil {
		return
	}
	switch rule.Direction {
	case enums.ChainPreRouting:
		if len(rule.NatPort) > 0 {
			port := param.AsUint16(rule.NatPort)
			err = nftablesutils.ValidatePort(port)
			if err != nil {
				return
			}
			args = args.Add(nftablesutils.RedirectTo(port)...)
			return
		}
		if len(rule.NatIP) > 0 {
			localIP := strings.SplitN(rule.LocalIP, `-`, 2)[0]
			ip := net.ParseIP(localIP)
			if a.isIPv4() {
				if ip == nil || ip.To4() == nil {
					err = fmt.Errorf(`%w: %s`, driver.ErrInvalidIPv4, localIP)
					return
				}
				args = args.Add(nftablesutils.DNAT(ip)...)
			} else {
				if ip == nil || ip.To4() != nil {
					err = fmt.Errorf(`%w: %s`, driver.ErrInvalidIPv6, localIP)
					return
				}
				args = args.Add(nftablesutils.DNATv6(ip)...)
			}
		} else {
			err = driver.ErrNatIPOrNatPortRequired
		}
	case enums.ChainPostRouting:
		if len(rule.NatIP) > 0 { // 发送给访客
			remoteIP := strings.SplitN(rule.NatIP, `-`, 2)[0]
			ip := net.ParseIP(remoteIP)
			if a.isIPv4() {
				if ip == nil || ip.To4() == nil {
					err = fmt.Errorf(`%w: %s`, driver.ErrInvalidIPv4, remoteIP)
					return
				}
				args = args.Add(nftablesutils.SNAT(ip)...)
			} else {
				if ip == nil || ip.To4() != nil {
					err = fmt.Errorf(`%w: %s`, driver.ErrInvalidIPv6, remoteIP)
					return
				}
				args = args.Add(nftablesutils.SNATv6(ip)...)
			}
		} else {
			args = args.Add(nftablesutils.ExprMasquerade(1, 0))
		}
	default:
		err = fmt.Errorf(`%w: %s (table=%v)`, driver.ErrUnsupportedChain, rule.Direction, rule.Type)
	}
	return
}
