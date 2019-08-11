package factory

import (
	"encoding/json"
	"fmt"
	"html/template"
	"strconv"
	"strings"
	"time"

	"github.com/webx-top/com"
	"github.com/webx-top/echo"
)

// HTML 表单输入框
func (f *FieldInfo) HTML(value interface{}, options echo.H) template.HTML {
	var input string
	name := com.LowerCaseFirst(f.GoName)
	val := fmt.Sprint(value)
	required := options.Bool(`required`)
	switch f.DataType {
	case `enum`:
		labels := options.Store(`optionLabels`)
		if len(val) == 0 && len(f.DefaultValue) > 0 {
			val = f.DefaultValue
		}
		for index, option := range f.Options {
			attrs := HTMLAttrs{}
			if (len(val) == 0 && index == 0) || option == val {
				attrs.Add(`checked`, `checked`)
			}
			label := option
			if v := labels.String(option); len(v) > 0 {
				label = v
			}
			data := echo.H{
				`theme`:     `primary`,
				`inline`:    true,
				`value`:     option,
				`name`:      name,
				`id`:        name + `-` + option,
				`attrs`:     attrs,
				`label`:     label,
				`helpBlock`: ``,
			}
			data.DeepMerge(options)
			input += DefaultHTMLTmpl.ToInput(`radio`, data)
		}
	case `set`:
		labels := options.Store(`optionLabels`)
		if len(val) == 0 && len(f.DefaultValue) > 0 {
			val = f.DefaultValue
		}
		var vals []string
		if len(val) > 0 {
			switch val[0] {
			case '[':
				json.Unmarshal(com.Str2bytes(val), &vals)
			default:
				vals = strings.Split(val, `,`)
			}
		}
		for index, option := range f.Options {
			attrs := HTMLAttrs{}
			if com.InSlice(option, vals) {
				attrs.Add(`checked`)
			}
			label := option
			if v := labels.String(option); len(v) > 0 {
				label = v
			}
			idValue := name + `-` + option + `-` + fmt.Sprint(index)
			data := echo.H{
				`theme`:     `primary`,
				`inline`:    true,
				`value`:     option,
				`name`:      name,
				`id`:        idValue,
				`attrs`:     attrs,
				`label`:     label,
				`helpBlock`: ``,
			}
			data.DeepMerge(options)
			input += DefaultHTMLTmpl.ToInput(`checkbox`, data)
		}
	case `date`:
		attrs := f.HTMLAttrBuilder(required)
		if isInteger(f.GoName) {
			if i, e := strconv.ParseInt(val, 10, 64); e == nil {
				val = time.Unix(i, 0).Format(`2006-01-02`)
			} else {
				val = ``
			}
		}
		data := echo.H{
			`type`:      `date`,
			`value`:     val,
			`name`:      name,
			`attrs`:     attrs,
			`helpBlock`: ``,
		}
		data.DeepMerge(options)
		input = DefaultHTMLTmpl.ToInput(`text`, data)
	case `time`:
		attrs := f.HTMLAttrBuilder(required)
		if isInteger(f.GoName) {
			if i, e := strconv.ParseInt(val, 10, 64); e == nil {
				val = time.Unix(i, 0).Format(`15:04:05`)
			} else {
				val = ``
			}
		}
		data := echo.H{
			`type`:      `time`,
			`value`:     val,
			`name`:      name,
			`attrs`:     attrs,
			`helpBlock`: ``,
		}
		data.DeepMerge(options)
		input = DefaultHTMLTmpl.ToInput(`text`, data)
	case `datetime`:
		attrs := f.HTMLAttrBuilder(required)
		val := fmt.Sprint(value)
		if isInteger(f.GoName) {
			if i, e := strconv.ParseInt(val, 10, 64); e == nil {
				val = time.Unix(i, 0).Format(`2006-01-02 15:04:05`)
			} else {
				val = ``
			}
		}
		data := echo.H{
			`type`:      `datetime`,
			`value`:     val,
			`name`:      name,
			`attrs`:     attrs,
			`helpBlock`: ``,
		}
		data.DeepMerge(options)
		input = DefaultHTMLTmpl.ToInput(`text`, data)
	case `text`, `longtext`, `tinytext`, `mediumtext`:
		attrs := f.HTMLAttrBuilder(required)
		data := echo.H{
			`value`:     val,
			`name`:      name,
			`attrs`:     attrs,
			`helpBlock`: ``,
		}
		data.DeepMerge(options)
		input = DefaultHTMLTmpl.ToInput(`textarea`, data)
	default:
		attrs := f.HTMLAttrBuilder(required)
		switch f.GoType {
		case `int`, `int64`, `uint`, `uint64`:
			data := echo.H{
				`type`:      `number`,
				`value`:     val,
				`name`:      name,
				`attrs`:     attrs,
				`helpBlock`: ``,
			}
			attrs.Add(`step`, `1`)
			data.DeepMerge(options)
			input = DefaultHTMLTmpl.ToInput(`text`, data)
		case `float64`, `float32`:
			if f.Precision > 0 {
				attrs.Add(`step`, fmt.Sprintf(`0.%0*d`, f.Precision, 1))
			} else {
				attrs.Add(`step`, "1")
			}
			data := echo.H{
				`type`:      `number`,
				`value`:     val,
				`name`:      name,
				`attrs`:     attrs,
				`helpBlock`: ``,
			}
			data.DeepMerge(options)
			input = DefaultHTMLTmpl.ToInput(`text`, data)
		case `bool`:
			labels := options.Store(`optionLabels`)
			if len(val) == 0 && len(f.DefaultValue) > 0 {
				val = f.DefaultValue
			}
			for index, option := range []string{`1`, `0`} {
				attrs := HTMLAttrs{}
				if (len(val) == 0 && index == 0) || option == val {
					attrs.Add(`checked`, `checked`)
				}
				label := option
				if v := labels.String(option); len(v) > 0 {
					label = v
				}
				data := echo.H{
					`theme`:     `primary`,
					`inline`:    true,
					`value`:     option,
					`name`:      name,
					`id`:        name + `-` + option,
					`attrs`:     attrs,
					`label`:     label,
					`helpBlock`: ``,
				}
				data.DeepMerge(options)
				input += DefaultHTMLTmpl.ToInput(`radio`, data)
			}
		case `[]byte`:
			data := echo.H{
				`type`:      `file`,
				`value`:     val,
				`name`:      name,
				`attrs`:     attrs,
				`helpBlock`: ``,
			}
			data.DeepMerge(options)
			input = DefaultHTMLTmpl.ToInput(`text`, data)
		case `string`:
			fallthrough
		default:
			data := echo.H{
				`type`:      `text`,
				`value`:     val,
				`name`:      name,
				`attrs`:     attrs,
				`helpBlock`: ``,
			}
			data.DeepMerge(options)
			input = DefaultHTMLTmpl.ToInput(`text`, data)
		}
	}
	return template.HTML(input)
}

// FormGroup 表单组，带标签(label)
func (f *FieldInfo) FormGroup(value interface{}, options echo.H, inputAndLabelCols ...int) template.HTML {
	labelCols := 2
	inputCols := 8
	switch len(inputAndLabelCols) {
	case 2:
		labelCols = inputAndLabelCols[1]
		inputCols = inputAndLabelCols[0]
	case 1:
		inputCols = inputAndLabelCols[0]
		labelCols = 10 - inputCols
	}
	var star string
	required := options.Bool(`required`)
	if required {
		star = DefaultHTMLTmpl.Required
	}
	data := echo.H{
		`labelCols`:   labelCols,
		`inputCols`:   inputCols,
		`label`:       f.Comment,
		`labelSuffix`: star,
		`input`:       f.HTML(value, options),
	}
	data.DeepMerge(options)
	return template.HTML(DefaultHTMLTmpl.ToGroup(data))
}
