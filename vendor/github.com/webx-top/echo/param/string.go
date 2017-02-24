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
package param

import (
	"strconv"
)

type String string

func (p String) String() string {
	return string(p)
}

func (p String) Int() int {
	if len(p) > 0 {
		r, _ := strconv.Atoi(p.String())
		return r
	}
	return 0
}

func (p String) Int64() int64 {
	if len(p) > 0 {
		r, _ := strconv.ParseInt(p.String(), 10, 64)
		return r
	}
	return 0
}

func (p String) Int32() int32 {
	if len(p) > 0 {
		r, _ := strconv.ParseInt(p.String(), 10, 32)
		return int32(r)
	}
	return 0
}

func (p String) Uint() uint {
	if len(p) > 0 {
		r, _ := strconv.ParseUint(p.String(), 10, 64)
		return uint(r)
	}
	return 0
}

func (p String) Uint64() uint64 {
	if len(p) > 0 {
		r, _ := strconv.ParseUint(p.String(), 10, 64)
		return r
	}
	return 0
}

func (p String) Uint32() uint32 {
	if len(p) > 0 {
		r, _ := strconv.ParseUint(p.String(), 10, 32)
		return uint32(r)
	}
	return 0
}

func (p String) Float32() float32 {
	if len(p) > 0 {
		r, _ := strconv.ParseFloat(p.String(), 32)
		return float32(r)
	}
	return 0
}

func (p String) Float64() float64 {
	if len(p) > 0 {
		r, _ := strconv.ParseFloat(p.String(), 64)
		return r
	}
	return 0
}

func (p String) Bool() bool {
	if len(p) > 0 {
		r, _ := strconv.ParseBool(p.String())
		return r
	}
	return false
}
