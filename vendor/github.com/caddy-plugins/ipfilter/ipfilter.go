package ipfilter

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/mholt/caddy"
	"github.com/mholt/caddy/caddyhttp/httpserver"
	"github.com/oschwald/maxminddb-golang"
)

// IPFilter is a middleware for filtering clients based on their ip or country's ISO code.
type IPFilter struct {
	Next   httpserver.Handler
	Config IPFConfig
}

// IPPath holds the configuration of a single ipfilter block.
type IPPath struct {
	PathScopes   []string
	BlockPage    string
	CountryCodes []string
	PrefixDir    string
	Nets         []*net.IPNet
	IsBlock      bool
	Strict       bool
}

// IPFConfig holds the configuration for the ipfilter middleware.
type IPFConfig struct {
	Paths     []IPPath
	DBHandler *maxminddb.Reader // Database's handler if it gets opened.
}

// OnlyCountry is used to fetch only the country's code from 'mmdb'.
type OnlyCountry struct {
	Country struct {
		ISOCode string `maxminddb:"iso_code"`
	} `maxminddb:"country"`
}

// Status is used to keep track of the status of the request.
type Status struct {
	countryMatch, inRange bool
}

// Any returns 'true' if we have a match on a country code or an IP in range.
func (s *Status) Any() bool {
	return s.countryMatch || s.inRange
}

// block will take care of blocking
func block(blockPage string, w http.ResponseWriter) (int, error) {
	if blockPage != "" {
		bp, err := os.Open(blockPage)
		if err != nil {
			return http.StatusInternalServerError, err
		}
		defer bp.Close()

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		if _, err := io.Copy(w, bp); err != nil {
			return http.StatusInternalServerError, err
		}
		// we wrote the blockpage, return OK.
		return http.StatusOK, nil
	}

	// if we don't have blockpage, return forbidden.
	return http.StatusForbidden, nil
}

// Init initializes the plugin
func init() {
	caddy.RegisterPlugin("ipfilter", caddy.Plugin{
		ServerType: "http",
		Action:     Setup,
	})
}

// Setup parses the ipfilter configuration and returns the middleware handler.
func Setup(c *caddy.Controller) error {
	ifconfig, err := ipfilterParse(c)
	if err != nil {
		return err
	}

	// Create new middleware
	newMiddleWare := func(next httpserver.Handler) httpserver.Handler {
		return &IPFilter{
			Next:   next,
			Config: ifconfig,
		}
	}
	// Add middleware
	cfg := httpserver.GetConfig(c)
	cfg.AddMiddleware(newMiddleWare)

	return nil
}

func getClientIP(r *http.Request, strict bool) (net.IP, error) {
	var ip string

	// Use the client ip from the 'X-Forwarded-For' header, if available.
	if fwdFor := r.Header.Get("X-Forwarded-For"); fwdFor != "" && !strict {
		ips := strings.Split(fwdFor, ", ")
		ip = ips[0]
	} else {
		// Otherwise, get the client ip from the request remote address.
		var err error
		ip, _, err = net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			return nil, err
		}
	}

	// Parse the ip address string into a net.IP.
	parsedIP := net.ParseIP(ip)
	if parsedIP == nil {
		return nil, errors.New("unable to parse address")
	}

	return parsedIP, nil
}

// ShouldAllow takes a path and a request and decides if it should be allowed
func (ipf IPFilter) ShouldAllow(path IPPath, r *http.Request) (bool, string, error) {
	allow := true
	scopeMatched := ""

	// check if we are in one of our scopes.
	for _, scope := range path.PathScopes {
		if httpserver.Path(r.URL.Path).Matches(scope) {
			// extract the client's IP and parse it.
			clientIP, err := getClientIP(r, path.Strict)
			if err != nil {
				return false, scope, err
			}

			// request status.
			var rs Status

			if len(path.CountryCodes) != 0 {
				// do the lookup.
				var result OnlyCountry
				if err = ipf.Config.DBHandler.Lookup(clientIP, &result); err != nil {
					return false, scope, err
				}

				// get only the ISOCode out of the lookup results.
				clientCountry := result.Country.ISOCode
				for _, c := range path.CountryCodes {
					if clientCountry == c {
						rs.countryMatch = true
						break
					}
				}
			}

			if len(path.Nets) != 0 {
				for _, rng := range path.Nets {
					if rng.Contains(clientIP) {
						rs.inRange = true
						break
					}
				}
			}

			if ipf.PrefixDirBlocked(clientIP, path) {
				rs.inRange = true
			}

			scopeMatched = scope
			if rs.Any() {
				// Rule matched, if the rule has IsBlock = true then we have to deny access
				allow = !path.IsBlock
			} else {
				// Rule did not match, if the rule has IsBlock = true then we have to allow access
				allow = path.IsBlock
			}

			// We only have to test the first path that matches because it is the most specific
			break
		}
	}

	// no scope match, pass-through.
	return allow, scopeMatched, nil
}

