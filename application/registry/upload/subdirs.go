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

package upload

import "strings"

var subdirs = map[string]bool{
	`user`:              true, //后台用户个人文件
	`customer`:          true, //前台客户个人文件
	`comment`:           true, //评论附件
	`customer-avatar`:   true, //前台客户头像
	`user-avatar`:       true, //后台用户头像
	`pay-product`:       true, //产品图片
	`product-version`:   true, //产品版本文件
	`site-announcement`: true, //站点公告图片
	`news`:              true, //新闻图片
}

func SubdirRegister(subdir string, allow bool) {
	subdirs[subdir] = allow
}

func SubdirAll() map[string]bool {
	return subdirs
}

func SubdirIsAllowed(subdir string) bool {
	allow, _ := subdirs[subdir]
	return allow
}

// CleanTempFile 清理临时文件
func CleanTempFile(prefix string, deleter func(folderPath string) error) error {
	if !strings.HasSuffix(prefix, `/`) {
		prefix += `/`
	}
	for subdir := range subdirs {
		err := deleter(prefix + subdir + `/0/`)
		if err != nil {
			return err
		}
	}
	return nil
}
