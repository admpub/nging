package config

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/admpub/log"
	"github.com/webx-top/echo"

	"github.com/admpub/nging/v4/application/library/webhook"

	"github.com/nging-plugins/ddnsmanager/application/library/ddnsmanager/domain/dnsdomain"
)

// step1. config.Commit()
// step2. domains, err := ParseDomain(conf *config.Config)
// step3. err = domains.Update()
// - step1. err = updater.Init(settings, domains)
// - step2. updater.Update(`A`, newIPv4) / updater.Update(`AAAA`, newIPv6)

func New() *Config {
	return &Config{
		IPv4: NewNetIPConfig(),
		IPv6: NewNetIPConfig(),
		NotifyTemplate: map[string]string{
			"html":     "",
			"markdown": "",
		},
		Interval: 5 * time.Minute,
	}
}

const (
	NotifyDisabled = 0
	NotifyIfError  = 1
	NotifyAll      = 2
)

type Config struct {
	Closed         bool
	IPv4           *NetIPConfig
	IPv6           *NetIPConfig
	DNSServices    []*DNSService
	DNSResolver    string            // example: 8.8.8.8
	NotifyMode     int               // 0-关闭通知; 1-仅仅出错时发送通知；2-发送全部通知
	NotifyTemplate map[string]string // 通知模板{html:"",markdown:""}
	Webhooks       []*webhook.Webhook
	Interval       time.Duration
}

var (
	_ echo.BinderKeyNormalizer = &Config{}
	_ echo.BeforeValidate      = &Config{}
	_ echo.AfterValidate       = &Config{}
)

func (c *Config) BinderKeyNormalizer(key string) string {
	return key
}

func (c *Config) BeforeValidate(ctx echo.Context) error {
	if c.IPv4.Enabled {
		switch c.IPv4.Type {
		case "netInterface":
			if len(c.IPv4.NetInterface.Name) == 0 {
				return errors.New(`请选择用于获取IPv4地址的网卡(如果没有出现选择项，说明不支持)`)
			}
		case "cmd":
			if len(c.IPv4.CommandLine.Command) == 0 {
				return errors.New(`请输入用于获取IPv4地址的可执行命令`)
			}
		default:
			// if len(c.IPv4.NetIPApiUrl) == 0 {
			// 	return errors.New(`请输入用于获取IPv4地址的API接口网址`)
			// }
		}
	}
	if c.IPv6.Enabled {
		switch c.IPv6.Type {
		case "netInterface":
			if len(c.IPv6.NetInterface.Name) == 0 {
				return errors.New(`请选择用于获取IPv6地址的网卡(如果没有出现选择项，说明不支持)`)
			}
		case "cmd":
			if len(c.IPv6.CommandLine.Command) == 0 {
				return errors.New(`请输入用于获取IPv4地址的可执行命令`)
			}
		default:
			// if len(c.IPv6.NetIPApiUrl) == 0 {
			// 	return errors.New(`请输入用于获取IPv4地址的API接口网址`)
			// }
		}
	}
	if !c.Closed {
		for _, dnsService := range c.DNSServices {
			if c.IPv4.Enabled {
				for _, domain := range dnsService.IPv4Domains {
					if domain == nil {
						continue
					}
					if len(domain.Domain) == 0 {
						return fmt.Errorf(`请设置“%s”的IPv4域名`, dnsService.Provider)
					}
				}
			}
			if c.IPv6.Enabled {
				for _, domain := range dnsService.IPv6Domains {
					if domain == nil {
						continue
					}
					if len(domain.Domain) == 0 {
						return fmt.Errorf(`请设置“%s”的IPv6域名`, dnsService.Provider)
					}
				}
			}
		}
	}
	for _, webhook := range c.Webhooks {
		if webhook == nil {
			continue
		}
		if err := webhook.Validate(); err != nil {
			return err
		}
	}
	return nil
}

func (c *Config) AfterValidate(ctx echo.Context) error {
	return nil
}

func (c *Config) FindService(provider string) *DNSService {
	for _, dnsService := range c.DNSServices {
		if dnsService == nil {
			continue
		}
		if dnsService.Provider == provider {
			return dnsService
		}
	}
	return nil
}

func (c *Config) Commit() error {
	if err := c.IPv4.NetInterface.Filter.Init(); err != nil {
		return err
	}
	return c.IPv6.NetInterface.Filter.Init()
}

func (c *Config) IsValid() bool {
	if c.Closed {
		return false
	}
	if !c.IPv4.Enabled && !c.IPv6.Enabled {
		return false
	}
	for _, srv := range c.DNSServices {
		if srv == nil {
			continue
		}
		if srv.Enabled {
			if c.IPv4.Enabled && len(srv.IPv4Domains) > 0 {
				return true
			}
			if c.IPv6.Enabled && len(srv.IPv6Domains) > 0 {
				return true
			}
		}
	}
	for _, webhook := range c.Webhooks {
		if webhook == nil {
			continue
		}
		if len(webhook.Url) > 0 {
			return true
		}
	}
	return false
}

func (c *Config) HasWebhook() bool {
	if len(c.Webhooks) == 0 {
		return false
	}
	for _, webhook := range c.Webhooks {
		if webhook == nil {
			continue
		}
		return true
	}
	return false
}

func (c *Config) ExecWebhooks(tagValues *dnsdomain.TagValues) error {
	var retErr error
	var errs []string
	for _, webhook := range c.Webhooks {
		if webhook == nil {
			continue
		}
		err := webhook.Exec(tagValues.Parse, tagValues.ParseQuery)
		if err != nil {
			log.Error(err)
			errs = append(errs, err.Error())
		}
	}
	if len(errs) > 0 {
		retErr = errors.New(strings.Join(errs, "\n"))
	}
	return retErr
}
