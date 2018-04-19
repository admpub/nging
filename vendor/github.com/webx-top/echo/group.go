package echo

type Group struct {
	prefix     string
	middleware []interface{}
	echo       *Echo
}

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
	g.Match(splitHTTPMethod.Split(methods, -1), path, h, middleware...)
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
		if g.echo.MiddlewareDebug {
			g.echo.logger.Debugf(`Middleware[Use](%p): [%s] -> %s`, m, g.prefix, HandlerName(m))
		}
	}
}

// Pre is an alias for `PreUse` function.
func (g *Group) Pre(middleware ...interface{}) {
	g.PreUse(middleware...)
}

func (g *Group) PreUse(middleware ...interface{}) {
	var middlewares []interface{}
	for _, m := range middleware {
		g.echo.ValidMiddleware(m)
		middlewares = append(middlewares, m)
		if g.echo.MiddlewareDebug {
			g.echo.logger.Debugf(`Middleware[Pre](%p): [%s] -> %s`, m, g.prefix, HandlerName(m))
		}
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

func (g *Group) Group(prefix string, middleware ...interface{}) *Group {
	m := []interface{}{}
	m = append(m, g.middleware...)
	m = append(m, middleware...)
	return g.echo.Group(g.prefix+prefix, m...)
}

// Static implements `Echo#Static()` for sub-routes within the Group.
func (g *Group) Static(prefix, root string) {
	static(g, prefix, root)
}

// File implements `Echo#File()` for sub-routes within the Group.
func (g *Group) File(path, file string) {
	g.echo.File(g.prefix+path, file)
}

func (g *Group) Prefix() string {
	return g.prefix
}

func (g *Group) Echo() *Echo {
	return g.echo
}

// MetaHandler Add meta information about endpoint
func (g *Group) MetaHandler(m H, handler interface{}) Handler {
	return &MetaHandler{m, g.echo.ValidHandler(handler)}
}

func (g *Group) add(method, path string, h interface{}, middleware ...interface{}) {
	// Combine into a new slice to avoid accidentally passing the same slice for
	// multiple routes, which would lead to later add() calls overwriting the
	// middleware from earlier calls.
	m := []interface{}{}
	m = append(m, g.middleware...)
	m = append(m, middleware...)

	g.echo.add(method, g.prefix, g.prefix+path, h, m...)
}