// PrefixDirBlocked takes an IP and a path and decides to allow or block based on prefix_dir.
func (ipf IPFilter) PrefixDirBlocked(clientIP net.IP, path IPPath) bool {
	if path.PrefixDir == "" {
		return false
	}

	fname := clientIP.String()
	fname_variant := ""
	is_ipv6 := clientIP.To4() == nil
	if is_ipv6 {
		fname_variant = strings.ReplaceAll(fname, ":", "=")
	}

	// Check the "flat" namespace.
	blacklistPath := filepath.Join(path.PrefixDir, fname)
	if _, err := os.Stat(blacklistPath); err == nil {
		return true
	}
	if is_ipv6 {
		blacklistPath := filepath.Join(path.PrefixDir, fname_variant)
		if _, err := os.Stat(blacklistPath); err == nil {
			return true
		}
	}

	// Check the "sharded" namespace.
	c := strings.SplitN(fname, ".", 3) // shard IPv4 address
	if len(c) != 3 {
		c = strings.SplitN(fname, ":", 3) // shard IPv6 address
		if len(c) != 3 {
			// This should be a "can't happen" situation. Perhaps there is an
			// IP address type we don't know how to shard. But rather than
			// blow up below just log the problem and grant access.
			log.Println("ipfilter: Could not shard address:", fname)
			return false
		}
	}
	blacklistPath = filepath.Join(path.PrefixDir, c[0], c[1], fname)
	if _, err := os.Stat(blacklistPath); err == nil {
		return true
	}
	if is_ipv6 {
		blacklistPath = filepath.Join(path.PrefixDir, c[0], c[1], fname_variant)
		if _, err := os.Stat(blacklistPath); err == nil {
			return true
		}
	}

	return false
}

func (ipf IPFilter) ServeHTTP(w http.ResponseWriter, r *http.Request) (int, error) {
	allow := true
	matchedPath := ""
	blockPage := ""

	// Loop over all IPPaths in the config
	for _, path := range ipf.Config.Paths {
		pathAllow, pathMathedPath, err := ipf.ShouldAllow(path, r)
		if err != nil {
			return http.StatusInternalServerError, err
		}

		if len(pathMathedPath) >= len(matchedPath) {
			allow = pathAllow
			matchedPath = pathMathedPath
			blockPage = path.BlockPage
		}
	}

	if !allow {
		return block(blockPage, w)
	}
	return ipf.Next.ServeHTTP(w, r)
}

// parseIP parses a string to an IP range.
func parseIP(ip string) ([]*net.IPNet, error) {
	// CIDR notation
	_, ipnet, err := net.ParseCIDR(ip)
	if err == nil {
		return []*net.IPNet{ipnet}, nil
	}

	// Singular IP
	parsedIP := net.ParseIP(ip)
	if parsedIP != nil {
		mask := len(parsedIP) * 8
		return []*net.IPNet{{
			IP:   parsedIP,
			Mask: net.CIDRMask(mask, mask),
		}}, nil
	}

	// for backward compatibility, convert ranges into CIDR notation.
	parseError := fmt.Errorf("Can't parse IP: %s", ip)
	// check if the ip isn't complete;
	// e.g. 192.168 -> Range{"192.168.0.0", "192.168.255.255"}
	dotSplit := strings.Split(ip, ".")
	if len(dotSplit) < 4 {
		startR := make([]string, len(dotSplit), 4)
		copy(startR, dotSplit)
		for len(dotSplit) < 4 {
			startR = append(startR, "0")
			dotSplit = append(dotSplit, "255")
		}
		start := net.ParseIP(strings.Join(startR, "."))
		end := net.ParseIP(strings.Join(dotSplit, "."))
		if start.To4() == nil || end.To4() == nil {
			return nil, parseError
		}

		return range2CIDRs(start, end), nil
	}

	// try to split on '-' to see if it is a range of ips e.g. 1.1.1.1-10
	splitted := strings.Split(ip, "-")
	if len(splitted) > 1 { // if more than one, then we got a range e.g. ["1.1.1.1", "10"]
		start := net.ParseIP(splitted[0])
		// make sure that we got a valid IPv4 IP.
		if start.To4() == nil {
			return nil, parseError
		}

		// split the start of the range on "." and switch the last field with splitted[1], e.g 1.1.1.1 -> 1.1.1.10
		fields := strings.Split(start.String(), ".")
		fields[3] = splitted[1]
		end := net.ParseIP(strings.Join(fields, "."))

		// parse the end range.
		if end.To4() == nil {
			return nil, parseError
		}

		return range2CIDRs(start, end), nil
	}

	// Failed to parse IP
	return nil, parseError
}

