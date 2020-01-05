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

package redis

import (
	"strconv"
	"strings"

	dbPagination "github.com/webx-top/db/lib/factory/pagination"
	"github.com/webx-top/echo"
	"github.com/webx-top/pagination"

	"github.com/admpub/nging/application/handler"
)

func NewValue(c echo.Context) *Value {
	return &Value{
		TotalRows:  -1,
		List:       []echo.KV{},
		NextOffset: `0`,
		context:    c,
	}
}

type Value struct {
	TotalRows  int
	List       []echo.KV
	NextOffset string
	context    echo.Context
	paging     *pagination.Pagination
	Text       string
}

func (a *Value) Add(key, value string) *Value {
	a.List = append(a.List, echo.KV{K: key, V: value})
	return a
}

func (a *Value) String() string {
	return a.Text
}

func (a *Value) SetPaging(paging *pagination.Pagination) *Value {
	a.paging = paging
	return a
}

func (a *Value) Paging(vkeys ...string) *pagination.Pagination {
	if a.paging != nil {
		return a.paging
	}
	var (
		pageKey = `vpage`
		sizeKey = `vsize`
		rowsKey = `vrows`
	)
	if len(vkeys) > 0 {
		pageKey = vkeys[0]
		if len(vkeys) > 1 {
			sizeKey = vkeys[1]
			if len(vkeys) > 2 {
				rowsKey = vkeys[2]
			}
		}
	}
	page := a.context.Formx(pageKey).Int()
	size := a.context.Formx(sizeKey).Int()
	if page < 1 {
		page = 1
	}
	if size < 1 || size > dbPagination.PageMaxSize {
		size = 50
	}
	delKeys := []string{}
	totalRows := a.context.Formx(rowsKey).Int()
	if totalRows < 1 {
		totalRows = a.TotalRows
	}
	pjax := a.context.PjaxContainer()
	if len(pjax) > 0 {
		delKeys = append(delKeys, `_pjax`)
	}
	a.paging = pagination.New(a.context).SetAll(``, totalRows, page, 10, size).SetURL(map[string]string{
		`rows`: rowsKey,
		`page`: pageKey,
		`size`: sizeKey,
	}, delKeys...)
	return a.paging
}

func (a *Value) CursorPaging(vkeys ...string) *pagination.Pagination {
	if a.paging != nil {
		return a.paging
	}
	var (
		currOffsetKey = `voffset`
		prevOffsetKey = `vprev`
	)
	if len(vkeys) > 0 {
		currOffsetKey = vkeys[0]
		if len(vkeys) > 1 {
			prevOffsetKey = vkeys[1]
		}
	}
	_, _, _, a.paging = handler.PagingWithPagination(a.context)
	prevOffset := a.context.Form(currOffsetKey, `0`)
	q := a.context.Request().URL().Query()
	q.Del(currOffsetKey)
	q.Del(prevOffsetKey)
	q.Del(`_pjax`)
	a.paging.SetURL(`/db?`+q.Encode()+`&`+currOffsetKey+`={curr}&`+prevOffsetKey+`={prev}`).SetPosition(prevOffset, a.NextOffset, a.NextOffset)
	return a.paging
}

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
