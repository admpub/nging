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
	"strings"
	"time"

	"github.com/admpub/nftablesutils"
	setutils "github.com/admpub/nftablesutils/set"
	"github.com/google/nftables"
	"github.com/google/nftables/expr"
	"github.com/webx-top/com"
	"github.com/webx-top/echo/param"

	"github.com/nging-plugins/firewallmanager/application/library/driver"
	"github.com/nging-plugins/firewallmanager/application/library/enums"
)

func (a *NFTables) buildCommonRule(c *nftables.Conn, rule *driver.Rule) (args nftablesutils.Exprs, err error) {
	args = args.Add(a.buildProtoRule(rule)...)
	if com.InSlice(rule.Direction, enums.InputIfaceChainList) && !enums.IsEmptyIface(rule.Interface) {
		args = args.Add(nftablesutils.SetIIF(rule.Interface)...)
	}

	if com.InSlice(`localIp`, enums.ChainParams[rule.Direction]) {
		_args, _err := a.buildLocalIPRule(c, rule)
		if _err != nil {
			return nil, _err
		}
		args = args.Add(_args...)
	}

	if com.InSlice(`localPort`, enums.ChainParams[rule.Direction]) {
		_args, _err := a.buildLocalPortRule(c, rule)
		if _err != nil {
			return nil, _err
		}
		args = args.Add(_args...)
	}

	if com.InSlice(rule.Direction, enums.OutputIfaceChainList) && !enums.IsEmptyIface(rule.Outerface) {
		args = args.Add(nftablesutils.SetOIF(rule.Outerface)...)
	}

	if com.InSlice(`remoteIp`, enums.ChainParams[rule.Direction]) {
		_args, _err := a.buildRemoteIPRule(c, rule)
		if _err != nil {
			return nil, _err
		}
		args = args.Add(_args...)
	}

	if com.InSlice(`remotePort`, enums.ChainParams[rule.Direction]) {
		_args, _err := a.buildRemotePortRule(c, rule)
		if _err != nil {
			return nil, _err
		}
		args = args.Add(_args...)
	}

	return
}

func (a *NFTables) buildProtoRule(rule *driver.Rule) (args nftablesutils.Exprs) {
	switch rule.Protocol {
	case enums.ProtocolTCP:
		args = nftablesutils.JoinExprs(args, nftablesutils.SetProtoTCP())
	case enums.ProtocolUDP:
		args = nftablesutils.JoinExprs(args, nftablesutils.SetProtoUDP())
	case enums.ProtocolICMP:
		if a.base.isIPv4() {
			args = nftablesutils.JoinExprs(args, nftablesutils.SetProtoICMP())
		} else {
			args = nftablesutils.JoinExprs(args, nftablesutils.SetProtoICMPv6())
		}
	default:
		// all
	}
	return
}

func (a *NFTables) buildLocalIPRule(c *nftables.Conn, rule *driver.Rule) (args nftablesutils.Exprs, err error) {
	if enums.IsEmptyIP(rule.LocalIP) {
		return
	}
	var neq bool
	if strings.HasPrefix(rule.LocalIP, `!`) {
		neq = true
		rule.LocalIP = strings.TrimPrefix(rule.LocalIP, `!`)
	}
	if strings.Contains(rule.LocalIP, `-`) {
		var ipSet *nftables.Set
		var elems []nftables.SetElement
		var eErr error
		if a.base.isIPv4() {
			ipSet = nftablesutils.GetIPv4AddrSet(a.base.TableFilter())
			elems, eErr = setutils.GenerateElementsFromIPv4Address([]string{rule.LocalIP})
		} else {
			ipSet = nftablesutils.GetIPv6AddrSet(a.base.TableFilter())
			elems, eErr = setutils.GenerateElementsFromIPv6Address([]string{rule.LocalIP})
		}
		if eErr != nil {
			err = eErr
			return
		}
		ipSet.Interval = true
		err = c.AddSet(ipSet, elems)
		if err != nil {
			return nil, err
		}
		args = args.Add(nftablesutils.SetDAddrSet(ipSet, !neq)...)
	} else {
		args = args.Add(nftablesutils.SetCIDRMatcher(nftablesutils.ExprDirectionDestination, rule.LocalIP, false, !neq)...)
	}
	return
}

