package domain

import (
	"context"
	"errors"
	"strings"

	"github.com/admpub/log"
	"github.com/nging-plugins/ddnsmanager/application/library/ddnsmanager"
	"github.com/nging-plugins/ddnsmanager/application/library/ddnsmanager/config"
	"github.com/nging-plugins/ddnsmanager/application/library/ddnsmanager/domain/dnsdomain"
	"github.com/nging-plugins/ddnsmanager/application/library/ddnsmanager/resolver"
	"github.com/nging-plugins/ddnsmanager/application/library/ddnsmanager/sender"
	"github.com/nging-plugins/ddnsmanager/application/library/ddnsmanager/utils"
)

func (domains *Domains) SetIPv4Addr(ipv4Addr string) {
	domains.IPv4Addr = ipv4Addr
	domains.SaveIP(4)
}

func (domains *Domains) updateIPv4(ctx context.Context, conf *config.Config, ipv4Addr string) (ipv4Changed bool, errs []error) {
	if len(ipv4Addr) == 0 || domains.IPv4Addr == ipv4Addr {
		return
	}
	log.Debugf(`[DDNS] 查询到ipv4变更: %s => %s`, domains.IPv4Addr, ipv4Addr)
	domains.SetIPv4Addr(ipv4Addr)
	for dnsProvider, dnsDomains := range domains.IPv4Domains {
		var _dnsDomains []*dnsdomain.Domain
		for _, dnsDomain := range dnsDomains {
			if dnsDomain == nil {
				continue
			}
			oldIP, err := resolver.ResolveDNS(dnsDomain.String(), conf.DNSResolver, `IPV4`)
			if err != nil {
				log.Errorf("[%s] ResolveDNS(%s): %s", dnsProvider, dnsDomain.String(), err.Error())
				//errs = append(errs, err)
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
	return
}

func (domains *Domains) SetIPv6Addr(ipv6Addr string) {
	domains.IPv6Addr = ipv6Addr
	domains.SaveIP(6)
}

func (domains *Domains) updateIPv6(ctx context.Context, conf *config.Config, ipv6Addr string) (ipv6Changed bool, errs []error) {
	if len(ipv6Addr) == 0 || domains.IPv6Addr == ipv6Addr {
		return
	}
	log.Debugf(`[DDNS] 查询到ipv6变更: %s => %s`, domains.IPv6Addr, ipv6Addr)
	domains.SetIPv6Addr(ipv6Addr)
	for dnsProvider, dnsDomains := range domains.IPv6Domains {
		var _dnsDomains []*dnsdomain.Domain
		for _, dnsDomain := range dnsDomains {
			if dnsDomain == nil {
				continue
			}
			oldIP, err := resolver.ResolveDNS(dnsDomain.String(), conf.DNSResolver, `IPV6`)
			if err != nil {
				log.Errorf("[%s] ResolveDNS(%s): %s", dnsProvider, dnsDomain.String(), err.Error())
				//errs = append(errs, err)
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
	return
}

func (domains *Domains) Update(ctx context.Context, conf *config.Config) error {

	var (
		errs        []error
		ipv4Changed bool
		ipv6Changed bool
	)

	// IPv4
	if conf.IPv4.Enabled {
		ipv4Addr, err := utils.GetIPv4Addr(conf.IPv4)
		if err != nil {
			log.Error(err)
		} else {
			ipv4Changed, errs = domains.updateIPv4(ctx, conf, ipv4Addr)
		}
	}
	// IPv6
	if conf.IPv6.Enabled {
		ipv6Addr, err := utils.GetIPv6Addr(conf.IPv6)
		if err != nil {
			log.Error(err)
		} else {
			var _errs []error
			ipv6Changed, _errs = domains.updateIPv6(ctx, conf, ipv6Addr)
			if len(_errs) > 0 {
				if len(errs) > 0 {
					errs = append(errs, _errs...)
				} else {
					errs = _errs
				}
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
		if err != nil {
			t.Error = err.Error()
		}
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
	if err := sender.Send(*tagValues(), conf.NotifyTemplate); err != nil {
		log.Errorf("[DDNS] sender.Send - %v", err)
	}
	return err
}
