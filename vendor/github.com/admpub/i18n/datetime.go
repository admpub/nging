package i18n

import (
	"fmt"
	"strings"
	"time"
)

// Standard Formats for Dates, Times & DateTimes
// These are the options to pass to the FormatDateTime method.
const (
	DateFormatFull = iota
	DateFormatLong
	DateFormatMedium
	DateFormatShort
	TimeFormatFull
	TimeFormatLong
	TimeFormatMedium
	TimeFormatShort
	DateTimeFormatFull
	DateTimeFormatLong
	DateTimeFormatMedium
	DateTimeFormatShort
)

// Characters with special meaning in a datetime string:
// Technically, all a-z,A-Z characters should be treated as if they represent a
// datetime unit - but not all actually do. Any a-z,A-Z character that is
// intended to be rendered as a literal a-z,A-Z character should be surrounded
// by single quotes. There is currently no support for rendering a single quote
// literal.
const (
	datetimeFormatUnitEra       = 'G'
	datetimeFormatUnitYear      = 'y'
	datetimeFormatUnitMonth     = 'M'
	datetimeFormatUnitDayOfWeek = 'E'
	datetimeFormatUnitDay       = 'd'
	datetimeFormatUnitHour12    = 'h'
	datetimeFormatUnitHour24    = 'H'
	datetimeFormatUnitMinute    = 'm'
	datetimeFormatUnitSecond    = 's'
	datetimeFormatUnitPeriod    = 'a'
	datetimeForamtUnitQuarter   = 'Q'
	datetimeFormatUnitTimeZone1 = 'z'
	datetimeFormatUnitTimeZone2 = 'v'

	datetimeFormatTimeSeparator = ':'
	datetimeFormatLiteral       = '\''
)

// The sequence length of datetime unit characters indicates how they should be
// rendered.
const (
	datetimeFormatLength1Plus       = 1
	datetimeFormatLength2Plus       = 2
	datetimeFormatLengthAbbreviated = 3
	datetimeFormatLengthWide        = 4
	datetimeFormatLengthNarrow      = 5
)

// datetime formats are a sequences off datetime components and string literals
const (
	datetimePatternComponentUnit = iota
	datetimePatternComponentLiteral
)

// A list of currently unsupported units:
// These still need to be implemented. For now they are ignored.
var (
	datetimeFormatUnitCutset = []rune{
		datetimeFormatUnitEra,
		datetimeForamtUnitQuarter,
		datetimeFormatUnitTimeZone1,
		datetimeFormatUnitTimeZone2,
	}
)

type datetimePatternComponent struct {
	pattern       string
	componentType int
}

// FormatDateTime takes a time struct and a format and returns a formatted
// string. Callers should use a DateFormat, TimeFormat, or DateTimeFormat
// constant.
func (t *Translator) FormatDateTime(format int, datetime time.Time) (string, error) {
	pattern := ""
	switch format {
	case DateFormatFull:
		pattern = t.rules.DateTime.Formats.Date.Full
	case DateFormatLong:
		pattern = t.rules.DateTime.Formats.Date.Long
	case DateFormatMedium:
		pattern = t.rules.DateTime.Formats.Date.Medium
	case DateFormatShort:
		pattern = t.rules.DateTime.Formats.Date.Short
	case TimeFormatFull:
		pattern = t.rules.DateTime.Formats.Time.Full
	case TimeFormatLong:
		pattern = t.rules.DateTime.Formats.Time.Long
	case TimeFormatMedium:
		pattern = t.rules.DateTime.Formats.Time.Medium
	case TimeFormatShort:
		pattern = t.rules.DateTime.Formats.Time.Short
	case DateTimeFormatFull:
		datePattern := strings.Trim(t.rules.DateTime.Formats.Date.Full, " ,")
		timePattern := strings.Trim(t.rules.DateTime.Formats.Time.Full, " ,")
		pattern = getDateTimePattern(t.rules.DateTime.Formats.DateTime.Full, datePattern, timePattern)
	case DateTimeFormatLong:
		datePattern := strings.Trim(t.rules.DateTime.Formats.Date.Long, " ,")
		timePattern := strings.Trim(t.rules.DateTime.Formats.Time.Long, " ,")
		pattern = getDateTimePattern(t.rules.DateTime.Formats.DateTime.Long, datePattern, timePattern)
	case DateTimeFormatMedium:
		datePattern := strings.Trim(t.rules.DateTime.Formats.Date.Medium, " ,")
		timePattern := strings.Trim(t.rules.DateTime.Formats.Time.Medium, " ,")
		pattern = getDateTimePattern(t.rules.DateTime.Formats.DateTime.Medium, datePattern, timePattern)
	case DateTimeFormatShort:
		datePattern := strings.Trim(t.rules.DateTime.Formats.Date.Short, " ,")
		timePattern := strings.Trim(t.rules.DateTime.Formats.Time.Short, " ,")
		pattern = getDateTimePattern(t.rules.DateTime.Formats.DateTime.Short, datePattern, timePattern)
	default:
		return "", translatorError{message: "unknown datetime format" + pattern[0:1]}
	}

	parsed, err := t.parseDateTimeFormat(pattern)
	if err != nil {
		return "", err
	}

	return t.formatDateTime(datetime, parsed)
}

