package route

var Default = &Collection{}

type Collection struct {
	Backend  IRegister
	Frontend IRegister
}

func (r *Collection) Clear() {
	if r.Backend != nil {
		r.Backend.Clear()
	}
	if r.Backend != nil {
		r.Backend.Clear()
	}
}
