package domain

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/admpub/log"
	"github.com/admpub/nging/v3/application/library/ddnsmanager"
	"github.com/admpub/nging/v3/application/library/ddnsmanager/config"
	"github.com/admpub/nging/v3/application/library/ddnsmanager/domain/dnsdomain"
	"github.com/admpub/nging/v3/application/library/ddnsmanager/resolver"
	"github.com/admpub/nging/v3/application/library/ddnsmanager/sender"
	"github.com/admpub/nging/v3/application/library/ddnsmanager/utils"
	"github.com/webx-top/echo"
	"golang.org/x/net/publicsuffix"
)

func NewDomains() *Domains {
	return &Domains{
		IPv4Domains: map[string][]*dnsdomain.Domain{},
		IPv6Domains: map[string][]*dnsdomain.Domain{},
	}
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

func (domains *Domains) TagValues(ipv4Changed, ipv6Changed bool) *dnsdomain.TagValues {
	t := dnsdomain.NewTagValues()
	if ipv4Changed {
		t.IPv4Addr = domains.IPv4Addr
		for _provider, _domains := range domains.IPv4Domains {
			for _, _domain := range _domains {
				t.IPv4Domains = append(t.IPv4Domains, _domain.String())
				t.IPv4Result.Add(_provider, _domain.Result())
			}
		}
	}
	if ipv6Changed {
		t.IPv6Addr = domains.IPv6Addr
		for _provider, _domains := range domains.IPv6Domains {
			for _, _domain := range _domains {
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

func (domains *Domains) Update(ctx context.Context, conf *config.Config) error {

	var (
		errs        []error
		ipv4Changed bool
		ipv6Changed bool
	)

	// IPv4
	if conf.IPv4.Enabled {
		ipv4Addr := utils.GetIPv4Addr(conf.IPv4)
		if len(ipv4Addr) > 0 && domains.IPv4Addr != ipv4Addr {
			log.Debugf(`[DDNS] 查询到ipv4变更: %s => %s`, domains.IPv4Addr, ipv4Addr)
			domains.IPv4Addr = ipv4Addr
			for dnsProvider, dnsDomains := range domains.IPv4Domains {
				var _dnsDomains []*dnsdomain.Domain
				for _, dnsDomain := range dnsDomains {
					if dnsDomain == nil {
						continue
					}
					oldIP, err := resolver.ResolveDNS(dnsDomain.String(), conf.DNSResolver, `IPV4`)
					if err != nil {
						log.Errorf("[%s] ResolveDNS(%s): %s", dnsProvider, dnsDomain.String(), err.Error())
						errs = append(errs, err)
						dnsDomain.UpdateStatus = dnsdomain.UpdatedIdle
						_dnsDomains = append(_dnsDomains, dnsDomain)
						continue
					}
					if oldIP != ipv4Addr {
						dnsDomain.UpdateStatus = dnsdomain.UpdatedIdle
						_dnsDomains = append(_dnsDomains, dnsDomain)
						continue
					}
					dnsDomain.UpdateStatus = dnsdomain.UpdatedNothing
					log.Infof("[%s] IP is the same as cached one (%s). Skip update (%s)", dnsProvider, ipv4Addr, dnsDomain.String())
				}
				if len(_dnsDomains) == 0 {
					continue
				}
				updater := ddnsmanager.Open(dnsProvider)
				if updater == nil {
					continue
				}
				dnsService := conf.FindService(dnsProvider)
				err := updater.Init(dnsService.Settings, _dnsDomains)
				if err != nil {
					errs = append(errs, err)
					continue
				}
				log.Infof("[%s] %s - Start to update record IP...", dnsProvider, ipv4Addr)
				err = updater.Update(ctx, `A`, ipv4Addr)
				if err != nil {
					errs = append(errs, err)
				}
				ipv4Changed = true
			}
		}
	}
	// IPv6
	if conf.IPv6.Enabled {
		ipv6Addr := utils.GetIPv6Addr(conf.IPv6)
		if len(ipv6Addr) > 0 && domains.IPv6Addr != ipv6Addr {
			log.Debugf(`[DDNS] 查询到ipv6变更: %s => %s`, domains.IPv6Addr, ipv6Addr)
			domains.IPv6Addr = ipv6Addr
			for dnsProvider, dnsDomains := range domains.IPv6Domains {
				var _dnsDomains []*dnsdomain.Domain
				for _, dnsDomain := range dnsDomains {
					if dnsDomain == nil {
						continue
					}
					oldIP, err := resolver.ResolveDNS(dnsDomain.String(), conf.DNSResolver, `IPV6`)
					if err != nil {
						log.Errorf("[%s] ResolveDNS(%s): %s", dnsProvider, dnsDomain.String(), err.Error())
						errs = append(errs, err)
						dnsDomain.UpdateStatus = dnsdomain.UpdatedIdle
						_dnsDomains = append(_dnsDomains, dnsDomain)
						continue
					}
					if oldIP != ipv6Addr {
						dnsDomain.UpdateStatus = dnsdomain.UpdatedIdle
						_dnsDomains = append(_dnsDomains, dnsDomain)
						continue
					}
					dnsDomain.UpdateStatus = dnsdomain.UpdatedNothing
					log.Infof("[%s] IP is the same as cached one (%s). Skip update (%s)", dnsProvider, ipv6Addr, dnsDomain.String())
				}
				if len(_dnsDomains) == 0 {
					continue
				}
				updater := ddnsmanager.Open(dnsProvider)
				if updater == nil {
					continue
				}
				dnsService := conf.FindService(dnsProvider)
				err := updater.Init(dnsService.Settings, _dnsDomains)
				if err != nil {
					errs = append(errs, err)
					continue
				}
				log.Infof("[%s] %s - Start to update record IP...", dnsProvider, ipv6Addr)
				err = updater.Update(ctx, `AAAA`, ipv6Addr)
				if err != nil {
					errs = append(errs, err)
				}
				ipv6Changed = true
			}
		}
	}
	if !conf.IPv4.Enabled && !conf.IPv6.Enabled {
		return nil
	}
	var err error
	if len(errs) > 0 {
		errMessages := make([]string, len(errs))
		for index, err := range errs {
			errMessages[index] = err.Error()
		}
		err = errors.New(strings.Join(errMessages, "\n"))
	}
	var t *dnsdomain.TagValues
	tagValues := func() *dnsdomain.TagValues {
		if t != nil {
			return t
		}
		t = domains.TagValues(ipv4Changed, ipv6Changed)
		return t
	}
	if conf.HasWebhook() {
		if err := conf.ExecWebhooks(tagValues()); err != nil {
			log.Errorf("[DDNS] webhook - %v", err)
		}
	}
	switch conf.NotifyMode {
	case config.NotifyDisabled:
		return err
	case config.NotifyIfError:
		if err == nil {
			return err
		}
	case config.NotifyAll:
	}
	if err != nil {
		t.Error = err.Error()
	}
	if err := sender.Send(*tagValues(), conf.NotifyTemplate); err != nil {
		log.Errorf("[DDNS] sender.Send - %v", err)
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
