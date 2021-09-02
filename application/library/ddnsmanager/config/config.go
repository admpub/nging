package config

// step1. config
// step2. domains := ParseDomain(conf *config.Config)
// step3. updater.Init(settings, domains)
// step4. updater.Update(`A`) / updater.Update(`AAAA`)

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
