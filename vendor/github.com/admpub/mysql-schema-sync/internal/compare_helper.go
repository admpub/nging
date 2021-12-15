package internal

import (
	"regexp"
)

var (
	intFindRegexp        = regexp.MustCompile(`(?i) [a-z]*int `)
	intReplaceRegexp     = regexp.MustCompile(`(?i) ([a-z]*int)[^ ]+ `)
	collateFind          = regexp.MustCompile(`(?i) COLLATE `)
	collateReplaceRegexp = regexp.MustCompile(`(?i) CHARACTER SET [^ ]+( COLLATE )`) //CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci

	collateWithCharsetReplaceRegexp = regexp.MustCompile(`(?i)( CHARACTER SET [^ ]+)? COLLATE [^ ]+`)
)

func isSameSchemaItem(src, dest string) bool {
	equal := src == dest
	if !equal {
		// 检查mysql8中版本差异的问题
		if intFindRegexp.MatchString(src) { // src整数无括号
			if !intFindRegexp.MatchString(dest) { // dest整数有括号
				dest = intReplaceRegexp.ReplaceAllString(dest, ` $1 `) // 清除dest括号
			}
		} else if intFindRegexp.MatchString(dest) { // src整数有括号; dest整数无括号
			src = intReplaceRegexp.ReplaceAllString(src, ` $1 `) // 清除src括号
		} else {
			if collateFind.MatchString(src) { // src含COLLATE
				if collateFind.MatchString(dest) { //dest含COLLATE
					dest = collateReplaceRegexp.ReplaceAllString(dest, `$1`)
					src = collateReplaceRegexp.ReplaceAllString(src, `$1`)
				} else { //dest不含COLLATE
					src = collateWithCharsetReplaceRegexp.ReplaceAllString(src, ``)
				}
			} else if collateFind.MatchString(dest) { // src不含COLLATE; dest含COLLATE
				dest = collateWithCharsetReplaceRegexp.ReplaceAllString(dest, ``)
			}
		}
		equal = src == dest
		// if !equal {
		// 	fmt.Println(`src:`, src)
		// 	fmt.Println(`dst:`, dest)
		// 	panic(``)
		// }
	}
	return equal
}
