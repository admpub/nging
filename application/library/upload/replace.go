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

package upload

import (
	"regexp"
)

var (
	fileURLDomainRegex = regexp.MustCompile(`^(?i)(http[s]?:)?//([^/]+)`)
	placeholderRegexp  = regexp.MustCompile(`\[storage:[\d]+\]`)
)

func CleanDomain(fileURL string) string {
	return fileURLDomainRegex.ReplaceAllString(fileURL, ``)
}

func ParseDomain(fileURL string) (scheme string, domain string) {
	matched := fileURLDomainRegex.FindStringSubmatch(fileURL)
	if len(matched) > 2 {
		scheme = matched[1]
		domain = matched[2]
	}
	return
}

// ReplacePlaceholder 从文本中替换占位符
var ReplacePlaceholder = func(s string, repl func(string) string) string {
	return placeholderRegexp.ReplaceAllStringFunc(s, func(find string) string {
		id := find[9 : len(find)-1]
		return repl(id)
	})
}
