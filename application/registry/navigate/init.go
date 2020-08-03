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

import "github.com/webx-top/echo"

// EmptyList 空菜单列表
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
			//元数据操作
			{
				DisplayOnMenu: true,
				Name:          `元数据`,
				Action:        `kv`,
			},
			{
				DisplayOnMenu: false,
				Name:          `添加元数据`,
				Action:        `kv_add`,
			},
			{
				DisplayOnMenu: false,
				Name:          `修改元数据`,
				Action:        `kv_edit`,
			},
			{
				DisplayOnMenu: false,
				Name:          `删除元数据`,
				Action:        `kv_delete`,
			},
			//警报收信账号操作
			{
				DisplayOnMenu: true,
				Name:          `警报收信账号`,
				Action:        `alert_recipient`,
			},
			{
				DisplayOnMenu: false,
				Name:          `添加警报收信账号`,
				Action:        `alert_recipient_add`,
			},
			{
				DisplayOnMenu: false,
				Name:          `修改警报收信账号`,
				Action:        `alert_recipient_edit`,
			},
			{
				DisplayOnMenu: false,
				Name:          `删除警报收信账号`,
				Action:        `alert_recipient_delete`,
			},
			//警报专题
			{
				DisplayOnMenu: false,
				Name:          `警报专题`,
				Action:        `alert_topic`,
			},
			{
				DisplayOnMenu: false,
				Name:          `关联收信账号`,
				Action:        `alert_topic_add`,
			},
			{
				DisplayOnMenu: false,
				Name:          `修改收信账号`,
				Action:        `alert_topic_edit`,
			},
			{
				DisplayOnMenu: false,
				Name:          `取消关联收信账号`,
				Action:        `alert_topic_delete`,
			},
			//用户管理
			{
				DisplayOnMenu: true,
				Name:          `用户管理`,
				Action:        `user`,
			},
			{
				DisplayOnMenu: false,
				Name:          `添加用户`,
				Action:        `user_add`,
				Icon:          `plus`,
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
			//角色管理
			{
				DisplayOnMenu: true,
				Name:          `角色管理`,
				Action:        `role`,
			},
			{
				DisplayOnMenu: false,
				Name:          `添加角色`,
				Action:        `role_add`,
				Icon:          `plus`,
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
			//邀请码管理
			{
				DisplayOnMenu: true,
				Name:          `邀请码`,
				Action:        `invitation`,
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
			//验证码管理
			{
				DisplayOnMenu: true,
				Name:          `验证码`,
				Action:        `verification`,
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
			{
				DisplayOnMenu: true,
				Name:          `本地附件`,
				Action:        `uploaded_file`,
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
			{
				DisplayOnMenu: true,
				Name:          `附件网址替换`,
				Action:        `replaceurl`,
			},
		},
	},
}

var (
	defaultNavigate = &List{}
	topNavURLs      = map[string]int{}
)

func TopNavURLs() map[string]int {
	return topNavURLs
}

//LeftNavigate 左边导航菜单
var LeftNavigate = defaultNavigate

func init() {
	echo.On(`beforeRun`, func(_ echo.H) error {
		ProjectInitURLsIdent()
		for index, urlPath := range TopNavigate.FullPath(``) {
			topNavURLs[urlPath] = index
		}
		return nil
	})
}
