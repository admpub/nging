package route

var Default = &Collection{}

type Collection struct {
	Backend  IRegister
	Frontend IRegister
}
