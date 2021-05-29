package factory

import (
	"net"
	"unicode"
)

type Timeout interface {
	Timeout() bool
}
type Temporary interface {
	Temporary() bool
}

func NetError(err error) net.Error {
	if netErr, ok := err.(net.Error); ok {
		return netErr
	}
	return nil
}

func IsTimeoutError(err error) bool {
	if e, ok := err.(Timeout); ok {
		return e.Timeout()
	}
	return false
}

func IsTemporaryError(err error) bool {
	if e, ok := err.(Temporary); ok {
		return e.Temporary()
	}
	return false
}

// ToSnakeCase : WebxTop => webx_top
func ToSnakeCase(name string) string {
	bytes := []rune{}
	for i, char := range name {
		if 'A' <= char && 'Z' >= char {
			char = unicode.ToLower(char)
			if i > 0 {
				bytes = append(bytes, '_')
			}
		}
		bytes = append(bytes, char)
	}
	return string(bytes)
}

// ToCamleCase : webx_top => WebxTop
func ToCamleCase(name string) string {
	underline := rune('_')
	isUnderline := false
	bytes := []rune{}
	for i, v := range name {
		if v == underline {
			isUnderline = true
			continue
		}
		if isUnderline || i == 0 {
			v = unicode.ToUpper(v)
		}
		isUnderline = false
		bytes = append(bytes, v)
	}
	return string(bytes)
}
