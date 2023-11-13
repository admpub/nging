package echo

import (
	"fmt"
	"strings"
	"time"

	"github.com/admpub/log"
	"github.com/webx-top/com"
	"github.com/webx-top/echo/param"
)

type (
	//FormDataFilter 将map映射到结构体时，对名称和值的过滤处理，如果返回的名称为空，则跳过本字段
	//这里的 key 为结构体字段或map的key层级路径，如果有父级则用星号“*”表示，例如：*.Name
	FormDataFilter func(key string, values []string) (string, []string)
)

var (
	//DefaultNopFilter 默认过滤器(map->struct)
	DefaultNopFilter FormDataFilter = func(k string, v []string) (string, []string) {
		return k, v
	}
	//DateToTimestamp 日期时间转时间戳
	DateToTimestamp = func(layouts ...string) FormDataFilter {
		layout := `2006-01-02`
		if len(layouts) > 0 && len(layouts[0]) > 0 {
			layout = layouts[0]
		}
		return func(k string, v []string) (string, []string) {
			if len(v) > 0 && len(v[0]) > 0 {
				t, e := time.ParseInLocation(layout, v[0], time.Local)
				if e != nil {
					log.Error(e)
					return k, []string{`0`}
				}
				return k, []string{fmt.Sprint(t.Unix())}
			}
			return k, []string{`0`}
		}
	}
	//TimestampToDate 时间戳转日期时间
	TimestampToDate = func(layouts ...string) FormDataFilter {
		layout := `2006-01-02 15:04:05`
		if len(layouts) > 0 && len(layouts[0]) > 0 {
			layout = layouts[0]
		}
		return func(k string, v []string) (string, []string) {
			if len(v) > 0 && len(v[0]) > 0 {
				tsi := strings.SplitN(v[0], `.`, 2)
				var sec, nsec int64
				switch len(tsi) {
				case 2:
					nsec = param.AsInt64(tsi[1])
					fallthrough
				case 1:
					sec = param.AsInt64(tsi[0])
				}
				t := time.Unix(sec, nsec)
				if t.IsZero() {
					return k, []string{``}
				}
				return k, []string{t.Format(layout)}
			}
			return k, v
		}
	}
	//JoinValues 组合数组为字符串
	JoinValues = func(seperators ...string) FormDataFilter {
		sep := `,`
		if len(seperators) > 0 {
			sep = seperators[0]
		}
		return func(k string, v []string) (string, []string) {
			return k, []string{strings.Join(v, sep)}
		}
	}
	//SplitValues 拆分字符串为数组
	SplitValues = func(seperators ...string) FormDataFilter {
		sep := `,`
		if len(seperators) > 0 {
			sep = seperators[0]
		}
		return func(k string, v []string) (string, []string) {
			if len(v) > 0 && len(v[0]) > 0 {
				v = strings.Split(v[0], sep)
			}
			return k, v
		}
	}
)

// FormatFieldValue 格式化字段值
func FormatFieldValue(formatters map[string]FormDataFilter, keyNormalizerArg ...func(string) string) FormDataFilter {
	newFormatters := map[string]FormDataFilter{}
	keyNormalizer := strings.Title
	if len(keyNormalizerArg) > 0 && keyNormalizerArg[0] != nil {
		keyNormalizer = keyNormalizerArg[0]
	}
	for k, v := range formatters {
		newFormatters[keyNormalizer(k)] = v
	}
	return func(k string, v []string) (string, []string) {
		tk := keyNormalizer(k)
		if formatter, ok := newFormatters[tk]; ok {
			return formatter(k, v)
		}
		return k, v
	}
}

// IncludeFieldName 包含字段
func IncludeFieldName(fieldNames ...string) FormDataFilter {
	for k, v := range fieldNames {
		fieldNames[k] = com.Title(v)
	}
	return func(k string, v []string) (string, []string) {
		tk := com.Title(k)
		for _, fv := range fieldNames {
			if fv == tk {
				return k, v
			}
		}
		return ``, v
	}
}

// ExcludeFieldName 排除字段
func ExcludeFieldName(fieldNames ...string) FormDataFilter {
	for k, v := range fieldNames {
		fieldNames[k] = com.Title(v)
	}
	return func(k string, v []string) (string, []string) {
		tk := com.Title(k)
		for _, fv := range fieldNames {
			if fv == tk {
				return ``, v
			}
		}
		return k, v
	}
}
