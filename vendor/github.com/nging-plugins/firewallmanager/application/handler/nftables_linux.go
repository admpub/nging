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
	pkgnftables "github.com/google/nftables"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/code"
	"github.com/webx-top/pagination"

	"github.com/admpub/nging/v5/application/handler"
	"github.com/admpub/nging/v5/application/library/common"
	"github.com/admpub/nging/v5/application/registry/navigate"
	"github.com/nging-plugins/firewallmanager/application/library/cmder"
	"github.com/nging-plugins/firewallmanager/application/library/driver"
	"github.com/nging-plugins/firewallmanager/application/library/driver/nftables"
	"github.com/nging-plugins/firewallmanager/application/library/enums"
	"github.com/nging-plugins/firewallmanager/application/library/firewall"
)

func init() {
	routeRegisters.Register(func(g echo.RouteRegister) {
		subG := g.Group(`/nftables`)
		subG.Route(`GET`, `/index`, nfTablesIndex)
		subG.Route(`GET`, `/delete`, nfTablesDelete)
	})
	LeftNavigate.Children.Add(-1, &navigate.Item{
		Display: true,
		Name:    `NFTables`,
		Action:  `nftables/index`,
	}, &navigate.Item{
		Display: false,
		Name:    `删除NFTables规则`,
		Action:  `nftables/delete`,
	})
}

var nfTablesFieldList = []string{
	`num`, `pkts`, `bytes`, `target`, `prot`, `opt`, `in`, `out`, `source`, `destination`, // original fields
	`options`, // custom fields
}

func nfTablesGetTableAndChain(ctx echo.Context) (ipVer string, table string, chain string, set string) {
	ipVer = ctx.Form(`ipVer`)
	table = ctx.Form(`table`)
	chain = ctx.Form(`chain`)
	set = ctx.Form(`set`)
	if ipVer != `4` && ipVer != `6` {
		ipVer = `4`
	}
	return
}

func nfTablesIndex(ctx echo.Context) error {
	if !nftables.IsSupported() {
		return ctx.NewError(code.Unsupported, `未安装 nftables`)
	}
	ipVer, table, chain, set := nfTablesGetTableAndChain(ctx)
	nft, ok := firewall.Engine(ipVer).(*nftables.NFTables)
	if !ok {
		return ctx.NewError(code.Unsupported, `不支持 nftables`)
	}
	var list interface{}
	var tableList []*pkgnftables.Table
	var chainList []*pkgnftables.Chain
	var setList []*pkgnftables.Set
	err := nft.Do(func(conn *pkgnftables.Conn) (err error) {
		var family pkgnftables.TableFamily
		if ipVer == `4` {
			family = pkgnftables.TableFamilyIPv4
		} else {
			family = pkgnftables.TableFamilyIPv6
		}
		tableList, err = conn.ListTablesOfFamily(family)
		if err != nil {
			return
		}
		if len(tableList) == 0 {
			return nil
		}
		var tableObj *pkgnftables.Table
		if len(table) == 0 && len(tableList) > 0 {
			table = tableList[0].Name
			tableObj = tableList[0]
		} else {
			for _, tb := range tableList {
				if tb.Name == table {
					tableObj = tb
					break
				}
			}
			if tableObj == nil {
				return nil
			}
		}
		setList, err = conn.GetSets(tableObj)
		if err != nil {
			return
		}
		var _chainList []*pkgnftables.Chain
		_chainList, err = conn.ListChainsOfTableFamily(family)
		if err != nil {
			return
		}
		for _, _chain := range _chainList {
			if _chain.Table.Name == table {
				chainList = append(chainList, _chain)
			}
		}
		if len(set) == 0 {
			if len(chain) == 0 && len(chainList) > 0 {
				chain = chainList[0].Name
			}
		}
		return
	})
	if err != nil {
		return err
	}
	pageNumber, pageSize := common.Paging(ctx)
	var hasMore bool
	if len(set) > 0 {
		list, hasMore, err = nft.ListSets(table, set, uint(pageNumber), uint(pageSize))
	} else {
		list, hasMore, err = nft.ListChainRules(table, chain, uint(pageNumber), uint(pageSize))
	}
	//echo.Dump(echo.H{`list`: rules, `hasMore`: hasMore, `err`: err})
	ctx.Set(`listData`, list)
	ctx.Set(`hasMore`, hasMore)

	q := ctx.Request().URL().Query()
	q.Del(`page`)
	q.Del(`rows`)
	q.Del(`_pjax`)
	paging := pagination.New(ctx).SetURL(`/firewall/nftables/index?` + q.Encode() + `&page={page}&rows={rows}`).SetPage(pageNumber)
	if hasMore {
		paging.SetRows((pageNumber + 1) * pageSize)
	} else {
		paging.SetRows(pageNumber * pageSize)
	}
	ctx.Set(`pagination`, paging)

	ctx.Set(`tableList`, tableList)
	ctx.Set(`chainList`, chainList)
	ctx.Set(`setList`, setList)
	ctx.Set(`ipVerList`, enums.IPProtocols.Slice())
	ctx.Set(`table`, table)
	ctx.Set(`chain`, chain)
	ctx.Set(`set`, set)
	ctx.Set(`ipVer`, ipVer)
	ctx.SetFunc(`canDelete`, nfTablesCanDelete)
	return ctx.Render(`firewall/nftables/index`, common.Err(ctx, err))
}

func nfTablesCanDelete(table string) bool {
	return table != cmder.DefaultTable4Name && table != cmder.DefaultTable6Name
}

func nfTablesDelete(ctx echo.Context) error {
	if !nftables.IsSupported() {
		return ctx.NewError(code.Unsupported, `未安装 nftables`)
	}
	id := ctx.Formx(`id`).Uint64()
	ipVer, table, chain, set := nfTablesGetTableAndChain(ctx)
	nft, ok := firewall.Engine(ipVer).(*nftables.NFTables)
	if !ok {
		return ctx.NewError(code.Unsupported, `不支持 nftables`)
	}
	var err error
	if len(set) > 0 {
		err = nft.DeleteElementInSetByHandleID(table, set, id)
	} else {
		//err = nft.DeleteRuleByHandleID(table, chain, id)
		err = firewall.Engine(ipVer).Delete(driver.Rule{
			Number:    id,
			Type:      table,
			Direction: chain,
			IPVersion: ipVer,
		})
	}
	if err == nil {
		handler.SendOk(ctx, ctx.T(`删除成功`))
	} else {
		handler.SendErr(ctx, err)
	}
	qs := `?ipVer=` + ipVer + `&table=` + table
	if len(set) > 0 {
		qs += `&set=` + set
	} else {
		qs += `&chain=` + chain
	}
	return ctx.Redirect(handler.URLFor(`/firewall/nftables/index`) + qs)
}
