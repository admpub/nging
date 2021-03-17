package extend

type Initer func() interface{}

var extendIniters = map[string]Initer{}

func Register(name string, initer Initer) {
	extendIniters[name] = initer
}

func Range(f func(string, interface{})) {
	for name, initer := range extendIniters {
		f(name, initer())
	}
}

func Get(name string) Initer {
	initer, _ := extendIniters[name]
	return initer
}

func Unregister(name string) {
	if _, ok := extendIniters[name]; ok {
		delete(extendIniters, name)
	}
}
