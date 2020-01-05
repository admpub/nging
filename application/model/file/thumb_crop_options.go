/*
   Nging is a toolbox for webmasters
   Copyright (C) 2018-present  Wenhui Shen <swh@admpub.com>

   This program is free software: you can redistribute it and/or modify
   it under the terms of the GNU Affero General Public License as published
   by the Free Software Foundation, either version 3 of the License, or
   (at your option) any later version.

   This program is distributed in the hope that it will be useful,
   but WITHOUT ANY WARRANTY; without even the implied warranty of
   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
   GNU Affero General Public License for more details.

   You should have received a copy of the GNU Affero General Public License
   along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/

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
