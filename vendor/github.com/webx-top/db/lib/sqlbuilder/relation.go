package sqlbuilder

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/webx-top/db"
	"github.com/webx-top/db/lib/reflectx"
	"github.com/webx-top/echo/param"
)

type (
	BuilderChainFunc      func(Selector) Selector
	DBConnFunc            func(name string) db.Database
	StructToTableNameFunc func(data interface{}, retry ...bool) (string, error)
	TableNameFunc         func(fieldInfo *reflectx.FieldInfo, data interface{}) (string, error)
	SQLBuilderFunc        func(fieldInfo *reflectx.FieldInfo, defaults ...SQLBuilder) SQLBuilder
)

const (
	// ForeignKeyIndex 外键名下标
	ForeignKeyIndex = 0
	// RelationKeyIndex 关联键名下标
	RelationKeyIndex = 1
)

func (sel *selector) Relation(name string, fn BuilderChainFunc) Selector {
	return sel.frame(func(sq *selectorQuery) error {
		if sq.relationMap == nil {
			sq.relationMap = map[string]BuilderChainFunc{}
		}
		sq.relationMap[name] = fn
		return nil
	})
}

func eachField(t reflect.Type, fn func(fieldInfo *reflectx.FieldInfo, relations []string, pipes []Pipe) error) error {
	typeMap := mapper.TypeMap(t)
	options, ok := typeMap.Options[`relation`]
	if !ok {
		return nil
	}
	for _, fieldInfo := range options {
		// `db:"-,relation=ForeignKey:RelationKey"`
		// `db:"-,relation=外键名:关联键名|gtZero|eq(field:value)"`
		// `db:"-,relation=ForeignKey:RelationKey,dbconn=link2,columns=colName&colName2"`
		rel, ok := fieldInfo.Options[`relation`]
		if !ok || len(rel) == 0 || rel == `-` {
			continue
		}
		relations := strings.SplitN(rel, `:`, 2)
		if len(relations) != 2 {
			return fmt.Errorf("wrong relation option, length must 2, but get %v. Reference format: `db:\"-,relation=ForeignKey:RelationKey\"`", relations)
		}
		rels := strings.Split(relations[1], `|`)
		var pipes []Pipe
		if len(rels) > 1 {
			relations[1] = rels[0]
			for _, pipeName := range rels[1:] {
				pipe := parsePipe(pipeName)
				if pipe == nil {
					continue
				}
				pipes = append(pipes, pipe)
			}
		}
		err := fn(fieldInfo, relations, pipes)
		if err != nil {
			return err
		}
	}
	return nil
}

type Name_ interface {
	Name_() string
}

func buildCond(refVal reflect.Value, relations []string, pipes []Pipe) interface{} {
	fieldName := relations[ForeignKeyIndex]
	rFieldName := relations[RelationKeyIndex]
	fieldValue := mapper.FieldByName(refVal, rFieldName).Interface()
	if len(pipes) == 0 {
		return db.Cond{
			fieldName: fieldValue,
		}
	}
	for _, pipe := range pipes {
		if fieldValue = pipe(refVal, fieldValue); fieldValue == nil {
			return nil
		}
	}
	var cond interface{}
	if v, y := fieldValue.([]interface{}); y {
		if len(v) == 0 {
			return nil
		}
		cond = db.Cond{
			fieldName: db.In(v),
		}
	} else {
		cond = db.Cond{
			fieldName: fieldValue,
		}
	}
	return cond
}

func buildSelector(fieldInfo *reflectx.FieldInfo, sel Selector, mustColumnName string, hasMustCol *bool, dataTypes *map[string]string) Selector {
	columns, ok := fieldInfo.Options[`columns`] // columns=col1:uint&col2:string&col3:uint64
	if !ok || len(columns) == 0 {
		return sel
	}
	cols := []interface{}{}
	var _hasMustCol bool
	if len(mustColumnName) == 0 {
		_hasMustCol = true
	}
	for _, colName := range strings.Split(columns, `&`) {
		colName = strings.TrimSpace(colName)
		if len(colName) > 0 {
			parts := strings.SplitN(colName, `:`, 2)
			colName = parts[0]
			if !_hasMustCol && colName == mustColumnName {
				_hasMustCol = true
			}
			cols = append(cols, colName)
			if len(parts) == 2 && dataTypes != nil {
				(*dataTypes)[colName] = parts[1]
			}
		}
	}
	if !_hasMustCol {
		cols = append(cols, mustColumnName)
	}
	if hasMustCol != nil {
		*hasMustCol = _hasMustCol
	}
	if len(cols) > 0 {
		return sel.Columns(cols...)
	}
	return sel
}

