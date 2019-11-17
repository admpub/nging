package sqlbuilder

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/admpub/errors"

	"github.com/webx-top/com"
	"github.com/webx-top/db"
	"github.com/webx-top/echo/param"
)

type BuilderChainFunc func(Selector) Selector

const (
	// ForeignKeyIndex 外键名下标
	ForeignKeyIndex = 0
	// RelationKeyIndex 关联键名下标
	RelationKeyIndex = 1
)

var (
	ErrUnableDetermineTableName = errors.New(`Unable to determine table name`)
	TableName                   = DefaultTableName
)

func (b *sqlBuilder) Relation(name string, fn BuilderChainFunc) SQLBuilder {
	if b.relationMap == nil {
		b.relationMap = make(map[string]BuilderChainFunc)
	}
	b.relationMap[name] = fn
	return b
}

func (b *sqlBuilder) RelationMap() map[string]BuilderChainFunc {
	return b.relationMap
}

func (sel *selector) Relation(name string, fn BuilderChainFunc) Selector {
	sel.SQLBuilder().Relation(name, fn)
	return sel
}

func eachField(t reflect.Type, fn func(field reflect.StructField, relations []string, pipes []Pipe) error) error {
	typeMap := mapper.TypeMap(t)
	options, ok := typeMap.Options[`relation`]
	if !ok {
		return nil
	}
	for _, fieldInfo := range options {
		//fmt.Println(`==>`, fieldInfo.Name, fieldInfo.Embedded, com.Dump(fieldInfo.Options, false))
		// `db:"-,relation=ForeignKey:RelationKey"`
		// `db:"-,relation=外键名:关联键名|gtZero|eq(field:value)"`
		rel, ok := fieldInfo.Options[`relation`]
		if !ok || len(rel) == 0 || rel == `-` {
			continue
		}
		relations := strings.SplitN(rel, `:`, 2)
		if len(relations) != 2 {
			return fmt.Errorf("Wrong relation option, length must 2, but get %v. Reference format: `db:\"-,relation=ForeignKey:RelationKey\"`", relations)
		}
		rels := strings.Split(relations[1], `|`) // `db:"-,relation=外键名:关联键名|neq(field,value)"`
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
		err := fn(fieldInfo.Field, relations, pipes)
		if err != nil {
			return err
		}
	}
	return nil
}

type Name_ interface {
	Name_() string
}

