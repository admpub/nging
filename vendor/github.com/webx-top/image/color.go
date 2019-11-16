package image

import (
	"image/color"
	"image/png"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	ID "github.com/admpub/identicon"

	"github.com/webx-top/com"
)

// ColorToRGB 颜色代码转换为RGB
// input int
// output int red, green, blue
func ColorToRGB(color int) (red, green, blue uint8) {
	red = uint8(color >> 16)
	green = uint8((color & 0x00FF00) >> 8)
	blue = uint8(color & 0x0000FF)
	return
}

func ColorHexToRGB(color string) (red, green, blue uint8, err error) {
	color = strings.TrimPrefix(color, "#") //过滤掉16进制前缀
	var color64 int64
	color64, err = strconv.ParseInt(color, 16, 32) //字串到整型数据
	if err != nil {
		return
	}
	color32 := int(color64) //类型强转
	red, green, blue = ColorToRGB(color32)
	return
}

func MakeAvatar(widthHeight int, data []byte, saveAs string) error {
	var red, green, blue uint8 = 0, 153, 204
	var err error
	/*
		var colorStr string = a.GetString("color")
		if len(colorStr)>0 {
			red, green, blue, err = ColorHexToRGB(colorStr)
			if err != nil {
				return err
			}
		}
		fmt.Println(red, green, blue)
	*/
	imger, err := ID.New(128,
		color.White,
		color.RGBA{red, green, blue, 100},
		color.RGBA{255, 153, 0, 100},
		color.RGBA{0, 51, 102, 100},
		color.RGBA{153, 0, 51, 100},
		color.RGBA{0, 153, 153, 100},
		color.RGBA{255, 255, 204, 100},
		color.RGBA{204, 255, 255, 100},
		color.RGBA{51, 51, 153, 100},
		color.RGBA{255, 102, 102, 100},
	)
	if err != nil {
		return err
	}
	img := imger.Make(data)
	dir := filepath.Dir(saveAs)
	if !com.IsDir(dir) {
		err = os.MkdirAll(dir, 0777)
		if err != nil {
			return err
		}
	}
	fi, err := os.Create(saveAs + ".png")
	if err != nil {
		return err
	}
	defer fi.Close()
	return png.Encode(fi, img)
}
