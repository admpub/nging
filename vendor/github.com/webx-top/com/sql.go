package com

import "regexp"

//HasTableAlias 检查sql语句中是否包含指定表别名
//sqlStr 可以由"<where子句>,<select子句>,<orderBy子句>,<groupBy子句>"组成
func HasTableAlias(alias string, sqlStr string, quotes ...string) (bool, error) {
	var left, right string
	switch len(quotes) {
	case 2:
		right = quotes[1]
		left = quotes[0]
	case 1:
		left = quotes[0]
		right = left
	default:
		left = "`"
		right = "`"
	}
	re, err := regexp.Compile("[ ,][" + left + "]?" + alias + "[" + right + "]?\\.")
	if err != nil {
		return false, err
	}
	return re.MatchString(` ` + sqlStr), err
}