func DefaultTableName(data interface{}, retry ...bool) (string, error) {
	switch m := data.(type) {
	case Name_:
		return m.Name_(), nil
	case db.TableName:
		return m.TableName(), nil
	default:
		if len(retry) > 0 && retry[0] {
			return ``, ErrUnableDetermineTableName
		}
	}
	value := reflect.ValueOf(data)
	if value.IsNil() {
		return ``, errors.WithMessagef(errors.New("model argument cannot be nil pointer passed"), `%T`, data)
	}
	tp := reflect.Indirect(value).Type()
	if tp.Kind() == reflect.Interface {
		tp = reflect.Indirect(value).Elem().Type()
	}

	if tp.Kind() != reflect.Slice {
		return ``, fmt.Errorf("model argument must slice, but get %T", data)
	}

	tpEl := tp.Elem()
	//Compatible with []*Struct or []Struct
	if tpEl.Kind() == reflect.Ptr {
		tpEl = tpEl.Elem()
	}
	//fmt.Printf("[TableName] %s ========>%[1]T, %[1]v\n", tpEl.Name(), reflect.New(tpEl).Interface())
	name, err := DefaultTableName(reflect.New(tpEl).Interface(), true)
	if err == ErrUnableDetermineTableName {
		name = com.SnakeCase(tpEl.Name())
		err = nil
	}
	return name, err
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

// RelationOne is get the associated relational data for a single piece of data
func RelationOne(builder SQLBuilder, data interface{}) error {
	refVal := reflect.Indirect(reflect.ValueOf(data))
	t := refVal.Type()

	return eachField(t, func(field reflect.StructField, relations []string, pipes []Pipe) error {
		name := field.Name
		var foreignModel reflect.Value
		// if field type is slice then one-to-many ,eg: []*Struct
		if field.Type.Kind() == reflect.Slice {
			foreignModel = reflect.New(field.Type)
			foreignIV := foreignModel.Interface()
			table, err := TableName(foreignIV)
			if err != nil {
				return err
			}
			// batch get field values
			// Since the structure is slice, there is no need to new Value
			cond := buildCond(refVal, relations, pipes)
			if cond == nil {
				return nil
			}
			sel := builder.SelectFrom(table).Where(cond)
			if chains := builder.RelationMap(); chains != nil {
				if chainFn, ok := chains[name]; ok {
					sel = chainFn(sel)
				}
			}
			err = sel.All(foreignIV)
			if err != nil && err != db.ErrNoMoreRows {
				return err
			}

			if reflect.Indirect(foreignModel).Len() == 0 {
				// If relation data is empty, must set empty slice
				// Otherwise, the JSON result will be null instead of []
				refVal.FieldByName(name).Set(reflect.MakeSlice(field.Type, 0, 0))
			} else {
				refVal.FieldByName(name).Set(foreignModel.Elem())
			}

		} else {
			// If field type is struct the one-to-one,eg: *Struct
			foreignModel = reflect.New(field.Type.Elem())
			foreignIV := foreignModel.Interface()
			table, err := TableName(foreignIV)
			if err != nil {
				return err
			}
			cond := buildCond(refVal, relations, pipes)
			if cond == nil {
				return nil
			}
			sel := builder.SelectFrom(table).Where(cond)
			if chains := builder.RelationMap(); chains != nil {
				if chainFn, ok := chains[name]; ok {
					sel = chainFn(sel)
				}
			}
			err = sel.One(foreignIV)
			// If one-to-one NoRows is not an error that needs to be terminated
			if err != nil && err != db.ErrNoMoreRows {
				return err
			}

			if err == nil {
				refVal.FieldByName(name).Set(foreignModel)
			}
		}
		return nil
	})
}

// RelationAll is gets the associated relational data for multiple pieces of data
func RelationAll(builder SQLBuilder, data interface{}) error {
	refVal := reflect.Indirect(reflect.ValueOf(data))

	l := refVal.Len()

	if l == 0 {
		return nil
	}

	// get the struct field in slice
	t := reflect.Indirect(refVal.Index(0)).Type()

	return eachField(t, func(field reflect.StructField, relations []string, pipes []Pipe) error {
		name := field.Name
		relVals := make([]interface{}, 0)
		relValsMap := make(map[interface{}]struct{}, 0)
		relValsMapx := make(map[int][]interface{}, 0)
		fieldName := relations[ForeignKeyIndex]
		rFieldName := relations[RelationKeyIndex]
		var rt reflect.Kind
		if l > 0 {
			rt = mapper.FieldByName(refVal.Index(0), rFieldName).Kind()
		}
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

		for k := range relValsMap {
			relVals = append(relVals, k)
		}

		var foreignModel reflect.Value
		// if field type is slice then one to many ,eg: []*Struct
		if field.Type.Kind() == reflect.Slice {
			foreignModel = reflect.New(field.Type)
			foreignIV := foreignModel.Interface()
			table, err := TableName(foreignIV)
			if err != nil {
				return err
			}
			// batch get field values
			// Since the structure is slice, there is no need to new Value
			sel := builder.SelectFrom(table).Where(db.Cond{
				fieldName: db.In(relVals),
			})
			if chains := builder.RelationMap(); chains != nil {
				if chainFn, ok := chains[name]; ok {
					sel = chainFn(sel)
				}
			}
			err = sel.All(foreignIV)
			if err != nil && err != db.ErrNoMoreRows {
				return err
			}

			fmap := make(map[interface{}]reflect.Value)

			// Combine relation data as a one-to-many relation
			// For example, if there are multiple images under an article
			// we use the article ID to associate the images, map[1][]*Images
			mlen := reflect.Indirect(foreignModel).Len()
			var ft reflect.Kind
			if mlen > 0 {
				ft = mapper.FieldByName(reflect.Indirect(foreignModel).Index(0), fieldName).Kind()
			}
			for n := 0; n < mlen; n++ {
				val := reflect.Indirect(foreignModel).Index(n)
				fid := mapper.FieldByName(val, fieldName)
				fv := fid.Interface()
				if _, has := fmap[fv]; !has {
					fmap[fv] = reflect.New(reflect.SliceOf(field.Type.Elem())).Elem()
				}
				fmap[fv] = reflect.Append(fmap[fv], val)
			}
			needConversion := rt != ft && ft != reflect.Invalid
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
							_v = param.AsType(ft.String(), _v)
							if value, has := fmap[_v]; has {
								slicev = reflect.AppendSlice(slicev, value)
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
			// If field type is struct the one to one,eg: *Struct
			foreignModel = reflect.New(field.Type.Elem())

			// Batch get field values, but must new slice []*Struct
			fi := reflect.New(reflect.SliceOf(foreignModel.Type()))

			foreignIV := fi.Interface()
			table, err := TableName(foreignIV)
			if err != nil {
				return err
			}
			sel := builder.SelectFrom(table).Where(db.Cond{
				fieldName: db.In(relVals),
			})
			if chains := builder.RelationMap(); chains != nil {
				if chainFn, ok := chains[name]; ok {
					sel = chainFn(sel)
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
			var ft reflect.Kind
			if mlen > 0 {
				ft = mapper.FieldByName(fval.Index(0), fieldName).Kind()
			}
			for n := 0; n < mlen; n++ {
				val := fval.Index(n)
				fid := mapper.FieldByName(val, fieldName)
				fmap[fid.Interface()] = val
			}
			needConversion := rt != ft && ft != reflect.Invalid
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
