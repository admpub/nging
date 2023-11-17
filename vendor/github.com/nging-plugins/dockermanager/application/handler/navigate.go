package handler

import "github.com/admpub/nging/v5/application/registry/navigate"

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
			{
				Display: true,
				Name:    `镜像管理`,
				Action:  `image/index`,
			},
			{
				Display: true,
				Name:    `容器管理`,
				Action:  `container/index`,
			},
			{
				Display: true,
				Name:    `网络管理`,
				Action:  `network/index`,
			},
			{
				Display: true,
				Name:    `存储卷管理`,
				Action:  `volume/index`,
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
				Display: true,
				Name:    `节点列表`,
				Action:  `node/index`,
			},
			{
				Display: true,
				Name:    `服务列表`,
				Action:  `service/index`,
			},
			{
				Display: true,
				Name:    `任务列表`,
				Action:  `task/index`,
			},
			{
				Display: true,
				Name:    `配置列表`,
				Action:  `config/index`,
			},
			{
				Display: true,
				Name:    `密钥列表`,
				Action:  `secret/index`,
			},
		},
	},
}

var Project = navigate.NewProject(`Docker`, `docker`, `/docker/base/index`, nav)

func init() {
	navigate.ProjectAdd(1, Project)
}
