package biz

import (
	utils "github.com/admpub/nftablesutils"
	"github.com/google/nftables"
	"github.com/google/nftables/expr"
)

// inputPublicRules to apply.
func (nft *NFTables) inputPublicRules(c *nftables.Conn, iface string) error {
	if nft.myPort <= 0 {
		return nil
	}
	// cmd: nft add rule ip filter input meta iifname "eth0" \
	// ip protocol udp udp dport 51820 accept
	// --
	// iifname "eth0" udp dport 51820 accept

	exprs := make([]expr.Any, 0, 9)
	exprs = append(exprs, utils.SetIIF(iface)...)
	exprs = append(exprs, utils.SetProtoUDP()...)
	exprs = append(exprs, utils.SetDPort(nft.myPort)...)
	exprs = append(exprs, utils.ExprAccept())
	rule := &nftables.Rule{
		Table: nft.tFilter,
		Chain: nft.cInput,
		Exprs: exprs,
	}
	c.AddRule(rule)

	return nil
}

// outputPublicRules to apply.
func (nft *NFTables) outputPublicRules(c *nftables.Conn, iface string) error {
	if nft.myPort <= 0 {
		return nil
	}
	// cmd: nft add rule ip filter output meta oifname "eth0" \
	// ip protocol udp udp sport 51820 accept
	// --
	// oifname "eth0" udp sport 51820 accept

	exprs := make([]expr.Any, 0, 10)
	exprs = append(exprs, utils.SetOIF(iface)...)
	exprs = append(exprs, utils.SetProtoUDP()...)
	exprs = append(exprs, utils.SetSPort(nft.myPort)...)
	exprs = append(exprs, utils.ExprAccept())
	rule := &nftables.Rule{
		Table: nft.tFilter,
		Chain: nft.cOutput,
		Exprs: exprs,
	}
	c.AddRule(rule)

	return nil
}
