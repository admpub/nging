package resolver

import (
	"net"
	"strings"

	"github.com/miekg/dns"
)

// ResolveDNS will query DNS for a given hostname.
// example: ResolveDNS(`webx.top`,`8.8.8.8`,`IPV4`)
func ResolveDNS(hostname, resolver, ipType string) (string, error) {
	var dnsType uint16
	if len(ipType) == 0 || strings.ToUpper(ipType) == `IPV4` {
		dnsType = dns.TypeA
	} else {
		dnsType = dns.TypeAAAA
	}

	// If no DNS server is set in config file, falls back to default resolver.
	if len(resolver) == 0 {
		dnsAdress, err := net.LookupHost(hostname)
		if err != nil {
			return "", err
		}

		return dnsAdress[0], nil
	}
	res := New([]string{resolver})
	// In case of i/o timeout
	res.RetryTimes = 5

	ip, err := res.LookupHost(hostname, dnsType)
	if err != nil {
		return "", err
	}

	if len(ip) > 0 {
		return ip[0].String(), nil
	}

	return ``, nil
}
