package code

import "net/http"

type TextHTTPCode struct {
	Text     string
	HTTPCode int
}

type CodeMap map[Code]TextHTTPCode

func (s CodeMap) Get(code Code) TextHTTPCode {
	v, _ := s[code]
	return v
}

func (s CodeMap) Set(code Code, text string, httpCodes ...int) CodeMap {
	httpCode := http.StatusOK
	if len(httpCodes) > 0 {
		httpCode = httpCodes[0]
		if httpCode <= 0 {
			httpCode = http.StatusOK
		}
	}
	s[code] = TextHTTPCode{Text: text, HTTPCode: httpCode}
	return s
}

func (s CodeMap) GetByInt(code int) TextHTTPCode {
	return s.Get(Code(code))
}
