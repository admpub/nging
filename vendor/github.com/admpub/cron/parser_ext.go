package cron

var withSecondParser = NewParser(
	Second | Minute | Hour | Dom | Month | Dow | Descriptor,
)

func Parse(standardSpec string) (Schedule, error) {
	return withSecondParser.Parse(standardSpec)
}

