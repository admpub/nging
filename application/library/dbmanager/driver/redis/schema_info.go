package redis

import (
	"strconv"
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
	Name           string
	Value          string
	parsedKeyspace map[string]int64
}

type Infos struct {
	Name  string
	Attrs []*InfoKV
}

func NewInfos(name string, attrs ...*InfoKV) *Infos {
	return &Infos{
		Name:  name,
		Attrs: attrs,
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

func (a *InfoKV) ParseKeyspace() map[string]int64 {
	if a.parsedKeyspace != nil {
		return a.parsedKeyspace
	}
	if !strings.HasPrefix(a.Name, `db`) {
		return nil
	}
	a.parsedKeyspace = map[string]int64{}
	//keys=5,expires=0,avg_ttl=0
	for _, v := range strings.Split(a.Value, `,`) {
		kv := strings.SplitN(v, `=`, 2)
		var n int64
		if len(kv) > 1 {
			n, _ = strconv.ParseInt(kv[1], 10, 64)
		}
		a.parsedKeyspace[kv[0]] = n
	}
	return a.parsedKeyspace
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
		infoList[index].Attrs = append(infoList[index].Attrs, &InfoKV{
			Name:  kv[0],
			Value: kv[1],
		})
	}
	return infoList
}
