package formbuilder

import (
	"errors"

	"github.com/webx-top/echo/formfilter"
	"github.com/webx-top/validation"
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
		if _, ok := hooks[method]; ok {
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
	if !ok {
		return nil
	}
	var err error
	for _, fn := range funcs {
		if err = fn(); err != nil {
			return err
		}
	}
	return err
}

func BindModel(form *FormBuilder) MethodHook {
	return func() error {
		opts := []formfilter.Options{formfilter.Include(form.Config().GetNames()...)}
		opts = append(opts, form.filters...)
		return form.ctx.MustBind(form.Model, formfilter.Build(opts...))
	}
}

func ValidModel(form *FormBuilder) MethodHook {
	return func() error {
		form.ValidFromConfig()
		err := form.Validate().Error()
		if !errors.Is(err, validation.NoError) {
			form.ctx.Data().SetInfo(err.Message, 0).SetZone(err.Field)
			return err
		}
		return nil
	}
}
