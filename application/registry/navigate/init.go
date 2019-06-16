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

package navigate

import (
	"github.com/admpub/events"
	"github.com/admpub/events/emitter"
)

var EmptyList = List{}

//TopNavigate 顶部导航菜单
var TopNavigate = &List{
	{
		DisplayOnMenu: true,
		Name:          `设置`,
		Action:        `manager`,
		Icon:          `gear`,
		Children: &List{
			{
				DisplayOnMenu: true,
				Name:          `系统设置`,
				Action:        `settings`,
			},
			{
				DisplayOnMenu: true,
				Name:          `用户管理`,
				Action:        `user`,
			},
			{
				DisplayOnMenu: true,
				Name:          `添加用户`,
				Action:        `user_add`,
				Icon:          `plus`,
			},
			{
				DisplayOnMenu: true,
				Name:          `角色管理`,
				Action:        `role`,
			},
			{
				DisplayOnMenu: true,
				Name:          `添加角色`,
				Action:        `role_add`,
				Icon:          `plus`,
			},
			{
				DisplayOnMenu: true,
				Name:          `邀请码`,
				Action:        `invitation`,
			},
			{
				DisplayOnMenu: true,
				Name:          `验证码`,
				Action:        `verification`,
			},
			{
				DisplayOnMenu: true,
				Name:          `清理缓存`,
				Action:        `clear_cache`,
				Target:        `ajax`,
			},
			{
				DisplayOnMenu: false,
				Name:          `修改用户`,
				Action:        `user_edit`,
			},
			{
				DisplayOnMenu: false,
				Name:          `删除用户`,
				Action:        `user_delete`,
			},
			{
				DisplayOnMenu: false,
				Name:          `修改角色`,
				Action:        `role_edit`,
			},
			{
				DisplayOnMenu: false,
				Name:          `删除角色`,
				Action:        `role_delete`,
			},
			{
				DisplayOnMenu: false,
				Name:          `添加邀请码`,
				Action:        `invitation_add`,
			},
			{
				DisplayOnMenu: false,
				Name:          `修改邀请码`,
				Action:        `invitation_edit`,
			},
			{
				DisplayOnMenu: false,
				Name:          `删除邀请码`,
				Action:        `invitation_delete`,
			},
			{
				DisplayOnMenu: false,
				Name:          `删除验证码`,
				Action:        `verification_delete`,
			},
			{
				DisplayOnMenu: false,
				Name:          `上传图片`,
				Action:        `upload/:type`,
			},
		},
	},
	{
		DisplayOnMenu: true,
		Name:          `工具箱`,
		Action:        `tool`,
		Icon:          `suitcase`,
		Children: &List{
			{
				DisplayOnMenu: true,
				Name:          `IP归属地`,
				Action:        `ip`,
			},
			{
				DisplayOnMenu: true,
				Name:          `Base64解码`,
				Action:        `base64`,
			},
			{
				DisplayOnMenu: true,
				Name:          `URL解码`,
				Action:        `url`,
			},
			{
				DisplayOnMenu: true,
				Name:          `时间戳转换`,
				Action:        `timestamp`,
			},
		},
	},
}

var defaultNavigate = &List{}

//LeftNavigate 左边导航菜单
var LeftNavigate = defaultNavigate

func init() {
	emitter.DefaultCondEmitter.On(`beforeRun`, events.Callback(func(_ events.Event) error {
		ProjectInitURLsIdent()
		return nil
	}))
}
