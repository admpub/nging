package storer

import (
	"github.com/webx-top/echo"
)

var Default = Info{Name: `local`}

func GetOk() (*Info, bool) {
	storerConfig, ok := echo.Get(StorerInfoKey).(*Info)
	if !ok {
		storerConfig, ok = echo.GetStore(`NgingConfig`).Store(`base`).Get(`storer`).(*Info)
	}
	return storerConfig, ok
}

func Get() Info {
	storerConfig, ok := GetOk()
	if !ok {
		return Default
	}
	return *storerConfig
}