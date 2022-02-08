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

func NewButton(content func(echo.Context) error, tmpl ...string) *Button {
	var t string
	if len(tmpl) > 0 {
		t = tmpl[0]
	}
	return &Button{Tmpl: t, content: content}
}

type Button struct {
	Tmpl    string //模板文件
	content func(echo.Context) error
}

func (c *Button) Ready(ctx echo.Context) error {
	if c.content != nil {
		return c.content(ctx)
	}
	return nil
}

func (c *Button) SetContentGenerator(content func(echo.Context) error) *Button {
	c.content = content
	return c
}

func (c *Button) SetTmpl(tmpl string) *Button {
	c.Tmpl = tmpl
	return c
}

type Buttons []*Button

func (c *Buttons) Ready(ctx echo.Context) error {
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
func (c *Buttons) Remove(index int) {
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

func (c *Buttons) Add(index int, list ...*Button) {
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
func (c *Buttons) Set(index int, list ...*Button) {
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

func (c *Buttons) Size() int {
	return len(*c)
}

func (c *Buttons) Search(cb func(*Button) bool) int {
	for index, button := range *c {
		if cb(button) {
			return index
		}
	}
	return -1
}

func (c *Buttons) FindTmpl(tmpl string) int {
	return c.Search(func(button *Button) bool {
		return button.Tmpl == tmpl
	})
}

func (c *Buttons) RemoveByTmpl(tmpl string) {
	index := c.FindTmpl(tmpl)
	if index > -1 {
		c.Remove(index)
	}
}

var topButtons = Buttons{
	{
		Tmpl: `manager/topbutton/donation`,
	},
	{
		Tmpl: `manager/topbutton/language`,
	},
	{
		Tmpl: `manager/topbutton/source`,
	},
	{
		Tmpl: `manager/topbutton/bug-report`,
	},
}

func TopButtonRegister(topButton ...*Button) {
	topButtons.Add(-1, topButton...)
}

func TopButtonAdd(index int, topButton ...*Button) {
	topButtons.Add(index, topButton...)
}

//TopButtonRemove 删除元素
func TopButtonRemove(index int) {
	topButtons.Remove(index)
}

func TopButtonSearch(cb func(*Button) bool) int {
	return topButtons.Search(cb)
}

func TopButtonFindTmpl(tmpl string) int {
	return topButtons.FindTmpl(tmpl)
}

func TopButtonRemoveByTmpl(tmpl string) {
	topButtons.RemoveByTmpl(tmpl)
}

//TopButtonSet 设置元素
func TopButtonSet(index int, list ...*Button) {
	topButtons.Set(index, list...)
}

func TopButtonAll(_ echo.Context) Buttons {
	return topButtons
}
