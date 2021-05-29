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
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/common/http_agent"
	"github.com/nacos-group/nacos-sdk-go/common/logger"
	"github.com/nacos-group/nacos-sdk-go/common/nacos_server"
	"github.com/nacos-group/nacos-sdk-go/model"
	"github.com/nacos-group/nacos-sdk-go/util"
	"github.com/nacos-group/nacos-sdk-go/vo"
)

type ConfigProxy struct {
	nacosServer  nacos_server.NacosServer
	clientConfig constant.ClientConfig
}

func NewConfigProxy(serverConfig []constant.ServerConfig, clientConfig constant.ClientConfig, httpAgent http_agent.IHttpAgent) (ConfigProxy, error) {
	proxy := ConfigProxy{}
	var err error
	proxy.nacosServer, err = nacos_server.NewNacosServer(serverConfig, clientConfig, httpAgent, clientConfig.TimeoutMs, clientConfig.Endpoint)
	proxy.clientConfig = clientConfig
	return proxy, err

}

func (cp *ConfigProxy) GetServerList() []constant.ServerConfig {
	return cp.nacosServer.GetServerList()
}

func (cp *ConfigProxy) GetConfigProxy(param vo.ConfigParam, tenant, accessKey, secretKey string) (string, error) {
	params := util.TransformObject2Param(param)
	if len(tenant) > 0 {
		params["tenant"] = tenant
	}

	var headers = map[string]string{}
	headers["accessKey"] = accessKey
	headers["secretKey"] = secretKey

	result, err := cp.nacosServer.ReqConfigApi(constant.CONFIG_PATH, params, headers, http.MethodGet, cp.clientConfig.TimeoutMs)
	return result, err
}

func (cp *ConfigProxy) SearchConfigProxy(param vo.SearchConfigParm, tenant, accessKey, secretKey string) (*model.ConfigPage, error) {
	params := util.TransformObject2Param(param)
	if len(tenant) > 0 {
		params["tenant"] = tenant
	}
	if _, ok := params["group"]; !ok {
		params["group"] = ""
	}
	if _, ok := params["dataId"]; !ok {
		params["dataId"] = ""
	}
	var headers = map[string]string{}
	headers["accessKey"] = accessKey
	headers["secretKey"] = secretKey
	result, err := cp.nacosServer.ReqConfigApi(constant.CONFIG_PATH, params, headers, http.MethodGet, cp.clientConfig.TimeoutMs)
	if err != nil {
		return nil, err
	}
	var configPage model.ConfigPage
	err = json.Unmarshal([]byte(result), &configPage)
	if err != nil {
		return nil, err
	}
	return &configPage, nil
}
func (cp *ConfigProxy) PublishConfigProxy(param vo.ConfigParam, tenant, accessKey, secretKey string) (bool, error) {
	params := util.TransformObject2Param(param)
	if len(tenant) > 0 {
		params["tenant"] = tenant
	}

	var headers = map[string]string{}
	headers["accessKey"] = accessKey
	headers["secretKey"] = secretKey
	result, err := cp.nacosServer.ReqConfigApi(constant.CONFIG_PATH, params, headers, http.MethodPost, cp.clientConfig.TimeoutMs)
	if err != nil {
		return false, errors.New("[client.PublishConfig] publish config failed:" + err.Error())
	}
	if strings.ToLower(strings.Trim(result, " ")) == "true" {
		return true, nil
	} else {
		return false, errors.New("[client.PublishConfig] publish config failed:" + result)
	}
}

func (cp *ConfigProxy) DeleteConfigProxy(param vo.ConfigParam, tenant, accessKey, secretKey string) (bool, error) {
	params := util.TransformObject2Param(param)
	if len(tenant) > 0 {
		params["tenant"] = tenant
	}
	var headers = map[string]string{}
	headers["accessKey"] = accessKey
	headers["secretKey"] = secretKey
	result, err := cp.nacosServer.ReqConfigApi(constant.CONFIG_PATH, params, headers, http.MethodDelete, cp.clientConfig.TimeoutMs)
	if err != nil {
		return false, errors.New("[client.DeleteConfig] deleted config failed:" + err.Error())
	}
	if strings.ToLower(strings.Trim(result, " ")) == "true" {
		return true, nil
	} else {
		return false, errors.New("[client.DeleteConfig] deleted config failed: " + string(result))
	}
}

func (cp *ConfigProxy) ListenConfig(params map[string]string, isInitializing bool, tenant, accessKey, secretKey string) (string, error) {
	//fixed at 30000ms，avoid frequent request on the server
	var listenInterval uint64 = 30000
	headers := map[string]string{
		"Content-Type":         "application/x-www-form-urlencoded;charset=utf-8",
		"Long-Pulling-Timeout": strconv.FormatUint(listenInterval, 10),
	}
	if isInitializing {
		headers["Long-Pulling-Timeout-No-Hangup"] = "true"
	}

	headers["accessKey"] = accessKey
	headers["secretKey"] = secretKey

	if len(tenant) > 0 {
		params["tenant"] = tenant
	}
	logger.Infof("[client.ListenConfig] request params:%+v header:%+v \n", params, headers)
	// In order to prevent the server from handling the delay of the client's long task,
	// increase the client's read timeout to avoid this problem.
	timeout := listenInterval + listenInterval/10
	result, err := cp.nacosServer.ReqConfigApi(constant.CONFIG_LISTEN_PATH, params, headers, http.MethodPost, timeout)
	return result, err
}
