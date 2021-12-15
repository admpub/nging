package echo

type (
	Host struct {
		head   Handler
		group  *Group
		groups map[string]*Group
		Router *Router
	}
	TypeHost struct {
		prefix string
		router *Router
		echo   *Echo
	}
)

func (t TypeHost) URI(handler interface{}, params ...interface{}) string {
	if t.router == nil || t.echo == nil {
		return ``
	}
	return t.prefix + t.echo.URI(handler, params...)
}

func (t TypeHost) String() string {
	return t.prefix
}

func (h *Host) Host(args ...interface{}) (r TypeHost) {
	if h.group == nil || h.group.host == nil {
		return
	}
	r.echo = h.group.echo
	r.router = h.Router
	if len(args) != 1 {
		r.prefix = h.group.host.Format(args...)
		return
	}
	switch v := args[0].(type) {
	case map[string]interface{}:
		r.prefix = h.group.host.FormatMap(v)
	case H:
		r.prefix = h.group.host.FormatMap(v)
	default:
		r.prefix = h.group.host.Format(args...)
	}
	return
}
