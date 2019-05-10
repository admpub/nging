package com

import (
	"fmt"
	"reflect"
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
		if vv.CanInterface() {
			fmt.Printf("Type: %v => %v;  Field name: %v\n", vv.Kind(), vv.Interface(), vt.Name)
		} else {
			fmt.Printf("Type: %v => %v;  Field name: %v\n", vv.Kind(), "<unexported>", vt.Name)
		}
	}
	fmt.Println("==================[/" + name + "]==================")
}
