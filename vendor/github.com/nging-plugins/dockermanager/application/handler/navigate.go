package handler

import (
	"github.com/admpub/nging/v5/application/library/role"
	"github.com/admpub/nging/v5/application/registry/navigate"
	"github.com/nging-plugins/dockermanager/application/library/dockerclient"
	"github.com/webx-top/echo"
)

var nav = &navigate.List{
	&navigate.Item{
		Display: true,
		Name:    `Docker管理`,
		Action:  `docker/base`,
		Icon:    `road`,
		Children: &navigate.List{
			{
				Display: true,
				Name:    `Docker信息`,
				Action:  `index`,
			},
			// -- registry --
			{
				Display: false,
				Name:    `仓库登录`,
				Action:  `registry/login`,
			},
			// -- image --
			{
				Display: true,
				Name:    `镜像管理`,
				Action:  `image/index`,
			},
			{
				Display: false,
				Name:    `拉取新镜像`,
				Action:  `image/add`,
			},
			/*
				{
					Display: false,
					Name:    `修改镜像`,
					Action:  `image/edit`,
				},
			*/
			{
				Display: false,
				Name:    `查看镜像详情`,
				Action:  `image/detail/:id`,
			},
			{
				Display: false,
				Name:    `删除镜像`,
				Action:  `image/delete/:id`,
			},
			{
				Display: false,
				Name:    `拉取镜像更新`,
				Action:  `image/pull`,
			},
			{
				Display: false,
				Name:    `清理镜像`,
				Action:  `image/prune`,
			},
			{
				Display: false,
				Name:    `下载(备份)镜像`,
				Action:  `image/download/:id`,
			},
			{
				Display: false,
				Name:    `载入(还原)镜像`,
				Action:  `image/load`,
			},
			{
				Display: false,
				Name:    `导入容器快照镜像`,
				Action:  `image/import`,
			},
			/*
				{
					Display: false,
					Name:    `构建镜像`,
					Action:  `image/build`,
				},
			*/
			// -- container --
			{
				Display: true,
				Name:    `容器管理`,
				Action:  `container/index`,
			},
			{
				Display: false,
				Name:    `新建容器`,
				Action:  `container/add`,
			},
			{
				Display: false,
				Name:    `修改容器`,
				Action:  `container/edit/:id`,
			},
			{
				Display: false,
				Name:    `查看容器详情`,
				Action:  `container/detail/:id`,
			},
			{
				Display: false,
				Name:    `删除容器`,
				Action:  `container/delete/:id`,
			},
			{
				Display: false,
				Name:    `清理容器`,
				Action:  `container/prune`,
			},
			{
				Display: false,
				Name:    `强行停止(kill)容器`,
				Action:  `container/kill/:id`,
			},
			{
				Display: false,
				Name:    `启动容器`,
				Action:  `container/start/:id`,
			},
			{
				Display: false,
				Name:    `停止容器`,
				Action:  `container/stop/:id`,
			},
			{
				Display: false,
				Name:    `重启容器`,
				Action:  `container/restart/:id`,
			},
			{
				Display: false,
				Name:    `暂停/恢复容器`,
				Action:  `container/pause/:id`,
			},
			// {
			// 	Display: false,
			// 	Name:    `容器TOP信息`,
			// 	Action:  `container/top/:id`,
			// },
			// {
			// 	Display: false,
			// 	Name:    `容器统计数据接口`,
			// 	Action:  `container/stats/:id`,
			// },
			{
				Display: false,
				Name:    `容器统计信息`, // 依赖 “容器TOP信息” 和 “容器统计数据接口”
				Action:  `container/statsPage/:id`,
			},
			{
				Display: false,
				Name:    `容器文件导入`,
				Action:  `container/fileImport/:id`,
			},
			{
				Display: false,
				Name:    `容器文件导出`,
				Action:  `container/fileExport/:id`,
			},
			{
				Display: false,
				Name:    `容器文件管理`,
				Action:  `container/file/:id`,
			},
			{
				Display: false,
				Name:    `导出容器快照镜像`,
				Action:  `container/export/:id`,
			},
			{
				Display: false,
				Name:    `容器终端`,
				Action:  `container/pty/:id`,
			},
			{
				Display: false,
				Name:    `查看容器日志`,
				Action:  `container/logs/:id`,
			},
			// -- network --
			{
				Display: true,
				Name:    `网络管理`,
				Action:  `network/index`,
			},
			{
				Display: false,
				Name:    `新建网络`,
				Action:  `network/add`,
			},
			/*
				{
					Display: false,
					Name:    `连接网络`,
					Action:  `network/connect`,
				},
				{
					Display: false,
					Name:    `断开网络`,
					Action:  `network/disconnect`,
				},
			*/
			{
				Display: false,
				Name:    `查看网络详情`,
				Action:  `network/detail/:id`,
			},
			{
				Display: false,
				Name:    `删除网络`,
				Action:  `network/delete/:id`,
			},
			// -- volume --
			{
				Display: true,
				Name:    `存储卷管理`,
				Action:  `volume/index`,
			},
			{
				Display: false,
				Name:    `新建存储卷`,
				Action:  `volume/add`,
			},
			{
				Display: false,
				Name:    `查看存储卷详情`,
				Action:  `volume/detail/:id`,
			},
			{
				Display: false,
				Name:    `删除存储卷`,
				Action:  `volume/delete/:id`,
			},
			{
				Display: false,
				Name:    `清理存储卷`,
				Action:  `volume/prune`,
			},
			// -- compose --
			{
				Display: true,
				Name:    `Compose管理`,
				Action:  `compose/index`,
			},
			{
				Display: false,
				Name:    `新建Compose项目`,
				Action:  `compose/add`,
			},
			{
				Display: false,
				Name:    `修改Compose项目`,
				Action:  `compose/edit/:id`,
			},
			{
				Display: false,
				Name:    `查看Compose项目详情`,
				Action:  `compose/detail/:id`,
			},
			{
				Display: false,
				Name:    `删除Compose项目`,
				Action:  `compose/delete/:id`,
			},
			{
				Display: false,
				Name:    `查看Compose项目容器列表`,
				Action:  `compose/listContainers/:id`,
			},
			{
				Display: false,
				Name:    `停止Compose项目`,
				Action:  `compose/stop/:id`,
			},
			{
				Display: false,
				Name:    `启动Compose项目`,
				Action:  `compose/start/:id`,
			},
			{
				Display: false,
				Name:    `Compose项目容器扩容`,
				Action:  `compose/scale/:id/:service`,
			},
		},
	},
	&navigate.Item{
		Display: true,
		Name:    `Swarm管理`,
		Action:  `docker/swarm`,
		Icon:    `road`,
		Children: &navigate.List{
			{
				Display: true,
				Name:    `Swarm信息`,
				Action:  `index`,
			},
			{
				Display: false,
				Name:    `初始化Swarm集群`,
				Action:  `init`,
			},
			{
				Display: false,
				Name:    `加入Swarm集群`,
				Action:  `join`,
			},
			{
				Display: false,
				Name:    `脱离Swarm集群`,
				Action:  `leave`,
			},
			// -- node --
			{
				Display: true,
				Name:    `节点列表`,
				Action:  `node/index`,
			},
			{
				Display: false,
				Name:    `添加节点`,
				Action:  `node/add`,
			},
			{
				Display: false,
				Name:    `修改节点`,
				Action:  `node/edit/:id`,
			},
			{
				Display: false,
				Name:    `查看节点详情`,
				Action:  `node/detail/:id`,
			},
			{
				Display: false,
				Name:    `删除节点`,
				Action:  `node/delete/:id`,
			},
			// -- service --
			{
				Display: true,
				Name:    `服务列表`,
				Action:  `service/index`,
			},
			{
				Display: false,
				Name:    `新建服务`,
				Action:  `service/add`,
			},
			{
				Display: false,
				Name:    `修改服务`,
				Action:  `service/edit/:id`,
			},
			{
				Display: false,
				Name:    `查看服务详情`,
				Action:  `service/detail/:id`,
			},
			{
				Display: false,
				Name:    `删除服务`,
				Action:  `service/delete/:id`,
			},
			{
				Display: false,
				Name:    `回滚服务`,
				Action:  `service/rollback/:id`,
			},
			{
				Display: false,
				Name:    `查看服务日志`,
				Action:  `service/logs/:id`,
			},
			// -- task --
			{
				Display: true,
				Name:    `任务列表`,
				Action:  `task/index`,
			},
			{
				Display: false,
				Name:    `查看任务详情`,
				Action:  `task/detail/:id`,
			},
			{
				Display: false,
				Name:    `查看任务日志`,
				Action:  `task/logs/:id`,
			},
			// -- config --
			{
				Display: true,
				Name:    `配置列表`,
				Action:  `config/index`,
			},
			{
				Display: false,
				Name:    `新建配置`,
				Action:  `config/add`,
			},
			{
				Display: false,
				Name:    `修改配置`,
				Action:  `config/edit/:id`,
			},
			{
				Display: false,
				Name:    `查看配置详情`,
				Action:  `config/detail/:id`,
			},
			{
				Display: false,
				Name:    `删除配置`,
				Action:  `config/delete/:id`,
			},
			// -- secret --
			{
				Display: true,
				Name:    `密钥列表`,
				Action:  `secret/index`,
			},
			{
				Display: false,
				Name:    `新建密钥`,
				Action:  `secret/add`,
			},

			{
				Display: false,
				Name:    `修改密钥`,
				Action:  `secret/edit/:id`,
			},
			{
				Display: false,
				Name:    `查看密钥详情`,
				Action:  `secret/detail/:id`,
			},
			{
				Display: false,
				Name:    `删除密钥`,
				Action:  `secret/delete/:id`,
			},
			// -- stack --
			{
				Display: true,
				Name:    `Stack管理`,
				Action:  `stack/index`,
			},
			{
				Display: false,
				Name:    `新建Stack项目`,
				Action:  `stack/add`,
			},
			{
				Display: false,
				Name:    `修改Stack项目`,
				Action:  `stack/edit/:id`,
			},
			{
				Display: false,
				Name:    `查看Stack项目详情`,
				Action:  `stack/detail/:id`,
			},
			{
				Display: false,
				Name:    `删除Stack项目`,
				Action:  `stack/delete/:id`,
			},
			{
				Display: false,
				Name:    `查看Stack项目服务列表`,
				Action:  `stack/listServices/:id`,
			},
			{
				Display: false,
				Name:    `查看Stack项目任务列表`,
				Action:  `stack/listTasks/:id`,
			},
			{
				Display: false,
				Name:    `停止Stack项目`,
				Action:  `stack/stop/:id`,
			},
			{
				Display: false,
				Name:    `启动Stack项目`,
				Action:  `stack/start/:id`,
			},
		},
	},
}

var Project = navigate.NewProject(`Docker`, `docker`, `/docker/base/index`, nav)

func init() {
	navigate.ProjectAdd(-1, Project)
	echo.Set(`DockerContainerExec`, dockerclient.Exec)
	role.RegisterAuthDependency(`/docker/base/container/statsPage/:id`, `/docker/base/container/top/:id`, `/docker/base/container/stats/:id`)
}
