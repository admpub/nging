package i18n

func (f *TranslatorFactory) Reload(localeCode string) (t *Translator, errors []error) {
	return f.load(localeCode)
}
