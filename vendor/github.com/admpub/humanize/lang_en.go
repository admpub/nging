package humanize

var lang_en = LanguageProvider{
	Times: Times{
		Ranges: []TimeRanges{
			{Minute, 1, []TimeRange{
				{2, "1 second"},
				{60, "%d seconds"},
			}},
			{Hour, Minute, []TimeRange{
				{2 * Minute, "1 minute"},
				{Hour, "%d minutes"},
			}},
			{Day, Hour, []TimeRange{
				{2 * Hour, "1 hour"},
				{Day, "%d hours"},
			}},
			{Week, Day, []TimeRange{
				{2 * Day, "1 day"},
				{Week, "%d days"},
			}},
			{Month, Week, []TimeRange{
				{2 * Week, "1 week"},
				{Month, "%d weeks"},
			}},
			{Year, Month, []TimeRange{
				{2 * Month, "1 month"},
				{Year, "%d months"},
			}},
			{LongTime, Year, []TimeRange{
				{2 * Year, "1 year"},
				{LongTime, "%d years"},
			}},
		},
		Future:       "in %s",
		Past:         "%s ago",
		Now:          "now",
		RemainderSep: " and ",
		Units:        DefaultTimeUnits,
	},
}
