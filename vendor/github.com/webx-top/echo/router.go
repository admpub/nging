package echo

import (
	"bytes"
)

var (
	NotFoundRoute = &Route{
		Handler: NotFoundHandler,
	}
	MethodNotAllowedRoute = &Route{
		Handler: MethodNotAllowedHandler,
	}
)

type (
	Router struct {
		tree   *node
		static map[string]*methodRoute
		routes []*Route
		nroute map[string][]int
		echo   *Echo
	}

	Route struct {
		Method      string
		Path        string
		Handler     Handler `json:"-" xml:"-"`
		HandlerName string
		Format      string
		Params      []string //param names
		Prefix      string
		Meta        H
	}

	node struct {
		kind        kind
		label       byte
		prefix      string
		parent      *node
		children    children
		ppath       string
		pnames      []string
		methodRoute *methodRoute
	}
	kind        uint8
	children    []*node
	methodRoute struct {
		connect *Route
		delete  *Route
		get     *Route
		head    *Route
		options *Route
		patch   *Route
		post    *Route
		put     *Route
		trace   *Route
	}
)

const (
	skind kind = iota
	pkind
	akind
)

func (m *methodRoute) addRoute(h *Route) {
	switch h.Method {
	case GET:
		m.get = h
	case POST:
		m.post = h
	case PUT:
		m.put = h
	case DELETE:
		m.delete = h
	case PATCH:
		m.patch = h
	case OPTIONS:
		m.options = h
	case HEAD:
		m.head = h
	case CONNECT:
		m.connect = h
	case TRACE:
		m.trace = h
	}
}

func (m *methodRoute) findRoute(method string) *Route {
	switch method {
	case GET:
		return m.get
	case POST:
		return m.post
	case PUT:
		return m.put
	case DELETE:
		return m.delete
	case PATCH:
		return m.patch
	case OPTIONS:
		return m.options
	case HEAD:
		return m.head
	case CONNECT:
		return m.connect
	case TRACE:
		return m.trace
	default:
		return nil
	}
}

func (m *methodRoute) check405() *Route {
	for _, method := range methods {
		if r := m.findRoute(method); r != nil {
			return MethodNotAllowedRoute
		}
	}
	return NotFoundRoute
}

func NewRouter(e *Echo) *Router {
	return &Router{
		tree: &node{
			methodRoute: new(methodRoute),
		},
		static: map[string]*methodRoute{},
		routes: []*Route{},
		nroute: map[string][]int{},
		echo:   e,
	}
}

func (r *Router) Handle(h Handler) Handler {
	return HandlerFunc(func(c Context) error {
		method := c.Request().Method()
		path := c.Request().URL().Path()
		r.Find(method, path, c)
		return c.Handle(c)
	})
}

func (r *Router) Add(method, path string, h Handler, name string, meta H, e *Echo) (route *Route) {
	ppath := path        // Pristine path
	pnames := []string{} // Param names
	uri := new(bytes.Buffer)
	newRoute := func(h Handler, name string, ppath string, pnames []string) *Route {
		return &Route{
			Method:      method,
			Path:        ppath,
			Handler:     h,
			HandlerName: name,
			Format:      uri.String(),
			Params:      pnames,
			Meta:        meta,
		}
	}
	for i, l := 0, len(path); i < l; i++ {
		if path[i] == ':' {
			uri.WriteString(`%v`)
			j := i + 1

			r.insert(path[:i], skind, newRoute(nil, ``, ``, nil), e)
			for ; i < l && path[i] != '/'; i++ {
			}

			pnames = append(pnames, path[j:i])
			path = path[:j] + path[i:]
			i, l = j, len(path)

			if i == l {
				route = newRoute(h, name, ppath, pnames)
				r.insert(path[:i], pkind, route, e)
				return
			}
			r.insert(path[:i], pkind, newRoute(nil, ``, ppath, pnames), e)
		} else if path[i] == '*' {
			uri.WriteString(`%v`)
			r.insert(path[:i], skind, newRoute(h, ``, ``, nil), e)
			pnames = append(pnames, "_*")
			route = newRoute(h, name, ppath, pnames)
			r.insert(path[:i+1], akind, route, e)
			return
		}

		if i < l {
			uri.WriteByte(path[i])
		}
	}

	route = newRoute(h, name, ppath, pnames)
	//static route
	if m, ok := r.static[path]; ok {
		m.addRoute(route)
	} else {
		m = &methodRoute{}
		m.addRoute(route)
		r.static[path] = m
	}
	//r.insert(path, skind, route, e)
	return
}

