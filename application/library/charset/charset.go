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
	"unicode/utf8"

	"github.com/admpub/mahonia"
	//"github.com/admpub/chardet"
	sc "github.com/admpub/mahonia"
	runewidth "github.com/mattn/go-runewidth"
)

func Convert(fromEnc string, toEnc string, b []byte) ([]byte, error) {
	if !Validate(fromEnc) {
		return nil, errors.New(`Unsuppored charset: ` + fromEnc)
	}
	if !Validate(toEnc) {
		return nil, errors.New(`Unsuppored charset: ` + toEnc)
	}
	dec := sc.NewDecoder(fromEnc)
	s := dec.ConvertString(string(b))
	enc := sc.NewEncoder(toEnc)
	s = enc.ConvertString(s)
	b = []byte(s)
	return b, nil
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
