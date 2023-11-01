package form

import (
	"fmt"
	"html/template"
	"math"
	"mime"
	"path"
	"sort"
	"strconv"
	"strings"

	"github.com/webx-top/com"
	"github.com/webx-top/echo/param"
)

type NginxDomainInfo struct {
	Port      int
	Args      []string
	Domains   []string
	CertsPath []*CertPath
}

var allPathFieldsForNginx = []string{
	`header_path`,
	//`fastcgi_path`,
	`proxy_from`,
	`browse_path`,
	`expires_match_k[]`,
	`basicauth_resources[]`,
	`ratelimit_resources[]`,
	`cors_path`,
}

type LocationDef struct {
	PathKey  string
	Module   string
	Location string
	IsRegexp bool
	Items    []*Item
}

type Item struct {
	Key  interface{}
	Val  interface{}
	Args []interface{}
}

type Locations struct {
	SortedStaticPath []string
	SortedRegexpPath []string
	GroupByPath      map[string][]*LocationDef
}

type SortByLen []string

func (s SortByLen) Len() int { return len(s) }
func (s SortByLen) Less(i, j int) bool {
	return len(s[i]) < len(s[j])
}
func (s SortByLen) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

func (v Values) NginxLocations() Locations {
	return v.GroupByLocations(allPathFieldsForNginx)
}

// {remote} - {user} [{when}] "{method} {scheme} {host} {uri} {proto}" {status} {size} "{>Referer}" "{>User-Agent}" {latency}
var NginxLogFormatReplacer = strings.NewReplacer(
	//`{method} {scheme} {host} {uri} {proto}`, `$request`,
	`{method} {uri} {proto}`, `$request`,
	`{remote}`, `$remote_addr`,
	`{user}`, `$remote_user`,
	`{when}`, `$time_local`,
	`{method}`, `$request_method`,
	`{scheme}`, `$scheme`,
	`{host}`, `$http_host`,
	`{uri}`, `$request_uri`,
	`{proto}`, `$server_protocol`,
	`{status}`, `$status`,
	`{size}`, `$body_bytes_sent`,
	`{>Referer}`, `$http_referer`,
	`{>User-Agent}`, `$http_user_agent`,
	`{latency}`, `$request_time`,
)

func AsNginxLogFormat(value string) string {
	value = ExplodeCombinedLogFormat(value)
	return NginxLogFormatReplacer.Replace(value)
}

func (v Values) AsNginxLogFormat() string {
	return AsNginxLogFormat(v.Get(`log_format`))
}

func (v Values) NginxLimitRateWithUnit() string {
	rate := param.AsInt64(v.Get(`ratelimit_rate`))
	switch v.Get(`ratelimit_unit`) {
	case `second`:
		return fmt.Sprintf(`%dr/s`, rate)
	case `minute`:
		return fmt.Sprintf(`%dr/m`, rate)
	case `hour`:
		return fmt.Sprintf(`%vr/m`, math.Floor(float64(rate)/60))
	case `day`:
		return fmt.Sprintf(`%vr/m`, math.Floor((float64(rate)/24)/60))
	case `week`:
		return fmt.Sprintf(`%vr/m`, math.Floor(((float64(rate)/7)/24)/60))
	default:
		return fmt.Sprintf(`%dr/s`, rate)
	}
}

func (v Values) ExtensionsToMime(value string) []string {
	extensions := SplitBySpace(value)
	mimes := make([]string, 0, len(extensions))
	for _, ext := range extensions {
		// text/css; charset=utf-8
		mimeType := strings.SplitN(mime.TypeByExtension(ext), `;`, 2)[0]
		if len(mimeType) > 0 {
			mimes = append(mimes, mimeType)
		}
	}
	return mimes
}

type CertPath struct {
	Cert   string
	Key    string
	Trust  string
	Domain string
}

func (v Values) GetCerts(domains []string) []*CertPath {
	res := make([]*CertPath, 0, len(domains))
	for _, domain := range domains {
		cert := v.Get(`tls/` + domain + `/cert`)
		certKey := v.Get(`tls/` + domain + `/cert_key`)
		certTrust := v.Get(`tls/` + domain + `/cert_trust`)
		if len(cert) > 0 && len(certKey) > 0 {
			res = append(res, &CertPath{
				Cert:   cert,
				Key:    certKey,
				Trust:  certTrust,
				Domain: domain,
			})
		}
	}
	return res
}

func (v Values) GetNginxDomainList() []NginxDomainInfo {
	domainList := v.GetDomainList()
	var list []NginxDomainInfo
	portsDomains := map[int][]string{}
	var ports []int
	for _, domain := range domainList {
		domain = com.ParseEnvVar(domain)
		parts := strings.SplitN(domain, `://`, 2)
		var scheme, host, port string
		if len(parts) == 2 {
			scheme = parts[1]
			host, port = com.SplitHostPort(parts[1])
		} else {
			host, port = com.SplitHostPort(domain)
		}
		portNumber, _ := strconv.ParseUint(port, 10, 16)
		portN := int(portNumber)
		if portN == 0 {
			switch scheme {
			case `http`:
				portN = 80
			case `https`:
				portN = 443
			default:
				portN = 80
			}
		}
		if len(host) == 0 {
			host = `127.0.0.1`
		}
		if _, ok := portsDomains[portN]; !ok {
			portsDomains[portN] = []string{}
			ports = append(ports, portN)
		}
		portsDomains[portN] = append(portsDomains[portN], host)
	}
	sort.Sort(sort.IntSlice(ports))
	isTLS := v.Values.Get(`tls`) == `1`
	for _, portN := range ports {
		info := NginxDomainInfo{
			Port:    portN,
			Domains: portsDomains[portN],
		}
		if isTLS {
			info.CertsPath = v.GetCerts(info.Domains)
			if len(info.CertsPath) > 0 {
				info.Args = append(info.Args, `ssl`, `http2`)
			}
		}
		list = append(list, info)
	}
	return list
}

