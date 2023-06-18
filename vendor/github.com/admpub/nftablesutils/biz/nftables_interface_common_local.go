package biz

import (
	utils "github.com/admpub/nftablesutils"
	"github.com/google/nftables"
	"github.com/google/nftables/expr"
	"golang.org/x/sys/unix"
)

// inputLocalIfaceRules to apply.
func (nft *NFTables) inputLocalIfaceRules(c *nftables.Conn) {
	// cmd: nft add rule ip filter input meta iifname "lo" accept
	// --
	// iifname "lo" accept
	exprs := make([]expr.Any, 0, 3)
	exprs = append(exprs, utils.SetIIF(loIface)...)
	exprs = append(exprs, utils.ExprAccept())
	rule := &nftables.Rule{
		Table: nft.tFilter,
		Chain: nft.cInput,
		Exprs: exprs,
	}
	c.AddRule(rule)

	// cmd: nft add rule ip filter input meta iifname != "lo" \
	// ip saddr 127.0.0.0/8 reject
	// --
	// iifname != "lo" ip saddr 127.0.0.0/8 reject with icmp type prot-unreachable
	exprs = make([]expr.Any, 0, 6)
	exprs = append(exprs, utils.SetNIIF(loIface)...)

	switch nft.tableFamily {
	case nftables.TableFamilyIPv4: //127.0.0.0/24
		exprs = append(exprs, utils.SetSourceIPv4Net([]byte{127, 0, 0, 0}, []byte{255, 255, 255, 0})...)
		exprs = append(exprs, utils.ExprReject(
			unix.NFT_REJECT_ICMP_UNREACH,
			unix.NFT_REJECT_ICMPX_UNREACH,
		))
	case nftables.TableFamilyIPv6:
		exprs = append(exprs, utils.SetCIDRMatcher(utils.ExprDirectionSource, `fe80::/10`, false)...)
		exprs = append(exprs, utils.ExprReject(
			unix.NFT_REJECT_ICMP_UNREACH,
			unix.NFT_REJECT_ICMPX_NO_ROUTE,
		))
	}
	rule = &nftables.Rule{
		Table: nft.tFilter,
		Chain: nft.cInput,
		Exprs: exprs,
	}
	c.AddRule(rule)
}

// outputLocalIfaceRules to apply.
func (nft *NFTables) outputLocalIfaceRules(c *nftables.Conn) {
	// cmd: nft add rule ip filter output meta oifname "lo" accept
	// --
	// oifname "lo" accept
	exprs := make([]expr.Any, 0, 3)
	exprs = append(exprs, utils.SetOIF(loIface)...)
	exprs = append(exprs, utils.ExprAccept())
	rule := &nftables.Rule{
		Table: nft.tFilter,
		Chain: nft.cOutput,
		Exprs: exprs,
	}
	c.AddRule(rule)
}
