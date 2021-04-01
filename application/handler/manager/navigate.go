package manager

import (
	"github.com/admpub/nging/application/handler/manager/file"
	"github.com/admpub/nging/application/registry/navigate"
)

//TopNavigate 顶部导航菜单
var TopNavigate = &navigate.List{
	{
		DisplayOnMenu: true,
		Name:          `设置`,
		Action:        `manager`,
		Icon:          `gear`,
		Children: (&navigate.List{
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
			{
				DisplayOnMenu: false,
				Name:          `测试发送警报信息`,
				Action:        `alert_recipient_test`,
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
			{
				DisplayOnMenu: true,
				Name:          `登录日志`,
				Action:        `login_log`,
			},
			{
				DisplayOnMenu: false,
				Name:          `删除登录日志`,
				Action:        `login_log_delete`,
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
				Action:        `upload`,
			},
			{
				DisplayOnMenu: true,
				Name:          `本地附件`,
				Action:        `uploaded/file`,
			},
			{
				DisplayOnMenu: false,
				Name:          `合并文件`,
				Action:        `uploaded/merged`,
			},
			{
				DisplayOnMenu: false,
				Name:          `分片文件`,
				Action:        `uploaded/chunk`,
			},
			{
				DisplayOnMenu: true,
				Name:          `清理缓存`,
				Action:        `clear_cache`,
				Target:        `ajax`,
			},
		}).Add(-1, file.TopNavigate...),
	},
	{
		DisplayOnMenu: true,
		Name:          `工具箱`,
		Action:        `tool`,
		Icon:          `suitcase`,
		Children: &navigate.List{
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
