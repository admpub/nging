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

func NewBlock(content func(echo.Context) error) *Block {
	return &Block{content: content}
}

type Block struct {
	Title   string `json:",omitempty" xml:",omitempty"` // 标题
	Ident   string `json:",omitempty" xml:",omitempty"` // 英文标识
	Extra   echo.H `json:",omitempty" xml:",omitempty"` // 附加数据
	Tmpl    string //模板文件
	Footer  string //末尾模版或JS代码
	content func(echo.Context) error
}

func (c *Block) Ready(ctx echo.Context) error {
	if c.content != nil {
		return c.content(ctx)
	}
	return nil
}

func (c *Block) SetTitle(title string) *Block {
	c.Title = title
	return c
}

func (c *Block) SetIdent(ident string) *Block {
	c.Ident = ident
	return c
}

func (c *Block) SetExtra(extra echo.H) *Block {
	c.Extra = extra
	return c
}

func (c *Block) SetExtraKV(key string, value interface{}) *Block {
	if c.Extra == nil {
		c.Extra = echo.H{}
	}
	c.Extra.Set(key, value)
	return c
}

func (c *Block) SetTmpl(tmpl string) *Block {
	c.Tmpl = tmpl
	return c
}

func (c *Block) SetFooter(footer string) *Block {
	c.Footer = footer
	return c
}

func (c *Block) SetContentGenerator(content func(echo.Context) error) *Block {
	c.content = content
	return c
}

type Blocks []*Block

func (c *Blocks) Ready(ctx echo.Context) error {
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
func (c *Blocks) Remove(index int) {
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

func (c *Blocks) Add(block ...*Block) {
	*c = append(*c, block...)
}

// Set 设置元素
func (c *Blocks) Set(index int, list ...*Block) {
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

func (c *Blocks) Size() int {
	return len(*c)
}

var blocks = Blocks{}

func BlockRegister(block ...*Block) {
	blocks.Add(block...)
}

//BlockRemove 删除元素
func BlockRemove(index int) {
	blocks.Remove(index)
}

//BlockSet 设置元素
func BlockSet(index int, list ...*Block) {
	blocks.Set(index, list...)
}

func BlockAll() Blocks {
	return blocks
}
