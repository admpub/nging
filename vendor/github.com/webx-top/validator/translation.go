package validator

import (
	"github.com/go-playground/locales"
	locale_en "github.com/go-playground/locales/en"
	locale_enUS "github.com/go-playground/locales/en_US"
	locale_fr "github.com/go-playground/locales/fr"
	locale_jaJP "github.com/go-playground/locales/ja_JP"
	locale_ru "github.com/go-playground/locales/ru"
	locale_zh "github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"

	// RegisterTranslation
	"github.com/go-playground/validator/v10"
	translation_en "github.com/go-playground/validator/v10/translations/en"
	translation_fr "github.com/go-playground/validator/v10/translations/fr"
	translation_ja "github.com/go-playground/validator/v10/translations/ja"
	translation_ru "github.com/go-playground/validator/v10/translations/ru"
	translation_zh "github.com/go-playground/validator/v10/translations/zh"
)

var SupportedLocales []locales.Translator = []locales.Translator{
	locale_zh.New(),
	locale_en.New(),
	locale_enUS.New(),
	locale_fr.New(),
	locale_ru.New(),
	locale_jaJP.New(),
}

type TranslationRegister func(v *validator.Validate, trans ut.Translator) (err error)

var Translations = map[string]TranslationRegister{
	`zh`: translation_zh.RegisterDefaultTranslations,
	`ru`: translation_ru.RegisterDefaultTranslations,
	`ja`: translation_ja.RegisterDefaultTranslations,
	`en`: translation_en.RegisterDefaultTranslations,
	`fr`: translation_fr.RegisterDefaultTranslations,
}

func RegisterTranslation(translator locales.Translator, register TranslationRegister, locales ...string) {
	SupportedLocales = append(SupportedLocales, translator)
	if len(locales) > 0 {
		Translations[locales[0]] = register
	} else {
		Translations[translator.Locale()] = register
	}
}

func UniversalTranslator() *ut.UniversalTranslator {
	fallback := SupportedLocales[0]
	return ut.New(fallback, SupportedLocales...)
}

// RegisterTranslation 添加额外翻译
func (v *Validate) RegisterTranslation(tag string, trans ut.Translator, registerFn validator.RegisterTranslationsFunc, translationFn validator.TranslationFunc) error {
	return v.validator.RegisterTranslation(tag, trans, registerFn, translationFn)
}
