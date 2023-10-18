package biz

import (
	utils "github.com/admpub/nftablesutils"
	"github.com/google/nftables"
	"github.com/google/nftables/expr"
)

func (nft *NFTables) blacklistRules(c *nftables.Conn) error {
	exprs := make([]expr.Any, 0, 3)
	switch nft.tableFamily {
	case nftables.TableFamilyIPv4:
		exprs = append(exprs, utils.SetSAddrSet(nft.filterSetBlacklistIP)...)
	case nftables.TableFamilyIPv6:
		exprs = append(exprs, utils.SetSAddrIPv6Set(nft.filterSetBlacklistIP)...)
	}
	exprs = append(exprs, utils.Reject())
	rule := &nftables.Rule{
		Table: nft.tFilter,
		Chain: nft.cInput,
		Exprs: exprs,
	}
	c.AddRule(rule)
	return nil
}
