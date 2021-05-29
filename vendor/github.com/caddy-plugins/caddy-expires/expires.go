package expires

import (
	"net/http"
	"regexp"
	"strconv"
	"time"

	"github.com/caddyserver/caddy"
	"github.com/caddyserver/caddy/caddyhttp/httpserver"
)

type matchDef struct {
	re       *regexp.Regexp
	duration time.Duration
}

func (m *matchDef) Match(header http.Header, request *http.Request) bool {
	return m.re.MatchString(request.URL.Path)
}

func (m *matchDef) Duration() time.Duration {
	return m.duration
}

func (m *matchDef) Parse(args []string) error {
	re, err := regexp.Compile(args[0])
	if err != nil {
		return err
	}
	m.re = re
	m.duration = parseDuration(args[1])
	return nil
}

type headerMatchDef struct {
	header   string
	re       *regexp.Regexp
	duration time.Duration
}

func (m *headerMatchDef) Match(header http.Header, request *http.Request) bool {
	return m.re.MatchString(header.Get(m.header))
}

func (m *headerMatchDef) Duration() time.Duration {
	return m.duration
}

func (m *headerMatchDef) Parse(args []string) error {
	m.header = args[0]
	re, err := regexp.Compile(args[1])
	if err != nil {
		return err
	}
	m.re = re
	m.duration = parseDuration(args[2])
	return nil
}

type matchRule interface {
	Duration() time.Duration
	Match(http.Header, *http.Request) bool
	Parse([]string) error
}

func init() {
	caddy.RegisterPlugin("expires", caddy.Plugin{
		ServerType: "http",
		Action:     setup,
	})
}

func setup(c *caddy.Controller) error {
	rules, err := parseRules(c)
	if err != nil {
		return err
	}

	cfg := httpserver.GetConfig(c)
	mid := func(next httpserver.Handler) httpserver.Handler {
		return expiresHandler{Next: next, Rules: rules}
	}
	cfg.AddMiddleware(mid)

	return nil
}

func parseRules(c *caddy.Controller) ([]matchRule, error) {
	rules := []matchRule{}

	for c.Next() {
		for c.NextBlock() {
			switch c.Val() {
			case "match":
				args := c.RemainingArgs()
				if len(args) != 2 {
					return nil, c.ArgErr()
				}
				rule := &matchDef{}
				rule.Parse(args)
				rules = append(rules, rule)
			case "match_header":
				args := c.RemainingArgs()
				if len(args) != 3 {
					return nil, c.ArgErr()
				}
				rule := &headerMatchDef{}
				rule.Parse(args)
				rules = append(rules, rule)
			default:
				return nil, c.SyntaxErr("match")
			}
		}
	}
	return rules, nil
}

func parseDuration(str string) time.Duration {
	durationRegex := regexp.MustCompile(`(?P<years>\d+y)?(?P<months>\d+m)?(?P<days>\d+d)?T?(?P<hours>\d+h)?(?P<minutes>\d+i)?(?P<seconds>\d+s)?`)
	matches := durationRegex.FindStringSubmatch(str)

	years := parseInt64(matches[1])
	months := parseInt64(matches[2])
	days := parseInt64(matches[3])
	hours := parseInt64(matches[4])
	minutes := parseInt64(matches[5])
	seconds := parseInt64(matches[6])

	hour := int64(time.Hour)
	minute := int64(time.Minute)
	second := int64(time.Second)
	return time.Duration(years*24*365*hour + months*30*24*hour + days*24*hour + hours*hour + minutes*minute + seconds*second)
}

func parseInt64(value string) int64 {
	if len(value) == 0 {
		return 0
	}
	parsed, err := strconv.Atoi(value[:len(value)-1])
	if err != nil {
		return 0
	}
	return int64(parsed)
}

type expiresHandler struct {
	Next  httpserver.Handler
	Rules []matchRule
}

func (h expiresHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) (int, error) {
	for _, rule := range h.Rules {
		if rule.Match(w.Header(), r) {
			w.Header().Set("Expires", time.Now().Add(rule.Duration()).UTC().Format(time.RFC1123))
			w.Header().Set("Cache-Control", "public, max-age="+strconv.Itoa(int(rule.Duration().Seconds())))
			break
		}
	}
	return h.Next.ServeHTTP(w, r)
}
