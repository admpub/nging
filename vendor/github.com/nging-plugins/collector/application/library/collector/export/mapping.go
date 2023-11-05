/*
   Nging is a toolbox for webmasters
   Copyright (C) 2018-present  Wenhui Shen <swh@admpub.com>

   This program is free software: you can redistribute it and/or modify
   it under the terms of the GNU Affero General Public License as published
   by the Free Software Foundation, either version 3 of the License, or
   (at your option) any later version.

   This program is distributed in the hope that it will be useful,
   but WITHOUT ANY WARRANTY; without even the implied warranty of
   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
   GNU Affero General Public License for more details.

   You should have received a copy of the GNU Affero General Public License
   along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/

package export

import (
	"errors"
	"strings"

	"github.com/webx-top/com"
)

var ErrSameTableName = errors.New(`same table name`)

type Mapping struct {
	FromParent int    //如果大于0，则代表从上一层或上几层页面中取值
	FromTable  string //如果设置，则会从FromTable指定的“表的上一个插入操作”之后的结果集字段中取值
	FromField  string //来自结果集字段
	FromPipe   string //管道函数
	FromInput  string //输入框数据的原始内容
	ToTable    string //保存到表
	ToField    string //保存到字段
	ToInput    string //输入框数据的原始内容
}

func NewMapping(fromInput, filter string, toInput string) (*Mapping, error) {
	mapping := &Mapping{
		FromInput: fromInput,
		FromPipe:  filter,
		ToInput:   toInput,
	}
	if arr := strings.SplitN(mapping.ToInput, `.`, 2); len(arr) > 1 {
		mapping.ToTable = arr[0]
		mapping.ToField = arr[1]
	} else {
		mapping.ToField = mapping.ToInput
	}
	fromInputField := mapping.FromInput
	if strings.HasPrefix(fromInputField, `@`) {
		mapping.FromField = strings.TrimPrefix(fromInputField, `@`)
		if arr := strings.SplitN(mapping.ToInput, `.`, 2); len(arr) > 1 {
			mapping.FromTable = arr[0]
			if mapping.FromTable == mapping.ToTable {
				return nil, ErrSameTableName
			}
			mapping.FromField = arr[1]
		}
	} else {
		mapping.FromField = strings.TrimLeft(fromInputField, `./`)
		for strings.HasPrefix(mapping.FromField, `../`) {
			mapping.FromParent++
			mapping.FromField = strings.TrimPrefix(mapping.FromField, `../`)
			mapping.FromField = strings.TrimLeft(mapping.FromField, `./`)
		}
	}
	return mapping, nil
}

type Mappings struct {
	TableKeys  map[string][]int
	Slice      []*Mapping
	TableNames []string
}

func NewMappings() *Mappings {
	return &Mappings{
		TableKeys:  map[string][]int{},
		Slice:      []*Mapping{},
		TableNames: []string{},
	}
}

func (m *Mappings) addTableName(tableName string) {
	var exists bool
	for _, table := range m.TableNames {
		if table == tableName {
			exists = true
			break
		}
	}
	if !exists {
		m.TableNames = append(m.TableNames, tableName)
	}
}

func (m *Mappings) Add(mp *Mapping) *Mappings {
	if len(mp.FromTable) > 0 {
		m.addTableName(mp.FromTable)
	}
	if _, ok := m.TableKeys[mp.ToTable]; !ok {
		m.TableKeys[mp.ToTable] = []int{}
		m.addTableName(mp.ToTable)
	}
	m.TableKeys[mp.ToTable] = append(m.TableKeys[mp.ToTable], len(m.Slice))
	m.Slice = append(m.Slice, mp)
	return m
}

func (m *Mappings) String() (string, error) {
	b, err := com.JSONEncode(m)
	if err == nil {
		return string(b), nil
	}
	return ``, err
}
