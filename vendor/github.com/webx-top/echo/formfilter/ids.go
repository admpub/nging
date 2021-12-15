package formfilter

import (
	"strings"
)

func JoinValues(field string) Options {
	return func() (string, Filter) {
		return field, func(data *Data) {
			data.Value = []string{strings.Join(data.Value, `,`)}
		}
	}
}

func SplitValues(field string, seperators ...string) Options {
	seperator := `,`
	if len(seperators) > 0 {
		seperator = seperators[0]
	}
	return func() (string, Filter) {
		return field, func(data *Data) {
			if len(data.Value) > 0 && len(data.Value[0]) > 0 {
				data.Value = strings.Split(data.Value[0], seperator)
			}
		}
	}
}
