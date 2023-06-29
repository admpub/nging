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
	"github.com/admpub/nftablesutils"
	"github.com/google/nftables"
	"github.com/google/nftables/expr"
	"github.com/webx-top/com"
	"golang.org/x/sys/unix"

	"github.com/nging-plugins/firewallmanager/application/library/driver"
	"github.com/nging-plugins/firewallmanager/application/library/enums"
)

func (a *NFTables) ruleFilterFrom(c *nftables.Conn, rule *driver.Rule) (args nftablesutils.Exprs, err error) {
	args, err = a.buildCommonRule(c, rule)
	if err != nil {
		return
	}
	if com.InSlice(`state`, enums.ChainParams[rule.Direction]) {
		_args, _err := a.buildStateRule(c, rule)
		if _err != nil {
			return nil, _err
		}
		args = args.Add(_args...)
	}

	if com.InSlice(`connLimit`, enums.ChainParams[rule.Direction]) {
		_args, _err := a.buildConnLimitRule(c, rule)
		if _err != nil {
			err = _err
			return
		}
		args = args.Add(_args...)
	}

	if com.InSlice(`rateLimit`, enums.ChainParams[rule.Direction]) {
		var _args []expr.Any
		var _err error
		if rule.Action == enums.TargetAccept {
			_args, _err = a.buildLimitRule(c, rule)
		} else {
			_args, _err = a.buildLimitRuleWithTimeout(c, rule)
		}
		if _err != nil {
			err = _err
			return
		}
		args = args.Add(_args...)
	}

	switch rule.Action {
	case enums.TargetAccept:
		args = args.Add(nftablesutils.ExprCounter())
		args = args.Add(nftablesutils.Accept())
	case enums.TargetDrop:
		args = args.Add(nftablesutils.ExprCounter())
		args = args.Add(nftablesutils.Drop())
	case enums.TargetReject:
		args = args.Add(nftablesutils.ExprCounter())
		args = args.Add(nftablesutils.Reject())
	case enums.TargetLog:
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
