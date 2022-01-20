package timeago

func getRussian() map[string]string {
	return map[string]string{
		"ago":    "назад",
		"online": "В сети",
		"now":    "сейчас",
		// Seconds
		"second":   "секунда",
		"seconds":  "секунды",
		"seconds2": "секунд",
		// Minutes
		"minute":   "минута",
		"minutes":  "минуты",
		"minutes2": "минут",
		// Hours
		"hour":   "час",
		"hours":  "часа",
		"hours2": "часов",
		// Days
		"day":   "день",
		"days":  "дня",
		"days2": "дней",
		// Weeks
		"week":   "неделя",
		"weeks":  "недели",
		"weeks2": "недель",
		// Months
		"month":   "месяц",
		"months":  "месяца",
		"months2": "месяцев",
		// Years
		"year":   "год",
		"years":  "года",
		"years2": "лет",
	}
}

func getRussianRule() Rule {
	return Rule{
		Special: func(number int64, lastDigit int) bool {
			return (number >= 5 && number <= 20) || lastDigit == 0 || (lastDigit >= 5 && lastDigit <= 9)
		},
		Single: func(number int64, lastDigit int) bool { return lastDigit == 1 || number == 0 },
		Plural: func(number int64, lastDigit int) bool { return lastDigit >= 2 && lastDigit < 5 },
	}
}