// formatDateTime takes a time.Time and a sequence of parsed pattern components
// and returns an internationalized string representation.
func (t *Translator) formatDateTime(datetime time.Time, pattern []*datetimePatternComponent) (string, error) {
	formatted := ""
	for _, component := range pattern {
		if component.componentType == datetimePatternComponentLiteral {
			formatted += component.pattern
		} else {
			f, err := t.formatDateTimeComponent(datetime, component.pattern)
			if err != nil {
				return "", err
			}
			formatted += f
		}
	}

	return strings.Trim(formatted, " ,"), nil
}

// formatDateTimeComponent renders a single component of a datetime format
// pattern.
func (t *Translator) formatDateTimeComponent(datetime time.Time, pattern string) (string, error) {

	switch pattern[0:1] {
	case string(datetimeFormatUnitEra):
		return t.formatDateTimeComponentEra(datetime, len(pattern))
	case string(datetimeFormatUnitYear):
		return t.formatDateTimeComponentYear(datetime, len(pattern))
	case string(datetimeFormatUnitMonth):
		return t.formatDateTimeComponentMonth(datetime, len(pattern))
	case string(datetimeFormatUnitDayOfWeek):
		return t.formatDateTimeComponentDayOfWeek(datetime, len(pattern))
	case string(datetimeFormatUnitDay):
		return t.formatDateTimeComponentDay(datetime, len(pattern))
	case string(datetimeFormatUnitHour12):
		return t.formatDateTimeComponentHour12(datetime, len(pattern))
	case string(datetimeFormatUnitHour24):
		return t.formatDateTimeComponentHour24(datetime, len(pattern))
	case string(datetimeFormatUnitMinute):
		return t.formatDateTimeComponentMinute(datetime, len(pattern))
	case string(datetimeFormatUnitSecond):
		return t.formatDateTimeComponentSecond(datetime, len(pattern))
	case string(datetimeFormatUnitPeriod):
		return t.formatDateTimeComponentPeriod(datetime, len(pattern))
	case string(datetimeForamtUnitQuarter):
		return t.formatDateTimeComponentQuarter(datetime, len(pattern))
	case string(datetimeFormatUnitTimeZone1):
		fallthrough
	case string(datetimeFormatUnitTimeZone2):
		return t.formatDateTimeComponentTimeZone(datetime, len(pattern))
	}

	return "", translatorError{message: "unknown datetime format unit: " + pattern[0:1]}
}

// formatDateTimeComponentEra renders an era component.
// TODO: not yet implemented
func (t *Translator) formatDateTimeComponentEra(datetime time.Time, length int) (string, error) {
	return "", nil
}

