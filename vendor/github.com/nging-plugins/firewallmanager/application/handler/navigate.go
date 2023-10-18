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

package handler

import "github.com/admpub/nging/v5/application/registry/navigate"

var LeftNavigate = &navigate.Item{
	Display: true,
	Name:    `防火墙`,
	Action:  `firewall`,
	Icon:    `shield`,
	Children: &navigate.List{
		{
			Display: true,
			Name:    `静态规则`,
			Action:  `rule/static`,
		},
		{
			Display: false,
			Name:    `添加静态规则`,
			Action:  `rule/static_add`,
		},
		{
			Display: false,
			Name:    `修改静态规则`,
			Action:  `rule/static_edit`,
		},
		{
			Display: false,
			Name:    `应用静态规则`,
			Action:  `rule/static_apply`,
		},
		{
			Display: false,
			Name:    `临时封IP`,
			Action:  `rule/static_ban`,
		},
		{
			Display: false,
			Name:    `删除静态规则`,
			Action:  `rule/static_delete`,
		},
		{
			Display: true,
			Name:    `动态规则`,
			Action:  `rule/dynamic`,
		},
		{
			Display: false,
			Name:    `添加动态规则`,
			Action:  `rule/dynamic_add`,
		},
		{
			Display: false,
			Name:    `修改动态规则`,
			Action:  `rule/dynamic_edit`,
		},
		{
			Display: false,
			Name:    `删除动态规则`,
			Action:  `rule/dynamic_delete`,
		},
		{
			Display: false,
			Name:    `重启服务`,
			Action:  `service/restart`,
			Icon:    ``,
		},
		{
			Display: false,
			Name:    `关闭服务`,
			Action:  `service/stop`,
			Icon:    ``,
		},
		{
			Display: false,
			Name:    `查看动态`,
			Action:  `service/log`,
			Icon:    ``,
		},
	},
}
