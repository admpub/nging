package timeago

func getEnglish() map[string]string {
	return map[string]string{
		"ago":    "ago",
		"online": "Online",
		"now":    "Now",
		// Seconds
		"second":   "second",
		"seconds":  "seconds",
		"seconds2": "seconds",
		// Minutes
		"minute":   "minute",
		"minutes":  "minutes",
		"minutes2": "minutes",
		// Hours
		"hour":   "hour",
		"hours":  "hours",
		"hours2": "hours",
		// Days
		"day":   "day",
		"days":  "days",
		"days2": "days",
		// Weeks
		"week":   "week",
		"weeks":  "weeks",
		"weeks2": "weeks",
		// Months
		"month":   "month",
		"months":  "months",
		"months2": "months",
		// Years
		"year":   "year",
		"years":  "years",
		"years2": "years",
	}
}

func getEnglishRule() Rule {
	return Rule{
		Single: func(number int64, lastDigit int) bool { return number == 1 },
		Plural: func(number int64, lastDigit int) bool { return number > 1 || number == 0 },
	}
}
