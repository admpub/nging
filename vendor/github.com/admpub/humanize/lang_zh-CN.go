package humanize

var lang_zh_cn = LanguageProvider{
	Times: Times{
		Ranges: []TimeRanges{
			{Minute, 1, []TimeRange{
				{2, "1秒"},
				{60, "%d秒"},
			}},
			{Hour, Minute, []TimeRange{
				{2 * Minute, "1分钟"},
				{Hour, "%d分钟"},
			}},
			{Day, Hour, []TimeRange{
				{2 * Hour, "1小时"},
				{Day, "%d小时"},
			}},
			{Week, Day, []TimeRange{
				{2 * Day, "1天"},
				{Week, "%d天"},
			}},
			{Month, Week, []TimeRange{
				{2 * Week, "1周"},
				{Month, "%d周"},
			}},
			{Year, Month, []TimeRange{
				{2 * Month, "1个月"},
				{Year, "%d个月"},
			}},
			{LongTime, Year, []TimeRange{
				{2 * Year, "1年"},
				{LongTime, "%d年"},
			}},
		},
		Future:       "%s之后",
		Past:         "%s以前",
		Now:          "刚刚",
		RemainderSep: ", ",
		Units:        DefaultTimeUnits,
	},
}
