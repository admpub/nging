package param

type AsMap interface {
	AsMap() Store
}

type AsPartialMap interface {
	AsMap(...string) Store
}
