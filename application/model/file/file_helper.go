package file

import (
	"regexp"

	"github.com/webx-top/echo/param"
)

var (
	//defaultResString = `["'\(]([^"'#\(\)]+)#FileID-(\d+)["'\)]`
	defaultResString = `["'\(]([^"'#\(\)]+)["'\)]`
	regexResRule     = regexp.MustCompile(defaultResString)
)

// ReplaceEmbeddedRes 替换正文中的资源网址
func ReplaceEmbeddedRes(v string, reses map[string]string) (r string) {
	for fid, rurl := range reses {
		re, _ := regexp.Compile(`["'\(][^"'#\(\)]+#FileID-` + fid + `["'\)]`)
		v = re.ReplaceAllString(v, `"`+rurl+`"`)
	}
	return v
}

// EmbeddedRes 获取正文中的资源
func EmbeddedRes(v string, fn func(string, int64)) [][]string {
	list := regexResRule.FindAllStringSubmatch(v, -1)
	if fn != nil {
		for _, a := range list {
			resource := a[1]
			var fileID int64
			if len(a) > 2 {
				fileID = param.AsInt64(a[2])
			}
			fn(resource, fileId)
		}
	}
	return list
}
