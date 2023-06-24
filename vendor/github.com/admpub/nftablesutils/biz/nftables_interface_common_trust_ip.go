package biz

import (
	utils "github.com/admpub/nftablesutils"
	"github.com/google/nftables"
	"github.com/google/nftables/expr"
)

var defaultStateWithNew = []string{utils.StateNew, utils.StateEstablished}
var defaultStateWithOld = []string{utils.StateEstablished, utils.StateRelated}

// inputTrustIPSetRules to apply.
func (nft *NFTables) inputTrustIPSetRules(c *nftables.Conn, iface string) error {
	if len(nft.cfg.TrustPorts) == 0 {
		return nil
	}
	// cmd: nft add rule ip filter input meta iifname "eth0" ip protocol icmp \
	// icmp type echo-request ip saddr @trust_ipset ct state new accept
	// --
	// iifname "eth0" icmp type echo-request ip saddr @trust_ipset ct state new accept
	exprs := make([]expr.Any, 0, 12)
	exprs = append(exprs, utils.SetIIF(iface)...)
	switch nft.tableFamily {
	case nftables.TableFamilyIPv4:
		exprs = append(exprs, utils.SetProtoICMP()...)
		exprs = append(exprs, utils.SetICMPTypeEchoRequest()...)
		exprs = append(exprs, utils.SetSAddrSet(nft.filterSetTrustIP)...)
	case nftables.TableFamilyIPv6:
		exprs = append(exprs, utils.SetProtoICMPv6()...)
		exprs = append(exprs, utils.SetICMPv6TypeEchoRequest()...)
		exprs = append(exprs, utils.SetSAddrIPv6Set(nft.filterSetTrustIP)...)
	}
	exprs = append(exprs, utils.SetConntrackStateNew()...)
	exprs = append(exprs, utils.ExprAccept())
	rule := &nftables.Rule{
		Table: nft.tFilter,
		Chain: nft.cInput,
		Exprs: exprs,
	}
	c.AddRule(rule)

	// cmd: nft add rule ip filter input meta iifname "eth0" \
	// ip protocol tcp tcp dport { 5522 } ip saddr @trust_ipset \
	// ct state { new, established } accept
	// --
	// iifname "eth0" tcp dport { 5522 } ip saddr @trust_ipset ct state { established, new } accept
	ctStateSet := utils.GetConntrackStateSet(nft.tFilter)
	elems := utils.GetConntrackStateSetElems(defaultStateWithNew)
	err := c.AddSet(ctStateSet, elems)
	if err != nil {
		return err
	}

	portSet := utils.GetPortSet(nft.tFilter)
	err = c.AddSet(portSet, nft.cfg.trustPorts())
	if err != nil {
		return err
	}

	exprs = make([]expr.Any, 0, 11)
	exprs = append(exprs, utils.SetIIF(iface)...)
	exprs = append(exprs, utils.SetProtoTCP()...)
	exprs = append(exprs, utils.SetDPortSet(portSet)...)
	switch nft.tableFamily {
	case nftables.TableFamilyIPv4:
		exprs = append(exprs, utils.SetSAddrSet(nft.filterSetTrustIP)...)
	case nftables.TableFamilyIPv6:
		exprs = append(exprs, utils.SetSAddrIPv6Set(nft.filterSetTrustIP)...)
	}
	exprs = append(exprs, utils.SetConntrackStateSet(ctStateSet)...)
	exprs = append(exprs, utils.ExprAccept())
	rule = &nftables.Rule{
		Table: nft.tFilter,
		Chain: nft.cInput,
		Exprs: exprs,
	}
	c.AddRule(rule)

	return nil
}

// outputTrustIPSetRules to apply.
func (nft *NFTables) outputTrustIPSetRules(c *nftables.Conn, iface string) error {
	if len(nft.cfg.TrustPorts) == 0 {
		return nil
	}
	// cmd: nft add rule ip filter output meta oifname "eth0" \
	// ip protocol tcp tcp sport { 5522 } ip daddr @trust_ipset \
	// ct state established accept
	// --
	// oifname "eth0" tcp sport { 5522 } ip daddr @trust_ipset ct state established accept
	portSet := utils.GetPortSet(nft.tFilter)
	err := c.AddSet(portSet, nft.cfg.trustPorts())
	if err != nil {
		return err
	}

	exprs := make([]expr.Any, 0, 12)
	exprs = append(exprs, utils.SetOIF(iface)...)
	exprs = append(exprs, utils.SetProtoTCP()...)
	exprs = append(exprs, utils.SetSPortSet(portSet)...)
	switch nft.tableFamily {
	case nftables.TableFamilyIPv4:
		exprs = append(exprs, utils.SetDAddrSet(nft.filterSetTrustIP)...)
	case nftables.TableFamilyIPv6:
		exprs = append(exprs, utils.SetDAddrIPv6Set(nft.filterSetTrustIP)...)
	}
	exprs = append(exprs, utils.SetConntrackStateEstablished()...)
	exprs = append(exprs, utils.ExprAccept())
	rule := &nftables.Rule{
		Table: nft.tFilter,
		Chain: nft.cOutput,
		Exprs: exprs,
	}
	c.AddRule(rule)

	return nil
}
