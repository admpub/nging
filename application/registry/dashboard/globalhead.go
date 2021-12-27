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

package dashboard

import (
	"github.com/webx-top/echo"
)

func NewGlobalHead(tmpl string, handle ...func(echo.Context) error) *GlobalHead {
	var content func(echo.Context) error
	if len(handle) > 0 {
		content = handle[0]
	}
	return &GlobalHead{Tmpl: tmpl, content: content}
}

type GlobalHead struct {
	Tmpl    string //模板文件
	content func(echo.Context) error
}

func (c *GlobalHead) Ready(ctx echo.Context) error {
	if c.content != nil {
		return c.content(ctx)
	}
	return nil
}

func (c *GlobalHead) SetContentGenerator(content func(echo.Context) error) *GlobalHead {
	c.content = content
	return c
}

func (c *GlobalHead) SetTmpl(tmpl string) *GlobalHead {
	c.Tmpl = tmpl
	return c
}

type GlobalHeads []*GlobalHead

func (c *GlobalHeads) Ready(ctx echo.Context) error {
	for _, blk := range *c {
		if blk != nil {
			if err := blk.Ready(ctx); err != nil {
				return err
			}
		}
	}
	return nil
}

// Remove 删除元素
func (c *GlobalHeads) Remove(index int) {
	if index < 0 {
		*c = (*c)[0:0]
		return
	}
	size := c.Size()
	if size > index {
		if size > index+1 {
			*c = append((*c)[0:index], (*c)[index+1:]...)
		} else {
			*c = (*c)[0:index]
		}
	}
}

func (c *GlobalHeads) Add(index int, list ...*GlobalHead) {
	if len(list) == 0 {
		return
	}
	if index < 0 {
		*c = append(*c, list...)
		return
	}
	size := c.Size()
	if size > index {
		list = append(list, (*c)[index])
		(*c)[index] = list[0]
		if len(list) > 1 {
			c.Add(index+1, list[1:]...)
		}
		return
	}
	for start, end := size, index-1; start < end; start++ {
		*c = append(*c, nil)
	}
	*c = append(*c, list...)
}

// Set 设置元素
func (c *GlobalHeads) Set(index int, list ...*GlobalHead) {
	if len(list) == 0 {
		return
	}
	if index < 0 {
		*c = append(*c, list...)
		return
	}
	size := c.Size()
	if size > index {
		(*c)[index] = list[0]
		if len(list) > 1 {
			c.Set(index+1, list[1:]...)
		}
		return
	}
	for start, end := size, index-1; start < end; start++ {
		*c = append(*c, nil)
	}
	*c = append(*c, list...)
}

func (c *GlobalHeads) Size() int {
	return len(*c)
}

func (c *GlobalHeads) Search(cb func(*GlobalHead) bool) int {
	for index, head := range *c {
		if cb(head) {
			return index
		}
	}
	return -1
}

func (c *GlobalHeads) FindTmpl(tmpl string) int {
	return c.Search(func(head *GlobalHead) bool {
		return head.Tmpl == tmpl
	})
}

func (c *GlobalHeads) RemoveByTmpl(tmpl string) {
	index := c.FindTmpl(tmpl)
	if index > -1 {
		c.Remove(index)
	}
}

var globalHeads = GlobalHeads{}

func GlobalHeadRegister(footer ...*GlobalHead) {
	globalHeads.Add(-1, footer...)
}

func GlobalHeadAdd(index int, footer ...*GlobalHead) {
	globalHeads.Add(index, footer...)
}

//GlobalHeadRemove 删除元素
func GlobalHeadRemove(index int) {
	globalHeads.Remove(index)
}

//GlobalHeadSet 设置元素
func GlobalHeadSet(index int, list ...*GlobalHead) {
	globalHeads.Set(index, list...)
}

func GlobalHeadAll(_ echo.Context) GlobalHeads {
	return globalHeads
}
