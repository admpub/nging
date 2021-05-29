package middleware

import (
	"regexp"
	"strings"

	"github.com/webx-top/echo"
)

type (
	// RewriteConfig defines the config for Rewrite middleware.
	RewriteConfig struct {
		// Skipper defines a function to skip middleware.
		Skipper echo.Skipper `json:"-"`

		// Rules defines the URL path rewrite rules. The values captured in asterisk can be
		// retrieved by index e.g. $1, $2 and so on.
		// Example:
		// "/old":              "/new",
		// "/api/*":            "/$1",
		// "/js/*":             "/public/javascripts/$1",
		// "/users/*/orders/*": "/user/$1/order/$2",
		// Required.
		Rules map[string]string `json:"rules"`

		addresses  map[string]*regexp.Regexp //"/old": *regexp.Regexp
		rulesRegex map[*regexp.Regexp]string //*regexp.Regexp: "/new"
	}
)

// Init Initialize
func (c *RewriteConfig) Init() *RewriteConfig {
	c.rulesRegex = map[*regexp.Regexp]string{}
	c.addresses = map[string]*regexp.Regexp{}

	for k, v := range c.Rules {
		c.Set(k, v)
	}
	return c
}

var queryParamRegex = regexp.MustCompile(`/:[^/]+`)

func QueryParamToRegexpRule(query string) string {
	query = strings.Replace(query, "*", "(\\S*)", -1)
	query = queryParamRegex.ReplaceAllString(query, "/([^/]+)")
	return query
}

// Set rule
func (c *RewriteConfig) Set(urlPath, newPath string) *RewriteConfig {
	re, ok := c.addresses[urlPath]
	if ok {
		delete(c.rulesRegex, re)
		r := QueryParamToRegexpRule(urlPath)
		re = regexp.MustCompile(r)
		c.rulesRegex[re] = newPath
		c.addresses[urlPath] = re
	} else {
		c.Add(urlPath, newPath)
	}
	return c
}

// Delete rule
func (c *RewriteConfig) Delete(urlPath string) *RewriteConfig {
	re, ok := c.addresses[urlPath]
	if ok {
		delete(c.rulesRegex, re)
		delete(c.addresses, urlPath)
	}
	return c
}

// Add rule
func (c *RewriteConfig) Add(urlPath, newPath string) *RewriteConfig {
	r := QueryParamToRegexpRule(urlPath)
	re := regexp.MustCompile(r)
	c.rulesRegex[re] = newPath
	c.addresses[urlPath] = re
	return c
}

// Rewrite url
func (c *RewriteConfig) Rewrite(urlPath string) string {
	for k, v := range c.rulesRegex {
		replacer := echo.CaptureTokens(k, urlPath)
		if replacer != nil {
			urlPath = replacer.Replace(v)
		}
	}
	return urlPath
}

var (
	// DefaultRewriteConfig is the default Rewrite middleware config.
	DefaultRewriteConfig = RewriteConfig{
		Skipper: echo.DefaultSkipper,
	}
)

// Rewrite returns a Rewrite middleware.
//
// Rewrite middleware rewrites the URL path based on the provided rules.
func Rewrite(rules map[string]string) echo.MiddlewareFuncd {
	c := DefaultRewriteConfig
	c.Rules = rules
	return RewriteWithConfig(c)
}

// RewriteWithConfig returns a Rewrite middleware with config.
// See: `Rewrite()`.
func RewriteWithConfig(config RewriteConfig) echo.MiddlewareFuncd {
	// Defaults
	if config.Rules == nil {
		panic("echo: rewrite middleware requires url path rewrite rules")
	}
	if config.Skipper == nil {
		config.Skipper = DefaultRewriteConfig.Skipper
	}
	config.Init()
	return func(next echo.Handler) echo.HandlerFunc {
		return func(c echo.Context) (err error) {
			if config.Skipper(c) {
				return next.Handle(c)
			}

			req := c.Request()
			req.URL().SetPath(config.Rewrite(req.URL().Path()))
			return next.Handle(c)
		}
	}
}
