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

	"github.com/coscms/forms/common"
)

// RangeField creates a default range field with the provided name. Min, max and step parameters define the expected behavior
// of the HTML field.
func RangeField(name string, min, max, step int) *Field {
	ret := FieldWithType(name, common.RANGE)
	ret.SetParam("min", fmt.Sprintf("%d", min))
	ret.SetParam("max", fmt.Sprintf("%d", max))
	ret.SetParam("step", fmt.Sprintf("%d", step))
	return ret
}

// NumberField craetes a default number field with the provided name.
func NumberField(name string) *Field {
	ret := FieldWithType(name, common.NUMBER)
	return ret
}

// NumberFieldFromInstance creates and initializes a number field based on its name, the reference object instance and field number.
// This method looks for "form_min", "form_max" and "form_value" tags to add additional parameters to the field.
func NumberFieldFromInstance(val reflect.Value, t reflect.Type, fieldNo int, name string, useFieldValue bool) *Field {
	ret := NumberField(name)
	// check tags
	if v := common.TagVal(t, fieldNo, "form_min"); v != "" {
		ret.SetParam("min", v)
	}
	if v := common.TagVal(t, fieldNo, "form_max"); v != "" {
		ret.SetParam("max", v)
	}
	if v := common.TagVal(t, fieldNo, "form_step"); v != "" {
		ret.SetParam("step", v)
	}
	ret.SetValue(defaultValue(val, t, fieldNo, useFieldValue))
	return ret
}

// RangeFieldFromInstance creates and initializes a range field based on its name, the reference object instance and field number.
// This method looks for "form_min", "form_max", "form_step" and "form_value" tags to add additional parameters to the field.
func RangeFieldFromInstance(val reflect.Value, t reflect.Type, fieldNo int, name string, useFieldValue bool) *Field {
	ret := RangeField(name, 0, 10, 1)
	// check tags
	if v := common.TagVal(t, fieldNo, "form_min"); v != "" {
		ret.SetParam("min", v)
	}
	if v := common.TagVal(t, fieldNo, "form_max"); v != "" {
		ret.SetParam("max", v)
	}
	if v := common.TagVal(t, fieldNo, "form_step"); v != "" {
		ret.SetParam("step", v)
	}
	ret.SetValue(defaultValue(val, t, fieldNo, useFieldValue))
	return ret
}
