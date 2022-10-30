package validator

import (
	"errors"
	"fmt"
	"strings"

	"github.com/admpub/log"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/webx-top/echo"
)

var DefaultLocale = `zh`

func NormalizeLocale(locale string) string {
	return strings.Replace(locale, `-`, `_`, -1)
}

func New(ctx echo.Context, locales ...string) *Validate {
	var locale string
	if len(locales) > 0 {
		locale = locales[0]
	}
	if len(locale) == 0 {
		locale = DefaultLocale
	}
	translator, _ := UniversalTranslator().GetTranslator(locale)
	validate := validator.New()
	transtation, ok := Translations[locale]
	if !ok {
		args := strings.SplitN(locale, `_`, 2)
		if len(args) == 2 {
			transtation, ok = Translations[args[0]]
		} else {
			log.Warnf(`[validator] not found translation: %s`, locale)
		}
	}
	if ok {
		transtation(validate, translator)
	}
	v := &Validate{
		validator:  validate,
		translator: translator,
		context:    ctx,
	}
	regiterCustomValidationTranslator(v, translator, locale)
	return v
}

type Validate struct {
	validator  *validator.Validate
	translator ut.Translator
	context    echo.Context
}

func (v *Validate) Object() *validator.Validate {
	return v.validator
}

func (v *Validate) Context() echo.Context {
	return v.context
}

// ValidateMap validates map data form a map of tags
func (v *Validate) ValidateMap(data map[string]interface{}, rules map[string]interface{}) map[string]interface{} {
	return v.validator.ValidateMapCtx(v.context, data, rules)
}

// Struct 接收的参数为一个struct
func (v *Validate) Struct(i interface{}) error {
	return v.Error(v.validator.StructCtx(v.context, i))
}

// StructExcept 校验struct中的选项，不过除了fields里所给的字段
func (v *Validate) StructExcept(s interface{}, fields ...string) error {
	return v.Error(v.validator.StructExceptCtx(v.context, s))
}

// StructFiltered 接收一个struct和一个函数，这个函数的返回值为bool，决定是否跳过该选项
func (v *Validate) StructFiltered(s interface{}, fn validator.FilterFunc) error {
	return v.Error(v.validator.StructFilteredCtx(v.context, s, fn))
}

// StructPartial 接收一个struct和fields，仅校验在fields里的值
func (v *Validate) StructPartial(s interface{}, fields ...string) error {
	return v.Error(v.validator.StructPartialCtx(v.context, s, fields...))
}

// Var 接收一个变量和一个tag的值，比如 validate.Var(i, "gt=1,lt=10")
func (v *Validate) Var(field interface{}, tag string) error {
	return v.Error(v.validator.VarCtx(v.context, field, tag))
}

// VarWithValue 将两个变量进行对比，比如 validate.VarWithValue(s1, s2, "eqcsfield")
func (v *Validate) VarWithValue(field interface{}, other interface{}, tag string) error {
	return v.Error(v.validator.VarWithValueCtx(v.context, field, other, tag))
}

// Validate 此处支持两种用法：
// 1. Validate(表单字段名, 表单值, 验证规则名)
// 2. Validate(结构体实例, 要验证的结构体字段1，要验证的结构体字段2)
// Validate(结构体实例) 代表验证所有带“valid”标签的字段
func (v *Validate) Validate(i interface{}, args ...interface{}) echo.ValidateResult {
	e := echo.NewValidateResult()
	var err error
	switch m := i.(type) {
	case string:
		field := m
		var value interface{}
		var rule string
		switch len(args) {
		case 2:
			rule = fmt.Sprint(args[1])
			fallthrough
		case 1:
			value = args[0]
		}
		if len(rule) == 0 {
			return e
		}
		err = v.validator.VarCtx(v.context, value, rule)
		if err != nil {
			e.SetField(field)
			e.SetRaw(err)
			return e.SetError(v.Error(err))
		}
	default:
		if len(args) > 0 {
			err = v.validator.StructPartialCtx(v.context, i, echo.InterfacesToStrings(args)...)
		} else {
			err = v.validator.StructCtx(v.context, i)
		}
		if err != nil {
			vErrors := err.(validator.ValidationErrors)
			e.SetField(vErrors[0].Field())
			e.SetRaw(vErrors)
			return e.SetError(v.Error(vErrors[0]))
		}
	}
	return e
}

func (v *Validate) Error(err error) error {
	if err == nil {
		return nil
	}
	switch rErr := err.(type) {
	case validator.FieldError:
		return errors.New(rErr.Translate(v.translator))
	case validator.ValidationErrors:
		return errors.New(rErr[0].Translate(v.translator))
	default:
		return err
	}
}
