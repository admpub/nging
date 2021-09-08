package config

import "github.com/webx-top/echo"

type DNSDomain struct {
	IPFormat string // IP格式模板(支持变量标签#{ip})
	Domain   string // 域名
	Line     string // 线路类型
	Extra    echo.H // 扩展数据
}

type DNSService struct {
	Enabled     bool
	Provider    string
	Settings    echo.H
	IPv4Domains []*DNSDomain
	IPv6Domains []*DNSDomain
}

func (d *DNSService) Clone() *DNSService {
	c := *d
	c.IPv4Domains = make([]*DNSDomain, 0, len(d.IPv4Domains))
	c.IPv6Domains = make([]*DNSDomain, 0, len(d.IPv6Domains))
	for _, v := range d.IPv4Domains {
		if v == nil {
			continue
		}
		_v := *v
		c.IPv4Domains = append(c.IPv4Domains, &_v)
	}
	for _, v := range d.IPv6Domains {
		if v == nil {
			continue
		}
		_v := *v
		c.IPv6Domains = append(c.IPv6Domains, &_v)
	}
	if d.Settings != nil {
		c.Settings = d.Settings.Clone()
	}
	return &c
}
