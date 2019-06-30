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

func (c GlobalFooters) Ready(block echo.Context) error {
	for _, blk := range c {
		if blk != nil {
			if err := blk.Ready(block); err != nil {
				return err
			}
		}
	}
	return nil
}

var globalFooters = GlobalFooters{}

func GlobalFooterRegister(topButton ...*GlobalFooter) {
	globalFooters = append(globalFooters, topButton...)
}

//GlobalFooterRemove 删除元素
func GlobalFooterRemove(index int) {
	if index < 0 {
		globalFooters = globalFooters[0:0]
		return
	}
	size := len(globalFooters)
	if size > index {
		if size > index+1 {
			globalFooters = append(globalFooters[0:index], globalFooters[index+1:]...)
		} else {
			globalFooters = globalFooters[0:index]
		}
	}
}

//GlobalFooterSet 设置元素
func GlobalFooterSet(index int, list ...*GlobalFooter) {
	if len(list) == 0 {
		return
	}
	if index < 0 {
		globalFooters = append(globalFooters, list...)
		return
	}
	size := len(globalFooters)
	if size > index {
		globalFooters[index] = list[0]
		if len(list) > 1 {
			GlobalFooterSet(index+1, list[1:]...)
		}
		return
	}
	for start, end := size, index-1; start < end; start++ {
		globalFooters = append(globalFooters, nil)
	}
	globalFooters = append(globalFooters, list...)
}

func GlobalFooterAll() GlobalFooters {
	return globalFooters
}
