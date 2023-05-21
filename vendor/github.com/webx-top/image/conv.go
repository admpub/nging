package image

import (
	"bytes"
	"encoding/base64"
	"image"
	"image/jpeg"
	"io"
	"io/ioutil"
	"mime/multipart"
	"os"
)

// Base64ToFile base64 -> file
func Base64ToFile(base64Data string, file string) error {
	b, err := base64.StdEncoding.DecodeString(base64Data)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(file, b, 0666)
}

// Base64ToBuffer base64 -> buffer
func Base64ToBuffer(base64Data string) (*bytes.Buffer, error) {
	b, err := base64.StdEncoding.DecodeString(base64Data) //成图片文件并把文件写入到buffer
	if err != nil {
		return nil, err
	}
	return bytes.NewBuffer(b), nil // 必须加一个buffer 不然没有read方法就会报错
}

// 转换成buffer之后里面就有Reader方法了。才能被图片API decode

// BufferToImageBuffer buffer-> ImageBuff（图片裁剪,代码接上面）
func BufferToImageBuffer(b *bytes.Buffer, x1 int, y1 int) (*image.YCbCr, error) {
	m, _, err := image.Decode(b) // 图片文件解码
	if err != nil {
		return nil, err
	}
	rgbImg := m.(*image.YCbCr)
	subImg := rgbImg.SubImage(image.Rect(0, 0, x1, y1)).(*image.YCbCr) //图片裁剪x0 y0 x1 y1
	return subImg, nil
}

type SubImage interface {
	SubImage(r image.Rectangle) image.Image
}

func BufferToImageXBuffer(b *bytes.Buffer, x1 int, y1 int) (image.Image, error) {
	src, _, err := image.Decode(b) // 图片文件解码
	if err != nil {
		return nil, err
	}
	var subImg image.Image
	switch img := src.(type) {
	case SubImage:
		subImg = img.SubImage(image.Rect(0, 0, x1, y1))
	default:
		err = ErrUnsupportedImageType
	}
	return subImg, err
}

// ToFile img -> file(代码接上面)
func ToFile(subImg *image.YCbCr, file string) error {
	f, err := os.Create(file) //创建文件
	if err != nil {
		return err
	}
	defer f.Close()                    //关闭文件
	return jpeg.Encode(f, subImg, nil) //写入文件
}

// ToBase64 img -> base64(代码接上面)
func ToBase64(subImg *image.YCbCr) ([]byte, error) {
	src := bytes.NewBuffer(nil)          //开辟一个新的空buff
	err := jpeg.Encode(src, subImg, nil) //img写入到buff
	if err != nil {
		return nil, err
	}
	dist := make([]byte, base64.StdEncoding.EncodedLen(src.Len()))
	base64.StdEncoding.Encode(dist, src.Bytes()) //buff转成base64
	return dist, nil
}

// FileToBase64 imgFile -> base64
func FileToBase64(srcFile string) ([]byte, error) {
	ff, err := ioutil.ReadFile(srcFile)
	if err != nil {
		return nil, err
	}
	dist := make([]byte, 5000000)       //数据缓存
	base64.StdEncoding.Encode(dist, ff) // 文件转base64
	//_ = ioutil.WriteFile("./output2.jpg.txt", dist, 0666) //直接写入到文件就ok完活了。
	return dist, nil
}

// Bytes2file 直接转 multipart.File
func Bytes2file(b []byte) multipart.File {
	return ReaderAt2file(bytes.NewReader(b), int64(len(b)))
}

// ReaderAt2file 直接转 multipart.File
func ReaderAt2file(readerAt io.ReaderAt, size int64) multipart.File {
	r := io.NewSectionReader(readerAt, 0, size)
	return SectionReadCloser{r}
}

type SectionReadCloser struct {
	*io.SectionReader
}

func (rc SectionReadCloser) Close() error {
	return nil
}
