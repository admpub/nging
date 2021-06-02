package formbuilder

import (
	"strings"

	"github.com/webx-top/echo"
)

func (f *FormBuilder) On(method string, funcs ...MethodHook) *FormBuilder {
	method = strings.ToUpper(method)
	f.on.On(method, funcs...)
	return f
}

func (f *FormBuilder) OnPost(funcs ...MethodHook) *FormBuilder {
	f.on.On(echo.POST, funcs...)
	return f
}

func (f *FormBuilder) OnPut(funcs ...MethodHook) *FormBuilder {
	f.on.On(echo.PUT, funcs...)
	return f
}

func (f *FormBuilder) OnDelete(funcs ...MethodHook) *FormBuilder {
	f.on.On(echo.DELETE, funcs...)
	return f
}

func (f *FormBuilder) OnGet(funcs ...MethodHook) *FormBuilder {
	f.on.On(echo.GET, funcs...)
	return f
}

func (f *FormBuilder) Off(methods ...string) *FormBuilder {
	upperedMethods := make([]string, len(methods))
	for index, method := range methods {
		upperedMethods[index] = strings.ToUpper(method)
	}
	f.on.Off(upperedMethods...)
	return f
}

func (f *FormBuilder) OffAll() *FormBuilder {
	f.on.OffAll()
	return f
}
