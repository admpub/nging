package echo

type (
	Group struct {
		prefix     string
		middleware []interface{}
		echo       *Echo
	}
)

func (g *Group) URL(h interface{}, params ...interface{}) string {
	return g.echo.URL(h, params...)
}

func (g *Group) SetRenderer(r Renderer) {
	g.echo.renderer = r
}

func (g *Group) Any(path string, h interface{}, middleware ...interface{}) {
	for _, m := range methods {
		g.add(m, path, h, middleware...)
	}
}

func (g *Group) Route(methods string, path string, h interface{}, middleware ...interface{}) {
	g.Match(httpMethodRegexp.Split(methods, -1), path, h, middleware...)
}

func (g *Group) Match(methods []string, path string, h interface{}, middleware ...interface{}) {
	for _, m := range methods {
		g.add(m, path, h, middleware...)
	}
}

func (g *Group) Use(middleware ...interface{}) {
	for _, m := range middleware {
		g.echo.ValidMiddleware(m)
		g.middleware = append(g.middleware, m)
	}
}

func (g *Group) Pre(middleware ...interface{}) {
	g.PreUse(middleware...)
}

func (g *Group) PreUse(middleware ...interface{}) {
	middlewares := make([]interface{}, 0)
	for _, m := range middleware {
		g.echo.ValidMiddleware(m)
		middlewares = append(middlewares, m)
	}
	g.middleware = append(middlewares, g.middleware...)
}

func (g *Group) Connect(path string, h interface{}, m ...interface{}) {
	g.add(CONNECT, path, h, m...)
}

func (g *Group) Delete(path string, h interface{}, m ...interface{}) {
	g.add(DELETE, path, h, m...)
}

func (g *Group) Get(path string, h interface{}, m ...interface{}) {
	g.add(GET, path, h, m...)
}

func (g *Group) Head(path string, h interface{}, m ...interface{}) {
	g.add(HEAD, path, h, m...)
}

func (g *Group) Options(path string, h interface{}, m ...interface{}) {
	g.add(OPTIONS, path, h, m...)
}

func (g *Group) Patch(path string, h interface{}, m ...interface{}) {
	g.add(PATCH, path, h, m...)
}

func (g *Group) Post(path string, h interface{}, m ...interface{}) {
	g.add(POST, path, h, m...)
}

func (g *Group) Put(path string, h interface{}, m ...interface{}) {
	g.add(PUT, path, h, m...)
}

func (g *Group) Trace(path string, h interface{}, m ...interface{}) {
	g.add(TRACE, path, h, m...)
}

func (g *Group) Group(prefix string, m ...interface{}) *Group {
	return g.echo.Group(g.prefix+prefix, m...)
}

func (g *Group) Prefix() string {
	return g.prefix
}

func (g *Group) add(method, path string, h interface{}, middleware ...interface{}) {
	// Combine into a new slice to avoid accidentally passing the same slice for
	// multiple routes, which would lead to later add() calls overwriting the
	// middleware from earlier calls.
	m := []interface{}{}
	m = append(m, g.middleware...)
	m = append(m, middleware...)

	g.echo.add(method, g.prefix+path, h, m...)
}
