package domain

import (
	"fmt"
	"log"
	"strings"

	"github.com/admpub/nging/v3/application/library/ddnsmanager/config"
	"github.com/admpub/nging/v3/application/library/ddnsmanager/utils"
	"golang.org/x/net/publicsuffix"
)

// UpdateStatusType 更新状态
type UpdateStatusType string

const (
	// UpdatedNothing 未改变
	UpdatedNothing UpdateStatusType = "未改变"
	// UpdatedFailed 更新失败
	UpdatedFailed UpdateStatusType = "失败"
	// UpdatedSuccess 更新成功
	UpdatedSuccess UpdateStatusType = "成功"
	UpdatedIdle    UpdateStatusType = ""
)

// Domains Ipv4/Ipv6 domains
type Domains struct {
	IPv4Addr    string
	IPv4Domains map[string][]*Domain // {dnspod:[]}
	IPv6Addr    string
	IPv6Domains map[string][]*Domain // {dnspod:[]}
}

// Domain 域名实体
type Domain struct {
	Port         int
	DomainName   string
	SubDomain    string
	UpdateStatus UpdateStatusType // 更新状态
}

func (d Domain) String() string {
	if len(d.SubDomain) > 0 {
		return d.SubDomain + "." + d.DomainName
	}
	return d.DomainName
}

func (d Domain) IP(ip string) string {
	if d.Port > 0 {
		return fmt.Sprintf(`%s:%d`, ip, d.Port)
	}
	return ip
}

// GetFullDomain 获得全部的，子域名
func (d Domain) GetFullDomain() string {
	if len(d.SubDomain) > 0 {
		return d.SubDomain + "." + d.DomainName
	}
	return "@." + d.DomainName
}

// GetSubDomain 获得子域名，为空返回@
// 阿里云，dnspod需要
func (d Domain) GetSubDomain() string {
	if len(d.SubDomain) > 0 {
		return d.SubDomain
	}
	return "@"
}

func NewDomains() *Domains {
	return &Domains{
		IPv4Domains: map[string][]*Domain{},
		IPv6Domains: map[string][]*Domain{},
	}
}

// ParseDomain 接口获得ip并校验用户输入的域名
func ParseDomain(conf *config.Config) *Domains {
	domains := NewDomains()
	// IPv4
	ipv4Addr := utils.GetIPv4Addr(conf.IPv4.NetInterface, conf.IPv4.NetIPApiUrl)
	if len(ipv4Addr) > 0 {
		domains.IPv4Addr = ipv4Addr
		for _, service := range conf.DNSServices {
			_, ok := domains.IPv4Domains[service.Provider]
			if !ok {
				domains.IPv4Domains[service.Provider] = []*Domain{}
			}
			domains.IPv4Domains[service.Provider] = parseDomainArr(service.IPv4Domains)
		}
	}
	// IPv6
	ipv6Addr := utils.GetIPv6Addr(conf.IPv6.NetInterface, conf.IPv6.NetIPApiUrl)
	if len(ipv6Addr) > 0 {
		domains.IPv6Addr = ipv6Addr
		for _, service := range conf.DNSServices {
			_, ok := domains.IPv6Domains[service.Provider]
			if !ok {
				domains.IPv6Domains[service.Provider] = []*Domain{}
			}
			domains.IPv6Domains[service.Provider] = parseDomainArr(service.IPv6Domains)
		}
	}
	return domains
}

// parseDomainArr 校验用户输入的域名
func parseDomainArr(dnsDomains []*config.DNSDomain) (domains []*Domain) {
	for _, dnsDomain := range dnsDomains {
		_domain := strings.TrimSpace(dnsDomain.Domain)
		if len(_domain) == 0 {
			continue
		}
		domain := &Domain{
			Port:         dnsDomain.Port,
			UpdateStatus: UpdatedIdle,
		}
		sp := strings.Split(_domain, ".")
		length := len(sp)
		if length <= 1 {
			log.Println(_domain, "域名不正确")
			continue
		}
		// 处理域名
		topLevelDomain, err := publicsuffix.EffectiveTLDPlusOne(_domain)
		if err != nil {
			log.Println(_domain, "域名不正确", err)
			continue
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

// ParseDomainResult 获得ParseDomain结果
func (domains *Domains) ParseDomainResult(recordType string, provider string) (ipAddr string, retDomains []*Domain) {
	if recordType == "AAAA" {
		return domains.IPv6Addr, domains.IPv6Domains[provider]
	}
	return domains.IPv4Addr, domains.IPv4Domains[provider]
}
