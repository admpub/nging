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
	opt := NewTextImageOptions()
	opt.Text = textContent
	opt.FontFile = fontFile
	opt.Width = width
	opt.Height = height
	if rlen > 0 {
		pngBgImgPath := args[0].(string)
		var img image.Image
		if len(pngBgImgPath) > 0 {
			pngFile, err := os.Open(pngBgImgPath)
			checkErr(err)
			defer pngFile.Close()
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
		opt.Width = width
		opt.Height = height
		textLayer = TextImage(opt)
		width = textLayer.Bounds().Max.X
		height = textLayer.Bounds().Max.Y
		c = image.NewRGBA(image.Rect(0, 0, width, height))
		if len(pngBgImgPath) > 0 {
			// draw bottom layer from file
			draw.Draw(c, c.Bounds(), img, image.Point{0, 0}, draw.Src)
		}
	} else {
		textLayer = TextImage(opt)
		width = textLayer.Bounds().Max.X
		height = textLayer.Bounds().Max.Y
		c = image.NewRGBA(image.Rect(0, 0, width, height))
	}

	// draw text layer on top
	draw.Draw(c, c.Bounds(), textLayer, image.Point{0, 0}, draw.Over)

	return c
}

func NewTextImageOptions() *TextImageOptions {
	return &TextImageOptions{
		FontSize: 8,
		PosX:     10,
		PosY:     10,
		Spacing:  1.5,
		DPI:      300,
	}
}

type TextImageOptions struct {
	Text     string
	FontFile string
	FontSize float64
	Width    int
	Height   int
	PosX     int
	PosY     int
	Spacing  float64
	DPI      float64
}

func TextImage(opt *TextImageOptions) *image.RGBA {
	lines := strings.Split(opt.Text, "\n")
	// read font
	fontBytes, err := ioutil.ReadFile(opt.FontFile)
	checkErr(err)
	font, err := freetype.ParseFont(fontBytes)
	checkErr(err)

	c := freetype.NewContext()
	c.SetDPI(opt.DPI)
	c.SetFont(font)
	c.SetFontSize(opt.FontSize)

	// Initialize the context.
	fg, bg := image.White, image.Transparent
	if opt.Width <= 0 {
		s := float64(c.PointToFixed(opt.FontSize) >> 8)
		opt.Width = int(math.Ceil((s-float64(s)/2.5)*float64(len(opt.Text)) + float64(opt.PosX)*2))
	}
	if opt.Height <= 0 {
		opt.Height = len(lines)*int(c.PointToFixed(opt.FontSize*opt.Spacing)>>8) + opt.PosY*2
	}
	rgba := image.NewRGBA(image.Rect(0, 0, opt.Width, opt.Height))

	draw.Draw(rgba, rgba.Bounds(), bg, image.ZP, draw.Src)
	c.SetClip(rgba.Bounds())
	c.SetDst(rgba)
	c.SetSrc(fg)

	// Draw the text
	pt := freetype.Pt(opt.PosX, opt.PosY+int(c.PointToFixed(opt.FontSize)>>8))
	for _, s := range lines {
		_, err = c.DrawString(s, pt)
		checkErr(err)
		pt.Y += c.PointToFixed(opt.FontSize * opt.Spacing)
	}

	return rgba
}
