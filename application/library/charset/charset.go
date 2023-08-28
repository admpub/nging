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

package charset

import (
	"bytes"
	"errors"
	"strings"
	"unicode/utf8"

	"github.com/admpub/mahonia"
	runewidth "github.com/mattn/go-runewidth"
	"github.com/webx-top/com"
)

func NewDecoderAndEncoder(fromEnc string, toEnc string) (mahonia.Decoder, mahonia.Encoder, error) {
	fromEnc = strings.ToLower(fromEnc)
	toEnc = strings.ToLower(toEnc)
	if fromEnc == toEnc {
		return nil, nil, nil
	}
	if fromEnc == `utf8` {
		fromEnc = `utf-8`
	}
	if toEnc == `utf8` {
		toEnc = `utf-8`
	}
	if !Validate(fromEnc) {
		return nil, nil, errors.New(`Unsupported charset: ` + fromEnc)
	}
	if !Validate(toEnc) {
		return nil, nil, errors.New(`Unsupported charset: ` + toEnc)
	}
	var dec mahonia.Decoder
	var enc mahonia.Encoder
	if fromEnc != `utf-8` {
		dec = mahonia.NewDecoder(fromEnc)
	}
	if toEnc != `utf-8` {
		enc = mahonia.NewEncoder(toEnc)
	}
	return dec, enc, nil
}

func Convert(fromEnc string, toEnc string, b []byte) ([]byte, error) {
	dec, enc, err := NewDecoderAndEncoder(fromEnc, toEnc)
	if err != nil {
		return nil, err
	}
	var s string
	if dec != nil {
		s = dec.ConvertString(com.Bytes2str(b))
	}
	if enc != nil {
		if len(s) > 0 {
			s = enc.ConvertString(s)
		} else {
			s = enc.ConvertString(com.Bytes2str(b))
		}
	}
	if len(s) > 0 {
		b = com.Str2bytes(s)
	}
	return b, nil
}

func NewConvertBytesFunc(fromEnc string, toEnc string) (func([]byte) []byte, error) {
	dec, enc, err := NewDecoderAndEncoder(fromEnc, toEnc)
	if err != nil {
		return nil, err
	}
	return func(b []byte) []byte {
		var s string
		if dec != nil {
			s = dec.ConvertString(com.Bytes2str(b))
		}
		if enc != nil {
			if len(s) > 0 {
				s = enc.ConvertString(s)
			} else {
				s = enc.ConvertString(com.Bytes2str(b))
			}
		}
		if len(s) > 0 {
			b = com.Str2bytes(s)
		}
		return b
	}, nil
}

func NewConvertFunc(fromEnc string, toEnc string) (func(string) string, error) {
	dec, enc, err := NewDecoderAndEncoder(fromEnc, toEnc)
	if err != nil {
		return nil, err
	}
	return func(s string) string {
		if dec != nil {
			s = dec.ConvertString(s)
		}
		if enc != nil {
			s = enc.ConvertString(s)
		}
		return s
	}, nil
}

func Validate(enc string) bool {
	return mahonia.GetCharset(enc) != nil
}

func Truncate(str string, width int) string {
	w := 0
	b := []byte(str)
	var buf bytes.Buffer
	for len(b) > 0 {
		r, size := utf8.DecodeRune(b)
		rw := runewidth.RuneWidth(r)
		if w+rw > width {
			break
		}
		buf.WriteRune(r)
		w += rw
		b = b[size:]
	}
	return buf.String()
}

func With(str string) int {
	return runewidth.StringWidth(str)
}

func RuneWith(str string) int {
	return utf8.RuneCountInString(str)
}
