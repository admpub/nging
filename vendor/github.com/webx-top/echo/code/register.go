package code

func Register(code Code, text string, httpCodes ...int) {
	CodeDict.Set(code, text, httpCodes...)
}

func Get(code Code) TextHTTPCode {
	return CodeDict.Get(code)
}
