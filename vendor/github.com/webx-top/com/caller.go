/*

   Copyright 2016 Wenhui Shen <www.webx.top>

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.

*/
package com

import (
	"errors"
	"path"
	"reflect"
	"runtime"
	"strings"
)

// FuncName 函数名
func FuncName(i interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
}

// FuncPath 返回包名、实例名和函数名
func FuncPath(i interface{}) (pkgName string, objName string, funcName string) {
	s := FuncName(i)
	_, file := path.Split(s)
	return ParseFuncName(file)
}

// FuncFullPath 返回完整路径包名、实例名和函数名
func FuncFullPath(i interface{}) (pkgName string, objName string, funcName string) {
	return ParseFuncName(FuncName(i))
}

// ParseFuncName 解析函数名
func ParseFuncName(funcString string) (pkgName string, objName string, funcName string) {
	if strings.HasSuffix(funcString, `-fm`) {
		funcString = strings.TrimSuffix(funcString, `-fm`)
		part := strings.Split(funcString, `.`)
		switch len(part) {
		case 3:
			funcName = part[2]
			fallthrough
		case 2:
			objName = part[1]
			if objName[0] == '(' {
				objName = objName[1 : len(objName)-1]
				objName = strings.TrimPrefix(objName, `*`)
			}
			fallthrough
		case 1:
			pkgName = part[0]
		}
		return
	}
	part := strings.Split(funcString, `.`)
	switch len(part) {
	case 2:
		funcName = part[1]
		fallthrough
	case 1:
		pkgName = part[0]
	}
	return
}

//CallFunc 根据名称调用对应函数
func CallFunc(getFuncByName func(string) (interface{}, bool), funcName string, params ...interface{}) ([]reflect.Value, error) {
	fn, ok := getFuncByName(funcName)
	if !ok {
		return nil, errors.New("no func found: " + funcName)
	}
	f := reflect.ValueOf(fn)
	if len(params) != f.Type().NumIn() {
		return nil, errors.New("parameter mismatched")
	}
	in := make([]reflect.Value, len(params))
	for k, param := range params {
		in[k] = reflect.ValueOf(param)
	}
	return f.Call(in), nil
}
