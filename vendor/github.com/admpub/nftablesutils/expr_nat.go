package nftablesutils

import (
	"net"

	"github.com/google/nftables"
	"github.com/google/nftables/binaryutil"
	"github.com/google/nftables/expr"
)

func SNAT(ip net.IP) []expr.Any {
	return []expr.Any{
		ExprImmediate(ip),
		ExprSNAT(defaultRegister, 0),
	}
}

func SNATv6(ip net.IP) []expr.Any {
	return []expr.Any{
		ExprImmediate(ip),
		ExprSNATv6(defaultRegister, 0),
	}
}

func DNAT(ip net.IP) []expr.Any {
	return []expr.Any{
		ExprImmediate(ip),
		ExprDNAT(defaultRegister, 0),
	}
}

func DNATv6(ip net.IP) []expr.Any {
	return []expr.Any{
		ExprImmediate(ip),
		ExprDNATv6(defaultRegister, 0),
	}
}

func RedirectTo(port uint16) []expr.Any {
	return []expr.Any{
		&expr.Immediate{
			Register: defaultRegister,
			Data:     binaryutil.BigEndian.PutUint16(port),
		},
		ExprRedirect(defaultRegister, 0),
	}
}

func ExprTproxy() *expr.TProxy {
	return &expr.TProxy{
		Family:      byte(nftables.TableFamilyIPv4),
		TableFamily: byte(nftables.TableFamilyIPv4),
		RegPort:     defaultRegister,
	}
}

func ExprTproxy6() *expr.TProxy {
	return &expr.TProxy{
		Family:      byte(nftables.TableFamilyIPv6),
		TableFamily: byte(nftables.TableFamilyIPv6),
		RegPort:     defaultRegister,
	}
}

// ExprSNAT wrapper
func ExprSNAT(addrMin, addrMax uint32) *expr.NAT {
	// [ nat snat ip addr_min reg 1 addr_max reg 0 ]
	return &expr.NAT{
		Type:       expr.NATTypeSourceNAT,
		Family:     uint32(nftables.TableFamilyIPv4),
		RegAddrMin: addrMin,
		RegAddrMax: addrMax,
	}
}

// ExprSNATv6 wrapper
func ExprSNATv6(addrMin, addrMax uint32) *expr.NAT {
	// [ nat snat ip addr_min reg 1 addr_max reg 0 ]
	return &expr.NAT{
		Type:       expr.NATTypeSourceNAT,
		Family:     uint32(nftables.TableFamilyIPv6),
		RegAddrMin: addrMin,
		RegAddrMax: addrMax,
	}
}

// ExprDNAT wrapper
func ExprDNAT(addrMin, addrMax uint32) *expr.NAT {
	// [ nat dnat ip addr_min reg 1 addr_max reg 0 ]
	return &expr.NAT{
		Type:       expr.NATTypeDestNAT,
		Family:     uint32(nftables.TableFamilyIPv4),
		RegAddrMin: addrMin,
		RegAddrMax: addrMax,
	}
}

// ExprDNATv6 wrapper
func ExprDNATv6(addrMin, addrMax uint32) *expr.NAT {
	// [ nat dnat ip addr_min reg 1 addr_max reg 0 ]
	return &expr.NAT{
		Type:       expr.NATTypeDestNAT,
		Family:     uint32(nftables.TableFamilyIPv6),
		RegAddrMin: addrMin,
		RegAddrMax: addrMax,
	}
}

// ExprMasquerade wrapper
func ExprMasquerade(protoMin, protoMax uint32) *expr.Masq {
	return &expr.Masq{
		Random:      false,
		FullyRandom: false,
		Persistent:  false,
		ToPorts:     false,
		RegProtoMin: protoMin,
		RegProtoMax: protoMax,
	}
}

// ExprRedirect wrapper
func ExprRedirect(protoMin, protoMax uint32) *expr.Redir {
	return &expr.Redir{
		Flags:            0,
		RegisterProtoMin: protoMin,
		RegisterProtoMax: protoMax,
	}
}
