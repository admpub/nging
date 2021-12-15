package i18n

import (
	"math"
)

// pluralRule is a function that takes a single float64 and returns an int.  Its
// intended use is to return an int index for what plural form to use for the
// given float
type pluralRule func(float64) int

// pluralRules contains the list of all pluralRule functions. The string map
// index is used when loading plural rules from yaml files
var pluralRules = map[string]pluralRule{
	"1":  pluralRule1,
	"2A": pluralRule2A,
	"2B": pluralRule2B,
	"2C": pluralRule2C,
	"2D": pluralRule2D,
	"2E": pluralRule2E,
	"2F": pluralRule2F,
	"3A": pluralRule3A,
	"3B": pluralRule3B,
	"3C": pluralRule3C,
	"3D": pluralRule3D,
	"3E": pluralRule3E,
	"3F": pluralRule3F,
	"3G": pluralRule3G,
	"3H": pluralRule3H,
	"3I": pluralRule3I,
	"4A": pluralRule4A,
	"4B": pluralRule4B,
	"4C": pluralRule4C,
	"4D": pluralRule4D,
	"4E": pluralRule4E,
	"4F": pluralRule4F,
	"5A": pluralRule5A,
	"5B": pluralRule5B,
	"6A": pluralRule6A,
	"6B": pluralRule6B,
}

// isInt checks if a float64 is an integer value
func isInt(n float64) bool {
	return n == float64(int64(n))
}

// pluralRule1:
// Logic for calculating the nth plural for languages with no plurals
//
// Plural Forms Rules Documented here:
// https://developer.mozilla.org/en/docs/Localization_and_Plurals
//
// This Plural Rule contains 1 form:
//     - other:
//         - rule:     everything
//         - examples: 0, 0.5, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, …
//
// Languages:
//     - ay:  Aymara
//     - az:  Azerbaijani
//     - bm:  Bambara
//     - bo:  Tibetan
//     - dz:  Dzongkha
//     - fa:  Persian
//     - id:  Indonesian
//     - ig:  Igbo
//     - ii:  Sichuan Yi
//     - hu:  Hungarian
//     - ja:  Japanese
//     - jbo: Lojban
//     - jv:  Javanese
//     - ka:  Georgian
//     - kde: Makonde
//     - kea: Kabuverdianu
//     - km:  Khmer
//     - kn:  Kannada
//     - ko:  Korean
//     - lo:  Lao
//     - ms:  Malay
//     - my:  Burmese
//     - sah: Sakha
//     - ses: Koyraboro Senni
//     - sg:  Sango
//     - su:  Sundanese
//     - th:  Thai
//     - to:  Tongan
//     - tr:  Turkish
//     - tt:  Tatar
//     - ug:  Uyghur
//     - vi:  Vietnamese
//     - wo:  Wolof
//     - yo:  Yoruba
//     - zh:  Chinese
func pluralRule1(n float64) int {
	return 0
}