// formatDateTimeComponentYear renders a year component.
func (t *Translator) formatDateTimeComponentYear(datetime time.Time, length int) (string, error) {
	year := datetime.Year()
	switch length {
	case datetimeFormatLength1Plus:
		return t.formatDateTimeComponentYearLengthWide(year), nil
	case datetimeFormatLength2Plus:
		return t.formatDateTimeComponentYearLength2Plus(year), nil
	case datetimeFormatLengthWide:
		return t.formatDateTimeComponentYearLengthWide(year), nil
	}

	return "", translatorError{message: fmt.Sprintf("unsupported year length: %d", length)}
}

// formatDateTimeComponentYearLength2Plus renders a 2-digit year component.
func (t *Translator) formatDateTimeComponentYearLength2Plus(year int) string {
	yearShort := year % 100

	if yearShort < 10 {
		return fmt.Sprintf("0%d", yearShort)
	}

	return fmt.Sprintf("%d", yearShort)
}

// formatDateTimeComponentYearLength2Plus renders a full-year component - for
// all modern dates, that's four digits.
func (t *Translator) formatDateTimeComponentYearLengthWide(year int) string {
	return fmt.Sprintf("%d", year)
}

// formatDateTimeComponentMonth renders a month component.
func (t *Translator) formatDateTimeComponentMonth(datetime time.Time, length int) (string, error) {

	month := int(datetime.Month())

	switch length {
	case datetimeFormatLength1Plus:
		return t.formatDateTimeComponentMonth1Plus(month), nil
	case datetimeFormatLength2Plus:
		return t.formatDateTimeComponentMonth2Plus(month), nil
	case datetimeFormatLengthAbbreviated:
		return t.formatDateTimeComponentMonthAbbreviated(month), nil
	case datetimeFormatLengthWide:
		return t.formatDateTimeComponentMonthWide(month), nil
	case datetimeFormatLengthNarrow:
		return t.formatDateTimeComponentMonthNarrow(month), nil
	}

	return "", translatorError{message: fmt.Sprintf("unsupported month length: %d", length)}
}

// formatDateTimeComponentMonth1Plus renders a numeric month component with 1 or
// 2 digits depending on value.
func (t *Translator) formatDateTimeComponentMonth1Plus(month int) string {
	return fmt.Sprintf("%d", month)
}

// formatDateTimeComponentMonth2Plus renders a numeric month component always
// with 2 digits.
func (t *Translator) formatDateTimeComponentMonth2Plus(month int) string {
	if month < 10 {
		return fmt.Sprintf("0%d", month)
	}
	return fmt.Sprintf("%d", month)
}

// formatDateTimeComponentMonthAbbreviated renders an abbreviated text month
// component.
func (t *Translator) formatDateTimeComponentMonthAbbreviated(month int) string {
	switch month {
	case 1:
		return t.rules.DateTime.FormatNames.Months.Abbreviated.Month1
	case 2:
		return t.rules.DateTime.FormatNames.Months.Abbreviated.Month2
	case 3:
		return t.rules.DateTime.FormatNames.Months.Abbreviated.Month3
	case 4:
		return t.rules.DateTime.FormatNames.Months.Abbreviated.Month4
	case 5:
		return t.rules.DateTime.FormatNames.Months.Abbreviated.Month5
	case 6:
		return t.rules.DateTime.FormatNames.Months.Abbreviated.Month6
	case 7:
		return t.rules.DateTime.FormatNames.Months.Abbreviated.Month7
	case 8:
		return t.rules.DateTime.FormatNames.Months.Abbreviated.Month8
	case 9:
		return t.rules.DateTime.FormatNames.Months.Abbreviated.Month9
	case 10:
		return t.rules.DateTime.FormatNames.Months.Abbreviated.Month10
	case 11:
		return t.rules.DateTime.FormatNames.Months.Abbreviated.Month11
	case 12:
		return t.rules.DateTime.FormatNames.Months.Abbreviated.Month12
	}

	return ""
}

