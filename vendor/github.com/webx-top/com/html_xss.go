package com

import (
	"fmt"
	"regexp"
	"strings"
)

var (
	defaultXSSStringForTag       = `(</?)(?i)(applet|meta|xml|blink|link|style|script|embed|object|iframe|frame|frameset|ilayer|layer|bgsound|title|base%v)(\b[^>]*>)`
	defaultXSSStringForEvent     = `(<[a-zA-Z][a-zA-Z0-9]*\b[^<>]+\b)(?i)on([a-zA-Z]+)(\s*=[^<>]+>)`
	defaultXSSStringForAttrName  = `(<[a-zA-Z][a-zA-Z0-9]*\b[^<>]*\b)(?i)(style%v)(\s*=[^<>]+>)`
	defaultXSSStringForAttrValue = `(<[a-zA-Z][a-zA-Z0-9]*\b[^<>]*\b)([a-zA-Z]+)(\s*=\s*["']*\s*)(?i)(javascript|vbscript%v)(\b[^<>]+>)`

	removeXSSForTag       = regexp.MustCompile(fmt.Sprintf(defaultXSSStringForTag, ""))
	removeXSSForEvent     = regexp.MustCompile(defaultXSSStringForEvent)
	removeXSSForAttrName  = regexp.MustCompile(fmt.Sprintf(defaultXSSStringForAttrName, ""))
	removeXSSForAttrValue = regexp.MustCompile(fmt.Sprintf(defaultXSSStringForAttrValue, ""))
)

// RemoveXSS 删除XSS代码
func RemoveXSS(v string) (r string) {
	r = strings.Replace(v, `<!---->`, ``, -1)

	//过滤HTML标签
	r = removeXSSForTag.ReplaceAllString(r, `${1}_$2$3`)

	//过滤事件属性
	r = removeXSSForEvent.ReplaceAllString(r, `${1}_on$2$3`)

	//过滤属性
	r = removeXSSForAttrName.ReplaceAllString(r, `${1}_$2$3`)

	//过滤属性值
	r = removeXSSForAttrValue.ReplaceAllString(r, `${1}_$2${3}_$4$5`)

	//fmt.Println("Execute the filter: RemoveXSS.")
	return
}
