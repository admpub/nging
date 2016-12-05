package i18n

import (
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"
)

// numberFormat is a struct that contains all the information about number
// formatting for a specific locale that we need to do number, currency, and
// percentage formatting
type numberFormat struct {
	positivePrefix   string
	positiveSuffix   string
	negativePrefix   string
	negativeSuffix   string
	multiplier       int
	minDecimalDigits int
	maxDecimalDigits int
	minIntegerDigits int
	groupSizeFinal   int // only the right-most (least significant) group
	groupSizeMain    int // all other groups
}

var (
	// numberFormats keeps a copy of all numberFormat instances that have been
	// loaded before, to prevent parsing a single number format string multiple
	// times. There is vey little danger of this list consuming too much memory,
	// since the data for each of these is pretty small in size, and the same
	// formats are used by multiple locales.
	numberFormats           = map[string]*numberFormat{}
	numberFormatsNoDecimals = map[string]*numberFormat{}

	// prefixSuffixRegex is a regular expression that is used to parse number
	// formats
	prefixSuffixRegex = regexp.MustCompile(`(.*?)[#,\.0]+(.*)`)
)

// FormatCurrency takes a float number and a currency key and returns a string
// with a properly formatted currency amount with the correct currency symbol.
// If a symbol cannot be found for the reqested currency, the the key is used
// instead. If the currency key requested is not recognized, it is used as the
// symbol, and an error is returned with the formatted string.
func (t *Translator) FormatCurrency(number float64, currency string) (formatted string, err error) {
	format := t.parseFormat(t.rules.Numbers.Formats.Currency, true)
	result := t.formatNumber(format, number)
	symbol := currency
	if c, ok := t.rules.Currencies[currency]; ok {
		symbol = c.Symbol
	} else {
		err = translatorError{translator: t, message: "unknown currency: " + currency}
	}
	formatted = strings.Replace(result, "¤", symbol, -1)
	return
}

// FormatCurrencyWhole does exactly what FormatCurrency does, but it leaves off
// any decimal places. AKA, it would return $100 rather than $100.00.
func (t *Translator) FormatCurrencyWhole(number float64, currency string) (formatted string, err error) {
	format := t.parseFormat(t.rules.Numbers.Formats.Currency, false)

	result := t.formatNumber(format, number)
	symbol := currency
	if c, ok := t.rules.Currencies[currency]; ok {
		symbol = c.Symbol
	} else {
		err = translatorError{translator: t, message: "unknown currency: " + currency}
	}
	formatted = strings.Replace(result, "¤", symbol, -1)
	return
}

// FormatNumber takes a float number and returns a properly formatted string
// representation of that number according to the locale's number format.
func (t *Translator) FormatNumber(number float64) string {
	return t.formatNumber(t.parseFormat(t.rules.Numbers.Formats.Decimal, true), number)
}

// FormatNumberWhole does exactly what FormatNumber does, but it leaves off any
// decimal places. AKA, it would return 100 rather than 100.01.
func (t *Translator) FormatNumberWhole(number float64) string {
	return t.formatNumber(t.parseFormat(t.rules.Numbers.Formats.Decimal, false), number)
}

// FormatPercent takes a float number and returns a properly formatted string
// representation of that number as a percentage according to the locale's
// percentage format.
func (t *Translator) FormatPercent(number float64) string {
	return t.formatNumber(t.parseFormat(t.rules.Numbers.Formats.Percent, true), number)
}

// parseFormat takes a format string and returns a numberFormat instance
func (t *Translator) parseFormat(pattern string, includeDecimalDigits bool) *numberFormat {

	processed := false
	if includeDecimalDigits {
		_, processed = numberFormats[pattern]
	} else {
		_, processed = numberFormatsNoDecimals[pattern]
	}

	if !processed {

		format := new(numberFormat)
		patterns := strings.Split(pattern, ";")

		matches := prefixSuffixRegex.FindAllStringSubmatch(patterns[0], -1)
		if len(matches) > 0 {
			if len(matches[0]) > 1 {
				format.positivePrefix = matches[0][1]
			}
			if len(matches[0]) > 2 {
				format.positiveSuffix = matches[0][2]
			}
		}

		// default values for negative prefix & suffix
		format.negativePrefix = string(t.rules.Numbers.Symbols.Negative) + string(format.positivePrefix)
		format.negativeSuffix = format.positiveSuffix

		// see if they are in the pattern
		if len(patterns) > 1 {
			matches = prefixSuffixRegex.FindAllStringSubmatch(patterns[1], -1)

			if len(matches) > 0 {
				if len(matches[0]) > 1 {
					format.negativePrefix = matches[0][1]
				}
				if len(matches[0]) > 2 {
					format.negativeSuffix = matches[0][2]
				}
			}
		}

		pat := patterns[0]

		if strings.Index(pat, "%") != -1 {
			format.multiplier = 100
		} else if strings.Index(pat, "‰") != -1 {
			format.multiplier = 1000
		} else {
			format.multiplier = 1
		}

		pos := strings.Index(pat, ".")

		if pos != -1 {
			pos2 := strings.LastIndex(pat, "0")
			if pos2 > pos {
				format.minDecimalDigits = pos2 - pos
			}

			pos3 := strings.LastIndex(pat, "#")
			if pos3 >= pos2 {
				format.maxDecimalDigits = pos3 - pos
			} else {
				format.maxDecimalDigits = format.minDecimalDigits
			}

			pat = pat[0:pos]
		}

		p := strings.Replace(pat, ",", "", -1)
		pos = strings.Index(p, "0")
		if pos != -1 {
			format.minIntegerDigits = strings.LastIndex(p, "0") - pos + 1
		}

		p = strings.Replace(pat, "#", "0", -1)
		pos = strings.LastIndex(pat, ",")
		if pos != -1 {
			format.groupSizeFinal = strings.LastIndex(p, "0") - pos
			pos2 := strings.LastIndex(p[0:pos], ",")
			if pos2 != -1 {
				format.groupSizeMain = pos - pos2 - 1
			} else {
				format.groupSizeMain = format.groupSizeFinal
			}
		}

		if includeDecimalDigits {
			numberFormats[pattern] = format
		} else {
			format.maxDecimalDigits = 0
			format.minDecimalDigits = 0
			numberFormatsNoDecimals[pattern] = format
		}

	}

	if includeDecimalDigits {
		return numberFormats[pattern]
	}

	return numberFormatsNoDecimals[pattern]
}

