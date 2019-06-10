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

func (c *Block) SetContentGenerator(content func(echo.Context) error) *Block {
	c.content = content
	return c
}

type Blocks []*Block

func (c Blocks) Ready(block echo.Context) error {
	for _, blk := range c {
		if err := blk.Ready(block); err != nil {
			return nil
		}
	}
	return nil
}

var blocks = Blocks{}

func BlockRegister(block ...*Block) {
	blocks = append(blocks, block...)
}

func BlockAll() Blocks {
	return blocks
}
