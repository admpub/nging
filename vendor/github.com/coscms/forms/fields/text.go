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
	"fmt"
	"reflect"
	"strconv"

	"github.com/coscms/forms/common"
)

func ColorField(name string) *Field {
	return FieldWithType(name, common.COLOR)
}

func EmailField(name string) *Field {
	return FieldWithType(name, common.EMAIL)
}

func FileField(name string) *Field {
	return FieldWithType(name, common.FILE)
}

func ImageField(name string) *Field {
	return FieldWithType(name, common.IMAGE)
}

func MonthField(name string) *Field {
	return FieldWithType(name, common.MONTH)
}

func SearchField(name string) *Field {
	return FieldWithType(name, common.SEARCH)
}

func TelField(name string) *Field {
	return FieldWithType(name, common.TEL)
}

func UrlField(name string) *Field {
	return FieldWithType(name, common.URL)
}

func WeekField(name string) *Field {
	return FieldWithType(name, common.WEEK)
}

// TextField creates a default text input field based on the provided name.
func TextField(name string, typ ...string) *Field {
	var t = common.TEXT
	if len(typ) > 0 {
		t = typ[0]
	}
	return FieldWithType(name, t)
}

// PasswordField creates a default password text input field based on the provided name.
func PasswordField(name string) *Field {
	return FieldWithType(name, common.PASSWORD)
}

// =========== TEXT AREA

// TextAreaField creates a default textarea input field based on the provided name and dimensions.
func TextAreaField(name string, rows, cols int) *Field {
	ret := FieldWithType(name, common.TEXTAREA)
	ret.SetParam("rows", fmt.Sprintf("%d", rows))
	ret.SetParam("cols", fmt.Sprintf("%d", cols))
	return ret
}

// ========================

// HiddenField creates a default hidden input field based on the provided name.
func HiddenField(name string) *Field {
	return FieldWithType(name, common.HIDDEN)
}

// TextFieldFromInstance creates and initializes a text field based on its name, the reference object instance and field number.
func TextFieldFromInstance(val reflect.Value, t reflect.Type, fieldNo int, name string, useFieldValue bool, typ ...string) *Field {
	ret := TextField(name, typ...)
	if useFieldValue {
		if dateFormat := common.TagVal(t, fieldNo, "form_format"); len(dateFormat) > 0 {
			if vt, isEmpty := ConvertTime(val.Field(fieldNo).Interface()); !vt.IsZero() {
				ret.SetValue(vt.Format(dateFormat))
			} else if isEmpty {
				ret.SetValue(``)
			}
		} else {
			ret.SetValue(fmt.Sprintf("%v", val.Field(fieldNo).Interface()))
		}
	} else if v := common.TagVal(t, fieldNo, "form_value"); len(v) > 0 {
		ret.SetValue(v)
	}
	return ret
}

// PasswordFieldFromInstance creates and initializes a password field based on its name, the reference object instance and field number.
func PasswordFieldFromInstance(val reflect.Value, t reflect.Type, fieldNo int, name string, useFieldValue bool) *Field {
	ret := PasswordField(name)
	if useFieldValue {
		ret.SetValue(fmt.Sprintf("%s", val.Field(fieldNo).String()))
	} else if v := common.TagVal(t, fieldNo, "form_value"); len(v) > 0 {
		ret.SetValue(v)
	}
	return ret
}

// TextFieldFromInstance creates and initializes a text field based on its name, the reference object instance and field number.
// This method looks for "form_rows" and "form_cols" tags to add additional parameters to the field.
func TextAreaFieldFromInstance(val reflect.Value, t reflect.Type, fieldNo int, name string, useFieldValue bool) *Field {
	var rows, cols int = 20, 50
	var err error
	if v := common.TagVal(t, fieldNo, "form_rows"); len(v) > 0 {
		rows, err = strconv.Atoi(v)
		if err != nil {
			return nil
		}
	}
	if v := common.TagVal(t, fieldNo, "form_cols"); len(v) > 0 {
		cols, err = strconv.Atoi(v)
		if err != nil {
			return nil
		}
	}
	ret := TextAreaField(name, rows, cols)
	if useFieldValue {
		ret.SetText(fmt.Sprintf("%s", val.Field(fieldNo).String()))
	} else if v := common.TagVal(t, fieldNo, "form_value"); len(v) > 0 {
		ret.SetText(v)
	}
	return ret
}

// HiddenFieldFromInstance creates and initializes a hidden field based on its name, the reference object instance and field number.
func HiddenFieldFromInstance(val reflect.Value, t reflect.Type, fieldNo int, name string, useFieldValue bool) *Field {
	ret := HiddenField(name)
	if useFieldValue {
		ret.SetValue(fmt.Sprintf("%v", val.Field(fieldNo).Interface()))
	} else if v := common.TagVal(t, fieldNo, "form_value"); len(v) > 0 {
		ret.SetValue(v)
	}
	return ret
}
