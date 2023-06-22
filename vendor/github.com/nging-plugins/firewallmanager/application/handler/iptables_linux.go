// go:build linux

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
	"github.com/webx-top/com"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/code"

	"github.com/admpub/nging/v5/application/handler"
	"github.com/admpub/nging/v5/application/library/common"
	"github.com/admpub/nging/v5/application/registry/navigate"
	"github.com/nging-plugins/firewallmanager/application/library/cmder"
	"github.com/nging-plugins/firewallmanager/application/library/driver"
	"github.com/nging-plugins/firewallmanager/application/library/driver/iptables"
	"github.com/nging-plugins/firewallmanager/application/library/enums"
	"github.com/nging-plugins/firewallmanager/application/library/firewall"
)

func init() {
	routeRegisters.Register(func(g echo.RouteRegister) {
		subG := g.Group(`/iptables`)
		subG.Route(`GET`, `/index`, ipTablesIndex)
		subG.Route(`GET`, `/delete`, ipTablesDelete)
	})
	LeftNavigate.Children.Add(-1, &navigate.Item{
		Display: false,
		Name:    `IPTables`,
		Action:  `iptables/index`,
	}, &navigate.Item{
		Display: false,
		Name:    `删除IPTables规则`,
		Action:  `iptables/delete`,
	})
}

var ipTablesFieldList = []string{
	`num`, `pkts`, `bytes`, `target`, `prot`, `opt`, `in`, `out`, `source`, `destination`, // original fields
	`options`, // custom fields
}

func ipTablesGetTableAndChain(ctx echo.Context) (ipVer string, table string, chain string, chainList []string) {
	ipVer = ctx.Form(`ipVer`, `4`)
	table = ctx.Form(`table`, enums.TableFilter)
	chain = ctx.Form(`chain`, enums.ChainInput)
	if !com.InSlice(table, enums.TableList) {
		table = enums.TableFilter
	}
	chainList = enums.TablesChains[table]
	if !com.InSlice(chain, chainList) {
		chain = chainList[0]
	}
	if ipVer != `4` && ipVer != `6` {
		ipVer = `4`
	}
	return
}

func ipTablesIndex(ctx echo.Context) error {
	if !iptables.IsSupported() {
		return ctx.NewError(code.Unsupported, `未安装 iptables`)
	}
	ipVer, table, chain, chainList := ipTablesGetTableAndChain(ctx)
	rules, err := firewall.Engine(ipVer).Stats(table, chain)
	if err != nil {
		return err
	}
	ctx.Set(`listData`, rules)
	ctx.Set(`tableList`, enums.TableList)
	ctx.Set(`chainList`, chainList)
	ctx.Set(`ipVerList`, enums.IPProtocols.Slice())
	// ctx.Set(`targetList`, enums.TargetList)
	// ctx.Set(`protocolList`, enums.ProtocolList)
	ctx.Set(`fieldList`, ipTablesFieldList)
	ctx.Set(`table`, table)
	ctx.Set(`chain`, chain)
	ctx.Set(`ipVer`, ipVer)
	ctx.Set(`lastModidyTs`, getStaticRuleLastModifyTs())
	if ctx.Form(`from`) == `dynamic` {
		ctx.Set(`activeURL`, `/firewall/rule/dynamic`)
	} else {
		ctx.Set(`activeURL`, `/firewall/rule/static`)
	}
	ctx.SetFunc(`canDelete`, ipTablesCanDelete)
	return ctx.Render(`firewall/iptables/index`, common.Err(ctx, err))
}

func ipTablesCanDelete(target string) bool {
	return target != cmder.DefaultChainName
}

func ipTablesDelete(ctx echo.Context) error {
	if !iptables.IsSupported() {
		return ctx.NewError(code.Unsupported, `未安装 iptables`)
	}
	id := ctx.Formx(`id`).Uint64()
	ts := ctx.Formx(`ts`).Uint64()
	ipVer, table, chain, _ := ipTablesGetTableAndChain(ctx)
	if ts != getStaticRuleLastModifyTs() {
		handler.SendErr(ctx, ctx.NewError(code.Failure, `操作失败，规则有更改，编号可能已经发生变化，请重新操作`))
		return ctx.Redirect(handler.URLFor(`/firewall/iptables/index`) + `?ipVer=` + ipVer + `&table=` + table + `&chain=` + chain)
	}
	err := firewall.Engine(ipVer).Delete(driver.Rule{
		Number:    id,
		Type:      table,
		Direction: chain,
	})

	if err == nil {
		handler.SendOk(ctx, ctx.T(`删除成功`))
	} else {
		handler.SendErr(ctx, err)
	}
	from := ctx.Form(`from`, `dynamic`)
	return ctx.Redirect(handler.URLFor(`/firewall/iptables/index`) + `?from=` + from + `&ipVer=` + ipVer + `&table=` + table + `&chain=` + chain)
}
