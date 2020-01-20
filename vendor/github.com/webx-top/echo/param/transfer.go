package param

type Transfer interface {
	Transform(interface{}, Store) interface{}
	Destination() string
}

type Transfers map[string]Transfer

func (t *Transfers) Add(name string, transfer Transfer) *Transfers {
	(*t)[name] = transfer
	return t
}

func (t *Transfers) Delete(names ...string) *Transfers {
	for _, name := range names {
		if _, ok := (*t)[name]; ok {
			delete(*t, name)
		}
	}
	return t
}

func NewTransform() *Transform {
	return &Transform{}
}

type Transform struct {
	Key  string
	Func func(interface{}, Store) interface{}
}

func (t *Transform) Transform(v interface{}, r Store) interface{} {
	if t.Func == nil {
		return v
	}
	return t.Func(v, r)
}

func (t *Transform) Destination() string {
	return t.Key
}

func (t *Transform) SetKey(key string) *Transform {
	t.Key = key
	return t
}

func (t *Transform) SetFunc(fn func(interface{}, Store) interface{}) *Transform {
	t.Func = fn
	return t
}
