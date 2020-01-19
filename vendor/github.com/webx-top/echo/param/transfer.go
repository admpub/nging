package param

type Transfer interface {
	Transform(interface{}) interface{}
	Destination() string
}

func NewTransform() *Transform {
	return &Transform{}
}

type Transform struct {
	Key  string
	Func func(interface{}) interface{}
}

func (t *Transform) Transform(v interface{}) interface{} {
	if t.Func == nil {
		return v
	}
	return t.Func(v)
}

func (t *Transform) Destination() string {
	return t.Key
}

func (t *Transform) SetKey(key string) *Transform {
	t.Key = key
	return t
}

func (t *Transform) SetFunc(fn func(interface{}) interface{}) *Transform {
	t.Func = fn
	return t
}
