package formbuilder

import (
	"strings"

	"github.com/webx-top/echo"
)

func (f *FormBuilder) On(method string, fn MethodHook) *FormBuilder {
	method = strings.ToUpper(method)
	f.on.On(method, fn)
	return f
}

func (f *FormBuilder) OnPost(fn MethodHook) *FormBuilder {
	f.on.On(echo.POST, fn)
	return f
}

func (f *FormBuilder) OnPut(fn MethodHook) *FormBuilder {
	f.on.On(echo.PUT, fn)
	return f
}

func (f *FormBuilder) OnDelete(fn MethodHook) *FormBuilder {
	f.on.On(echo.DELETE, fn)
	return f
}

func (f *FormBuilder) OnGet(fn MethodHook) *FormBuilder {
	f.on.On(echo.GET, fn)
	return f
}

func (f *FormBuilder) Off(method string) *FormBuilder {
	method = strings.ToUpper(method)
	f.on.Off(method)
	return f
}

func (f *FormBuilder) OffAll() *FormBuilder {
	f.on.OffAll()
	return f
}
