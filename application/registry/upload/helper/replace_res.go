/*
   Nging is a toolbox for webmasters
   Copyright (C) 2018-present Wenhui Shen <swh@admpub.com>

   This program is free software: you can redistribute it and/or modify
   it under the terms of the GNU Affero General Public License as published
   by the Free Software Foundation, either version 3 of the License, or
   (at your option) any later version.

   This program is distributed in the hope that it will be useful,
   but WITHOUT ANY WARRANTY; without even the implied warranty of
   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
   GNU Affero General Public License for more details.

   You should have received a copy of the GNU Affero General Public License
   along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/

package helper

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/webx-top/com"
	"github.com/webx-top/echo/param"
)

const (
	fileNamePattern  = `(?:[^"'#\(\)]+)` //文件名
	fileStartPattern = `["'\(]`          //起始符号
	fileEndPattern   = `["'\)]`          //终止符号
)

var (
	filePattern = fileStartPattern + `(` + fileNamePattern + `\.(?:[\w]+)` + fileNamePattern + `?)` + fileEndPattern
	fileRGX     = regexp.MustCompile(filePattern)
)

// ReplaceEmbeddedResID 替换正文中的资源网址
func ReplaceEmbeddedResID(v string, reses map[uint64]string) (r string) {
	for fid, rurl := range reses {
		re := regexp.MustCompile(`(` + fileStartPattern + `)` + fileNamePattern + `#FileID-` + fmt.Sprint(fid) + `(` + fileEndPattern + `)`)
		v = re.ReplaceAllString(v, `${1}`+rurl+`${2}`)
	}
	return v
}

// ReplaceRelatedResID 替换字段中的资源网址
func ReplaceRelatedResID(v string, reses map[uint64]string, seperator ...string) (r string) {
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

// ReplaceEmbeddedRes 替换正文中的资源网址
func ReplaceEmbeddedRes(v string, reses map[string]string) (r string) {
	for furl, rurl := range reses {
		re := regexp.MustCompile(`(` + fileStartPattern + `)` + regexp.QuoteMeta(furl) + `(` + fileEndPattern + `)`)
		v = re.ReplaceAllString(v, `${1}`+rurl+`${2}`)
	}
	return v
}

// ReplaceRelatedRes 替换字段中的资源网址
func ReplaceRelatedRes(v string, reses map[string]string, seperator ...string) (r string) {
	var fileList []string
	var sep string
	if len(seperator) > 0 && len(seperator[0]) > 0 {
		sep = seperator[0]
		fileList = strings.Split(v, sep)
	} else {
		fileList = append(fileList, v)
	}
	for key, file := range fileList {
		rurl, ok := reses[file]
		if !ok {
			continue
		}
		fileList[key] = rurl
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
			continue
		}
		p := strings.LastIndex(file, `#FileID-`)
		if p < 0 {
			if com.StrIsNumeric(file) {
				fn(``, com.Int64(file))
			} else {
				fn(file, 0)
			}
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
