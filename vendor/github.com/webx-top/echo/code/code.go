package code

import (
	"net/http"
)

type Code int

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

// Ok 是否是成功状态
func (c Code) Ok() bool {
	return c == Success
}

// Fail 是否是失败状态
func (c Code) Fail() bool {
	return c != Success
}

// Is 是否是期望状态
func (c Code) Is(expected Code) bool {
	return c == expected
}

// In 是否属于期望状态列表中的任一状态
func (c Code) In(expected ...Code) bool {
	for _, code := range expected {
		if c == code {
			return true
		}
	}
	return false
}

// Between 是否属于某个区间的状态码
func (c Code) Between(start Code, end Code) bool {
	return c >= start && c <= end
}
