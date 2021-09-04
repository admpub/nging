package config

func NewNetInterface() *NetInterface {
	return &NetInterface{
		Filter: &Filter{},
	}
}

type NetInterface struct {
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
	Type         string // netInterface / api
	NetInterface *NetInterface
	NetIPApiUrl  string
}
