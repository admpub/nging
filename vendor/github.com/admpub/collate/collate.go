package collate

import (
	"sort"
	"strings"

	"golang.org/x/text/collate"
	"golang.org/x/text/language"
)

// SupportedLangs returns all of the languages that Index() supports.
func SupportedLangs() []string {
	var langs []string
	for _, tag := range langMap {
		langs = append(langs, tag.name)
	}
	sort.Strings(langs)
	return langs
}

// IndexString returns a Less function that can be used to compare if
// string "a" is less than string "b".
// The "name" parameter should be a valid collate definition.
//
//	Examples of collation names
//	--------------------------------------------------------------------
//	ENGLISH, EN                -- English
//	AMERICANENGLISH, EN-US     -- English US
//	FRENCH, FR                 -- French
//	CHINESE, ZH                -- Chinese
//	SIMPLIFIEDCHINESE, ZH-HANS -- Simplified Chinese
//	...
//
//	Case insensitive: add the CI tag to the name
//	--------------------------------------------------------------------
//	ENGLISH_CI
//	FR_CI
//	ZH-HANS_CI
//	...
//
//	Case sensitive: add the CS tag to the name
//	--------------------------------------------------------------------
//	ENGLISH_CS
//	FR_CS
//	ZH-HANS_CS
//	...
//
//	For numerics: add the NUM tag to the name
//	Specifies that numbers should sort numerically ("2" < "12")
//	--------------------------------------------------------------------
//	DUTCH_NUM
//	JAPANESE_NUM
//	...
//
//	For loosness: add the LOOSE tag to the name
//	Ignores diacritics, case and weight
//	--------------------------------------------------------------------
//	JA_LOOSE
//	CHINESE_LOOSE
//	...
func IndexString(name string) (less func(a, b string) bool) {
	t, opts := parseCollation(name)
	c := collate.New(t, opts...)
	return func(a, b string) bool {
		return c.CompareString(a, b) == -1
	}
}

func parseCollation(s string) (tag language.Tag, opts []collate.Option) {
	parts := strings.Split(s, "_")
	if lt, ok := langMap[strings.ToLower(parts[0])]; ok {
		tag = lt.tag
	} else {
		tag = language.Make(parts[0])
	}
	if tag == language.Und {
		tag = language.English
	}
	opts = append(opts, collate.OptionsFromTag(tag))
	for i := 1; i < len(parts); i++ {
		switch strings.ToLower(parts[i]) {
		case "ci":
			opts = append(opts, collate.IgnoreCase)
		case "num":
			opts = append(opts, collate.Numeric)
		case "loose":
			opts = append(opts, collate.Loose)
		}
	}
	return tag, opts
}

type tlang struct {
	name string
	tag  language.Tag
}

