package config

func NewNetInterface() *NetInterface {
	return &NetInterface{
		Filter: &Filter{},
	}
}

type NetInterface struct {
	Type   string // netInterface / api
	Name   string
	Filter *Filter
}

func NewNetIPConfig() *NetIPConfig {
	return &NetIPConfig{
		NetInterface: NewNetInterface(),
	}
}

type NetIPConfig struct {
	Enabled      bool
	NetInterface *NetInterface
	NetIPApiUrl  string
}
