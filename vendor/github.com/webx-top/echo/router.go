package echo

import (
	"bytes"
)

var defaultRoute = &Route{}

type (
	Router struct {
		tree   *node
		static map[string]*methodHandler
		routes []*Route
		nroute map[string][]int
		echo   *Echo
	}

	meta struct {
		name string
		meta H
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
		handler     interface{}   //原始handler
		middleware  []interface{} //中间件
	}

	endpoint struct {
		handler Handler
		rid     int //routes index
	}

	node struct {
		kind          kind
		label         byte
		prefix        string
		parent        *node
		children      children
		ppath         string
		pnames        []string
		methodHandler *methodHandler
	}
	kind          uint8
	children      []*node
	methodHandler struct {
		connect *endpoint
		delete  *endpoint
		get     *endpoint
		head    *endpoint
		options *endpoint
		patch   *endpoint
		post    *endpoint
		put     *endpoint
		trace   *endpoint
	}
)

const (
	skind kind = iota
	pkind
	akind
)

func (r *Route) IsZero() bool {
	return r.Handler == nil
}

func (r *Route) apply(e *Echo) *Route {
	handler := e.ValidHandler(r.handler)
	middleware := r.middleware
	if hn, ok := handler.(Name); ok {
		r.HandlerName = hn.Name()
	} else {
		r.HandlerName = HandlerName(handler)
	}
	if mt, ok := handler.(Meta); ok {
		r.Meta = mt.Meta()
	} else {
		r.Meta = H{}
	}
	for i := len(middleware) - 1; i >= 0; i-- {
		m := middleware[i]
		mw := e.ValidMiddleware(m)
		handler = mw.Handle(handler)
	}
	r.Handler = handler
	return r
}

func (m *methodHandler) addHandler(method string, h Handler, rid int) {
	endpoint := &endpoint{handler: h, rid: rid}
	switch method {
	case GET:
		m.get = endpoint
	case POST:
		m.post = endpoint
	case PUT:
		m.put = endpoint
	case DELETE:
		m.delete = endpoint
	case PATCH:
		m.patch = endpoint
	case OPTIONS:
		m.options = endpoint
	case HEAD:
		m.head = endpoint
	case CONNECT:
		m.connect = endpoint
	case TRACE:
		m.trace = endpoint
	}
}

func (m *methodHandler) findHandler(method string) Handler {
	endpoint := m.find(method)
	if endpoint == nil {
		return nil
	}
	return endpoint.handler
}

