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

package settings

import (
	"errors"
	"strings"

	"github.com/webx-top/com"
	"github.com/webx-top/echo"

	"github.com/admpub/nging/application/dbschema"
	"github.com/admpub/nging/application/library/common"
)

type Encoder func(v *dbschema.NgingConfig, r echo.H) ([]byte, error)

var encoders = map[string]Encoder{}

func Encoders() map[string]Encoder {
	return encoders
}

func GetEncoder(group string) Encoder {
	ps, _ := encoders[group]
	return ps
}

// RegisterEncoder 注册配置值编码器（用于客户端提交表单数据之后的编码操作，编码结果保存到数据库）
// 名称支持"group"或"group.key"两种格式，例如:
// settings.RegisterDecoder(`sms`,...)对整个sms组的配置有效
// settings.RegisterDecoder(`sms.twilio`,...)对sms组内key为twilio的配置有效
func RegisterEncoder(group string, encoder Encoder) {
	encoders[group] = encoder
}

var ErrNotExists = errors.New(`Not exists`)

func EncodeConfigValue(_v *echo.Mapx, v *dbschema.NgingConfig, encoder Encoder) (value string, err error) {
	if _v.IsMap() {
		var b []byte
		var e error
		store := _v.AsStore()
		if encoder != nil {
			b, e = encoder(v, store)
			if e != nil {
				err = e
				return
			}
		}
		if subEncoder := GetEncoder(v.Group + `.` + v.Key); subEncoder != nil {
			b, e = subEncoder(v, store)
		} else {
			b, e = com.JSONEncode(store)
		}
		if e != nil {
			err = e
			return
		}
		value = string(b)
	} else if _v.IsSlice() {
		items := []string{}
		for _, item := range _v.AsFlatSlice() {
			item = strings.TrimSpace(item)
			if len(item) == 0 {
				return
			}
			items = append(items, item)
		}
		value = strings.Join(items, `,`)
	} else {
		value = _v.Value() //c.Form(group + `[` + v.Key + `]`)
	}
	value = DefaultEncoder(v, value)
	return
}

func ContentEncode(content string, contypes ...string) string {
	var contype string
	if len(contypes) > 0 {
		contype = contypes[0]
	}
	switch contype {
	case `html`:
		content = common.RemoveXSS(content)

	case `url`, `image`, `video`, `audio`, `file`:
		content = common.MyCleanText(content)

	case `id`, `text`:
		content = com.StripTags(content)

	case `json`:
		// pass

	case `markdown`:
		// pass

	case `list`:
		content = strings.TrimSpace(content)
		content = strings.Trim(content, `,`)

	default:
		content = com.StripTags(content)
	}
	return content
}

func DefaultEncoder(v *dbschema.NgingConfig, value string) string {
	value = ContentEncode(value, v.Type)
	if v.Encrypted == `Y` {
		value = echo.Get(`DefaultConfig`).(Codec).Encode(value)
	}
	return value
}
