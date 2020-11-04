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
	"context"
	"encoding/gob"
	"fmt"
	"strconv"

	pkgCode "github.com/webx-top/echo/code"
)

func init() {
	gob.Register(&RawData{})
	gob.Register(H{})
}

func AsError(code pkgCode.Code) *Error {
	return NewError(pkgCode.CodeDict.Get(code).Text, code)
}

//Data 响应数据
type Data interface {
	Assign(key string, val interface{})
	Assignx(values *map[string]interface{})
	SetTmplFuncs()
	SetContext(ctx Context) Data
	String() string
	Set(code int, args ...interface{}) Data
	Reset() Data
	SetByMap(Store) Data
	SetError(err error, args ...int) Data
	SetCode(code int) Data
	SetURL(url string, args ...int) Data
	SetInfo(info interface{}, args ...int) Data
	SetZone(zone interface{}) Data
	SetData(data interface{}, args ...int) Data
	Gets() (code pkgCode.Code, info interface{}, zone interface{}, data interface{})
	GetCode() pkgCode.Code
	GetInfo() interface{}
	GetZone() interface{}
	GetData() interface{}
	GetURL() string
}

type RawData struct {
	context Context
	Code    pkgCode.Code
	State   string `json:",omitempty" xml:",omitempty"`
	Info    interface{}
	URL     string      `json:",omitempty" xml:",omitempty"`
	Zone    interface{} `json:",omitempty" xml:",omitempty"`
	Data    interface{} `json:",omitempty" xml:",omitempty"`
}

func (d *RawData) Error() string {
	return fmt.Sprintf(`%v`, d.Info)
}

func (d *RawData) Reset() Data {
	d.Code = pkgCode.Code(0)
	d.State = ``
	d.Info = nil
	d.URL = ``
	d.Zone = nil
	d.Data = nil
	return d
}

func (d *RawData) copyFrom(v *RawData) Data {
	d.SetCode(v.Code.Int())
	d.Info = v.Info
	if len(v.URL) > 0 {
		d.URL = v.URL
	}
	if v.Zone != nil {
		d.Zone = v.Zone
	}
	if v.Data != nil {
		d.Data = v.Data
	}
	return d
}

func (d *RawData) String() string {
	return fmt.Sprintf(`%v`, d.Info)
}

//Gets 获取全部数据
func (d *RawData) Gets() (pkgCode.Code, interface{}, interface{}, interface{}) {
	return d.Code, d.Info, d.Zone, d.Data
}

func (d *RawData) GetCode() pkgCode.Code {
	return d.Code
}

func (d *RawData) GetInfo() interface{} {
	return d.Info
}

func (d *RawData) GetZone() interface{} {
	return d.Zone
}

func (d *RawData) GetURL() string {
	return d.URL
}

//GetData 获取数据
func (d *RawData) GetData() interface{} {
	return d.Data
}

//SetError 设置错误
func (d *RawData) SetError(err error, args ...int) Data {
	if err == nil {
		return d.SetCode(pkgCode.Success.Int())
	}
	switch v := err.(type) {
	case *Error:
		d.SetInfo(v.Message, v.Code.Int()).SetByMap(v.Extra).SetZone(v.Zone)
	case *RawData:
		if v != d {
			d.copyFrom(v)
		}
	default:
		d.SetCode(pkgCode.Failure.Int())
		d.Info = err.Error()
	}
	if len(args) > 0 {
		d.SetCode(args[0])
	}
	return d
}

//SetCode 设置状态码
func (d *RawData) SetCode(code int) Data {
	d.Code = pkgCode.Code(code)
	d.State = d.Code.String()
	return d
}

//SetURL 设置跳转网址
func (d *RawData) SetURL(url string, args ...int) Data {
	d.URL = url
	if len(args) > 0 {
		d.SetCode(args[0])
	}
	return d
}

