package config

import "github.com/webx-top/echo"

type DNSDomain struct {
	Port   int
	Domain string
}

type DNSService struct {
	Provider    string
	Settings    echo.H
	IPv4Domains []*DNSDomain
	IPv6Domains []*DNSDomain
}
