package echo

import (
	"bytes"
	"fmt"
	"net/url"
	"regexp"
	"sort"
	"strings"

	"github.com/webx-top/echo/param"
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

	Rewriter interface {
		Rewrite(string) string
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
		kind           kind
		label          byte
		prefix         string
		parent         *node
		staticChildren children
		ppath          string
		pnames         []string
		methodHandler  *methodHandler
		regExp         *regexp.Regexp
		regexChild     *node
		paramChild     *node
		anyChild       *node
		// isLeaf indicates that node does not have child routes
		isLeaf bool
		// isHandler indicates that node has at least one handler registered to it
		isHandler bool
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
	staticKind kind = iota // static kind
	regexKind              // regexp kind
	paramKind              // param kind
	anyKind                // any kind

	regexLabel = byte('<')
	paramLabel = byte(':')
	anyLabel   = byte('*')
)

func (r Routes) SetName(name string) IRouter {
	for _, route := range r {
		route.SetName(name)
	}
	return r
}

func (r Routes) SetMeta(meta param.Store) IRouter {
	for _, route := range r {
		route.Meta = meta
	}
	return r
}

func (r Routes) SetMetaKV(key string, value interface{}) IRouter {
	for _, route := range r {
		if route.Meta == nil {
			route.Meta = param.Store{}
		}
		route.Meta[key] = value
	}
	return r
}

func (r Routes) GetName() string {
	for _, route := range r {
		return route.GetName()
	}
	return ``
}

func (r Routes) GetMeta() param.Store {
	for _, route := range r {
		return route.Meta
	}
	return nil
}

func (r *Route) SetName(name string) IRouter {
	r.Name = name
	return r
}

func (r *Route) SetMeta(meta param.Store) IRouter {
	r.Meta = meta
	return r
}

func (r *Route) SetMetaKV(key string, value interface{}) IRouter {
	if r.Meta == nil {
		r.Meta = param.Store{}
	}
	r.Meta[key] = value
	return r
}

func (r *Route) GetName() string {
	return r.Name
}

func (r *Route) GetMeta() param.Store {
	return r.Meta
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

func (r *Route) MakeURI(e *Echo, params ...interface{}) (uri string) {
	length := len(params)
	if length != 1 {
		uri = fmt.Sprintf(r.Format, params...)
		uri = e.wrapURI(uri)
		return
	}
	switch val := params[0].(type) {
	case url.Values:
		uri = r.Path
		if len(r.Params) > 0 {
			values := make([]interface{}, len(r.Params))
			for index, name := range r.Params {
				values[index] = val.Get(name)
				val.Del(name)
			}
			uri = fmt.Sprintf(r.Format, values...)
		}
		uri = e.wrapURI(uri)
		q := val.Encode()
		if len(q) > 0 {
			uri += `?` + q
		}
	case param.Store:
		uri = r.Path
		if len(r.Params) > 0 {
			values := make([]interface{}, len(r.Params))
			for index, name := range r.Params {
				var ok bool
				values[index], ok = val[name]
				if ok {
					delete(val, name)
				}
			}
			uri = fmt.Sprintf(r.Format, values...)
		}
		uri = e.wrapURI(uri)
		sep := `?`
		keys := make([]string, 0, len(val))
		for k := range val {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			uri += sep + url.QueryEscape(k) + `=` + url.QueryEscape(val.String(k))
			sep = `&`
		}
	case map[string]interface{}:
		uri = r.Path
		if len(r.Params) > 0 {
			values := make([]interface{}, len(r.Params))
			for index, name := range r.Params {
				var ok bool
				values[index], ok = val[name]
				if ok {
					delete(val, name)
				}
			}
			uri = fmt.Sprintf(r.Format, values...)
		}
		uri = e.wrapURI(uri)
		sep := `?`
		keys := make([]string, 0, len(val))
		for k := range val {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			uri += sep + url.QueryEscape(k) + `=` + url.QueryEscape(param.AsString(val[k]))
			sep = `&`
		}
	case map[string]string:
		uri = r.Path
		if len(r.Params) > 0 {
			values := make([]interface{}, len(r.Params))
			for index, name := range r.Params {
				var ok bool
				values[index], ok = val[name]
				if ok {
					delete(val, name)
				}
			}
			uri = fmt.Sprintf(r.Format, values...)
		}
		uri = e.wrapURI(uri)
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
		uri = e.wrapURI(uri)
	default:
		uri = fmt.Sprintf(r.Format, val)
		uri = e.wrapURI(uri)
	}
	return
}

func (r *Route) apply(e *Echo) *Route {
	handler := e.ValidHandler(r.handler)
	middleware := r.middleware
	if len(r.Name) == 0 {
		if hn, ok := handler.(Name); ok {
			r.Name = hn.Name()
		}
		if len(r.Name) == 0 {
			r.Name = HandlerName(handler)
		}
	}
	if len(r.Meta) == 0 {
		if mt, ok := handler.(Meta); ok {
			r.Meta = mt.Meta()
		} else if r.Meta == nil {
			r.Meta = H{}
		}
	}
	for i := len(middleware) - 1; i >= 0; i-- {
		m := middleware[i]
		mw := e.ValidMiddleware(m)
		handler = mw.Handle(handler)
	}
	r.Handler = handler
	return r
}

func (e *endpoint) Map() H {
	return H{`handler`: HandlerName(e.handler), `index`: e.rid}
}

func (m *methodHandler) isHandler() bool {
	return m.connect != nil ||
		m.delete != nil ||
		m.get != nil ||
		m.head != nil ||
		m.options != nil ||
		m.patch != nil ||
		m.post != nil ||
		m.put != nil ||
		m.trace != nil
}

func (m *methodHandler) Map() H {
	r := H{}
	if m.get != nil {
		r[`get`] = m.get.Map()
	}
	if m.post != nil {
		r[`post`] = m.post.Map()
	}
	if m.put != nil {
		r[`put`] = m.put.Map()
	}
	if m.delete != nil {
		r[`delete`] = m.delete.Map()
	}
	if m.patch != nil {
		r[`patch`] = m.patch.Map()
	}
	if m.options != nil {
		r[`options`] = m.options.Map()
	}
	if m.head != nil {
		r[`head`] = m.head.Map()
	}
	if m.connect != nil {
		r[`connect`] = m.connect.Map()
	}
	if m.trace != nil {
		r[`trace`] = m.trace.Map()
	}
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

func (m *methodHandler) checkMethodNotAllowed() HandlerFunc {
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
	return r.Dispatch(c, c.Request().URL().Path())
}

func (r *Router) Dispatch(c Context, path string, _method ...string) Handler {
	method := c.Request().Method()
	if len(_method) > 0 && len(_method[0]) > 0 {
		method = _method[0]
	}
	found := r.Find(method, path, c)
	if !found {
		ext := c.DefaultExtension()
		if len(ext) == 0 {
			return c
		}
		if strings.HasSuffix(path, ext) {
			path = strings.TrimSuffix(path, ext)
			r.Find(method, path, c)
		}
	}
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
		if path[i] == paramLabel {
			if i > 0 && path[i-1] == '\\' {
				continue
			}
			uri.WriteString(`%v`)
			j := i + 1
			r.insert(rt.Method, path[:i], nil, staticKind, "", nil, -1)
			for ; i < l && path[i] != '/'; i++ {
			}
			pname := path[j:i]
			pnames = append(pnames, pname)
			path = path[:j] + path[i:]
			i, l = j, len(path)

			if i == l {
				// path node is last fragment of route path. ie. `/users/:id`
				r.insert(rt.Method, path[:i], rt.Handler, paramKind, ppath, pnames, rid)
			} else {
				r.insert(rt.Method, path[:i], nil, paramKind, "", nil, -1)
			}
		} else if path[i] == regexLabel {
			if i > 0 && path[i-1] == '\\' {
				continue
			}
			uri.WriteString(`%v`)
			j := i + 1
			r.insert(rt.Method, path[:i], nil, staticKind, "", nil, -1)
			for ; i < l && path[i] != '>'; i++ {
			}
			pname := path[j:i]
			parts := strings.SplitN(pname, `:`, 2)
			var regExpr string
			if len(parts) == 2 {
				pname = parts[0]
				regExpr = `(` + parts[1] + `)`
			} else {
				regExpr = `([^/]+)`
			}
			pnames = append(pnames, pname)
			if path[i] == '>' {
				i++
			}
			if len(path) > i {
				path = path[:j] + path[i:]
			} else {
				path = path[:j]
			}
			i, l = j, len(path)

			r.insert(rt.Method, path[:i], rt.Handler, regexKind, ppath, pnames, rid, regExpr)

		} else if path[i] == '*' {
			uri.WriteString(`%v`)
			r.insert(rt.Method, path[:i], nil, staticKind, "", nil, -1)
			pnames = append(pnames, "*")
			r.insert(rt.Method, path[:i+1], rt.Handler, anyKind, ppath, pnames, rid)
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
	r.insert(rt.Method, path, rt.Handler, staticKind, ppath, pnames, rid)
}

func (r *Router) insert(method, path string, h Handler, t kind, ppath string, pnames []string, rid int, regExpr ...string) {
	e := r.echo
	// Adjust max param
	paramLen := len(pnames)
	if *e.maxParam < paramLen {
		*e.maxParam = paramLen
	}

	currentNode := r.tree // Current node as root
	if currentNode == nil {
		panic("echo: invalid method")
	}
	search := path

	for {
		searchLen := len(search)
		prefixLen := len(currentNode.prefix)
		lcpLen := 0

		// LCP
		max := prefixLen
		if searchLen < max {
			max = searchLen
		}
		for ; lcpLen < max && search[lcpLen] == currentNode.prefix[lcpLen]; lcpLen++ {
		}

		if lcpLen == 0 {
			// At root node
			if len(search) > 0 {
				currentNode.label = search[0]
			}
			currentNode.prefix = search
			if h != nil {
				currentNode.kind = t
				currentNode.addHandler(method, h, rid)
				currentNode.ppath = ppath
				currentNode.pnames = pnames
			}
			currentNode.isLeaf = currentNode.IsLeaf()
		} else if lcpLen < prefixLen {
			// Split node
			n := newNode(
				currentNode.kind,
				currentNode.prefix[lcpLen:],
				currentNode,
				currentNode.staticChildren,
				currentNode.methodHandler,
				currentNode.ppath,
				currentNode.pnames,
				currentNode.regexChild,
				currentNode.paramChild,
				currentNode.anyChild,
				regExpr...,
			)
			// Update parent path for all children to new node
			for _, child := range currentNode.staticChildren {
				child.parent = n
			}
			if currentNode.paramChild != nil {
				currentNode.paramChild.parent = n
			}
			if currentNode.anyChild != nil {
				currentNode.anyChild.parent = n
			}

			// Reset parent node
			currentNode.kind = staticKind
			currentNode.label = currentNode.prefix[0]
			currentNode.prefix = currentNode.prefix[:lcpLen]
			currentNode.staticChildren = nil
			currentNode.methodHandler = new(methodHandler)
			currentNode.ppath = ""
			currentNode.pnames = nil
			currentNode.regExp = nil
			currentNode.paramChild = nil
			currentNode.anyChild = nil
			currentNode.isLeaf = false
			currentNode.isHandler = false

			currentNode.addStaticChild(n)

			if lcpLen == searchLen {
				// At parent node
				currentNode.kind = t
				currentNode.addHandler(method, h, rid)
				currentNode.ppath = ppath
				currentNode.pnames = pnames
			} else {
				// Create child node
				n = newNode(t, search[lcpLen:], currentNode, nil, new(methodHandler), ppath, pnames, nil, nil, nil, regExpr...)
				n.addHandler(method, h, rid)
				// Only Static children could reach here
				currentNode.addStaticChild(n)
			}
			currentNode.isLeaf = currentNode.IsLeaf()
		} else if lcpLen < searchLen {
			search = search[lcpLen:]
			c := currentNode.findChildWithLabel(search[0])
			if c != nil {
				// Go deeper
				currentNode = c
				continue
			}
			// Create child node
			n := newNode(t, search, currentNode, nil, new(methodHandler), ppath, pnames, nil, nil, nil, regExpr...)
			n.addHandler(method, h, rid)
			switch t {
			case staticKind:
				currentNode.addStaticChild(n)
			case regexKind:
				currentNode.regexChild = n
			case paramKind:
				currentNode.paramChild = n
			case anyKind:
				currentNode.anyChild = n
			}
			currentNode.isLeaf = currentNode.IsLeaf()
		} else {
			// Node already exists
			if h != nil {
				currentNode.addHandler(method, h, rid)
				currentNode.ppath = ppath
				if len(currentNode.pnames) == 0 {
					currentNode.pnames = pnames
				}
			}
		}
		return
	}
}

func newNode(t kind, pre string, p *node, sc children, mh *methodHandler, ppath string, pnames []string, regexChildren, paramChildren, anyChildren *node, regExpr ...string) *node {
	n := &node{
		kind:           t,
		label:          pre[0],
		prefix:         pre,
		parent:         p,
		staticChildren: sc,
		ppath:          ppath,
		pnames:         pnames,
		methodHandler:  mh,
		regexChild:     regexChildren,
		paramChild:     paramChildren,
		anyChild:       anyChildren,
		isHandler:      mh.isHandler(),
	}
	n.isLeaf = n.IsLeaf()
	if len(regExpr) > 0 && len(regExpr[0]) > 0 {
		if n.isLeaf && strings.HasSuffix(n.ppath, `:`+regExpr[0]+`>`) { // <name:regExpr>
			n.regExp = regexp.MustCompile(`^` + regExpr[0] + `$`)
		} else {
			n.regExp = regexp.MustCompile(`^` + regExpr[0])
		}
	}
	return n
}

func (n *node) String() string {
	return Dump(n.Tree(), false)
}

func (n *node) IsLeaf() bool {
	return n.staticChildren == nil && n.regexChild == nil && n.paramChild == nil && n.anyChild == nil
}

func (n *node) Tree() H {
	children := make([]H, len(n.staticChildren))
	for k, v := range n.staticChildren {
		children[k] = v.Tree()
	}
	var (
		regExpr    string
		regexChild H
		paramChild H
		anyChild   H
	)
	if n.regExp != nil {
		regExpr = n.regExp.String()
	}
	if n.regexChild != nil {
		regexChild = n.regexChild.Tree()
	}
	if n.paramChild != nil {
		paramChild = n.paramChild.Tree()
	}
	if n.anyChild != nil {
		anyChild = n.anyChild.Tree()
	}
	return H{
		"kind":           n.kind,
		"label":          string([]byte{n.label}),
		"prefix":         n.prefix,
		"parent":         n.parent,
		"staticChildren": children,
		"ppath":          n.ppath,
		"pnames":         n.pnames,
		"methodHandler":  n.methodHandler.Map(),
		"regExpr":        regExpr,
		"regexChild":     regexChild,
		"paramChild":     paramChild,
		"anyChild":       anyChild,
		"isLeaf":         n.isLeaf,
		"isHandler":      n.isHandler,
	}
}

func (n *node) addStaticChild(c *node) {
	n.staticChildren = append(n.staticChildren, c)
}

func (n *node) findStaticChild(l byte) *node {
	for _, c := range n.staticChildren {
		if c.label == l {
			return c
		}
	}
	return nil
}

func (n *node) findChildWithLabel(l byte) *node {
	for _, c := range n.staticChildren {
		if c.label == l {
			return c
		}
	}
	if l == regexLabel {
		return n.regexChild
	}
	if l == paramLabel {
		return n.paramChild
	}
	if l == anyLabel {
		return n.anyChild
	}
	return nil
}

func (n *node) findRegexChild(search string) (*node, []int) {
	var matchIndex []int
	c := n.regexChild
	matchIndex = c.regExp.FindStringSubmatchIndex(search)
	//fmt.Printf("%s => %s: %#v\n", search, c.regExp.String(), matchIndex)
	if len(matchIndex) > 3 {
		return c, matchIndex
	}
	return nil, matchIndex
}

func (n *node) addHandler(method string, h Handler, rid int) {
	n.methodHandler.addHandler(method, h, rid)
	if h != nil {
		n.isHandler = true
	} else {
		n.isHandler = n.methodHandler.isHandler()
	}
}

func (n *node) findHandler(method string) Handler {
	return n.methodHandler.findHandler(method)
}

func (n *node) find(method string) *endpoint {
	return n.methodHandler.find(method)
}

func (n *node) checkMethodNotAllowed() Handler {
	return n.methodHandler.checkMethodNotAllowed()
}

func (n *node) applyHandler(method string, ctx *xContext) {
	n.methodHandler.applyHandler(method, ctx)
	ctx.path = n.ppath
	ctx.pnames = n.pnames
}

func (r *Router) Tree() H {
	return r.tree.Tree()
}

func (r *Router) String() string {
	return r.tree.String()
}

func (r *Router) Find(method, path string, context Context) (found bool) {
	ctx := context.Object()
	ctx.path = path
	currentNode := r.tree // Current node as root

	if m, ok := r.static[path]; ok {
		m.applyHandler(method, ctx)
		if ctx.handler == nil {
			ctx.handler = m.checkMethodNotAllowed()
			return false
		}
		return true
	}

	var (
		previousBestMatchNode *node
		matchedEndpoint       *endpoint
		// search stores the remaining path to check for match. By each iteration we move from start of path to end of the path
		// and search value gets shorter and shorter.
		search      = path
		searchIndex = 0
		paramIndex  int // Param counter
		paramValues = context.ParamValues()
		matchIndex  []int
	)

	// Backtracking is needed when a dead end (leaf node) is reached in the router tree.
	// To backtrack the current node will be changed to the parent node and the next kind for the
	// router logic will be returned based on fromKind or kind of the dead end node (static > param > any).
	// For example if there is no static node match we should check parent next sibling by kind (param).
	// Backtracking itself does not check if there is a next sibling, this is done by the router logic.
	backtrackToNextNodeKind := func(fromKind kind) (nextNodeKind kind, valid bool) {
		previous := currentNode
		currentNode = previous.parent
		valid = currentNode != nil

		// Next node type by priority
		if previous.kind == anyKind {
			nextNodeKind = staticKind
		} else {
			nextNodeKind = previous.kind + 1
		}

		if fromKind == staticKind {
			// when backtracking is done from static kind block we did not change search so nothing to restore
			return
		}

		// restore search to value it was before we move to current node we are backtracking from.
		if previous.kind == staticKind {
			searchIndex -= len(previous.prefix)
		} else {
			paramIndex--
			// for param/any node.prefix value is always `:` so we can not deduce searchIndex from that and must use pValue
			// for that index as it would also contain part of path we cut off before moving into node we are backtracking from
			searchIndex -= len(paramValues[paramIndex])
			paramValues[paramIndex] = ""
		}
		search = path[searchIndex:]
		return
	}

	// Router tree is implemented by longest common prefix array (LCP array) https://en.wikipedia.org/wiki/LCP_array
	// Tree search is implemented as for loop where one loop iteration is divided into 3 separate blocks
	// Each of these blocks checks specific kind of node (static/param/any). Order of blocks reflex their priority in routing.
	// Search order/priority is: static > param > any.
	//
	// Note: backtracking in tree is implemented by replacing/switching currentNode to previous node
	// and hoping to (goto statement) next block by priority to check if it is the match.
	for {
		prefixLen := 0 // Prefix length
		lcpLen := 0    // LCP length

		if currentNode.kind == staticKind {
			searchLen := len(search)
			prefixLen = len(currentNode.prefix)

			// LCP - Longest Common Prefix (https://en.wikipedia.org/wiki/LCP_array)
			max := prefixLen
			if searchLen < max {
				max = searchLen
			}
			for ; lcpLen < max && search[lcpLen] == currentNode.prefix[lcpLen]; lcpLen++ {
			}
		}

		if lcpLen != prefixLen {
			// No matching prefix, let's backtrack to the first possible alternative node of the decision path
			nk, ok := backtrackToNextNodeKind(staticKind)
			if !ok {
				return // No other possibilities on the decision path
			}
			if nk == regexKind {
				goto Regex
			} else if nk == paramKind {
				goto Param
				// NOTE: this case (backtracking from static node to previous any node) can not happen by current any matching logic. Any node is end of search currently
				//} else if nk == anyKind {
				//	goto Any
			} else {
				// Not found (this should never be possible for static node we are looking currently)
				break
			}
		}

		// The full prefix has matched, remove the prefix from the remaining search
		search = search[lcpLen:]
		searchIndex = searchIndex + lcpLen
		//fmt.Println(`search:`, search, currentNode.String())
		// Finish routing if no remaining search and we are on a node with handler and matching method type
		if search == "" && currentNode.isHandler {
			// check if current node has handler registered for http method we are looking for. we store currentNode as
			// best matching in case we do no find no more routes matching this path+method
			if previousBestMatchNode == nil {
				previousBestMatchNode = currentNode
			}
			if matchedEndpoint = currentNode.find(method); matchedEndpoint != nil {
				break
			}
		}

		// Static node
		if search != "" {
			if child := currentNode.findStaticChild(search[0]); child != nil {
				currentNode = child
				continue
			}
		}

	Regex:
		// Regex node
		if child := currentNode.regexChild; search != "" && child != nil {
			child, matchIndex = currentNode.findRegexChild(search)
			if child != nil {
				currentNode = child
				startIndex := matchIndex[2]
				endIndex := matchIndex[3]
				paramValues[paramIndex] = search[startIndex:endIndex]
				paramIndex++
				endIndex = matchIndex[1]
				search = search[endIndex:]
				searchIndex = searchIndex + endIndex
				continue
			}
		}

	Param:
		// Param node
		if child := currentNode.paramChild; search != "" && child != nil {
			currentNode = child
			i := 0
			l := len(search)
			if currentNode.isLeaf {
				// when param node does not have any children then param node should act similarly to any node - consider all remaining search as match
				i = l
			} else {
				for ; i < l && search[i] != '/'; i++ {
				}
			}

			paramValues[paramIndex] = search[:i]
			paramIndex++
			search = search[i:]
			searchIndex = searchIndex + i
			continue
		}

	Any:
		// Any node
		if child := currentNode.anyChild; child != nil {
			// If any node is found, use remaining path for paramValues
			currentNode = child
			paramValues[len(currentNode.pnames)-1] = search
			// update indexes/search in case we need to backtrack when no handler match is found
			paramIndex++
			searchIndex += +len(search)
			search = ""

			// check if current node has handler registered for http method we are looking for. we store currentNode as
			// best matching in case we do no find no more routes matching this path+method
			if previousBestMatchNode == nil {
				previousBestMatchNode = currentNode
			}
			if matchedEndpoint = currentNode.find(method); matchedEndpoint != nil {
				break
			}
		}

		// Let's backtrack to the first possible alternative node of the decision path
		nk, ok := backtrackToNextNodeKind(anyKind)
		if !ok {
			break // No other possibilities on the decision path
		} else if nk == regexKind {
			goto Regex
		} else if nk == paramKind {
			goto Param
		} else if nk == anyKind {
			goto Any
		} else {
			// Not found
			break
		}
	}

	if currentNode == nil && previousBestMatchNode == nil {
		return // nothing matched at all
	}

	if matchedEndpoint != nil {
		ctx.handler = matchedEndpoint.handler
		ctx.rid = matchedEndpoint.rid
		found = true
	} else if previousBestMatchNode != nil {
		// use previous match as basis. although we have no matching handler we have path match.
		// so we can send http.StatusMethodNotAllowed (405) instead of http.StatusNotFound (404)
		currentNode = previousBestMatchNode
		ctx.handler = currentNode.checkMethodNotAllowed()
	}
	ctx.path = currentNode.ppath
	ctx.pnames = currentNode.pnames
	//currentNode.applyHandler(method, ctx)
	if ctx.handler == nil {
		ctx.handler = currentNode.checkMethodNotAllowed()
	}
	return
}
