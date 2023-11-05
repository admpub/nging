package gopiper

import (
	"strings"

	"github.com/webx-top/com"
)

func _filterValue(src interface{}, fn func(v string) (interface{}, error), fnDefaults ...func(interface{}) (interface{}, error)) (interface{}, error) {

	switch vt := src.(type) {
	case string:
		return fn(vt)

	case []string:
		for i, v := range vt {
			_v, _e := fn(v)
			if _e != nil {
				vt[i] = _e.Error()
				continue
			}
			vt[i], _ = _v.(string)
		}
		return vt, nil

	case map[string]string:
		for i, v := range vt {
			_v, _e := fn(v)
			if _e != nil {
				vt[i] = _e.Error()
				continue
			}
			vt[i], _ = _v.(string)
		}
		return vt, nil

	case []interface{}:
		for i, v := range vt {
			vt[i], _ = _filterValue(v, fn)
		}
		return vt, nil

	case map[string]interface{}:
		for i, v := range vt {
			vt[i], _ = _filterValue(v, fn)
		}
		return vt, nil

	case bool, float32, float64, int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		if len(fnDefaults) > 0 && fnDefaults[0] != nil {
			return fnDefaults[0](vt)
		}
		return fn(com.ToStr(vt))

	default:
		if len(fnDefaults) > 0 && fnDefaults[0] != nil {
			return fnDefaults[0](vt)
		}
		return vt, nil
	}
}

func SplitParams(params string, separators ...string) []string {
	if len(params) == 0 {
		return []string{}
	}
	separator := `,`
	if len(separators) > 0 {
		separator = separators[0]
		if len(separator) < 1 {
			return strings.Split(params, separator)
		}
		if len(separator) > 1 {
			separator = separator[0:1]
		}
	}
	vt := strings.Split(params, separator)
	var (
		lastEnd string
		results []string
	)
	for k, v := range vt {
		lastKey := k - 1
		if lastEnd == `\` {
			lastVal := vt[lastKey]
			vt[lastKey] = lastVal[0:len(lastVal)-1] + separator + v
			resultLen := len(results)
			if resultLen > 0 {
				results[resultLen-1] = vt[lastKey]
			}
			lastEnd = v[len(v)-1:]
			continue
		}
		lastEnd = v[len(v)-1:]
		results = append(results, v)
	}
	return results
}
