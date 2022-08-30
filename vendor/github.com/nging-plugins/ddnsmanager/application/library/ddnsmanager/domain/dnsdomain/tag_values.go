package dnsdomain

import (
	"net/url"
	"strings"
)

func NewTagValues() *TagValues {
	return &TagValues{
		IPv4Result: Results{},
		IPv6Result: Results{},
	}
}

type TagValues struct {
	IPv4Addr    string
	IPv4Result  Results
	IPv4Domains []string
	IPv6Addr    string
	IPv6Result  Results
	IPv6Domains []string
	IPAddr      string
	Error       string

	escapedIPv4Addr    string
	escapedIPv4Result  string
	escapedIPv4Domains string
	escapedIPv6Addr    string
	escapedIPv6Result  string
	escapedIPv6Domains string
	escapedIPAddr      string
	escapedError       string
}

func (t *TagValues) Parse(content string) string {
	content = strings.ReplaceAll(content, Tag(`ip`), t.IPAddr)
	content = strings.ReplaceAll(content, Tag(`ipv4Addr`), t.IPv4Addr)                           // 新的IPv4地址
	content = strings.ReplaceAll(content, Tag(`ipv4Result`), t.IPv4Result.String())              // IPv4地址更新结果: `未改变` `失败` `成功`
	content = strings.ReplaceAll(content, Tag(`ipv4Domains`), strings.Join(t.IPv4Domains, `, `)) // IPv4的域名，多个以`,`分割
	content = strings.ReplaceAll(content, Tag(`ipv6Addr`), t.IPv6Addr)                           // 新的IPv6地址
	content = strings.ReplaceAll(content, Tag(`ipv6Result`), t.IPv6Result.String())              // IPv6地址更新结果: `未改变` `失败` `成功`
	content = strings.ReplaceAll(content, Tag(`ipv6Domains`), strings.Join(t.IPv6Domains, `, `)) // IPv6的域名，多个以`,`分割
	content = strings.ReplaceAll(content, Tag(`error`), t.Error)
	return content
}

var tagNames = []string{`ip`, `ipv4Addr`, `ipv4Result`, `ipv4Domains`, `ipv6Addr`, `ipv6Result`, `ipv6Domains`, `error`}

func (t *TagValues) ParseQuery(urlQuery string) string {
	for _, key := range tagNames {
		urlQuery = strings.ReplaceAll(urlQuery, Tag(key), t.urlEscape(key))
	}
	return urlQuery
}

func (t *TagValues) urlEscape(key string) string {
	switch key {
	case `ip`:
		if len(t.escapedIPAddr) == 0 && len(t.IPAddr) > 0 {
			t.escapedIPAddr = url.QueryEscape(t.IPAddr)
		}
		return t.escapedIPAddr
	case `ipv4Addr`:
		if len(t.escapedIPv4Addr) == 0 && len(t.IPv4Addr) > 0 {
			t.escapedIPv4Addr = url.QueryEscape(t.IPv4Addr)
		}
		return t.escapedIPv4Addr
	case `ipv4Result`:
		if len(t.escapedIPv4Result) == 0 {
			t.escapedIPv4Result = url.QueryEscape(t.IPv4Result.String())
		}
		return t.escapedIPv4Result
	case `ipv4Domains`:
		if len(t.escapedIPv4Domains) == 0 && len(t.IPv4Domains) > 0 {
			t.escapedIPv4Domains = url.QueryEscape(strings.Join(t.IPv4Domains, `, `))
		}
		return t.escapedIPv4Domains
	case `ipv6Addr`:
		if len(t.escapedIPv6Addr) == 0 && len(t.IPv6Addr) > 0 {
			t.escapedIPv6Addr = url.QueryEscape(t.IPv6Addr)
		}
		return t.escapedIPv6Addr
	case `ipv6Result`:
		if len(t.escapedIPv6Result) == 0 {
			t.escapedIPv6Result = url.QueryEscape(t.IPv6Result.String())
		}
		return t.escapedIPv6Result
	case `ipv6Domains`:
		if len(t.escapedIPv6Domains) == 0 && len(t.IPv6Domains) > 0 {
			t.escapedIPv6Domains = url.QueryEscape(strings.Join(t.IPv6Domains, `, `))
		}
		return t.escapedIPv6Domains
	case `error`:
		if len(t.escapedError) == 0 && len(t.Error) > 0 {
			t.escapedError = url.QueryEscape(t.Error)
		}
		return t.escapedError
	}
	return ``
}
