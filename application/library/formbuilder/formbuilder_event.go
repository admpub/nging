package formbuilder

import (
	"strings"

	"github.com/webx-top/echo"
)

// On 注册事件
func (f *FormBuilder) On(method string, funcs ...MethodHook) *FormBuilder {
	method = strings.ToUpper(method)
	f.on.On(method, funcs...)
	return f
}

// - 常用事件注册快捷函数

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

// - 不常用事件注册快捷函数

func (f *FormBuilder) OnConnect(funcs ...MethodHook) *FormBuilder {
	f.on.On(echo.CONNECT, funcs...)
	return f
}

func (f *FormBuilder) OnHead(funcs ...MethodHook) *FormBuilder {
	f.on.On(echo.HEAD, funcs...)
	return f
}

func (f *FormBuilder) OnOptions(funcs ...MethodHook) *FormBuilder {
	f.on.On(echo.OPTIONS, funcs...)
	return f
}

func (f *FormBuilder) OnPatch(funcs ...MethodHook) *FormBuilder {
	f.on.On(echo.PATCH, funcs...)
	return f
}

func (f *FormBuilder) OnTrace(funcs ...MethodHook) *FormBuilder {
	f.on.On(echo.TRACE, funcs...)
	return f
}

func (f *FormBuilder) OnGlobal(funcs ...MethodHook) *FormBuilder {
	f.on.On(`*`, funcs...)
	return f
}

// Off 撤销事件注册
func (f *FormBuilder) Off(methods ...string) *FormBuilder {
	upperedMethods := make([]string, len(methods))
	for index, method := range methods {
		upperedMethods[index] = strings.ToUpper(method)
	}
	f.on.Off(upperedMethods...)
	return f
}

// Off 撤销所有事件注册
func (f *FormBuilder) OffAll() *FormBuilder {
	f.on.OffAll()
	return f
}
