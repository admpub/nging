package request

import (
	"context"
	"regexp"

	validatorPkg "github.com/go-playground/validator/v10"
	"github.com/webx-top/com"
	"github.com/webx-top/validation"
	"github.com/webx-top/validator"
)

func init() {
	validator.RegisterCustomValidation(`username`, func(_ context.Context, f validatorPkg.FieldLevel) bool {
		return com.IsUsername(f.Field().String())
	}, validator.OptTranslations(map[string]*validator.Translation{
		`zh`: {Text: `输入的用户名无效 (用户名只能由字母、数字、下划线或汉字组成)`},
		`en`: {Text: `invalid username. {0} can only consist of letters, numbers, underscores or Chinese characters`},
	}))
	validator.RegisterCustomValidation(`alphanum_`, func(_ context.Context, f validatorPkg.FieldLevel) bool {
		return com.IsAlphaNumericUnderscore(f.Field().String())
	}, validator.OptTranslations(map[string]*validator.Translation{
		`zh`: {Text: `{0}的值无效 (只能由字母、数字或下划线组成)`},
		`en`: {Text: `invalid parameter. {0} can only consist of letters, numbers or underscores`},
	}))
	regexpMobile := regexp.MustCompile(validation.DefaultRule.Mobile)
	validator.RegisterCustomValidation(`mobile`, func(_ context.Context, f validatorPkg.FieldLevel) bool {
		return regexpMobile.MatchString(f.Field().String())
	}, validator.OptTranslations(map[string]*validator.Translation{
		`zh`: {Text: `{0}的值无效 (只能由字母、数字或下划线组成)`},
		`en`: {Text: `invalid parameter. {0} can only consist of letters, numbers or underscores`},
	}))
}
