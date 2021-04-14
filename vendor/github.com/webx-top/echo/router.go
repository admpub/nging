package echo

import (
	"bytes"
	"fmt"
	"net/url"
	"sort"
	"strings"
)

var (
	defaultRoute = &Route{}
	emptyMeta    = H{}
)

type (
	Router struct {
		tree   *node
		static map[string]*methodHandler
		routes []*Route
		nroute map[string][]int
		echo   *Echo
	}

	Route struct {
		Host       string
		Method     string
		Path       string
		Handler    Handler `json:"-" xml:"-"`
		Name       string
		Format     string
		Params     []string //param names
		Prefix     string
		Meta       H
		handler    interface{}   //原始handler
		middleware []interface{} //中间件
	}

	Routes []*Route

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

func (r Routes) SetName(name string) IRouter {
	for _, route := range r {
		route.SetName(name)
	}
	return r
}

func (r *Route) SetName(name string) IRouter {
	r.Name = name
	return r
}

func (r *Route) IsZero() bool {
	return r.Handler == nil
}

func (r *Route) Bool(name string, defaults ...interface{}) bool {
	if r.Meta == nil {
		return false
	}
	return r.Meta.Bool(name, defaults...)
}

func (r *Route) String(name string, defaults ...interface{}) string {
	if r.Meta == nil {
		return ``
	}
	return r.Meta.String(name, defaults...)
}

func (r *Route) Float64(name string, defaults ...interface{}) float64 {
	if r.Meta == nil {
		return 0
	}
	return r.Meta.Float64(name, defaults...)
}

func (r *Route) Float32(name string, defaults ...interface{}) float32 {
	if r.Meta == nil {
		return 0
	}
	return r.Meta.Float32(name, defaults...)
}

func (r *Route) Uint64(name string, defaults ...interface{}) uint64 {
	if r.Meta == nil {
		return 0
	}
	return r.Meta.Uint64(name, defaults...)
}

func (r *Route) Uint32(name string, defaults ...interface{}) uint32 {
	if r.Meta == nil {
		return 0
	}
	return r.Meta.Uint32(name, defaults...)
}

func (r *Route) Uint(name string, defaults ...interface{}) uint {
	if r.Meta == nil {
		return 0
	}
	return r.Meta.Uint(name, defaults...)
}

func (r *Route) Int64(name string, defaults ...interface{}) int64 {
	if r.Meta == nil {
		return 0
	}
	return r.Meta.Int64(name, defaults...)
}

func (r *Route) Int32(name string, defaults ...interface{}) int32 {
	if r.Meta == nil {
		return 0
	}
	return r.Meta.Int32(name, defaults...)
}

func (r *Route) Int(name string, defaults ...interface{}) int {
	if r.Meta == nil {
		return 0
	}
	return r.Meta.Int(name, defaults...)
}

func (r *Route) Get(name string, defaults ...interface{}) interface{} {
	if r.Meta == nil {
		return nil
	}
	return r.Meta.Get(name, defaults...)
}

func (r *Route) GetStore(names ...string) H {
	if r.Meta == nil {
		return emptyMeta
	}
	res := r.Meta
	for _, name := range names {
		res = res.GetStore(name)
	}
	return res
}

func (r *Route) MakeURI(params ...interface{}) (uri string) {
	length := len(params)
	if length == 1 {
		switch val := params[0].(type) {
		case url.Values:
			uri = r.Path
			for _, name := range r.Params {
				tag := `:` + name
				v := val.Get(name)
				uri = strings.Replace(uri, tag+`/`, v+`/`, -1)
				if strings.HasSuffix(uri, tag) {
					uri = strings.TrimSuffix(uri, tag) + v
				}
				val.Del(name)
			}
			q := val.Encode()
			if len(q) > 0 {
				uri += `?` + q
			}
		case map[string]string:
			uri = r.Path
			for _, name := range r.Params {
				tag := `:` + name
				v, y := val[name]
				if y {
					delete(val, name)
				}
				uri = strings.Replace(uri, tag+`/`, v+`/`, -1)
				if strings.HasSuffix(uri, tag) {
					uri = strings.TrimSuffix(uri, tag) + v
				}
			}
			sep := `?`
			keys := make([]string, 0, len(val))
			for k := range val {
				keys = append(keys, k)
			}
			sort.Strings(keys)
			for _, k := range keys {
				uri += sep + url.QueryEscape(k) + `=` + url.QueryEscape(val[k])
				sep = `&`
			}
		case []interface{}:
			uri = fmt.Sprintf(r.Format, val...)
		default:
			uri = fmt.Sprintf(r.Format, val)
		}
	} else {
		uri = fmt.Sprintf(r.Format, params...)
	}
	return
}

func (r *Route) apply(e *Echo) *Route {
	handler := e.ValidHandler(r.handler)
	middleware := r.middleware
	if hn, ok := handler.(Name); ok {
		r.Name = hn.Name()
	}
	if len(r.Name) == 0 {
		r.Name = HandlerName(handler)
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

func (r *Router) Handle(c Context) Handler {
	r.Find(c.Request().Method(), c.Request().URL().Path(), c)
	return c
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
		//Dump(rt)
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
			} else {
				r.insert(rt.Method, path[:i], nil, pkind, "", nil, -1)
			}
		} else if path[i] == '*' {
			uri.WriteString(`%v`)
			r.insert(rt.Method, path[:i], nil, skind, "", nil, -1)
			pnames = append(pnames, "*")
			r.insert(rt.Method, path[:i+1], rt.Handler, akind, ppath, pnames, rid)
			continue
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
	r.insert(rt.Method, path, rt.Handler, skind, ppath, pnames, rid)
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
		panic("echo: invalid method")
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
			if len(search) > 0 {
				cn.label = search[0]
			}
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
				if len(cn.pnames) == 0 {
					cn.pnames = pnames
				}
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

func (n *node) String() string {
	return Dump(n.Tree(), false)
}

func (n *node) Tree() H {
	children := make([]H, len(n.children))
	for k, v := range n.children {
		children[k] = v.Tree()
	}
	return H{
		"kind":          n.kind,
		"label":         string([]byte{n.label}),
		"prefix":        n.prefix,
		"parent":        n.parent,
		"children":      children,
		"ppath":         n.ppath,
		"pnames":        n.pnames,
		"methodHandler": n.methodHandler,
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
			break
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
			break
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
				nn = cn.parent // Next
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
		break
	}

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
