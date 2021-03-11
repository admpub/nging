package convert

import (
	"bytes"
	"io"
)

// Convert 图片格式转换
type Convert func(r io.Reader, quality int) (*bytes.Buffer, error)

var formatConverters = map[string]Convert{}

// Register 注册图片格式转换功能
func Register(extension string, convert Convert) {
	formatConverters[extension] = convert
}

// GetConverter 获取图片格式转换功能
func GetConverter(extension string) (Convert, bool) {
	convert, ok := formatConverters[extension]
	return convert, ok
}

// Unregister 取消注册
func Unregister(extension string) {
	if _, ok := formatConverters[extension]; ok {
		delete(formatConverters, extension)
	}
}

// Extensions 所有可转换格式的扩展名
func Extensions() []string {
	extensions := make([]string, 0)
	for extension := range formatConverters {
		extensions = append(extensions, extension)
	}
	return extensions
}
