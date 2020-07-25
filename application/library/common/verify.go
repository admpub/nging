package common

import "github.com/webx-top/com"

var (
	boolFlags = []string{`Y`, `N`}
	contypes  = []string{`html`, `markdown`, `text`}
)

func GetBoolFlag(value string, defaults ...string) string {
	if len(value) == 0 || !com.InSlice(value, boolFlags) {
		if len(defaults) > 0 {
			return defaults[0]
		}
		return `N`
	}
	return value
}

func GetContype(value string, defaults ...string) string {
	if len(value) == 0 || !com.InSlice(value, contypes) {
		if len(defaults) > 0 {
			return defaults[0]
		}
		return `text`
	}
	return value
}

func GetEnumValue(enums []string, value string, defaults string) string {
	if len(value) == 0 || !com.InSlice(value, enums) {
		return defaults
	}
	return value
}