// formatDateTimeComponentMonthWide renders a full text month component.
func (t *Translator) formatDateTimeComponentMonthWide(month int) string {
	switch month {
	case 1:
		return t.rules.DateTime.FormatNames.Months.Wide.Month1
	case 2:
		return t.rules.DateTime.FormatNames.Months.Wide.Month2
	case 3:
		return t.rules.DateTime.FormatNames.Months.Wide.Month3
	case 4:
		return t.rules.DateTime.FormatNames.Months.Wide.Month4
	case 5:
		return t.rules.DateTime.FormatNames.Months.Wide.Month5
	case 6:
		return t.rules.DateTime.FormatNames.Months.Wide.Month6
	case 7:
		return t.rules.DateTime.FormatNames.Months.Wide.Month7
	case 8:
		return t.rules.DateTime.FormatNames.Months.Wide.Month8
	case 9:
		return t.rules.DateTime.FormatNames.Months.Wide.Month9
	case 10:
		return t.rules.DateTime.FormatNames.Months.Wide.Month10
	case 11:
		return t.rules.DateTime.FormatNames.Months.Wide.Month11
	case 12:
		return t.rules.DateTime.FormatNames.Months.Wide.Month12
	}

	return ""
}

// formatDateTimeComponentMonthNarrow renders a super-short month compontent -
// not guaranteed to be unique for different months.
func (t *Translator) formatDateTimeComponentMonthNarrow(month int) string {
	switch month {
	case 1:
		return t.rules.DateTime.FormatNames.Months.Narrow.Month1
	case 2:
		return t.rules.DateTime.FormatNames.Months.Narrow.Month2
	case 3:
		return t.rules.DateTime.FormatNames.Months.Narrow.Month3
	case 4:
		return t.rules.DateTime.FormatNames.Months.Narrow.Month4
	case 5:
		return t.rules.DateTime.FormatNames.Months.Narrow.Month5
	case 6:
		return t.rules.DateTime.FormatNames.Months.Narrow.Month6
	case 7:
		return t.rules.DateTime.FormatNames.Months.Narrow.Month7
	case 8:
		return t.rules.DateTime.FormatNames.Months.Narrow.Month8
	case 9:
		return t.rules.DateTime.FormatNames.Months.Narrow.Month9
	case 10:
		return t.rules.DateTime.FormatNames.Months.Narrow.Month10
	case 11:
		return t.rules.DateTime.FormatNames.Months.Narrow.Month11
	case 12:
		return t.rules.DateTime.FormatNames.Months.Narrow.Month12
	}

	return ""
}

// formatDateTimeComponentDayOfWeek renders a day-of-week component.
func (t *Translator) formatDateTimeComponentDayOfWeek(datetime time.Time, length int) (string, error) {
	switch length {
	case datetimeFormatLength1Plus:
		return t.formatDateTimeComponentDayOfWeekWide(datetime.Weekday()), nil
	case datetimeFormatLength2Plus:
		return t.formatDateTimeComponentDayOfWeekShort(datetime.Weekday()), nil
	case datetimeFormatLengthAbbreviated:
		return t.formatDateTimeComponentDayOfWeekAbbreviated(datetime.Weekday()), nil
	case datetimeFormatLengthWide:
		return t.formatDateTimeComponentDayOfWeekWide(datetime.Weekday()), nil
	case datetimeFormatLengthNarrow:
		return t.formatDateTimeComponentDayOfWeekNarrow(datetime.Weekday()), nil
	}

	return "", translatorError{message: fmt.Sprintf("unsupported year day-of-week: %d", length)}
}

// formatDateTimeComponentDayOfWeekAbbreviated renders an abbreviated text
// day-of-week component.
func (t *Translator) formatDateTimeComponentDayOfWeekAbbreviated(dayOfWeek time.Weekday) string {
	switch dayOfWeek {
	case time.Sunday:
		return t.rules.DateTime.FormatNames.Days.Abbreviated.Sun
	case time.Monday:
		return t.rules.DateTime.FormatNames.Days.Abbreviated.Mon
	case time.Tuesday:
		return t.rules.DateTime.FormatNames.Days.Abbreviated.Tue
	case time.Wednesday:
		return t.rules.DateTime.FormatNames.Days.Abbreviated.Wed
	case time.Thursday:
		return t.rules.DateTime.FormatNames.Days.Abbreviated.Thu
	case time.Friday:
		return t.rules.DateTime.FormatNames.Days.Abbreviated.Fri
	case time.Saturday:
		return t.rules.DateTime.FormatNames.Days.Abbreviated.Sat
	}

	return ""
}