func (a *NFTables) buildRemoteIPRule(c *nftables.Conn, rule *driver.Rule) (args nftablesutils.Exprs, err error) {
	if enums.IsEmptyIP(rule.RemoteIP) {
		return
	}
	var neq bool
	if strings.HasPrefix(rule.RemoteIP, `!`) {
		neq = true
		rule.RemoteIP = strings.TrimPrefix(rule.RemoteIP, `!`)
	}
	if strings.Contains(rule.RemoteIP, `-`) {
		var ipSet *nftables.Set
		var elems []nftables.SetElement
		var eErr error
		if a.base.isIPv4() {
			ipSet = nftablesutils.GetIPv4AddrSet(a.base.TableFilter())
			elems, eErr = setutils.GenerateElementsFromIPv4Address([]string{rule.RemoteIP})
		} else {
			ipSet = nftablesutils.GetIPv6AddrSet(a.base.TableFilter())
			elems, eErr = setutils.GenerateElementsFromIPv6Address([]string{rule.RemoteIP})
		}
		if eErr != nil {
			err = eErr
			return
		}
		ipSet.Interval = true
		err = c.AddSet(ipSet, elems)
		if err != nil {
			return nil, err
		}
		args = args.Add(nftablesutils.SetSAddrSet(ipSet, !neq)...)
	} else {
		args = args.Add(nftablesutils.SetCIDRMatcher(nftablesutils.ExprDirectionSource, rule.RemoteIP, false, !neq)...)
	}
	return
}

func (a *NFTables) buildLocalPortRule(c *nftables.Conn, rule *driver.Rule) (args nftablesutils.Exprs, err error) {
	if enums.IsEmptyPort(rule.LocalPort) {
		return
	}
	var neq bool
	if strings.HasPrefix(rule.LocalPort, `!`) {
		neq = true
		rule.LocalPort = strings.TrimPrefix(rule.LocalPort, `!`)
	}
	if strings.Contains(rule.LocalPort, `,`) {
		ports := param.Split(rule.LocalPort, `,`).Unique().Uint16(func(_ int, v uint16) bool {
			return nftablesutils.ValidatePort(v) == nil
		})
		if len(ports) > 0 {
			portSet := nftablesutils.GetPortSet(a.base.TableFilter())
			portsUint16 := make([]uint16, len(ports))
			for k, v := range ports {
				portsUint16[k] = uint16(v)
			}
			elems := nftablesutils.GetPortElems(portsUint16)
			//portSet.Interval = true
			err = c.AddSet(portSet, elems)
			if err != nil {
				return nil, err
			}
			args = args.Add(nftablesutils.SetDPortSet(portSet, !neq)...)
		}
	} else {
		ss := param.StringSlice(notNumberRegexp.Split(rule.LocalPort, -1)).Unique()
		ports := ss.Uint16(func(_ int, v uint16) bool {
			return nftablesutils.ValidatePort(v) == nil
		})

		if len(ports) > 0 {
			portsUint16 := make([]uint16, len(ports))
			portsString := make([]string, len(ports))
			for k, v := range ports {
				portsUint16[k] = uint16(v)
				portsString[k] = param.AsString(v)
			}
			if len(portsUint16) >= 2 {
				err = nftablesutils.ValidatePortRange(portsUint16[0], portsUint16[1])
				if err != nil {
					return
				}
				if neq {
					portSet := nftablesutils.GetPortSet(a.base.TableFilter())
					portSet.Interval = true
					elems, err := setutils.GenerateElementsFromPort([]string{portsString[0] + `-` + portsString[1]})
					if err != nil {
						return nil, err
					}
					err = c.AddSet(portSet, elems)
					if err != nil {
						return nil, err
					}
					args = args.Add(nftablesutils.SetDPortSet(portSet, !neq)...)
				} else {
					args = args.Add(nftablesutils.SetDPortRange(portsUint16[0], portsUint16[1])...)
				}
			} else {
				args = args.Add(nftablesutils.SetDPort(portsUint16[0], !neq)...)
			}
		}
	}
	return
}

