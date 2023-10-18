package biz

import (
	utils "github.com/admpub/nftablesutils"
	"github.com/google/nftables"
	"github.com/google/nftables/expr"
)

func (nft *NFTables) forwardInterfaceRules(c *nftables.Conn) error {
	if len(nft.myIface) == 0 {
		return nil
	}

	// cmd: nft add rule ip filter forward \
	// meta iifname "wg0" \
	// ip saddr @wgforward_ipset \
	// meta oifname "eth0" \
	// accept
	// --
	// iifname "wg0" oifname "eth0" accept;
	{
		exprs := make([]expr.Any, 0, 10)
		exprs = append(exprs, utils.SetIIF(nft.myIface)...)
		switch nft.tableFamily {
		case nftables.TableFamilyIPv4:
			exprs = append(exprs, utils.SetSAddrSet(nft.filterSetForwardIP)...)
		case nftables.TableFamilyIPv6:
			exprs = append(exprs, utils.SetSAddrIPv6Set(nft.filterSetForwardIP)...)
		}
		exprs = append(exprs, utils.SetOIF(nft.wanIface)...)
		exprs = append(exprs, utils.ExprAccept())
		rule := &nftables.Rule{
			Table: nft.tFilter,
			Chain: nft.cForward,
			Exprs: exprs,
		}
		c.AddRule(rule)
	}

	// cmd: nft add rule ip filter forward \
	// ct state { established, related } accept
	// --
	// ct state { established, related } accept;
	ctStateSet := utils.GetConntrackStateSet(nft.tFilter)
	elems := utils.GetConntrackStateSetElems(defaultStateWithOld)
	err := c.AddSet(ctStateSet, elems)
	if err != nil {
		return err
	}

	{
		exprs := make([]expr.Any, 0, 10)
		exprs = append(exprs, utils.SetIIF(nft.wanIface)...)
		switch nft.tableFamily {
		case nftables.TableFamilyIPv4:
			exprs = append(exprs, utils.SetDAddrSet(nft.filterSetForwardIP)...)
		case nftables.TableFamilyIPv6:
			exprs = append(exprs, utils.SetDAddrIPv6Set(nft.filterSetForwardIP)...)
		}
		exprs = append(exprs, utils.SetOIF(nft.myIface)...)
		exprs = append(exprs, utils.SetConntrackStateSet(ctStateSet)...)
		exprs = append(exprs, utils.ExprAccept())
		rule := &nftables.Rule{
			Table: nft.tFilter,
			Chain: nft.cForward,
			Exprs: exprs,
		}
		c.AddRule(rule)
	}

	// cmd: nft add rule ip filter forward \
	// meta iifname "wg0" \
	// meta oifname "wg0" \
	// accept
	// --
	// iifname "wg0" oifname "wg0" accept;
	exprs := make([]expr.Any, 0, 10)
	exprs = append(exprs, utils.SetIIF(nft.myIface)...)
	exprs = append(exprs, utils.SetOIF(nft.myIface)...)
	exprs = append(exprs, utils.ExprAccept())
	rule := &nftables.Rule{
		Table: nft.tFilter,
		Chain: nft.cForward,
		Exprs: exprs,
	}
	c.AddRule(rule)
	return nil
}
