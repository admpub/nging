package file

import (
	"fmt"
	"regexp"

	"github.com/webx-top/echo/param"
)

var (
	//defaultResString = `["'\(]([^"'#\(\)]+)#FileID-(\d+)["'\)]`
	defaultResString = `["'\(]([^"'#\(\)]+)["'\)]`
	regexResRule     = regexp.MustCompile(defaultResString)
)

// ReplaceEmbeddedRes 替换正文中的资源网址
func ReplaceEmbeddedRes(v string, reses map[string]interface{}) (r string) {
	for fid, rurl := range reses {
		re := regexp.MustCompile(`["'\(][^"'#\(\)]+#FileID-` + fid + `["'\)]`)
		v = re.ReplaceAllString(v, `"`+fmt.Sprint(rurl)+`"`)
	}
	return v
}

// EmbeddedRes 获取正文中的资源
func EmbeddedRes(v string, fn func(string, int64)) [][]string {
	list := regexResRule.FindAllStringSubmatch(v, -1)
	if fn == nil {
		return list
	}
	for _, a := range list {
		resource := a[1]
		var fileID int64
		if len(a) > 2 {
			fileID = param.AsInt64(a[2])
		}
		fn(resource, fileID)
	}
	return list
}
