package locale

import (
	"strings"

	"github.com/admpub/caddy"
	"github.com/admpub/caddy/caddyhttp/httpserver"

	"github.com/caddy-plugins/caddy-locale/method"
)

func init() {
	caddy.RegisterPlugin("locale", caddy.Plugin{
		ServerType: "http",
		Action:     setup,
	})
}

func setup(c *caddy.Controller) error {
	l, err := parseLocale(c)
	if err != nil {
		return err
	}

	siteConfig := httpserver.GetConfig(c)

	siteConfig.AddMiddleware(func(next httpserver.Handler) httpserver.Handler {
		l.Next = next
		return l
	})

	return nil
}

func parseLocale(c *caddy.Controller) (*Middleware, error) {
	result := &Middleware{
		AvailableLocales: []string{},
		Methods:          []method.Method{},
		PathScope:        "/",
		Configuration: &method.Configuration{
			CookieName: "locale",
		},
	}

	for c.Next() {
		args := c.RemainingArgs()

		if len(args) > 0 {
			result.AvailableLocales = append(result.AvailableLocales, args...)
		}

		for c.NextBlock() {
			switch c.Val() {
			case "available":
				result.AvailableLocales = append(result.AvailableLocales, c.RemainingArgs()...)
			case "detect":
				detectArgs := c.RemainingArgs()
				if len(detectArgs) == 0 {
					return nil, c.ArgErr()
				}
				for _, detectArg := range detectArgs {
					method, found := method.Names[strings.ToLower(strings.TrimSpace(detectArg))]
					if !found {
						return nil, c.Errf("could not find detect method [%s]", detectArg)
					}
					result.Methods = append(result.Methods, method)
				}
			case "cookie":
				if !c.NextArg() {
					return nil, c.ArgErr()
				}
				if value := strings.TrimSpace(c.Val()); value != "" {
					result.Configuration.CookieName = value
				}
			case "path":
				if !c.NextArg() {
					return nil, c.ArgErr()
				}
				if value := strings.TrimSpace(c.Val()); value != "" {
					result.PathScope = value
				}
			default:
				return nil, c.ArgErr()
			}
		}
	}

	if len(result.AvailableLocales) == 0 {
		return nil, c.Errf("no available locales specified")
	}

	if len(result.Methods) == 0 {
		result.Methods = append(result.Methods, method.Names["header"])
	}

	return result, nil
}
