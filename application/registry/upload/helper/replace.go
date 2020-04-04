/*
   Nging is a toolbox for webmasters
   Copyright (C) 2018-present  Wenhui Shen <swh@admpub.com>

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
	"regexp"
)

var (
	temporaryFileRegexp  *regexp.Regexp
	persistentFileRegexp *regexp.Regexp
	anyFileRegexp        *regexp.Regexp
	placeholderRegexp = regexp.MustCompile(`\[storage:[\d]+\]`)
)

func init() {
	Init()
}

func Init() {
	ruleEnd := ExtensionRegexpEnd()
	temporaryFileRegexp = regexp.MustCompile(UploadURLPath + `[\w-]+/0/[\w]+` + ruleEnd)
	persistentFileRegexp = regexp.MustCompile(UploadURLPath + `[\w-]+/([^0]|[0-9]{2,})/[\w]+` + ruleEnd)
	anyFileRegexp = regexp.MustCompile(UploadURLPath + `[\w-]+/([\w-]+/)+[\w-]+` + ruleEnd)
}

// ParseTemporaryFileName 从文本中解析出临时文件名称
var ParseTemporaryFileName = func(s string) []string {
	files := temporaryFileRegexp.FindAllString(s, -1)
	return files
}

// ParsePersistentFileName 从文本中解析出正式文件名称
var ParsePersistentFileName = func(s string) []string {
	files := persistentFileRegexp.FindAllString(s, -1)
	return files
}

// ParseAnyFileName 从文本中解析出任意上传文件名称
var ParseAnyFileName = func(s string) []string {
	files := anyFileRegexp.FindAllString(s, -1)
	return files
}

// ReplaceAnyFileName 从文本中替换任意上传文件名称
var ReplaceAnyFileName = func(s string, repl func(string) string) string {
	return anyFileRegexp.ReplaceAllStringFunc(s, repl)
}

// ReplacePlaceholder 从文本中替换占位符
var ReplacePlaceholder = func(s string, repl func(string) string) string {
	return placeholderRegexp.ReplaceAllStringFunc(s, func(find string) string{
		id := find[9:len(find)-1]
		return repl(id)
	})
}
