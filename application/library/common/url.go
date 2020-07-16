package common

import (
	"net/url"
	"strings"

	"github.com/webx-top/echo"
)

var DefaultReturnToURLVarName = `return_to`

func GetReturnURL(ctx echo.Context, varNames ...string) string {
	varName := DefaultReturnToURLVarName
	if len(varNames) > 0 && len(varNames[0]) > 0 {
		varName = varNames[0]
	}
	returnTo := ctx.Form(varName)
	if returnTo == ctx.Request().URL().Path() {
		returnTo = ``
	}
	return returnTo
}

func ReturnToCurrentURL(ctx echo.Context, varNames ...string) string {
	varName := DefaultReturnToURLVarName
	if len(varNames) > 0 && len(varNames[0]) > 0 {
		varName = varNames[0]
	}
	returnTo := ctx.Form(varName)
	if len(returnTo) == 0 {
		returnTo = ctx.Request().URI()
	}
	return returnTo
}

func WithURLParams(urlStr string, key string, value string, args ...string) string {
	if strings.Contains(urlStr, `?`) {
		urlStr += `&`
	} else {
		urlStr += `?`
	}
	urlStr += key + `=` + url.QueryEscape(value)
	var k string
	for i, j := 0, len(args); i < j; i++ {
		if i%2 == 0 {
			k = args[i]
			continue
		}
		urlStr += `&` + k + `=` + url.QueryEscape(args[i])
		k = ``
	}
	if len(k) > 0 {
		urlStr += `&` + k + `=`
	}
	return urlStr
}

func WithReturnURL(ctx echo.Context, urlStr string, varNames ...string) string {
	varName := DefaultReturnToURLVarName
	if len(varNames) > 0 && len(varNames[0]) > 0 {
		varName = varNames[0]
	}
	withVarName := varName
	if len(varNames) > 1 && len(varNames[1]) > 0 {
		withVarName = varNames[1]
	}

	returnTo := GetReturnURL(ctx, varName)
	if len(returnTo) == 0 {
		return urlStr
	}
	return WithURLParams(urlStr, withVarName, returnTo)
}

func GetOtherURL(ctx echo.Context, returnTo string) string {
	if len(returnTo) == 0 {
		return returnTo
	}
	urlInfo, _ := url.Parse(returnTo)
	if urlInfo == nil || urlInfo.Path == ctx.Request().URL().Path() {
		returnTo = ``
	}
	return returnTo
}
