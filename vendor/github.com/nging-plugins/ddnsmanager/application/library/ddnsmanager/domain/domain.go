package domain

import (
	"fmt"
	"strings"

	"github.com/webx-top/echo"
	"golang.org/x/net/publicsuffix"

	"github.com/admpub/nging/v4/application/library/common"
	"github.com/nging-plugins/ddnsmanager/application/library/ddnsmanager/config"
	"github.com/nging-plugins/ddnsmanager/application/library/ddnsmanager/domain/dnsdomain"
)

func NewDomains() *Domains {
	d := &Domains{
		IPv4Domains: map[string][]*dnsdomain.Domain{},
		IPv6Domains: map[string][]*dnsdomain.Domain{},
	}
	d.RestoreIP()
	return d
}

// ParseDomain 接口获得ip并校验用户输入的域名
func ParseDomain(conf *config.Config) (*Domains, error) {
	domains := NewDomains()
	err := domains.Init(conf)
	return domains, err
}

// Domains Ipv4/Ipv6 domains
type Domains struct {
	IPv4Addr    string
	IPv4Domains map[string][]*dnsdomain.Domain // {dnspod:[]}
	IPv6Addr    string
	IPv6Domains map[string][]*dnsdomain.Domain // {dnspod:[]}
}

func (domains *Domains) RestoreIP() error {
	b, err := common.ReadCache(`ip`, `ddns_ipv4`)
	if err != nil {
		return err
	}
	domains.IPv4Addr = string(b)
	b, err = common.ReadCache(`ip`, `ddns_ipv6`)
	if err != nil {
		return err
	}
	domains.IPv6Addr = string(b)
	return nil
}

func (domains *Domains) SaveIP(ver int) error {
	switch ver {
	case 4:
		return common.WriteCache(`ip`, `ddns_ipv4`, []byte(domains.IPv4Addr))
	case 6:
		return common.WriteCache(`ip`, `ddns_ipv6`, []byte(domains.IPv6Addr))
	default:
		return nil
	}
}

func (domains *Domains) TagValues(ipv4Changed, ipv6Changed bool) *dnsdomain.TagValues {
	t := dnsdomain.NewTagValues()
	if ipv4Changed {
		t.IPv4Addr = domains.IPv4Addr
		for _provider, _domains := range domains.IPv4Domains {
			for _, _domain := range _domains {
				if _domain == nil {
					continue
				}
				t.IPv4Domains = append(t.IPv4Domains, _domain.String())
				t.IPv4Result.Add(_provider, _domain.Result())
			}
		}
	}
	if ipv6Changed {
		t.IPv6Addr = domains.IPv6Addr
		for _provider, _domains := range domains.IPv6Domains {
			for _, _domain := range _domains {
				if _domain == nil {
					continue
				}
				t.IPv6Domains = append(t.IPv6Domains, _domain.String())
				t.IPv6Result.Add(_provider, _domain.Result())
			}
		}
	}
	t.IPAddr = domains.IPv6Addr
	if len(t.IPAddr) == 0 || len(t.IPv6Domains) == 0 {
		t.IPAddr = domains.IPv4Addr
	}
	return t
}

func (domains *Domains) Init(conf *config.Config) error {
	var err error
	for _, service := range conf.DNSServices {
		if service == nil {
			continue
		}
		if !service.Enabled {
			continue
		}
		if service.Settings == nil {
			service.Settings = echo.H{}
		}
		_, ok := domains.IPv6Domains[service.Provider]
		if !ok {
			domains.IPv6Domains[service.Provider] = []*dnsdomain.Domain{}
		}
		domains.IPv6Domains[service.Provider], err = parseDomainArr(service.IPv6Domains)
		if err != nil {
			return err
		}
		_, ok = domains.IPv4Domains[service.Provider]
		if !ok {
			domains.IPv4Domains[service.Provider] = []*dnsdomain.Domain{}
		}
		domains.IPv4Domains[service.Provider], err = parseDomainArr(service.IPv4Domains)
	}
	return err
}

// parseDomainArr 校验用户输入的域名
func parseDomainArr(dnsDomains []*config.DNSDomain) (domains []*dnsdomain.Domain, err error) {
	for _, dnsDomain := range dnsDomains {
		_domain := strings.TrimSpace(dnsDomain.Domain)
		if len(_domain) == 0 {
			continue
		}
		domain := &dnsdomain.Domain{
			IPFormat:     dnsDomain.IPFormat,
			UpdateStatus: dnsdomain.UpdatedIdle,
			Line:         dnsDomain.Line,
		}
		if dnsDomain.Extra != nil {
			domain.Extra = dnsDomain.Extra
		} else {
			domain.Extra = echo.H{}
		}
		sp := strings.Split(_domain, ".")
		length := len(sp)
		if length <= 1 {
			err = fmt.Errorf(`域名不正确: %s`, _domain)
			return
		}
		var topLevelDomain string
		// 处理域名
		topLevelDomain, err = publicsuffix.EffectiveTLDPlusOne(_domain)
		if err != nil {
			err = fmt.Errorf(`域名不正确: %w`, err)
			return
		}
		domain.DomainName = topLevelDomain
		domainLen := len(_domain) - len(domain.DomainName)
		if domainLen > 0 {
			domain.SubDomain = _domain[:domainLen-1]
		} else {
			domain.SubDomain = _domain[:domainLen]
		}
		domains = append(domains, domain)
	}
	return
}
