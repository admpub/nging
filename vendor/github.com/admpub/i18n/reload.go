package i18n

import (
	"fmt"
	"strings"
)

func (f *TranslatorFactory) Reload(localeCode string) (t *Translator, errors []error) {

	fallback := f.getFallback(localeCode)

	exists, errs := f.LocaleExists(localeCode)
	if !exists {
		errors = append(errors, translatorError{message: "could not find rules and messages for locale " + localeCode})
	}
	for _, e := range errs {
		errors = append(errors, e)
	}

	rules := new(TranslatorRules)
	files := []string{}

	// TODO: the rules loading logic is fairly complex, and there are some
	// specific cases we are not testing for yet. We need to test that the
	// fallback locale rules do not influence the rules loaded, and that the
	// base rules do.

	// the load the base (default) rule values
	// the step above
	for _, p := range f.rulesPaths {
		p = strings.TrimRight(p, pathSeparator)
		files = append(files, p+pathSeparator+"root.yaml")
	}

	// load less specific fallback locale rules
	parts := strings.Split(localeCode, "-")
	if len(parts) > 1 {
		for i := range parts {
			fb := strings.Join(parts[0:i+1], "-")
			for _, p := range f.rulesPaths {
				p = strings.TrimRight(p, pathSeparator)
				files = append(files, p+pathSeparator+fb+".yaml")
			}
		}
	}

	// finally load files for this specific locale
	for _, p := range f.rulesPaths {
		p = strings.TrimRight(p, pathSeparator)
		files = append(files, p+pathSeparator+localeCode+".yaml")
	}

	errs = rules.load(files, f.getFileSystem)
	for _, err := range errs {
		errors = append(errors, err)
	}

	messages, errs := loadMessages(localeCode, f.messagesPaths, f.getFileSystem)
	for _, err := range errs {
		fmt.Println(err)
		errors = append(errors, err)
	}

	t = new(Translator)
	t.locale = localeCode
	t.messages = messages
	t.fallback = fallback
	t.rules = rules

	f.translators[localeCode] = t

	return
}
