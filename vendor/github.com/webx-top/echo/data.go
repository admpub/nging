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
package echo

import (
	"encoding/gob"
	"fmt"
)

func init() {
	gob.Register(&Data{})
}

type Data struct {
	context Context
	Code    int
	Info    interface{}
	Zone    interface{} `json:",omitempty" xml:",omitempty"`
	Data    interface{} `json:",omitempty" xml:",omitempty"`
}

func (d *Data) Error() string {
	return fmt.Sprintf(`%v`, d.Info)
}

func (d *Data) String() string {
	return fmt.Sprintf(`%v`, d.Info)
}

func (d *Data) Render(tmpl string, code ...int) error {
	return d.context.Render(tmpl, d.Data, code...)
}

func (d *Data) Gets() (int, interface{}, interface{}, interface{}) {
	return d.Code, d.Info, d.Zone, d.Data
}

func (d *Data) GetData() interface{} {
	return d.Data
}

func (d *Data) SetError(err error, args ...int) *Data {
	if err != nil {
		if len(args) > 0 {
			d.Code = args[0]
		} else {
			d.Code = 0
		}
		d.Info = err.Error()
	} else {
		d.Code = 1
	}
	return d
}

func (d *Data) SetCode(code int) *Data {
	d.Code = code
	return d
}

func (d *Data) SetInfo(info interface{}, args ...int) *Data {
	d.Info = info
	if len(args) > 0 {
		d.Code = args[0]
	}
	return d
}

func (d *Data) SetZone(zone interface{}) *Data {
	d.Zone = zone
	return d
}

func (d *Data) SetData(data interface{}, args ...int) *Data {
	d.Data = data
	if len(args) > 0 {
		d.Code = args[0]
	} else {
		d.Code = 1
	}
	return d
}

func (d *Data) SetContext(ctx Context) *Data {
	d.context = ctx
	return d
}

func (c *Data) Assign(key string, val interface{}) {
	data, _ := c.Data.(H)
	if data == nil {
		data = H{}
	}
	data[key] = val
	c.Data = data
}

func (c *Data) Assignx(values *map[string]interface{}) {
	if values == nil {
		return
	}
	data, _ := c.Data.(H)
	if data == nil {
		data = H{}
	}
	for key, val := range *values {
		data[key] = val
	}
	c.Data = data
}

func (c *Data) SetTmplFuncs() {
	flash, ok := c.context.Session().Get(`webx:flash`).(*Data)
	if ok {
		c.context.Session().Delete(`webx:flash`).Save()
		c.context.SetFunc(`Code`, func() int {
			return flash.Code
		})
		c.context.SetFunc(`Info`, func() interface{} {
			return flash.Info
		})
		c.context.SetFunc(`Zone`, func() interface{} {
			return flash.Zone
		})
	} else {
		c.context.SetFunc(`Code`, func() int {
			return c.Code
		})
		c.context.SetFunc(`Info`, func() interface{} {
			return c.Info
		})
		c.context.SetFunc(`Zone`, func() interface{} {
			return c.Zone
		})
	}
}

// Set 设置输出(code,info,zone,data)
func (c *Data) Set(code int, args ...interface{}) {
	c.Code = code
	var hasData bool
	switch len(args) {
	case 3:
		c.Data = args[2]
		hasData = true
		fallthrough
	case 2:
		c.Zone = args[1]
		fallthrough
	case 1:
		c.Info = args[0]
		if !hasData {
			flash := &Data{
				context: c.context,
				Code:    c.Code,
				Info:    c.Info,
				Zone:    c.Zone,
				Data:    nil,
			}
			c.context.Session().Set(`webx:flash`, flash).Save()
		}
	}
}

// NewData params: Code,Info,Zone,Data
func NewData(ctx Context, code int, args ...interface{}) *Data {
	var info, zone, data interface{}
	switch len(args) {
	case 3:
		data = args[2]
		fallthrough
	case 2:
		zone = args[1]
		fallthrough
	case 1:
		info = args[0]
	}
	return &Data{
		context: ctx,
		Code:    code,
		Info:    info,
		Zone:    zone,
		Data:    data,
	}
}

type KV struct {
	K string
	V string
}

func NewKVData() *KVData {
	return &KVData{
		slice: []*KV{},
		index: map[string][]int{},
	}
}

type KVData struct {
	slice []*KV
	index map[string][]int
}

func (a *KVData) Slice() []*KV {
	return a.slice
}

func (a *KVData) Index(k string) []int {
	v, _ := a.index[k]
	return v
}

func (a *KVData) Indexes() map[string][]int {
	return a.index
}

func (a *KVData) Reset() *KVData {
	a.index = map[string][]int{}
	a.slice = []*KV{}
	return a
}

func (a *KVData) Add(k, v string) *KVData {
	if _, y := a.index[k]; !y {
		a.index[k] = []int{}
	}
	a.index[k] = append(a.index[k], len(a.slice))
	a.slice = append(a.slice, &KV{K: k, V: v})
	return a
}

func (a *KVData) Set(k, v string) *KVData {
	a.index[k] = []int{0}
	a.slice = []*KV{&KV{K: k, V: v}}
	return a
}

func (a *KVData) Delete(ks ...string) *KVData {
	indexes := []int{}
	for _, k := range ks {
		v, y := a.index[k]
		if !y {
			continue
		}
		for _, key := range v {
			indexes = append(indexes, key)
		}
	}
	newSlice := []*KV{}
	a.index = map[string][]int{}
	for i, v := range a.slice {
		var exists bool
		for _, idx := range indexes {
			if i != idx {
				continue
			}
			exists = true
			break
		}
		if exists {
			continue
		}
		if _, y := a.index[v.K]; !y {
			a.index[v.K] = []int{}
		}
		a.index[v.K] = append(a.index[v.K], len(newSlice))
		newSlice = append(newSlice, v)
	}
	a.slice = newSlice
	return a
}
