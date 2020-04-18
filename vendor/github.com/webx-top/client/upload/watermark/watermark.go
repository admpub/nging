package watermark

import (
	"bytes"
	"io"
	"io/ioutil"
	"mime/multipart"
	"os"

	"github.com/admpub/errors"
	rifs "github.com/admpub/go-utility/filesystem"

	"github.com/webx-top/image"
)

func Write(rw io.ReadWriteSeeker, ext string, opt *image.WatermarkOptions) ([]byte, error) {
	wm, err := opt.CreateInstance()
	if err != nil {
		return nil, errors.WithMessage(err, `NewWatermark`)
	}
	err = wm.Mark(rw, ext)
	rw.Seek(0, os.SEEK_SET)
	if err != nil {
		return nil, errors.WithMessage(err, `Mark`)
	}
	return ioutil.ReadAll(rw)
}

// Bytes 添加水印到图片字节数据中
func Bytes(b []byte, ext string, opt *image.WatermarkOptions) ([]byte, error) {
	sb, err := Bytes2readWriteSeeker(b)
	if err != nil {
		return nil, errors.WithMessage(err, `Bytes2readWriteSeeker`)
	}
	return Write(sb, ext, opt)
}

func Bytes2readWriteSeeker(b []byte) (io.ReadWriteSeeker, error) {
	sb := rifs.NewSeekableBuffer()
	_, err := sb.Write(b)
	if err != nil {
		return sb, err
	}
	sb.Seek(0, os.SEEK_SET)
	return sb, err
}

func Bytes2file(b []byte) multipart.File {
	r := io.NewSectionReader(bytes.NewReader(b), 0, int64(len(b)))
	return sectionReadCloser{r}
}

type sectionReadCloser struct {
	*io.SectionReader
}

func (rc sectionReadCloser) Close() error {
	return nil
}