func (m *methodHandler) find(method string) *endpoint {
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

func (m *methodHandler) check405() HandlerFunc {
	for _, method := range methods {
		if r := m.findHandler(method); r != nil {
			return MethodNotAllowedHandler
		}
	}
	return NotFoundHandler
}

func (m *methodHandler) applyHandler(method string, ctx *xContext) {
	endpoint := m.find(method)
	if endpoint != nil {
		ctx.handler = endpoint.handler
		ctx.rid = endpoint.rid
	} else {
		ctx.handler = nil
	}
}

func NewRouter(e *Echo) *Router {
	return &Router{
		tree: &node{
			methodHandler: new(methodHandler),
		},
		static: map[string]*methodHandler{},
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

// Add 添加路由
// method: 方法(GET/POST/PUT/DELETE/PATCH/OPTIONS/HEAD/CONNECT/TRACE)
// prefix: group前缀
// path: 路径(含前缀)
// h: Handler
// name: Handler名
// meta: meta数据
func (r *Router) Add(rt *Route, rid int) {
	path := rt.Path
	ppath := path        // Pristine path
	pnames := []string{} // Param names
	uri := new(bytes.Buffer)
	defer func() {
		rt.Format = uri.String()
		rt.Params = pnames
	}()
	for i, l := 0, len(path); i < l; i++ {
		if path[i] == ':' {
			uri.WriteString(`%v`)
			j := i + 1

			r.insert(rt.Method, path[:i], nil, skind, "", nil, -1)
			for ; i < l && path[i] != '/'; i++ {
			}

			pnames = append(pnames, path[j:i])
			path = path[:j] + path[i:]
			i, l = j, len(path)

			if i == l {
				r.insert(rt.Method, path[:i], rt.Handler, pkind, ppath, pnames, rid)
				return
			}
			r.insert(rt.Method, path[:i], nil, pkind, ppath, pnames, -1)
		} else if path[i] == '*' {
			uri.WriteString(`%v`)
			r.insert(rt.Method, path[:i], nil, skind, "", nil, -1)
			pnames = append(pnames, "*")
			r.insert(rt.Method, path[:i+1], rt.Handler, akind, ppath, pnames, rid)
			return
		}

		if i < l {
			uri.WriteByte(path[i])
		}
	}

	//static route
	if m, ok := r.static[path]; ok {
		m.addHandler(rt.Method, rt.Handler, rid)
	} else {
		m = &methodHandler{}
		m.addHandler(rt.Method, rt.Handler, rid)
		r.static[path] = m
	}
	//r.insert(method, path, h, skind, ppath, pnames, e)
	return
}

func (r *Router) insert(method, path string, h Handler, t kind, ppath string, pnames []string, rid int) {
	e := r.echo
	// Adjust max param
	l := len(pnames)
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
			if h != nil {
				cn.kind = t
				cn.addHandler(method, h, rid)
				cn.ppath = ppath
				cn.pnames = pnames
			}
		} else if l < pl {
			// Split node
			n := newNode(cn.kind, cn.prefix[l:], cn, cn.children, cn.methodHandler, cn.ppath, cn.pnames)

			// Reset parent node
			cn.kind = skind
			cn.label = cn.prefix[0]
			cn.prefix = cn.prefix[:l]
			cn.children = nil
			cn.methodHandler = new(methodHandler)
			cn.ppath = ""
			cn.pnames = nil

			cn.addChild(n)

			if l == sl {
				// At parent node
				cn.kind = t
				cn.addHandler(method, h, rid)
				cn.ppath = ppath
				cn.pnames = pnames
			} else {
				// Create child node
				n = newNode(t, search[l:], cn, nil, new(methodHandler), ppath, pnames)
				n.addHandler(method, h, rid)
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
			n := newNode(t, search, cn, nil, new(methodHandler), ppath, pnames)
			n.addHandler(method, h, rid)
			cn.addChild(n)
		} else {
			// Node already exists
			if h != nil {
				cn.addHandler(method, h, rid)
				cn.ppath = ppath
				cn.pnames = pnames
			}
		}
		return
	}
}

func newNode(t kind, pre string, p *node, c children, mh *methodHandler, ppath string, pnames []string) *node {
	return &node{
		kind:          t,
		label:         pre[0],
		prefix:        pre,
		parent:        p,
		children:      c,
		ppath:         ppath,
		pnames:        pnames,
		methodHandler: mh,
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

func (n *node) addHandler(method string, h Handler, rid int) {
	n.methodHandler.addHandler(method, h, rid)
}

func (n *node) findHandler(method string) Handler {
	return n.methodHandler.findHandler(method)
}

func (n *node) find(method string) *endpoint {
	return n.methodHandler.find(method)
}

func (n *node) check405() HandlerFunc {
	return n.methodHandler.check405()
}

func (n *node) applyHandler(method string, ctx *xContext) {
	n.methodHandler.applyHandler(method, ctx)
	ctx.path = n.ppath
	ctx.pnames = n.pnames
}

func (r *Router) Find(method, path string, context Context) {
	ctx := context.Object()
	ctx.path = path
	cn := r.tree // Current node as root

	if m, ok := r.static[path]; ok {
		m.applyHandler(method, ctx)
		if ctx.handler == nil {
			ctx.handler = m.check405()
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
	cn.applyHandler(method, ctx)

	// NOTE: Slow zone...
	if ctx.handler == nil {
		// Dig further for any, might have an empty value for *, e.g.
		// serving a directory. Issue #207.
		if child := cn.findChildByKind(akind); child != nil {
			child.applyHandler(method, ctx)
			if ctx.handler == nil {
				ctx.handler = child.check405()
			}
			pvalues[len(child.pnames)-1] = ""
			return
		}
		ctx.handler = cn.check405()
	}
	return
}
