package file

import (
	"io"

	imageproxy "github.com/admpub/imageproxy"
	"github.com/admpub/nging/application/dbschema"
	"github.com/admpub/nging/application/registry/upload/driver"
)

// CropOptions 图片裁剪选项
type CropOptions struct {
	Options       *imageproxy.Options //裁剪方式设置
	File          *dbschema.File      //原图信息
	SrcReader     io.Reader           //原图reader
	Storer        driver.Storer       //存储器
	DestFile      string              //保存文件路径
	FileMD5       string              //原图MD5
	WatermarkFile string              //水印图片文件
}
