package ip2region

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/admpub/nging/v3/application/library/common"
	"github.com/admpub/nging/v3/application/library/restclient"
)

const (
	IPv4Rule = `((?:(?:25[0-5]|(?:2[0-4]|1{0,1}[0-9]){0,1}[0-9])\.){3,3}(?:25[0-5]|(?:2[0-4]|1{0,1}[0-9]){0,1}[0-9]))`
	IPv6Rule = `((?:(?:(?:[0-9A-Fa-f]{1,4}:){7}(?:[0-9A-Fa-f]{1,4}|:))|(?:(?:[0-9A-Fa-f]{1,4}:){6}(?::[0-9A-Fa-f]{1,4}|(?:(?:25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(?:\.(25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3})|:))|(?:(?:[0-9A-Fa-f]{1,4}:){5}(?:(?:(?::[0-9A-Fa-f]{1,4}){1,2})|:(?:(?:25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(?:\.(?:25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3})|:))|(?:(?:[0-9A-Fa-f]{1,4}:){4}(?:(?:(?::[0-9A-Fa-f]{1,4}){1,3})|(?:(?::[0-9A-Fa-f]{1,4})?:(?:(?:25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(?:\.(?:25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3}))|:))|(?:(?:[0-9A-Fa-f]{1,4}:){3}(?:(?:(?::[0-9A-Fa-f]{1,4}){1,4})|(?:(?::[0-9A-Fa-f]{1,4}){0,2}:(?:(?:25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(?:\.(?:25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3}))|:))|(?:(?:[0-9A-Fa-f]{1,4}:){2}(?:(?:(?::[0-9A-Fa-f]{1,4}){1,5})|(?:(?::[0-9A-Fa-f]{1,4}){0,3}:(?:(?:25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(?:\.(?:25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3}))|:))|(?:(?:[0-9A-Fa-f]{1,4}:){1}(?:(?:(:[0-9A-Fa-f]{1,4}){1,6})|(?:(?::[0-9A-Fa-f]{1,4}){0,4}:(?:(?:25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(?:\.(?:25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3}))|:))|(?::(?:(?:(?::[0-9A-Fa-f]{1,4}){1,7})|(?:(?::[0-9A-Fa-f]{1,4}){0,5}:(?:(?:25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(?:\.(?:25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3}))|:))))`
)

var ipv6Regexp = regexp.MustCompile(IPv6Rule)
var ipv4Regexp = regexp.MustCompile(IPv4Rule)

func FindIPv4(content string) string {
	matches := ipv4Regexp.FindAllStringSubmatch(content, 1)
	if len(matches) > 0 && len(matches[0]) > 1 {
		return matches[0][1]
	}
	return ``
}

func FindIPv6(content string) string {
	matches := ipv6Regexp.FindAllStringSubmatch(content, 1)
	if len(matches) > 0 && len(matches[0]) > 1 {
		return matches[0][1]
	}
	return ``
}

type WANIP struct {
	IPv4      string
	IPv6      string
	QueryTime time.Time
}

func GetWANIP(cachedSeconds float64) (wanIP WANIP, err error) {
	var (
		ipv4 string
		ipv6 string
	)
	if cachedSeconds > 0 {
		var valid bool
		if m, e := common.ModTimeCache(`ip`, `wan`); e == nil {
			wanIP.QueryTime = m
			if time.Since(m).Seconds() < cachedSeconds { // 缓存1小时(3600秒)
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