func (r *Router) insert(path string, t kind, route *Route, e *Echo) {
	// Adjust max param
	l := len(route.Params)
	if *e.maxParam < l {
		*e.maxParam = l
	}

	cn := r.tree // Current node as root
	if cn == nil {
		panic("echo => invalid method")
	}
	search := path

	for {
		sl := len(search)
		pl := len(cn.prefix)
		l := 0

		// LCP
		max := pl
		if sl < max {
			max = sl
		}
		for ; l < max && search[l] == cn.prefix[l]; l++ {
		}

		if l == 0 {
			// At root node
			cn.label = search[0]
			cn.prefix = search
			if route.Handler != nil {
				cn.kind = t
				cn.addRoute(route)
				cn.ppath = route.Path
				cn.pnames = route.Params
			}
		} else if l < pl {
			// Split node
			n := newNode(cn.kind, cn.prefix[l:], cn, cn.children, cn.methodRoute, cn.ppath, cn.pnames)

			// Reset parent node
			cn.kind = skind
			cn.label = cn.prefix[0]
			cn.prefix = cn.prefix[:l]
			cn.children = nil
			cn.methodRoute = new(methodRoute)
			cn.ppath = ""
			cn.pnames = nil

			cn.addChild(n)

			if l == sl {
				// At parent node
				cn.kind = t
				cn.addRoute(route)
				cn.ppath = route.Path
				cn.pnames = route.Params
			} else {
				// Create child node
				n = newNode(t, search[l:], cn, nil, new(methodRoute), route.Path, route.Params)
				n.addRoute(route)
				cn.addChild(n)
			}
		} else if l < sl {
			search = search[l:]
			c := cn.findChildWithLabel(search[0])
			if c != nil {
				// Go deeper
				cn = c
				continue
			}
			// Create child node
			n := newNode(t, search, cn, nil, new(methodRoute), route.Path, route.Params)
			n.addRoute(route)
			cn.addChild(n)
		} else {
			// Node already exists
			if route.Handler != nil {
				cn.addRoute(route)
				cn.ppath = route.Path
				cn.pnames = route.Params
			}
		}
		return
	}
}

func newNode(t kind, pre string, p *node, c children, mh *methodRoute, ppath string, pnames []string) *node {
	return &node{
		kind:        t,
		label:       pre[0],
		prefix:      pre,
		parent:      p,
		children:    c,
		ppath:       ppath,
		pnames:      pnames,
		methodRoute: mh,
	}
}

func (n *node) addChild(c *node) {
	n.children = append(n.children, c)
}

func (n *node) findChild(l byte, t kind) *node {
	for _, c := range n.children {
		if c.label == l && c.kind == t {
			return c
		}
	}
	return nil
}

func (n *node) findChildWithLabel(l byte) *node {
	for _, c := range n.children {
		if c.label == l {
			return c
		}
	}
	return nil
}

func (n *node) findChildByKind(t kind) *node {
	for _, c := range n.children {
		if c.kind == t {
			return c
		}
	}
	return nil
}

func (n *node) addRoute(route *Route) {
	n.methodRoute.addRoute(route)
}

func (n *node) findRoute(method string) *Route {
	return n.methodRoute.findRoute(method)
}

func (n *node) check405() *Route {
	return n.methodRoute.check405()
}

func (r *Router) Find(method, path string, context Context) {
	ctx := context.Object()
	cn := r.tree // Current node as root

	if m, ok := r.static[path]; ok {
		ctx.route = m.findRoute(method)
		if ctx.route == nil {
			ctx.route = m.check405()
		}
		return
	}

	var (
		search  = path
		c       *node  // Child node
		n       int    // Param counter
		nk      kind   // Next kind
		nn      *node  // Next node
		ns      string // Next search
		pvalues = context.ParamValues()
	)

	// Search order static > param > any
	for {
		if search == "" {
			goto End
		}

		pl := 0 // Prefix length
		l := 0  // LCP length

		if cn.label != ':' {
			sl := len(search)
			pl = len(cn.prefix)

			// LCP
			max := pl
			if sl < max {
				max = sl
			}
			for ; l < max && search[l] == cn.prefix[l]; l++ {
			}
		}

		if l == pl {
			// Continue search
			search = search[l:]
		} else {
			cn = nn
			search = ns
			if nk == pkind {
				goto Param
			} else if nk == akind {
				goto Any
			}
			// Not found
			return
		}

		if search == "" {
			goto End
		}

		// Static node
		if c = cn.findChild(search[0], skind); c != nil {
			// Save next
			if cn.prefix[len(cn.prefix)-1] == '/' {
				nk = pkind
				nn = cn
				ns = search
			}
			cn = c
			continue
		}

		// Param node
	Param:
		if c = cn.findChildByKind(pkind); c != nil {

			if len(pvalues) == n {
				continue
			}

			// Save next
			if cn.prefix[len(cn.prefix)-1] == '/' {
				nk = akind
				nn = cn
				ns = search
			}

			cn = c
			i, l := 0, len(search)
			for ; i < l && search[i] != '/'; i++ {
			}
			pvalues[n] = search[:i]
			n++
			search = search[i:]
			continue
		}

		// Any node
	Any:
		if cn = cn.findChildByKind(akind); cn == nil {
			if nn != nil {
				cn = nn
				nn = nil // Next
				search = ns
				if nk == pkind {
					goto Param
				} else if nk == akind {
					goto Any
				}
			}
			// Not found
			return
		}
		pvalues[len(cn.pnames)-1] = search
		goto End
	}

End:
	ctx.route = cn.findRoute(method)

	// NOTE: Slow zone...
	if ctx.route == nil {
		check405 := cn.check405
		// Dig further for any, might have an empty value for *, e.g.
		// serving a directory. Issue #207.
		if cn = cn.findChildByKind(akind); cn != nil {
			ctx.route = cn.findRoute(method)
			pvalues[len(cn.pnames)-1] = ""
		}
		if ctx.route == nil {
			ctx.route = check405()
		}
	}
	return
}
