/*
 * Copyright 1999-2020 Alibaba Group Holding Ltd.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package util

import (
	"encoding/json"
	"reflect"
	"strconv"
	"strings"

	"github.com/nacos-group/nacos-sdk-go/common/logger"
)

func TransformObject2Param(object interface{}) (params map[string]string) {
	params = make(map[string]string)
	if object != nil {
		valueOf := reflect.ValueOf(object)
		typeOf := reflect.TypeOf(object)
		if reflect.TypeOf(object).Kind() == reflect.Ptr {
			valueOf = reflect.ValueOf(object).Elem()
			typeOf = reflect.TypeOf(object).Elem()
		}
		numField := valueOf.NumField()
		for i := 0; i < numField; i++ {
			tag := typeOf.Field(i).Tag.Get("param")
			if len(tag) > 0 && tag != "-" {
				switch valueOf.Field(i).Kind() {
				case reflect.Int, reflect.Int8, reflect.Int16,
					reflect.Int32, reflect.Int64:
					params[tag] = strconv.FormatInt(valueOf.Field(i).Int(), 10)
				case reflect.Uint, reflect.Uint8, reflect.Uint16,
					reflect.Uint32, reflect.Uint64:
					params[tag] = strconv.FormatUint(valueOf.Field(i).Uint(), 10)
				case reflect.Float32, reflect.Float64:
					params[tag] = strconv.FormatFloat(valueOf.Field(i).Float(), 'f', -1, 64)
				case reflect.Bool:
					params[tag] = strconv.FormatBool(valueOf.Field(i).Bool())
				case reflect.String:
					if len(valueOf.Field(i).String()) > 0 {
						params[tag] = valueOf.Field(i).String()
					}
				case reflect.Map:
					if !valueOf.Field(i).IsNil() {
						bytes, err := json.Marshal(valueOf.Field(i).Interface())
						if err != nil {
							logger.Errorf("[TransformObject2Param] json.Marshal err:%+v", err)
						} else {
							params[tag] = string(bytes)
						}
					}
				case reflect.Slice:
					if ss, ok := valueOf.Field(i).Interface().([]string); ok {
						var pv string
						for _, sv := range ss {
							pv += sv + ","
						}
						if strings.HasSuffix(pv, ",") {
							pv = pv[:len(pv)-1]
						}
						if len(pv) > 0 {
							params[tag] = pv
						}
					}
				}
			}
		}
	}
	return
}
