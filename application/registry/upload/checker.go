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
	"fmt"
	"time"

	"github.com/webx-top/echo"
)

// UploadLinkLifeTime 上传链接生存时间
var UploadLinkLifeTime int64 = 86400

type Checker func(echo.Context) (subdir string, name string, err error)

var DefaultChecker = func(ctx echo.Context) (subdir string, name string, err error) {
	refid := ctx.Formx(`refid`).Uint64()
	timestamp := ctx.Formx(`time`).Int64()
	// 验证签名（避免上传接口被滥用）
	if ctx.Form(`token`) != Token(`refid`, refid, `time`, timestamp) {
		err = ctx.E(`令牌错误`)
		return
	}
	if time.Now().Local().Unix()-timestamp > UploadLinkLifeTime {
		err = ctx.E(`上传网址已过期`)
		return
	}
	subdir = fmt.Sprint(refid) + `/`
	//subdir = time.Now().Format(`2006/01/02/`)
	return
}

var checkers = map[string]Checker{
	`customer-avatar`: func(ctx echo.Context) (subdir string, name string, err error) {
		customerID := ctx.Formx(`customerId`).Uint64()
		timestamp := ctx.Formx(`time`).Int64()
		// 验证签名（避免上传接口被滥用）
		if ctx.Form(`token`) != Token(`customerId`, customerID, `time`, timestamp) {
			err = ctx.E(`令牌错误`)
			return
		}
		if time.Now().Local().Unix()-timestamp > UploadLinkLifeTime {
			err = ctx.E(`上传网址已过期`)
			return
		}
		if customerID > 0 {
			name = `avatar`
		}
		subdir = fmt.Sprint(customerID) + `/`
		return
	},
	`user-avatar`: func(ctx echo.Context) (subdir string, name string, err error) {
		userID := ctx.Formx(`userId`).Uint64()
		timestamp := ctx.Formx(`time`).Int64()
		// 验证签名（避免上传接口被滥用）
		if ctx.Form(`token`) != Token(`userId`, userID, `time`, timestamp) {
			err = ctx.E(`令牌错误`)
			return
		}
		if time.Now().Local().Unix()-timestamp > UploadLinkLifeTime {
			err = ctx.E(`上传网址已过期`)
			return
		}
		if userID > 0 {
			name = `avatar`
		}
		subdir = fmt.Sprint(userID) + `/`
		return
	},
}

func CheckerRegister(typ string, checker Checker) {
	checkers[typ] = checker
}

func CheckerAll() map[string]Checker {
	return checkers
}

func CheckerGet(typ string) Checker {
	checker, ok := checkers[typ]
	if !ok {
		return DefaultChecker
	}
	return checker
}
