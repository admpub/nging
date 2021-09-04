package config

import "github.com/webx-top/echo"

type DNSDomain struct {
	IPFormat string
	Domain   string
	Extra    echo.H
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
	c.IPv4Domains = make([]*DNSDomain, len(d.IPv4Domains))
	c.IPv6Domains = make([]*DNSDomain, len(d.IPv6Domains))
	for i, v := range d.IPv4Domains {
		_v := *v
		c.IPv4Domains[i] = &_v
	}
	for i, v := range d.IPv6Domains {
		_v := *v
		c.IPv6Domains[i] = &_v
	}
	if d.Settings != nil {
		c.Settings = d.Settings.Clone()
	}
	return &c
}
