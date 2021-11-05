package factory

import (
	"reflect"

	"github.com/webx-top/db"
	"github.com/webx-top/db/lib/sqlbuilder"
)

func (p *Param) UsingStructField(bean interface{}, fields ...string) error {
	if bean == nil {
		bean = p.save
	}
	itemV := reflect.ValueOf(bean)
	if !itemV.IsValid() {
		return ErrExpectingStruct
	}
	itemT := itemV.Type()
	if itemT.Kind() == reflect.Ptr {
		item := itemV.Elem().Interface()
		itemV = reflect.ValueOf(item)
		itemT = itemV.Type()
	}

	switch itemT.Kind() {
	case reflect.Struct:
		mapFields := sqlbuilder.Mapper().TypeMap(itemT).Index
		data := db.NewKeysValues()
		for _, field := range fields {
			for _, mapField := range mapFields {
				if mapField.Field.Name != field { //结构体中字段名称
					continue
				}
				value := itemV.FieldByIndex(mapField.Index)
				if value.CanAddr() {
					value = value.Addr()
				}
				if !value.CanInterface() {
					break
				}
				var key string //数据库中字段名称
				if len(mapField.Name) > 0 {
					key = mapField.Name //结构体tag中db属性所指定的名称
				} else {
					key = mapField.Field.Name //结构体字段名称
				}
				data.Add(key, value.Interface()) //数据库中字段的值
				break
			}
		}
		p.SetSend(data)
	default:
		return ErrExpectingStruct
	}
	return nil
}
