package perm

import "strings"

//ParseCmd 解析用户获取的权限
func (m *Map) ParseCmd(permCmds string) *Map {
	perms := strings.Split(permCmds, `,`)
	result := m.V
	for _, a := range perms {
		if _, y := result[a]; !y {
			result[a] = NewMap(m.cached)
		}
	}
	return m
}

//CheckCmd 检测权限
func (m *Map) CheckCmd(perm string) bool {
	_, y := m.V[perm]
	return y
}
