/*
   Nging is a toolbox for webmasters
   Copyright (C) 2018-present Wenhui Shen <swh@admpub.com>

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
	"strings"

	"github.com/webx-top/com"
	"github.com/webx-top/echo"

	"github.com/admpub/nging/v4/application/dbschema"
)

type Decoder func(v *dbschema.NgingConfig, dbschemaMap echo.H) error

var decoders = map[string]Decoder{}

func Decoders() map[string]Decoder {
	return decoders
}

func GetDecoder(group string) Decoder {
	ps, _ := decoders[group]
	return ps
}

type Codec interface {
	Decode(string, ...string) string
	Encode(string, ...string) string
}

// RegisterDecoder 注册配置值解码器（用于从数据库读出来之后的解码操作）
// 名称支持"group"或"group.key"两种格式，例如:
// settings.RegisterDecoder(`sms`,...)对整个sms组的配置有效
// settings.RegisterDecoder(`sms.twilio`,...)对sms组内key为twilio的配置有效
func RegisterDecoder(group string, decoder Decoder) {
	decoders[group] = decoder
}

func DecodeConfigValue(v *dbschema.NgingConfig, decoder Decoder) (echo.H, error) {
	if v.Encrypted == `Y` {
		v.Value = echo.Get(`FromFile()`).(Codec).Decode(v.Value)
	}
	r := echo.H(v.AsMap())
	var err error
	if decoder != nil {
		err = decoder(v, r)
		if err != nil {
			return nil, err
		}
	}
	if subDecoder := GetDecoder(v.Group + `.` + v.Key); subDecoder != nil {
		err = subDecoder(v, r)
	} else {
		err = DefaultDecoder(v, r)
	}
	return r, err
}

func DecodeConfig(v *dbschema.NgingConfig, cfg echo.H, decoder Decoder) (echo.H, error) {
	r, e := DecodeConfigValue(v, decoder)
	if e != nil {
		return cfg, e
	}
	cfg.Set(v.Key, r)
	return cfg, nil
}

func DefaultDecoder(v *dbschema.NgingConfig, r echo.H) error {
	if r.Has(`ValueObject`) {
		return nil
	}
	switch v.Type {
	case `json`:
		jsonData := echo.H{}
		if len(v.Value) > 0 {
			err := com.JSONDecode([]byte(v.Value), &jsonData)
			if err != nil {
				return err
			}
		}
		r[`ValueObject`] = jsonData
	case `list`:
		jsonData := echo.H{}
		if len(v.Value) > 0 {
			com.JSONDecode([]byte(v.Value), &jsonData)
		}
		v.Value = strings.Trim(v.Value, `,`)
		if len(v.Value) > 0 {
			r[`ValueObject`] = strings.Split(v.Value, `,`)
		} else {
			r[`ValueObject`] = []string{}
		}
	default:
		r[`ValueObject`] = nil
	}
	return nil
}
