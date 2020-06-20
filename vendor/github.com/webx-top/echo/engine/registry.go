package engine

var engines = map[string]func(*Config) Engine{}

func Register(name string, newEngine func(*Config) Engine) {
	engines[name] = newEngine
}

func Get(name string) func(*Config) Engine {
	return engines[name]
}

func New(name string, config *Config) Engine {
	return Get(name)(config)
}
