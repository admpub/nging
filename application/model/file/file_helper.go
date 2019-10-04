package file

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/webx-top/com"
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
		v = re.ReplaceAllString(v, `${1}`+rurl+`${2}`)
	}
	return v
}

// ReplaceRelatedRes 替换字段中的资源网址
func ReplaceRelatedRes(v string, reses map[uint64]string, seperator ...string) (r string) {
	var fileList []string
	var sep string
	if len(seperator) > 0 && len(seperator[0]) > 0 {
		sep = seperator[0]
		fileList = strings.Split(v, sep)
	} else {
		fileList = append(fileList, v)
	}
	replaced := map[int]struct{}{}
	for fid, rurl := range reses {
		suffix := `#FileID-` + fmt.Sprint(fid)
		for key, file := range fileList {
			if _, ok := replaced[key]; ok {
				continue
			}
			if strings.HasSuffix(file, suffix) {
				fileList[key] = rurl
				replaced[key] = struct{}{}
			}
		}
	}
	v = strings.Join(fileList, sep)
	return v
}

// EmbeddedRes 获取正文中的资源
func EmbeddedRes(v string, fn func(string, int64)) [][]string {
	if len(v) == 0 {
		return nil
	}
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

// RelatedRes 获取字段中关联的资源
func RelatedRes(v string, fn func(string, int64), seperator ...string) {
	if len(v) == 0 {
		return
	}
	var fileList []string
	if len(seperator) > 0 && len(seperator[0]) > 0 {
		fileList = strings.Split(v, seperator[0])
	} else {
		fileList = append(fileList, v)
	}
	for _, file := range fileList {
		file = strings.TrimSpace(file)
		if len(file) == 0 {
			fn(file, 0)
			continue
		}
		p := strings.LastIndex(file, `#FileID-`)
		if p < 0 {
			fn(file, 0)
			continue
		}
		var fid int64
		fileID := file[p+8:]
		if len(fileID) > 0 {
			fid = com.Int64(fileID)
		}
		file = file[0:p]
		fn(file, fid)
	}
}
