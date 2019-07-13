package reflectx

import (
	"reflect"
	"strings"
)

// StructMap returns a mapping of field strings to int slices representing
// the traversal down the struct to reach the field.
func (m *Mapper) StructMap(bean interface{}) *StructMap {
	return m.TypeMap(reflect.ValueOf(bean).Type())
}

// Find Find("user.profile")
func (f StructMap) Find(fieldPath string) (tree *FieldInfo, exists bool) {
	tree = f.Tree
	for _, field := range strings.Split(fieldPath, `.`) {
		field = strings.Title(field)
		var found bool
		for _, fieldInfo := range tree.Children {
			if fieldInfo.Field.Name == field {
				tree = fieldInfo
				found = true
				break
			}
		}
		if !found {
			return nil, false
		}
		exists = true
	}
	return
}