//SetInfo 设置提示信息
func (d *RawData) SetInfo(info interface{}, args ...int) Data {
	d.Info = info
	if len(args) > 0 {
		d.SetCode(args[0])
	}
	return d
}

//SetByMap 批量设置属性
func (d *RawData) SetByMap(s Store) Data {
	if len(s) == 0 {
		return d
	}
	if v, y := s["Data"]; y {
		d.Data = v
	}
	if v, y := s["Zone"]; y {
		d.Zone = v
	}
	if v, y := s["Info"]; y {
		d.Info = v
	}
	if v, y := s["URL"]; y {
		d.URL, _ = v.(string)
	}
	var code pkgCode.Code
	if v, y := s["Code"]; y {
		switch c := v.(type) {
		case pkgCode.Code:
			code = c
		case int:
			code = pkgCode.Code(c)
		case string:
			i, _ := strconv.Atoi(c)
			code = pkgCode.Code(i)
		default:
			s := fmt.Sprint(c)
			i, _ := strconv.Atoi(s)
			code = pkgCode.Code(i)
		}
	}
	d.Code = code
	return d
}

//SetZone 设置提示区域
func (d *RawData) SetZone(zone interface{}) Data {
	d.Zone = zone
	return d
}

//SetData 设置正常数据
func (d *RawData) SetData(data interface{}, args ...int) Data {
	d.Data = data
	if len(args) > 0 {
		d.SetCode(args[0])
	} else {
		d.SetCode(1)
	}
	return d
}

//SetContext 设置Context
func (d *RawData) SetContext(ctx Context) Data {
	d.context = ctx
	return d
}

//Assign 赋值
func (d *RawData) Assign(key string, val interface{}) {
	data, _ := d.Data.(H)
	if data == nil {
		data = H{}
	}
	data[key] = val
	d.Data = data
}

//Assignx 批量赋值
func (d *RawData) Assignx(values *map[string]interface{}) {
	if values == nil {
		return
	}
	data, _ := d.Data.(H)
	if data == nil {
		data = H{}
	}
	for key, val := range *values {
		data[key] = val
	}
	d.Data = data
}

//SetTmplFuncs 设置模板函数
func (d *RawData) SetTmplFuncs() {
	flash, ok := d.context.Flash().(*RawData)
	if ok {
		d.context.Session().Save()
	} else {
		flash = d
	}
	d.context.SetFunc(`Code`, func() pkgCode.Code {
		return flash.Code
	})
	d.context.SetFunc(`Info`, func() interface{} {
		return flash.Info
	})
	d.context.SetFunc(`Zone`, func() interface{} {
		return flash.Zone
	})
	d.context.SetFunc(`FURL`, func() interface{} {
		return flash.URL
	})
}

// Set 设置输出(code,info,zone,RawData)
func (d *RawData) Set(code int, args ...interface{}) Data {
	d.SetCode(code)
	var hasData bool
	switch len(args) {
	case 3:
		d.Data = args[2]
		hasData = true
		fallthrough
	case 2:
		d.Zone = args[1]
		fallthrough
	case 1:
		d.Info = args[0]
		if !hasData {
			flash := &RawData{
				context: d.context,
				Code:    d.Code,
				State:   d.State,
				Info:    d.Info,
				URL:     d.URL,
				Zone:    d.Zone,
				Data:    nil,
			}
			d.context.Session().AddFlash(flash).Save()
		}
	}
	return d
}

func NewData(ctx Context) *RawData {
	c := pkgCode.Success
	return &RawData{
		context: ctx,
		Code:    c,
		State:   c.String(),
	}
}

func NewKV(k, v string) *KV {
	return &KV{K: k, V: v}
}

//KV 键值对
type KV struct {
	K  string
	V  string
	H  H           `json:",omitempty" xml:",omitempty"`
	X  interface{} `json:",omitempty" xml:",omitempty"`
	fn func(context.Context) interface{}
}

func (a *KV) SetK(k string) *KV {
	a.K = k
	return a
}

