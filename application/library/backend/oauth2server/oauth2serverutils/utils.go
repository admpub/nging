package oauth2serverutils

import "strings"

func MatchDomain(domain string, allowedDomains string) bool {
	if len(allowedDomains) == 0 || len(domain) == 0 {
		return true
	}
	outsiteDomains := strings.Split(allowedDomains, `,`)
	size := len(domain)
	for _, allowedDomain := range outsiteDomains {
		allowedDomainSize := len(allowedDomain)
		if allowedDomainSize == 0 || size < allowedDomainSize {
			continue
		}
		if size == allowedDomainSize {
			if allowedDomain == domain {
				return true
			}
			continue
		}
		if !strings.HasSuffix(domain, allowedDomain) {
			continue
		}
		if strings.HasPrefix(allowedDomain, `.`) {
			return true
		}
		if domain[(size-allowedDomainSize)-1] == '.' {
			return true
		}
	}
	return false
}
