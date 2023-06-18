package biz

import (
	utils "github.com/admpub/nftablesutils"
	"github.com/google/nftables"
	"github.com/google/nftables/expr"
)

func (nft *NFTables) inputHTTPServerRules(c *nftables.Conn, iface string) error {
	if !nft.cfg.CanApply(ApplyTypeHTTP) {
		return nil
	}
	// cmd: nft add rule ip filter input meta iifname "eth0" \
	// ip protocol tcp tcp sport { 80, 443 } \
	// ct state established accept
	// --
	// iifname "eth0" tcp sport { http, https } ct state established accept
	portSet := utils.GetPortSet(nft.tFilter)
	// portSet := &nftables.Set{Anonymous: true, Constant: true,
	// 	Table: nft.tFilter, KeyType: nftables.TypeInetService}
	elems := utils.GetPortElems([]uint16{80, 443})
	err := c.AddSet(portSet, elems)
	if err != nil {
		return err
	}

	exprs := make([]expr.Any, 0, 10)
	exprs = append(exprs, utils.SetIIF(iface)...)
	exprs = append(exprs, utils.SetProtoTCP()...)
	exprs = append(exprs, utils.SetSPortSet(portSet)...)
	exprs = append(exprs, utils.SetConntrackStateEstablished()...)
	exprs = append(exprs, utils.ExprAccept())

	rule := &nftables.Rule{
		Table: nft.tFilter,
		Chain: nft.cInput,
		Exprs: exprs,
	}
	c.AddRule(rule)
	return nil
}

func (nft *NFTables) outputHTTPServerRules(c *nftables.Conn, iface string) error {
	if !nft.cfg.CanApply(ApplyTypeHTTP) {
		return nil
	}
	// cmd: nft add rule ip filter output meta oifname "eth0" \
	// ip protocol tcp tcp dport { 80, 443 } \
	// ct state { new, established } accept
	// --
	// oifname "eth0" tcp dport { http, https } ct state { established, new } accept
	portSet := utils.GetPortSet(nft.tFilter)
	elems := utils.GetPortElems([]uint16{80, 443})
	err := c.AddSet(portSet, elems)
	if err != nil {
		return err
	}

	ctStateSet := utils.GetConntrackStateSet(nft.tFilter)
	elems = utils.GetConntrackStateSetElems(defaultStateWithNew)
	err = c.AddSet(ctStateSet, elems)
	if err != nil {
		return err
	}

	exprs := make([]expr.Any, 0, 10)
	exprs = append(exprs, utils.SetOIF(iface)...)
	exprs = append(exprs, utils.SetProtoTCP()...)
	exprs = append(exprs, utils.SetDPortSet(portSet)...)
	exprs = append(exprs, utils.SetConntrackStateSet(ctStateSet)...)
	exprs = append(exprs, utils.ExprAccept())

	rule := &nftables.Rule{
		Table: nft.tFilter,
		Chain: nft.cOutput,
		Exprs: exprs,
	}
	c.AddRule(rule)
	return nil
}
