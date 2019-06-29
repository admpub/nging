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

const (
	// StatusCaptchaError 验证码错误
	StatusCaptchaError = -9
	// StatusNonPrivileged 无权限
	StatusNonPrivileged = -2
	// StatusNotLoggedIn 未登录
	StatusNotLoggedIn = -1
	// StatusFailure 操作失败
	StatusFailure = 0
	// StatusSuccess 操作成功
	StatusSuccess = 1
)

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