func convertRelationMapDataType(row reflect.Value, dataTypes map[string]string) {
	if len(dataTypes) == 0 {
		return
	}
	iter := row.MapRange()
	for iter.Next() {
		k := iter.Key()
		v := iter.Value()
		kStr := param.AsString(k.Interface())
		dataType, ok := dataTypes[kStr]
		if ok {
			v = reflect.ValueOf(param.AsType(dataType, v.Interface()))
			row.SetMapIndex(k, v)
		}
	}
}

func deleteRelationMapElement(row reflect.Value, fieldName string) {
	switch rawMap := row.Interface().(type) {
	case map[string]interface{}:
		delete(rawMap, fieldName)
	case *map[string]interface{}:
		delete(*rawMap, fieldName)
	case param.Store:
		rawMap.Delete(fieldName)
	default:
		//fmt.Printf("=======================%T\n", row.Interface())
	}
}

// RelationOne is get the associated relational data for a single piece of data
func RelationOne(builder SQLBuilder, data interface{}, relationMap map[string]BuilderChainFunc) error {
	refVal := reflect.Indirect(reflect.ValueOf(data))
	t := refVal.Type()

	return eachField(t, func(fieldInfo *reflectx.FieldInfo, relations []string, pipes []Pipe) error {
		field := fieldInfo.Field
		name := field.Name
		b := GetSQLBuilder(fieldInfo, builder)
		var foreignModel reflect.Value
		// if field type is slice then one-to-many ,eg: []*Struct
		if field.Type.Kind() == reflect.Slice {
			foreignModel = reflect.New(field.Type)
			foreignIV := foreignModel.Interface()
			table, err := GetTableName(fieldInfo, foreignIV)
			if err != nil {
				return err
			}
			// batch get field values
			// Since the structure is slice, there is no need to new Value
			cond := buildCond(refVal, relations, pipes)
			if cond == nil {
				return nil
			}
			dataTypes := map[string]string{}
			sel := buildSelector(fieldInfo, b.SelectFrom(table).Where(cond), ``, nil, &dataTypes)
			if relationMap != nil {
				if chainFn, ok := relationMap[name]; ok {
					if sel = chainFn(sel); sel == nil {
						return nil
					}
				}
			}
			err = sel.All(foreignIV)
			if err != nil && err != db.ErrNoMoreRows {
				return err
			}
			sliceLen := reflect.Indirect(foreignModel).Len()
			if sliceLen == 0 {
				// If relation data is empty, must set empty slice
				// Otherwise, the JSON result will be null instead of []
				refVal.FieldByName(name).Set(reflect.MakeSlice(field.Type, 0, 0))
			} else if len(dataTypes) == 0 {
				refVal.FieldByName(name).Set(foreignModel.Elem())
			} else {
				childElem := field.Type.Elem()
				if childElem.Kind() == reflect.Ptr {
					childElem = childElem.Elem()
				}
				isMap := childElem.Kind() == reflect.Map
				if !isMap || len(dataTypes) == 0 {
					refVal.FieldByName(name).Set(foreignModel.Elem())
					return nil
				}
				recvVal := reflect.Indirect(foreignModel)
				for n := 0; n < sliceLen; n++ {
					row := recvVal.Index(n)
					convertRelationMapDataType(row, dataTypes)
				}
				refVal.FieldByName(name).Set(foreignModel.Elem())
			}
		} else {
			// If field type is struct the one-to-one,eg: *Struct
			if field.Type.Kind() == reflect.Ptr {
				foreignModel = reflect.New(field.Type.Elem())
			} else {
				foreignModel = reflect.New(field.Type)
			}
			foreignIV := foreignModel.Interface()

			table, err := GetTableName(fieldInfo, foreignIV)
			if err != nil {
				return err
			}

			cond := buildCond(refVal, relations, pipes)
			if cond == nil {
				return nil
			}
			dataTypes := map[string]string{}
			sel := buildSelector(fieldInfo, b.SelectFrom(table).Where(cond), ``, nil, &dataTypes)
			if relationMap != nil {
				if chainFn, ok := relationMap[name]; ok {
					if sel = chainFn(sel); sel == nil {
						return nil
					}
				}
			}
			err = sel.One(foreignIV)
			// If one-to-one NoRows is not an error that needs to be terminated
			if err != nil {
				if err != db.ErrNoMoreRows {
					return err
				}
			} else {
				if field.Type.Kind() != reflect.Ptr {
					foreignModel = foreignModel.Elem()
				}
				if foreignModel.Kind() == reflect.Map {
					convertRelationMapDataType(foreignModel, dataTypes)
				}
				refVal.FieldByName(name).Set(foreignModel)
			}
		}
		return nil
	})
}

