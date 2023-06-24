package biz

import (
	utils "github.com/admpub/nftablesutils"
	"github.com/google/nftables"
	"github.com/google/nftables/expr"
)

// inputHostBaseRules to apply.
func (nft *NFTables) inputHostBaseRules(c *nftables.Conn, iface string) error {
	// cmd: nft add rule ip filter input meta iifname "eth0" ip protocol icmp \
	// ct state { established, related } accept
	// --
	// iifname "eth0" ip protocol icmp ct state { established, related } accept
	ctStateSet := utils.GetConntrackStateSet(nft.tFilter)
	elems := utils.GetConntrackStateSetElems(defaultStateWithOld)
	err := c.AddSet(ctStateSet, elems)
	if err != nil {
		return err
	}

	exprs := make([]expr.Any, 0, 7)
	exprs = append(exprs, utils.SetIIF(iface)...)
	switch nft.tableFamily {
	case nftables.TableFamilyIPv4:
		exprs = append(exprs, utils.SetProtoICMP()...)
	case nftables.TableFamilyIPv6:
		exprs = append(exprs, utils.SetProtoICMPv6()...)
	}
	exprs = append(exprs, utils.SetConntrackStateSet(ctStateSet)...)
	exprs = append(exprs, utils.ExprAccept())

	rule := &nftables.Rule{
		Table: nft.tFilter,
		Chain: nft.cInput,
		Exprs: exprs,
	}
	c.AddRule(rule)

	// DNS
	if err = nft.inputDNSRules(c, iface); err != nil {
		return err
	}
	// HTTP Server
	if err = nft.inputHTTPServerRules(c, iface); err != nil {
		return err
	}

	return nil
}

// outputHostBaseRules to apply.
func (nft *NFTables) outputHostBaseRules(c *nftables.Conn, iface string) error {
	// cmd: nft add rule ip filter output meta oifname "eth0" ip protocol icmp \
	// ct state { new, established } accept
	// --
	// oifname "eth0" ip protocol icmp ct state { established, new } accept
	ctStateSet := utils.GetConntrackStateSet(nft.tFilter)
	elems := utils.GetConntrackStateSetElems(defaultStateWithNew)
	err := c.AddSet(ctStateSet, elems)
	if err != nil {
		return err
	}

	exprs := make([]expr.Any, 0, 7)
	exprs = append(exprs, utils.SetOIF(iface)...)
	switch nft.tableFamily {
	case nftables.TableFamilyIPv4:
		exprs = append(exprs, utils.SetProtoICMP()...)
	case nftables.TableFamilyIPv6:
		exprs = append(exprs, utils.SetProtoICMPv6()...)
	}
	exprs = append(exprs, utils.SetConntrackStateSet(ctStateSet)...)
	exprs = append(exprs, utils.ExprAccept())

	rule := &nftables.Rule{
		Table: nft.tFilter,
		Chain: nft.cOutput,
		Exprs: exprs,
	}
	c.AddRule(rule)

	// DNS
	if err = nft.outputDNSRules(c, iface); err != nil {
		return err
	}
	// HTTP Server
	if err = nft.outputHTTPServerRules(c, iface); err != nil {
		return err
	}

	return nil
}
