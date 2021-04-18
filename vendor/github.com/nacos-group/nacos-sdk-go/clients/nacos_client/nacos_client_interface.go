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

package nacos_client

import (
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/common/http_agent"
)

//go:generate mockgen -destination mock_nacos_client_interface.go -package nacos_client -source=./nacos_client_interface.go

type INacosClient interface {

	//SetClientConfig is use to set nacos client config
	SetClientConfig(constant.ClientConfig) error
	//SetServerConfig is use to set nacos server config
	SetServerConfig([]constant.ServerConfig) error
	//GetClientConfig use to get client config
	GetClientConfig() (constant.ClientConfig, error)
	//GetServerConfig use to get server config
	GetServerConfig() ([]constant.ServerConfig, error)
	//SetHttpAgent use to set http agent
	SetHttpAgent(http_agent.IHttpAgent) error
	//GetHttpAgent use to get http agent
	GetHttpAgent() (http_agent.IHttpAgent, error)
}