// pluralRule2A:
// Logic for calculating the nth plural for Spanish or languages who share the same rules as Spanish
//
// Plural Forms Rules Documented here:
// https://developer.mozilla.org/en/docs/Localization_and_Plurals
//
// This Plural Rule contains 2 forms:
//     - one:
//         - rule:     is 1
//         - examples: 1
//     - other:
//         - rule:     everythng else
//         - examples: 0, 0.5, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, …
//
// Languages:
//     - af:  Afrikaans
//     - an:  Aragonese
//     - asa: Asu
//     - ast: Asturian
//     - bem: Bemba
//     - bez: Bena
//     - bg:  Bulgarian
//     - bn:  Bengali
//     - brx: Bodo
//     - ca:  Catalan
//     - cgg: Chiga
//     - chr: Cherokee
//     - ckb: Sorani Kurdish
//     - da:  Danish
//     - de:  German
//     - doi: Dogri
//     - dv:  Divehi
//     - ee:  Ewe
//     - el:  Greek
//     - en:  English
//     - eo:  Esperanto
//     - es:  Spanish
//     - et:  Estonian
//     - eu:  Basque
//     - fi:  Finnish
//     - fo:  Faroese
//     - fur: Friulian
//     - fy:  Western Frisian
//     - gl:  Galician
//     - gsw: Swiss German
//     - gu:  Gujarati
//     - ha:  Hausa
//     - haw: Hawaiian
//     - hne: Chhattisgarhi
//     - hy:  Armenian
//     - ia:  Interlingua
//     - is:  Icelandic
//     - it:  Italian
//     - jgo: Ngomba
//     - jmc: Machame
//     - kaj: Jju
//     - kcg: Tyap
//     - kk:  Kazakh
//     - kkj: Kako
//     - kl:  Kalaallisut
//     - ks:  Kashmiri
//     - ksb: Shambala
//     - ku:  Kurdish
//     - ky:  Kirghiz
//     - lb:  Luxembourgish
//     - lg:  Ganda
//     - mai: Maithili
//     - mas: Masai
//     - mgo: Meta'
//     - ml:  Malayalam
//     - mn:  Mongolian
//     - mni: Manipuri
//     - mr:  Marathi
//     - nah: Nahuatl
//     - nap: Neapolitan
//     - nb:  Norwegian Bokmål
//     - nd:  North Ndebele
//     - ne:  Nepali
//     - nl:  Dutch
//     - nn:  Norwegian Nynorsk
//     - nnh: Ngiemboon
//     - no:  Norwegian
//     - nr:  South Ndebele
//     - ny:  Nyanja
//     - nyn: Nyankole
//     - om:  Oromo
//     - or:  Oriya
//     - os:  Ossetic
//     - pa:  Punjabi
//     - pap: Papiamento
//     - pms: Piemontese
//     - ps:  Pashto
//     - pt:  Portuguese
//     - rof: Rombo
//     - rm:  Romansh
//     - rw:  Kinyarwanda
//     - rwk: Rwa
//     - saq: Samburu
//     - sat: Santali
//     - sco: Scots
//     - sd:  Sindhi
//     - seh: Sena
//     - si:  Sinhala
//     - sn:  Shona
//     - so:  Somali
//     - son: Songhai
//     - sq:  Albanian
//     - ss:  Swati
//     - ssy: Saho
//     - st:  Southern Sotho
//     - sv:  Swedish
//     - sw:  Swahili
//     - syr: Syriac
//     - ta:  Tamil
//     - te:  Telugu
//     - teo: Teso
//     - tig: Tigre
//     - tk:  Turkmen
//     - tn:  Tswana
//     - ts:  Tsonga
//     - ur:  Urdu
//     - ve:  Venda
//     - vo:  Volapük
//     - vun: Vunjo
//     - wae: Walser
//     - xh:  Xhosa
//     - xog: Soga
//     - zu:  Zulu
func pluralRule2A(n float64) int {

	isInt := isInt(n)
	i := int64(math.Abs(n))

	switch {
	case isInt && i == 1:
		return 0
	}
	return 1
}

// pluralRule2B:
// Logic for calculating the nth plural for Hindi or languages who share the same rules as Hindi
//
// Plural Forms Rules Documented here:
// https://developer.mozilla.org/en/docs/Localization_and_Plurals
//
// This Plural Rule contains 2 forms:
//     - one:
//         - rule:     is 0 or 1
//         - examples: 0, 1
//     - other:
//         - rule:     everythng else
//         - examples: 0.5, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, …
//
// Languages:
//     - ach: Acholi
//     - ak:  Akan
//     - am:  Amharic
//     - arn: Mapudungun
//     - bh:  Bihari
//     - fil: Filipino
//     - guw: Gun
//     - hi:  Hindi
//     - ln:  Lingala
//     - mfe: Mauritian Creole
//     - mg:  Malagasy
//     - mi:  Maori
//     - nso: Northern Sotho
//     - oc:  Occitan
//     - tg:  Tajic
//     - ti:  Tigrinya
//     - tl:  Tagalog
//     - uz:  Uzbek
//     - wa:  Walloon
func pluralRule2B(n float64) int {

	isInt := isInt(n)
	i := int64(math.Abs(n))

	switch {
	case isInt && (i == 0 || i == 1):
		return 0
	}
	return 1
}

// pluralRule2C:
// Logic for calculating the nth plural for French or languages who share the same rules as French
//
// Plural Forms Rules Documented here:
// https://developer.mozilla.org/en/docs/Localization_and_Plurals
//
// This Plural Rule contains 2 forms:
//     - one:
//         - rule:     n within 0..2 and n is not 2
//         - examples: 0, 0.5, 1, 1.5, …
//     - other:
//         - rule:     everythng else
//         - examples: 2, 2.5, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, …
//
// Languages:
//     - ff:  Fulah
//     - fr:  French
//     - kab: Kabyle
func pluralRule2C(n float64) int {
	abs := math.Abs(n)

	switch {
	case abs >= 0 && abs < 2:
		return 0
	}
	return 1
}

