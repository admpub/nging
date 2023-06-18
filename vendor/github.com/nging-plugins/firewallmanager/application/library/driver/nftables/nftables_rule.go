package nftables

import (
	"strings"

	"github.com/admpub/nftablesutils"
	setutils "github.com/admpub/nftablesutils/set"
	"github.com/google/nftables"
	"github.com/nging-plugins/firewallmanager/application/library/driver"
	"github.com/webx-top/echo/param"
)

func (a *NFTables) buildProtoRule(rule *driver.Rule) (args nftablesutils.Exprs) {
	switch rule.Protocol {
	case `tcp`:
		args = nftablesutils.JoinExprs(args, nftablesutils.SetProtoTCP())
	case `udp`:
		args = nftablesutils.JoinExprs(args, nftablesutils.SetProtoUDP())
	case `icmp`:
		if a.isIPv4() {
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
	if len(rule.LocalIP) == 0 {
		return
	}
	if strings.Contains(rule.LocalIP, `-`) {
		var ipSet *nftables.Set
		var elems []nftables.SetElement
		var eErr error
		if a.isIPv4() {
			ipSet = nftablesutils.GetIPv4AddrSet(a.NFTables.TableFilter())
			elems, eErr = setutils.GenerateElementsFromIPv4Address([]string{rule.LocalIP})
		} else {
			ipSet = nftablesutils.GetIPv6AddrSet(a.NFTables.TableFilter())
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
		args = args.Add(nftablesutils.SetDAddrSet(ipSet)...)
	} else {
		args = args.Add(nftablesutils.SetCIDRMatcher(nftablesutils.ExprDirectionDestination, rule.LocalIP, false)...)
	}
	return
}

func (a *NFTables) buildRemoteIPRule(c *nftables.Conn, rule *driver.Rule) (args nftablesutils.Exprs, err error) {
	if len(rule.RemoteIP) == 0 {
		return
	}
	if strings.Contains(rule.RemoteIP, `-`) {
		var ipSet *nftables.Set
		var elems []nftables.SetElement
		var eErr error
		if a.isIPv4() {
			ipSet = nftablesutils.GetIPv4AddrSet(a.NFTables.TableFilter())
			elems, eErr = setutils.GenerateElementsFromIPv4Address([]string{rule.RemoteIP})
		} else {
			ipSet = nftablesutils.GetIPv6AddrSet(a.NFTables.TableFilter())
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
		args = args.Add(nftablesutils.SetSAddrSet(ipSet)...)
	} else {
		args = args.Add(nftablesutils.SetCIDRMatcher(nftablesutils.ExprDirectionSource, rule.RemoteIP, false)...)
	}
	return
}

func (a *NFTables) buildLocalPortRule(c *nftables.Conn, rule *driver.Rule) (args nftablesutils.Exprs, err error) {
	if len(rule.LocalPort) == 0 {
		return
	}
	if strings.Contains(rule.LocalPort, `,`) {
		ports := param.Split(rule.LocalPort, `,`).Unique().Uint16(func(_ int, v uint16) bool {
			return nftablesutils.ValidatePort(v) == nil
		})
		if len(ports) > 0 {
			portSet := nftablesutils.GetPortSet(a.NFTables.TableFilter())
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
			args = args.Add(nftablesutils.SetDPortSet(portSet)...)
		}
	} else {
		ports := param.StringSlice(notNumberRegexp.Split(rule.LocalPort, -1)).Unique().Uint16(func(_ int, v uint16) bool {
			return nftablesutils.ValidatePort(v) == nil
		})

		if len(ports) > 0 {
			portsUint16 := make([]uint16, len(ports))
			for k, v := range ports {
				portsUint16[k] = uint16(v)
			}
			if len(portsUint16) >= 2 {
				err = nftablesutils.ValidatePortRange(portsUint16[0], portsUint16[1])
				if err != nil {
					return
				}
				args = args.Add(nftablesutils.SetDPortRange(portsUint16[0], portsUint16[1])...)
			} else {
				args = args.Add(nftablesutils.SetDPort(portsUint16[0])...)
			}
		}
	}
	return
}

func (a *NFTables) buildRemotePortRule(c *nftables.Conn, rule *driver.Rule) (args nftablesutils.Exprs, err error) {
	if len(rule.RemotePort) == 0 {
		return
	}
	if strings.Contains(rule.RemotePort, `,`) {
		ports := param.Split(rule.RemotePort, `,`).Unique().Uint16(func(_ int, v uint16) bool {
			return nftablesutils.ValidatePort(v) == nil
		})
		if len(ports) > 0 {
			portSet := nftablesutils.GetPortSet(a.NFTables.TableFilter())
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
			args = args.Add(nftablesutils.SetSPortSet(portSet)...)
		}
	} else {
		ports := param.StringSlice(notNumberRegexp.Split(rule.RemotePort, -1)).Unique().Uint16(func(_ int, v uint16) bool {
			return nftablesutils.ValidatePort(v) == nil
		})

		if len(ports) > 0 {
			portsUint16 := make([]uint16, len(ports))
			for k, v := range ports {
				portsUint16[k] = uint16(v)
			}
			if len(portsUint16) >= 2 {
				err = nftablesutils.ValidatePortRange(portsUint16[0], portsUint16[1])
				if err != nil {
					return
				}
				args = args.Add(nftablesutils.SetSPortRange(portsUint16[0], portsUint16[1])...)
			} else {
				args = args.Add(nftablesutils.SetSPort(portsUint16[0])...)
			}
		}
	}
	return
}

func (a *NFTables) buildStateRule(c *nftables.Conn, rule *driver.Rule) (args nftablesutils.Exprs, err error) {
	if len(rule.State) == 0 {
		return
	}
	stateSet := nftablesutils.GetConntrackStateSet(a.NFTables.TableFilter())
	states := strings.SplitN(rule.State, ` `, 2) // "target1,target2 allow1,allow2"
	if len(states) != 2 {
		states = strings.Split(rule.State, `,`)
	} else {
		states = strings.Split(states[1], `,`)
	}
	states = param.StringSlice(states).Unique().Filter().String()
	if len(states) == 0 {
		states = []string{nftablesutils.StateNew, nftablesutils.StateEstablished}
	}
	elems := nftablesutils.GetConntrackStateSetElems(states)
	err = c.AddSet(stateSet, elems)
	if err != nil {
		return nil, err
	}
	args = args.Add(nftablesutils.SetConntrackStateSet(stateSet)...)
	return
}
