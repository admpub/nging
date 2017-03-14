package collector

func New(engine string) Collector {
	switch engine {
	case `goquery`:
		return &GoQuery{}
	default:
		return &Regexp{}
	}
}

type Collector interface {
}

type Browser interface {
	Open(string) string
}