func (a *NFTables) buildRemotePortRule(c *nftables.Conn, rule *driver.Rule) (args nftablesutils.Exprs, err error) {
	if enums.IsEmptyPort(rule.RemotePort) {
		return
	}
	var neq bool
	if strings.HasPrefix(rule.RemotePort, `!`) {
		neq = true
		rule.RemotePort = strings.TrimPrefix(rule.RemotePort, `!`)
	}
	if strings.Contains(rule.RemotePort, `,`) {
		ports := param.Split(rule.RemotePort, `,`).Unique().Uint16(func(_ int, v uint16) bool {
			return nftablesutils.ValidatePort(v) == nil
		})
		if len(ports) > 0 {
			portSet := nftablesutils.GetPortSet(a.base.TableFilter())
			portsUint16 := make([]uint16, len(ports))
			for k, v := range ports {
				portsUint16[k] = uint16(v)
			}
			elems := nftablesutils.GetPortElems(portsUint16)
			//portSet.Interval = true
			err = c.AddSet(portSet, elems)
			if err != nil {
				return nil, err
			}
			args = args.Add(nftablesutils.SetSPortSet(portSet, !neq)...)
		}
	} else {
		ports := param.StringSlice(notNumberRegexp.Split(rule.RemotePort, -1)).Unique().Uint16(func(_ int, v uint16) bool {
			return nftablesutils.ValidatePort(v) == nil
		})

		if len(ports) > 0 {
			portsUint16 := make([]uint16, len(ports))
			portsString := make([]string, len(ports))
			for k, v := range ports {
				portsUint16[k] = uint16(v)
				portsString[k] = param.AsString(v)
			}
			if len(portsUint16) >= 2 {
				err = nftablesutils.ValidatePortRange(portsUint16[0], portsUint16[1])
				if err != nil {
					return
				}
				if neq {
					portSet := nftablesutils.GetPortSet(a.base.TableFilter())
					portSet.Interval = true
					elems, err := setutils.GenerateElementsFromPort([]string{portsString[0] + `-` + portsString[1]})
					if err != nil {
						return nil, err
					}
					err = c.AddSet(portSet, elems)
					if err != nil {
						return nil, err
					}
					args = args.Add(nftablesutils.SetSPortSet(portSet, !neq)...)
				} else {
					args = args.Add(nftablesutils.SetSPortRange(portsUint16[0], portsUint16[1])...)
				}
			} else {
				args = args.Add(nftablesutils.SetSPort(portsUint16[0], !neq)...)
			}
		}
	}
	return
}

func (a *NFTables) buildStateRule(c *nftables.Conn, rule *driver.Rule) (args nftablesutils.Exprs, err error) {
	if len(rule.State) == 0 {
		return
	}
	stateSet := nftablesutils.GetConntrackStateSet(a.base.TableFilter())
	states := strings.Split(rule.State, `,`)
	states = param.StringSlice(states).Unique().Filter().String()
	if len(states) == 0 {
		states = []string{nftablesutils.StateNew, nftablesutils.StateEstablished}
	} else {
		for index, state := range states {
			states[index] = strings.ToLower(state)
		}
	}
	elems := nftablesutils.GetConntrackStateSetElems(states)
	err = c.AddSet(stateSet, elems)
	if err != nil {
		return nil, err
	}
	args = args.Add(nftablesutils.SetConntrackStateSet(stateSet)...)
	return
}

func (a *NFTables) buildConnLimitRule(_ *nftables.Conn, rule *driver.Rule) (args nftablesutils.Exprs, err error) {
	if len(rule.ConnLimit) == 0 {
		return
	}
	var m *expr.Connlimit
	m, err = nftablesutils.ParseConnLimit(rule.ConnLimit)
	if err != nil {
		return
	}
	args = args.Add(m)
	return
}

func (a *NFTables) buildLimitRule(c *nftables.Conn, rule *driver.Rule) (args nftablesutils.Exprs, err error) {
	if len(rule.RateLimit) == 0 {
		setName := LimitSetNamePrefix + param.AsString(rule.ID)
		existSet, existErr := c.GetSetByName(a.base.TableFilter(), setName)
		if existErr != nil {
			return
		}
		c.DelSet(existSet)
		return
	}
	var exp *expr.Limit
	exp, err = nftablesutils.ParseLimits(rule.RateLimit, uint32(rule.RateBurst))
	if err != nil {
		return
	}
	args = args.Add(exp)
	return
}

func (a *NFTables) buildLimitRuleWithTimeout(c *nftables.Conn, rule *driver.Rule) (args nftablesutils.Exprs, err error) {
	setName := LimitSetNamePrefix + param.AsString(rule.ID)
	if len(rule.RateLimit) == 0 {
		existSet, existErr := c.GetSetByName(a.base.TableFilter(), setName)
		if existErr != nil {
			return
		}
		c.DelSet(existSet)
		return
	}
	var set *nftables.Set
	if a.base.isIPv4() {
		set = nftablesutils.GetIPv4AddrSet(a.base.TableFilter())
	} else {
		set = nftablesutils.GetIPv6AddrSet(a.base.TableFilter())
	}
	set.Anonymous = false
	set.Constant = false
	set.Dynamic = true
	set.HasTimeout = true
	set.Timeout = time.Duration(rule.RateExpires) * time.Second
	set.Name = setName

	var existSet *nftables.Set
	existSet, err = c.GetSetByName(a.base.TableFilter(), set.Name)
	if err == nil {
		if existSet.Timeout != set.Timeout {
			c.DelSet(existSet)
		} else {
			goto END
		}
	}
	err = c.AddSet(set, []nftables.SetElement{})
	if err != nil {
		return
	}

END:
	var exprs []expr.Any
	exprs, err = nftablesutils.SetDynamicLimitSet(set, rule.RateLimit, uint32(rule.RateBurst))
	if err != nil {
		return
	}
	args = args.Add(exprs...)
	return
}
