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

import (
	"github.com/webx-top/echo"
)

type Responser func(ctx echo.Context, field string, err error, imageURLs []string) (result echo.H, embed bool)

var responsers = map[string]Responser{
	`editormd-image-file`: func(ctx echo.Context, field string, err error, imageURLs []string) (result echo.H, embed bool) {
		var code int
		var msg string
		if err == nil {
			code = 1
		} else {
			msg = err.Error()
		}
		result = echo.H{
			`success`: code, // 0 表示上传失败，1 表示上传成功
			`message`: msg,
			//`url`     : imageURLs        // 上传成功时才返回
		}
		if len(imageURLs) > 0 {
			result[`url`] = imageURLs[0]
		}
		return
	},
}

var DefaultResponser = func(ctx echo.Context, field string, err error, imageURLs []string) (result echo.H, embed bool) {
	return echo.H{
		`files`: imageURLs,
	}, true
}

func ResponserRegister(field string, responser Responser) {
	responsers[field] = responser
}

func ResponserAll() map[string]Responser {
	return responsers
}

func ResponserGet(field string) Responser {
	responser, ok := responsers[field]
	if !ok {
		return DefaultResponser
	}
	return responser
}
