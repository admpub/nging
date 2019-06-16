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

package perm

import (
	"regexp"
	"strings"

	"github.com/admpub/nging/application/registry/navigate"
)

type Map struct {
	V   map[string]*Map
	Nav *navigate.Item
}

//Import 导入菜单（用户缓存结果）
func (m *Map) Import(navList *navigate.List) *Map {
	if navList == nil {
		return m
	}
	for _, nav := range *navList {
		item := NewMap()
		item.Nav = nav
		m.V[nav.Action] = item
		item.Import(nav.Children)
	}
	return m
}

func BuildPermActions(values []string) string {
	var permActions string

	if len(values) > 0 && values[0] == `*` {
		permActions = `*`
		return permActions
	}
	var prefix string
	for _, v := range values {
		length := len(v)
		var suffix string
		if length > 2 {
			suffix = v[length-2:]
		}
		if suffix == `/*` {
			if len(prefix) > 0 {
				prefix += `|`
			}
			prefix += regexp.QuoteMeta(v[0 : length-2])
			if len(permActions) > 0 {
				permActions += `,`
			}
			permActions += v
			continue
		}
		if len(prefix) > 0 {
			re := regexp.MustCompile(`^(` + prefix + `)`)
			if re.MatchString(v) {
				continue
			}
		}
		if len(permActions) > 0 {
			permActions += `,`
		}
		permActions += v
	}
	return permActions
}

//Parse 解析用户获取的权限
func (m *Map) Parse(permActions string, navTree *Map) *Map {
	perms := strings.Split(permActions, `,`)
	for _, perm := range perms {
		arr := strings.Split(perm, `/`)
		amap := m
		result := m.V
		var spath string
		for _, a := range arr {
			if len(spath) == 0 {
				spath = a
			}
			if mp, y := navTree.V[spath]; y {
				amap.Nav = m.Nav
				spath = ``
				amap = mp
			} else {
				if spath != a {
					spath += `/` + a
				}
			}
			if _, y := result[a]; !y {
				result[a] = NewMap()
			}
			result = result[a].V
		}
	}
	return m
}

//Check 检测权限
func (m *Map) Check(perm string) bool {
	if m.Nav != nil && m.Nav.Unlimited {
		return true
	}
	arr := strings.Split(perm, `/`)
	result := m.V
	for _, a := range arr {
		v, y := result[a]
		if !y {
			return false
		}
		if v.Nav != nil && v.Nav.Unlimited {
			return true
		}
		if _, y := v.V[`*`]; y {
			return true
		}
		result = v.V
	}
	return true
}

func NewMap() *Map {
	return &Map{V: map[string]*Map{}}
}
