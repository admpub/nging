package collate

import (
	"sort"
	"sync"
)

const defaultSortLang = `Chinese_NUM`

var (
	defaultSortCompares = map[string]func(a, b string) bool{
		defaultSortLang: defaultSortCompare,
	}
	defaultSortMutex = sync.RWMutex{}
)

func Preinit(langs ...string) {
	defaultSortMutex.Lock()
	for _, lang := range langs {
		defaultSortCompares[lang] = IndexString(lang)
	}
	defaultSortMutex.Unlock()
}

func GetPreinited(lang string) func(a, b string) bool {
	defaultSortMutex.RLock()
	less := defaultSortCompares[lang]
	defaultSortMutex.RUnlock()
	return less
}

var defaultSortCompare = IndexString(defaultSortLang)

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
		less = GetPreinited(lng)
		if less == nil {
			less = IndexString(lng)
			defaultSortMutex.Lock()
			defaultSortCompares[lng] = less
			defaultSortMutex.Unlock()
		}
	}
	sort.SliceStable(arr, func(i, j int) bool {
		return less(arr[i], arr[j])
	})
}
