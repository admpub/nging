package biz

import (
	utils "github.com/admpub/nftablesutils"
	"github.com/google/nftables"
	"github.com/google/nftables/expr"
)

func (nft *NFTables) forwardSMTPRules(c *nftables.Conn) error {
	if !nft.cfg.CanApply(ApplyTypeSMTP) {
		return nil
	}
	// cmd: nft add rule ip filter forward \
	// ip protocol tcp tcp sport 25 drop
	// --
	// tcp sport smtp drop;
	exprs := make([]expr.Any, 0, 10)
	exprs = append(exprs, utils.SetProtoTCP()...)
	exprs = append(exprs, utils.SetSPort(25)...)
	exprs = append(exprs, utils.ExprDrop())
	rule := &nftables.Rule{
		Table: nft.tFilter,
		Chain: nft.cForward,
		Exprs: exprs,
	}
	c.AddRule(rule)
	return nil
}