// formatDateTimeComponentDayOfWeekAbbreviated renders a
// shorter-then-abbreviated but still unique text day-of-week component.
func (t *Translator) formatDateTimeComponentDayOfWeekShort(dayOfWeek time.Weekday) string {
	switch dayOfWeek {
	case time.Sunday:
		return t.rules.DateTime.FormatNames.Days.Short.Sun
	case time.Monday:
		return t.rules.DateTime.FormatNames.Days.Short.Mon
	case time.Tuesday:
		return t.rules.DateTime.FormatNames.Days.Short.Tue
	case time.Wednesday:
		return t.rules.DateTime.FormatNames.Days.Short.Wed
	case time.Thursday:
		return t.rules.DateTime.FormatNames.Days.Short.Thu
	case time.Friday:
		return t.rules.DateTime.FormatNames.Days.Short.Fri
	case time.Saturday:
		return t.rules.DateTime.FormatNames.Days.Short.Sat
	}

	return ""
}

// formatDateTimeComponentDayOfWeekWide renders a full text day-of-week
// component.
func (t *Translator) formatDateTimeComponentDayOfWeekWide(dayOfWeek time.Weekday) string {
	switch dayOfWeek {
	case time.Sunday:
		return t.rules.DateTime.FormatNames.Days.Wide.Sun
	case time.Monday:
		return t.rules.DateTime.FormatNames.Days.Wide.Mon
	case time.Tuesday:
		return t.rules.DateTime.FormatNames.Days.Wide.Tue
	case time.Wednesday:
		return t.rules.DateTime.FormatNames.Days.Wide.Wed
	case time.Thursday:
		return t.rules.DateTime.FormatNames.Days.Wide.Thu
	case time.Friday:
		return t.rules.DateTime.FormatNames.Days.Wide.Fri
	case time.Saturday:
		return t.rules.DateTime.FormatNames.Days.Wide.Sat
	}

	return ""
}

// formatDateTimeComponentDayOfWeekNarrow renders a super-short day-of-week
// compontent - not guaranteed to be unique for different days.
func (t *Translator) formatDateTimeComponentDayOfWeekNarrow(dayOfWeek time.Weekday) string {
	switch dayOfWeek {
	case time.Sunday:
		return t.rules.DateTime.FormatNames.Days.Narrow.Sun
	case time.Monday:
		return t.rules.DateTime.FormatNames.Days.Narrow.Mon
	case time.Tuesday:
		return t.rules.DateTime.FormatNames.Days.Narrow.Tue
	case time.Wednesday:
		return t.rules.DateTime.FormatNames.Days.Narrow.Wed
	case time.Thursday:
		return t.rules.DateTime.FormatNames.Days.Narrow.Thu
	case time.Friday:
		return t.rules.DateTime.FormatNames.Days.Narrow.Fri
	case time.Saturday:
		return t.rules.DateTime.FormatNames.Days.Narrow.Sat
	}

	return ""
}

// formatDateTimeComponentDay renders a day-of-year component.
func (t *Translator) formatDateTimeComponentDay(datetime time.Time, length int) (string, error) {
	day := datetime.Day()

	switch length {
	case datetimeFormatLength1Plus:
		return fmt.Sprintf("%d", day), nil
	case datetimeFormatLength2Plus:
		if day < 10 {
			return fmt.Sprintf("0%d", day), nil
		}
		return fmt.Sprintf("%d", day), nil
	}

	return "", translatorError{message: fmt.Sprintf("unsupported day-of-year: %d", length)}
}

