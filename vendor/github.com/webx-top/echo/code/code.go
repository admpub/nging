package code

import (
	"net/http"
)

type (
	Code         int
	TextHTTPCode struct {
		Text     string
		HTTPCode int
	}
	CodeMap map[Code]TextHTTPCode
)

const (
	// - 操作状态

	RequestFailure   Code = -205  //提交失败
	RequestTimeout   Code = -204  //提交超时
	AbnormalResponse Code = -203  //响应异常
	OperationTimeout Code = -202  //操作超时
	Unsupported      Code = -201  //不支持的操作
	RepeatOperation  Code = -200  //重复操作

	// - 数据状态

	DataNotFound     Code = -100  //数据未找到

	// - 用户状态

	BalanceNoEnough  Code = -5  //余额不足
	UserDisabled     Code = -4  //用户被禁用
	UserNotFound     Code = -3  //用户未找到
	NonPrivileged    Code = -2  //无权限
	Unauthenticated  Code = -1  //未登录

	// - 通用

	Failure          Code = 0   //操作失败
	Success          Code = 1   //操作成功
)

// CodeDict 状态码字典
var CodeDict = CodeMap{
	BalanceNoEnough:  {"BalanceNoEnough", http.StatusOK},
	RequestTimeout:   {"RequestTimeout", http.StatusOK},
	AbnormalResponse: {"AbnormalResponse", http.StatusOK},
	OperationTimeout: {"OperationTimeout", http.StatusOK},
	Unsupported:      {"Unsupported", http.StatusOK},
	RepeatOperation:  {"RepeatOperation", http.StatusOK},
	DataNotFound:     {"DataNotFound", http.StatusOK},
	UserNotFound:     {"UserNotFound", http.StatusOK},
	NonPrivileged:    {"NonPrivileged", http.StatusOK},
	Unauthenticated:  {"Unauthenticated", http.StatusOK},
	Failure:          {"Failure", http.StatusOK},
	Success:          {"Success", http.StatusOK},
}

func (c Code) String() string {
	if v, y := CodeDict[c]; y {
		return v.Text
	}
	return `Undefined`
}

// Int 返回int类型的自定义状态码
func (c Code) Int() int {
	return int(c)
}

// HTTPCode 返回HTTP状态码
func (c Code) HTTPCode() int {
	if v, y := CodeDict[c]; y {
		return v.HTTPCode
	}
	return http.StatusOK
}

func (s CodeMap) Get(code Code) TextHTTPCode {
	v, _ := s[code]
	return v
}

func (s CodeMap) GetByInt(code int) TextHTTPCode {
	return s.Get(Code(code))
}
