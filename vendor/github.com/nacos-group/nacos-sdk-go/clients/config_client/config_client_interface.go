/*
 * Copyright 1999-2020 Alibaba Group Holding Ltd.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package config_client

import (
	"github.com/nacos-group/nacos-sdk-go/model"
	"github.com/nacos-group/nacos-sdk-go/vo"
)

//go:generate mockgen -destination ../../mock/mock_config_client_interface.go -package mock -source=./config_client_interface.go

type IConfigClient interface {
	// GetConfig use to get config from nacos server
	// dataId  require
	// group   require
	// tenant ==>nacos.namespace optional
	GetConfig(param vo.ConfigParam) (string, error)

	// PublishConfig use to publish config to nacos server
	// dataId  require
	// group   require
	// content require
	// tenant ==>nacos.namespace optional
	PublishConfig(param vo.ConfigParam) (bool, error)

	// DeleteConfig use to delete config
	// dataId  require
	// group   require
	// tenant ==>nacos.namespace optional
	DeleteConfig(param vo.ConfigParam) (bool, error)

	// ListenConfig use to listen config change,it will callback OnChange() when config change
	// dataId  require
	// group   require
	// onchange require
	// tenant ==>nacos.namespace optional
	ListenConfig(params vo.ConfigParam) (err error)

	//CancelListenConfig use to cancel listen config change
	// dataId  require
	// group   require
	// tenant ==>nacos.namespace optional
	CancelListenConfig(params vo.ConfigParam) (err error)

	// SearchConfig use to search nacos config
	// search  require search=accurate--精确搜索  search=blur--模糊搜索
	// group   option
	// dataId  option
	// tenant ==>nacos.namespace optional
	// pageNo  option,default is 1
	// pageSize option,default is 10
	SearchConfig(param vo.SearchConfigParm) (*model.ConfigPage, error)
}
