package copier

import "reflect"

// DeepFindFields deep finds all fields
func DeepFindFields(reflectType reflect.Type, reflectValue reflect.Value, prefix string, needInitFields map[string]struct{}) []reflect.StructField {
	return deepFieldsEx(reflectType, reflectValue, prefix, needInitFields)
}

func deepFieldsEx(reflectType reflect.Type, reflectValue reflect.Value, prefix string, needInitFields map[string]struct{}) []reflect.StructField {
	var fields []reflect.StructField

	if reflectType, _ = indirectType(reflectType); reflectType.Kind() == reflect.Struct {
		for i := 0; i < reflectType.NumField(); i++ {
			v := reflectType.Field(i)
			if v.Anonymous {
				value := indirect(reflectValue).Field(i)
				if value.Kind() == reflect.Ptr {
					if value.IsNil() {
						continue
					}
					if needInitFields != nil {
						needInitFields[prefix+v.Name] = struct{}{}
					}
				}
				prefix += v.Name + `.`
				fields = append(fields, deepFieldsEx(v.Type, value, prefix, needInitFields)...)
			} else {
				fields = append(fields, v)
			}
		}
	}

	return fields
}

// AllNilFields 初始化所有nil字段
var AllNilFields = map[string]struct{}{
	`*`: struct{}{},
}

// InitNilFields initializes nil fields
func InitNilFields(reflectType reflect.Type, reflectValue reflect.Value, prefix string, needInitFields map[string]struct{}) {
	if needInitFields == nil {
		return
	}
	reflectType, _ = indirectType(reflectType)
	if reflectType.Kind() != reflect.Struct {
		return
	}
	for i := 0; i < reflectType.NumField(); i++ {
		v := reflectType.Field(i)
		_, ok := needInitFields[prefix+v.Name]
		if !ok {
			_, ok = needInitFields[`*`]
			if !ok {
				continue
			}
		}
		if !v.Anonymous {
			continue
		}
		value := indirect(reflectValue).Field(i)
		if value.Kind() != reflect.Ptr {
			continue
		}
		if !value.IsNil() {
			continue
		}
		if !value.CanSet() {
			continue
		}
		value.Set(reflect.New(v.Type.Elem()))
		prefix += v.Name + `.`
		InitNilFields(v.Type, value, prefix, needInitFields)
	}
}