type UpstreamInfo struct {
	Scheme       string
	Host         string
	Path         string
	UpstreamName string
	Rewrite      string // rewrite ^/user/(.*)$ /$1 break;
	withQuote    bool
}

func (u UpstreamInfo) String() string {
	host := u.Host
	if len(u.UpstreamName) > 0 {
		host = u.UpstreamName
	}
	value := u.Scheme
	if u.Scheme == `unix` {
		value += `:/`
	} else {
		value += `://`
	}
	value += host + u.Path
	if u.withQuote {
		value = `"` + com.AddCSlashes(value, '"') + `"`
	}
	return value
}

func (v Values) ServerGroup(key string, customHost string, withQuotes ...bool) interface{} {
	var withQuote bool
	if len(withQuotes) > 0 {
		withQuote = withQuotes[0]
	}
	val := v.Get(key)
	sh := strings.SplitN(val, `://`, 2)
	var scheme string
	var host string
	var ppath string
	if len(sh) == 2 {
		scheme = sh[0]
		host = sh[1]
	} else {
		parts := strings.SplitN(val, `:`, 2)
		if len(parts) == 2 {
			scheme = parts[0]
			host = strings.TrimLeft(parts[1], "/")
		} else {
			scheme = `http`
			host = val
		}
	}
	hp := strings.SplitN(host, `/`, 2)
	if len(hp) == 2 {
		host = hp[0]
		ppath = `/` + hp[1]
	}
	var rewrite string
	if scheme == `http` || scheme == `https` {
		stripPrefix := v.Get(`proxy_without`)
		if len(stripPrefix) > 0 {
			proxyPath := v.Get(`proxy_from`)
			if stripPrefix == proxyPath {
				ppath = `/`
			} else {
				if !strings.HasPrefix(stripPrefix, `/`) {
					stripPrefix = `/` + stripPrefix
				}
				var targetPrefix string
				if len(ppath) > 0 && path.Join(proxyPath, ppath) != strings.TrimSuffix(stripPrefix, `/`) {
					targetPrefix = strings.TrimSuffix(ppath, `/`)
				}
				rewrite = fmt.Sprintf(
					`rewrite %q %s/$1 break;`,
					`^`+v.AddSlashes(stripPrefix)+`(.*)`,
					targetPrefix,
				)
			}
		}
	}
	return UpstreamInfo{
		Scheme:       scheme,
		Host:         host,
		Path:         ppath,
		UpstreamName: customHost,
		Rewrite:      rewrite,
		withQuote:    withQuote,
	}
}

func (v Values) IteratorHeaderKV(addon string, item string, plusPrefix string, minusPrefix string, withValueAndQuotes ...bool) interface{} {
	if len(addon) > 0 && len(item) > 0 {
		addon += `_`
	}
	k := addon + item + `_k`
	keys, _ := v.Values[k]

	k = addon + item + `_v`
	values, _ := v.Values[k]

	var r, t string
	var withValueAndQuote bool
	if len(withValueAndQuotes) > 0 {
		withValueAndQuote = withValueAndQuotes[0]
	}
	l := len(values)
	var suffix string
	if v.Config.GetEngine() == `nginx` {
		suffix = `;`
	}
	for i, k := range keys {
		if i < l {
			prefix := plusPrefix
			if strings.HasPrefix(k, `-`) {
				if len(minusPrefix) == 0 {
					continue
				}
				k = strings.TrimPrefix(k, `-`)
				prefix = minusPrefix
			} else {
				if len(plusPrefix) == 0 {
					continue
				}
				k = strings.TrimPrefix(k, `+`)
			}
			if withValueAndQuote {
				v := values[i]
				v = `"` + com.AddCSlashes(v, '"') + `"`
				r += t + prefix + k + `   ` + v + suffix
			} else {
				r += t + prefix + k + suffix
			}
			t = "\n"
		}
	}
	if withValueAndQuote {
		return template.HTML(r)
	}
	return r
}

func (v Values) IteratorNginxProxyHeaderKV() interface{} {
	addon := `proxy`
	item := `header_downstream`
	if len(addon) > 0 && len(item) > 0 {
		addon += `_`
	}
	k := addon + item + `_k`
	keys, _ := v.Values[k]

	k = addon + item + `_v`
	values, _ := v.Values[k]

	var r, t string
	l := len(values)
	suffix := `;`
	for i, k := range keys {
		if i < l {
			var prefix string
			if strings.HasPrefix(k, `-`) {
				k = strings.TrimPrefix(k, `-`)
				prefix = `proxy_hide_header`
				r += t + prefix + ` ` + k + suffix
			} else {
				k = strings.TrimPrefix(k, `+`)
				prefix = `proxy_pass_header`
				v := values[i]
				if len(v) > 0 {
					prefix = `add_header`
					v = `"` + com.AddCSlashes(v, '"') + `"`
					r += t + prefix + ` ` + k + ` ` + v + suffix
				} else {
					r += t + prefix + ` ` + k + suffix
				}
			}
			t = "\n"
		}
	}
	return template.HTML(r)
}
