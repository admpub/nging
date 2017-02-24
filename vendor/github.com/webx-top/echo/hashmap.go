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
	"encoding/json"
	"encoding/xml"
	"fmt"
	"strconv"
)

// Dump 输出对象和数组的结构信息
func Dump(m interface{}, args ...bool) (r string) {
	v, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		fmt.Printf("%v\n", err)
	}
	r = string(v)
	l := len(args)
	if l < 1 || args[0] {
		fmt.Println(r)
	}
	return
}

type H map[string]interface{}

// MarshalXML allows type H to be used with xml.Marshal
func (h H) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	start.Name = xml.Name{
		Space: ``,
		Local: `Map`,
	}
	if err := e.EncodeToken(start); err != nil {
		return err
	}
	for key, value := range h {
		elem := xml.StartElement{
			Name: xml.Name{Space: ``, Local: key},
			Attr: []xml.Attr{},
		}
		if err := e.EncodeElement(value, elem); err != nil {
			return err
		}
	}
	if err := e.EncodeToken(xml.EndElement{Name: start.Name}); err != nil {
		return err
	}
	return nil
}

// ToData conversion to *Data
func (h H) ToData() *Data {
	var info, zone, data interface{}
	if v, y := h["Data"]; y {
		data = v
	}
	if v, y := h["Zone"]; y {
		zone = v
	}
	if v, y := h["Info"]; y {
		info = v
	}
	var code int
	if v, y := h["Code"]; y {
		if c, y := v.(int); y {
			code = c
		}
	}
	return &Data{
		Code: code,
		Info: info,
		Zone: zone,
		Data: data,
	}
}

func (h H) DeepMerge(source H) {
	for k, value := range source {
		var (
			destinationValue interface{}
			ok               bool
		)
		if destinationValue, ok = h[k]; !ok {
			h[k] = value
			continue
		}
		sourceM, sourceOk := value.(H)
		destinationM, destinationOk := destinationValue.(H)
		if sourceOk && sourceOk == destinationOk {
			destinationM.DeepMerge(sourceM)
		} else {
			h[k] = value
		}
	}
}

type Mapx struct {
	Map   map[string]*Mapx `json:",omitempty"`
	Slice []*Mapx          `json:",omitempty"`
	Val   []string         `json:",omitempty"`
}

func NewMapx(data map[string][]string) *Mapx {
	m := &Mapx{}
	return m.Parse(data)
}

func (m *Mapx) Parse(data map[string][]string) *Mapx {
	m.Map = map[string]*Mapx{}
	for name, values := range data {
		names := FormNames(name)
		end := len(names) - 1
		v := m
		for idx, key := range names {
			if len(key) == 0 {

				if v.Slice == nil {
					v.Slice = []*Mapx{}
				}

				if idx == end {
					v.Slice = append(v.Slice, &Mapx{Val: values})
					continue
				}
				mapx := &Mapx{
					Map: map[string]*Mapx{},
				}
				v.Slice = append(v.Slice, mapx)
				v = mapx
				continue
			}
			if _, ok := v.Map[key]; !ok {
				if idx == end {
					v.Map[key] = &Mapx{Val: values}
					continue
				}
				v.Map[key] = &Mapx{
					Map: map[string]*Mapx{},
				}
				v = v.Map[key]
				continue
			}

			if idx == end {
				v.Map[key] = &Mapx{Val: values}
			} else {
				v = v.Map[key]
			}
		}
	}
	return m
}

func (m *Mapx) Value(names ...string) string {
	v := m.Values(names...)
	if v != nil {
		if len(v) > 0 {
			return v[0]
		}
	}
	return ``
}

func (m *Mapx) ValueOk(names ...string) (string, bool) {
	v, y := m.ValuesOk(names...)
	if y && v != nil {
		if len(v) > 0 {
			return v[0], true
		}
	}
	return ``, false
}

func (m *Mapx) ValuesOk(names ...string) ([]string, bool) {
	if len(names) == 0 {
		if m.Val == nil {
			return []string{}, false
		}
		return m.Val, true
	}
	v := m.Get(names...)
	if v != nil {
		return v.Val, true
	}
	return []string{}, false
}

func (m *Mapx) Values(names ...string) []string {
	if len(names) == 0 {
		if m.Val == nil {
			return []string{}
		}
		return m.Val
	}
	v := m.Get(names...)
	if v != nil {
		return v.Val
	}
	return []string{}
}

func (m *Mapx) Get(names ...string) *Mapx {
	v := m
	end := len(names) - 1
	for idx, key := range names {
		_, ok := v.Map[key]
		if !ok {
			if v.Slice == nil {
				return nil
			}
			i, err := strconv.Atoi(key)
			if err != nil {
				return nil
			}
			if i < 0 {
				return nil
			}
			if i < len(v.Slice) {
				v = v.Slice[i]
				continue
			}
			return nil
		}
		v = v.Map[key]

		if idx == end {
			return v
		}
	}
	return nil
}
