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

package forms

import (
	"html/template"
	"strings"

	"github.com/coscms/forms/fields"
	"github.com/webx-top/validation"
)

func ValidationEngine(valid string, f fields.FieldInterface) {
	//for jQuery-Validation-Engine
	validFuncs := strings.Split(valid, ";")
	var validClass string
	for _, v := range validFuncs {
		pos := strings.Index(v, "(")
		var fn string
		if pos > -1 {
			fn = v[0:pos]
		} else {
			fn = v
		}
		switch fn {
		case "required":
			validClass += "," + strings.ToLower(fn)
		case "min", "max":
			val := v[pos+1:]
			val = strings.TrimSuffix(val, ")")
			validClass += "," + strings.ToLower(fn) + "[" + val + "]"
		case "range":
			val := v[pos+1:]
			val = strings.TrimSuffix(val, ")")
			rangeVals := strings.SplitN(val, ",", 2)
			validClass += ",min[" + strings.TrimSpace(rangeVals[0]) + "],max[" + strings.TrimSpace(rangeVals[1]) + "]"
		case "minSize":
			val := v[pos+1:]
			val = strings.TrimSuffix(val, ")")
			validClass += ",minSize[" + val + "]"
		case "maxSize":
			val := v[pos+1:]
			val = strings.TrimSuffix(val, ")")
			validClass += ",maxSize[" + val + "]"
		case "mumeric":
			validClass += ",number"
		case "alphaNumeric":
			validClass += ",custom[onlyLetterNumber]"
		/*
			case "Length":
				validClass += ",length"
			case "Match":
				val := v[pos+1:]
				val = strings.TrimSuffix(val, ")")
				val = strings.Trim(val, "/")
				validClass += ",match[]"
		*/
		case "alphaDash":
			validClass += ",custom[onlyLetterNumber]"
		case "ip":
			validClass += ",custom[ipv4]"
		case "alpha", "email", "base64", "mobile", "tel", "phone":
			validClass += ",custom[" + strings.ToLower(fn) + "]"
		case "zipCode":
			validClass += ",custom[zip]"
		}
	}
	if len(validClass) > 0 {
		validClass = strings.TrimPrefix(validClass, ",")
		validClass = "validate[" + validClass + "]"
		f.AddClass(validClass)
	}
}

func Html5Validate(valid string, f fields.FieldInterface) {
	validFuncs := strings.Split(valid, ";")
	for _, v := range validFuncs {
		pos := strings.Index(v, "(")
		var fn string
		if pos > -1 {
			fn = v[0:pos]
		} else {
			fn = v
		}
		switch fn {
		case "required":
			f.AddTag(strings.ToLower(fn))
		case "min", "max":
			val := v[pos+1:]
			val = strings.TrimSuffix(val, ")")
			f.SetParam(strings.ToLower(fn), val)
		case "range":
			val := v[pos+1:]
			val = strings.TrimSuffix(val, ")")
			rangeVals := strings.SplitN(val, ",", 2)
			f.SetParam("min", strings.TrimSpace(rangeVals[0]))
			f.SetParam("max", strings.TrimSpace(rangeVals[1]))
		case "minSize":
			val := v[pos+1:]
			val = strings.TrimSuffix(val, ")")
			f.SetParam("data-min", val)
		case "maxSize":
			val := v[pos+1:]
			val = strings.TrimSuffix(val, ")")
			f.SetParam("maxlength", val)
			f.SetParam("data-max", val)
		case "numeric":
			f.SetParam("pattern", template.HTML("^\\-?\\d+(\\.\\d+)?$"))
		case "alphaNumeric":
			f.SetParam("pattern", template.HTML("^[a-zA-Z\\d]+$"))
		case "length":
			val := v[pos+1:]
			val = strings.TrimSuffix(val, ")")
			f.SetParam("pattern", ".{"+val+"}")
		case "match":
			val := v[pos+1:]
			val = strings.TrimSuffix(val, ")")
			val = strings.Trim(val, "/")
			f.SetParam("pattern", template.HTML(val))

		case "alphaDash":
			f.SetParam("pattern", template.HTML(validation.DefaultRule.AlphaDash))
		case "ip":
			f.SetParam("pattern", template.HTML(validation.DefaultRule.IPv4))
		case "alpha":
			f.SetParam("pattern", template.HTML("^[a-zA-Z]+$"))
		case "email":
			f.SetParam("pattern", template.HTML(validation.DefaultRule.Email))
		case "base64":
			f.SetParam("pattern", template.HTML(validation.DefaultRule.Base64))
		case "mobile":
			f.SetParam("pattern", template.HTML(validation.DefaultRule.Mobile))
		case "tel":
			f.SetParam("pattern", template.HTML(validation.DefaultRule.Telephone))
		case "phone":
			f.SetParam("pattern", template.HTML(validation.DefaultRule.GetPhone()))
		case "zipCode":
			f.SetParam("pattern", template.HTML(validation.DefaultRule.ZipCode))
		}
	}
}
