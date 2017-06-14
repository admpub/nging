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

import "fmt"

type Translator interface {
	T(format string, args ...interface{}) string
	Lang() string
}

var DefaultNopTranslate Translator = &NopTranslate{language: `en`}

type NopTranslate struct {
	language string
}

func (n *NopTranslate) T(format string, args ...interface{}) string {
	if len(args) > 0 {
		return fmt.Sprintf(format, args...)
	}
	return format
}

func (n *NopTranslate) Lang() string {
	return n.language
}
