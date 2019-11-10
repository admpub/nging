package com

import (
	"fmt"
	"regexp"
	"strings"
)

var (
	defaultXSSStringForTag       = `(</?)(?i)(applet|meta|xml|blink|link|style|script|embed|object|iframe|frame|frameset|ilayer|layer|bgsound|title|base%v)(\b[^>]*>)`
	defaultXSSStringForEvent     = `(<[a-zA-Z][a-zA-Z0-9]*\b[^<>]+\b)(?i)on(abort|activate|afterprint|afterupdate|beforeactivate|beforecopy|beforecut|beforedeactivate|beforeeditfocus|beforepaste|beforeprint|beforeunload|beforeupdate|blur|bounce|cellchange|change|click|contextmenu|controlselect|copy|cut|dataavailable|datasetchanged|datasetcomplete|dblclick|deactivate|drag|dragend|dragenter|dragleave|dragover|dragstart|drop|error|errorupdate|filterchange|finish|focus|focusin|focusout|help|keydown|keypress|keyup|layoutcomplete|load|losecapture|message|mousedown|mouseenter|mouseleave|mousemove|mouseout|mouseover|mouseup|mousewheel|move|moveend|movestart|paste|propertychange|readystatechange|reset|resize|resizeend|resizestart|rowenter|rowexit|rowsdelete|rowsinserted|scroll|select|selectionchange|selectstart|start|stop|submit|unload%v)(\s*=[^<>]+>)`
	defaultXSSStringForAttrName  = `(<[a-zA-Z][a-zA-Z0-9]*\b[^<>]*\b)(?i)(style%v)(\s*=[^<>]+>)`
	defaultXSSStringForAttrValue = `(<[a-zA-Z][a-zA-Z0-9]*\b[^<>]*\b)([a-zA-Z]+)(\s*=\s*["']*\s*)(?i)(javascript|vbscript%v)(\b[^<>]+>)`
	defaultAllowHTMLTag          = `(<[a-zA-Z][a-zA-Z0-9]*\b[^<>]*\b)([a-zA-Z]+)(\s*=\s*["']*\s*)(?i)(javascript|vbscript%v)(\b[^<>]+>)`

	removeXSSForTag       = regexp.MustCompile(fmt.Sprintf(defaultXSSStringForTag, ""))
	removeXSSForEvent     = regexp.MustCompile(fmt.Sprintf(defaultXSSStringForEvent, ""))
	removeXSSForAttrName  = regexp.MustCompile(fmt.Sprintf(defaultXSSStringForAttrName, ""))
	removeXSSForAttrValue = regexp.MustCompile(fmt.Sprintf(defaultXSSStringForAttrValue, ""))

	cleanHTMLTag = regexp.MustCompile(fmt.Sprintf(defaultXSSStringForAttrValue, ""))
)

// RemoveXSS 删除XSS代码
func RemoveXSS(v string) (r string) {
	r = strings.Replace(v, `<!---->`, ``, -1)

	//过滤HTML标签
	r = removeXSSForTag.ReplaceAllString(r, `${1}k_$2$3`)

	//过滤事件属性
	r = removeXSSForEvent.ReplaceAllString(r, `${1}k_on$2$3`)

	//过滤属性
	r = removeXSSForAttrName.ReplaceAllString(r, `${1}k_$2$3`)

	//过滤属性值
	r = removeXSSForAttrValue.ReplaceAllString(r, `${1}k_$2${3}k_$4$5`)

	//fmt.Println("Execute the fiter: RemoveXSS.")
	return
}
