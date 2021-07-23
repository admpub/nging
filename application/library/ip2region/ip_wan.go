package ip2region

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/admpub/nging/v3/application/library/common"
	"github.com/admpub/nging/v3/application/library/config/extend"
	"github.com/admpub/nging/v3/application/library/restclient"
)

type WANIPProvider struct {
	Name        string
	Description string
	URL         string
	Method      string
	IP4Rule     string
	IP6Rule     string
	ip4regexp   *regexp.Regexp
	ip6regexp   *regexp.Regexp
	Disabled    bool
}

type WANIPProviders map[string]*WANIPProvider

func (w *WANIPProviders) Reload() error {
	if w == nil {
		return nil
	}
	for key, value := range *w {
		if value != nil && len(value.Name) > 0 && len(value.URL) > 0 {
			Register(value)
		} else {
			Unregister(key)
		}
	}
	return nil
}

var (
	wanIPProviders   = map[string]*WANIPProvider{}
	defaultProviders = []*WANIPProvider{
		{
			Name:        `sohu`,
			Description: `搜狐`,
			URL:         `https://pv.sohu.com/cityjson`,
			IP4Rule:     `"` + IPv4Rule + `"`,
		}, {
			Name:        `ip-api.com`,
			Description: ``,
			URL:         `http://ip-api.com/json/?fields=query`,
			IP4Rule:     `"query":"` + IPv4Rule + `"`,
		}, {
			Name:        `ip.sb`,
			Description: ``,
			URL:         `https://api.ip.sb/ip`,
			IP4Rule:     ``,
		}, {
			Name:        `ipconfig.io`,
			Description: ``,
			URL:         `https://ipconfig.io/ip`,
			IP4Rule:     ``,
		},
	}
)

const (
	IPv4Rule = `([0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3})`
)

func init() {
	extend.Register(`wanIPProvider`, func() interface{} {
		return &WANIPProviders{}
	})
	for _, provider := range defaultProviders {
		if err := Register(provider); err != nil {
			panic(err)
		}
	}
}

func Register(p *WANIPProvider) (err error) {
	if len(p.IP4Rule) > 0 && p.IP4Rule != `=` {
		p.ip4regexp, err = regexp.Compile(p.IP4Rule)
		if err != nil {
			return
		}
	}
	if len(p.IP6Rule) > 0 && p.IP6Rule != `=` {
		p.ip6regexp, err = regexp.Compile(p.IP6Rule)
		if err != nil {
			return
		}
	}
	if len(p.Method) == 0 {
		p.Method = `GET`
	}
	wanIPProviders[p.Name] = p
	return
}

func Get(name string) *WANIPProvider {
	p, _ := wanIPProviders[name]
	return p
}

func Unregister(names ...string) {
	for _, name := range names {
		delete(wanIPProviders, name)
	}
}

type WANIP struct {
	IPv4      string
	IPv6      string
	QueryTime time.Time
}

func GetWANIP(noCaches ...bool) (wanIP WANIP, err error) {
	var noCache bool
	if len(noCaches) > 0 {
		noCache = noCaches[0]
	}
	var (
		ipv4 string
		ipv6 string
	)
	if !noCache {
		var valid bool
		if m, e := common.ModTimeCache(`ip`, `wan`); e == nil {
			wanIP.QueryTime = m
			if time.Since(m).Seconds() < 3600 { // 缓存1小时(3600秒)
				valid = true
			}
		}
		if valid {
			if b, e := common.ReadCache(`ip`, `wan`); e == nil {
				c := strings.Split(string(b), "\n")
				if len(c) > 0 {
					ipv4 = strings.TrimSpace(c[0])
				}
				if len(c) > 1 {
					ipv6 = strings.TrimSpace(c[1])
				}
				wanIP.IPv4 = ipv4
				wanIP.IPv6 = ipv6
				return
			}
		}
	}
	var errs []string
	for _, provider := range wanIPProviders {
		if provider == nil || provider.Disabled || len(provider.URL) == 0 {
			continue
		}
		client := restclient.Resty()
		resp, err := client.Execute(provider.Method, provider.URL)
		if err != nil {
			errs = append(errs, `[`+provider.Name+`] `+err.Error())
			continue
		}
		if !resp.IsSuccess() {
			errs = append(errs, `[`+provider.Name+`] `+strconv.Itoa(resp.StatusCode())+`: `+resp.Status())
			continue
		}
		body := resp.Body()
		if len(body) == 0 {
			continue
		}
		if provider.ip6regexp != nil {
			matches := provider.ip6regexp.FindAllStringSubmatch(string(body), 1)
			if len(matches) > 0 && len(matches[0]) > 1 {
				ipv6 = matches[0][1]
			}
		} else if provider.IP6Rule == `=` {
			ipv6 = string(body)
			continue
		}
		if provider.ip4regexp != nil {
			matches := provider.ip4regexp.FindAllStringSubmatch(string(body), 1)
			//com.Dump(matches)
			if len(matches) > 0 && len(matches[0]) > 1 {
				ipv4 = matches[0][1]
			}
		} else {
			ipv4 = string(body)
		}
		break
	}
	if len(ipv4) > 0 || len(ipv6) > 0 {
		if err := common.WriteCache(`ip`, `wan`, []byte(ipv4+"\n"+ipv6)); err != nil {
			errs = append(errs, err.Error())
		}
	}
	if len(errs) > 0 {
		err = errors.New(strings.Join(errs, "\n"))
	}
	wanIP.QueryTime = time.Now()
	wanIP.IPv4 = ipv4
	wanIP.IPv6 = ipv6
	return
}
