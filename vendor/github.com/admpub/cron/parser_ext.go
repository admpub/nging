package cron

var withSecondParser = NewParser(
	Second | Minute | Hour | Dom | Month | Dow | Descriptor,
)

func Parse(spec string) (Schedule, error) {
	return withSecondParser.Parse(spec)
}
