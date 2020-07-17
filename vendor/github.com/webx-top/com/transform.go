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
	"strings"
)

func Int64(i interface{}) int64 {
	switch v := i.(type) {
	case int64:
		return v
	case int32:
		return int64(v)
	case uint32:
		return int64(v)
	case int:
		return int64(v)
	case uint:
		return int64(v)
	case float32:
		return int64(v)
	case float64:
		return int64(v)
	case string:
		out, _ := strconv.ParseInt(v, 10, 64)
		return out
	case nil:
		return 0
	default:
		in := Str(i)
		if len(in) == 0 {
			return 0
		}
		out, err := strconv.ParseInt(in, 10, 64)
		if err != nil {
			log.Printf("string[%s] covert int64 fail. %s", in, err)
			return 0
		}
		return out
	}
}

func Int(i interface{}) int {
	switch v := i.(type) {
	case int32:
		return int(v)
	case int:
		return v
	case float32:
		return int(v)
	case float64:
		return int(v)
	case string:
		out, _ := strconv.Atoi(v)
		return out
	case nil:
		return 0
	default:
		in := Str(i)
		if len(in) == 0 {
			return 0
		}
		out, err := strconv.Atoi(in)
		if err != nil {
			log.Printf("string[%s] covert int fail. %s", in, err)
			return 0
		}
		return out
	}
}

func Int32(i interface{}) int32 {
	switch v := i.(type) {
	case int32:
		return v
	case float32:
		return int32(v)
	case float64:
		return int32(v)
	case string:
		out, _ := strconv.ParseInt(v, 10, 32)
		return int32(out)
	case nil:
		return 0
	default:
		in := Str(i)
		if len(in) == 0 {
			return 0
		}
		out, err := strconv.ParseInt(in, 10, 32)
		if err != nil {
			log.Printf("string[%s] covert int32 fail. %s", in, err)
			return 0
		}
		return int32(out)
	}
}

func Uint64(i interface{}) uint64 {
	switch v := i.(type) {
	case uint32:
		return uint64(v)
	case uint:
		return uint64(v)
	case uint64:
		return v
	case float32:
		if v > 0 {
			return uint64(v)
		}
		return 0
	case float64:
		if v > 0 {
			return uint64(v)
		}
		return 0
	case string:
		out, _ := strconv.ParseUint(v, 10, 64)
		return out
	case nil:
		return 0
	default:
		in := Str(i)
		if len(in) == 0 {
			return 0
		}
		out, err := strconv.ParseUint(in, 10, 64)
		if err != nil {
			log.Printf("string[%s] covert uint64 fail. %s", in, err)
			return 0
		}
		return out
	}
}

func Uint(i interface{}) uint {
	switch v := i.(type) {
	case uint32:
		return uint(v)
	case uint:
		return v
	case float32:
		if v > 0 {
			return uint(v)
		}
		return 0
	case float64:
		if v > 0 {
			return uint(v)
		}
		return 0
	case string:
		out, _ := strconv.ParseUint(v, 10, 32)
		return uint(out)
	case nil:
		return 0
	default:
		in := Str(i)
		if len(in) == 0 {
			return 0
		}
		out, err := strconv.ParseUint(in, 10, 32)
		if err != nil {
			log.Printf("string[%s] covert uint fail. %s", in, err)
			return 0
		}
		return uint(out)
	}
}

func Uint32(i interface{}) uint32 {
	switch v := i.(type) {
	case uint32:
		return v
	case uint:
		return uint32(v)
	case uint64:
		return uint32(v)
	case float32:
		if v > 0 {
			return uint32(v)
		}
		return 0
	case float64:
		if v > 0 {
			return uint32(v)
		}
		return 0
	case string:
		out, _ := strconv.ParseUint(v, 10, 32)
		return uint32(out)
	case nil:
		return 0
	default:
		in := Str(i)
		if len(in) == 0 {
			return 0
		}
		out, err := strconv.ParseUint(in, 10, 32)
		if err != nil {
			log.Printf("string[%s] covert uint32 fail. %s", in, err)
			return 0
		}
		return uint32(out)
	}
}

func Float32(i interface{}) float32 {
	switch v := i.(type) {
	case float32:
		return v
	case float64:
		return float32(v)
	case int8:
		return float32(v)
	case uint8:
		return float32(v)
	case int16:
		return float32(v)
	case uint16:
		return float32(v)
	case int32:
		return float32(v)
	case uint32:
		return float32(v)
	case int:
		return float32(v)
	case uint:
		return float32(v)
	case int64:
		return float32(v)
	case uint64:
		return float32(v)
	case string:
		out, _ := strconv.ParseFloat(v, 32)
		return float32(out)
	case nil:
		return 0
	default:
		in := Str(i)
		if len(in) == 0 {
			return 0
		}
		out, err := strconv.ParseFloat(in, 32)
		if err != nil {
			log.Printf("string[%s] covert float32 fail. %s", in, err)
			return 0
		}
		return float32(out)
	}
}

func Float64(i interface{}) float64 {
	switch v := i.(type) {
	case float32:
		return float64(v)
	case float64:
		return v
	case int8:
		return float64(v)
	case uint8:
		return float64(v)
	case int16:
		return float64(v)
	case uint16:
		return float64(v)
	case int32:
		return float64(v)
	case uint32:
		return float64(v)
	case int:
		return float64(v)
	case uint:
		return float64(v)
	case int64:
		return float64(v)
	case uint64:
		return float64(v)
	case string:
		out, _ := strconv.ParseFloat(v, 64)
		return out
	case nil:
		return 0
	default:
		in := Str(i)
		if len(in) == 0 {
			return 0
		}
		out, err := strconv.ParseFloat(in, 64)
		if err != nil {
			log.Printf("string[%s] covert float64 fail. %s", in, err)
			return 0
		}
		return out
	}
}

func Bool(i interface{}) bool {
	switch v := i.(type) {
	case bool:
		return v
	case nil:
		return false
	default:
		in := Str(i)
		if len(in) == 0 {
			return false
		}
		out, err := strconv.ParseBool(in)
		if err != nil {
			log.Printf("string[%s] covert bool fail. %s", in, err)
			return false
		}
		return out
	}
}

func Str(i interface{}) string {
	switch v := i.(type) {
	case string:
		return v
	case nil:
		return ``
	default:
		return fmt.Sprint(v)
	}
}

func String(v interface{}) string {
	return Str(v)
}

// SeekRangeNumbers 遍历范围数值，支持设置步进值。格式例如：1-2,2-3:2
func SeekRangeNumbers(expr string, fn func(int) bool) {
	expa := strings.SplitN(expr, ":", 2)
	step := 1
	switch len(expa) {
	case 2:
		if i, e := strconv.Atoi(strings.TrimSpace(expa[1])); e == nil {
			step = i
		}
		fallthrough
	case 1:
		for _, exp := range strings.Split(strings.TrimSpace(expa[0]), `,`) {
			exp = strings.TrimSpace(exp)
			if len(exp) == 0 {
				continue
			}
			expb := strings.SplitN(exp, `-`, 2)
			var minN, maxN int
			switch len(expb) {
			case 2:
				maxN, _ = strconv.Atoi(strings.TrimSpace(expb[1]))
				fallthrough
			case 1:
				minN, _ = strconv.Atoi(strings.TrimSpace(expb[0]))
			}
			if maxN == 0 {
				if !fn(minN) {
					return
				}
				continue
			}
			for ; minN <= maxN; minN += step {
				if !fn(minN) {
					return
				}
			}
		}
	}
}
