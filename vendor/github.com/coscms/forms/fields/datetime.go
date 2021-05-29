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
	"reflect"
	"time"

	"github.com/coscms/forms/common"
)

// Datetime format string to convert from time.Time objects to HTML fields and viceversa.
const (
	DATETIME_FORMAT = "2006-01-02 15:05"
	DATE_FORMAT     = "2006-01-02"
	TIME_FORMAT     = "15:05"
)

func ConvertTime(v interface{}) (time.Time, bool) {
	t, ok := v.(time.Time)
	var isEmpty bool
	if !ok {
		var timestamp int64
		switch i := v.(type) {
		case int:
			timestamp = int64(i)
		case int64:
			timestamp = i
		}
		if timestamp > 0 {
			t = time.Unix(timestamp, 0)
		} else {
			isEmpty = true
		}
	}
	return t, isEmpty
}

// DatetimeField creates a default datetime input field with the given name.
func DatetimeField(name string) *Field {
	return FieldWithType(name, common.DATETIME)
}

// DateField creates a default date input field with the given name.
func DateField(name string) *Field {
	return FieldWithType(name, common.DATE)
}

// TimeField creates a default time input field with the given name.
func TimeField(name string) *Field {
	return FieldWithType(name, common.TIME)
}

// DatetimeFieldFromInstance creates and initializes a datetime field based on its name, the reference object instance and field number.
// This method looks for "form_min", "form_max" and "form_value" tags to add additional parameters to the field.
func DatetimeFieldFromInstance(val reflect.Value, t reflect.Type, fieldNo int, name string, useFieldValue bool) *Field {
	ret := DatetimeField(name)
	dateFormat := DATETIME_FORMAT
	if v := common.TagVal(t, fieldNo, "form_format"); len(v) > 0 {
		dateFormat = v
	}
	ret.Format = dateFormat
	// check tags
	if v := common.TagVal(t, fieldNo, "form_min"); len(v) > 0 {
		if !validateDateformat(v, dateFormat) {
			panic("Invalid date value (min) for field: " + name)
		}
		ret.SetParam("min", v)
	}
	if v := common.TagVal(t, fieldNo, "form_max"); len(v) > 0 {
		if !validateDateformat(v, dateFormat) {
			panic("Invalid date value (max) for field: " + name)
		}
		ret.SetParam("max", v)
	}

	if useFieldValue {
		if vt, isEmpty := ConvertTime(val.Field(fieldNo).Interface()); !vt.IsZero() {
			ret.SetValue(vt.Format(dateFormat))
		} else if isEmpty {
			ret.SetValue(``)
		}
	} else if v := common.TagVal(t, fieldNo, "form_value"); len(v) > 0 {
		ret.SetValue(v)
	}
	return ret
}

// DateFieldFromInstance creates and initializes a date field based on its name, the reference object instance and field number.
// This method looks for "form_min", "form_max" and "form_value" tags to add additional parameters to the field.
func DateFieldFromInstance(val reflect.Value, t reflect.Type, fieldNo int, name string, useFieldValue bool) *Field {
	ret := DateField(name)
	dateFormat := DATE_FORMAT
	if v := common.TagVal(t, fieldNo, "form_format"); len(v) > 0 {
		dateFormat = v
	}
	ret.Format = dateFormat
	// check tags
	if v := common.TagVal(t, fieldNo, "form_min"); len(v) > 0 {
		if !validateDateformat(v, dateFormat) {
			panic("Invalid date value (min) for field: " + name)
		}
		ret.SetParam("min", v)
	}
	if v := common.TagVal(t, fieldNo, "form_max"); len(v) > 0 {
		if !validateDateformat(v, dateFormat) {
			panic("Invalid date value (max) for field: " + name)
		}
		ret.SetParam("max", v)
	}

	if useFieldValue {
		if vt, isEmpty := ConvertTime(val.Field(fieldNo).Interface()); !vt.IsZero() {
			ret.SetValue(vt.Format(dateFormat))
		} else if isEmpty {
			ret.SetValue(``)
		}
	} else if v := common.TagVal(t, fieldNo, "form_value"); len(v) > 0 {
		ret.SetValue(v)
	}
	return ret
}

// TimeFieldFromInstance creates and initializes a time field based on its name, the reference object instance and field number.
// This method looks for "form_min", "form_max" and "form_value" tags to add additional parameters to the field.
func TimeFieldFromInstance(val reflect.Value, t reflect.Type, fieldNo int, name string, useFieldValue bool) *Field {
	ret := TimeField(name)
	dateFormat := TIME_FORMAT
	if v := common.TagVal(t, fieldNo, "form_format"); len(v) > 0 {
		dateFormat = v
	}
	ret.Format = dateFormat
	// check tags
	if v := common.TagVal(t, fieldNo, "form_min"); len(v) > 0 {
		if !validateDateformat(v, dateFormat) {
			panic("Invalid time value (min) for field: " + name)
		}
		ret.SetParam("min", v)
	}
	if v := common.TagVal(t, fieldNo, "form_max"); len(v) > 0 {
		if !validateDateformat(v, dateFormat) {
			panic("Invalid time value (max) for field: " + name)
		}
		ret.SetParam("max", v)
	}
	if useFieldValue {
		if v, isEmpty := ConvertTime(val.Field(fieldNo).Interface()); !v.IsZero() {
			ret.SetValue(v.Format(dateFormat))
		} else if isEmpty {
			ret.SetValue(``)
		}
	} else if v := common.TagVal(t, fieldNo, "form_value"); len(v) > 0 {
		ret.SetValue(v)
	}
	return ret
}

func validateDateformat(v string, format string) bool {
	_, err := time.Parse(format, v)
	return err == nil
}

func validateDatetime(v string) bool {
	_, err := time.Parse(DATETIME_FORMAT, v)
	return err == nil
}

func validateDate(v string) bool {
	_, err := time.Parse(DATE_FORMAT, v)
	return err == nil
}

func validateTime(v string) bool {
	_, err := time.Parse(TIME_FORMAT, v)
	return err == nil
}
