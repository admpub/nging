package formbuilder

import (
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/formfilter"
)

type MethodHook func(echo.Context, *FormBuilder) error
type MethodHooks map[string][]MethodHook

func (hooks MethodHooks) On(method string, fn MethodHook) {
	if _, ok := hooks[method]; !ok {
		hooks[method] = []MethodHook{}
	}
	hooks[method] = append(hooks[method], fn)
}

func (hooks MethodHooks) Off(method string) {
	if _, ok := hooks[method]; !ok {
		delete(hooks, method)
	}
}

func (hooks MethodHooks) Fire(method string, ctx echo.Context, f *FormBuilder) error {
	funcs, ok := hooks[method]
	if ok {
		for _, fn := range funcs {
			if err := fn(f.ctx, f); err != nil {
				return err
			}
		}
	}
	return nil
}

func BindModel(ctx echo.Context, form *FormBuilder) error {
	opts := []formfilter.Options{formfilter.Include(form.Config().GetNames()...)}
	if customs, ok := ctx.Internal().Get(`formfilter.Options`).([]formfilter.Options); ok {
		opts = append(opts, customs...)
	}
	return ctx.MustBind(form.Model, formfilter.Build(opts...))
}

func ValidModel(ctx echo.Context, form *FormBuilder) error {
	form.ValidFromConfig()
	valid := form.Validate()
	var err error
	if valid.HasError() {
		err = valid.Errors[0]
		ctx.Data().SetInfo(valid.Errors[0].Message, 0).SetZone(valid.Errors[0].Field)
	}
	return err
}