// formatDateTimeComponentHour12 renders an hour-component using a 12-hour
// clock.
func (t *Translator) formatDateTimeComponentHour12(datetime time.Time, length int) (string, error) {
	hour := datetime.Hour()
	if hour > 12 {
		hour = hour - 12
	}

	switch length {
	case datetimeFormatLength1Plus:
		return fmt.Sprintf("%d", hour), nil
	case datetimeFormatLength2Plus:
		if hour < 10 {
			return fmt.Sprintf("0%d", hour), nil
		}
		return fmt.Sprintf("%d", hour), nil
	}

	return "", translatorError{message: fmt.Sprintf("unsupported hour-12: %d", length)}
}

// formatDateTimeComponentHour24 renders an hour-component using a 24-hour
// clock.
func (t *Translator) formatDateTimeComponentHour24(datetime time.Time, length int) (string, error) {
	hour := datetime.Hour()

	switch length {
	case datetimeFormatLength1Plus:
		return fmt.Sprintf("%d", hour), nil
	case datetimeFormatLength2Plus:
		if hour < 10 {
			return fmt.Sprintf("0%d", hour), nil
		}
		return fmt.Sprintf("%d", hour), nil
	}

	return "", translatorError{message: fmt.Sprintf("unsupported hour-24: %d", length)}
}

// formatDateTimeComponentMinute renders a minute component.
func (t *Translator) formatDateTimeComponentMinute(datetime time.Time, length int) (string, error) {
	minute := datetime.Minute()

	switch length {
	case datetimeFormatLength1Plus:
		return fmt.Sprintf("%d", minute), nil
	case datetimeFormatLength2Plus:
		if minute < 10 {
			return fmt.Sprintf("0%d", minute), nil
		}
		return fmt.Sprintf("%d", minute), nil
	}

	return "", translatorError{message: fmt.Sprintf("unsupported minute: %d", length)}
}

// formatDateTimeComponentSecond renders a second component
func (t *Translator) formatDateTimeComponentSecond(datetime time.Time, length int) (string, error) {
	second := datetime.Second()

	switch length {
	case datetimeFormatLength1Plus:
		return fmt.Sprintf("%d", second), nil
	case datetimeFormatLength2Plus:
		if second < 10 {
			return fmt.Sprintf("0%d", second), nil
		}
		return fmt.Sprintf("%d", second), nil
	}

	return "", translatorError{message: fmt.Sprintf("unsupported second: %d", length)}
}

// formatDateTimeComponentPeriod renders a period component (AM/PM).
func (t *Translator) formatDateTimeComponentPeriod(datetime time.Time, length int) (string, error) {
	hour := datetime.Hour()

	switch length {
	case datetimeFormatLength1Plus:
		return t.formatDateTimeComponentPeriodWide(hour), nil
	case datetimeFormatLengthAbbreviated:
		return t.formatDateTimeComponentPeriodAbbreviated(hour), nil
	case datetimeFormatLengthWide:
		return t.formatDateTimeComponentPeriodWide(hour), nil
	case datetimeFormatLengthNarrow:
		return t.formatDateTimeComponentPeriodNarrow(hour), nil
	}

	return "", translatorError{message: fmt.Sprintf("unsupported day-period: %d", length)}
}

// formatDateTimeComponentPeriodAbbreviated renders an abbreviated period
// component.
func (t *Translator) formatDateTimeComponentPeriodAbbreviated(hour int) string {
	if hour < 12 {
		return t.rules.DateTime.FormatNames.Periods.Abbreviated.AM
	}

	return t.rules.DateTime.FormatNames.Periods.Abbreviated.PM
}

// formatDateTimeComponentPeriodWide renders a full period component.
func (t *Translator) formatDateTimeComponentPeriodWide(hour int) string {
	if hour < 12 {
		return t.rules.DateTime.FormatNames.Periods.Wide.AM
	}

	return t.rules.DateTime.FormatNames.Periods.Wide.PM
}

// formatDateTimeComponentPeriodNarrow renders a super-short period component.
func (t *Translator) formatDateTimeComponentPeriodNarrow(hour int) string {
	if hour < 12 {
		return t.rules.DateTime.FormatNames.Periods.Narrow.AM
	}

	return t.rules.DateTime.FormatNames.Periods.Narrow.PM
}

