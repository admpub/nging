package factory

import (
	"fmt"
	"html"
	"html/template"

	"github.com/webx-top/com"
	"github.com/webx-top/echo"
)

// ListCol 列表时列值显示方式
func (f *FieldInfo) ListCol(index int, value interface{}, options echo.H) template.HTML {
	var input string
	val := fmt.Sprint(value)
	switch f.DataType {
	case `enum`:
		openValue := `Y`
		if len(f.Options) == 2 && com.InSlice(openValue, f.Options) {
			name := com.LowerCaseFirst(f.GoName)
			attrs := HTMLAttrs{}
			checked := val == openValue
			label := val
			data := echo.H{
				`field`:     f.Name,
				`index`:     index,
				`label`:     label,
				`name`:      name,
				`value`:     val,
				`attrs`:     attrs,
				`openValue`: openValue,
				`checked`:   checked,
			}
			data.DeepMerge(options)
			input = GetHTMLTmpl(options).ToListCol(`switch`, data)
			return template.HTML(input)
		}
		fallthrough
	default:
		input = html.EscapeString(val)
	}
	return template.HTML(input)
}
