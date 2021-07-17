package locale

import (
	"net/http"
	"strings"

	"github.com/admpub/caddy/caddyhttp/httpserver"

	"github.com/caddy-plugins/caddy-locale/method"
)

// Middleware is a httpserver to detect the user's locale.
type Middleware struct {
	Next             httpserver.Handler
	AvailableLocales []string
	Methods          []method.Method
	PathScope        string
	Configuration    *method.Configuration
}

// ServeHTTP implements the httpserver.Handler interface.
func (l *Middleware) ServeHTTP(w http.ResponseWriter, r *http.Request) (int, error) {
	if !httpserver.Path(r.URL.Path).Matches(l.PathScope) {
		return l.Next.ServeHTTP(w, r)
	}

	candidates := []string{}
	for _, method := range l.Methods {
		candidates = append(candidates, method(r, l.Configuration)...)
	}

	locale := l.firstValid(candidates)
	if locale == "" {
		locale = l.defaultLocale()
	}
	r.Header.Set("Detected-Locale", locale)

	return l.Next.ServeHTTP(w, r)
}

func (l *Middleware) defaultLocale() string {
	return l.AvailableLocales[0]
}

func (l *Middleware) firstValid(candidates []string) string {
	for _, candidate := range candidates {
		candidate = strings.TrimSpace(candidate)
		if val := l.validAvailableLocale(candidate); val != "" {
			return val
		}
	}
	return ""
}

func (l *Middleware) validAvailableLocale(locale string) string {
	locale = strings.ToLower(locale)
	for _, validLocale := range l.AvailableLocales {
		if locale == strings.ToLower(validLocale) {
			return validLocale
		}
	}
	return ""
}
