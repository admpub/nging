package echo

type RouteRegister interface {
	Any(path string, h interface{}, middleware ...interface{})
	Route(methods string, path string, h interface{}, middleware ...interface{})
	Match(methods []string, path string, h interface{}, middleware ...interface{})
	Connect(path string, h interface{}, m ...interface{})
	Delete(path string, h interface{}, m ...interface{})
	Get(path string, h interface{}, m ...interface{})
	Head(path string, h interface{}, m ...interface{})
	Options(path string, h interface{}, m ...interface{})
	Patch(path string, h interface{}, m ...interface{})
	Post(path string, h interface{}, m ...interface{})
	Put(path string, h interface{}, m ...interface{})
	Trace(path string, h interface{}, m ...interface{})
}

type MiddlewareRegister interface {
	Use(middleware ...interface{})
	Pre(middleware ...interface{})
}

type URLBuilder interface {
	URL(interface{}, ...interface{}) string
}

type ICore interface {
	RouteRegister
	MiddlewareRegister
	URLBuilder
}
