/*
   Nging is a toolbox for webmasters
   Copyright (C) 2018-present Wenhui Shen <swh@admpub.com>

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

package navigate

import (
	"path"

	"github.com/webx-top/echo"
)

//Item 操作
type Item struct {
	Display    bool        `json:",omitempty" xml:",omitempty"` //是否在菜单上显示
	Name       string      `json:",omitempty" xml:",omitempty"` //名称
	Action     string      `json:",omitempty" xml:",omitempty"` //操作(一般为网址)
	Icon       string      `json:",omitempty" xml:",omitempty"` //图标
	Badge      string      `json:",omitempty" xml:",omitempty"` // <sup class="badge badge-danger">123</sup>
	Target     string      `json:",omitempty" xml:",omitempty"` //打开方式
	Unlimited  bool        `json:",omitempty" xml:",omitempty"` //是否不限制权限
	Attributes echo.KVList `json:",omitempty" xml:",omitempty"` //HTML标签a属性
	Children   *List       `json:",omitempty" xml:",omitempty"` //子菜单
}

func (a *Item) FullPath(parentPath string) string {
	if a == nil {
		return parentPath
	}
	return path.Join(parentPath, a.Action)
}

//List 操作列表
type List []*Item

func (a *List) FullPath(parentPath string) []string {
	var r []string
	if a == nil {
		return r
	}
	for _, nav := range *a {
		urlPath := path.Join(parentPath, nav.Action)
		//fmt.Println(`<FullPath>`, urlPath)
		if nav.Children == nil {
			r = append(r, urlPath)
			continue
		}
		r = append(r, nav.Children.FullPath(urlPath)...)
	}
	return r
}

//Remove 删除元素
func (a *List) Remove(index int) *List {
	if index < 0 {
		*a = (*a)[0:0]
		return a
	}
	size := len(*a)
	if size > index {
		if size > index+1 {
			*a = append((*a)[0:index], (*a)[index+1:]...)
		} else {
			*a = (*a)[0:index]
		}
	}
	return a
}

//Set 设置元素
func (a *List) Set(index int, list ...*Item) *List {
	if len(list) == 0 {
		return a
	}
	if index < 0 {
		*a = append(*a, list...)
		return a
	}
	size := len(*a)
	if size > index {
		(*a)[index] = list[0]
		if len(list) > 1 {
			a.Set(index+1, list[1:]...)
		}
		return a
	}
	for start, end := size, index-1; start < end; start++ {
		*a = append(*a, nil)
	}
	*a = append(*a, list...)
	return a
}

//Add 添加列表项
func (a *List) Add(index int, list ...*Item) *List {
	if len(list) == 0 {
		return a
	}
	if index < 0 {
		*a = append(*a, list...)
		return a
	}
	size := len(*a)
	if size > index {
		list = append(list, (*a)[index])
		(*a)[index] = list[0]
		if len(list) > 1 {
			a.Add(index+1, list[1:]...)
		}
		return a
	}
	for start, end := size, index-1; start < end; start++ {
		*a = append(*a, nil)
	}
	*a = append(*a, list...)
	return a
}

//Get 添加列表项
func (a *List) Get(index int) *Item {
	if len(*a) > index {
		return (*a)[index]
	}
	return nil
}

//Size 子项数量
func (a *List) Size() int {
	return len(*a)
}

//ChildrenBy 添加列表项
func (a *List) ChildrenBy(index int) *List {
	ls := a.Get(index)
	if ls == nil {
		return nil
	}
	return ls.Children
}

func (a *List) AddChild(action string, index int, list ...*Item) {
	for _, item := range *a {
		if item.Action == action {
			item.Children.Add(index, list...)
			break
		}
	}
}
