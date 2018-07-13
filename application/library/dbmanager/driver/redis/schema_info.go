package redis

import (
	"strings"
)

type InfoSection struct {
	Map map[string]string
	Idx []string
}

func (a *InfoSection) Add(key string, val string) *InfoSection {
	a.Idx = append(a.Idx, key)
	a.Map[key] = val
	return a
}

type Info struct {
	Map map[string]*InfoSection
	Idx []string
}

func (a *Info) Add(sectionName string, sectionData *InfoSection) *Info {
	a.Idx = append(a.Idx, sectionName)
	a.Map[sectionName] = sectionData
	return a
}

func (a *Info) MustSection(sectionName string) *InfoSection {
	section, exists := a.Map[sectionName]
	if !exists {
		section = NewInfoSection()
		a.Map[sectionName] = section
	}
	return section
}

func NewInfo() *Info {
	info := &Info{
		Map: map[string]*InfoSection{},
		Idx: []string{},
	}
	return info
}

func NewInfoSection() *InfoSection {
	section := &InfoSection{
		Map: map[string]string{},
		Idx: []string{},
	}
	return section
}

type InfoKV struct {
	Name  string
	Value string
}

type Infos struct {
	Name    string
	Configs []*InfoKV
}

func NewInfos(name string, configs ...*InfoKV) *Infos {
	return &Infos{
		Name:    name,
		Configs: configs,
	}
}

func ParseInfo(infoText string) *Info {
	info := NewInfo()
	infoText = strings.TrimSpace(infoText)
	rows := strings.Split(infoText, "\n")
	var sectionName string
	sectionPrefix := `# `
	for _, row := range rows {
		row = strings.TrimSpace(row)
		if len(row) == 0 {
			continue
		}
		if strings.HasPrefix(row, sectionPrefix) {
			sectionName = strings.TrimPrefix(row, sectionPrefix)
			section := NewInfoSection()
			info.Add(sectionName, section)
			continue
		}
		kv := strings.SplitN(row, `:`, 2)
		if len(kv) < 2 {
			kv = append(kv, ``)
		}
		info.MustSection(sectionName).Add(kv[0], kv[1])
	}
	return info
}

func ParseInfos(infoText string) []*Infos {
	infoList := []*Infos{}
	infoText = strings.TrimSpace(infoText)
	rows := strings.Split(infoText, "\n")
	var sectionName string
	sectionPrefix := `# `
	indexes := map[string]int{}
	for _, row := range rows {
		row = strings.TrimSpace(row)
		if len(row) == 0 {
			continue
		}
		if strings.HasPrefix(row, sectionPrefix) {
			sectionName = strings.TrimPrefix(row, sectionPrefix)
			indexes[sectionName] = len(infoList)
			infoList = append(infoList, NewInfos(sectionName))
			continue
		}
		kv := strings.SplitN(row, `:`, 2)
		if len(kv) < 2 {
			kv = append(kv, ``)
		}
		index, ok := indexes[sectionName]
		if !ok {
			index = len(infoList)
			indexes[sectionName] = index
			infoList = append(infoList, NewInfos(sectionName))
		}
		infoList[index].Configs = append(infoList[index].Configs, &InfoKV{
			Name:  kv[0],
			Value: kv[1],
		})
	}
	return infoList
}
