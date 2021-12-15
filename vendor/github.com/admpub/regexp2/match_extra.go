package regexp2

// Slice 匹配所有结果到一个切片中
func (m *Match) Slice() (sv []string) {
	groups := m.Groups()
	for index, group := range groups {
		if index == 0 {
			continue
		}
		for _, str := range group.Captures {
			sv = append(sv, str.String())
		}
	}
	return
}

// Slice2 将匹配到的所有结果按组分配到切片中
func (m *Match) Slice2() (sv [][]string) {
	groups := m.Groups()
	for index, group := range groups {
		if index == 0 {
			continue
		}
		svc := make([]string, len(group.Captures))
		for idx, str := range group.Captures {
			svc[idx] = str.String()
		}
		sv = append(sv, svc)
	}
	return
}
