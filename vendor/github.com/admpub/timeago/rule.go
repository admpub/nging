package timeago

var rules = map[string]Rule{
	`ru`:    getRussianRule(),
	`en`:    getEnglishRule(),
	`zh-cn`: getZhCNRule(),
}

func RegisterRules(lang string, rule Rule) {
	rules[lang] = rule
}

func getRules(lang string) Rule {
	return rules[lang]
}

type Detector func(number int64, lastDigit int) bool

type Rule struct {
	Single  Detector
	Plural  Detector
	Special Detector
}

func (r Rule) String(number int64, lastDigit int) string {
	switch {
	case r.Special != nil && r.Special(number, lastDigit):
		return "special"
	case r.Single != nil && r.Single(number, lastDigit):
		return "single"
	case r.Plural != nil && r.Plural(number, lastDigit):
		return "plural"
	default:
		return "single"
	}
}
