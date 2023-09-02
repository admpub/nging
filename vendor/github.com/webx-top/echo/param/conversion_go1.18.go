//go:build go1.18

package param

import "reflect"

// Converts 转换 Slice (不支持指针类型元素)
func Converts[V Scalar, T Scalar](p []T, converter ...func(s T) V) []V {
	result := make([]V, len(p))
	if len(p) == 0 {
		return result
	}
	var convert func(s T) V
	if len(converter) > 0 {
		convert = converter[0]
	} else {
		rv := reflect.ValueOf(result[0])
		typeName := rv.Kind().String()
		convert = func(s T) V {
			return AsType(typeName, s).(V)
		}
	}
	for i, s := range p {
		result[i] = convert(s)
	}
	return result
}

func InterfacesTo[T Scalar](p []any, converter ...func(s any) T) []T {
	result := make([]T, len(p))
	if len(p) == 0 {
		return result
	}
	var convert func(s any) T
	if len(converter) > 0 {
		convert = converter[0]
	} else {
		rv := reflect.ValueOf(result[0])
		typeName := rv.Kind().String()
		convert = func(s any) T {
			return AsType(typeName, s).(T)
		}
	}
	for i, s := range p {
		result[i] = convert(s)
	}
	return result
}

func AsInterfaces[T any](p []T, converter ...func(s T) any) []any {
	result := make([]any, len(p))
	if len(p) == 0 {
		return result
	}
	var convert func(s T) any
	if len(converter) > 0 {
		convert = converter[0]
	} else {
		convert = func(s T) any {
			return any(s)
		}
	}
	for i, s := range p {
		result[i] = convert(s)
	}
	return result
}
