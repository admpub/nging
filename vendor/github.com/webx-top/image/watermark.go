package image

import (
	"fmt"
	"image"
	"image/draw"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/admpub/errors"
	godl "github.com/admpub/go-download/v2"
	"github.com/webx-top/com"
	"golang.org/x/image/bmp"
	_ "golang.org/x/image/webp"
)

// Pos 水印的位置
type Pos int

// 水印的位置
const (
	TopLeft Pos = iota
	TopRight
	BottomLeft
	BottomRight
	Center
)

const sniffLen = 512

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
	retry   int
	debug   bool
}

func NewWatermarkData(f []byte) *WatermarkData {
	return &WatermarkData{Data: f, Time: time.Now()}
}

type WatermarkData struct {
	Data []byte
	Time time.Time
}

// FileReader 文件读取接口
type FileReader interface {
	io.Closer
	io.Reader
	io.Seeker
}

var (
	// ErrUnsupportedWatermarkType 水印图片格式不支持
	ErrUnsupportedWatermarkType = errors.New(`水印图片格式不支持`)
	// ErrInvalidPos 水印位置不正确
	ErrInvalidPos = errors.New(`水印位置不正确`)
	// DefaultHTTPSystemOpen 水印文件默认打开方式
	DefaultHTTPSystemOpen = func(name string) (FileReader, error) {
		if strings.Contains(name, `://`) {
			return GetRemoteWatermarkFileData(name)
		}
		fp, err := os.Open(name) // 用完后别忘了关闭 fp.Close()
		return fp, err
	}
	// WatermarkOpen 水印文件打开方式
	WatermarkOpen = DefaultHTTPSystemOpen
	// DefaultWatermarkFileDownloadClient 默认水印文件下载客户端
	DefaultWatermarkFileDownloadClient = com.HTTPClientWithTimeout(30 * time.Second)
	// DefaultWatermarkFileDownloadOptions 默认水印文件下载选项
	DefaultWatermarkFileDownloadOptions = &godl.Options{
		Client: func() http.Client {
			return *DefaultWatermarkFileDownloadClient
		},
	}
)

// NewWatermark 设置水印的相关参数。
// path为水印文件的路径；
// padding为水印在目标不图像上的留白大小；
// pos水印的位置。
func NewWatermark(path string, padding int, pos Pos) (*Watermark, error) {
	if pos < TopLeft || pos > Center {
		return nil, ErrInvalidPos
	}
	f, err := WatermarkOpen(path)
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
		return nil, errors.WithMessage(ErrUnsupportedWatermarkType, path)
	}
	if err != nil {
		return nil, err
	}

	return &Watermark{
		image:   img,
		padding: padding,
		pos:     pos,
	}, nil
}

// IsAllowExt 该扩展名的图片是否允许使用水印
func (w *Watermark) IsAllowExt(ext string) bool {
	for _, e := range watermarkExts {
		if e == ext {
			return true
		}
	}
	return false
}

// MarkFile 给指定的文件打上水印
func (w *Watermark) MarkFile(path string, dest ...string) error {
	file, err := os.OpenFile(path, os.O_RDWR, os.ModePerm)
	if err != nil {
		return err
	}
	defer file.Close()

	return w.Mark(file, filepath.Ext(path), dest...)
}

func GetContentTypeByContent(buffer []byte) string {
	// Use the net/http package's handy DetectContentType function. Always returns a valid
	// content-type by returning "application/octet-stream" if no others seemed to match.
	contentType := http.DetectContentType(buffer)
	return contentType
}

func GetExtensionByContentType(contentType string) string {
	switch {
	case strings.Contains(contentType, "image/jpeg"):
		return ".jpg"
	case strings.Contains(contentType, "image/png"):
		return ".png"
	case strings.Contains(contentType, "image/bmp"):
		return ".bmp"
	case strings.Contains(contentType, "image/gif"):
		return ".gif"
	default:
		return ""
	}
}

func IsFormatError(err error) bool {
	var isFormatError bool
	switch err.(type) {
	case png.FormatError, jpeg.FormatError:
		isFormatError = true
	default:
		if err == bmp.ErrUnsupported || strings.Contains(err.Error(), `can't recognize format`) {
			isFormatError = true
		}
	}
	return isFormatError
}

// Mark 将水印写入src中，由ext确定当前图片的类型。
func (w *Watermark) Mark(src io.ReadWriteSeeker, ext string, dest ...string) error {
	ext = strings.ToLower(ext)
	var srcImg image.Image
	var err error
	switch ext {
	case ".jpg", ".jpeg":
		srcImg, err = jpeg.Decode(src)
	case ".png":
		srcImg, err = png.Decode(src)
	case ".gif":
		srcImg, err = gif.Decode(src)
	case ".bmp":
		srcImg, err = bmp.Decode(src)
	default:
		return errors.WithMessage(ErrUnsupportedWatermarkType, ext)
	}
	if err != nil {
		if w.retry < 1 && IsFormatError(err) {
			body := make([]byte, sniffLen)
			src.Seek(0, 0)
			if _, err := io.ReadFull(src, body); err != nil {
				return err
			}
			contentType := GetContentTypeByContent(body)
			newExt := GetExtensionByContentType(contentType)
			if len(newExt) == 0 || ext == newExt {
				return err
			}
			src.Seek(0, 0)
			w.retry++
			return w.Mark(src, newExt)
		}
		return err
	}

	var point image.Point
	srcw := srcImg.Bounds().Dx()
	srch := srcImg.Bounds().Dy()
	if srcw/w.image.Bounds().Dx() < 3 || srch/w.image.Bounds().Dy() < 3 {
		if w.debug {
			fmt.Println(`[skip] image is too small: `, fmt.Sprintf("(watermark:%dx%d\t=>\tsrc-image:%dx%d)", w.image.Bounds().Dx(), w.image.Bounds().Dy(), srcw, srch))
		}
		return err
	}
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
	draw.Draw(dstImg, dstImg.Bounds(), srcImg, image.Point{}, draw.Src)
	draw.Draw(dstImg, dstImg.Bounds(), w.image, point, draw.Over)
	if len(dest) > 0 && len(dest[0]) > 0 {
		var file *os.File
		file, err = os.Create(dest[0])
		if err != nil {
			return err
		}
		defer file.Close()
		src = file
	} else {
		_, err = src.Seek(0, 0)
		if err != nil {
			return err
		}
	}
	switch ext {
	case ".jpg", ".jpeg":
		err = jpeg.Encode(src, dstImg, nil)
	case ".png":
		err = png.Encode(src, dstImg)
	case ".gif":
		err = gif.Encode(src, dstImg, nil)
	case ".bmp":
		err = bmp.Encode(src, dstImg)
		// default: // 由前一个Switch确保此处没有default的出现。
	}
	return err
}
