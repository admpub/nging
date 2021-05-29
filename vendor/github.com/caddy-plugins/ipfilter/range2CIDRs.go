package ipfilter

// https://groups.google.com/forum/m/#!topic/golang-nuts/rJvVwk4jwjQ

import (
	"fmt"
	"net"
)

var allFF = net.ParseIP("255.255.255.255").To4()

func range2CIDRs(a1, a2 net.IP) (r []*net.IPNet) {
	// Warn users that they're using an old method of ranging
	fmt.Println("ipfilter Warning: You are using an old method of ranging over IPs, it's highly recommended to switch over to CIDR notation, For more: https://caddyserver.com/docs/ipfilter")
	maxLen := 32
	a1 = a1.To4()
	a2 = a2.To4()
	for cmp(a1, a2) <= 0 {
		l := 32
		for l > 0 {
			m := net.CIDRMask(l-1, maxLen)
			if cmp(a1, first(a1, m)) != 0 || cmp(last(a1, m), a2) > 0 {
				break
			}
			l--
		}
		r = append(r, &net.IPNet{IP: a1, Mask: net.CIDRMask(l, maxLen)})
		a1 = last(a1, net.CIDRMask(l, maxLen))
		if cmp(a1, allFF) == 0 {
			break
		}
		a1 = next(a1)
	}
	return r
}

func next(ip net.IP) net.IP {
	n := len(ip)
	out := make(net.IP, n)
	copy := false
	for n > 0 {
		n--
		if copy {
			out[n] = ip[n]
			continue
		}
		if ip[n] < 255 {
			out[n] = ip[n] + 1
			copy = true
			continue
		}
		out[n] = 0
	}
	return out
}

func cmp(ip1, ip2 net.IP) int {
	l := len(ip1)
	for i := 0; i < l; i++ {
		if ip1[i] == ip2[i] {
			continue
		}
		if ip1[i] < ip2[i] {
			return -1
		}
		return 1
	}
	return 0
}

func first(ip net.IP, mask net.IPMask) net.IP {
	return ip.Mask(mask)
}

func last(ip net.IP, mask net.IPMask) net.IP {
	n := len(ip)
	out := make(net.IP, n)
	for i := 0; i < n; i++ {
		out[i] = ip[i] | ^mask[i]
	}
	return out
}
