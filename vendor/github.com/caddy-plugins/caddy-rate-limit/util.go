package ratelimit

import (
	"net"
	"net/http"
	"strings"
)

// IsWhitelistIPAddress check whether an ip is in whitelist
func IsWhitelistIPAddress(address string, localIPNets []*net.IPNet) bool {

	ip := net.ParseIP(address)
	if ip != nil {
		for _, ipNet := range localIPNets {
			if ipNet.Contains(ip) {
				return true
			}
		}
	}

	return false
}

// GetRemoteIP returns the ip of requester
// Doesn't care if the ip is real or not
func GetRemoteIP(r *http.Request) (string, error) {
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	return host, err
}

// MatchMethod check whether the request method is in the methods list
func MatchMethod(methods, method string) bool {
	methods = strings.ToUpper(methods)
	if methods == "*" || strings.Contains(methods, method) {
		return true
	}
	return false
}

// MatchStatus check whether the upstream response status code is  in the status list
func MatchStatus(status, s string) bool {
	return strings.Contains(status, s)
}
