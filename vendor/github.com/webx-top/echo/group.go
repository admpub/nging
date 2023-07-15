package echo

type Group struct {
	host       *host
	prefix     string
	middleware []interface{}
	echo       *Echo
	meta       H
}

func (g *Group) URL(h interface{}, params ...interface{}) string {
	return g.echo.URL(h, params...)
}

func (g *Group) SetAlias(alias string) *Group {
	if g.host != nil {
		g.host.alias = alias
		for a, v := range g.echo.hostAlias {
			if v == g.host.name {
				delete(g.echo.hostAlias, a)
			}
		}
		if len(alias) > 0 {
			g.echo.hostAlias[alias] = g.host.name
		}
	}
	return g
}

func (g *Group) Alias(alias string) Hoster {
	if name, ok := g.echo.hostAlias[alias]; ok {
		hs, ok := g.echo.hosts[name]
		if !ok || hs == nil || hs.group == nil {
			return nil
		}
		return hs.group.host
	}
	return nil
}

func (g *Group) SetRenderer(r Renderer) {
	g.echo.renderer = r
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

// Pre adds handler to the middleware chain.
func (g *Group) Pre(middleware ...interface{}) {
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

func (g *Group) Connect(path string, h interface{}, m ...interface{}) IRouter {
	return g.Add(CONNECT, path, h, m...)
}

func (g *Group) Delete(path string, h interface{}, m ...interface{}) IRouter {
	return g.Add(DELETE, path, h, m...)
}

func (g *Group) Get(path string, h interface{}, m ...interface{}) IRouter {
	return g.Add(GET, path, h, m...)
}

func (g *Group) Head(path string, h interface{}, m ...interface{}) IRouter {
	return g.Add(HEAD, path, h, m...)
}

func (g *Group) Options(path string, h interface{}, m ...interface{}) IRouter {
	return g.Add(OPTIONS, path, h, m...)
}

func (g *Group) Patch(path string, h interface{}, m ...interface{}) IRouter {
	return g.Add(PATCH, path, h, m...)
}

func (g *Group) Post(path string, h interface{}, m ...interface{}) IRouter {
	return g.Add(POST, path, h, m...)
}

func (g *Group) Put(path string, h interface{}, m ...interface{}) IRouter {
	return g.Add(PUT, path, h, m...)
}

func (g *Group) Trace(path string, h interface{}, m ...interface{}) IRouter {
	return g.Add(TRACE, path, h, m...)
}

func (g *Group) Any(path string, h interface{}, middleware ...interface{}) IRouter {
	routes := Routes{}
	for _, m := range methods {
		routes = append(routes, g.Add(m, path, h, middleware...))
	}
	return routes
}

func (g *Group) Route(methods string, path string, h interface{}, middleware ...interface{}) IRouter {
	return g.Match(splitHTTPMethod.Split(methods, -1), path, h, middleware...)
}

func (g *Group) Match(methods []string, path string, h interface{}, middleware ...interface{}) IRouter {
	routes := Routes{}
	for _, m := range methods {
		routes = append(routes, g.Add(m, path, h, middleware...))
	}
	return routes
}

func (g *Group) Group(prefix string, middleware ...interface{}) *Group {
	m := []interface{}{}
	m = append(m, g.middleware...)
	m = append(m, middleware...)
	if g.host != nil {
		subG, y := g.echo.hosts[g.host.name].groups[prefix]
		if !y {
			subG = &Group{host: g.host, prefix: prefix, echo: g.echo, meta: H{}}
			g.echo.hosts[g.host.name].groups[prefix] = subG
			if len(g.meta) > 0 {
				subG.meta.DeepMerge(g.meta)
			}
		}
		if len(m) > 0 {
			subG.Use(m...)
		}
		return subG
	}
	return g.echo.subgroup(g, prefix, m...)
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
func (g *Group) MetaHandler(m H, handler interface{}, requests ...interface{}) Handler {
	return g.echo.MetaHandler(m, handler, requests...)
}

func (g *Group) Add(method, path string, h interface{}, middleware ...interface{}) *Route {
	// Combine into a new slice to avoid accidentally passing the same slice for
	// multiple routes, which would lead to later add() calls overwriting the
	// middleware from earlier calls.
	m := []interface{}{}
	m = append(m, g.middleware...)
	m = append(m, middleware...)
	var host string
	if g.host != nil {
		host = g.host.name
	}
	r := g.echo.add(host, method, g.prefix, g.prefix+path, h, m...)
	if len(g.meta) > 0 {
		r.Meta = H{}
		r.Meta.DeepMerge(g.meta)
	}
	return r
}

func (g *Group) SetMeta(meta H) *Group {
	g.meta = meta
	return g
}

func (g *Group) SetMetaKV(key string, value interface{}) *Group {
	if g.meta == nil {
		g.meta = H{}
	}
	g.meta[key] = value
	return g
}
