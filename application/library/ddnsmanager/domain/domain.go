package domain

import (
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
	IPv4Domains []*Domain
	IPv6Addr    string
	IPv6Domains []*Domain
}

// Domain 域名实体
type Domain struct {
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
	return &Domains{}
}

// ParseDomain 接口获得ip并校验用户输入的域名
func (domains *Domains) ParseDomain(conf *config.Config) *Domains {
	// IPv4
	ipv4Addr := utils.GetIPv4Addr(conf.IPv4.InterfaceType, conf.IPv4.InterfaceName, conf.IPv4.NetIPApiUrl)
	if len(ipv4Addr) > 0 {
		domains.IPv4Addr = ipv4Addr
		domains.IPv4Domains = parseDomainArr(conf.IPv4.Domains)
	}
	// IPv6
	ipv6Addr := utils.GetIPv6Addr(conf.IPv4.InterfaceType, conf.IPv4.InterfaceName, conf.IPv4.NetIPApiUrl)
	if len(ipv6Addr) > 0 {
		domains.IPv6Addr = ipv6Addr
		domains.IPv6Domains = parseDomainArr(conf.IPv6.Domains)
	}
	return domains
}

// parseDomainArr 校验用户输入的域名
func parseDomainArr(domainArr []string) (domains []*Domain) {
	for _, _domain := range domainArr {
		_domain = strings.TrimSpace(_domain)
		if len(_domain) == 0 {
			continue
		}
		domain := &Domain{
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
func (domains *Domains) ParseDomainResult(recordType string) (ipAddr string, retDomains []*Domain) {
	if recordType == "AAAA" {
		return domains.IPv6Addr, domains.IPv6Domains
	}
	return domains.IPv4Addr, domains.IPv4Domains
}
