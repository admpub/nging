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
}
