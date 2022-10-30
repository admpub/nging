package common

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/webx-top/db"
	"github.com/webx-top/echo"
)

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
