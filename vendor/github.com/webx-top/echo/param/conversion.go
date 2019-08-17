package param

import (
	"fmt"
	"html/template"
	"strconv"
	"strings"
	"time"
)

const (
	EmptyHTML      = template.HTML(``)
	EmptyJS        = template.JS(``)
	EmptyCSS       = template.CSS(``)
	EmptyHTMLAttr  = template.HTMLAttr(``)
	DateTimeLayout = `2006-01-02 15:04:05`
	DateTimeShort  = `2006-01-02 15:04`
	DateLayout     = `2006-01-02`
	TimeLayout     = `15:04:05`
	DateMd         = `01-02`
	DateShort      = `06-01-02`
	TimeShort      = `15:04`
)

func AsString(val interface{}) string {
	switch v := val.(type) {
	case string:
		return v
	case nil:
		return ``
	default:
		return fmt.Sprint(val)
	}
}

func Split(val interface{}, sep string, limit ...int) StringSlice {
	str := AsString(val)
	if len(str) == 0 {
		return StringSlice{}
	}
	if len(limit) > 0 {
		return strings.SplitN(str, sep, limit[0])
	}
	return strings.Split(str, sep)
}

func Trim(val interface{}) String {
	return String(strings.TrimSpace(AsString(val)))
}

func AsHTML(val interface{}) template.HTML {
	switch v := val.(type) {
	case template.HTML:
		return v
	case string:
		return template.HTML(v)
	case nil:
		return EmptyHTML
	default:
		return template.HTML(fmt.Sprint(v))
	}
}

func AsHTMLAttr(val interface{}) template.HTMLAttr {
	switch v := val.(type) {
	case template.HTMLAttr:
		return v
	case string:
		return template.HTMLAttr(v)
	case nil:
		return EmptyHTMLAttr
	default:
		return template.HTMLAttr(fmt.Sprint(v))
	}
}

func AsJS(val interface{}) template.JS {
	switch v := val.(type) {
	case template.JS:
		return v
	case string:
		return template.JS(v)
	case nil:
		return EmptyJS
	default:
		return template.JS(fmt.Sprint(v))
	}
}

func AsCSS(val interface{}) template.CSS {
	switch v := val.(type) {
	case template.CSS:
		return v
	case string:
		return template.CSS(v)
	case nil:
		return EmptyCSS
	default:
		return template.CSS(fmt.Sprint(v))
	}
}

func AsBool(val interface{}) bool {
	switch v := val.(type) {
	case bool:
		return v
	case string:
		if len(v) > 0 {
			r, _ := strconv.ParseBool(v)
			return r
		}
		return false
	case nil:
		return false
	default:
		p := fmt.Sprint(v)
		if len(p) > 0 {
			r, _ := strconv.ParseBool(p)
			return r
		}
	}
	return false
}

func AsFloat64(val interface{}) float64 {
	switch v := val.(type) {
	case float64:
		return v
	case int64:
		return float64(v)
	case uint64:
		return float64(v)
	case float32:
		return float64(v)
	case int32:
		return float64(v)
	case uint32:
		return float64(v)
	case int:
		return float64(v)
	case uint:
		return float64(v)
	case string:
		i, _ := strconv.ParseFloat(v, 64)
		return i
	case nil:
		return 0
	default:
		s := fmt.Sprint(v)
		i, _ := strconv.ParseFloat(s, 64)
		return i
	}
}

func AsFloat32(val interface{}) float32 {
	switch v := val.(type) {
	case float32:
		return v
	case int32:
		return float32(v)
	case uint32:
		return float32(v)
	case string:
		f, _ := strconv.ParseFloat(v, 32)
		return float32(f)
	case nil:
		return 0
	default:
		s := fmt.Sprint(val)
		f, _ := strconv.ParseFloat(s, 32)
		return float32(f)
	}
}

func AsInt8(val interface{}) int8 {
	switch v := val.(type) {
	case int8:
		return v
	case string:
		i, _ := strconv.ParseInt(v, 10, 8)
		return int8(i)
	case nil:
		return 0
	default:
		s := fmt.Sprint(val)
		i, _ := strconv.ParseInt(s, 10, 8)
		return int8(i)
	}
}

