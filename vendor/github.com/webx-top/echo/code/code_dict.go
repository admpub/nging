package code

import "net/http"

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

	CaptchaCodeRequired: {"CaptchaCodeRequired", http.StatusOK},
	CaptchaIdMissing:    {"CaptchaIdMissing", http.StatusOK},
	CaptchaError:        {"CaptchaError", http.StatusOK},
	BalanceNoEnough:     {"BalanceNoEnough", http.StatusOK},
	UserDisabled:        {"UserDisabled", http.StatusOK},
	UserNotFound:        {"UserNotFound", http.StatusOK},
	NonPrivileged:       {"NonPrivileged", http.StatusOK},
	Unauthenticated:     {"Unauthenticated", http.StatusOK},

	// - 通用

	Failure: {"Failure", http.StatusOK},
	Success: {"Success", http.StatusOK},
}