// pluralRule2D:
// Logic for calculating the nth plural for Macedonian or languages who share the same rules as
// Macedonian
//
// Plural Forms Rules Documented here:
// https://developer.mozilla.org/en/docs/Localization_and_Plurals
//
// This Plural Rule contains 2 forms:
//     - one:
//         - rule:     n mod 10 is 1 and n is not 11
//         - examples: 1, 21, 31, 41, 51, 61, 71, 81, 91, 101, 111, 121, 131, 141, 151, 161, 171, …
//     - other:
//         - rule:     everythng else
//         - examples: 0, 0.5, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, …
//
// Languages:
//     - mk: Macedonian
func pluralRule2D(n float64) int {

	isInt := isInt(n)
	i := int64(math.Abs(n))

	mod10 := i % 10

	switch {
	case isInt && mod10 == 1 && i != 11:
		return 0
	}

	return 1
}

// pluralRule2E:
// Logic for calculating the nth plural for Central Atlas Tamazight or languages who share the same
// rules as Central Atlas Tamazight
//
// Plural Forms Rules Documented here:
// https://developer.mozilla.org/en/docs/Localization_and_Plurals
//
// This Plural Rule contains 2 forms:
//     - one:
//         - rule:     n in 0..1 or n in 11..99
//         - examples: 0, 1, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, …
//     - other:
//         - rule:     everythng else
//         - examples: 0.5, 1.5, 2, 3, 4, 5, 6, 7, 8, 9, 10, 100, 101, 102, 103, 104, 105, 106, …
//
// Languages:
//     - tzm: Central Atlas Tamazight
func pluralRule2E(n float64) int {

	isInt := isInt(n)
	i := int64(math.Abs(n))

	switch {
	case isInt && (i == 0 || i == 1 || (i >= 11 && i <= 99)):
		return 0
	}

	return 1
}

// pluralRule2F:
// Logic for calculating the nth plural for Manx or languages who share the same rules as Manx
//
// Plural Forms Rules Documented here:
// https://developer.mozilla.org/en/docs/Localization_and_Plurals
//
// This Plural Rule contains 2 forms:
//     - one:
//         - rule:     n mod 10 in 1..2 or n mod 20 is 0
//         - examples: 0, 1, 2, 11, 12, 20, 21, 22, 31, 32, 40, 41, 42, 51, 52, 60, 61, 62, 71, …
//     - other:
//         - rule:     everythng else
//         - examples: 0.5, 1.5, 3, 3.5, 4, 5, 6, 7, 8, 9, 10, 13, 14, 15, 16, 17, 18, 19, 23, 24, …
//
// Languages:
//     - gv: Manx
func pluralRule2F(n float64) int {

	isInt := isInt(n)
	i := int64(math.Abs(n))

	mod10 := i % 10
	mod20 := i % 20

	switch {
	case isInt && (mod10 == 1 || mod10 == 2 || mod20 == 0):
		return 0
	}

	return 1
}

// pluralRule3A:
// Logic for calculating the nth plural for Latvian or languages who share the same rules as Latvian
//
// Plural Forms Rules Documented here:
// https://developer.mozilla.org/en/docs/Localization_and_Plurals
//
// This Plural Rule contains 3 forms:
//     - zero:
//         - rule:     n is 0
//         - examples: 0
//     - one:
//         - rule:     n mod 10 is 1 and n mod 100 is not 11
//         - examples: 1, 21, 31, 41, 51, 61, 71, 81, 91, 101, 121, 131, 141, 151, 161, 171, 181, …
//     - other:
//         - rule:     everythng else
//         - examples: 2, 2.5, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, …
//
// Languages:
//     - lv: Latvian
func pluralRule3A(n float64) int {

	isInt := isInt(n)
	i := int64(math.Abs(n))

	switch {
	case isInt && i == 0:
		return 0
	case isInt && i%10 == 1 && i%100 != 11:
		return 1
	}
	return 2
}

// pluralRule3B:
// Logic for calculating the nth plural for Nama or languages who share the same rules as Nama
//
// Plural Forms Rules Documented here:
// https://developer.mozilla.org/en/docs/Localization_and_Plurals
//
// This Plural Rule contains 3 forms:
//     - one:
//         - rule:     n is 1
//         - examples: 1
//     - two:
//         - rule:     n is 2
//         - examples: 2
//     - other:
//         - rule:     everythng else
//         - examples: 0, 0.5, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, …
//
// Languages:
//     - iu:  Inuktitut
//     - kw:  Cornish
//     - naq: Nama
//     - se:  Northern Sami
//     - sma: Southern Sami
//     - smi: Sami Language
//     - smj: Lule Sami
//     - smn: Inari Sami
//     - sms: Skolt Sami
func pluralRule3B(n float64) int {

	isInt := isInt(n)
	i := int64(math.Abs(n))

	switch {
	case isInt && i == 1:
		return 0
	case isInt && i == 2:
		return 1
	}
	return 2
}

