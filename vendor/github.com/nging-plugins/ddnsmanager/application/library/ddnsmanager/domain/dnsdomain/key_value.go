package dnsdomain

import "github.com/webx-top/echo"

var TagValueDescs = echo.NewKVData().
	Add(Tag(`ipv4Addr`), `新的IPv4地址`).
	Add(Tag(`ipv4Result`), `IPv4地址更新结果(JSON格式)`).
	Add(Tag(`ipv4Domains`), `IPv4的域名，多个以","分割`).
	Add(Tag(`ipv6Addr`), `新的IPv6地址`).
	Add(Tag(`ipv6Result`), `IPv6地址更新结果(JSON格式)`).
	Add(Tag(`ipv6Domains`), `IPv6的域名，多个以","分割`).
	Add(Tag(`error`), `错误信息`).
	Add(Tag(`ip`), `新的IP地址(自动选择IPv4/IPv6)`)

var TrackerTypes = echo.NewKVData().
	Add(`api`, `通过接口获取`).
	Add(`netInterface`, `通过网卡获取`).
	Add(`cmd`, `通过命令行获取`)
