package storer

import (
	"path/filepath"
	
	"github.com/webx-top/echo"
	"github.com/webx-top/image"
)

var (
	Default = Info{Name: `local`}
	DefaultWatermarkOptions = &image.WatermarkOptions{
		Watermark:filepath.Join(echo.Wd(), `public/assets/backend/images/nging-gear.png`),
		Type:"image",
		On:true,
	}
)

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

func GetWatermarkOptions() *image.WatermarkOptions {
	options, ok := echo.GetStore(`NgingConfig`).Store(`base`).Get(`watermark`).(*image.WatermarkOptions)
	if ok {
		return options
	}
	return DefaultWatermarkOptions
}
