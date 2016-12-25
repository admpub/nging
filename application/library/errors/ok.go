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
package errors

import (
	"encoding/gob"
)

func init() {
	gob.Register(&Success{})
}

func NewOk(v string) Successor {
	return &Success{
		Value: v,
	}
}

type Success struct {
	Value string
}

func (s *Success) Success() string {
	return s.Value
}

func (s *Success) String() string {
	return s.Value
}

type Successor interface {
	Success() string
}

func IsOk(err interface{}) bool {
	_, y := err.(Successor)
	return y
}

func Ok(err interface{}) string {
	if v, y := err.(Successor); y {
		return v.Success()
	}
	return ``
}
