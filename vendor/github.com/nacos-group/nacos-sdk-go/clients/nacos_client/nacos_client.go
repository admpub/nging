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
	"errors"
	"log"
	"os"
	"strconv"

	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/common/file"
	"github.com/nacos-group/nacos-sdk-go/common/http_agent"
)

type NacosClient struct {
	clientConfigValid  bool
	serverConfigsValid bool
	agent              http_agent.IHttpAgent
	clientConfig       constant.ClientConfig
	serverConfigs      []constant.ServerConfig
}

//SetClientConfig is use to set nacos client Config
func (client *NacosClient) SetClientConfig(config constant.ClientConfig) (err error) {
	if config.TimeoutMs <= 0 {
		err = errors.New("[client.SetClientConfig] config.TimeoutMs should > 0")
		return
	}

	if config.BeatInterval <= 0 {
		config.BeatInterval = 5 * 1000
	}

	if config.UpdateThreadNum <= 0 {
		config.UpdateThreadNum = 20
	}

	if len(config.RotateTime) == 0 {
		config.RotateTime = "24h"
	}

	if len(config.LogLevel) == 0 {
		config.LogLevel = "info"
	}

	if config.MaxAge <= 0 {
		config.MaxAge = 3
	}

	if config.CacheDir == "" {
		config.CacheDir = file.GetCurrentPath() + string(os.PathSeparator) + "cache"
	}

	if config.LogDir == "" {
		config.LogDir = file.GetCurrentPath() + string(os.PathSeparator) + "log"
	}

	log.Printf("[INFO] logDir:<%s>   cacheDir:<%s>", config.LogDir, config.CacheDir)
	client.clientConfig = config
	client.clientConfigValid = true

	return
}

//SetServerConfig is use to set nacos server config
func (client *NacosClient) SetServerConfig(configs []constant.ServerConfig) (err error) {
	if len(configs) <= 0 {
		//it's may be use endpoint to get nacos server address
		client.serverConfigsValid = true
		return
	}

	for i := 0; i < len(configs); i++ {
		if len(configs[i].IpAddr) <= 0 || configs[i].Port <= 0 || configs[i].Port > 65535 {
			err = errors.New("[client.SetServerConfig] configs[" + strconv.Itoa(i) + "] is invalid")
			return
		}
		if len(configs[i].ContextPath) <= 0 {
			configs[i].ContextPath = constant.DEFAULT_CONTEXT_PATH
		}
	}
	client.serverConfigs = configs
	client.serverConfigsValid = true
	return
}

//GetClientConfig use to get client config
func (client *NacosClient) GetClientConfig() (config constant.ClientConfig, err error) {
	config = client.clientConfig
	if !client.clientConfigValid {
		err = errors.New("[client.GetClientConfig] invalid client config")
	}
	return
}

//GetServerConfig use to get server config
func (client *NacosClient) GetServerConfig() (configs []constant.ServerConfig, err error) {
	configs = client.serverConfigs
	if !client.serverConfigsValid {
		err = errors.New("[client.GetServerConfig] invalid server configs")
	}
	return
}

//SetHttpAgent use to set http agent
func (client *NacosClient) SetHttpAgent(agent http_agent.IHttpAgent) (err error) {
	if agent == nil {
		err = errors.New("[client.SetHttpAgent] http agent can not be nil")
	} else {
		client.agent = agent
	}
	return
}

//GetHttpAgent use to get http agent
func (client *NacosClient) GetHttpAgent() (agent http_agent.IHttpAgent, err error) {
	if client.agent == nil {
		err = errors.New("[client.GetHttpAgent] invalid http agent")
	} else {
		agent = client.agent
	}
	return
}