// pluralRule3C:
// Logic for calculating the nth plural for Romanian or languages who share the same rules as
// Romanian
//
// Plural Forms Rules Documented here:
// https://developer.mozilla.org/en/docs/Localization_and_Plurals
//
// This Plural Rule contains 3 forms:
//     - one:
//         - rule:     n is 1
//         - examples: 1
//     - few:
//         - rule:     n is 0 OR n is not 1 AND n mod 100 in 1..19
//         - examples: 0, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 101, …
//     - other:
//         - rule:     everythng else
//         - examples: 0.5, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32, 33, 34, 35, 36, …
//
// Languages:
//     - ro: Romanian
//     - mo: Moldavian
func pluralRule3C(n float64) int {

	isInt := isInt(n)
	i := int64(math.Abs(n))

	switch {
	case isInt && i == 1:
		return 0
	case isInt && (i == 0 || (i%100 >= 1 && i%100 <= 19)):
		return 1
	}
	return 2
}

// pluralRule3D:
// Logic for calculating the nth plural for Lithuanian or languages who share the same rules as
// Lithuanian
//
// Plural Forms Rules Documented here:
// https://developer.mozilla.org/en/docs/Localization_and_Plurals
//
// This Plural Rule contains 3 forms:
//     - one:
//         - rule:     n mod 10 is 1 and n mod 100 not in 11..19
//         - examples: 1, 21, 31, 41, 51, 61, 71, 81, 91, 101, 121, 131, 141, 151, 161 171, 181, …
//     - few:
//         - rule:     n mod 10 in 2..9 and n mod 100 not in 11..19
//         - examples: 2, 3, 4, 5, 6, 7, 8, 9, 22, 23, 24, 25, 26, 27, 28, 29, 32, 33, 34, 35, 36, …
//     - other:
//         - rule:     everythng else
//         - examples: 0, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 30, 40, 50, 60, 70, 80, 90, …
//
// Languages:
//     - lt: Lithuanian
func pluralRule3D(n float64) int {

	isInt := isInt(n)
	i := int64(math.Abs(n))

	mod10 := i % 10
	mod100 := i % 100

	switch {
	// the part in the parentheses should be replaced with just mod100 !- 11, right?  I've seen this
	// same logically expression so many places that I'm doubting my own logical thinking, so I've
	// implemented it like I've seen it implemented elsewhere.
	case isInt && mod10 == 1 && (mod100 < 11 || mod100 > 19):
		return 0
	case isInt && mod10 >= 2 && mod10 <= 9 && (mod100 < 11 || mod100 > 19):
		return 1
	}

	return 2
}

// pluralRule3E:
// Logic for calculating the nth plural for Czech or languages who share the same rules as Czech
//
// Plural Forms Rules Documented here:
// https://developer.mozilla.org/en/docs/Localization_and_Plurals
//
// This Plural Rule contains 3 forms:
//     - one:
//         - rule:     n is 1
//         - examples: 1
//     - few:
//         - rule:     n in 2..4
//         - examples: 2, 3, 4
//     - other:
//         - rule:     everythng else
//         - examples: 0, 0.5, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, …
//
// Languages:
//     - cs: Czech
//     - sk: Slovak
func pluralRule3E(n float64) int {

	isInt := isInt(n)
	i := int64(math.Abs(n))

	switch {
	case isInt && i == 1:
		return 0
	case isInt && i >= 2 && i <= 4:
		return 1
	}

	return 2
}

// pluralRule3F:
// Logic for calculating the nth plural for Langi or languages who share the same rules as Langi
//
// Plural Forms Rules Documented here:
// https://developer.mozilla.org/en/docs/Localization_and_Plurals
//
// This Plural Rule contains 3 forms:
//     - zero:
//         - rule:     n is 0
//         - examples: 0
//     - one:
//         - rule:     n within 0..2 and n is not 0 and n is not 2
//         - examples: 0.5, 1, 1.5, …
//     - other:
//         - rule:     everythng else
//         - examples: 2, 2.5,  3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, …
//
// Languages:
//     - lag: Langi
func pluralRule3F(n float64) int {

	isInt := isInt(n)
	abs := math.Abs(n)
	i := int64(abs)

	switch {
	case isInt && i == 0:
		return 0
	case abs > 0 && abs < 2:
		return 1
	}

	return 2
}

