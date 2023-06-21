package nftablesutils

import (
	"net"

	"github.com/google/nftables"
	"github.com/google/nftables/expr"
)

func SetNATWithIPAndPort(
	dir ExprDirection, isIPv6 bool,
	ipStart net.IP, ipEnd net.IP,
	portMinAndMax ...uint16) []expr.Any {
	var regAddrMax, regPortMin, regPortMax uint32
	var incr uint32 = 1
	r := make([]expr.Any, 0, 5)
	r = append(r, ExprImmediate(defaultRegister, ipStart))
	if ipEnd != nil {
		regAddrMax = defaultRegister + incr
		incr++
		r = append(r, ExprImmediate(regAddrMax, ipEnd))
	}
	if len(portMinAndMax) > 0 && portMinAndMax[0] > 0 {
		regPortMin = defaultRegister + incr
		incr++
		r = append(r, ExprImmediateWithPort(regPortMin, portMinAndMax[0]))
	}
	if len(portMinAndMax) > 1 && portMinAndMax[1] > 0 {
		regPortMax = defaultRegister + incr
		r = append(r, ExprImmediateWithPort(regPortMax, portMinAndMax[1]))
	}
	if dir == ExprDirectionSource {
		if isIPv6 {
			r = append(r, ExprSNATv6(defaultRegister, regAddrMax, regPortMin, regPortMax))
		} else {
			r = append(r, ExprSNAT(defaultRegister, regAddrMax, regPortMin, regPortMax))
		}
	} else {
		if isIPv6 {
			r = append(r, ExprDNATv6(defaultRegister, regAddrMax, regPortMin, regPortMax))
		} else {
			r = append(r, ExprDNAT(defaultRegister, regAddrMax, regPortMin, regPortMax))
		}
	}
	return r
}

func SetSNAT(ip net.IP, portMinAndMax ...uint16) []expr.Any {
	return SetNATWithIPAndPort(ExprDirectionSource, false, ip, nil, portMinAndMax...)
}

func SetSNATRange(ipStart net.IP, ipEnd net.IP, portMinAndMax ...uint16) []expr.Any {
	return SetNATWithIPAndPort(ExprDirectionSource, false, ipStart, ipEnd, portMinAndMax...)
}

func SetSNATv6(ip net.IP, portMinAndMax ...uint16) []expr.Any {
	return SetNATWithIPAndPort(ExprDirectionSource, true, ip, nil, portMinAndMax...)
}

func SetSNATv6Range(ipStart net.IP, ipEnd net.IP, portMinAndMax ...uint16) []expr.Any {
	return SetNATWithIPAndPort(ExprDirectionSource, true, ipStart, ipEnd, portMinAndMax...)
}

func SetDNAT(ip net.IP, portMinAndMax ...uint16) []expr.Any {
	return SetNATWithIPAndPort(ExprDirectionDestination, false, ip, nil, portMinAndMax...)
}

func SetDNATRange(ipStart net.IP, ipEnd net.IP, portMinAndMax ...uint16) []expr.Any {
	return SetNATWithIPAndPort(ExprDirectionDestination, false, ipStart, ipEnd, portMinAndMax...)
}

func SetDNATv6(ip net.IP, portMinAndMax ...uint16) []expr.Any {
	return SetNATWithIPAndPort(ExprDirectionDestination, true, ip, nil, portMinAndMax...)
}

func SetDNATv6Range(ipStart net.IP, ipEnd net.IP, portMinAndMax ...uint16) []expr.Any {
	return SetNATWithIPAndPort(ExprDirectionDestination, true, ipStart, ipEnd, portMinAndMax...)
}

