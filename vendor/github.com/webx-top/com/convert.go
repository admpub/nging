// Copyright 2014 com authors
//
// Licensed under the Apache License, Version 2.0 (the "License"): you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations
// under the License.

package com

import (
	"fmt"
	"strconv"
)

// ToStr Convert any type to string.
func ToStr(value interface{}, args ...int) (s string) {
	switch v := value.(type) {
	case bool:
		s = strconv.FormatBool(v)
	case float32:
		a := IntArgs(args)
		s = strconv.FormatFloat(float64(v), 'f', a.Get(0, -1), a.Get(1, 32))
	case float64:
		a := IntArgs(args)
		s = strconv.FormatFloat(v, 'f', a.Get(0, -1), a.Get(1, 64))
	case int:
		s = strconv.FormatInt(int64(v), IntArgs(args).Get(0, 10))
	case int8:
		s = strconv.FormatInt(int64(v), IntArgs(args).Get(0, 10))
	case int16:
		s = strconv.FormatInt(int64(v), IntArgs(args).Get(0, 10))
	case int32:
		s = strconv.FormatInt(int64(v), IntArgs(args).Get(0, 10))
	case int64:
		s = strconv.FormatInt(v, IntArgs(args).Get(0, 10))
	case uint:
		s = strconv.FormatUint(uint64(v), IntArgs(args).Get(0, 10))
	case uint8:
		s = strconv.FormatUint(uint64(v), IntArgs(args).Get(0, 10))
	case uint16:
		s = strconv.FormatUint(uint64(v), IntArgs(args).Get(0, 10))
	case uint32:
		s = strconv.FormatUint(uint64(v), IntArgs(args).Get(0, 10))
	case uint64:
		s = strconv.FormatUint(v, IntArgs(args).Get(0, 10))
	case string:
		s = v
	case []byte:
		s = string(v)
	case []rune:
		s = string(v)
	case nil:
		return
	default:
		s = fmt.Sprintf("%v", v)
	}
	return
}

type IntArgs []int

func (a IntArgs) Get(index int, defaults ...int) (r int) {
	if index >= 0 && index < len(a) {
		r = a[index]
		return
	}
	if len(defaults) > 0 {
		r = defaults[0]
	}
	return
}
