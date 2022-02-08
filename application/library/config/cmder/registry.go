package cmder

var cmders = map[string]Cmder{}

func Register(name string, cmder Cmder) {
	cmders[name] = cmder
}

func Get(name string) Cmder {
	cmder := cmders[name]
	return cmder
}

func Has(name string) bool {
	_, ok := cmders[name]
	return ok
}

func Unregister(name string) {
	delete(cmders, name)
}
