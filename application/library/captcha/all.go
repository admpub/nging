package captcha

import "sort"

const (
	TypeDefault = `default`
	TypeAPI     = `api`
)

var drivers = map[string]func() ICaptcha{
	TypeDefault: func() ICaptcha { return dflt },
	TypeAPI:     newCaptchaAPI,
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
