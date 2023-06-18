// Package extnetip is an extension to net/netip with
// a few missing but important auxiliary functions for
// converting IP-prefixes to IP-ranges and vice versa.
//
// The functions are effectively performed in uint128 space,
// no conversions from/to bytes are performed.
//
// With these extensions to net/netip, third-party IP-range
// libraries become easily possible.
package extnetip

import "net/netip"

// Range returns the inclusive range of IP addresses that p covers.
//
// If p is invalid, Range returns the zero values.
func Range(p netip.Prefix) (first, last netip.Addr) {
	if !p.IsValid() {
		return
	}

	// peek the internals, do math in uint128
	pa := peek(p.Addr())
	z := pa.z

	bits := p.Bits()
	if z == z4 {
		bits += 96
	}
	mask := mask6(bits)

	first128 := pa.ip.and(mask)
	last128 := first128.or(mask.not())

	// convert back to netip.Addr
	first = back(addr{first128, z})
	last = back(addr{last128, z})

	return
}

// Prefix returns the netip.Prefix from first to last and ok=true,
// if it can be presented exactly as such.
//
// If first or last are not valid, in the wrong order or not exactly
// equal to one prefix, ok is false.
func Prefix(first, last netip.Addr) (prefix netip.Prefix, ok bool) {
	if !(first.IsValid() && last.IsValid()) {
		return
	}
	if last.Less(first) {
		return
	}

	// peek the internals, do math in uint128
	pFirst := peek(first)
	pLast := peek(last)

	// IP versions differ?
	if pFirst.z != pLast.z {
		return
	}

	// do math in uint128
	bits, ok := pFirst.ip.prefixOK(pLast.ip)
	if !ok {
		return
	}

	if pFirst.z == z4 {
		bits -= 96
	}

	// make prefix, possible zone gets dropped
	return netip.PrefixFrom(first, bits), ok
}

// Prefixes returns the set of netip.Prefix entries that covers the
// IP range from first to last.
//
// If first or last are invalid, in the wrong order, or if they're of different
// address families, then Prefixes returns nil.
//
// Prefixes necessarily allocates. See PrefixesAppend for a version that
// uses memory you provide.
func Prefixes(first, last netip.Addr) []netip.Prefix {
	return PrefixesAppend(nil, first, last)
}

// PrefixesAppend is an append version of Prefixes. It appends
// the netip.Prefix entries to dst that covers the IP range from first to last.
func PrefixesAppend(dst []netip.Prefix, first, last netip.Addr) []netip.Prefix {
	if !(first.IsValid() && last.IsValid()) {
		return nil
	}
	if last.Less(first) {
		return nil
	}

	// peek the internals, do math in uint128
	pFirst := peek(first)
	pLast := peek(last)

	// different IP versions
	if pFirst.z != pLast.z {
		return nil
	}

	return prefixesAppendRec(dst, pFirst, pLast)
}

// append prefix if (first, last) represents a whole CIDR, like 10.0.0.0/8
// (first being 10.0.0.0 and last being 10.255.255.255)
//
// Otherwise recursively do both halves.
//
// Recursion is here faster than an iterative algo, no bounds checking and no heap escape and
// btw. the recursion level is max. 254 deep, for IP range: ::1-ffff:ffff:ffff:ffff:ffff:ffff:ffff:fffe
func prefixesAppendRec(dst []netip.Prefix, first, last addr) []netip.Prefix {
	// are first-last already representing a prefix?
	bits, ok := first.ip.prefixOK(last.ip)
	if ok {
		if first.z == z4 {
			bits -= 96
		}
		// convert back to netip
		pfx := netip.PrefixFrom(back(first), bits)

		return append(dst, pfx)
	}

	// otherwise split the range, make two halves and do both halves recursively
	mask := mask6(bits + 1)

	// make middle last, set hostbits
	midOne := addr{first.ip.or(mask.not()), first.z}

	// make middle next, clear hostbits
	midTwo := addr{last.ip.and(mask), first.z}

	// ... do both halves recursively
	dst = prefixesAppendRec(dst, first, midOne)
	dst = prefixesAppendRec(dst, midTwo, last)

	return dst
}
