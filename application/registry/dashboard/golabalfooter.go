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

package dashboard

import (
	"github.com/webx-top/echo"
)

func NewGlobalFooter(content func(echo.Context) error) *GlobalFooter {
	return &GlobalFooter{content: content}
}

type GlobalFooter struct {
	Tmpl    string //模板文件
	content func(echo.Context) error
}

func (c *GlobalFooter) Ready(ctx echo.Context) error {
	if c.content != nil {
		return c.content(ctx)
	}
	return nil
}

func (c *GlobalFooter) SetContentGenerator(content func(echo.Context) error) *GlobalFooter {
	c.content = content
	return c
}

type GlobalFooters []*GlobalFooter

func (c *GlobalFooters) Ready(ctx echo.Context) error {
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
func (c *GlobalFooters) Remove(index int) {
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

func (c *GlobalFooters) Add(index int, list ...*GlobalFooter) {
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
func (c *GlobalFooters) Set(index int, list ...*GlobalFooter) {
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

func (c *GlobalFooters) Size() int {
	return len(*c)
}

var globalFooters = GlobalFooters{}

func GlobalFooterRegister(footer ...*GlobalFooter) {
	globalFooters.Add(-1, footer...)
}

func GlobalFooterAdd(index int, footer ...*GlobalFooter) {
	globalFooters.Add(index, footer...)
}

//GlobalFooterRemove 删除元素
func GlobalFooterRemove(index int) {
	globalFooters.Remove(index)
}

//GlobalFooterSet 设置元素
func GlobalFooterSet(index int, list ...*GlobalFooter) {
	globalFooters.Set(index, list...)
}

func GlobalFooterAll(_ echo.Context) GlobalFooters {
	return globalFooters
}