// ipfilterParseSingle parses a single ipfilter {} block from the caddy config.
func ipfilterParseSingle(config *IPFConfig, c *caddy.Controller) (IPPath, error) {
	var cPath IPPath
	ruleTypeSpecified := false

	// Get PathScopes
	cPath.PathScopes = c.RemainingArgs()
	if len(cPath.PathScopes) == 0 {
		return cPath, c.ArgErr()
	}

	// Sort PathScopes by length (the longest is always the most specific so should be tested first)
	sort.Sort(sort.Reverse(ByLength(cPath.PathScopes)))

	for c.NextBlock() {
		value := c.Val()

		switch value {
		case "rule":
			if !c.NextArg() {
				return cPath, c.ArgErr()
			}
			if ruleTypeSpecified {
				return cPath, c.Err("ipfilter: Only one 'rule' directive per block allowed")
			}

			rule := c.Val()
			if rule == "block" {
				cPath.IsBlock = true
			} else if rule != "allow" {
				return cPath, c.Err("ipfilter: Rule should be 'block' or 'allow'")
			}
			ruleTypeSpecified = true
		case "database":
			if !c.NextArg() {
				return cPath, c.ArgErr()
			}
			// Check if a database has already been opened
			if config.DBHandler != nil {
				return cPath, c.Err("ipfilter: A database is already opened")
			}

			database := c.Val()

			// Open the database.
			var err error
			config.DBHandler, err = maxminddb.Open(database)
			if err != nil {
				return cPath, c.Err("ipfilter: Can't open database: " + database)
			}
		case "blockpage":
			if !c.NextArg() {
				return cPath, c.ArgErr()
			}

			// check if blockpage exists.
			blockpage := c.Val()
			if _, err := os.Stat(blockpage); os.IsNotExist(err) {
				return cPath, c.Err("ipfilter: No such file: " + blockpage)
			}
			cPath.BlockPage = blockpage
		case "country":
			countryCodes := c.RemainingArgs()
			if len(countryCodes) == 0 {
				return cPath, c.ArgErr()
			}
			cPath.CountryCodes = append(cPath.CountryCodes, countryCodes...)
		case "ip":
			ips := c.RemainingArgs()
			if len(ips) == 0 {
				return cPath, c.ArgErr()
			}

			for _, ip := range ips {
				ipRange, err := parseIP(ip)
				if err != nil {
					return cPath, c.Err("ipfilter: " + err.Error())
				}

				cPath.Nets = append(cPath.Nets, ipRange...)
			}
		case "strict":
			if c.NextArg() {
				return cPath, c.ArgErr()
			}
			cPath.Strict = true
		case "prefix_dir":
			if !c.NextArg() || cPath.PrefixDir != "" {
				return cPath, c.ArgErr()
			}
			// Verify the IP address path prefix exists and is a directory.
			prefixDir := c.Val()
			if statb, err := os.Stat(prefixDir); os.IsNotExist(err) || !statb.IsDir() {
				return cPath, c.Err("ipfilter: No such blacklist prefix dir: " + prefixDir)
			}
			cPath.PrefixDir = prefixDir
		}
	}

	if !ruleTypeSpecified {
		return cPath, c.Err("ipfilter: There must be one 'rule' directive per block")
	}
	return cPath, nil
}

// ipfilterParse parses all ipfilter {} blocks to an IPFConfig
func ipfilterParse(c *caddy.Controller) (IPFConfig, error) {
	var config IPFConfig

	var hasCountryCodes, hasRanges, hasPrefixDir bool

	for c.Next() {
		path, err := ipfilterParseSingle(&config, c)
		if err != nil {
			return config, err
		}

		if len(path.CountryCodes) != 0 {
			hasCountryCodes = true
		}
		if len(path.Nets) != 0 {
			hasRanges = true
		}
		if path.PrefixDir != "" {
			hasPrefixDir = true
		}

		config.Paths = append(config.Paths, path)
	}

	// having a database is mandatory if you are blocking by country codes.
	if hasCountryCodes && config.DBHandler == nil {
		return config, c.Err("ipfilter: Database is required to block/allow by country")
	}

	// Must specify at least one of these subdirectives.
	if !hasCountryCodes && !hasRanges && !hasPrefixDir {
		return config, c.Err("ipfilter: No IPs, Country codes, or prefix dir has been provided")
	}

	return config, nil
}

// ByLength sorts strings by length and alphabetically (if same length)
type ByLength []string

func (s ByLength) Len() int      { return len(s) }
func (s ByLength) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

func (s ByLength) Less(i, j int) bool {
	if len(s[i]) < len(s[j]) {
		return true
	} else if len(s[i]) == len(s[j]) {
		return s[i] < s[j] // Compare alphabetically in ascending order
	}
	return false
}
