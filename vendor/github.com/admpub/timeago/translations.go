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
	dist   map[string][]string
}

func (t *Translations) init() {
	t.dist = map[string][]string{
		"seconds": {t.T("second"), t.T("seconds"), t.T("seconds2")},
		"minutes": {t.T("minute"), t.T("minutes"), t.T("minutes2")},
		"hours":   {t.T("hour"), t.T("hours"), t.T("hours2")},
		"days":    {t.T("day"), t.T("days"), t.T("days2")},
		"weeks":   {t.T("week"), t.T("weeks"), t.T("weeks2")},
		"months":  {t.T("month"), t.T("months"), t.T("months2")},
		"years":   {t.T("year"), t.T("years"), t.T("years2")},
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

func RegisterTranslations(lang string, trans map[string]string) {
	translations[lang] = NewTranslations(trans)
}

// getTimeTranslations returns array of translations for different
// cases. For example `1 second` must not have `s` at the end
// but `2 seconds` requires `s`. So this method keeps all
// possible options for the translated word.
func getTimeTranslations(lang string) map[string][]string {
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
