package sqlbuilder

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/admpub/errors"
	"github.com/webx-top/com"
	"github.com/webx-top/db"
)

type BuilderChainFunc func(Selector) Selector

const (
	// ForeignKeyIndex 外键名下标
	ForeignKeyIndex = 0
	// RelationKeyIndex 关联键名下标
	RelationKeyIndex = 1
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
		// `db:"-,relation=外键名:关联键名"`
		rel, ok := fieldInfo.Options[`relation`]
		if !ok || len(rel) == 0 || rel == `-` {
			continue
		}
		relations := strings.SplitN(rel, `:`, 2)
		if len(relations) != 2 {
			return fmt.Errorf("Wrong relation option, length must 2, but get %v. Reference format: `db:\"-,relation=ForeignKey:RelationKey\"`", relations)
		}
		rels := strings.Split(relations[1], `|`)
		var pipes []Pipe
		if len(rels) > 1 {
			relations[1] = rels[0]
			for _, pipeName := range rels[1:] {
				pipe, ok := PipeList[pipeName]
				if !ok {
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

type Pipe func(interface{}) interface{}
type Pipes map[string]Pipe

func (pipes *Pipes) Add(name string, pipe Pipe) {
	(*pipes)[name] = pipe
}

var (
	ErrUnableDetermineTableName = errors.New(`Unable to determine table name`)
	TableName                   = DefaultTableName
	PipeList                    = Pipes{
		`split`: func(v interface{}) interface{} {
			items := strings.Split(v.(string), `,`)
			result := []interface{}{}
			for _, item := range items {
				item = strings.TrimSpace(item)
				if len(item) == 0 {
					continue
				}
				result = append(result, item)
			}
			return result
		},
	}
)

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
		fieldValue = pipe(fieldValue)
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
		relValsMap := make(map[interface{}]interface{}, 0)
		fieldName := relations[ForeignKeyIndex]
		rFieldName := relations[RelationKeyIndex]
		// get relation field values and unique
		if len(pipes) == 0 {
			for j := 0; j < l; j++ {
				v := mapper.FieldByName(refVal.Index(j), rFieldName).Interface()
				relValsMap[v] = nil
			}
		} else {
			for j := 0; j < l; j++ {
				v := mapper.FieldByName(refVal.Index(j), rFieldName).Interface()
				for _, pipe := range pipes {
					v = pipe(v)
				}
				if vs, ok := v.([]interface{}); ok {
					for _, vv := range vs {
						relValsMap[vv] = nil
					}
					continue
				}
				relValsMap[v] = nil
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
			for n := 0; n < mlen; n++ {
				val := reflect.Indirect(foreignModel).Index(n)
				fid := mapper.FieldByName(val, fieldName)
				fv := fid.Interface()
				if _, has := fmap[fv]; !has {
					fmap[fv] = reflect.New(reflect.SliceOf(field.Type.Elem())).Elem()
				}
				fmap[fv] = reflect.Append(fmap[fv], val)
			}

			// Set the result to the model
			for j := 0; j < l; j++ {
				v := refVal.Index(j)
				fid := mapper.FieldByName(v, rFieldName)
				if value, has := fmap[fid.Interface()]; has {
					reflect.Indirect(v).FieldByName(name).Set(value)
				} else {
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
			for n := 0; n < mlen; n++ {
				val := fval.Index(n)
				fid := mapper.FieldByName(val, fieldName)
				fmap[fid.Interface()] = val
			}

			// Set the result to the model
			for j := 0; j < l; j++ {
				v := refVal.Index(j)
				fid := mapper.FieldByName(v, rFieldName)
				if value, has := fmap[fid.Interface()]; has {
					reflect.Indirect(v).FieldByName(name).Set(value)
				}
			}
		}

		return nil
	})
}
