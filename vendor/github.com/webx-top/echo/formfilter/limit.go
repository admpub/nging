package formfilter

import "strings"

//Format 格式化字段值
func Format(field string, formatter func([]string) []string) Options {
	return func() (string, Filter) {
		return field, func(data *Data) {
			data.Value = formatter(data.Value)
		}
	}
}

//Include 包含字段
func Include(fieldNames ...string) Options {
	for k, v := range fieldNames {
		fieldNames[k] = strings.Title(v)
	}
	return func() (string, Filter) {
		return All, func(data *Data) {
			for _, fv := range fieldNames {
				if fv == data.NormalizedKey() {
					return
				}
			}
			data.Key = ``
		}
	}
}

//Exclude 排除字段
func Exclude(fieldNames ...string) Options {
	for k, v := range fieldNames {
		fieldNames[k] = strings.Title(v)
	}
	return func() (string, Filter) {
		return All, func(data *Data) {
			for _, fv := range fieldNames {
				if fv == data.NormalizedKey() {
					data.Key = ``
					return
				}
			}
		}
	}
}
