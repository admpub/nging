package config

type NetIPConfig struct {
	InterfaceType string // netInterface / api
	InterfaceName string
	NetIPApiUrl   string
	Domains       []string
	Enabled       bool
}

type Config struct {
	IPv4 NetIPConfig
	IPv6 NetIPConfig
}
