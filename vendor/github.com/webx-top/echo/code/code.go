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
	// - 系统状态

	SystemUnauthorized Code = -301 // 系统未获得授权
	SystemNotInstalled Code = -300 // 系统未安装

	// - 操作状态

	FrequencyTooFast    Code = -207 //操作频率太快
	OperationProcessing Code = -206 //操作处理中
	RequestFailure      Code = -205 //提交失败
	RequestTimeout      Code = -204 //提交超时
	AbnormalResponse    Code = -203 //响应异常
	OperationTimeout    Code = -202 //操作超时
	Unsupported         Code = -201 //不支持的操作
	RepeatOperation     Code = -200 //重复操作

	// - 数据状态

	InvalidToken Code = -151 //令牌错误
	InvalidAppID Code = -150 //AppID不正确

	DataSizeTooBig      Code = -110 //数据尺寸太大
	DataAlreadyExists   Code = -109 //数据已经存在
	DataFormatIncorrect Code = -108 //数据格式不正确
	DataStatusIncorrect Code = -107 //数据状态不正确
	DataProcessing      Code = -106 //数据未处理中状态
	DataUnavailable     Code = -105 //尚未启用
	DataHasExpired      Code = -104 //数据已经过期
	InvalidType         Code = -103 //类型不正确
	InvalidSignature    Code = -102 //无效的签名
	InvalidParameter    Code = -101 //无效的参数
	DataNotFound        Code = -100 //数据未找到

	// - 用户状态

	CaptchaError    Code = -9 //验证码错误
	BalanceNoEnough Code = -5 //余额不足
	UserDisabled    Code = -4 //用户被禁用
	UserNotFound    Code = -3 //用户未找到
	NonPrivileged   Code = -2 //无权限
	Unauthenticated Code = -1 //未登录

	// - 通用

	Failure Code = 0 //操作失败
	Success Code = 1 //操作成功
)

// CodeDict 状态码字典
var CodeDict = CodeMap{

	// - 系统状态

	SystemUnauthorized: {"SystemUnauthorized", http.StatusOK},
	SystemNotInstalled: {"SystemNotInstalled", http.StatusOK},

	// - 操作状态

	OperationProcessing: {"OperationProcessing", http.StatusOK},
	FrequencyTooFast:    {"FrequencyTooFast", http.StatusOK},
	RequestFailure:      {"RequestFailure", http.StatusOK},
	RequestTimeout:      {"RequestTimeout", http.StatusOK},
	AbnormalResponse:    {"AbnormalResponse", http.StatusOK},
	OperationTimeout:    {"OperationTimeout", http.StatusOK},
	Unsupported:         {"Unsupported", http.StatusOK},
	RepeatOperation:     {"RepeatOperation", http.StatusOK},

	// - 数据状态

	InvalidAppID: {"InvalidAppID", http.StatusOK},
	InvalidToken: {"InvalidToken", http.StatusOK},

	DataSizeTooBig:      {"DataSizeTooBig", http.StatusOK},
	DataAlreadyExists:   {"DataAlreadyExists", http.StatusOK},
	DataFormatIncorrect: {"DataFormatIncorrect", http.StatusOK},
	DataStatusIncorrect: {"DataStatusIncorrect", http.StatusOK},
	DataHasExpired:      {"DataHasExpired", http.StatusOK},
	DataProcessing:      {"DataProcessing", http.StatusOK},
	DataUnavailable:     {"DataUnavailable", http.StatusOK},
	InvalidType:         {"InvalidType", http.StatusOK},
	InvalidSignature:    {"InvalidSignature", http.StatusOK},
	InvalidParameter:    {"InvalidParameter", http.StatusOK},
	DataNotFound:        {"DataNotFound", http.StatusOK},

	// - 用户状态

	CaptchaError:    {"CaptchaError", http.StatusOK},
	BalanceNoEnough: {"BalanceNoEnough", http.StatusOK},
	UserDisabled:    {"UserDisabled", http.StatusOK},
	UserNotFound:    {"UserNotFound", http.StatusOK},
	NonPrivileged:   {"NonPrivileged", http.StatusOK},
	Unauthenticated: {"Unauthenticated", http.StatusOK},

	// - 通用

	Failure: {"Failure", http.StatusOK},
	Success: {"Success", http.StatusOK},
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

func (s CodeMap) Set(code Code, text string, httpCodes ...int) CodeMap {
	httpCode := http.StatusOK
	if len(httpCodes) > 0 {
		httpCode = httpCodes[0]
	}
	s[code] = TextHTTPCode{Text: text, HTTPCode: httpCode}
	return s
}

func (s CodeMap) GetByInt(code int) TextHTTPCode {
	return s.Get(Code(code))
}
