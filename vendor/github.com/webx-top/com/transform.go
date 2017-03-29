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
	"fmt"
	"log"
	"strconv"
)

func Int64(i interface{}) int64 {
	in := Str(i)
	if in == "" {
		return 0
	}
	out, err := strconv.ParseInt(in, 10, 64)
	if err != nil {
		log.Printf("string[%s] covert int64 fail. %s", in, err)
		return 0
	}
	return out
}

func Int(i interface{}) int {
	in := Str(i)
	if in == "" {
		return 0
	}
	out, err := strconv.Atoi(in)
	if err != nil {
		log.Printf("string[%s] covert int fail. %s", in, err)
		return 0
	}
	return out
}

func Int32(i interface{}) int32 {
	in := Str(i)
	if in == "" {
		return 0
	}
	out, err := strconv.ParseInt(in, 10, 32)
	if err != nil {
		log.Printf("string[%s] covert int32 fail. %s", in, err)
		return 0
	}
	return int32(out)
}

func Uint64(i interface{}) uint64 {
	in := Str(i)
	if in == "" {
		return 0
	}
	out, err := strconv.ParseUint(in, 10, 64)
	if err != nil {
		log.Printf("string[%s] covert uint64 fail. %s", in, err)
		return 0
	}
	return out
}

func Uint(i interface{}) uint {
	in := Str(i)
	if in == "" {
		return 0
	}
	out, err := strconv.ParseUint(in, 10, 32)
	if err != nil {
		log.Printf("string[%s] covert uint fail. %s", in, err)
		return 0
	}
	return uint(out)
}

func Uint32(i interface{}) uint32 {
	in := Str(i)
	if in == "" {
		return 0
	}
	out, err := strconv.ParseUint(in, 10, 32)
	if err != nil {
		log.Printf("string[%s] covert uint32 fail. %s", in, err)
		return 0
	}
	return uint32(out)
}

func Float32(i interface{}) float32 {
	in := Str(i)
	if in == "" {
		return 0
	}
	out, err := strconv.ParseFloat(in, 32)
	if err != nil {
		log.Printf("string[%s] covert float32 fail. %s", in, err)
		return 0
	}
	return float32(out)
}

func Float64(i interface{}) float64 {
	in := Str(i)
	if in == "" {
		return 0
	}
	out, err := strconv.ParseFloat(in, 64)
	if err != nil {
		log.Printf("string[%s] covert float64 fail. %s", in, err)
		return 0
	}
	return out
}

func Bool(i interface{}) bool {
	in := Str(i)
	if in == "" {
		return false
	}
	out, err := strconv.ParseBool(in)
	if err != nil {
		log.Printf("string[%s] covert bool fail. %s", in, err)
		return false
	}
	return out
}

func Str(v interface{}) string {
	return fmt.Sprintf("%v", v)
}

func String(v interface{}) string {
	return Str(v)
}
