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

package perm

import (
	"regexp"
	"strings"

	"github.com/admpub/nging/v5/application/registry/navigate"
)

type Map struct {
	V      map[string]*Map `json:",omitempty" xml:",omitempty"`
	Nav    *navigate.Item  `json:",omitempty" xml:",omitempty"`
	cached *Map
}

// Import 导入菜单（用户缓存结果）
func (m *Map) Import(navList *navigate.List) *Map {
	if navList == nil {
		return m
	}
	for _, nav := range *navList {
		item := NewMap(m.cached)
		item.Nav = nav
		if _, ok := m.V[nav.Action]; ok {
			panic(`The navigate name conflicts. Already existed name: ` + nav.Action)
		}
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

// Parse 解析用户获取的权限
func (m *Map) Parse(permActions string) *Map {
	perms := strings.Split(permActions, `,`)
	for _, perm := range perms {
		arr := strings.Split(perm, `/`)
		tree := m.cached
		result := m.V
		var spath string
		for _, a := range arr {
			if _, y := result[a]; !y {
				result[a] = NewMap(m.cached)
			}
			if a == `*` { //"*"是最后一个字符
				break
			}
			if mp, y := tree.V[a]; y {
				result[a].Nav = mp.Nav
				spath = ``
				tree = mp
			} else {
				if len(spath) > 0 {
					spath += `/`
				}
				spath += a
				if mp, y := tree.V[spath]; y {
					result[a].Nav = mp.Nav
					spath = ``
					tree = mp
				}
			}
			result = result[a].V
		}
	}
	return m
}

func (m *Map) checkByNav(perm string) bool {
	if m.Nav != nil && m.Nav.Unlimited {
		return true
	}
	arr := strings.Split(perm, `/`)
	navResult := m.V
	var prefix string
	for _, a := range arr {
		key := prefix + a
		navV, hasNav := navResult[key]
		if !hasNav {
			prefix += key + `/`
			continue
		}
		prefix = ``
		if navV.Nav != nil && navV.Nav.Unlimited {
			return true
		}
		navResult = navV.V
	}
	return false
}

// Check 检测权限
// perm: /a/b/c
func (m *Map) Check(perm string) bool {
	if m.cached == nil || m == m.cached {
		return m.checkByNav(perm)
	}
	return m.checkByChecked(perm, m.cached)
}

func (m *Map) checkByChecked(perm string, nav *Map) bool {
	if m.Nav != nil && m.Nav.Unlimited {
		return true
	}
	arr := strings.Split(perm, `/`)
	result := m.V
	navResult := nav.V
	hasPerm := true
	var prefix string
	for _, a := range arr {
		key := prefix + a
		navV, hasNav := navResult[key]
		if !hasNav {
			if hasPerm {
				var v *Map
				v, hasPerm = result[a]
				if hasPerm {
					if v.Nav != nil && v.Nav.Unlimited {
						return true
					}
					if _, y := v.V[`*`]; y {
						return true
					}
					result = v.V
				}
			}
			prefix += key + `/`
			continue
		}
		prefix = ``
		if !hasPerm {
			if navV.Nav != nil && navV.Nav.Unlimited {
				return true
			}
			navResult = navV.V
			continue
		}
		var v *Map
		v, hasPerm = result[a]
		if hasPerm {
			if v.Nav != nil && v.Nav.Unlimited {
				return true
			}
			if _, y := v.V[`*`]; y {
				return true
			}
			result = v.V
		}
		if navV.Nav != nil && navV.Nav.Unlimited {
			return true
		}
		navResult = navV.V
	}
	return hasPerm
}

func NewMap(cached *Map) *Map {
	return &Map{V: map[string]*Map{}, cached: cached}
}
