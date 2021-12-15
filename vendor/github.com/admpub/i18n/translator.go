package i18n

import "strings"

// Translator is a struct which contains all the rules and messages necessary
// to do internationalization for a specific locale. Most functionality in this
// package is accessed through a Translator instance.
type Translator struct {
	messages map[string]string
	locale   string
	rules    *TranslatorRules
	fallback *Translator
}

func (f *Translator) Messages() map[string]string {
	return f.messages
}

func (f *Translator) Locale() string {
	return f.locale
}

func (f *Translator) Fallback() *Translator {
	return f.fallback
}

// Translate returns the translated message, performang any substitutions
// requested in the substitutions map. If neither this translator nor its
// fallback translator (or the fallback's fallback and so on) have a translation
// for the requested key, and empty string and an error will be returned.
func (t *Translator) Translate(key string, substitutions map[string]string) (translation string, errors []error) {
	if _, ok := t.messages[key]; !ok {
		if t.fallback != nil && t.fallback != t {
			return t.fallback.Translate(key, substitutions)
		}

		errors = append(errors, translatorError{translator: t, message: "key not found: " + key})
		return
	}

	translation, errors = t.substitute(t.messages[key], substitutions)
	return
}

// Pluralize returns the translation for a message containing a plural. The
// plural form used is based on the number float64 and the number displayed in
// the translated string is the numberStr string. If neither this translator nor
// its fallback translator (or the fallback's fallback and so on) have a
// translation for the requested key, and empty string and an error will be
// returned.
func (t *Translator) Pluralize(key string, number float64, numberStr string) (translation string, errors []error) {

	// TODO: errors are returned when there isn't a substitution - but it is
	// valid to not have a substitution in cases where there's only one number
	// for a single plural form. In these cases, no error should be returned.

	if _, ok := t.messages[key]; !ok {
		if t.fallback != nil && t.fallback != t {
			return t.fallback.Pluralize(key, number, numberStr)
		}

		errors = append(errors, translatorError{translator: t, message: "key not found: " + key})
		return
	}

	form := (t.rules.PluralRuleFunc)(number)

	parts := strings.Split(t.messages[key], "|")

	if form > len(parts)-1 {
		errors = append(errors, translatorError{translator: t, message: "too few plural variations: " + key})
		form = len(parts) - 1
	}

	var errs []error
	translation, errs = t.substitute(parts[form], map[string]string{"n": numberStr})
	errors = append(errors, errs...)
	return
}

// Rules Translate returns the translated message, performang any substitutions
// requested in the substitutions map. If neither this translator nor its
// fallback translator (or the fallback's fallback and so on) have a translation
// for the requested key, and empty string and an error will be returned.
func (t *Translator) Rules() TranslatorRules {

	rules := *t.rules

	return rules
}

// Direction returns the text directionality of the locale's writing system
func (t *Translator) Direction() (direction string) {
	return t.rules.Direction
}

// substitute returns a string copy of the input str string will all keys in the
// substitutions map replaced with their value.
func (t *Translator) substitute(str string, substitutions map[string]string) (substituted string, errors []error) {

	substituted = str

	for find, replace := range substitutions {
		if !strings.Contains(str, "{"+find+"}") {
			errors = append(errors, translatorError{translator: t, message: "substitution not found: " + str + ", " + replace})
		}
		substituted = strings.Replace(substituted, "{"+find+"}", replace, -1)
	}

	return
}
