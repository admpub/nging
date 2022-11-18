// Copyright (c) 2012-present The upper.io/db authors. All rights reserved.
//
// Permission is hereby granted, free of charge, to any person obtaining
// a copy of this software and associated documentation files (the
// "Software"), to deal in the Software without restriction, including
// without limitation the rights to use, copy, modify, merge, publish,
// distribute, sublicense, and/or sell copies of the Software, and to
// permit persons to whom the Software is furnished to do so, subject to
// the following conditions:
//
// The above copyright notice and this permission notice shall be
// included in all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
// EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
// MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
// NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE
// LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION
// OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION
// WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

package sqlbuilder

import (
	"database/sql"
	"database/sql/driver"
	"reflect"

	"github.com/webx-top/db"
	"github.com/webx-top/db/lib/reflectx"
	"github.com/webx-top/echo/param"
)

type hasConvertValues interface {
	ConvertValues(values []interface{}) []interface{}
}

var mapper = reflectx.NewMapper("db")

// fetchRow receives a *sql.Rows value and tries to map all the rows into a
// single struct given by the pointer `dst`.
func fetchRow(iter *iterator, dst interface{}) error {
	var columns []string
	var err error

	rows := iter.cursor

	dstv := reflect.ValueOf(dst)

	if dstv.IsNil() || dstv.Kind() != reflect.Ptr {
		return ErrExpectingPointer
	}

	itemV := dstv.Elem()

	if columns, err = rows.Columns(); err != nil {
		return err
	}

	reset(dst)

	next := rows.Next()

	if !next {
		if err = rows.Err(); err != nil {
			return err
		}
		return db.ErrNoMoreRows
	}

	itemT := itemV.Type()
	item, err := fetchResult(iter, itemT, columns)

	if err != nil {
		return err
	}

	if itemT.Kind() == reflect.Ptr {
		itemV.Set(item)
	} else {
		itemV.Set(reflect.Indirect(item))
	}

	return RelationOne(iter.SQLBuilder, dst, iter.relationMap)
}

// fetchRows receives a *sql.Rows value and tries to map all the rows into a
// slice of structs given by the pointer `dst`.
func fetchRows(iter *iterator, dst interface{}) error {
	var err error
	rows := iter.cursor
	defer rows.Close()

	// Destination.
	dstv := reflect.ValueOf(dst)

	if dstv.IsNil() || dstv.Kind() != reflect.Ptr {
		return ErrExpectingPointer
	}

	if dstv.Elem().Kind() != reflect.Slice {
		return ErrExpectingSlicePointer
	}

	if dstv.Kind() != reflect.Ptr || dstv.Elem().Kind() != reflect.Slice || dstv.IsNil() {
		return ErrExpectingSliceMapStruct
	}

	var columns []string
	if columns, err = rows.Columns(); err != nil {
		return err
	}

	slicev := dstv.Elem()
	itemT := slicev.Type().Elem()

	reset(dst)

	for rows.Next() {
		item, err := fetchResult(iter, itemT, columns)
		if err != nil {
			return err
		}
		if itemT.Kind() == reflect.Ptr {
			slicev = reflect.Append(slicev, item)
		} else {
			slicev = reflect.Append(slicev, reflect.Indirect(item))
		}
	}

	dstv.Elem().Set(slicev)
	if err := rows.Err(); err != nil {
		return err
	}
	return RelationAll(iter.SQLBuilder, dst, iter.relationMap)
}

