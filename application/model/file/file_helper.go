package file

import (
	"fmt"
	"regexp"

	"github.com/webx-top/echo/param"
)

const (
	fileCSPattern    = `(?:[^"'#\(\)]+)`
	fileStartPattern = `["'\(]`
	fileEndPattern   = `["'\)]`
)

var (
	//defaultFilePattern = `["'\(]([^"'#\(\)]+)#FileID-(\d+)["'\)]`
	filePattern = fileStartPattern + `(` + fileCSPattern + `\.(?:[\w]+)` + fileCSPattern + `?)` + fileEndPattern
	fileRGX     = regexp.MustCompile(filePattern)
)

// ReplaceEmbeddedRes 替换正文中的资源网址
func ReplaceEmbeddedRes(v string, reses map[uint64]string) (r string) {
	for fid, rurl := range reses {
		re := regexp.MustCompile(`(` + fileStartPattern + `)` + fileCSPattern + `#FileID-` + fmt.Sprint(fid) + `(` + fileEndPattern + `)`)
		fmt.Println(`(` + fileStartPattern + `)` + fileCSPattern + `#FileID-` + fmt.Sprint(fid) + `(` + fileEndPattern + `)`)
		v = re.ReplaceAllString(v, `${1}`+rurl+`${2}`)
	}
	return v
}

// EmbeddedRes 获取正文中的资源
func EmbeddedRes(v string, fn func(string, int64)) [][]string {
	list := fileRGX.FindAllStringSubmatch(v, -1)
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
