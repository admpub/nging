package captcha

import "sort"

var drivers = map[string]func() ICaptcha{
	`default`: func() ICaptcha { return dflt },
	`api`:     newCaptchaAPI,
}

func Register(name string, ic func() ICaptcha) {
	drivers[name] = ic
}

func Get(name string) func() ICaptcha {
	return drivers[name]
}

func GetOk(name string) (func() ICaptcha, bool) {
	ic, ok := drivers[name]
	return ic, ok
}

func GetAllNames() []string {
	names := make([]string, 0, len(drivers))
	for name := range drivers {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

func Has(name string) bool {
	_, ok := drivers[name]
	return ok
}
