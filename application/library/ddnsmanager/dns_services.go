package ddnsmanager

import (
	"sort"

	"github.com/admpub/nging/v3/application/library/ddnsmanager/interfaces"
)

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

func All() []interfaces.Updater {
	services := map[string]interfaces.Updater{}
	names := make([]string, len(dnsServices))
	for name, c := range dnsServices {
		services[name] = c()
		names = append(names, name)
	}
	sort.Strings(names)
	r := make([]interfaces.Updater, len(names))
	for i, name := range names {
		r[i] = services[name]
	}
	return r
}
