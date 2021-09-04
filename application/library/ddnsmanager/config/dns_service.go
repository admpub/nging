package config

import "github.com/webx-top/echo"

type DNSDomain struct {
	IPFormat string
	Domain   string
	Extra    echo.H
}

type DNSService struct {
	Provider    string
	Settings    echo.H
	IPv4Domains []*DNSDomain
	IPv6Domains []*DNSDomain
}
