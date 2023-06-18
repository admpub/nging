package nftables

import (
	"github.com/admpub/nftablesutils"
	"github.com/google/nftables"
	"github.com/google/nftables/expr"
	"github.com/nging-plugins/firewallmanager/application/library/driver"
	"golang.org/x/sys/unix"
)

func (a *NFTables) ruleFilterFrom(c *nftables.Conn, rule *driver.Rule) (args nftablesutils.Exprs, err error) {
	args = args.Add(a.buildProtoRule(rule)...)
	if len(rule.Interface) > 0 {
		args = args.Add(nftablesutils.SetIIF(rule.Interface)...) // 只能用于 PREROUTING、INPUT、FORWARD
	} else if len(rule.Outerface) > 0 {
		args = args.Add(nftablesutils.SetOIF(rule.Outerface)...) // 只能用于 FORWARD、OUTPUT、POSTROUTING
	}
	if len(rule.RemoteIP) > 0 {
		_args, _err := a.buildRemoteIPRule(c, rule)
		if _err != nil {
			return nil, _err
		}
		args = args.Add(_args...)
	} else if len(rule.LocalIP) > 0 {
		_args, _err := a.buildLocalIPRule(c, rule)
		if _err != nil {
			return nil, _err
		}
		args = args.Add(_args...)
	}
	if len(rule.RemotePort) > 0 {
		_args, _err := a.buildRemotePortRule(c, rule)
		if _err != nil {
			return nil, _err
		}
		args = args.Add(_args...)
	} else if len(rule.LocalPort) > 0 {
		_args, _err := a.buildLocalPortRule(c, rule)
		if _err != nil {
			return nil, _err
		}
		args = args.Add(_args...)
	}
	if len(rule.State) > 0 {
		_args, _err := a.buildStateRule(c, rule)
		if _err != nil {
			return nil, _err
		}
		args = args.Add(_args...)
	}
	switch rule.Action {
	case `accept`, `ACCEPT`:
		args = args.Add(nftablesutils.Accept())
	case `drop`, `DROP`:
		args = args.Add(nftablesutils.ExprCounter())
		args = args.Add(nftablesutils.Drop())
	case `reject`, `REJECT`:
		args = args.Add(nftablesutils.ExprCounter())
		args = args.Add(nftablesutils.Reject())
	case `log`, `LOG`:
		args = args.Add(&expr.Log{
			Level: expr.LogLevelAlert,
			Flags: expr.LogFlagsNFLog, //expr.LogFlagsIPOpt | expr.LogFlagsTCPOpt,
			Key:   1 << unix.NFTA_LOG_PREFIX,
			Data:  []byte(`nging_`),
		})
	default:
		args = args.Add(nftablesutils.ExprCounter())
		args = args.Add(nftablesutils.Drop())
	}
	return args, nil
}
