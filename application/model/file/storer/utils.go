package storer

import (
	"path/filepath"

	"github.com/admpub/nging/v4/application/library/common"
	"github.com/webx-top/echo"
	"github.com/webx-top/image"
)

var (
	Default                 = Info{Name: `local`}
	DefaultWatermarkOptions = &image.WatermarkOptions{
		Watermark: filepath.Join(echo.Wd(), `public/assets/backend/images/nging-gear.png`),
		Type:      "image",
		On:        true,
	}
)

func GetOk() (*Info, bool) {
	storerConfig, ok := echo.Get(StorerInfoKey).(*Info)
	if !ok {
		storerConfig, ok = echo.GetStore(common.SettingName).GetStore(`base`).Get(`storer`).(*Info)
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
	options, ok := echo.GetStore(common.SettingName).GetStore(`base`).Get(`watermark`).(*image.WatermarkOptions)
	if ok {
		return options
	}
	return DefaultWatermarkOptions
}

// SaveFilename SaveFilename(`0/`,``,`img.jpg`)
func SaveFilename(subdir, name, postFilename string) (string, error) {
	ext := filepath.Ext(postFilename)
	fname := name
	if len(fname) == 0 {
		var err error
		fname, err = common.UniqueID()
		if err != nil {
			return ``, err
		}
	}
	fname += ext
	return subdir + fname, nil
}
