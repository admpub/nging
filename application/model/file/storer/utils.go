package storer

import (
	"github.com/webx-top/echo"
)

func Get() (*Info, bool) {
	storerConfig, ok := echo.Get(StorerInfoKey).(*Info)
	if !ok {
		storerConfig, ok = echo.GetStore(`NgingConfig`).Store(`base`).Get(`storer`).(*Info)
	}
	return storerConfig, ok
}