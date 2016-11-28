package echo

import "fmt"

type Translator interface {
	T(format string, args ...interface{}) string
	Lang() string
}

var DefaultNopTranslate Translator = &NopTranslate{language: `en`}

type NopTranslate struct {
	language string
}

func (n *NopTranslate) T(format string, args ...interface{}) string {
	if len(args) > 0 {
		return fmt.Sprintf(format, args...)
	}
	return format
}

func (n *NopTranslate) Lang() string {
	return n.language
}
