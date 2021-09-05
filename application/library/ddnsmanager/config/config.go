package config

import (
	"time"

	"github.com/webx-top/echo"
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

type Config struct {
	Closed         bool
	IPv4           *NetIPConfig
	IPv6           *NetIPConfig
	DNSServices    []*DNSService
	DNSResolver    string            // example: 8.8.8.8
	NotifyTemplate map[string]string // 通知模板{html:"",markdown:""}
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
	return nil
}

func (c *Config) AfterValidate(ctx echo.Context) error {
	return nil
}

func (c *Config) FindService(provider string) *DNSService {
	for _, dnsService := range c.DNSServices {
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
		if srv.Enabled {
			if c.IPv4.Enabled && len(srv.IPv4Domains) > 0 {
				return true
			}
			if c.IPv6.Enabled && len(srv.IPv6Domains) > 0 {
				return true
			}
			return true
		}
	}
	return false
}
