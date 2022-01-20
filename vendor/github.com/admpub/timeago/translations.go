package timeago

func NewTranslations(source map[string]string) *Translations {
	tr := &Translations{
		source: source,
	}
	tr.init()
	return tr
}

type Translations struct {
	source map[string]string
	dist   map[string]map[string]string
}

func (t *Translations) init() {
	t.dist = map[string]map[string]string{
		"seconds": {"single": t.T("second"), "plural": t.T("seconds"), "special": t.T("seconds2")},
		"minutes": {"single": t.T("minute"), "plural": t.T("minutes"), "special": t.T("minutes2")},
		"hours":   {"single": t.T("hour"), "plural": t.T("hours"), "special": t.T("hours2")},
		"days":    {"single": t.T("day"), "plural": t.T("days"), "special": t.T("days2")},
		"weeks":   {"single": t.T("week"), "plural": t.T("weeks"), "special": t.T("weeks2")},
		"months":  {"single": t.T("month"), "plural": t.T("months"), "special": t.T("months2")},
		"years":   {"single": t.T("year"), "plural": t.T("years"), "special": t.T("years2")},
	}
}

func (t *Translations) T(key string) string {
	if t, ok := t.source[key]; ok {
		return t
	}
	return key
}

var translations = map[string]*Translations{
	`ru`:    NewTranslations(getRussian()),
	`en`:    NewTranslations(getEnglish()),
	`zh-cn`: NewTranslations(getZhCN()),
}

func RegisterTranslations(lang string, trans map[string]string, rule ...Rule) {
	translations[lang] = NewTranslations(trans)
	if len(rule) > 0 {
		RegisterRules(lang, rule[0])
	}
}

// getTimeTranslations returns array of translations for different
// cases. For example `1 second` must not have `s` at the end
// but `2 seconds` requires `s`. So this method keeps all
// possible options for the translated word.
func getTimeTranslations(lang string) map[string]map[string]string {
	t := getTranslations(lang)
	if t == nil {
		return nil
	}
	return t.dist
}

func getTranslations(lang string) *Translations {
	t, ok := translations[lang]
	if ok {
		return t
	}
	if lang == language {
		return nil
	}
	return translations[language]
}

func trans(key string, langs ...string) string {
	lang := language
	if len(langs) > 0 && len(langs[0]) > 0 {
		lang = langs[0]
	}

	if t := getTranslations(lang); t != nil {
		return t.T(key)
	}

	return key
}

func getLanguageForm(num int64, lang string) string {
	lastDigit := getLastNumber(num)
	rule := getRules(lang)
	return rule.String(num, lastDigit)
}
