package ip2region

import (
	"regexp"

	"github.com/admpub/nging/v5/application/library/config/extend"
)

type WANIPProvider struct {
	Name        string
	Description string
	URL         string
	Method      string
	IP4Rule     string // 如果是IPv4规则，此项可以为空
	IP6Rule     string // 如果是IPv6规则，此项不能为空（“=”代表整个body数据）
	ip4regexp   *regexp.Regexp
	ip6regexp   *regexp.Regexp
	Disabled    bool
}

func (w *WANIPProvider) HasIPv6Rule() bool {
	return len(w.IP6Rule) > 0
}

func (w *WANIPProvider) HasIPv4Rule() bool {
	return len(w.IP4Rule) > 0 || !w.HasIPv6Rule()
}

type WANIPProviders map[string]*WANIPProvider

func (w *WANIPProviders) Reload() error {
	if w == nil {
		return nil
	}
	for key, value := range *w {
		if value != nil && len(value.Name) > 0 && len(value.URL) > 0 {
			Register(value)
		} else {
			Unregister(key)
		}
	}
	return nil
}

var (
	wanIPProviders   = map[string]*WANIPProvider{}
	defaultProviders = []*WANIPProvider{
		// IPv4
		{
			Name:        `oray.com`,
			Description: `oray`,
			URL:         `https://ddns.oray.com/checkip`,
			IP4Rule:     IPv4Rule,
		}, {
			Name:        `ip-api.com`,
			Description: ``,
			URL:         `http://ip-api.com/json/?fields=query`,
			IP4Rule:     `"query":"` + IPv4Rule + `"`,
		}, {
			Name:        `myip.la`,
			Description: ``,
			URL:         `https://api.myip.la/`,
			IP4Rule:     ``,
		}, {
			Name:        `ipify.org`,
			Description: ``,
			URL:         `https://api.ipify.org`,
			IP4Rule:     ``,
		}, {
			Name:        `3322.org`,
			Description: ``,
			URL:         `http://members.3322.org/dyndns/getip`,
			IP4Rule:     ``,
		}, {
			Name:        `ip.sb`,
			Description: ``,
			URL:         `https://api.ip.sb/ip`,
			IP4Rule:     ``,
		}, {
			Name:        `ipconfig.io`,
			Description: ``,
			URL:         `https://ipconfig.io/ip`,
			IP4Rule:     ``,
		},
		// IPv6
		{
			Name:        `ident.me`,
			Description: ``,
			URL:         `https://v6.ident.me`,
			IP6Rule:     `=`,
		}, {
			Name:        `api-ipv6.ip.sb`,
			Description: ``,
			URL:         `https://api-ipv6.ip.sb/ip`,
			IP6Rule:     `=`,
		}, {
			Name:        `v6.myip.la`,
			Description: ``,
			URL:         `https://v6.myip.la/`,
			IP6Rule:     `=`,
		},
	}
)

func init() {
	extend.Register(`wanIPProvider`, func() interface{} {
		return &WANIPProviders{} // 更新时会自动调用 WANIPProviders.Reload()
	})
	for _, provider := range defaultProviders {
		if err := Register(provider); err != nil {
			panic(err)
		}
	}
}

func Register(p *WANIPProvider) (err error) {
	if len(p.IP4Rule) > 0 && p.IP4Rule != `=` {
		p.ip4regexp, err = regexp.Compile(p.IP4Rule)
		if err != nil {
			return
		}
	}
	if len(p.IP6Rule) > 0 && p.IP6Rule != `=` {
		p.ip6regexp, err = regexp.Compile(p.IP6Rule)
		if err != nil {
			return
		}
	}
	if len(p.Method) == 0 {
		p.Method = `GET`
	}
	wanIPProviders[p.Name] = p
	return
}

func Get(name string) *WANIPProvider {
	p, _ := wanIPProviders[name]
	return p
}

func Unregister(names ...string) {
	for _, name := range names {
		delete(wanIPProviders, name)
	}
}