var langMap = map[string]tlang{
	"afrikaans":            tlang{"Afrikaans", language.Afrikaans},
	"amharic":              tlang{"Amharic", language.Amharic},
	"arabic":               tlang{"Arabic", language.Arabic},
	"modernstandardarabic": tlang{"ModernStandardArabic", language.ModernStandardArabic},
	"azerbaijani":          tlang{"Azerbaijani", language.Azerbaijani},
	"bulgarian":            tlang{"Bulgarian", language.Bulgarian},
	"bengali":              tlang{"Bengali", language.Bengali},
	"catalan":              tlang{"Catalan", language.Catalan},
	"czech":                tlang{"Czech", language.Czech},
	"danish":               tlang{"Danish", language.Danish},
	"german":               tlang{"German", language.German},
	"greek":                tlang{"Greek", language.Greek},
	"english":              tlang{"English", language.English},
	"americanenglish":      tlang{"AmericanEnglish", language.AmericanEnglish},
	"britishenglish":       tlang{"BritishEnglish", language.BritishEnglish},
	"spanish":              tlang{"Spanish", language.Spanish},
	"europeanspanish":      tlang{"EuropeanSpanish", language.EuropeanSpanish},
	"latinamericanspanish": tlang{"LatinAmericanSpanish", language.LatinAmericanSpanish},
	"estonian":             tlang{"Estonian", language.Estonian},
	"persian":              tlang{"Persian", language.Persian},
	"finnish":              tlang{"Finnish", language.Finnish},
	"filipino":             tlang{"Filipino", language.Filipino},
	"french":               tlang{"French", language.French},
	"canadianfrench":       tlang{"CanadianFrench", language.CanadianFrench},
	"gujarati":             tlang{"Gujarati", language.Gujarati},
	"hebrew":               tlang{"Hebrew", language.Hebrew},
	"hindi":                tlang{"Hindi", language.Hindi},
	"croatian":             tlang{"Croatian", language.Croatian},
	"hungarian":            tlang{"Hungarian", language.Hungarian},
	"armenian":             tlang{"Armenian", language.Armenian},
	"indonesian":           tlang{"Indonesian", language.Indonesian},
	"icelandic":            tlang{"Icelandic", language.Icelandic},
	"italian":              tlang{"Italian", language.Italian},
	"japanese":             tlang{"Japanese", language.Japanese},
	"georgian":             tlang{"Georgian", language.Georgian},
	"kazakh":               tlang{"Kazakh", language.Kazakh},
	"khmer":                tlang{"Khmer", language.Khmer},
	"kannada":              tlang{"Kannada", language.Kannada},
	"korean":               tlang{"Korean", language.Korean},
	"kirghiz":              tlang{"Kirghiz", language.Kirghiz},
	"lao":                  tlang{"Lao", language.Lao},
	"lithuanian":           tlang{"Lithuanian", language.Lithuanian},
	"latvian":              tlang{"Latvian", language.Latvian},
	"macedonian":           tlang{"Macedonian", language.Macedonian},
	"malayalam":            tlang{"Malayalam", language.Malayalam},
	"mongolian":            tlang{"Mongolian", language.Mongolian},
	"marathi":              tlang{"Marathi", language.Marathi},
	"malay":                tlang{"Malay", language.Malay},
	"burmese":              tlang{"Burmese", language.Burmese},
	"nepali":               tlang{"Nepali", language.Nepali},
	"dutch":                tlang{"Dutch", language.Dutch},
	"norwegian":            tlang{"Norwegian", language.Norwegian},
	"punjabi":              tlang{"Punjabi", language.Punjabi},
	"polish":               tlang{"Polish", language.Polish},
	"portuguese":           tlang{"Portuguese", language.Portuguese},
	"brazilianportuguese":  tlang{"BrazilianPortuguese", language.BrazilianPortuguese},
	"europeanportuguese":   tlang{"EuropeanPortuguese", language.EuropeanPortuguese},
	"romanian":             tlang{"Romanian", language.Romanian},
	"russian":              tlang{"Russian", language.Russian},
	"sinhala":              tlang{"Sinhala", language.Sinhala},
	"slovak":               tlang{"Slovak", language.Slovak},
	"slovenian":            tlang{"Slovenian", language.Slovenian},
	"albanian":             tlang{"Albanian", language.Albanian},
	"serbian":              tlang{"Serbian", language.Serbian},
	"serbianlatin":         tlang{"SerbianLatin", language.SerbianLatin},
	"swedish":              tlang{"Swedish", language.Swedish},
	"swahili":              tlang{"Swahili", language.Swahili},
	"tamil":                tlang{"Tamil", language.Tamil},
	"telugu":               tlang{"Telugu", language.Telugu},
	"thai":                 tlang{"Thai", language.Thai},
	"turkish":              tlang{"Turkish", language.Turkish},
	"ukrainian":            tlang{"Ukrainian", language.Ukrainian},
	"urdu":                 tlang{"Urdu", language.Urdu},
	"uzbek":                tlang{"Uzbek", language.Uzbek},
	"vietnamese":           tlang{"Vietnamese", language.Vietnamese},
	"chinese":              tlang{"Chinese", language.Chinese},
	"simplifiedchinese":    tlang{"SimplifiedChinese", language.SimplifiedChinese},
	"traditionalchinese":   tlang{"TraditionalChinese", language.TraditionalChinese},
	"zulu":                 tlang{"Zulu", language.Zulu},
}
