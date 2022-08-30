package ddnsmanager

import (
	"sort"

	"github.com/nging-plugins/ddnsmanager/application/library/ddnsmanager/config"
	"github.com/nging-plugins/ddnsmanager/application/library/ddnsmanager/domain/dnsdomain"
	"github.com/nging-plugins/ddnsmanager/application/library/ddnsmanager/interfaces"
	"github.com/webx-top/echo"
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

type ProviderMeta struct {
	Name        string
	Description string
	SignUpURL   string
	LineTypeURL string
	Support     dnsdomain.Support
	ConfigItems echo.KVList
	DNSService  *config.DNSService
}

func AllProvoderMeta(cfgServices []*config.DNSService) []*ProviderMeta {
	services := map[string]interfaces.Updater{}
	var names []string
	for name, c := range dnsServices {
		services[name] = c()
		names = append(names, name)
	}
	sort.Strings(names)
	r := make([]*ProviderMeta, len(names))
	for i, name := range names {
		sv := services[name]
		mt := &ProviderMeta{
			Name:        sv.Name(),
			Description: sv.Description(),
			SignUpURL:   sv.SignUpURL(),
			LineTypeURL: sv.LineTypeURL(),
			Support:     sv.Support(),
			ConfigItems: sv.ConfigItems(),
		}
		for _, cfgSrv := range cfgServices {
			if cfgSrv.Provider == name {
				mt.DNSService = cfgSrv.Clone()
				break
			}
		}
		if mt.DNSService == nil {
			mt.DNSService = &config.DNSService{
				Provider: name,
			}
		}
		if mt.DNSService.Settings == nil {
			mt.DNSService.Settings = echo.H{}
		}
		r[i] = mt
	}
	return r
}