func SetRedirect(portMin uint16, portMax ...uint16) []expr.Any {
	if len(portMax) > 0 && portMax[0] > 0 {
		return []expr.Any{
			ExprImmediateWithPort(defaultRegister, portMin),
			ExprImmediateWithPort(defaultRegister+1, portMax[0]),
			ExprRedirect(defaultRegister, defaultRegister+1),
		}
	}
	return []expr.Any{
		ExprImmediateWithPort(defaultRegister, portMin),
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
func ExprSNAT(regAddrMin, regAddrMax uint32, regPortMinAndMax ...uint32) *expr.NAT {
	var regPortMin uint32
	var regPortMax uint32
	if len(regPortMinAndMax) > 0 {
		regPortMin = regPortMinAndMax[0]
	}
	if len(regPortMinAndMax) > 1 {
		regPortMax = regPortMinAndMax[1]
	}
	// [ nat snat ip addr_min reg 1 addr_max reg 0 ]
	return &expr.NAT{
		Type:        expr.NATTypeSourceNAT,
		Family:      uint32(nftables.TableFamilyIPv4),
		RegAddrMin:  regAddrMin, // 起始IP地址注册号(即Immediate的Register值，此Immediate的Data中保存有起始IP地址)
		RegAddrMax:  regAddrMax, // 终止IP地址注册号(即Immediate的Register值，此Immediate的Data中保存有终止IP地址)
		RegProtoMin: regPortMin, // 起始端口注册号(即Immediate的Register值，此Immediate的Data中保存有起始端口号)
		RegProtoMax: regPortMax, //终止端口注册号(即Immediate的Register值，此Immediate的Data中保存有终止端口号)
	}
}

// ExprSNATv6 wrapper
func ExprSNATv6(regAddrMin, regAddrMax uint32, regPortMinAndMax ...uint32) *expr.NAT {
	var regPortMin uint32
	var regPortMax uint32
	if len(regPortMinAndMax) > 0 {
		regPortMin = regPortMinAndMax[0]
	}
	if len(regPortMinAndMax) > 1 {
		regPortMax = regPortMinAndMax[1]
	}
	// [ nat snat ip addr_min reg 1 addr_max reg 0 ]
	return &expr.NAT{
		Type:        expr.NATTypeSourceNAT,
		Family:      uint32(nftables.TableFamilyIPv6),
		RegAddrMin:  regAddrMin,
		RegAddrMax:  regAddrMax,
		RegProtoMin: regPortMin,
		RegProtoMax: regPortMax,
	}
}

// ExprDNAT wrapper
func ExprDNAT(regAddrMin, regAddrMax uint32, regPortMinAndMax ...uint32) *expr.NAT {
	var regPortMin uint32
	var regPortMax uint32
	if len(regPortMinAndMax) > 0 {
		regPortMin = regPortMinAndMax[0]
	}
	if len(regPortMinAndMax) > 1 {
		regPortMax = regPortMinAndMax[1]
	}
	// [ nat dnat ip addr_min reg 1 addr_max reg 0 ]
	return &expr.NAT{
		Type:        expr.NATTypeDestNAT,
		Family:      uint32(nftables.TableFamilyIPv4),
		RegAddrMin:  regAddrMin,
		RegAddrMax:  regAddrMax,
		RegProtoMin: regPortMin,
		RegProtoMax: regPortMax,
	}
}

// ExprDNATv6 wrapper
func ExprDNATv6(regAddrMin, regAddrMax uint32, regPortMinAndMax ...uint32) *expr.NAT {
	var regPortMin uint32
	var regPortMax uint32
	if len(regPortMinAndMax) > 0 {
		regPortMin = regPortMinAndMax[0]
	}
	if len(regPortMinAndMax) > 1 {
		regPortMax = regPortMinAndMax[1]
	}
	// [ nat dnat ip addr_min reg 1 addr_max reg 0 ]
	return &expr.NAT{
		Type:        expr.NATTypeDestNAT,
		Family:      uint32(nftables.TableFamilyIPv6),
		RegAddrMin:  regAddrMin,
		RegAddrMax:  regAddrMax,
		RegProtoMin: regPortMin,
		RegProtoMax: regPortMax,
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
		Flags:            1,
		RegisterProtoMin: protoMin,
		RegisterProtoMax: protoMax,
	}
}
