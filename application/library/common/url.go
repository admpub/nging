package common

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/webx-top/db"
	"github.com/webx-top/echo"
)

var DefaultNextURLVarName = `next`

func GetNextURL(ctx echo.Context, varNames ...string) string {
	varName := DefaultNextURLVarName
	if len(varNames) > 0 && len(varNames[0]) > 0 {
		varName = varNames[0]
	}
	next := ctx.Form(varName)
	if next == ctx.Request().URL().Path() {
		next = ``
	}
	return next
}

func ReturnToCurrentURL(ctx echo.Context, varNames ...string) string {
	varName := DefaultNextURLVarName
	if len(varNames) > 0 && len(varNames[0]) > 0 {
		varName = varNames[0]
	}
	next := ctx.Form(varName)
	if len(next) == 0 {
		next = ctx.Request().URI()
	}
	return next
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

func WithNextURL(ctx echo.Context, urlStr string, varNames ...string) string {
	varName := DefaultNextURLVarName
	if len(varNames) > 0 && len(varNames[0]) > 0 {
		varName = varNames[0]
	}
	withVarName := varName
	if len(varNames) > 1 && len(varNames[1]) > 0 {
		withVarName = varNames[1]
	}

	next := GetNextURL(ctx, varName)
	if len(next) == 0 || next == urlStr {
		return urlStr
	}
	if next[0] == '/' {
		if len(urlStr) > 8 {
			var urlCopy string
			switch strings.ToLower(urlStr[0:7]) {
			case `https:/`:
				urlCopy = urlStr[8:]
			case `http://`:
				urlCopy = urlStr[7:]
			}
			if len(urlCopy) > 0 {
				p := strings.Index(urlCopy, `/`)
				if p > 0 && urlCopy[p:] == next {
					return urlStr
				}
			}
		}
	}

	return WithURLParams(urlStr, withVarName, next)
}

func GetOtherURL(ctx echo.Context, next string) string {
	if len(next) == 0 {
		return next
	}
	urlInfo, _ := url.Parse(next)
	if urlInfo == nil || urlInfo.Path == ctx.Request().URL().Path() {
		next = ``
	}
	return next
}

func FullURL(domianURL string, myURL string) string {
	if IsFullURL(myURL) {
		return myURL
	}
	if !strings.HasPrefix(myURL, `/`) && !strings.HasSuffix(domianURL, `/`) {
		myURL = `/` + myURL
	}
	myURL = domianURL + myURL
	return myURL
}

func IsFullURL(purl string) bool {
	if len(purl) == 0 {
		return false
	}
	if purl[0] == '/' {
		return false
	}
	// find "://"
	firstPos := strings.Index(purl, `/`)
	if firstPos < 0 || firstPos == len(purl)-1 {
		return false
	}
	if firstPos > 1 && purl[firstPos-1] == ':' && purl[firstPos+1] == '/' {
		return true
	}
	return false
}

func NewSortedURLValues(query string) SortedURLValues {
	r := SortedURLValues{}
	r.ParseQuery(query)
	return r
}

type SortedURLValues []*URLValues

type URLValues struct {
	Key    string
	Values []string
}

// ParseQuery 解析 URL Query
// copy from standard library src/net/url/url.go: func parseQuery(m Values, query string) (err error)
func (s *SortedURLValues) ParseQuery(query string) (err error) {
	indexes := map[string]int{}
	if len(*s) > 0 {
		for k, v := range *s {
			indexes[v.Key] = k
		}
	}
	for query != "" {
		var key string
		key, query, _ = strings.Cut(query, "&")
		if strings.Contains(key, ";") {
			err = fmt.Errorf("invalid semicolon separator in query")
			continue
		}
		if key == "" {
			continue
		}
		key, value, _ := strings.Cut(key, "=")
		key, err1 := url.QueryUnescape(key)
		if err1 != nil {
			if err == nil {
				err = err1
			}
			continue
		}
		value, err1 = url.QueryUnescape(value)
		if err1 != nil {
			if err == nil {
				err = err1
			}
			continue
		}
		index, ok := indexes[key]
		if ok {
			(*s)[index].Values = append((*s)[index].Values, value)
		} else {
			indexes[key] = len(*s)
			*s = append(*s, &URLValues{
				Key:    key,
				Values: []string{value},
			})
		}
	}
	return err
}

func (s SortedURLValues) ApplyCond(cond *db.Compounds) {
	for _, v := range s {
		cond.AddKV(v.Key, v.Values[0])
	}
}

func (s SortedURLValues) Get(key string) string {
	for _, v := range s {
		if v.Key == key {
			if len(v.Values) == 0 {
				return ""
			}
			return v.Values[0]
		}
	}
	return ""
}

func (s *SortedURLValues) Set(key, value string) {
	for _, v := range *s {
		if v.Key == key {
			v.Values = []string{value}
			return
		}
	}
	*s = append(*s, &URLValues{
		Key:    key,
		Values: []string{value},
	})
}

func (s *SortedURLValues) Add(key, value string) {
	for _, v := range *s {
		if v.Key == key {
			v.Values = append(v.Values, value)
			return
		}
	}
	*s = append(*s, &URLValues{
		Key:    key,
		Values: []string{value},
	})
}

func (s *SortedURLValues) Del(key string) {
	delIndex := -1
	for i, v := range *s {
		if v.Key == key {
			delIndex = i
			break
		}
	}
	if delIndex > -1 {
		switch delIndex {
		case 0:
			if len(*s) > 1 {
				*s = (*s)[1:]
			} else {
				*s = (*s)[0:0]
			}
		case len(*s) - 1:
			*s = (*s)[0:delIndex]
		default:
			*s = append((*s)[0:delIndex], (*s)[delIndex+1:]...)
		}
	}
}

func (s SortedURLValues) Has(key string) bool {
	for _, v := range s {
		if v.Key == key {
			return true
		}
	}
	return false
}
