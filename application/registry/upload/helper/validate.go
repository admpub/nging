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
	"path"
	"strings"

	"github.com/webx-top/echo"
)

// IsRightUploadFile 是否是正确的上传文件
var IsRightUploadFile = func(ctx echo.Context, src string) error {
	src = path.Clean(src)
	ext := strings.ToLower(path.Ext(src))
	var ok bool
	var invalidExt string
	for _, ex := range AllowedUploadFileExtensions {
		if ext == ex {
			ok = true
			invalidExt = ext
			break
		}
	}
	if !ok {
		return ctx.E(`不支持的文件扩展名: %s`, invalidExt)
	}
	if !strings.Contains(src, UploadURLPath) {
		return ctx.E(`路径不合法`)
	}
	return nil
}
