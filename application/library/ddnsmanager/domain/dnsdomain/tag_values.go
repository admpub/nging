package dnsdomain

import (
	"strings"

	"github.com/webx-top/echo"
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

var TagValueDescs = echo.NewKVData().
	Add(Tag(`ipv4Addr`), `新的IPv4地址`).
	Add(Tag(`ipv4Result`), `IPv4地址更新结果(JSON格式)`).
	Add(Tag(`ipv4Domains`), `IPv4的域名，多个以","分割`).
	Add(Tag(`ipv6Addr`), `新的IPv6地址`).
	Add(Tag(`ipv6Result`), `IPv6地址更新结果(JSON格式)`).
	Add(Tag(`ipv6Domains`), `IPv6的域名，多个以","分割`).
	Add(Tag(`error`), `错误信息`).
	Add(Tag(`ip`), `新的IP地址(自动选择IPv4/IPv6)`)