// pluralRule3G:
// Logic for calculating the nth plural for Tachelhit or languages who share the same rules as
// Tachelhit
//
// Plural Forms Rules Documented here:
// https://developer.mozilla.org/en/docs/Localization_and_Plurals
//
// This Plural Rule contains 3 forms:
//     - one:
//         - rule:     n within 0..1
//         - examples: 0, 0.5, 1
//     - few:
//         - rule:     n in 2..10
//         - examples: 2, 3, 4, 5, 6, 7, 8, 9, 10
//     - other:
//         - rule:     everythng else
//         - examples: 1.5, 10.5, 11, 11.5, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, …
//
// Languages:
//     - shi: Tachelhit
func pluralRule3G(n float64) int {

	isInt := isInt(n)
	abs := math.Abs(n)
	i := int64(abs)

	switch {
	case abs >= 0 && abs <= 1:
		return 0
	case isInt && i >= 2 && i <= 10:
		return 1
	}

	return 2
}

// pluralRule3H:
// Logic for calculating the nth plural for Colognian or languages who share the same rules as
// Colognian
//
// Plural Forms Rules Documented here:
// https://developer.mozilla.org/en/docs/Localization_and_Plurals
//
// This Plural Rule contains 3 forms:
//     - zero:
//         - rule:     n is 0
//         - examples: 0
//     - one:
//         - rule:     n is 1
//         - examples: 1
//     - other:
//         - rule:     everythng else
//         - examples: 0.5, 1.5, 2, 2.5, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, …
//
// Languages:
//     - ksh: Colognian
//     - mnk: Mandinka
func pluralRule3H(n float64) int {

	isInt := isInt(n)
	i := int64(math.Abs(n))

	switch {
	case isInt && i == 0:
		return 0
	case isInt && i == 1:
		return 1
	}

	return 2
}

// pluralRule3I:
// Logic for calculating the nth plural for Kashubian or languages who share the same rules as
// Kashubian
//
// Plural Forms Rules Documented here:
// https://developer.mozilla.org/en/docs/Localization_and_Plurals
//
// This Plural Rule contains 3 forms:
//     - one:
//         - rule:     n is 1
//         - examples: 1
//     - few:
//         - rule:     n mod 10 in 2..4 and n mod 100 not in 10..19
//         - examples: 2, 3, 4, 22, 23, 24, 32, 33, 34, 42, 43, 44, 52, 53, 54, 62, 63, 64, 72, …
//     - other:
//         - rule:     everythng else
//         - examples: 0, 0.5, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 25, …
//
// Languages:
//     - csb: Kashubian
func pluralRule3I(n float64) int {

	isInt := isInt(n)
	i := int64(math.Abs(n))

	mod10 := i % 10
	mod100 := i % 100

	switch {
	case isInt && i == 1:
		return 0
	case isInt && (mod10 >= 2 && mod10 <= 4 && (mod100 < 10 || mod100 > 19)):
		return 1
	}

	return 2
}

// pluralRule4A:
// Logic for calculating the nth plural for Hebrew or languages who share the same rules as Hebrew
//
// Plural Forms Rules Documented here:
// https://developer.mozilla.org/en/docs/Localization_and_Plurals
//
// This Plural Rule contains 4 forms:
//     - one:
//         - rule:     n is 1
//         - examples: 1
//     - two:
//         - rule:     n is 2
//         - examples: 2
//     - many:
//         - rule:     n is not 0 AND n mod 10 is 0
//         - examples: 10, 20, 30, 40, 50, 60, 70, 80, 90, 100, 110, 120, 130, 140, 150, 160, 170, …
//     - other:
//         - rule:     everythng else
//         - examples: 0, 0.5, 3, 4, 5, 6, 7, 8, 9, 11, 12, 13, 14, 15, 16, 17, 18, 19, 21, 22, …
//
// Languages:
//     - he: Hebrew
func pluralRule4A(n float64) int {

	isInt := isInt(n)
	i := int64(math.Abs(n))

	switch {
	case isInt && i == 1:
		return 0
	case isInt && i == 2:
		return 1
	case isInt && i != 0 && i%10 == 0:
		return 2
	}
	return 3
}

