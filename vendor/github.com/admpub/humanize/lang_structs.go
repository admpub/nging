package humanize

import "strings"

// Language definition structures.

// List all the existing language providers here.
var languages = map[string]LanguageProvider{
	"pl":    lang_pl,
	"en":    lang_en,
	"zh-cn": lang_zh_cn,
}

func Register(language string, provider LanguageProvider) {
	language = strings.ToLower(language)
	languages[language] = provider
}

func HasLanguage(language string) bool {
	language = strings.ToLower(language)
	_, ok := languages[language]
	return ok
}

func AllLanguages() map[string]LanguageProvider {
	return languages
}

// LanguageProvider is a struct defining all the needed language elements.
type LanguageProvider struct {
	Times Times
}

var (
	DefaultLanguage  = `en`
	DefaultTimeUnits = TimeUnits{
		"second": 1,
		"minute": Minute,
		"hour":   Hour,
		"day":    Day,
		"week":   Week,
		"month":  Month,
		"year":   Year,
	}
)

// Times Time related language elements.
type Times struct {
	// Time ranges to humanize time.
	Ranges []TimeRanges
	// String for formatting time in the future.
	Future string
	// String for formatting time in the past.
	Past string
	// String to humanize now.
	Now string
	// Remainder separator
	RemainderSep string
	// Unit values for matching the input. Partial matches are ok.
	Units TimeUnits
}

// TimeUnits Time unit definitions for input parsing. Use partial matches.
type TimeUnits map[string]int64

// TimeRanges Definition of time ranges to match against.
type TimeRanges struct {
	UpperLimit int64
	DivideBy   int64
	Ranges     []TimeRange
}
type TimeRange struct {
	UpperLimit int64
	Format     string
}
