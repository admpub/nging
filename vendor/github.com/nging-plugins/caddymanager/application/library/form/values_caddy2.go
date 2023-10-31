package form

import (
	"strings"
	"time"

	"github.com/nging-plugins/caddymanager/application/library/htpasswd"
)

// {remote} - {user} [{when}] "{method} {scheme} {host} {uri} {proto}" {status} {size} "{>Referer}" "{>User-Agent}" {latency}
var Caddy2LogFormatReplacer = strings.NewReplacer(
	`{method} {uri} {proto}`, `{request>method} {request>uri} {request>proto}`,
	`{remote}`, `{request>remote_ip}`,
	`{user}`, `{request>user_id}`,
	`{when}`, `{ts}`,
	`{method}`, `{request>method}`,
	//`{scheme}`, `{scheme}`, TODO
	//`{host}`, `{host}`, TODO
	`{uri}`, `{request>uri}`,
	`{proto}`, `{request>proto}`,
	//`{status}`, `{status}`, 相同
	//`{size}`, `{size}`, 相同
	`{>Referer}`, `{request>headers>Referer>[0]}`,
	`{>User-Agent}`, `{request>headers>User-Agent>[0]}`,
	`{latency}`, `{duration}`, //TODO
)

func AsCaddy2LogFormat(value string) string {
	value = ExplodeCombinedLogFormat(value)
	return Caddy2LogFormatReplacer.Replace(value)
}

func (v Values) AsDuration(str string) time.Duration {
	return ParseDuration(str)
}

func (v Values) AsCaddy2LogFormat() string {
	return AsCaddy2LogFormat(v.Get(`log_format`))
}

// 添加路径通配符后缀
func (v Values) AddPathWildcardSuffix(path string) string {
	if strings.HasSuffix(path, `*`) {
		return path
	}
	if !strings.HasSuffix(path, `/`) {
		return path
	}
	return path + `*`
}

func (v Values) BCrypt(password string) string {
	password, _ = htpasswd.HashBCrypt(password)
	return password
}

// 添加扩展名通配符前缀
func (v Values) AddExtWildcardPrefix(path string) string {
	if strings.HasPrefix(path, `*`) {
		return path
	}
	if !strings.HasPrefix(path, `.`) {
		path = `.` + path
	}
	return `*` + path
}

func (v Values) AddonAttrFullKey(fullKey string, item string, defaults ...string) string {
	val := v.Values.Get(fullKey)
	if len(val) == 0 {
		if len(defaults) > 0 && len(defaults[0]) > 0 {
			return item + `   ` + defaults[0]
		}
		return ``
	}
	return item + `   ` + val
}