// pluralRule4B:
// Logic for calculating the nth plural for Russian or languages who share the same rules as Russian
//
// Plural Forms Rules Documented here:
// https://developer.mozilla.org/en/docs/Localization_and_Plurals
//
// This Plural Rule contains 4 forms:
//     - one:
//         - rule:     n mod 10 is 1 and n mod 100 is not 11
//         - examples: 1, 21, 31, 41, 51, 61, 71, 81, 91, 101, 121, 131, 141, 151, 161, 171, 181, …
//     - few:
//         - rule:     n mod 10 in 2..4 and n mod 100 not in 12..14
//         - examples: 2, 3, 4, 22, 23, 24, 32, 33, 34, 42, 43, 44, 52, 53, 54, 62, 63, 64, 72, …
//     - many:
//         - rule:     n mod 10 is 0 or n mod 10 in 5..9 or n mod 100 in 11..14
//         - examples: 0, 5, 6, 7, 8, 9, 11, 12, 13, 14, 25, 26, 27, 28, 29, 35, 36, 37, 38, 39, …
//     - other:
//         - rule:     everythng else
//         - examples: 0.1, 0.2, 0.3, 0.4, 0.5, 0.6, 0.7, 0.8, 0.9, 1.1, 1.2, 1.3, 1.4, 1.5, 1.6, …
//
// Languages:
//     - be: Belarusian
//     - bs: Bosnian
//     - hr: Croatian
//     - ru: Russian
//     - sh: Serbo-Croatian
//     - sr: Serbian
//     - uk: Ukrainian
func pluralRule4B(n float64) int {

	isInt := isInt(n)
	i := int64(math.Abs(n))

	mod10 := i % 10
	mod100 := i % 100

	switch {
	case isInt && mod10 == 1 && mod100 != 11:
		return 0
	case isInt && (mod10 >= 2 && mod10 <= 4) && (mod100 < 12 || mod100 > 14):
		return 1
	case isInt && (mod10 == 0 || (mod10 >= 5 && mod10 <= 9) || (mod100 >= 11 && mod100 <= 14)):
		return 2
	}

	return 3
}

// pluralRule4C:
// Logic for calculating the nth plural for Polish or languages who share the same rules as Polish
//
// Plural Forms Rules Documented here:
// https://developer.mozilla.org/en/docs/Localization_and_Plurals
//
// This Plural Rule contains 4 forms:
//     - one:
//         - rule:     n is 1
//         - examples: 1
//     - few:
//         - rule:     n mod 10 in 2..4 and n mod 100 not in 12..14
//         - examples: 2, 3, 4, 22, 23, 24, 32, 33, 34, 42, 43, 44, 52, 53, 54, 62, 63, 64, 72, …
//     - many:
//         - rule:     n is not 1 and n mod 10 in 0..1 or n mod 10 in 5..9 or n mod 100 in 12..14
//         - examples: 0, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 25, 26, …
//     - other:
//         - rule:     everythng else
//         - examples: 0.1, 0.2, 0.3, 0.4, 0.5, 0.6, 0.7, 0.8, 0.9, 1.1, 1.2, 1.3, 1.4, 1.5, 1.6, …
//
// Languages:
//     - pl: Polish
func pluralRule4C(n float64) int {

	isInt := isInt(n)
	i := int64(math.Abs(n))

	mod10 := i % 10
	mod100 := i % 100

	switch {
	case isInt && i == 1:
		return 0
	case isInt && mod10 >= 2 && mod10 <= 4 && (mod100 < 12 || mod100 > 14):
		return 1
	case isInt && ((mod10 >= 0 && mod10 <= 1) || (mod10 >= 5 && mod10 <= 9) || (mod100 >= 12 && mod100 <= 14)):
		return 2
	}

	return 3
}

// pluralRule4D:
// Logic for calculating the nth plural for Slovenian or languages who share the same rules as
// Slovenian
//
// Plural Forms Rules Documented here:
// https://developer.mozilla.org/en/docs/Localization_and_Plurals
//
// This Plural Rule contains 4 forms:
//     - one:
//         - rule:     n mod 100 is 1
//         - examples: 1, 11, 21, 31, 41, 51, 61, 71, 81, 91, 101, 111, 121, 131, 141, 151, 161, …
//     - two:
//         - rule:     n mod 100 is 2
//         - examples: 2, 12, 22, 32, 42, 52, 62, 72, 82, 92, 102, 112, 122, 132, 142, 152, 162, …
//     - few:
//         - rule:     n mod 100 in 3..4
//         - examples: 3, 4, 13, 14, 23, 24 33, 34, 43, 44, 53, 54, 63, 64, 73, 74, 83, 84, 93, …
//     - other:
//         - rule:     everythng else
//         - examples: 0, 0.5, 5, 6, 7, 8, 9, 10, 15, 16, 17, 18, 19, 20, 25, 26, 27, 28, 29, 30, …
//
// Languages:
//     - dsb: Lower Sorbian
//     - hsb: Upper Sorbian
//     - sl:  Slovenian
//     - wen: Sorbian Language
func pluralRule4D(n float64) int {

	isInt := isInt(n)
	i := int64(math.Abs(n))

	mod100 := i % 100

	switch {
	case isInt && mod100 == 1:
		return 0
	case isInt && mod100 == 2:
		return 1
	case isInt && mod100 >= 3 && mod100 <= 4:
		return 2
	}

	return 3
}

