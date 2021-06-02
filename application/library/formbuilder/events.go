package formbuilder

import (
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/formfilter"
)

type MethodHook func() error
type MethodHooks map[string][]MethodHook

func (hooks MethodHooks) On(method string, funcs ...MethodHook) {
	if _, ok := hooks[method]; !ok {
		hooks[method] = []MethodHook{}
	}
	hooks[method] = append(hooks[method], funcs...)
}

func (hooks MethodHooks) Off(methods ...string) {
	for _, method := range methods {
		if _, ok := hooks[method]; !ok {
			delete(hooks, method)
		}
	}
}

func (hooks MethodHooks) OffAll() {
	for method := range hooks {
		delete(hooks, method)
	}
}

func (hooks MethodHooks) Fire(method string) error {
	funcs, ok := hooks[method]
	if ok {
		for _, fn := range funcs {
			if err := fn(); err != nil {
				return err
			}
		}
	}
	return nil
}

func BindModel(ctx echo.Context, form *FormBuilder) MethodHook {
	return func() error {
		opts := []formfilter.Options{formfilter.Include(form.Config().GetNames()...)}
		if customs, ok := ctx.Internal().Get(`formfilter.Options`).([]formfilter.Options); ok {
			opts = append(opts, customs...)
		}
		return ctx.MustBind(form.Model, formfilter.Build(opts...))
	}
}

func ValidModel(ctx echo.Context, form *FormBuilder) MethodHook {
	return func() error {
		form.ValidFromConfig()
		valid := form.Validate()
		var err error
		if valid.HasError() {
			err = valid.Errors[0]
			ctx.Data().SetInfo(valid.Errors[0].Message, 0).SetZone(valid.Errors[0].Field)
		}
		return err
	}
}
