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

package common

import (
	stdCode "github.com/webx-top/echo/code"
)

// IsCaptchaErrCode 是否验证码错误码
func IsCaptchaErrCode(code stdCode.Code) bool {
	return code == stdCode.CaptchaError
}

// IsCaptchaError 用户验证码错误
func IsCaptchaError(err error) bool {
	return err == ErrCaptcha
}

// IsUserNotLoggedIn 用户是否未登录
func IsUserNotLoggedIn(err error) bool {
	return err == ErrUserNotLoggedIn
}

// IsUserNotFound 用户是否不存在
func IsUserNotFound(err error) bool {
	return err == ErrUserNotFound
}

// IsUserNoPerm 用户是否没有操作权限
func IsUserNoPerm(err error) bool {
	return err == ErrUserNoPerm
}

// IsUserDisabled 用户是否被禁用
func IsUserDisabled(err error) bool {
	return err == ErrUserDisabled
}
