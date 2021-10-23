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

func ParseURIRegExp(uriRegexp string, dflRegexp string) (names []string, format string, regExp *regexp.Regexp) {
	matches := hostRegExp.FindAllStringSubmatchIndex(uriRegexp, -1)
	if len(matches) == 0 {
		return
	}
	if len(dflRegexp) == 0 {
		dflRegexp = `[^.]+`
	}
	var regExpr string
	var lastPosition int
	for _, matchIndex := range matches {
		if matchIndex[0] > 0 {
			v := uriRegexp[lastPosition:matchIndex[0]]
			format += v
			regExpr += regexp.QuoteMeta(v)
		}
		lastPosition = matchIndex[1]
		format += `%v`
		name := uriRegexp[matchIndex[2]:matchIndex[3]]
		if matchIndex[4] > 0 {
			regExpr += `(` + uriRegexp[matchIndex[4]:matchIndex[5]] + `)`
		} else {
			regExpr += `(` + dflRegexp + `)`
		}
		names = append(names, name)
	}
	if lastPosition > 0 {
		if lastPosition < len(uriRegexp) {
			v := uriRegexp[lastPosition:]
			format += v
			regExpr += regexp.QuoteMeta(v)
		}
	}
	regExp = regexp.MustCompile(`^` + regExpr + `$`)
	return
}

func (h *host) Parse() *host {
	h.names, h.format, h.regExp = ParseURIRegExp(h.name, ``)
	return h
}
