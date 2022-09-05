package captcha

var drivers = map[string]ICaptcha{
	`default`: dflt,
}

func Register(name string, ic ICaptcha) {
	drivers[name] = ic
}

func Get(name string) ICaptcha {
	return drivers[name]
}

func GetOk(name string) (ICaptcha, bool) {
	ic, ok := drivers[name]
	return ic, ok
}

func Has(name string) bool {
	_, ok := drivers[name]
	return ok
}
