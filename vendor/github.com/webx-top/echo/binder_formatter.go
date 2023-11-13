package echo

import "unicode"

type (
	//FieldNameFormatter 结构体字段值映射到表单时，结构体字段名称格式化处理
	FieldNameFormatter func(topName, fieldName string) string
)

var (
	//DefaultFieldNameFormatter 默认格式化函数(struct->form)
	DefaultFieldNameFormatter FieldNameFormatter = func(topName, fieldName string) string {
		var fName string
		if len(topName) == 0 {
			fName = fieldName
		} else {
			fName = topName + "." + fieldName
		}
		return fName
	}
	//ArrayFieldNameFormatter 格式化函数(struct->form)
	ArrayFieldNameFormatter FieldNameFormatter = func(topName, fieldName string) string {
		var fName string
		if len(topName) == 0 {
			fName = fieldName
		} else {
			fName = topName + `[` + fieldName + `]`
		}
		return fName
	}
	//LowerCaseFirstLetter 小写首字母(struct->form)
	LowerCaseFirstLetter FieldNameFormatter = func(topName, fieldName string) string {
		var fName string
		s := []rune(fieldName)
		if len(s) > 0 {
			s[0] = unicode.ToLower(s[0])
			fieldName = string(s)
		}
		if len(topName) == 0 {
			fName = fieldName
		} else {
			fName = topName + "." + fieldName
		}
		return fName
	}
)

func MakeArrayFieldNameFormatter(keyFormatter func(string) string) FieldNameFormatter {
	return func(topName, fieldName string) string {
		if keyFormatter != nil {
			fieldName = keyFormatter(fieldName)
		}
		return ArrayFieldNameFormatter(topName, fieldName)
	}
}
