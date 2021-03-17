package cmder

var cmders = map[string]Cmder{}

func Register(name string, cmder Cmder) {
	cmders[name] = cmder
}

func Get(name string) Cmder {
	cmder, _ := cmders[name]
	return cmder
}

func Unregister(name string) {
	if _, ok := cmders[name]; ok {
		delete(cmders, name)
	}
}
