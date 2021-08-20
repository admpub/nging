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

import "time"

func ToStringMap(m map[string]string) StringMap {
	r := StringMap{}
	for k, v := range m {
		r[k] = String(v)
	}
	return r
}

type StringMap map[string]String

func (p StringMap) String(key string) string {
	return p[key].String()
}

func (p StringMap) Raw(key string) string {
	return p[key].Raw()
}

func (p StringMap) Split(key string, sep string, limit ...int) StringSlice {
	return p[key].Split(sep, limit...)
}

func (p StringMap) SplitAny(key string, sep string, limit ...int) StringSlice {
	return p[key].SplitAny(sep, limit...)
}

func (p StringMap) Interface(key string) interface{} {
	return p[key].Interface()
}

func (p StringMap) Interfaces() map[string]interface{} {
	r := map[string]interface{}{}
	for k, v := range p {
		r[k] = interface{}(v)
	}
	return r
}

func (p StringMap) Int(key string) int {
	return p[key].Int()
}

func (p StringMap) Int64(key string) int64 {
	return p[key].Int64()
}

func (p StringMap) Int32(key string) int32 {
	return p[key].Int32()
}

func (p StringMap) Uint(key string) uint {
	return p[key].Uint()
}

func (p StringMap) Uint64(key string) uint64 {
	return p[key].Uint64()
}

func (p StringMap) Uint32(key string) uint32 {
	return p[key].Uint32()
}

func (p StringMap) Float32(key string) float32 {
	return p[key].Float32()
}

func (p StringMap) Float64(key string) float64 {
	return p[key].Float64()
}

func (p StringMap) Bool(key string) bool {
	return p[key].Bool()
}

func (p StringMap) Duration(key string, defaults ...time.Duration) time.Duration {
	return p[key].Duration(defaults...)
}

func (p StringMap) Timestamp(key string) time.Time {
	return p[key].Timestamp()
}

func (p StringMap) DateTime(key string) time.Time {
	return p[key].DateTime()
}
