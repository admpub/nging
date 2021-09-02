package ddnsmanager

import "github.com/admpub/nging/v3/application/library/ddnsmanager/interfaces"

var dnsServices = map[string]func() interfaces.Updater{}

func Register(provider string, updater func() interfaces.Updater) {
	dnsServices[provider] = updater
}

func Open(provider string) interfaces.Updater {
	c, ok := dnsServices[provider]
	if ok {
		return c()
	}
	return nil
}
