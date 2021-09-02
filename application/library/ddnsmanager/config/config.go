package config

// step1. config.Commit()
// step2. domains, err := ParseDomain(conf *config.Config)
// step3. err = domains.Update()
// - step1. err = updater.Init(settings, domains)
// - step2. updater.Update(`A`, newIPv4) / updater.Update(`AAAA`, newIPv6)

func New() *Config {
	return &Config{
		IPv4: NewNetIPConfig(),
		IPv6: NewNetIPConfig(),
	}
}

type Config struct {
	IPv4        *NetIPConfig
	IPv6        *NetIPConfig
	DNSServices []*DNSService
	DNSResolver string // example: 8.8.8.8
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