// formatDateTimeComponentQuarter renders a calendar quarter component - this
// is calendar quarters and not fiscal quarters.
//  - Q1: Jan-Mar
//  - Q2: Apr-Jun
//  - Q3: Jul-Sep
//  - Q4: Oct-Dec
// TODO: not yet implemented
func (t *Translator) formatDateTimeComponentQuarter(datetime time.Time, length int) (string, error) {
	return "", nil
}

// formatDateTimeComponentTimeZone renders a time zone component.
// TODO: this has not yet been implemented
func (t *Translator) formatDateTimeComponentTimeZone(datetime time.Time, length int) (string, error) {
	return "", nil
}

// parseDateTimeFormat takes a format pattern string and returns a sequence of
// components.
func (t *Translator) parseDateTimeFormat(pattern string) ([]*datetimePatternComponent, error) {
	// every thing between single quotes should become a literal
	// all non a-z, A-Z characters become a literal
	// everything else, repeat character sequences become a component
	format := []*datetimePatternComponent{}
	for i := 0; i < len(pattern); {
		char := pattern[i : i+1]

		skip := false
		// for units we don't support yet, just skip over them
		for _, r := range datetimeFormatUnitCutset {
			if char == string(r) {
				skip = true
				break
			}
		}

		if skip {
			i++
			continue
		}

		if char == string(datetimeFormatLiteral) {
			// find the next single quote
			// create a literal out of everything between the quotes
			// and set i to the position after the second quote

			if i == len(pattern)-1 {
				return []*datetimePatternComponent{}, translatorError{message: "malformed datetime format"}
			}

			nextQuote := strings.Index(pattern[i+1:], string(datetimeFormatLiteral))
			if nextQuote == -1 {
				return []*datetimePatternComponent{}, translatorError{message: "malformed datetime format"}
			}

			component := &datetimePatternComponent{
				pattern:       pattern[i+1 : nextQuote+i+1],
				componentType: datetimePatternComponentLiteral,
			}

			format = append(format, component)
			i = nextQuote + i + 2
			continue

		}
		if (char >= "a" && char <= "z") || (char >= "A" && char <= "Z") {
			// this represents a format unit
			// find the entire sequence of the same character
			endChar := lastSequenceIndex(pattern[i:]) + i

			component := &datetimePatternComponent{
				pattern:       pattern[i : endChar+1],
				componentType: datetimePatternComponentUnit,
			}

			format = append(format, component)
			i = endChar + 1
			continue

		}
		if char == string(datetimeFormatTimeSeparator) {
			component := &datetimePatternComponent{
				pattern:       t.rules.DateTime.TimeSeparator,
				componentType: datetimePatternComponentLiteral,
			}
			format = append(format, component)
			i++
			continue

		}

		component := &datetimePatternComponent{
			pattern:       char,
			componentType: datetimePatternComponentLiteral,
		}

		format = append(format, component)
		i++
		continue

	}

	return format, nil
}

// getDateTimePattern combines a date pattern and a time pattern into a datetime
// pattern. The datetimePattern argument includes a {0} placeholder for the time
// pattern, and a {1} placeholder for the date component.
func getDateTimePattern(datetimePattern, datePattern, timePattern string) string {
	return strings.Replace(strings.Replace(datetimePattern, "{1}", datePattern, 1), "{0}", timePattern, 1)
}

// lastSequenceIndex looks at the first character in a string and returns the
// last digits of the first sequence of that character. For example:
//  - ABC: 0
//  - AAB: 1
//  - ABA: 0
//  - AAA: 2
func lastSequenceIndex(str string) int {
	if len(str) == 0 {
		return -1
	}

	if len(str) == 1 {
		return 0
	}

	sequenceChar := str[0:1]
	lastPos := 0
	for i := 1; i < len(str); i++ {
		if str[i:i+1] != sequenceChar {
			break
		}

		lastPos = i
	}

	return lastPos
}
