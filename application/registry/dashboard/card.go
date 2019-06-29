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
	"github.com/admpub/nging/application/model"
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
	content   func(echo.Context) interface{}
}

func (c *Card) Build(ctx echo.Context) *Card {
	if c.Content == nil && c.content != nil {
		c.Content = c.content(ctx)
	}
	return c
}

func (c *Card) SetContentGenerator(content func(echo.Context) interface{}) *Card {
	c.content = content
	return c
}

type Cards []*Card

func (c Cards) Build(ctx echo.Context) Cards {
	for _, card := range c {
		card.Build(ctx)
	}
	return c
}

var cards = Cards{
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
	cards = append(cards, card...)
}

//CardRemove 删除元素
func CardRemove(index int) {
	if index < 0 {
		cards = cards[0:0]
		return
	}
	size := len(cards)
	if size > index {
		if size > index+1 {
			cards = append(cards[0:index], cards[index+1:]...)
		} else {
			cards = cards[0:index]
		}
	}
}

//CardSet 设置元素
func CardSet(index int, list ...*Card) {
	if len(list) == 0 {
		return
	}
	if index < 0 {
		cards = append(cards, list...)
		return
	}
	size := len(cards)
	if size > index {
		cards[index] = list[0]
		if len(list) > 1 {
			CardSet(index+1, list[1:]...)
		}
		return
	}
	for start, end := size, index-1; start < end; start++ {
		cards = append(cards, nil)
	}
	cards = append(cards, list...)
}

func CardAll() Cards {
	return cards
}