// pluralRule4E:
// Logic for calculating the nth plural for Maltese or languages who share the same rules as Maltese
//
// Plural Forms Rules Documented here:
// https://developer.mozilla.org/en/docs/Localization_and_Plurals
//
// This Plural Rule contains 4 forms:
//     - one:
//         - rule:     n is 1
//         - examples: 1
//     - few:
//         - rule:     n is 0 or n mod 100 in 2..10
//         - examples: 0, 2, 3, 4, 5, 6, 7, 8, 9, 10, 102, 103, 104, 105, 106, 107, 108, 109, 110, …
//     - many:
//         - rule:     n mod 100 in 11..19
//         - examples: 11, 12, 13, 14, 15, 16, 17, 18, 19, 111, 112, 113, 114, 115, 116, 117, 118, …
//     - other:
//         - rule:     everythng else
//         - examples: 0.5, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32, 33, 34, 35, 36, …
//
// Languages:
//     - mt: Maltese
func pluralRule4E(n float64) int {

	isInt := isInt(n)
	i := int64(math.Abs(n))

	mod100 := i % 100

	switch {
	case isInt && i == 1:
		return 0
	case isInt && (i == 0 || (mod100 >= 2 && mod100 <= 10)):
		return 1
	case isInt && mod100 >= 11 && mod100 <= 19:
		return 2
	}

	return 3
}

// pluralRule4F:
// Logic for calculating the nth plural for Scottish Gaelic or languages who share the same rules as
// Scottish Gaelic
//
// Plural Forms Rules Documented here:
// https://developer.mozilla.org/en/docs/Localization_and_Plurals
//
// This Plural Rule contains 4 forms:
//     - one:
//         - rule:     n in 1,11
//         - examples: 1, 11
//     - two:
//         - rule:     n in 2,12
//         - examples: 2, 12
//     - few:
//         - rule:     n in 3..10,13..19
//         - examples: 3, 4, 5, 6, 7, 8, 9, 10, 13, 14, 15, 16, 17, 18, 19
//     - other:
//         - rule:     everythng else
//         - examples: 0, 0.5, 1.5, 2, 2.5, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, …
//
// Languages:
//     - gd: Scottish Gaelic
func pluralRule4F(n float64) int {

	isInt := isInt(n)
	i := int64(math.Abs(n))

	switch {
	case isInt && (i == 1 || i == 11):
		return 0
	case isInt && (i == 2 || i == 12):
		return 1
	case isInt && ((i >= 3 && i <= 10) || (i >= 13 && i <= 19)):
		return 2
	}

	return 3
}

// pluralRule5A:
// Logic for calculating the nth plural for Irish or languages who share the same rules as Irish
//
// Plural Forms Rules Documented here:
// https://developer.mozilla.org/en/docs/Localization_and_Plurals
//
// This Plural Rule contains 5 forms:
//     - one:
//         - rule:     n is 1
//         - examples: 1
//     - two:
//         - rule:     n is 2
//         - examples: 2
//     - few:
//         - rule:     n in 3..6
//         - examples: 3, 4, 5, 6
//     - many:
//         - rule:     n in 7..10
//         - examples: 7, 8, 9, 10
//     - other:
//         - rule:     everythng else
//         - examples: 0, 0.5, 5, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, …
//
// Languages:
//     - ga: Irish
func pluralRule5A(n float64) int {

	isInt := isInt(n)
	i := int64(math.Abs(n))

	switch {
	case isInt && i == 1:
		return 0
	case isInt && i == 2:
		return 1
	case isInt && i >= 3 && i <= 6:
		return 2
	case isInt && i >= 7 && i <= 10:
		return 3
	}
	return 4
}

