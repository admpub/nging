package echo

import (
	"fmt"
	"regexp"
)

type Hoster interface {
	Name() string
	Alias() string
	Format(args ...interface{}) string
	FormatMap(params H) string
	RegExp() *regexp.Regexp
	Match(host string) (r []string, hasExpr bool)
}

type host struct {
	name   string
	alias  string
	format string
	regExp *regexp.Regexp
	names  []string
}

func (h *host) Name() string {
	return h.name
}

func (h *host) Alias() string {
	return h.alias
}

func (h *host) Format(args ...interface{}) string {
	if len(args) > 0 {
		return fmt.Sprintf(h.format, args...)
	}
	return h.format
}

func (h *host) FormatMap(params H) string {
	if len(params) > 0 {
		args := make([]interface{}, len(h.names))
		for index, name := range h.names {
			v, y := params[name]
			if y {
				args[index] = v
			} else {
				args[index] = ``
			}
		}
		return fmt.Sprintf(h.format, args...)
	}
	return h.format
}

func (h *host) RegExp() *regexp.Regexp {
	return h.regExp
}

func (h *host) Match(host string) (r []string, hasExpr bool) {
	if h.regExp != nil {
		match := h.regExp.FindStringSubmatch(host)
		if len(match) > 0 {
			return match[1:], true
		}
		return nil, true
	}
	return nil, false
}

func NewHost(name string) *host {
	return &host{name: name}
}

var hostRegExp = regexp.MustCompile(`<([^:]+)(?:\:(.+?))?>`)

func (h *host) Parse() *host {
	matches := hostRegExp.FindAllStringSubmatchIndex(h.name, -1)
	if len(matches) == 0 {
		return h
	}
	var format string
	var regExpr string
	var lastPosition int
	for _, matchIndex := range matches {
		if matchIndex[0] > 0 {
			v := h.name[lastPosition:matchIndex[0]]
			format += v
			regExpr += regexp.QuoteMeta(v)
		}
		lastPosition = matchIndex[1]
		format += `%v`
		name := h.name[matchIndex[2]:matchIndex[3]]
		if matchIndex[4] > 0 {
			regExpr += `(` + h.name[matchIndex[4]:matchIndex[5]] + `)`
		} else {
			regExpr += `([^.]+)`
		}
		h.names = append(h.names, name)
	}
	if lastPosition > 0 {
		if lastPosition < len(h.name) {
			v := h.name[lastPosition:]
			format += v
			regExpr += regexp.QuoteMeta(v)
		}
	}
	h.format = format
	h.regExp = regexp.MustCompile(`^` + regExpr + `$`)
	return h
}
