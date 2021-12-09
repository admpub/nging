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
	"database/sql"

	"github.com/admpub/nging/v3/application/model"
	"github.com/webx-top/echo"
)

func NewCard(content func(echo.Context) interface{}) *Card {
	return &Card{content: content}
}

type Card struct {
	IconName  string      //图标名称：fa-tasks
	IconColor string      //图标颜色：primary、success、danger、warning、info
	Short     string      //简称
	Name      string      //中文名称
	Summary   string      //说明
	Content   interface{} //数字等内容
	Hidden    sql.NullBool
	content   func(echo.Context) interface{}
	hidden    func(echo.Context) bool
}

func (c *Card) Build(ctx echo.Context) *Card {
	if c.Content == nil && c.content != nil {
		c.Content = c.content(ctx)
	}
	return c
}

func (c *Card) IsHidden(ctx echo.Context) (hidden bool) {
	v, ok := ctx.Internal().GetOk(`registry.dashboard.card.` + c.Short)
	if ok {
		hidden = v.(bool)
		return
	}
	if c.hidden != nil {
		hidden = c.hidden(ctx)
	} else {
		hidden = c.Hidden.Bool
	}
	ctx.Internal().Set(`registry.dashboard.card.`+c.Short, hidden)
	return
}

func (c *Card) SetHidden(hidden bool) *Card {
	c.Hidden.Bool = hidden
	c.Hidden.Valid = true
	return c
}

func (c *Card) SetContentGenerator(content func(echo.Context) interface{}) *Card {
	c.content = content
	return c
}

func (c *Card) SetHiddenDetector(hidden func(echo.Context) bool) *Card {
	c.hidden = hidden
	return c
}

type Cards []*Card

func (c *Cards) Build(ctx echo.Context) Cards {
	for _, card := range *c {
		card.Build(ctx)
	}
	return *c
}

// Remove 删除元素
func (c *Cards) Remove(index int) {
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

func (c *Cards) Add(index int, list ...*Card) {
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
func (c *Cards) Set(index int, list ...*Card) {
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

func (c *Cards) Size() int {
	return len(*c)
}

func (c *Cards) Search(cb func(*Card) bool) int {
	for index, card := range *c {
		if cb(card) {
			return index
		}
	}
	return -1
}

func (c *Cards) FindShort(short string) int {
	return c.Search(func(card *Card) bool {
		return card.Short == short
	})
}

func (c *Cards) RemoveByShort(tmpl string) {
	index := c.FindShort(tmpl)
	if index > -1 {
		c.Remove(index)
	}
}

var cards = &Cards{
	{
		IconName:  `fa-user`,
		IconColor: `success`,
		Short:     `USERS`,
		Name:      `用户数量`,
		Summary:   ``,
		content: func(ctx echo.Context) interface{} {
			//用户统计
			userMdl := model.NewUser(ctx)
			userCount, _ := userMdl.Count(nil)
			return userCount
		},
	},
}

func CardRegister(card ...*Card) {
	cards.Add(-1, card...)
}

func CardAdd(index int, card ...*Card) {
	cards.Add(index, card...)
}

//CardRemove 删除元素
func CardRemove(index int) {
	cards.Remove(index)
}

//CardSet 设置元素
func CardSet(index int, list ...*Card) {
	cards.Set(index, list...)
}

func CardAll(_ echo.Context) *Cards {
	return cards
}
