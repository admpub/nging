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

package manager

import (
	"github.com/admpub/nging/application/cmd/event"
	_ "github.com/admpub/nging/application/handler/caddy"
	_ "github.com/admpub/nging/application/handler/collector"
	_ "github.com/admpub/nging/application/handler/database"
	_ "github.com/admpub/nging/application/handler/download"
	_ "github.com/admpub/nging/application/handler/frp"
	_ "github.com/admpub/nging/application/handler/ftp"
	_ "github.com/admpub/nging/application/handler/server"
	_ "github.com/admpub/nging/application/handler/task"
	_ "github.com/admpub/nging/application/handler/term"
	"github.com/admpub/nging/application/registry/navigate"
)

func init() {
	event.SupportManager = true
	navigate.LeftNavigate = navigate.List{
		{
			DisplayOnMenu: true,
			Name:          `网站管理`,
			Action:        `caddy`,
			Icon:          `sitemap`,
			Children: navigate.List{
				{
					DisplayOnMenu: false,
					Name:          `Caddy日志`,
					Action:        `log_show`,
				},
				{
					DisplayOnMenu: true,
					Name:          `网站列表`,
					Action:        `vhost`,
				},
				{
					DisplayOnMenu: true,
					Name:          `添加网站`,
					Action:        `vhost_add`,
					Icon:          `plus`,
				},
				{
					DisplayOnMenu: false,
					Name:          `重启Caddy`,
					Action:        `restart`,
					Icon:          ``,
				},
				{
					DisplayOnMenu: false,
					Name:          `停止Caddy`,
					Action:        `stop`,
					Icon:          ``,
				},
				{
					DisplayOnMenu: false,
					Name:          `查看网站日志`,
					Action:        `vhost_log`,
				},
				{
					DisplayOnMenu: false,
					Name:          `配置表单`,
					Action:        `addon_form`,
					Icon:          ``,
				},
				{
					DisplayOnMenu: false,
					Name:          `修改网站`,
					Action:        `vhost_edit`,
					Icon:          ``,
				},
				{
					DisplayOnMenu: false,
					Name:          `删除网站`,
					Action:        `vhost_delete`,
					Icon:          ``,
				},
				{
					DisplayOnMenu: false,
					Name:          `管理网站文件`,
					Action:        `vhost_file`,
					Icon:          ``,
				},
				{
					DisplayOnMenu: false,
					Name:          `生成Caddyfile`,
					Action:        `vhost_build`,
					Icon:          ``,
				},

				{
					DisplayOnMenu: true,
					Name:          `分组管理`,
					Action:        `group`,
				},
				{
					DisplayOnMenu: true,
					Name:          `添加分组`,
					Action:        `group_add`,
					Icon:          `plus`,
				},
				{
					DisplayOnMenu: false,
					Name:          `修改分组`,
					Action:        `group_edit`,
					Icon:          ``,
				},
				{
					DisplayOnMenu: false,
					Name:          `删除分组`,
					Action:        `group_delete`,
					Icon:          ``,
				},
			},
		},
		{
			DisplayOnMenu: true,
			Name:          `服务器`,
			Action:        `server`,
			Icon:          `desktop`,
			Children: navigate.List{
				{
					DisplayOnMenu: true,
					Name:          `服务器信息`,
					Action:        `sysinfo`,
				},
				{
					DisplayOnMenu: true,
					Name:          `网络端口`,
					Action:        `netstat`,
				},
				{
					DisplayOnMenu: true,
					Name:          `执行命令`,
					Action:        `cmd`,
				},
				{
					DisplayOnMenu: true,
					Name:          `服务管理`,
					Action:        `service`,
				},
				{
					DisplayOnMenu: false,
					Name:          `查看Nging日志`,
					Action:        `log`,
				},
				{
					DisplayOnMenu: false,
					Name:          `查看进程详情`,
					Action:        `process/:pid`,
				},
				{
					DisplayOnMenu: false,
					Name:          `杀死进程`,
					Action:        `procskill/:pid`,
				},
				{
					DisplayOnMenu: false,
					Name:          `命令对话`,
					Action:        `cmdSend/*`,
				},
				{
					DisplayOnMenu: false,
					Name:          `发送命令`,
					Action:        `cmdSendWS`,
				},
				{
					DisplayOnMenu: true,
					Name:          `进程值守`,
					Action:        `daemon_index`,
				},
				{
					DisplayOnMenu: false,
					Name:          `进程值守日志`,
					Action:        `daemon_log`,
				},
				{
					DisplayOnMenu: false,
					Name:          `添加值守配置`,
					Action:        `daemon_add`,
				},
				{
					DisplayOnMenu: false,
					Name:          `修改值守配置`,
					Action:        `daemon_edit`,
				},
				{
					DisplayOnMenu: false,
					Name:          `删除值守配置`,
					Action:        `daemon_delete`,
				},
			},
		},
		{
			DisplayOnMenu: true,
			Name:          `FTP账号`,
			Action:        `ftp`,
			Icon:          `users`,
			Children: navigate.List{
				{
					DisplayOnMenu: true,
					Name:          `账号管理`,
					Action:        `account`,
				},
				{
					DisplayOnMenu: true,
					Name:          `添加账号`,
					Action:        `account_add`,
					Icon:          `plus`,
				},
				{
					DisplayOnMenu: true,
					Name:          `用户组`,
					Action:        `group`,
				},
				{
					DisplayOnMenu: true,
					Name:          `添加用户组`,
					Action:        `group_add`,
					Icon:          `plus`,
				},
				{
					DisplayOnMenu: false,
					Name:          `修改账号`,
					Action:        `account_edit`,
					Icon:          ``,
				},
				{
					DisplayOnMenu: false,
					Name:          `删除账号`,
					Action:        `account_delete`,
					Icon:          ``,
				},
				{
					DisplayOnMenu: false,
					Name:          `修改分组`,
					Action:        `group_edit`,
					Icon:          ``,
				},
				{
					DisplayOnMenu: false,
					Name:          `删除分组`,
					Action:        `group_delete`,
					Icon:          ``,
				},
				{
					DisplayOnMenu: false,
					Name:          `重启FTP服务`,
					Action:        `restart`,
					Icon:          ``,
				},
				{
					DisplayOnMenu: false,
					Name:          `关闭FTP服务`,
					Action:        `stop`,
					Icon:          ``,
				},
				{
					DisplayOnMenu: false,
					Name:          `查看FTP日志`,
					Action:        `log`,
					Icon:          ``,
				},
			},
		},
		{
			DisplayOnMenu: true,
			Name:          `数据采集`,
			Action:        `collector`,
			Icon:          `truck`,
			Children: navigate.List{
				{
					DisplayOnMenu: true,
					Name:          `规则管理`,
					Action:        `rule`,
				},
				{
					DisplayOnMenu: true,
					Name:          `新建规则`,
					Action:        `rule_add`,
					Icon:          `plus`,
				},
				{
					DisplayOnMenu: true,
					Name:          `历史记录`,
					Action:        `history`,
				},
				{
					DisplayOnMenu: false,
					Name:          `查看历史内容`,
					Action:        `history_view`,
				},
				{
					DisplayOnMenu: false,
					Name:          `删除历史记录`,
					Action:        `history_delete`,
				},
				{
					DisplayOnMenu: true,
					Name:          `数据导出`,
					Action:        `export`,
				},
				{
					DisplayOnMenu: false,
					Name:          `添加导出规则`,
					Action:        `export_add`,
				},
				{
					DisplayOnMenu: false,
					Name:          `修改导出规则`,
					Action:        `export_edit`,
				},
				{
					DisplayOnMenu: false,
					Name:          `删除导出规则`,
					Action:        `export_delete`,
				},
				{
					DisplayOnMenu: false,
					Name:          `导出日志管理`,
					Action:        `export_log`,
				},
				{
					DisplayOnMenu: false,
					Name:          `查看导出日志`,
					Action:        `export_log_view/:id`,
				},
				{
					DisplayOnMenu: false,
					Name:          `删除导出日志`,
					Action:        `export_log_delete`,
				},
				{
					DisplayOnMenu: false,
					Name:          `更改导出日志状态`,
					Action:        `export_edit_status`,
				},
				{
					DisplayOnMenu: true,
					Name:          `分组管理`,
					Action:        `group`,
				},
				{
					DisplayOnMenu: true,
					Name:          `添加分组`,
					Action:        `group_add`,
					Icon:          `plus`,
				},
				{
					DisplayOnMenu: false,
					Name:          `修改规则`,
					Action:        `rule_edit`,
					Icon:          ``,
				},
				{
					DisplayOnMenu: false,
					Name:          `删除规则`,
					Action:        `rule_delete`,
					Icon:          ``,
				},
				{
					DisplayOnMenu: false,
					Name:          `规则测试`,
					Action:        `rule_collect`,
					Icon:          ``,
				},
				{
					DisplayOnMenu: false,
					Name:          `修改分组`,
					Action:        `group_edit`,
					Icon:          ``,
				},
				{
					DisplayOnMenu: false,
					Name:          `删除分组`,
					Action:        `group_delete`,
					Icon:          ``,
				},
				{
					DisplayOnMenu: false,
					Name:          `测试正则表达式`,
					Action:        `regexp_test`,
					Icon:          ``,
				},
			},
		},
		{
			DisplayOnMenu: true,
			Name:          `计划任务`,
			Action:        `task`,
			Icon:          `clock-o`,
			Children: navigate.List{
				{
					DisplayOnMenu: true,
					Name:          `任务管理`,
					Action:        `index`,
				},
				{
					DisplayOnMenu: true,
					Name:          `新建任务`,
					Action:        `add`,
					Icon:          `plus`,
				},
				{
					DisplayOnMenu: true,
					Name:          `分组管理`,
					Action:        `group`,
				},
				{
					DisplayOnMenu: true,
					Name:          `添加分组`,
					Action:        `group_add`,
					Icon:          `plus`,
				},
				{
					DisplayOnMenu: true,
					Name:          `测试邮件`,
					Action:        `email_test`,
				},
				{
					DisplayOnMenu: false,
					Name:          `修改任务`,
					Action:        `edit`,
				},
				{
					DisplayOnMenu: false,
					Name:          `删除任务`,
					Action:        `delete`,
				},
				{
					DisplayOnMenu: false,
					Name:          `启动任务`,
					Action:        `start`,
				},
				{
					DisplayOnMenu: false,
					Name:          `暂停任务`,
					Action:        `pause`,
				},
				{
					DisplayOnMenu: false,
					Name:          `立即执行任务`,
					Action:        `run`,
				},
				{
					DisplayOnMenu: false,
					Name:          `退出任务`,
					Action:        `exit`,
				},
				{
					DisplayOnMenu: false,
					Name:          `启动历史任务`,
					Action:        `start_history`,
				},
				{
					DisplayOnMenu: false,
					Name:          `修改分组`,
					Action:        `group_edit`,
				},
				{
					DisplayOnMenu: false,
					Name:          `删除分组`,
					Action:        `group_delete`,
				},
				{
					DisplayOnMenu: false,
					Name:          `日志列表`,
					Action:        `log`,
				},
				{
					DisplayOnMenu: false,
					Name:          `日志详情`,
					Action:        `log_view/:id`,
				},
				{
					DisplayOnMenu: false,
					Name:          `删除日志`,
					Action:        `log_delete`,
				},
			},
		},
		{
			DisplayOnMenu: true,
			Name:          `离线下载`,
			Action:        `download`,
			Icon:          `download`,
			Children: navigate.List{
				{
					DisplayOnMenu: true,
					Name:          `下载管理`,
					Action:        `index.html`,
				},
				{
					DisplayOnMenu: true,
					Name:          `文件管理`,
					Action:        `file`,
				},
				{
					DisplayOnMenu: false,
					Name:          `总进度信息`,
					Action:        `progress.json`,
				},
				{
					DisplayOnMenu: false,
					Name:          `添加任务`,
					Action:        `add_task`,
				},
				{
					DisplayOnMenu: false,
					Name:          `删除任务`,
					Action:        `remove_task`,
				},
				{
					DisplayOnMenu: false,
					Name:          `启动任务`,
					Action:        `start_task`,
				},
				{
					DisplayOnMenu: false,
					Name:          `停止任务`,
					Action:        `stop_task`,
				},
				{
					DisplayOnMenu: false,
					Name:          `启动所有任务`,
					Action:        `start_all_task`,
				},
				{
					DisplayOnMenu: false,
					Name:          `停止所有任务`,
					Action:        `stop_all_task`,
				},
				{
					DisplayOnMenu: false,
					Name:          `单个文件进度信息`,
					Action:        `progress/*`,
				},
			},
		},
		{
			DisplayOnMenu: true,
			Name:          `数据库`,
			Action:        `db`,
			Icon:          `table`,
			Children: navigate.List{
				{
					DisplayOnMenu: true,
					Name:          `数据库账号`,
					Action:        `account`,
				},
				{
					DisplayOnMenu: true,
					Name:          `添加账号`,
					Action:        `account_add`,
					Icon:          `plus`,
				},
				{
					DisplayOnMenu: false,
					Name:          `修改账号`,
					Action:        `account_edit`,
				},
				{
					DisplayOnMenu: false,
					Name:          `删除账号`,
					Action:        `account_delete`,
				},
				{
					DisplayOnMenu: true,
					Name:          `连接数据库`,
					Action:        ``,
				},
				{
					DisplayOnMenu: true,
					Name:          `表结构同步`,
					Action:        `schema_sync`,
				},
				{
					DisplayOnMenu: false,
					Name:          `新增同步方案`,
					Action:        `schema_sync_add`,
				},
				{
					DisplayOnMenu: false,
					Name:          `修改同步方案`,
					Action:        `schema_sync_edit`,
				},
				{
					DisplayOnMenu: false,
					Name:          `删除同步方案`,
					Action:        `schema_sync_delete`,
				},
				{
					DisplayOnMenu: false,
					Name:          `预览表结构差异`,
					Action:        `schema_sync_preview`,
				},
				{
					DisplayOnMenu: false,
					Name:          `执行表结构同步`,
					Action:        `schema_sync_run`,
				},
				{
					DisplayOnMenu: false,
					Name:          `表结构同步日志列表`,
					Action:        `schema_sync_log/:id`,
				},
				{
					DisplayOnMenu: false,
					Name:          `表结构同步日志详情`,
					Action:        `schema_sync_log_view/:id`,
				},
				{
					DisplayOnMenu: false,
					Name:          `删除表结构同步日志`,
					Action:        `schema_sync_log_delete`,
				},
			},
		},
		{
			DisplayOnMenu: true,
			Name:          `内网穿透`,
			Action:        `frp`,
			Icon:          `road`,
			Children: navigate.List{
				{
					DisplayOnMenu: true,
					Name:          `服务端配置`,
					Action:        `server_index`,
				},
				{
					DisplayOnMenu: false,
					Name:          `查看服务端日志`,
					Action:        `server_log`,
				},
				{
					DisplayOnMenu: false,
					Name:          `添加服务端配置`,
					Action:        `server_add`,
				},
				{
					DisplayOnMenu: false,
					Name:          `修改服务端配置`,
					Action:        `server_edit`,
				},
				{
					DisplayOnMenu: false,
					Name:          `删除服务端配置`,
					Action:        `server_delete`,
				},
				{
					DisplayOnMenu: true,
					Name:          `客户端配置`,
					Action:        `client_index`,
				},
				{
					DisplayOnMenu: false,
					Name:          `查看客户端日志`,
					Action:        `client_log`,
				},
				{
					DisplayOnMenu: false,
					Name:          `添加客户端配置`,
					Action:        `client_add`,
				},
				{
					DisplayOnMenu: false,
					Name:          `修改客户端配置`,
					Action:        `client_edit`,
				},
				{
					DisplayOnMenu: false,
					Name:          `删除客户端配置`,
					Action:        `client_delete`,
				},
				{
					DisplayOnMenu: true,
					Name:          `分组管理`,
					Action:        `group_index`,
				},
				{
					DisplayOnMenu: false,
					Name:          `添加分组`,
					Action:        `group_add`,
				},
				{
					DisplayOnMenu: false,
					Name:          `修改分组`,
					Action:        `group_edit`,
				},
				{
					DisplayOnMenu: false,
					Name:          `删除分组`,
					Action:        `group_delete`,
				},
				{
					DisplayOnMenu: true,
					Name:          `统计图表`,
					Action:        `dashboard`,
					Target:        `_blank`,
				},
				{
					DisplayOnMenu: false,
					Name:          `重启服务端`,
					Action:        `server_restart`,
				},
				{
					DisplayOnMenu: false,
					Name:          `关闭服务端`,
					Action:        `server_stop`,
				},
				{
					DisplayOnMenu: false,
					Name:          `重启客户端`,
					Action:        `client_restart`,
				},
				{
					DisplayOnMenu: false,
					Name:          `关闭客户端`,
					Action:        `client_stop`,
				},
				{
					DisplayOnMenu: false,
					Name:          `配置表单`,
					Action:        `addon_form`,
				},
			},
		},
		{
			DisplayOnMenu: true,
			Name:          `SSH管理`,
			Action:        `term`,
			Icon:          `terminal`,
			Children: navigate.List{
				{
					DisplayOnMenu: true,
					Name:          `账号管理`,
					Action:        `account`,
				},
				{
					DisplayOnMenu: true,
					Name:          `添加账号`,
					Action:        `account_add`,
					Icon:          `plus`,
				},
				{
					DisplayOnMenu: true,
					Name:          `分组管理`,
					Action:        `group`,
				},
				{
					DisplayOnMenu: true,
					Name:          `添加分组`,
					Action:        `group_add`,
					Icon:          `plus`,
				},
				{
					DisplayOnMenu: false,
					Name:          `修改账号`,
					Action:        `account_edit`,
					Icon:          ``,
				},
				{
					DisplayOnMenu: false,
					Name:          `删除账号`,
					Action:        `account_delete`,
					Icon:          ``,
				},
				{
					DisplayOnMenu: false,
					Name:          `修改分组`,
					Action:        `group_edit`,
					Icon:          ``,
				},
				{
					DisplayOnMenu: false,
					Name:          `删除分组`,
					Action:        `group_delete`,
					Icon:          ``,
				},
				{
					DisplayOnMenu: false,
					Name:          `SSH操作`,
					Action:        `client`,
					Icon:          ``,
				},
				{
					DisplayOnMenu: false,
					Name:          `SFTP操作`,
					Action:        `sftp`,
					Icon:          ``,
				},
			},
		},
	}
}
