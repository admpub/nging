package reflectx

import (
	"reflect"
	"strings"
)

// ReflectSliceElem *[]struct / []struct / *[]*struct / []*struct
func ReflectSliceElem(bean interface{}) reflect.Type {
	itemV := reflect.ValueOf(bean)
	itemT := itemV.Type()
	if itemT.Kind() == reflect.Ptr {
		itemT = itemT.Elem()
	}
	if itemT.Kind() == reflect.Slice {
		itemT = itemT.Elem()
		if itemT.Kind() == reflect.Ptr {
			itemT = itemT.Elem()
		}
	}
	return itemT
}

// StructMap returns a mapping of field strings to int slices representing
// the traversal down the struct to reach the field.
func (m *Mapper) StructMap(bean interface{}) *StructMap {
	return m.TypeMap(ReflectSliceElem(bean))
}

func (f *FieldInfo) Find(field string, isStructField bool) (found *FieldInfo) {
	if isStructField {
		field = strings.Title(field)
	}
	for _, fieldInfo := range f.Children {
		if fieldInfo == nil {
			continue
		}
		var equal bool
		if isStructField {
			equal = fieldInfo.Field.Name == field
		} else {
			equal = fieldInfo.Name == field
		}
		if equal {
			found = fieldInfo
			break
		}
	}
	return
}

func (f *FieldInfo) FindTableField(field string, isStructField bool, aliasOptionName string) (tableField string, found *FieldInfo) {
	if isStructField {
		field = strings.Title(field)
	}
	for _, fieldInfo := range f.Children {
		if fieldInfo == nil {
			continue
		}
		var equal bool
		if isStructField {
			equal = fieldInfo.Field.Name == field
		} else {
			equal = fieldInfo.Name == field
		}
		//fmt.Println(fieldInfo.Field.Name, `[=]`, field)
		if equal {
			tableField = getTableFieldName(fieldInfo, aliasOptionName)
			found = fieldInfo
			break
		}
	}
	return
}

// Find Find("user.profile")
func (f StructMap) Find(fieldPath string, isStructField bool) (tree *FieldInfo, exists bool) {
	tree = f.Tree
	for _, field := range strings.Split(fieldPath, `.`) {
		if len(field) == 0 {
			return nil, false
		}
		tree = tree.Find(field, isStructField)
		if tree == nil {
			return nil, false
		}
		exists = true
	}
	return
}

func getTableFieldName(fieldInfo *FieldInfo, aliasOptionName string) string {
	alias, _ := fieldInfo.Options[aliasOptionName]
	if len(alias) == 0 {
		if len(fieldInfo.Name) > 0 {
			alias = fieldInfo.Name
		} else {
			alias = fieldInfo.Field.Name
		}
	}
	return alias
}

// FindTableField Find("User.Profile")
func (f StructMap) FindTableField(fieldPath string, isStructField bool, aliasOptionNames ...string) (tableFieldPath string, exists bool) {
	tree := f.Tree
	var aliasOptionName string
	if len(aliasOptionNames) > 0 {
		aliasOptionName = aliasOptionNames[0]
	}
	if len(aliasOptionName) == 0 {
		aliasOptionName = `alias`
	}
	for _, field := range strings.Split(fieldPath, `.`) {
		if len(field) == 0 {
			return strings.TrimPrefix(tableFieldPath, `.`), false
		}
		var tableField string
		tableField, tree = tree.FindTableField(field, isStructField, aliasOptionName)
		if tree == nil {
			return strings.TrimPrefix(tableFieldPath, `.`), false
		}
		tableFieldPath += `.` + tableField
		exists = true
	}
	tableFieldPath = strings.TrimPrefix(tableFieldPath, `.`)
	return
}

type FindResult struct {
	RawData   interface{}
	FieldInfo *FieldInfo
	Parents   []*FieldInfo
	RawPath   []string
}

func (f *FindResult) Parent(index int) *FieldInfo {
	if index >= len(f.Parents) {
		return nil
	}
	return f.Parents[index]
}

func (f StructMap) FindTableFieldByMap(fieldPaths map[string]map[string]interface{}, isStructField bool, aliasOptionNames ...string) (tableFieldPaths map[string]*FindResult, pk []string) {
	tree := f.Tree
	var aliasOptionName string
	if len(aliasOptionNames) > 0 {
		aliasOptionName = aliasOptionNames[0]
	}
	if len(aliasOptionName) == 0 {
		aliasOptionName = `alias`
	}
	tableFieldPaths = map[string]*FindResult{}
	for parent, fields := range fieldPaths {
		if len(parent) == 0 {
			continue
		}
		parentRaw := parent
		parentTree := tree
		parent, parentTree = parentTree.FindTableField(parent, isStructField, aliasOptionName)
		if parentTree == nil {
			continue
		}
		for field, rawData := range fields {
			if len(field) == 0 {
				tableFieldPaths[parent] = &FindResult{
					RawData:   rawData,
					FieldInfo: parentTree,
					Parents:   []*FieldInfo{},
					RawPath:   []string{parentRaw, field},
				}
				continue
			}
			fieldRaw := field
			var info *FieldInfo
			field, info = parentTree.FindTableField(field, isStructField, aliasOptionName)
			if info == nil {
				continue
			}
			if _, exists := info.Options[`pk`]; exists {
				pk = append(pk, parent+`.`+field)
			}
			tableFieldPaths[parent+`.`+field] = &FindResult{
				RawData:   rawData,
				FieldInfo: info,
				Parents:   []*FieldInfo{parentTree},
				RawPath:   []string{parentRaw, fieldRaw},
			}
		}
	}
	return
}
