/*

   Copyright 2016-present Wenhui Shen <www.webx.top>

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

package fields

import (
	"github.com/coscms/forms/common"
)

// SubmitButton creates a default button with the provided name and text.
func SubmitButton(name string, text string) *Field {
	ret := FieldWithType(name, common.SUBMIT)
	ret.SetText(text)
	return ret
}

// ResetButton creates a default reset button with the provided name and text.
func ResetButton(name string, text string) *Field {
	ret := FieldWithType(name, common.RESET)
	ret.SetText(text)
	return ret
}

// Button creates a default generic button
func Button(name string, text string) *Field {
	ret := FieldWithType(name, common.BUTTON)
	ret.SetText(text)
	return ret
}
