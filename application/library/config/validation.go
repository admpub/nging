package config

import (
	"context"
	"regexp"

	"github.com/admpub/log"
	"github.com/admpub/nging/v5/application/library/errorslice"
	validatorPkg "github.com/go-playground/validator/v10"
	"github.com/webx-top/validator"
)

type Validations map[string]*Validation

func (v Validations) Register() (err error) {
	errs := errorslice.New()
	for _k, _v := range v {
		err := _v.Register(_k)
		if err != nil {
			errs.Add(err)
			log.Errorf(`failed to parse validation regexp %q: %v`, _v.Regexp, err)
		}
	}
	return errs
}

type Validation struct {
	Regexp       string            `json:"regexp"`
	Translations map[string]string `json:"translations"`
}

func (v *Validation) Register(name string) error {
	re, err := regexp.Compile(v.Regexp)
	if err != nil {
		return err
	}
	tr := map[string]*validator.Translation{}
	for k, v := range v.Translations {
		tr[k] = &validator.Translation{Text: v}
	}
	validator.RegisterCustomValidation(name, func(_ context.Context, f validatorPkg.FieldLevel) bool {
		return re.MatchString(f.Field().String())
	}, validator.OptTranslations(tr))
	return nil
}
