package param

import (
	"strings"

	"github.com/webx-top/com"
)

var DefaultTransferFuncs = map[string]func(value interface{}, row Store) interface{}{}

func RegisterTransferFunc(name string, fn func(value interface{}, row Store) interface{}) {
	DefaultTransferFuncs[name] = fn
}

func TransformCamelCase(keys ...string) Transfers {
	return TransformCase(com.CamelCase, keys...)
}

func TransformSnakeCase(keys ...string) Transfers {
	return TransformCase(com.SnakeCase, keys...)
}

func TransformLowerCaseFirst(keys ...string) Transfers {
	return TransformCase(com.LowerCaseFirst, keys...)
}

func TransformCase(keyTransfer func(string) string, keys ...string) Transfers {
	transfers := Transfers{}
	for _, k := range keys {
		ps := strings.SplitN(k, `:`, 2)
		k = ps[0]
		tr := NewTransform().SetKey(keyTransfer(k))
		if len(ps) == 2 {
			if fn, ok := DefaultTransferFuncs[ps[1]]; ok {
				tr.SetFunc(fn)
			}
		}
		transfers[k] = tr
	}
	return transfers
}
