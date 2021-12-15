package humanize

var lang_pl = LanguageProvider{
	Times: Times{
		Ranges: []TimeRanges{
			{Minute, 1, []TimeRange{
				{2, "1 sekundę"},
				{5, "%d sekundy"},
				{Minute, "%d sekund"},
			}},
			{Hour, Minute, []TimeRange{
				{2 * Minute, "minutę"},
				{5 * Minute, "%d minuty"},
				{Hour, "%d minut"},
			}},
			{Day, Hour, []TimeRange{
				{2 * Hour, "1 godzinę"},
				{5 * Hour, "%d godziny"},
				{Day, "%d godzin"},
			}},
			{Week, Day, []TimeRange{
				{2 * Day, "1 dzień"},
				{Week, "%d dni"},
			}},
			{Month, Week, []TimeRange{
				{2 * Week, "1 tydzień"},
				{5 * Week, "%d tygodnie"},
				{Month, "%d tygodni"},
			}},
			{Year, Month, []TimeRange{
				{2 * Month, "1 miesiąc"},
				{5 * Month, "%d miesiące"},
				{Year, "%d miesięcy"},
			}},
			{LongTime, Year, []TimeRange{
				{2 * Year, "1 rok"},
				{5 * Year, "%d lata"},
				{LongTime, "%d lat"},
			}},
		},
		Future:       "za %s",
		Past:         "%s temu",
		Now:          "teraz",
		RemainderSep: " i ",
		Units:        DefaultTimeUnits,
	},
}
