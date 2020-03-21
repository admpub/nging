package convert

import (
	"bytes"
	"io"
)

type Convert func(r io.Reader, quality int)(*bytes.Buffer, error)

var formatConverters = map[string]Convert {}

func Register(extension string, convert Convert) {
	formatConverters[extension] = convert
}

func Unregister(extension string) {
	if _, ok := formatConverters[extension]; ok {
		delete(formatConverters, extension)
	}
}

func Extensions() []string {
	extensions := make([]string, 0)
	for extension := range formatConverters {
		extensions = append(extensions, extension)
	}
	return extensions
}