func AsInt16(val interface{}) int16 {
	switch v := val.(type) {
	case int16:
		return v
	case string:
		i, _ := strconv.ParseInt(v, 10, 16)
		return int16(i)
	case nil:
		return 0
	default:
		s := fmt.Sprint(v)
		i, _ := strconv.ParseInt(s, 10, 16)
		return int16(i)
	}
}

func AsInt(val interface{}) int {
	switch v := val.(type) {
	case int:
		return v
	case string:
		i, _ := strconv.Atoi(v)
		return i
	case nil:
		return 0
	default:
		s := fmt.Sprint(v)
		i, _ := strconv.Atoi(s)
		return i
	}
}

func AsInt32(val interface{}) int32 {
	switch v := val.(type) {
	case int32:
		return v
	case string:
		i, _ := strconv.ParseInt(v, 10, 32)
		return int32(i)
	case nil:
		return 0
	default:
		s := fmt.Sprint(v)
		i, _ := strconv.ParseInt(s, 10, 32)
		return int32(i)
	}
}

func AsInt64(val interface{}) int64 {
	switch v := val.(type) {
	case int64:
		return v
	case int32:
		return int64(v)
	case uint32:
		return int64(v)
	case int:
		return int64(v)
	case uint:
		return int64(v)
	case string:
		i, _ := strconv.ParseInt(v, 10, 64)
		return i
	case nil:
		return 0
	default:
		s := fmt.Sprint(v)
		i, _ := strconv.ParseInt(s, 10, 64)
		return i
	}
}

func Decr(val interface{}, n int64) int64 {
	v, _ := val.(int64)
	v -= n
	return v
}

func Incr(val interface{}, n int64) int64 {
	v, _ := val.(int64)
	v += n
	return v
}

func AsUint8(val interface{}) uint8 {
	switch v := val.(type) {
	case uint8:
		return v
	case string:
		i, _ := strconv.ParseUint(v, 10, 8)
		return uint8(i)
	case nil:
		return 0
	default:
		s := fmt.Sprint(v)
		i, _ := strconv.ParseUint(s, 10, 8)
		return uint8(i)
	}
}

func AsUint16(val interface{}) uint16 {
	switch v := val.(type) {
	case uint16:
		return v
	case string:
		i, _ := strconv.ParseUint(v, 10, 16)
		return uint16(i)
	case nil:
		return 0
	default:
		s := fmt.Sprint(v)
		i, _ := strconv.ParseUint(s, 10, 16)
		return uint16(i)
	}
}

func AsUint(val interface{}) uint {
	switch v := val.(type) {
	case uint:
		return v
	case string:
		i, _ := strconv.ParseUint(v, 10, 32)
		return uint(i)
	case nil:
		return 0
	default:
		s := fmt.Sprint(v)
		i, _ := strconv.ParseUint(s, 10, 32)
		return uint(i)
	}
}

func AsUint32(val interface{}) uint32 {
	switch v := val.(type) {
	case uint32:
		return v
	case string:
		i, _ := strconv.ParseUint(v, 10, 32)
		return uint32(i)
	case nil:
		return 0
	default:
		s := fmt.Sprint(v)
		i, _ := strconv.ParseUint(s, 10, 32)
		return uint32(i)
	}
}

func AsUint64(val interface{}) uint64 {
	switch v := val.(type) {
	case uint64:
		return v
	case string:
		i, _ := strconv.ParseUint(v, 10, 64)
		return i
	case nil:
		return 0
	default:
		s := fmt.Sprint(v)
		i, _ := strconv.ParseUint(s, 10, 64)
		return i
	}
}

func AsTimestamp(val interface{}) time.Time {
	p := AsString(val)
	if len(p) > 0 {
		s := strings.SplitN(p, `.`, 2)
		var sec int64
		var nsec int64
		switch len(s) {
		case 2:
			nsec = String(s[1]).Int64()
			fallthrough
		case 1:
			sec = String(s[0]).Int64()
		}
		return time.Unix(sec, nsec)
	}
	return emptyTime
}

func AsDateTime(val interface{}, layouts ...string) time.Time {
	p := AsString(val)
	if len(p) > 0 {
		layout := DateTimeLayout
		if len(layouts) > 0 {
			layout = layouts[0]
		}
		t, _ := time.Parse(layout, p)
		return t
	}
	return emptyTime
}
