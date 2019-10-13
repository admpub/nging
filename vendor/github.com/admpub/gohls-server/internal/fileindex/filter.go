package fileindex

import (
	"os"
	"strings"
)

type Filter func(fi os.FileInfo) bool

var NoFilter = func(fi os.FileInfo) bool {
	return false
}

var HiddenFilter = PrefixFilter(".")

func SuffixFilter(suffixes ...string) Filter {
	return func(fi os.FileInfo) bool {
		for _, suffix := range suffixes {
			if strings.HasSuffix(fi.Name(), suffix) {
				return false
			}
		}
		return true
	}
}

func PrefixFilter(prefixes ...string) Filter {
	return func(fi os.FileInfo) bool {
		for _, prefix := range prefixes {
			if strings.HasPrefix(fi.Name(), prefix) {
				return false
			}
		}
		return true
	}
}

func AndFilter(filters ...Filter) Filter {
	return func(fi os.FileInfo) bool {
		for _, filter := range filters {
			if filter(fi) {
				return true
			}
		}
		return false
	}
}

func OrFilter(filters ...Filter) Filter {
	return func(fi os.FileInfo) bool {
		for _, filter := range filters {
			if filter(fi) {
				return false
			}
		}
		return true
	}
}
