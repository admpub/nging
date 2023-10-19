package nftablesutils

import (
	"net"
	"net/netip"
	"strings"

	"github.com/google/nftables"
	"github.com/google/nftables/expr"
)

// Returns a IPv4 source address payload expression
func IPv4SourceAddress(reg uint32) *expr.Payload {
	return ExprPayloadNetHeader(reg, IPv4SrcOffset, IPv4AddrLen)
}

// Returns a IPv6 source address payload expression
func IPv6SourceAddress(reg uint32) *expr.Payload {
	return ExprPayloadNetHeader(reg, IPv6SrcOffset, IPv6AddrLen)
}

// Returns a IPv4 destination address payload expression
func IPv4DestinationAddress(reg uint32) *expr.Payload {
	return ExprPayloadNetHeader(reg, IPv4DstOffset, IPv4AddrLen)
}

// Returns a IPv6 destination address payload expression
func IPv6DestinationAddress(reg uint32) *expr.Payload {
	return ExprPayloadNetHeader(reg, IPv6DstOffset, IPv6AddrLen)
}

// SetCIDRMatcherIngoreError generates nftables expressions that matches a CIDR
// SetCIDRMatcherIngoreError(ExprDirectionSource, `127.0.0.0/24`)
func SetCIDRMatcherIngoreError(direction ExprDirection, cidr string, isINet bool, isEq ...bool) []expr.Any {
	exprs, _ := SetCIDRMatcher(direction, cidr, isINet, isEq...)
	return exprs
}

// SetCIDRMatcher generates nftables expressions that matches a CIDR
// SetCIDRMatcher(ExprDirectionSource, `127.0.0.0/24`)
func SetCIDRMatcher(direction ExprDirection, cidr string, isINet bool, isEq ...bool) ([]expr.Any, error) {
	if !strings.Contains(cidr, `/`) {
		if strings.Contains(cidr, `:`) {
			cidr += `/128`
		} else {
			cidr += `/32`
		}
	}
	ip, network, err := net.ParseCIDR(cidr)
	if err != nil {
		return nil, err
	}
	ipToAddr, _ := netip.AddrFromSlice(ip)
	addr := ipToAddr.Unmap()

	offSet, packetLen, zeroXor := GetPayloadDirectives(direction, addr.Is4(), addr.Is6())

	exprs := make([]expr.Any, 0, 5)
	if isINet {
		var family nftables.TableFamily
		if addr.Is4() {
			family = nftables.TableFamilyIPv4
		} else {
			family = nftables.TableFamilyIPv6
		}
		exprs = append(exprs, CompareProtocolFamily(family)...)
	}

	exprs = append(
		exprs,
		// fetch src add
		&expr.Payload{
			DestRegister: defaultRegister,
			Base:         expr.PayloadBaseNetworkHeader,
			Offset:       offSet,
			Len:          packetLen,
		},
		// net mask
		&expr.Bitwise{
			DestRegister:   defaultRegister,
			SourceRegister: defaultRegister,
			Len:            packetLen,
			Mask:           network.Mask,
			Xor:            zeroXor,
		},
		// net address
		&expr.Cmp{
			Op:       GetCmpOp(isEq...),
			Register: defaultRegister,
			Data:     addr.AsSlice(),
		},
	)
	return exprs, err
}

// SetSourceIPv4Net helper.
func SetSourceIPv4Net(addr []byte, mask []byte, isEq ...bool) Exprs {
	exprs := []expr.Any{
		IPv4SourceAddress(defaultRegister),
		ExprBitwise(defaultRegister, defaultRegister, IPv4AddrLen,
			mask,
			make([]byte, IPv4AddrLen),
		),
		ExprCmp(GetCmpOp(isEq...), addr),
	}
	return exprs
}

// SetSAddrSet helper.
func SetSAddrSet(s *nftables.Set, isEq ...bool) Exprs {
	exprs := []expr.Any{
		IPv4SourceAddress(defaultRegister),
		ExprLookupSet(defaultRegister, s.Name, s.ID, isEq...),
	}
	return exprs
}

// SetDAddrSet helper.
func SetDAddrSet(s *nftables.Set, isEq ...bool) Exprs {
	exprs := []expr.Any{
		IPv4DestinationAddress(defaultRegister),
		ExprLookupSet(defaultRegister, s.Name, s.ID, isEq...),
	}
	return exprs
}

// SetSAddrIPv6Set helper.
func SetSAddrIPv6Set(s *nftables.Set, isEq ...bool) Exprs {
	exprs := []expr.Any{
		IPv6SourceAddress(defaultRegister),
		ExprLookupSet(defaultRegister, s.Name, s.ID, isEq...),
	}
	return exprs
}

// SetDAddrIPv6Set helper.
func SetDAddrIPv6Set(s *nftables.Set, isEq ...bool) Exprs {
	exprs := []expr.Any{
		IPv6DestinationAddress(defaultRegister),
		ExprLookupSet(defaultRegister, s.Name, s.ID, isEq...),
	}
	return exprs
}

// GetIPv4AddrSet helper.
func GetIPv4AddrSet(t *nftables.Table, isInterval ...bool) *nftables.Set {
	s := &nftables.Set{
		Anonymous: true,
		Constant:  true,
		Table:     t,
		KeyType:   nftables.TypeIPAddr,
		Interval:  len(isInterval) > 0 && isInterval[0],
	}
	return s
}

// GetIPv6AddrSet helper.
func GetIPv6AddrSet(t *nftables.Table, isInterval ...bool) *nftables.Set {
	s := &nftables.Set{
		Anonymous: true,
		Constant:  true,
		Table:     t,
		KeyType:   nftables.TypeIP6Addr,
		Interval:  len(isInterval) > 0 && isInterval[0],
	}
	return s
}
