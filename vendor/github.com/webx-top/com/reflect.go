package com

import (
	"fmt"
	"reflect"
	"strings"
)

//TestReflect 测试反射：显示字段和方法信息
func TestReflect(v interface{}) {
	val := reflect.ValueOf(v)
	typ := val.Type()
	if typ.Kind() == reflect.Ptr && typ.Elem().Kind() == reflect.Struct {
		typ = typ.Elem()
		val = val.Elem()
	}
	t := reflect.Indirect(val).Type()
	name := t.Name()
	fmt.Println("==================[" + name + ":]==================")
	fmt.Println("PkgPath:", typ.PkgPath())
	fmt.Printf("Methods total: %v\n", typ.NumMethod())
	for i := 0; i < typ.NumMethod(); i++ {
		vt := typ.Method(i)
		vv := val.Method(i)
		fmt.Printf("Type: %v => %v;  Method name: %v\n", vv.Kind(), vv.Kind() == reflect.Func, vt.Name)
	}
	fmt.Println("-----------------------------------------------")
	fmt.Printf("Fields total: %v\n", val.NumField())
	for i := 0; i < val.NumField(); i++ {
		vt := typ.Field(i)
		vv := val.Field(i)
		var extInfo []string
		if vt.Anonymous {
			extInfo = append(extInfo, `anonymous`)
		}
		if vv.CanInterface() {
			fmt.Printf("Type: %v => %v (%v); Index: %v;  Field name: %v\n", vv.Kind(), vv.Interface(), strings.Join(extInfo, `,`), i, vt.Name)
		} else {
			fmt.Printf("Type: %v => %v (%v); Index: %v;  Field name: %v\n", vv.Kind(), "<unexported>", strings.Join(extInfo, `,`), i, vt.Name)
		}
	}
	fmt.Println("==================[/" + name + "]==================")
}
