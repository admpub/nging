package watermark

import (
	"fmt"
	"io"

	rifs "github.com/admpub/go-utility/filesystem"

	"github.com/webx-top/image"
)

func Write(rw io.ReadWriteSeeker, ext string, opt *image.WatermarkOptions) ([]byte, error) {
	wm, err := opt.CreateInstance()
	if err != nil {
		return nil, fmt.Errorf(`NewWatermark: %w`, err)
	}
	err = wm.Mark(rw, ext)
	rw.Seek(0, io.SeekStart)
	if err != nil {
		return nil, fmt.Errorf(`Mark: %w`, err)
	}
	return io.ReadAll(rw)
}

// Bytes 添加水印到图片字节数据中
func Bytes(b []byte, ext string, opt *image.WatermarkOptions) ([]byte, error) {
	sb, err := Bytes2readWriteSeeker(b)
	if err != nil {
		return nil, fmt.Errorf(`Bytes2readWriteSeeker: %w`, err)
	}
	return Write(sb, ext, opt)
}

func Bytes2readWriteSeeker(b []byte) (io.ReadWriteSeeker, error) {
	sb := rifs.NewSeekableBuffer()
	_, err := sb.Write(b)
	if err != nil {
		return sb, err
	}
	sb.Seek(0, io.SeekStart)
	return sb, err
}

var Bytes2file = image.Bytes2file
