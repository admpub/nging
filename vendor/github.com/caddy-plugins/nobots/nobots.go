package nobots

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"

	"github.com/caddyserver/caddy"
	"github.com/caddyserver/caddy/caddyhttp/httpserver"
)

// botUA config representation
type botUA struct {
	uas    []string
	bomb   string
	re     []*regexp.Regexp
	public []*regexp.Regexp
}

// BotUA plugin struct
type BotUA struct {
	Next httpserver.Handler
	UA   *botUA
}

func init() {
	caddy.RegisterPlugin("nobots", caddy.Plugin{
		ServerType: "http",
		Action:     setup,
	})
}

// setup callback for parsing the config
func setup(c *caddy.Controller) error {
	ua, err := parseUA(c)
	if err != nil {
		return err
	}

	// Verfies whether bomb exist
	if _, err := os.Stat(ua.bomb); os.IsNotExist(err) {
		return fmt.Errorf("Bomb %s not found.", ua.bomb)
	}

	// Setup de middleware
	cfg := httpserver.GetConfig(c)
	mid := func(next httpserver.Handler) httpserver.Handler {
		return BotUA{Next: next, UA: ua}
	}
	cfg.AddMiddleware(mid)

	return nil
}

// parseUA propper config parser that generates a botUA object
func parseUA(c *caddy.Controller) (*botUA, error) {
	var ua botUA
	for c.Next() {
		if !c.NextArg() {
			return nil, c.ArgErr()
		}
		ua.bomb = c.Val()
		for c.NextBlock() {
			switch c.Val() {
			case "regexp":
				if !c.NextArg() {
					return nil, c.ArgErr()
				}
				re, err := regexp.Compile(c.Val())
				if err != nil {
					return nil, fmt.Errorf("%s", err)
				}
				ua.re = append(ua.re, re)
			case "public":
				if !c.NextArg() {
					return nil, c.ArgErr()
				}
				re, err := regexp.Compile(c.Val())
				if err != nil {
					return nil, fmt.Errorf("%s", err)
				}
				ua.public = append(ua.public, re)
			default:
				ua.uas = append(ua.uas, c.Val())
			}
		}
	}
	return &ua, nil
}

func (b BotUA) ServeHTTP(w http.ResponseWriter, r *http.Request) (int, error) {
	// Get request UA
	rua := r.UserAgent()

	// Avoid ban UA for public URI
	if !b.IsPublicURI(r.URL.Path) {
		// Check if the UA is a evil one
		if b.IsEvil(rua) {
			return serveBomb(w, r, b.UA.bomb)
		}
	}
	// Nothing happens carry on with next stuff
	return b.Next.ServeHTTP(w, r)
}

// IsEvil check the remote UA against evil UAs
func (b BotUA) IsEvil(rua string) bool {
	// In case there are regexp
	if len(b.UA.re) > 0 {
		for _, re := range b.UA.re {
			if re.MatchString(rua) {
				return true
			}
		}
	}
	// In case there are strings
	if len(b.UA.uas) > 0 {
		for _, ua := range b.UA.uas {
			if ua == rua {
				return true
			}
		}
	}
	// UA is not evil
	return false
}

// IsPublicURI check if the requested URI is defined as public or not
func (b BotUA) IsPublicURI(uri string) bool {
	if len(b.UA.public) > 0 {
		for _, re := range b.UA.public {
			if re.MatchString(uri) {
				return true
			}
		}
	}
	return false
}

// serveBomb provides the bomb to front-end
func serveBomb(w http.ResponseWriter, r *http.Request, bomb string) (int, error) {
	file, err := ioutil.ReadFile(bomb)
	if err != nil {
		return http.StatusNotFound, nil
	}

	w.Header().Set("Content-Type", "text/html; charset=UTF-8")
	w.Header().Set("Content-Encoding", "gzip")
	w.Header().Set("Content-Length", fmt.Sprintf("%d", len(file)))
	w.Write(file)
	return 200, nil
}
