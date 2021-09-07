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
		CommandLine:  &CommandLine{},
	}
}

type NetIPConfig struct {
	Enabled      bool
	Type         string // netInterface / api / cmd
	NetInterface *NetInterface
	NetIPApiUrl  string
	CommandLine  *CommandLine
}
