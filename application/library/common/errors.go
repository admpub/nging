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
	"encoding/gob"
	"errors"
	"strings"

	"github.com/webx-top/echo"
	"github.com/webx-top/echo/code"
)

func init() {
	gob.Register(&Success{})
}

var (
	// - JSON

	//ErrUserNotLoggedIn 用户未登录
	ErrUserNotLoggedIn = echo.NewError(`User not logged in`, code.Unauthenticated)
	//ErrUserNotFound 用户不存在
	ErrUserNotFound = echo.NewError(`User does not exist`, code.UserNotFound)
	//ErrUserNoPerm 用户无权限
	ErrUserNoPerm = echo.NewError(`User has no permission`, code.NonPrivileged)
	//ErrUserDisabled 用户已被禁用
	ErrUserDisabled = echo.NewError(`User has been disabled`, code.UserDisabled)
	//ErrBalanceNoEnough 余额不足
	ErrBalanceNoEnough = echo.NewError(`Balance is not enough`, code.BalanceNoEnough)
	//ErrCaptcha 验证码错误
	ErrCaptcha = echo.NewError(`Captcha is incorrect`, code.CaptchaError)
	//ErrCaptchaIdMissing 缺少captchaId
	ErrCaptchaIdMissing = echo.NewError(`Missing captchaId`, code.CaptchaIdMissing).SetZone(`captchaId`)
	//ErrInvalidAppID App ID 无效
	ErrInvalidAppID = echo.NewError(`Invalid app id`, code.InvalidAppID)
	//ErrInvalidSign 无效签名
	ErrInvalidSign = echo.NewError(`Invalid sign`, code.InvalidSignature)
	//ErrInvalidToken 令牌无效
	ErrInvalidToken = echo.NewError(`Invalid token`, code.InvalidToken)

	// - Operation

	//ErrRepeatOperation 重复操作
	ErrRepeatOperation = echo.NewError(`Repeat operation`, code.RepeatOperation)
	//ErrUnsupported 不支持
	ErrUnsupported = echo.NewError(`Unsupported`, code.Unsupported)
	//ErrOperationTimeout 操作超时
	ErrOperationTimeout = echo.NewError(`Operation timeout`, code.OperationTimeout)
	//ErrOperationFail 操作失败
	ErrOperationFail = echo.NewError(`Operation fail`, code.Failure)

	// - HTTP

	//ErrResponseFormatError 响应格式错误
	ErrResponseFormatError = echo.NewError(`Response format error`, code.AbnormalResponse)
	//ErrRequestTimeout 提交超时
	ErrRequestTimeout = echo.NewError(`Request timeout`, code.RequestTimeout)
	//ErrRequestFail 提交失败
	ErrRequestFail = echo.NewError(`Request fail`, code.RequestFailure)

	// - Watcher

	// ErrIgnoreConfigChange 忽略配置文件更改
	ErrIgnoreConfigChange = errors.New(`Ignore configuration file changes`)

	// - Checker

	// ErrNext 需要继续向下检查
	ErrNext           = errors.New("Next")
	ErrConcurrentLock = errors.New("Concurrent lock has been triggered")
)

// DefaultNopMessage 默认空消息
var DefaultNopMessage Messager = &NopMessage{}

// Errors 多个错误信息
type Errors []error

func (e Errors) Error() string {
	s := make([]string, len(e))
	for k, v := range e {
		s[k] = v.Error()
	}
	return strings.Join(s, "\n")
}

func (e Errors) String() string {
	return e.Error()
}

// NopMessage 空消息
type NopMessage struct {
}

// Error 错误信息
func (n *NopMessage) Error() string {
	return ``
}

// Success 成功信息
func (n *NopMessage) Success() string {
	return ``
}

// String 信息字符串
func (n *NopMessage) String() string {
	return ``
}

// Messager 信息接口
type Messager interface {
	Successor
	error
}

// IsMessage 判断err是否为Message
func IsMessage(err interface{}) bool {
	_, y := err.(Messager)
	return y
}

// Message 获取err中的信息接口
func Message(err interface{}) Messager {
	if v, y := err.(Messager); y {
		return v
	}
	return DefaultNopMessage
}

// NewOk 创建成功信息
func NewOk(v string) Successor {
	return &Success{
		Value: v,
	}
}

// Success 成功信息
type Success struct {
	Value string
}

// Success 成功信息
func (s *Success) Success() string {
	return s.Value
}

func (s *Success) String() string {
	return s.Value
}

// Successor 成功信息接口
type Successor interface {
	Success() string
}

// IsError 是否是错误信息
func IsError(err interface{}) bool {
	_, y := err.(error)
	return y
}

// IsOk 是否是成功信息
func IsOk(err interface{}) bool {
	_, y := err.(Successor)
	return y
}

// OkString 获取成功信息
func OkString(err interface{}) string {
	if v, y := err.(Successor); y {
		return v.Success()
	}
	return ``
}
