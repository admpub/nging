package common

import "github.com/webx-top/com"

const (
	BoolY = `Y`
	BoolN = `N`
)

const (
	ContentTypeHTML     = `html`
	ContentTypeMarkdown = `markdown`
	ContentTypeText     = `text`
)

var (
	boolFlags = []string{BoolY, BoolN}
	contypes  = []string{ContentTypeHTML, ContentTypeMarkdown, ContentTypeText}
)

func GetBoolFlag(value string, defaults ...string) string {
	if len(value) == 0 || !com.InSlice(value, boolFlags) {
		if len(defaults) > 0 {
			return defaults[0]
		}
		return BoolN
	}
	return value
}

func IsBoolFlag(value string) bool {
	return com.InSlice(value, boolFlags)
}

func BoolToFlag(v bool) string {
	if v {
		return BoolY
	}
	return BoolN
}

func FlagToBool(v string) bool {
	return v == BoolY
}

func GetContype(value string, defaults ...string) string {
	if len(value) == 0 || !com.InSlice(value, contypes) {
		if len(defaults) > 0 {
			return defaults[0]
		}
		return ContentTypeText
	}
	return value
}

func IsContype(value string) bool {
	return com.InSlice(value, contypes)
}

func GetEnumValue(enums []string, value string, defaults string) string {
	if len(value) == 0 || !com.InSlice(value, enums) {
		return defaults
	}
	return value
}
