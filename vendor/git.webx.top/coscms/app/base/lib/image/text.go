package image

import (
	"bufio"
	"image"
	"image/draw"
	"image/png"
	"io/ioutil"
	"math"
	"net/http"
	"os"
	"strings"

	"github.com/golang/freetype"
)

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func MergeTextImageServeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "image/png")
	// figure out the Content-Length of this new combined image somehow
	// w.Header().Set("Content-Length", fmt.Sprint(pngImage.ContentLength))
	textContent := "coscms.com"
	pngBgImgPath := "data/fonts/background.png"
	img := TextToImage(textContent, "data/fonts/Courier New.ttf", pngBgImgPath)
	b := bufio.NewWriter(w)
	png.Encode(b, img)
}

//TextToImage textContent string, fontFile string, pngBgImgPath string, width int, height int
func TextToImage(textContent string, fontFile string, args ...interface{}) *image.RGBA {
	var c, textLayer *image.RGBA
	rlen := len(args)
	var width, height int
	if rlen > 0 {
		pngBgImgPath := args[0].(string)
		var img image.Image
		if pngBgImgPath != "" {
			pngFile, err := os.Open(pngBgImgPath)
			checkErr(err)
			img, err = png.Decode(pngFile)
			checkErr(err)
			width = img.Bounds().Max.X
			height = img.Bounds().Max.Y
		}
		if rlen > 1 {
			width = args[1].(int)
		}
		if rlen > 2 {
			height = args[2].(int)
		}
		textLayer = TextImage(textContent, fontFile, width, height)
		width = textLayer.Bounds().Max.X
		height = textLayer.Bounds().Max.Y
		c = image.NewRGBA(image.Rect(0, 0, width, height))
		if pngBgImgPath != "" {
			// draw bottom layer from file
			draw.Draw(c, c.Bounds(), img, image.Point{0, 0}, draw.Src)
		}
	} else {
		textLayer = TextImage(textContent, fontFile, width, height)
		width = textLayer.Bounds().Max.X
		height = textLayer.Bounds().Max.Y
		c = image.NewRGBA(image.Rect(0, 0, width, height))
	}

	// draw text layer on top
	draw.Draw(c, c.Bounds(), textLayer, image.Point{0, 0}, draw.Over)

	return c
}

func TextImage(textContent string, fontFile string, width int, height int) *image.RGBA {
	lines := strings.Split(textContent, "\n")
	posX := 10
	posY := 10
	spacing := 1.5
	var fontsize float64 = 8

	// read font
	fontBytes, err := ioutil.ReadFile(fontFile)
	checkErr(err)
	font, err := freetype.ParseFont(fontBytes)
	checkErr(err)

	c := freetype.NewContext()
	c.SetDPI(300)
	c.SetFont(font)
	c.SetFontSize(fontsize)

	// Initialize the context.
	fg, bg := image.White, image.Transparent
	if width <= 0 {
		s := float64(c.PointToFixed(fontsize) >> 8)
		width = int(math.Ceil((s-float64(s)/2.5)*float64(len(textContent)) + float64(posX)*2))
	}
	if height <= 0 {
		height = len(lines)*int(c.PointToFixed(fontsize*spacing)>>8) + posY*2
	}
	rgba := image.NewRGBA(image.Rect(0, 0, width, height))

	draw.Draw(rgba, rgba.Bounds(), bg, image.ZP, draw.Src)
	c.SetClip(rgba.Bounds())
	c.SetDst(rgba)
	c.SetSrc(fg)

	// Draw the text
	pt := freetype.Pt(posX, posY+int(c.PointToFixed(fontsize)>>8))
	for _, s := range lines {
		_, err = c.DrawString(s, pt)
		checkErr(err)
		pt.Y += c.PointToFixed(fontsize * spacing)
	}

	return rgba
}