func (a *KV) SetV(v string) *KV {
	a.V = v
	return a
}

func (a *KV) SetKV(k, v string) *KV {
	a.K = k
	a.V = v
	return a
}

func (a *KV) SetH(h H) *KV {
	a.H = h
	return a
}

func (a *KV) SetHKV(k string, v interface{}) *KV {
	if a.H == nil {
		a.H = H{}
	}
	a.H.Set(k, v)
	return a
}

func (a *KV) SetX(x interface{}) *KV {
	a.X = x
	return a
}

func (a *KV) SetFn(fn func(context.Context) interface{}) *KV {
	a.fn = fn
	return a
}

func (a *KV) Fn() func(context.Context) interface{} {
	return a.fn
}

type KVList []*KV

func (list *KVList) Add(k, v string) {
	*list = append(*list, &KV{K: k, V: v})
}

func (list *KVList) AddItem(item *KV) {
	*list = append(*list, item)
}

func (list *KVList) Delete(i int) {
	n := len(*list)
	if i+1 < n {
		*list = append((*list)[0:i], (*list)[i+1:]...)
	} else if i < n {
		*list = (*list)[0:i]
	}
}

func (list *KVList) Reset() {
	*list = (*list)[0:0]
}

//NewKVData 键值对数据
func NewKVData() *KVData {
	return &KVData{
		slice: []*KV{},
		index: map[string][]int{},
	}
}

//KVData 键值对数据（保持顺序）
type KVData struct {
	slice []*KV
	index map[string][]int
}

//Slice 返回切片
func (a *KVData) Slice() []*KV {
	return a.slice
}

//Keys 返回所有K值
func (a *KVData) Keys() []string {
	keys := make([]string, len(a.slice))
	for i, v := range a.slice {
		if v == nil {
			continue
		}
		keys[i] = v.K
	}
	return keys
}

//Index 返回某个key的所有索引值
func (a *KVData) Index(k string) []int {
	v, _ := a.index[k]
	return v
}

//Indexes 返回所有索引值
func (a *KVData) Indexes() map[string][]int {
	return a.index
}

//Reset 重置
func (a *KVData) Reset() *KVData {
	a.index = map[string][]int{}
	a.slice = []*KV{}
	return a
}

//Add 添加键值
func (a *KVData) Add(k, v string) *KVData {
	if _, y := a.index[k]; !y {
		a.index[k] = []int{}
	}
	a.index[k] = append(a.index[k], len(a.slice))
	a.slice = append(a.slice, &KV{K: k, V: v})
	return a
}

func (a *KVData) AddItem(item *KV) *KVData {
	if _, y := a.index[item.K]; !y {
		a.index[item.K] = []int{}
	}
	a.index[item.K] = append(a.index[item.K], len(a.slice))
	a.slice = append(a.slice, item)
	return a
}

//Set 设置首个键值
func (a *KVData) Set(k, v string) *KVData {
	a.index[k] = []int{0}
	a.slice = []*KV{{K: k, V: v}}
	return a
}

func (a *KVData) SetItem(item *KV) *KVData {
	a.index[item.K] = []int{0}
	a.slice = []*KV{item}
	return a
}

func (a *KVData) Get(k string, defaults ...string) string {
	if indexes, ok := a.index[k]; ok {
		if len(indexes) > 0 {
			return a.slice[indexes[0]].V
		}
	}
	if len(defaults) > 0 {
		return defaults[0]
	}
	return ``
}

func (a *KVData) GetItem(k string, defaults ...func() *KV) *KV {
	if indexes, ok := a.index[k]; ok {
		if len(indexes) > 0 {
			return a.slice[indexes[0]]
		}
	}
	if len(defaults) > 0 {
		return defaults[0]()
	}
	return nil
}

func (a *KVData) Has(k string) bool {
	if _, ok := a.index[k]; ok {
		return true
	}
	return false
}

//Delete 设置某个键的所有值
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
