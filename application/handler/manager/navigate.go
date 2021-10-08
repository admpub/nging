package manager

import (
	"github.com/admpub/nging/v3/application/handler/manager/file"
	"github.com/admpub/nging/v3/application/handler/tool"
	"github.com/admpub/nging/v3/application/registry/navigate"
)

//TopNavigate 顶部导航菜单
var TopNavigate = &navigate.List{
	{
		Display: true,
		Name:    `设置`,
		Action:  `manager`,
		Icon:    `gear`,
		Children: (&navigate.List{
			{
				Display: true,
				Name:    `系统设置`,
				Action:  `settings`,
			},
			//元数据操作
			{
				Display: true,
				Name:    `元数据`,
				Action:  `kv`,
			},
			{
				Display: false,
				Name:    `添加元数据`,
				Action:  `kv_add`,
			},
			{
				Display: false,
				Name:    `修改元数据`,
				Action:  `kv_edit`,
			},
			{
				Display: false,
				Name:    `删除元数据`,
				Action:  `kv_delete`,
			},
			//警报收信账号操作
			{
				Display: true,
				Name:    `警报收信账号`,
				Action:  `alert_recipient`,
			},
			{
				Display: false,
				Name:    `添加警报收信账号`,
				Action:  `alert_recipient_add`,
			},
			{
				Display: false,
				Name:    `修改警报收信账号`,
				Action:  `alert_recipient_edit`,
			},
			{
				Display: false,
				Name:    `删除警报收信账号`,
				Action:  `alert_recipient_delete`,
			},
			//警报专题
			{
				Display: false,
				Name:    `警报专题`,
				Action:  `alert_topic`,
			},
			{
				Display: false,
				Name:    `关联收信账号`,
				Action:  `alert_topic_add`,
			},
			{
				Display: false,
				Name:    `修改收信账号`,
				Action:  `alert_topic_edit`,
			},
			{
				Display: false,
				Name:    `取消关联收信账号`,
				Action:  `alert_topic_delete`,
			},
			{
				Display: false,
				Name:    `测试发送警报信息`,
				Action:  `alert_recipient_test`,
			},
			//用户管理
			{
				Display: true,
				Name:    `用户管理`,
				Action:  `user`,
			},
			{
				Display: false,
				Name:    `添加用户`,
				Action:  `user_add`,
				Icon:    `plus`,
			},
			{
				Display: false,
				Name:    `修改用户`,
				Action:  `user_edit`,
			},
			{
				Display: false,
				Name:    `删除用户`,
				Action:  `user_delete`,
			},
			{
				Display: false,
				Name:    `踢下线`,
				Action:  `user_kick`,
			},
			//角色管理
			{
				Display: true,
				Name:    `角色管理`,
				Action:  `role`,
			},
			{
				Display: false,
				Name:    `添加角色`,
				Action:  `role_add`,
				Icon:    `plus`,
			},
			{
				Display: false,
				Name:    `修改角色`,
				Action:  `role_edit`,
			},
			{
				Display: false,
				Name:    `删除角色`,
				Action:  `role_delete`,
			},
			{
				Display: true,
				Name:    `登录日志`,
				Action:  `login_log`,
			},
			{
				Display: false,
				Name:    `删除登录日志`,
				Action:  `login_log_delete`,
			},
			//邀请码管理
			{
				Display: true,
				Name:    `邀请码`,
				Action:  `invitation`,
			},
			{
				Display: false,
				Name:    `添加邀请码`,
				Action:  `invitation_add`,
			},
			{
				Display: false,
				Name:    `修改邀请码`,
				Action:  `invitation_edit`,
			},
			{
				Display: false,
				Name:    `删除邀请码`,
				Action:  `invitation_delete`,
			},
			//验证码管理
			{
				Display: true,
				Name:    `验证码`,
				Action:  `verification`,
			},
			{
				Display: false,
				Name:    `删除验证码`,
				Action:  `verification_delete`,
			},
			{
				Display: false,
				Name:    `上传图片`,
				Action:  `upload`,
			},
			{
				Display: true,
				Name:    `本地附件`,
				Action:  `uploaded/file`,
			},
			{
				Display: false,
				Name:    `合并文件`,
				Action:  `uploaded/merged`,
			},
			{
				Display: false,
				Name:    `分片文件`,
				Action:  `uploaded/chunk`,
			},
			{
				Display: true,
				Name:    `清理缓存`,
				Action:  `clear_cache`,
				Target:  `ajax`,
			},
			{
				Display: true,
				Name:    `重载环境变量`,
				Action:  `reload_env`,
				Target:  `ajax`,
			},
		}).Add(-1, file.TopNavigate...),
	},
	{
		Display:  true,
		Name:     `工具箱`,
		Action:   `tool`,
		Icon:     `suitcase`,
		Children: (&navigate.List{}).Add(-1, tool.TopNavigate...),
	},
}
