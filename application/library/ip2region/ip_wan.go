package ip2region

import (
	"errors"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/admpub/nging/v5/application/library/common"
	"github.com/admpub/nging/v5/application/library/restclient"
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
	IP        string
	QueryTime time.Time
}

func GetWANIP(cachedSeconds float64, ipVers ...int) (wanIP WANIP, err error) {
	var (
		ipv4  string
		ipv6  string
		ipVer = 4
	)
	if len(ipVers) > 0 && ipVers[0] == 6 {
		ipVer = ipVers[0]
	}
	cacheFile := `v` + strconv.Itoa(ipVer)
	if cachedSeconds > 0 {
		var valid bool
		if m, e := common.ModTimeCache(`ip`, cacheFile); e == nil {
			wanIP.QueryTime = m
			if time.Since(m).Seconds() < cachedSeconds { // 缓存1小时(3600秒)
				valid = true
			}
		}
		if valid {
			if b, e := common.ReadCache(`ip`, cacheFile); e == nil {
				wanIP.IP = strings.TrimSpace(string(b))
				return
			}
		}
	}
	var errs []string
	for _, provider := range wanIPProviders {
		if provider == nil || provider.Disabled || len(provider.URL) == 0 {
			continue
		}
		if ipVer == 4 {
			if !provider.HasIPv4Rule() {
				continue
			}
		} else {
			if !provider.HasIPv6Rule() {
				continue
			}
		}
		client := restclient.Resty()
		resp, err := client.Execute(provider.Method, provider.URL)
		if err != nil {
			errs = append(errs, `[`+provider.Name+`] `+err.Error())
			continue
		}
		if !resp.IsSuccess() {
			if resp.StatusCode() == http.StatusNotFound {
				provider.Disabled = true
			}
			errs = append(errs, `[`+provider.Name+`] `+strconv.Itoa(resp.StatusCode())+`: `+resp.Status())
			continue
		}
		body := resp.Body()
		if len(body) == 0 {
			continue
		}
		if provider.ip6regexp != nil {
			if ipVer == 6 {
				matches := provider.ip6regexp.FindAllStringSubmatch(string(body), 1)
				if len(matches) > 0 && len(matches[0]) > 1 {
					ipv6 = strings.TrimSpace(matches[0][1])
					if len(ipv6) > 0 {
						if ipv6 = FindIPv6(ipv6); len(ipv6) > 0 {
							break
						}
					}
				}
			}
		} else if provider.IP6Rule == `=` {
			if ipVer == 6 {
				ipv6 = strings.TrimSpace(string(body))
				if len(ipv6) > 0 {
					if ipv6 = FindIPv6(ipv6); len(ipv6) > 0 {
						break
					}
				}
			}
			continue // 返回内容是IPv6，则没有必要再找IPv4了
		}
		if ipVer == 4 {
			if provider.ip4regexp != nil {
				matches := provider.ip4regexp.FindAllStringSubmatch(string(body), 1)
				//com.Dump(matches)
				if len(matches) > 0 && len(matches[0]) > 1 {
					ipv4 = strings.TrimSpace(matches[0][1])
				}
			} else {
				ipv4 = strings.TrimSpace(string(body))
			}
			if len(ipv4) > 0 {
				if ipv4 = FindIPv4(ipv4); len(ipv4) > 0 {
					break
				}
			}
		}
	}
	wanIP.QueryTime = time.Now()
	if ipVer == 4 {
		wanIP.IP = ipv4
	} else {
		wanIP.IP = ipv6
	}
	if len(wanIP.IP) > 0 {
		if err := common.WriteCache(`ip`, cacheFile, []byte(wanIP.IP)); err != nil {
			errs = append(errs, err.Error())
		}
	}
	if len(errs) > 0 {
		err = errors.New(strings.Join(errs, "\n"))
	}
	return
}
