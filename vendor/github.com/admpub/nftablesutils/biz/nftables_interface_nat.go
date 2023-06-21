package biz

import (
	utils "github.com/admpub/nftablesutils"
	"github.com/google/nftables"
	"github.com/google/nftables/expr"
)

func (nft *NFTables) natInterfaceRules(c *nftables.Conn) error {
	if len(nft.wanIface) == 0 || len(nft.wanIP) == 0 || nft.wanIP.IsUnspecified() {
		return nil
	}

	// cmd: nft add rule ip nat postrouting meta oifname "eth0" \
	// snat 192.168.0.1
	// --
	// oifname "eth0" snat to 192.168.15.11
	exprs := make([]expr.Any, 0, 10)
	exprs = append(exprs, utils.SetOIF(nft.wanIface)...)
	exprs = append(exprs, utils.ExprImmediate(1, nft.wanIP))
	switch nft.tNAT.Family {
	case nftables.TableFamilyIPv4:
		exprs = append(exprs, utils.ExprSNAT(1, 0))
	case nftables.TableFamilyIPv6:
		exprs = append(exprs, utils.ExprSNATv6(1, 0))
	}
	rule := &nftables.Rule{
		Table: nft.tNAT,
		Chain: nft.cPostrouting,
		Exprs: exprs,
	}
	c.AddRule(rule)
	return nil
}
