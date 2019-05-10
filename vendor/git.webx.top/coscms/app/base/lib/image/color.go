package image

import (
	"fmt"
	"image/color"
	"image/png"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	ID "git.webx.top/coscms/app/base/lib/image/identicon"
	"github.com/webx-top/com"
)

/**
 * 颜色代码转换为RGB
 * input int
 * output int red, green, blue
 **/
func ColorToRGB(color int) (red, green, blue uint8) {
	red = uint8(color >> 16)
	green = uint8((color & 0x00FF00) >> 8)
	blue = uint8(color & 0x0000FF)
	return
}

func ColorHexToRGB(color string) (red, green, blue uint8) {
	color = strings.TrimPrefix(color, "#")          //过滤掉16进制前缀
	color64, err := strconv.ParseInt(color, 16, 32) //字串到整型数据
	if err != nil {
		fmt.Println(err)
		return
	}
	color32 := int(color64) //类型强转
	red, green, blue = ColorToRGB(color32)
	return
}

func MakeAvatar(widthHeight int, data []byte, saveAs string) bool {
	var red, green, blue uint8 = 0, 153, 204
	/*
		var colorStr string = a.GetString("color")
		if colorStr != "" {
			red, green, blue = ColorHexToRGB(colorStr)
		}
		fmt.Println(red, green, blue)
	*/
	imger, _ := ID.New(128,
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
	img := imger.Make(data)
	dir := filepath.Dir(saveAs)
	if com.IsDir(dir) == false {
		os.MkdirAll(dir, 0777)
	}
	fi, _ := os.Create(saveAs + ".png")
	png.Encode(fi, img)
	fi.Close()
	return true
}
