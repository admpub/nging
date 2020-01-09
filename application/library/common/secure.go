package common

import (
	"bytes"
	"io"

	"github.com/microcosm-cc/bluemonday"
	"github.com/webx-top/com"
)

// ClearHTML 清除所有HTML标签及其属性，一般用处理文章标题等不含HTML标签的字符串
func ClearHTML(title string) string {
	p := bluemonday.StrictPolicy()
	return p.Sanitize(title)
}

// RemoveXSS 清除不安全的HTML标签和属性，一般用于处理文章内容
func RemoveXSS(content string) string {
	p := bluemonday.UGCPolicy()
	return p.Sanitize(content)
}

func RemoveBytesXSS(content []byte) []byte {
	p := bluemonday.UGCPolicy()
	return p.SanitizeBytes(content)
}

func RemoveReaderXSS(reader io.Reader) *bytes.Buffer {
	p := bluemonday.UGCPolicy()
	return p.SanitizeReader(reader)
}

// HTMLFilter 构建自定义的HTML标签过滤器
func HTMLFilter() *bluemonday.Policy {
	return bluemonday.NewPolicy()
}

func MyRemoveXSS(content string) string {
	return com.RemoveXSS(content)
}

func MyCleanText(value string) string {
	value = com.StripTags(value)
	value = com.RemoveEOL(value)
	return value
}

func MyCleanTags(value string) string {
	value = com.StripTags(value)
	return value
}
