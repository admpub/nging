package echo

import (
	"bytes"
)

type (
	Router struct {
		tree   *node
		static map[string]*methodHandler
		routes []*Route
		nroute map[string][]int
		echo   *Echo
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
		connect Handler
		delete  Handler
		get     Handler
		head    Handler
		options Handler
		patch   Handler
		post    Handler
		put     Handler
		trace   Handler
	}
)

const (
	skind kind = iota
	pkind
	akind
)

func (m *methodHandler) addHandler(method string, h Handler) {
	switch method {
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

func (m *methodHandler) findHandler(method string) Handler {
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

func (r *Router) Add(method, path string, h Handler, e *Echo) (fpath string, pnames []string) {
	ppath := path       // Pristine path
	pnames = []string{} // Param names
	uri := new(bytes.Buffer)
	defer func() {
		fpath = uri.String()
	}()

	for i, l := 0, len(path); i < l; i++ {
		if path[i] == ':' {
			uri.WriteString(`%v`)
			j := i + 1

			r.insert(method, path[:i], nil, skind, "", nil, e)
			for ; i < l && path[i] != '/'; i++ {
			}

			pnames = append(pnames, path[j:i])
			path = path[:j] + path[i:]
			i, l = j, len(path)

			if i == l {
				r.insert(method, path[:i], h, pkind, ppath, pnames, e)
				return
			}
			r.insert(method, path[:i], nil, pkind, ppath, pnames, e)
		} else if path[i] == '*' {
			uri.WriteString(`%v`)
			r.insert(method, path[:i], nil, skind, "", nil, e)
			pnames = append(pnames, "_*")
			r.insert(method, path[:i+1], h, akind, ppath, pnames, e)
			return
		}

		if i < l {
			uri.WriteByte(path[i])
		}
	}

	//static route
	if m, ok := r.static[path]; ok {
		m.addHandler(method, h)
	} else {
		m = &methodHandler{}
		m.addHandler(method, h)
		r.static[path] = m
	}
	//r.insert(method, path, h, skind, ppath, pnames, e)
	return
}

func (r *Router) insert(method, path string, h Handler, t kind, ppath string, pnames []string, e *Echo) {
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
				cn.addHandler(method, h)
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
				cn.addHandler(method, h)
				cn.ppath = ppath
				cn.pnames = pnames
			} else {
				// Create child node
				n = newNode(t, search[l:], cn, nil, new(methodHandler), ppath, pnames)
				n.addHandler(method, h)
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
			n.addHandler(method, h)
			cn.addChild(n)
		} else {
			// Node already exists
			if h != nil {
				cn.addHandler(method, h)
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

func (n *node) addHandler(method string, h Handler) {
	n.methodHandler.addHandler(method, h)
}

func (n *node) findHandler(method string) Handler {
	return n.methodHandler.findHandler(method)
}

func (n *node) check405() HandlerFunc {
	for _, m := range methods {
		if h := n.findHandler(m); h != nil {
			return MethodNotAllowedHandler
		}
	}
	return NotFoundHandler
}

func (r *Router) Find(method, path string, context Context) {
	ctx := context.Object()
	cn := r.tree // Current node as root

	if m, ok := r.static[path]; ok {
		ctx.handler = m.findHandler(method)
		ctx.path = path
		ctx.pnames = []string{}
		if ctx.handler == nil {
			ctx.handler = cn.check405()
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
	ctx.path = cn.ppath
	ctx.pnames = cn.pnames
	ctx.handler = cn.findHandler(method)

	// NOTE: Slow zone...
	if ctx.handler == nil {
		ctx.handler = cn.check405()

		// Dig further for any, might have an empty value for *, e.g.
		// serving a directory. Issue #207.
		if cn = cn.findChildByKind(akind); cn == nil {
			return
		}
		if ctx.handler = cn.findHandler(method); ctx.handler == nil {
			ctx.handler = cn.check405()
		}
		ctx.path = cn.ppath
		ctx.pnames = cn.pnames
		pvalues[len(cn.pnames)-1] = ""
	}
	return
}
