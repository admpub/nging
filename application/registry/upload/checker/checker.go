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

package checker

import (
	"time"

	"github.com/webx-top/echo"
	"github.com/webx-top/echo/code"
)

// UploadURLMaxAge 上传链接生存时间
var UploadURLMaxAge int64 = 86400

// Checker 验证并生成子文件夹名称和文件名称
type Checker func(echo.Context) (subdir string, name string, err error)

// Default 默认Checker
var Default = func(ctx echo.Context) (subdir string, name string, err error) {
	timestamp := ctx.Formx(`time`).Int64()
	// 验证签名（避免上传接口被滥用）
	if ctx.Form(`token`) != Token(ctx.Queries()) {
		err = ctx.NewError(code.InvalidParameter, ctx.T(`令牌错误`))
		return
	}
	if time.Now().Unix()-timestamp > UploadURLMaxAge {
		err = ctx.NewError(code.DataHasExpired, ctx.T(`上传网址已过期`))
		return
	}
	subdir = time.Now().Format(`2006/01/02/`)
	return
}
