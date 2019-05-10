package image

import (
	"errors"
	"image"
	"image/draw"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"os"
	"path/filepath"
	"strings"
)

//Pos 水印的位置
type Pos int

// 水印的位置
const (
	TopLeft Pos = iota
	TopRight
	BottomLeft
	BottomRight
	Center
)

// 允许做水印的图片
var watermarkExts = []string{
	".gif", ".jpg", ".jpeg", ".png",
}

// Watermark 用于给图片添加水印功能。
// 目前支持gif,jpeg和png三种图片格式。
// 若是gif图片，则只取图片的第一帧；png支持透明背景。
type Watermark struct {
	image   image.Image // 水印图片
	padding int         // 水印留的边白
	pos     Pos         // 水印的位置
}

var ErrUnsupportedWatermarkType = errors.New(`水印图片格式不支持`)
var ErrInvalidPos = errors.New(`水印位置不正确`)

//NewWatermark 设置水印的相关参数。
// path为水印文件的路径；
// padding为水印在目标不图像上的留白大小；
// pos水印的位置。
func NewWatermark(path string, padding int, pos Pos) (*Watermark, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var img image.Image
	switch strings.ToLower(filepath.Ext(path)) {
	case ".jpg", ".jpeg":
		img, err = jpeg.Decode(f)
	case ".png":
		img, err = png.Decode(f)
	case ".gif":
		img, err = gif.Decode(f)
	default:
		return nil, ErrUnsupportedWatermarkType
	}
	if err != nil {
		return nil, err
	}

	if pos < TopLeft || pos > Center {
		return nil, ErrInvalidPos
	}

	return &Watermark{
		image:   img,
		padding: padding,
		pos:     pos,
	}, nil
}

//IsAllowExt 该扩展名的图片是否允许使用水印
func (w *Watermark) IsAllowExt(ext string) bool {
	for _, e := range watermarkExts {
		if e == ext {
			return true
		}
	}
	return false
}

//MarkFile 给指定的文件打上水印
func (w *Watermark) MarkFile(path string) error {
	// 此处不能使用os.O_APPEND 在osx下会造成seek失效。
	// TODO:验证其它系统的正确性
	file, err := os.OpenFile(path, os.O_RDWR, os.ModePerm)
	if err != nil {
		return err
	}
	defer file.Close()

	return w.Mark(file, strings.ToLower(filepath.Ext(path)))
}

//Mark 将水印写入src中，由ext确定当前图片的类型。
func (w *Watermark) Mark(src io.ReadWriteSeeker, ext string) error {
	var srcImg image.Image
	var err error

	switch ext {
	case ".jpg", ".jpeg":
		srcImg, err = jpeg.Decode(src)
	case ".png":
		srcImg, err = png.Decode(src)
	case ".gif":
		srcImg, err = gif.Decode(src)
	default:
		return ErrUnsupportedWatermarkType
	}
	if err != nil {
		return err
	}

	var point image.Point
	srcw := srcImg.Bounds().Dx()
	srch := srcImg.Bounds().Dy()
	switch w.pos {
	case TopLeft:
		point = image.Point{X: -w.padding, Y: -w.padding}
	case TopRight:
		point = image.Point{
			X: -(srcw - w.padding - w.image.Bounds().Dx()),
			Y: -w.padding,
		}
	case BottomLeft:
		point = image.Point{
			X: -w.padding,
			Y: -(srch - w.padding - w.image.Bounds().Dy()),
		}
	case BottomRight:
		point = image.Point{
			X: -(srcw - w.padding - w.image.Bounds().Dx()),
			Y: -(srch - w.padding - w.image.Bounds().Dy()),
		}
	case Center:
		point = image.Point{
			X: -(srcw - w.padding - w.image.Bounds().Dx()) / 2,
			Y: -(srch - w.padding - w.image.Bounds().Dy()) / 2,
		}
	}

	dstImg := image.NewNRGBA64(srcImg.Bounds())
	draw.Draw(dstImg, dstImg.Bounds(), srcImg, image.ZP, draw.Src)
	draw.Draw(dstImg, dstImg.Bounds(), w.image, point, draw.Over)

	_, err = src.Seek(0, 0)
	if err != nil {
		return err
	}
	switch ext {
	case ".jpg", ".jpeg":
		err = jpeg.Encode(src, dstImg, nil)
	case ".png":
		err = png.Encode(src, dstImg)
	case ".gif":
		err = gif.Encode(src, dstImg, nil)
		// default: // 由前一个Switch确保此处没有default的出现。
	}
	return err
}