// pluralRule5B:
// Logic for calculating the nth plural for Breton or languages who share the same rules as Breton
//
// Plural Forms Rules Documented here:
// https://developer.mozilla.org/en/docs/Localization_and_Plurals
//
// This Plural Rule contains 5 forms:
//     - one:
//         - rule:     n mod 10 is 1 and n mod 100 not in 11,71,91
//         - examples: 1, 21, 31, 41, 51, 61, 81, 101, 121, 131, 141, 151, 161, 181, 201, 221, …
//     - two:
//         - rule:     n mod 10 is 2 and n mod 100 not in 12,72,92
//         - examples: 2, 22, 32, 42, 52, 62, 82, 102, 122, 132, 142, 152, 162, 182, 202, 222, …
//     - few:
//         - rule:     n mod 10 in 3..4,9 and n mod 100 not in 10..19,70..79,90..99
//         - examples: 3, 4, 9, 23, 24, 29, 33, 34, 39, 43, 44, 49, 53, 54, 59, 63, 64, 69, 83, …
//     - many:
//         - rule:     n is not 0 and n mod 1000000 is 0
//         - examples: 1000000, 2000000, 3000000, 4000000, 5000000, 6000000, 7000000, 8000000, …
//     - other:
//         - rule:     everythng else
//         - examples: 0, 0.5, 10, 50, 100, 500, 1000, 5000, 10000, 50000, 100000, 500000, …
//
// Languages:
//     - br: Breton
func pluralRule5B(n float64) int {

	isInt := isInt(n)
	i := int64(math.Abs(n))

	mod10 := i % 10
	mod100 := i % 100

	switch {
	case isInt && mod10 == 1 && mod100 != 11 && mod100 != 71 && mod100 != 91:
		return 0
	case isInt && mod10 == 2 && mod100 != 12 && mod100 != 72 && mod100 != 92:
		return 1
	case isInt && (mod10 == 3 || mod10 == 4 || mod10 == 9) && (mod100 < 10 || mod100 > 19) && (mod100 < 70 || mod100 > 79) && (mod100 < 90 || mod100 > 99):
		return 2
	case isInt && i != 0 && i%1000000 == 0:
		return 3
	}

	return 4
}

// pluralRule6A:
// Logic for calculating the nth plural for Arabic or languages who share the same rules as Arabic
//
// Plural Forms Rules Documented here:
// https://developer.mozilla.org/en/docs/Localization_and_Plurals
//
// This Plural Rule contains 6 forms:
//     - zero:
//         - rule:     n is 0
//         - examples: 0
//     - one:
//         - rule:     n is 1
//         - examples: 1
//     - two:
//         - rule:     n is 2
//         - examples: 2
//     - few:
//         - rule:     n mod 100 in 3..10
//         - examples: 3, 4, 5, 6, 7, 8, 9, 10, 103, 104, 105, 106, 107, 108, 109, 110, 203, 204, …
//     - many:
//         - rule:     n mod 100 in 11..99
//         - examples: 11, 12, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, …
//     - other:
//         - rule:     everythng else
//         - examples: 0.5, 1.5, 100, 101, 102, 200, 201, 202, 300, 301, 302, 400, 401, 402, 500, …
//
// Languages:
//     - ar: Arabic
func pluralRule6A(n float64) int {

	isInt := isInt(n)
	i := int64(math.Abs(n))

	switch {
	case isInt && i == 0:
		return 0
	case isInt && i == 1:
		return 1
	case isInt && i == 2:
		return 2
	case isInt && i%100 >= 3 && i%100 <= 10:
		return 3
	case isInt && i%100 >= 11:
		return 4
	}
	return 5
}

// pluralRule6B:
// Logic for calculating the nth plural for Welsh or languages who share the same rules as Welsh
//
// Plural Forms Rules Documented here:
// https://developer.mozilla.org/en/docs/Localization_and_Plurals
//
// This Plural Rule contains 6 forms:
//     - zero:
//         - rule:     n is 0
//         - examples: 0
//     - one:
//         - rule:     n is 1
//         - examples: 1
//     - two:
//         - rule:     n is 2
//         - examples: 2
//     - few:
//         - rule:     n is 3
//         - examples: 3
//     - many:
//         - rule:     n is 6
//         - examples: 6
//     - other:
//         - rule:     everythng else
//         - examples: 0.5, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, …
//
// Languages:
//     - cy: Welsh
func pluralRule6B(n float64) int {

	isInt := isInt(n)
	i := int64(math.Abs(n))

	switch {
	case isInt && i == 0:
		return 0
	case isInt && i == 1:
		return 1
	case isInt && i == 2:
		return 2
	case isInt && i == 3:
		return 3
	case isInt && i == 6:
		return 4
	}

	return 5
}