// RelationAll is gets the associated relational data for multiple pieces of data
func RelationAll(builder SQLBuilder, data interface{}, relationMap map[string]BuilderChainFunc) error {
	refVal := reflect.Indirect(reflect.ValueOf(data))

	l := refVal.Len()

	if l == 0 {
		return nil
	}

	// get the struct field in slice
	t := reflect.Indirect(refVal.Index(0)).Type()

	return eachField(t, func(fieldInfo *reflectx.FieldInfo, relations []string, pipes []Pipe) error {
		field := fieldInfo.Field
		name := field.Name
		relVals := make([]interface{}, 0)
		relValsMap := make(map[interface{}]struct{})
		relValsMapx := make(map[int][]interface{})
		fieldName := relations[ForeignKeyIndex]
		rFieldName := relations[RelationKeyIndex]
		relValKind := mapper.FieldByName(refVal.Index(0), rFieldName).Kind()
		// get relation field values and unique
		if len(pipes) == 0 {
			for j := 0; j < l; j++ {
				v := mapper.FieldByName(refVal.Index(j), rFieldName).Interface()
				relValsMap[v] = struct{}{}
			}
		} else {
			for j := 0; j < l; j++ {
				row := refVal.Index(j)
				v := mapper.FieldByName(row, rFieldName).Interface()
				for _, pipe := range pipes {
					if v = pipe(row, v); v == nil {
						break
					}
				}
				if v == nil {
					continue
				}
				if vs, ok := v.([]interface{}); ok {
					if _, ok := relValsMapx[j]; !ok {
						relValsMapx[j] = []interface{}{}
					}
					for _, vv := range vs {
						relValsMap[vv] = struct{}{}
						relValsMapx[j] = append(relValsMapx[j], vv)
					}
					continue
				}
				relValsMap[v] = struct{}{}
			}
		}
		if len(relValsMap) == 0 {
			return nil
		}
		for k := range relValsMap {
			relVals = append(relVals, k)
		}

		b := GetSQLBuilder(fieldInfo, builder)
		var foreignModel reflect.Value
		// if field type is slice then one to many ,eg: []*Struct
		if field.Type.Kind() == reflect.Slice {
			childElem := field.Type.Elem()
			if childElem.Kind() == reflect.Ptr {
				childElem = childElem.Elem()
			}
			isMap := childElem.Kind() == reflect.Map
			foreignModel = reflect.New(field.Type)
			foreignIV := foreignModel.Interface()
			table, err := GetTableName(fieldInfo, foreignIV)
			if err != nil {
				return err
			}
			var hasMustCol bool
			dataTypes := map[string]string{}
			// batch get field values
			// Since the structure is slice, there is no need to new Value
			sel := buildSelector(fieldInfo, b.SelectFrom(table).Where(db.Cond{
				fieldName: db.In(relVals),
			}), fieldName, &hasMustCol, &dataTypes)
			if relationMap != nil {
				if chainFn, ok := relationMap[name]; ok {
					if sel = chainFn(sel); sel == nil {
						return nil
					}
				}
			}
			err = sel.All(foreignIV)
			if err != nil && err != db.ErrNoMoreRows {
				return err
			}

			fmap := make(map[interface{}]reflect.Value)

			var nameV reflect.Value
			if isMap {
				nameV = reflect.ValueOf(fieldName)
			}
			// Combine relation data as a one-to-many relation
			// For example, if there are multiple images under an article
			// we use the article ID to associate the images, map[1][]*Images
			mlen := reflect.Indirect(foreignModel).Len()
			recvVal := reflect.Indirect(foreignModel)
			var fmapKeyKind reflect.Kind
			for n := 0; n < mlen; n++ {
				row := recvVal.Index(n)
				if isMap {
					fid := row.MapIndex(nameV)
					if !fid.CanInterface() {
						continue
					}
					val := param.AsType(relValKind.String(), fid.Interface())
					if !hasMustCol {
						deleteRelationMapElement(row, fieldName)
					}
					convertRelationMapDataType(row, dataTypes)
					if _, has := fmap[val]; !has {
						fmap[val] = reflect.New(reflect.SliceOf(field.Type.Elem())).Elem()
					}
					fmap[val] = reflect.Append(fmap[val], row)
					continue
				}
				fid := mapper.FieldByName(row, fieldName)
				fv := fid.Interface()
				if _, has := fmap[fv]; !has {
					fmap[fv] = reflect.New(reflect.SliceOf(field.Type.Elem())).Elem()
				}
				if fmapKeyKind == reflect.Invalid {
					fmapKeyKind = fid.Type().Kind()
				}
				fmap[fv] = reflect.Append(fmap[fv], row)
				if !hasMustCol {
					fid.Set(reflect.Zero(fid.Type()))
				}
			}
			var ft reflect.Kind
			if mlen > 0 && foreignModel.Type().Kind() == reflect.Struct {
				ft = mapper.FieldByName(reflect.Indirect(foreignModel).Index(0), fieldName).Kind()
			}
			needConversion := relValKind != ft && ft != reflect.Invalid
			// Set the result to the model
			for j := 0; j < l; j++ {
				v := refVal.Index(j)
				fid := mapper.FieldByName(v, rFieldName)
				val := fid.Interface()
				if needConversion {
					val = param.AsType(ft.String(), val)
				}
				if value, has := fmap[val]; has {
					reflect.Indirect(v).FieldByName(name).Set(value)
				} else {
					if idxList, ok := relValsMapx[j]; ok {
						slicev := reflect.New(reflect.SliceOf(field.Type.Elem())).Elem()
						for _, _v := range idxList {
							if fmapKeyKind != reflect.Invalid {
								_v = param.AsType(fmapKeyKind.String(), _v)
								if value, has := fmap[_v]; has {
									slicev = reflect.AppendSlice(slicev, value)
								}
							} else {
								_v = param.AsType(ft.String(), _v)
								if value, has := fmap[_v]; has {
									slicev = reflect.AppendSlice(slicev, value)
								}
							}
						}
						reflect.Indirect(v).FieldByName(name).Set(slicev)
						continue
					}
					// If relation data is empty, must set empty slice
					// Otherwise, the JSON result will be null instead of []
					reflect.Indirect(v).FieldByName(name).Set(reflect.MakeSlice(field.Type, 0, 0))
				}
			}
		} else {
			var sliceT reflect.Type
			var isMap bool
			// If field type is struct the one to one,eg: *Struct
			if field.Type.Kind() == reflect.Ptr {
				fieldT := field.Type.Elem()
				isMap = fieldT.Kind() == reflect.Map
				foreignModel = reflect.New(fieldT)
				sliceT = reflect.SliceOf(foreignModel.Type())
			} else {
				fieldT := field.Type
				isMap = fieldT.Kind() == reflect.Map
				foreignModel = reflect.New(fieldT)
				sliceT = reflect.SliceOf(foreignModel.Type().Elem())
			}

			// Batch get field values, but must new slice []*Struct
			fi := reflect.New(sliceT)
			foreignIV := fi.Interface()

			table, err := GetTableName(fieldInfo, foreignIV)
			if err != nil {
				return err
			}
			var hasMustCol bool
			dataTypes := map[string]string{}
			sel := buildSelector(fieldInfo, b.SelectFrom(table).Where(db.Cond{
				fieldName: db.In(relVals),
			}), fieldName, &hasMustCol, &dataTypes)
			if relationMap != nil {
				if chainFn, ok := relationMap[name]; ok {
					if sel = chainFn(sel); sel == nil {
						return nil
					}
				}
			}
			err = sel.All(foreignIV)
			if err != nil && err != db.ErrNoMoreRows {
				return err
			}

			// Combine relation data as a one-to-one relation
			fmap := make(map[interface{}]reflect.Value)
			fval := reflect.Indirect(fi)
			mlen := fval.Len()
			var nameV reflect.Value
			if isMap {
				nameV = reflect.ValueOf(fieldName)
			}
			for n := 0; n < mlen; n++ {
				row := fval.Index(n)
				if isMap {
					fid := row.MapIndex(nameV)
					if !fid.CanInterface() {
						continue
					}
					val := param.AsType(relValKind.String(), fid.Interface())
					if !hasMustCol {
						deleteRelationMapElement(row, fieldName)
					}
					convertRelationMapDataType(row, dataTypes)
					fmap[val] = row
					continue
				}
				fid := mapper.FieldByName(row, fieldName)
				fmap[fid.Interface()] = row
				if !hasMustCol {
					fid.Set(reflect.Zero(fid.Type()))
				}
			}
			var ft reflect.Kind
			if mlen > 0 && !isMap {
				ft = mapper.FieldByName(fval.Index(0), fieldName).Kind()
			}
			needConversion := relValKind != ft && ft != reflect.Invalid
			// Set the result to the model
			for j := 0; j < l; j++ {
				v := refVal.Index(j)
				fid := mapper.FieldByName(v, rFieldName)
				val := fid.Interface()
				if needConversion {
					val = param.AsType(ft.String(), val)
				}
				if value, has := fmap[val]; has {
					reflect.Indirect(v).FieldByName(name).Set(value)
				}
			}
		}

		return nil
	})
}
