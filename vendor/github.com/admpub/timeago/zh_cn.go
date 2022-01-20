package timeago

func getZhCN() map[string]string {
	return map[string]string{
		"format": "%s%s%s",
		"ago":    "以前",
		"online": "在线",
		"now":    "刚刚",
		// Seconds
		"second":   "秒",
		"seconds":  "秒",
		"seconds2": "秒",
		// Minutes
		"minute":   "分钟",
		"minutes":  "分钟",
		"minutes2": "分钟",
		// Hours
		"hour":   "小时",
		"hours":  "小时",
		"hours2": "小时",
		// Days
		"day":   "天",
		"days":  "天",
		"days2": "天",
		// Weeks
		"week":   "周",
		"weeks":  "周",
		"weeks2": "周",
		// Months
		"month":   "个月",
		"months":  "个月",
		"months2": "个月",
		// Years
		"year":   "年",
		"years":  "年",
		"years2": "年",
	}
}

func getZhCNRule() Rule {
	return Rule{}
}
