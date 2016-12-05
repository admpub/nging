// Copyright 2014 Vubeology, Inc. All rights reserved.
// http://vubeology.com
// http://vube.com

/*
Package i18n offers the following basic internationalization functionality:

	- message translation
		- with placeholder support
		- with plural support
	- number formatting
		- with currency support
		- with percentage support
	- locale-aware string sorting

There's more we'd like to add in the future, including:

	- datetime formatting
	- ordinals
	- CLDR xml to yaml rules generation
	- data caching with size limitations
	- nestable message categories
	- small-string number formatting
		- 7m   : about 7 million
		- 1.2k : about 1,200
		- 253  : exactly 253
		- 25b  : about 25 billion
		- etc.
	- out-of-the-box CLDR messages
		- date/time units
		- calendar/month/day names
		- languages
		- geographic region and country names
		- currencies
		- etc.

How the i18n Package Works

In order to interact with this package, you must first get a TranslatorFactory
instace. Through the TranslatorFactory, you can get a Translator instance.
Almost everything in this package is accessed through methods on the Translator
struct.

About the rules and messages paths: This package ships with built-in rules, and
you are welcome to use those directly. However, if there are locales or rules
that are missing from what ships directly with this package, or if you desire to
use different rules than those that ship with this package, then you can specify
additional rules paths. At this time, this package does not ship with built-in
messages, other than a few used for the unit tests. You will need to specify
your own messages path(s). For both rules and messages paths, you can specify
multiple. Paths later in the slice take precedence over packages earlier in the
slice.

For a basic example of getting a TranslatorFactory instance:

	func main() {

		rulesPath := "/usr/local/lib/i18n/locales/rules"
		messagesPath := "/usr/local/lib/i18n/locales/messages"

		f, _ := i18n.NewTranslatorFactory(
			[]string{rulesPath},
			[]string{messagesPath},
			"en",
		)

		tEn, _ := i18n.GetTranslator("en")

		_ = tEn
	}

Simple Message Translation

For simple message translation, use the Translate function, and send an empty
map as the second argument (we'll explain that argument in the next section).

	func main() {

		rulesPath := "/usr/local/lib/i18n/locales/rules"
		messagesPath := "/usr/local/lib/i18n/locales/messages"

		f, _ := i18n.NewTranslatorFactory(
			[]string{rulesPath},
			[]string{messagesPath},
			"en",
		)

		tEn, _ := i18n.GetTranslator("en")

		// WELCOME_MSG => "Welcome!"
		translation, _ := tEn.Translate("WELCOME_MSG", map[string]string{})

		_ = translation
	}

Message Translation with Placeholders

You can also pass placeholder values to the translate function.  That's what the
second argument is for.  In this example, we will inject a username into the
translation.

	func main() {

		rulesPath := "/usr/local/lib/i18n/locales/rules"
		messagesPath := "/usr/local/lib/i18n/locales/messages"

		f, _ := i18n.NewTranslatorFactory(
			[]string{rulesPath},
			[]string{messagesPath},
			"en",
		)

		tEn, _ := i18n.GetTranslator("en")

		// WELCOME_USER => "Welcome, {user}!"
		username := "Mother Goose"
		translation, _ := tEn.Translate("WELCOME_USER", map[string]string{
			"user" : username
		})

		// results in "Welcome, Mother Goose!"

		_ = translation
	}


Plural Message Translation

You can also translate strings with plurals. However, any one message can
contain at most one plural. If you want to translate "I need 5 apples and 3
oranges" you are out of luck.

The Pluralize method takes 3 arguments. The first is the message key - just like
the Translate method. The second argument is a float which is used to determine
which plural form to use. The third is a string representation of the number.
Why two arguments for the number instead of one? This allows you ultimate
flexibility in number formatting to use in the translation while eliminating the
need for string number parsing.

	func main() {

		rulesPath := "/usr/local/lib/i18n/locales/rules"
		messagesPath := "/usr/local/lib/i18n/locales/messages"

		f, _ := i18n.NewTranslatorFactory(
			[]string{rulesPath},
			[]string{messagesPath},
			"en",
		)

		tEn, _ := i18n.GetTranslator("en")

		// DAYS_AGO => "{n} day ago|{n} days ago"
		translation1, _ := tEn.Pluralize("DAYS_AGO", 1, "1")
		translation2, _ := tEn.Pluralize("DAYS_AGO", 2, "two")

		// results in "1 day ago" and "two days ago"

		_ = translation1
		_ = translation2
	}


Number Formatting

You can use the "FomatNumber", "FormatCurrency" and "FormatPercent" methods to
do locale-based number formatting for numbers, currencies and percentages.

	func main() {

		rulesPath := "/usr/local/lib/i18n/locales/rules"
		messagesPath := "/usr/local/lib/i18n/locales/messages"

		f, _ := i18n.NewTranslatorFactory(
			[]string{rulesPath},
			[]string{messagesPath},
			"en",
		)

		tEn, _ := i18n.GetTranslator("en")

		number := float64(1234.5678)
		numberStr := tEn.FormatNumber(number)
		currencyStr := tEn.FormatCurrency(number, "USD")
		percentStr := tEn.FormatPercent(number)

		// results in 1,234.567, $1,234.56, 123,456%

		_ = numberStr
		_ = currencyStr
		_ = percentStr
	}

Alphabetic String Sorting

If you need to sort a list of strings alphabetically, then you should not use
a simple string comparison to do so - this will often result in incorrect
results.  "ȧ" would normally evaluate as greater than "z", which is not correct
in any latin writing system alphabet.  Use can use the Sort method on the
Translator struct to do an alphabetic sorting that is correct for that locale.
Alternatively, you can access the SortUniversal and the SortLocale functions
directly without a Translator instance.  SortUniversal does not take a specific
locale into account when doing the alphabetic sorting, which means it might be
slightly less accurate than the SortLocal function.  However, there are cases
in which the collation rules for a specific locale are unknown, or the sorting
needs to be done in a local-agnostic way.  For these cases, the SortUniversal
function performs a unicode normalization in order to best sort the strings.

In order to be flexible, these functions take a generic interface slice and a
function for retrieving the value on which to perform the sorting.  For example:

	type Food struct {
		Name string
	}

	func main() {

		rulesPath := "/usr/local/lib/i18n/locales/rules"
		messagesPath := "/usr/local/lib/i18n/locales/messages"

		f, _ := i18n.NewTranslatorFactory(
			[]string{rulesPath},
			[]string{messagesPath},
			"en",
		)

		tEn, _ := i18n.GetTranslator("en")

		toSort := []interface{}{
			Food{Name: "apple"},
			Food{Name: "beet"},
			Food{Name: "carrot"},
			Food{Name: "ȧpricot"},
			Food{Name: "ḃanana"},
			Food{Name: "ċlementine"},
		}

		tEn.Sort(toSort1, func(i interface{}) string {
			if food, ok := i.(Food); ok {
				return food.Name
			}
			return ""
		})

		// results in "apple", "ȧpricot", "ḃanana", "beet", "carrot", "ċlementine"

		// Can also do this:
		i18n.SortLocal("en", toSort1, func(i interface{}) string {
			if food, ok := i.(Food); ok {
				return food.Name
			}
			return ""
		})

		// Or this:
		i18n.SortUniversal(toSort1, func(i interface{}) string {
			if food, ok := i.(Food); ok {
				return food.Name
			}
			return ""
		})

		_ = toSort
	}


Fallback Translators

When getting a Translator instance, the TranslatorFactory will automatically
attempt to determine an appropriate fallback Translator for the locale you
specify. For locales with specific "flavors", like "en-au" or "zh-hans", the
"vanilla" version of that locale will be used if it exists. In these cases that
would be "en" and "zh".

When creating a TranslatorFactory instance, you can optionally specify a
final fallback locale. This will be used if it exists.

When determining a fallback, the the factory first checks the less specific
versions of the specified locale, if they exist and will ultimate fallback
to the global fallback if specified.

	func main() {

		rulesPath := "/usr/local/lib/i18n/locales/rules"
		messagesPath := "/usr/local/lib/i18n/locales/messages"

		f, _ := i18n.NewTranslatorFactory(
			[]string{rulesPath},
			[]string{messagesPath},
			"en",
		)

		tEn, _ := i18n.GetTranslator("en") // no fallback

		tPt, _ := i18n.GetTranslator("pt") // fallback is "en"

		tPtBr, _ := i18n.GetTranslator("pt-br") // fallback is "pt"

		_, _, _ = tEn, tPt, tPtBr
	}


Handling Errors

All of the examples above conveniently ignore errors.  We recommend that you
DO handle errors.  The system is designed to give you a valid result if at all
possible, even in errors occur in the process. However, the errors are still
returned and may provide you helpful information you might otherwise miss - like
missing files, file permissions problems, yaml format problems, missing
translations, etc.  We recommend that you do some sort of logging of these
errors.

	func main() {

		f, errs := i18n.NewTranslatorFactory(
			[]string{rulesPath},
			[]string{messagesPath},
			"en",
		)

		for _, err := range errs {
			Log(err)
		}

		tEn, errs := i18n.GetTranslator("en")

		for _, err := range errs {
			Log(err)
		}

		translation1, err := tEn.Translate("WELCOME_MSG", map[string]string{})

		for _, err := range errs {
			Log(err)
		}

		translation2, err := tEn.Pluralize("DAYS_AGO", 1, "1")

		for _, err := range errs {
			Log(err)
		}

		number := float64(1234.5678)
		currencyStr, err := tEn.FormatCurrency(number, "USD")

		if err != nil {
			Log(err)
		}

		_ = translation1
		_ = translation2
		_ = currencyStr
	}

*/
package i18n
