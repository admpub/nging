package biz

import (
	utils "github.com/admpub/nftablesutils"
	"github.com/google/nftables"
	"github.com/google/nftables/expr"
)

func (nft *NFTables) inputDNSRules(c *nftables.Conn, iface string) error {
	if !nft.cfg.CanApply(ApplyTypeDNS) {
		return nil
	}
	// cmd: nft add rule ip filter input meta iifname "eth0" \
	// ip protocol udp udp sport 53 \
	// ct state established accept
	// --
	// iifname "eth0" udp sport domain ct state established accept
	exprs := make([]expr.Any, 0, 10)
	exprs = append(exprs, utils.SetIIF(iface)...)
	exprs = append(exprs, utils.SetProtoUDP()...)
	exprs = append(exprs, utils.SetSPort(53)...)
	exprs = append(exprs, utils.SetConntrackStateEstablished()...)
	exprs = append(exprs, utils.ExprAccept())

	rule := &nftables.Rule{
		Table: nft.tFilter,
		Chain: nft.cInput,
		Exprs: exprs,
	}
	c.AddRule(rule)

	// cmd: nft add rule ip filter input meta iifname "eth0" \
	// ip protocol tcp tcp sport 53 \
	// ct state established accept
	// --
	// iifname "eth0" tcp sport domain ct state established accept
	exprs = make([]expr.Any, 0, 10)
	exprs = append(exprs, utils.SetIIF(iface)...)
	exprs = append(exprs, utils.SetProtoTCP()...)
	exprs = append(exprs, utils.SetSPort(53)...)
	exprs = append(exprs, utils.SetConntrackStateEstablished()...)
	exprs = append(exprs, utils.ExprAccept())

	rule = &nftables.Rule{
		Table: nft.tFilter,
		Chain: nft.cInput,
		Exprs: exprs,
	}
	c.AddRule(rule)
	return nil
}

func (nft *NFTables) outputDNSRules(c *nftables.Conn, iface string) error {
	if !nft.cfg.CanApply(ApplyTypeDNS) {
		return nil
	}
	// cmd: nft add rule ip filter output meta oifname "eth0" \
	// ip protocol udp udp dport 53 \
	// ct state { new, established } accept
	// --
	// oifname "eth0" udp dport domain ct state { established, new } accept
	ctStateSet := utils.GetConntrackStateSet(nft.tFilter)
	elems := utils.GetConntrackStateSetElems(defaultStateWithNew)
	err := c.AddSet(ctStateSet, elems)
	if err != nil {
		return err
	}
	exprs := make([]expr.Any, 0, 10)
	exprs = append(exprs, utils.SetOIF(iface)...)
	exprs = append(exprs, utils.SetProtoUDP()...)
	exprs = append(exprs, utils.SetDPort(53)...)
	exprs = append(exprs, utils.SetConntrackStateSet(ctStateSet)...)
	exprs = append(exprs, utils.ExprAccept())

	rule := &nftables.Rule{
		Table: nft.tFilter,
		Chain: nft.cOutput,
		Exprs: exprs,
	}
	c.AddRule(rule)

	// cmd: nft add rule ip filter output meta oifname "eth0" \
	// ip protocol tcp tcp dport 53 \
	// ct state { new, established } accept
	// --
	// oifname "eth0" tcp dport domain ct state { established, new } accept
	ctStateSet = utils.GetConntrackStateSet(nft.tFilter)
	elems = utils.GetConntrackStateSetElems(defaultStateWithNew)
	err = c.AddSet(ctStateSet, elems)
	if err != nil {
		return err
	}
	exprs = make([]expr.Any, 0, 10)
	exprs = append(exprs, utils.SetOIF(iface)...)
	exprs = append(exprs, utils.SetProtoTCP()...)
	exprs = append(exprs, utils.SetDPort(53)...)
	exprs = append(exprs, utils.SetConntrackStateSet(ctStateSet)...)
	exprs = append(exprs, utils.ExprAccept())

	rule = &nftables.Rule{
		Table: nft.tFilter,
		Chain: nft.cOutput,
		Exprs: exprs,
	}
	c.AddRule(rule)
	return nil
}
