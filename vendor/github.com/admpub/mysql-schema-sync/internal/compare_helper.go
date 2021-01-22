package internal

import "regexp"

var (
	intFindRegexp        = regexp.MustCompile(`(?i) [a-z]*int `)
	charsetFind          = regexp.MustCompile(`(?i) COLLATE `)
	intReplaceRegexp     = regexp.MustCompile(`(?i) ([a-z]*int)[^ ]+ `)
	charsetReplaceRegexp = regexp.MustCompile(`(?i) CHARACTER SET [^ ]+( COLLATE )`)
)

func isSameSchemaItem(src, dest string) bool {
	equal := src == dest
	if !equal {
		// 检查mysql8中版本差异的问题
		if intFindRegexp.MatchString(src) {
			if !intFindRegexp.MatchString(dest) {
				dest = intReplaceRegexp.ReplaceAllString(dest, ` $1 `)
			}
		} else if intFindRegexp.MatchString(dest) {
			if !intFindRegexp.MatchString(src) {
				src = intReplaceRegexp.ReplaceAllString(src, ` $1 `)
			}
		} else {
			if charsetFind.MatchString(src) {
				src = charsetReplaceRegexp.ReplaceAllString(src, `$1`)
			}
			if charsetFind.MatchString(dest) {
				dest = charsetReplaceRegexp.ReplaceAllString(dest, `$1`)
			}
		}
		equal = src == dest
	}
	return equal
}
