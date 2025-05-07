package manager

import (
	"github.com/admpub/nging/v5/application/handler/manager/file"
	"github.com/admpub/nging/v5/application/handler/tool"
	"github.com/coscms/webcore/library/navigate"
	"github.com/webx-top/echo"
)

// TopNavigate 顶部导航菜单
var TopNavigate = &navigate.List{
	{
		Display: true,
		Name:    echo.T(`设置`),
		Action:  `manager`,
		Icon:    `gear`,
		Children: (&navigate.List{
			{
				Display: true,
				Name:    echo.T(`系统设置`),
				Action:  `settings`,
			},
			//元数据操作
			{
				Display: true,
				Name:    echo.T(`元数据`),
				Action:  `kv`,
			},
			{
				Display: false,
				Name:    echo.T(`添加元数据`),
				Action:  `kv_add`,
			},
			{
				Display: false,
				Name:    echo.T(`修改元数据`),
				Action:  `kv_edit`,
			},
			{
				Display: false,
				Name:    echo.T(`删除元数据`),
				Action:  `kv_delete`,
			},
			//警报收信账号操作
			{
				Display: true,
				Name:    echo.T(`告警通知`),
				Action:  `alert_recipient`,
			},
			{
				Display: false,
				Name:    echo.T(`添加告警通知账号`),
				Action:  `alert_recipient_add`,
			},
			{
				Display: false,
				Name:    echo.T(`修改告警通知账号`),
				Action:  `alert_recipient_edit`,
			},
			{
				Display: false,
				Name:    echo.T(`删除告警通知账号`),
				Action:  `alert_recipient_delete`,
			},
			//告警通知专题
			{
				Display: false,
				Name:    echo.T(`告警通知专题`),
				Action:  `alert_topic`,
			},
			{
				Display: false,
				Name:    echo.T(`关联收信账号`),
				Action:  `alert_topic_add`,
			},
			{
				Display: false,
				Name:    echo.T(`修改收信账号`),
				Action:  `alert_topic_edit`,
			},
			{
				Display: false,
				Name:    echo.T(`取消关联收信账号`),
				Action:  `alert_topic_delete`,
			},
			{
				Display: false,
				Name:    echo.T(`测试发送告警通知`),
				Action:  `alert_recipient_test`,
			},
			//用户管理
			{
				Display: true,
				Name:    echo.T(`用户管理`),
				Action:  `user`,
				Group:   `admin`,
			},
			{
				Display: false,
				Name:    echo.T(`添加用户`),
				Action:  `user_add`,
				Icon:    `plus`,
				Group:   `admin`,
			},
			{
				Display: false,
				Name:    echo.T(`修改用户`),
				Action:  `user_edit`,
				Group:   `admin`,
			},
			{
				Display: false,
				Name:    echo.T(`删除用户`),
				Action:  `user_delete`,
				Group:   `admin`,
			},
			{
				Display: false,
				Name:    echo.T(`踢下线`),
				Action:  `user_kick`,
				Group:   `admin`,
			},
			//角色管理
			{
				Display: true,
				Name:    echo.T(`角色管理`),
				Action:  `role`,
				Group:   `admin`,
			},
			{
				Display: false,
				Name:    echo.T(`添加角色`),
				Action:  `role_add`,
				Icon:    `plus`,
				Group:   `admin`,
			},
			{
				Display: false,
				Name:    echo.T(`修改角色`),
				Action:  `role_edit`,
				Group:   `admin`,
			},
			{
				Display: false,
				Name:    echo.T(`删除角色`),
				Action:  `role_delete`,
				Group:   `admin`,
			},
			{
				Display: true,
				Name:    echo.T(`登录日志`),
				Action:  `login_log`,
				Group:   `admin`,
			},
			{
				Display: false,
				Name:    echo.T(`删除登录日志`),
				Action:  `login_log_delete`,
				Group:   `admin`,
			},
			//邀请码管理
			{
				Display: true,
				Name:    echo.T(`邀请码`),
				Action:  `invitation`,
				Group:   `admin`,
			},
			{
				Display: false,
				Name:    echo.T(`添加邀请码`),
				Action:  `invitation_add`,
				Group:   `admin`,
			},
			{
				Display: false,
				Name:    echo.T(`修改邀请码`),
				Action:  `invitation_edit`,
				Group:   `admin`,
			},
			{
				Display: false,
				Name:    echo.T(`删除邀请码`),
				Action:  `invitation_delete`,
				Group:   `admin`,
			},
			//验证码管理
			{
				Display: true,
				Name:    echo.T(`验证码`),
				Action:  `verification`,
				Group:   `admin`,
			},
			{
				Display: false,
				Name:    echo.T(`删除验证码`),
				Action:  `verification_delete`,
				Group:   `admin`,
			},
			{
				Display: false,
				Name:    echo.T(`上传图片`),
				Action:  `upload`,
				Group:   `admin`,
			},
			{
				Display: true,
				Name:    echo.T(`本地附件`),
				Action:  `uploaded/file`,
				Group:   `file`,
			},
			{
				Display: false,
				Name:    echo.T(`合并文件`),
				Action:  `uploaded/merged`,
				Group:   `file`,
			},
			{
				Display: false,
				Name:    echo.T(`分片文件`),
				Action:  `uploaded/chunk`,
				Group:   `file`,
			},
			{
				Display: true,
				Name:    echo.T(`清理缓存`),
				Action:  `clear_cache`,
				Target:  `ajax`,
				Group:   `cache`,
			},
			{
				Display: true,
				Name:    echo.T(`重载环境变量`),
				Action:  `reload_env`,
				Target:  `ajax`,
				Group:   `cache`,
			},
			{
				Display: true,
				Name:    echo.T(`后台oAuth应用`),
				Action:  `oauth_app/index`,
				Group:   `admin`,
			},
			{
				Display: false,
				Name:    echo.T(`添加后台oAuth应用`),
				Action:  `oauth_app/add`,
				Group:   `admin`,
			},
			{
				Display: false,
				Name:    echo.T(`修改后台oAuth应用`),
				Action:  `oauth_app/edit/:id`,
				Group:   `admin`,
			},
			{
				Display: false,
				Name:    echo.T(`删除后台oAuth应用`),
				Action:  `oauth_app/delete/:id`,
				Group:   `admin`,
			},
			{
				Display: false,
				Name:    echo.T(`程序升级`),
				Action:  `upgrade`,
			},
		}).Add(-1, file.TopNavigate...),
	},
	{
		Display:  true,
		Name:     echo.T(`工具箱`),
		Action:   `tool`,
		Icon:     `suitcase`,
		Children: (&navigate.List{}).Add(-1, tool.TopNavigate...),
	},
}
