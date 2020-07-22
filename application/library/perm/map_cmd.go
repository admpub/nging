package perm

import "strings"

//ParseCmd 解析用户获取的权限
func (m *Map) ParseCmd(permCmds string) *Map {
	perms := strings.Split(permCmds, `,`)
	for _, _perm := range perms {
		arr := strings.Split(_perm, `/`)
		result := m.V
		for _, a := range arr {
			if _, y := result[a]; !y {
				result[a] = NewMap()
			}
			result = result[a].V
		}
	}
	return m
}

//CheckCmd 检测权限
func (m *Map) CheckCmd(perm string) bool {
	arr := strings.Split(perm, `/`)
	result := m.V
	for _, a := range arr {
		v, y := result[a]
		if !y {
			return false
		}
		result = v.V
	}
	return true
}
