package validate_dns

import (
	"regexp"
	"strings"
)

// ä, ö, ü, etc. is not supported.
const dnsLabelFmt = `[0-9a-z]([-0-9a-z]*[0-9a-z])?`

// Note, that in RFC 1123, applications must handle length up to 63 and should handle length up to 255.
// So this limit is pretty strict.
const maxDNSLabelLen = 63

var regDNSLabel *regexp.Regexp = regexp.MustCompile(`^` + dnsLabelFmt + `$`)

// IsDNSLabel returns true when the given string
// resembles a valid DNS label e.g. "de", but not "-de"
func IsDNSLabel(l string) bool {
	if len(l) > 63 {
		return false
	}
	if !regDNSLabel.MatchString(l) {
		return false
	}
	return true
}

const dnsFmt = `^` + dnsLabelFmt + `(\.` + dnsLabelFmt + `)+$`

const maxDNSLen = 255

var regDNSSub *regexp.Regexp = regexp.MustCompile(dnsFmt)

// IsDNS returns true, if given string resembles a valid DNS name.
// "example.com" will yield true, but not "ex_mple.com" or "example"
func IsDNS(s string) bool {
	if len(s) > maxDNSLen {
		return false
	}
	labels := strings.Split(s, ".")
	for _, l := range labels {
		if !IsDNSLabel(l) {
			return false
		}
	}
	if !regDNSSub.MatchString(s) {
		return false
	}
	return true
}

// IsDNSOrHostname returns true, if given string is a valid
// DNS name or hostname.
func IsDNSOrHostname(s string) bool {
	if len(s) > maxDNSLen {
		return false
	}
	labels := strings.Split(s, ".")
	for _, l := range labels {
		if !IsDNSLabel(l) {
			return false
		}
	}
	return true
}
