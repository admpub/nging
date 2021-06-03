/*

   Copyright 2016 Wenhui Shen <www.webx.top>

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.

*/

package echo

import (
	"errors"
	"fmt"
	"strings"
)

// T 标记为多语言文本
func T(format string, args ...interface{}) string {
	if len(args) > 0 {
		return fmt.Sprintf(format, args...)
	}
	return format
}

func E(format string, args ...interface{}) error {
	if len(args) > 0 {
		return fmt.Errorf(format, args...)
	}
	return errors.New(format)
}

type Translator interface {
	T(format string, args ...interface{}) string
	E(format string, args ...interface{}) error
	Lang() LangCode
}

func NewLangCode(language string, separator ...string) LangCode {
	l := LangCode{}
	sep := `-`
	if len(separator) > 0 && len(separator[0]) > 0 {
		sep = separator[0]
	}
	lg := strings.SplitN(language, sep, 2)
	switch len(lg) {
	case 2:
		l.CountryLower = strings.ToLower(lg[1])
		l.CountryUpper = strings.ToUpper(lg[1])
		fallthrough
	case 1:
		l.Language = strings.ToLower(lg[0])
	}
	return l
}

type LangCode struct {
	Language     string
	CountryLower string
	CountryUpper string
}

func (l LangCode) String() string {
	if len(l.CountryLower) > 0 {
		return l.Language + `-` + l.CountryLower
	}
	return l.Language
}

func (l LangCode) Format(upperCountry bool, separator ...string) string {
	var country string
	if upperCountry {
		country = l.CountryUpper
	} else {
		country = l.CountryLower
	}
	if len(country) > 0 {
		if len(separator) > 0 {
			return l.Language + separator[0] + country
		}
		return l.Language + `-` + country
	}
	return l.Language
}

var DefaultNopTranslate Translator = &NopTranslate{
	code: LangCode{
		Language: `en`,
	},
}

type NopTranslate struct {
	code LangCode
}

func (n *NopTranslate) T(format string, args ...interface{}) string {
	return T(format, args...)
}

func (n *NopTranslate) E(format string, args ...interface{}) error {
	return E(format, args...)
}

func (n *NopTranslate) Lang() LangCode {
	return n.code
}
