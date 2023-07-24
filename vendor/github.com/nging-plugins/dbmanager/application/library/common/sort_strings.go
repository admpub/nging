package common

import (
	"sort"

	"github.com/tidwall/collate"
)

const defaultSortLang = `Chinese_NUM`

var defaultSortCompare = collate.IndexString(defaultSortLang)

func SortStrings(arr []string, lang ...string) {
	var lng string
	if len(lang) > 0 {
		lng = lang[0]
	} else {
		lng = defaultSortLang
	}
	var less func(a, b string) bool
	if lng == defaultSortLang {
		less = defaultSortCompare
	} else {
		less = collate.IndexString(lng)
	}
	sort.SliceStable(arr, func(i, j int) bool {
		return less(arr[i], arr[j])
	})
}