func fetchResult(iter *iterator, itemT reflect.Type, columns []string) (reflect.Value, error) {
	var item reflect.Value
	var err error
	rows := iter.cursor

	objT := itemT

	switch objT.Kind() {
	case reflect.Map:
		item = reflect.MakeMap(objT)
	case reflect.Struct:
		item = reflect.New(objT)
	case reflect.Ptr:
		objT = itemT.Elem()
		switch objT.Kind() {
		case reflect.Struct:
			item = reflect.New(objT)
		case reflect.Map:
			item = reflect.MakeMap(objT)
		default:
			return item, ErrExpectingMapOrStruct
		}
	default:
		return item, ErrExpectingMapOrStruct
	}

	switch objT.Kind() {
	case reflect.Struct:

		values := make([]interface{}, len(columns))
		typeMap := mapper.TypeMap(itemT)
		fieldMap := typeMap.Names

		for i, k := range columns {
			fi, ok := fieldMap[k]
			if !ok {
				values[i] = new(interface{})
				continue
			}

			// Check for deprecated jsonb tag.
			if _, hasJSONBTag := fi.Options["jsonb"]; hasJSONBTag {
				return item, errDeprecatedJSONBTag
			}

			f := reflectx.FieldByIndexes(item, fi.Index)
			values[i] = f.Addr().Interface()

			if u, ok := values[i].(db.Unmarshaler); ok {
				values[i] = scanner{u}
			}
		}

		if converter, ok := iter.SQLBuilder.(*sqlBuilder).sess.(hasConvertValues); ok {
			values = converter.ConvertValues(values)
		}

		if err = rows.Scan(values...); err != nil {
			return item, err
		}
	case reflect.Map:

		columns, err := rows.Columns()
		if err != nil {
			return item, err
		}
		columnTypes, err := rows.ColumnTypes()
		if err != nil {
			return item, err
		}

		values := make([]interface{}, len(columns))
		prepareValues(values, columnTypes, columns) // [SWH|+]

		if err = rows.Scan(values...); err != nil {
			return item, err
		}

		elemType := itemT.Elem()

		for i, column := range columns {
			// [SWH|+]------\
			reflectValue := reflect.Indirect(reflect.Indirect(reflect.ValueOf(values[i])))
			if reflectValue.IsValid() && reflectValue.CanInterface() {
				switch rawValue := reflectValue.Interface().(type) {
				case driver.Valuer:
					val, _ := rawValue.Value()
					reflectValue = reflect.ValueOf(val)
				case sql.RawBytes:
					reflectValue = reflect.ValueOf(string(rawValue))
				}
			}
			var elemValue reflect.Value
			if elemType.Kind() != reflectValue.Kind() && elemType.Kind() != reflect.Interface {
				isPtr := elemType.Kind() == reflect.Ptr
				nonPtrType := elemType
				if isPtr {
					nonPtrType = itemT.Elem().Elem()
				}
				elemValue = reflect.New(nonPtrType)
				if nonPtrType.Kind() == reflect.Struct {
					if sc, ok := elemValue.Interface().(sql.Scanner); ok {
						sc.Scan(reflectValue.Interface())
					}
				} else {
					value := param.AsType(nonPtrType.Kind().String(), reflectValue.Interface())
					elemValue.Set(reflect.ValueOf(value))
				}
				if !isPtr {
					elemValue = elemValue.Elem()
				}
			} else {
				elemValue = reflectValue
			}
			// [SWH|+]------/
			item.SetMapIndex(reflect.ValueOf(column), elemValue)
		}
	}

	return item, nil
}

func reset(data interface{}) error {
	// Resetting element.
	v := reflect.ValueOf(data).Elem()
	t := v.Type()

	var z reflect.Value

	switch v.Kind() {
	case reflect.Slice:
		z = reflect.MakeSlice(t, 0, v.Cap())
	default:
		z = reflect.Zero(t)
	}

	v.Set(z)
	return nil
}

// prepareValues prepare values slice
func prepareValues(values []interface{}, columnTypes []*sql.ColumnType, columns []string) {
	if len(columnTypes) > 0 {
		for idx, columnType := range columnTypes {
			if columnType.ScanType() != nil {
				values[idx] = reflect.New(reflect.PtrTo(columnType.ScanType())).Interface()
			} else {
				values[idx] = new(interface{})
			}
		}
	} else {
		for idx := range columns {
			values[idx] = new(interface{})
		}
	}
}
