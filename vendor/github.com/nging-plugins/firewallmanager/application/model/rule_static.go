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
	"github.com/admpub/nging/v5/application/library/common"
	"github.com/webx-top/com"
	"github.com/webx-top/db"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/code"

	"github.com/nging-plugins/firewallmanager/application/dbschema"
	"github.com/nging-plugins/firewallmanager/application/library/driver"
	"github.com/nging-plugins/firewallmanager/application/library/enums"
	"github.com/nging-plugins/firewallmanager/application/library/netutils"
)

func NewRuleStatic(ctx echo.Context) *RuleStatic {
	return &RuleStatic{
		NgingFirewallRuleStatic: dbschema.NewNgingFirewallRuleStatic(ctx),
	}
}

type RuleStatic struct {
	*dbschema.NgingFirewallRuleStatic
}

func (r *RuleStatic) check() error {
	ctx := r.Context()
	if !enums.Types.Has(r.Type) {
		return ctx.NewError(code.InvalidParameter, `类型无效`).SetZone(`type`)
	}
	if !com.InSlice(r.Direction, enums.TablesChains[r.Type]) {
		return ctx.NewError(code.InvalidParameter, `类型“%v”不支持设置“%s”规则`, r.Type, ctx.T(enums.Directions.Get(r.Direction))).SetZone(`direction`)
	}
	if !enums.IPProtocols.Has(r.IpVersion) {
		return ctx.NewError(code.InvalidParameter, `IP版本值“%s”无效`, r.IpVersion).SetZone(`ipVersion`)
	}
	if len(r.State) > 0 && !com.InSlice(r.State, enums.StateList) {
		return ctx.NewError(code.InvalidParameter, `网络连接状态值“%s”无效`, r.State).SetZone(`state`)
	}
	if r.Type != enums.TableNAT {
		if !enums.Actions.Has(r.Action) {
			return ctx.NewError(code.InvalidParameter, `操作值“%s”无效`, r.Action).SetZone(`action`)
		}
		if len(r.NatIp) > 0 {
			r.NatIp = ``
		}
		if len(r.NatPort) > 0 {
			r.NatPort = ``
		}
	} else {
		if len(r.Action) > 0 && !enums.Actions.Has(r.Action) {
			return ctx.NewError(code.InvalidParameter, `操作值“%s”无效`, r.Action).SetZone(`action`)
		}
		switch r.Direction {
		case enums.ChainPreRouting:
			if len(r.NatIp) == 0 && len(r.NatPort) == 0 {
				return ctx.NewError(code.InvalidParameter, `NAT IP 和 NAT 端口 不能同时为空`).SetZone(`natPort`)
			}
		case enums.ChainPostRouting:
		}
	}
	if (len(r.Protocol) == 0 || r.Protocol == enums.ProtocolAll) && (len(r.LocalPort) > 0 || len(r.RemotePort) > 0 || len(r.NatPort) > 0) {
		return ctx.NewError(code.InvalidParameter, `当指定了端口时，必须明确的指定网络协议`).SetZone(`protocol`)
	}
	if len(r.LocalPort) > 0 {
		if err := netutils.ValidatePort(ctx, r.LocalPort); err != nil {
			return ctx.NewError(code.InvalidParameter, `本机%v`, err.Error()).SetZone(`localPort`)
		}
	}
	if len(r.RemotePort) > 0 {
		if err := netutils.ValidatePort(ctx, r.RemotePort); err != nil {
			return ctx.NewError(code.InvalidParameter, `远程%v`, err.Error()).SetZone(`remotePort`)
		}
	}
	if len(r.NatPort) > 0 {
		if err := netutils.ValidatePort(ctx, r.NatPort); err != nil {
			return ctx.NewError(code.InvalidParameter, `NAT %v`, err.Error()).SetZone(`natPort`)
		}
	}
	if len(r.LocalIp) > 0 {
		if err := netutils.ValidateIP(ctx, r.LocalIp); err != nil {
			return ctx.NewError(code.InvalidParameter, `本机%v`, err.Error()).SetZone(`localIp`)
		}
	}
	if len(r.RemoteIp) > 0 {
		if err := netutils.ValidateIP(ctx, r.RemoteIp); err != nil {
			return ctx.NewError(code.InvalidParameter, `远程%v`, err.Error()).SetZone(`remoteIp`)
		}
	}
	if len(r.NatIp) > 0 {
		if err := netutils.ValidateIP(ctx, r.NatIp); err != nil {
			return ctx.NewError(code.InvalidParameter, `NAT %v`, err.Error()).SetZone(`natIp`)
		}
	}
	r.Disabled = common.GetBoolFlag(r.Disabled)
	return nil
}

func (r *RuleStatic) Add() (interface{}, error) {
	if err := r.check(); err != nil {
		return nil, err
	}
	return r.NgingFirewallRuleStatic.Insert()
}

func (r *RuleStatic) Edit(mw func(db.Result) db.Result, args ...interface{}) error {
	if err := r.check(); err != nil {
		return err
	}
	return r.NgingFirewallRuleStatic.Update(mw, args...)
}

func (r *RuleStatic) ListPage(cond *db.Compounds, sorts ...interface{}) ([]*dbschema.NgingFirewallRuleStatic, error) {
	err := r.NgingFirewallRuleStatic.ListPage(cond, sorts...)
	if err != nil {
		return nil, err
	}
	return r.Objects(), nil
}

func (r *RuleStatic) AsRule(row ...*dbschema.NgingFirewallRuleStatic) driver.Rule {
	m := r.NgingFirewallRuleStatic
	if len(row) > 0 && row[0] != nil {
		m = row[0]
	}
	return AsRule(m)
}