// formatNumber takes an arbitrary numberFormat and a number and applies that
// format to that number, returning the resulting string
func (t *Translator) formatNumber(format *numberFormat, number float64) string {
	negative := number < 0

	// apply the multiplier first - this is mainly used for percents
	value := math.Abs(number * float64(format.multiplier))
	stringValue := ""

	// get the initial string value, with the maximum # decimal digits
	if format.maxDecimalDigits >= 0 {
		stringValue = numberRound(value, format.maxDecimalDigits)
	} else {
		stringValue = fmt.Sprintf("%f", value)
	}

	// separate the integer from the decimal parts
	pos := strings.Index(stringValue, ".")
	integer := stringValue
	decimal := ""
	if pos != -1 {
		integer = stringValue[:pos]
		decimal = stringValue[pos+1:]
	}

	// make sure the minimum # decimal digits are there
	for len(decimal) < format.minDecimalDigits {
		decimal = decimal + "0"
	}

	// make sure the minimum # integer digits are there
	for len(integer) < format.minIntegerDigits {
		integer = "0" + integer
	}

	// if there's a decimal portion, prepend the decimal point symbol
	if len(decimal) > 0 {
		decimal = string(t.rules.Numbers.Symbols.Decimal) + decimal
	}

	// put the integer portion into properly sized groups
	if format.groupSizeFinal > 0 && len(integer) > format.groupSizeFinal {
		if len(integer) > format.groupSizeMain {
			groupFinal := integer[len(integer)-format.groupSizeFinal:]
			groupFirst := integer[:len(integer)-format.groupSizeFinal]
			integer = strings.Join(chunkString(groupFirst, format.groupSizeMain), t.rules.Numbers.Symbols.Group) + t.rules.Numbers.Symbols.Group + groupFinal
		}
	}

	// append/prepend negative/positive prefix/suffix
	formatted := ""
	if negative {
		formatted = format.negativePrefix + integer + decimal + format.negativeSuffix
	} else {
		formatted = format.positivePrefix + integer + decimal + format.positiveSuffix
	}

	// replace percents and permilles with the local symbols (likely to be exactly the same)
	formatted = strings.Replace(formatted, "%", string(t.rules.Numbers.Symbols.Percent), -1)
	formatted = strings.Replace(formatted, "‰", string(t.rules.Numbers.Symbols.Permille), -1)

	return formatted
}

// chunkString takes a string and chunks it into size-sized pieces in a slice.
// If the length of the string is not divisible by the size, then the first
// chunk in the slice will be padded to compensate.
func chunkString(str string, size int) []string {
	if str == "" {
		return []string{}
	}

	if size == 0 {
		return []string{str}
	}

	chunks := make([]string, int64(math.Ceil(float64(len(str))/float64(size))))

	for len(str) < len(chunks)*size {
		str = " " + str
	}

	for i := 0; i < len(chunks); i++ {
		start := i * size
		stop := int64(math.Min(float64(start+size), float64(len(str))))
		chunks[i] = str[start:stop]
	}

	chunks[0] = strings.TrimLeft(chunks[0], " ")

	return chunks
}

// numberRound takes a number and returns a string containing a rounded to the
// even with the number of decimal places requested.  If this would result in
// the right most decimal place(s) containing "0"s, then all "0"s on the end of
// the decimal portion will be truncated.
func numberRound(number float64, decimals int) string {

	if number == float64(int64(number)) {
		return strconv.FormatInt(int64(number), 10)
	}

	str := fmt.Sprintf("%f", number)
	pos := strings.Index(str, ".")

	if pos != -1 && len(str) > (pos+decimals) {
		str = str[0 : pos+decimals+1]
	}

	backToNum, _ := strconv.ParseFloat(str, 64)
	difference := number - backToNum
	half := 0.5
	for i := 0; i < decimals; i++ {
		half = half / 10
	}

	roundUp := false
	if difference > half {
		roundUp = true
	} else if difference == half {
		// for halfs, round to even
		lastDigit := str[:len(str)-1]
		roundUp = lastDigit == "1" || lastDigit == "3" || lastDigit == "5" || lastDigit == "7" || lastDigit == "9"
	}

	if roundUp {
		// multiply, then ceil, then divide
		multiplier := math.Pow(float64(10), float64(decimals))
		multiplied := strconv.FormatFloat(math.Ceil(number*multiplier), 'f', 0, 64)

		str = multiplied[:len(multiplied)-decimals] + "." + multiplied[len(multiplied)-decimals:]
	}

	str = strings.TrimRight(str, "0")
	str = strings.TrimRight(str, ".")

	return str
}
