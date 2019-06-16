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

package navigate

import "strings"

var projects = NewProjects()

func ProjectAddNavList(name string, ident string, navList *List) {
	proj := ProjectGet(ident)
	if proj == nil {
		ProjectAdd(-1, NewProject(name, ident, navList))
		return
	}
	proj.NavList.Add(-1, *navList...)
}

func ProjectAdd(index int, list ...*ProjectItem) {
	projects.Add(index, list...)
}

func ProjectIdent(urlPath string) string {
	if projects.List == nil {
		return ``
	}
	arr := strings.Split(urlPath, `/`)
	var spath string
	for _, a := range arr {
		if len(spath) == 0 {
			spath = a
		}
		var hasPrefix bool
		for _, proj := range *projects.List {
			if proj.NavList == nil {
				continue
			}
			for _, nav := range *proj.NavList {
				if len(nav.Action) < len(spath) {
					continue
				}
				if nav.Action == spath {
					return proj.Ident
				}
				if strings.HasPrefix(nav.Action, spath+`/`) {
					hasPrefix = true
				}
			}
		}
		if hasPrefix && spath != a {
			spath += `/` + a
		}
	}
	return ``
}

func ProjectRemove(index int) {
	projects.Remove(index)
}

func ProjectSet(index int, list ...*ProjectItem) {
	projects.Set(index, list...)
}

func ProjectListAll() ProjectList {
	return *projects.List
}

func ProjectGet(ident string) *ProjectItem {
	return projects.Get(ident)
}

func NewProjects() *Projects {
	return &Projects{
		List: &ProjectList{},
		Hash: map[string]*ProjectItem{},
	}
}

func NewProject(name string, ident string, navList *List) *ProjectItem {
	return &ProjectItem{
		Name:    name,
		Ident:   ident,
		NavList: navList,
	}
}

type Projects struct {
	List *ProjectList
	Hash map[string]*ProjectItem
}

func (p *Projects) Get(ident string) *ProjectItem {
	if item, ok := p.Hash[ident]; ok {
		return item
	}
	return nil
}
func (p *Projects) Remove(index int) *Projects {
	if len(*p.List) <= index {
		return p
	}
	ident := (*p.List)[index].Ident
	p.List.Remove(index)
	if _, ok := p.Hash[ident]; ok {
		delete(p.Hash, ident)
	}
	return p
}
func (p *Projects) Add(index int, list ...*ProjectItem) *Projects {
	for _, item := range list {
		ident := item.Ident
		if _, ok := p.Hash[ident]; ok {
			panic(`Project already exists: ` + item.Ident)
		}
		p.Hash[ident] = item
	}
	p.List.Add(index, list...)
	return p
}
func (p *Projects) Set(index int, list ...*ProjectItem) *Projects {
	p.List.Set(index, list...)
	for _, item := range list {
		p.Hash[item.Ident] = item
	}
	return p
}

type ProjectList []*ProjectItem

type ProjectItem struct {
	Name    string
	Ident   string
	NavList *List
}

//Remove 删除元素
func (a *ProjectList) Remove(index int) *ProjectList {
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
func (a *ProjectList) Set(index int, list ...*ProjectItem) *ProjectList {
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
	for start := size; start < index; start++ {
		*a = append(*a, nil)
	}
	*a = append(*a, list...)
	return a
}

//Add 添加列表项
func (a *ProjectList) Add(index int, list ...*ProjectItem) *ProjectList {
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
	for start := size; start < index; start++ {
		*a = append(*a, nil)
	}
	*a = append(*a, list...)
	return a
}
