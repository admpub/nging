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

package security

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"sync/atomic"
	"time"

	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/common/http_agent"
	"github.com/nacos-group/nacos-sdk-go/common/logger"
)

type AuthClient struct {
	username           string
	password           string
	accessToken        *atomic.Value
	tokenTtl           int64
	lastRefreshTime    int64
	tokenRefreshWindow int64
	agent              http_agent.IHttpAgent
	clientCfg          constant.ClientConfig
	serverCfgs         []constant.ServerConfig
}

func NewAuthClient(clientCfg constant.ClientConfig, serverCfgs []constant.ServerConfig, agent http_agent.IHttpAgent) AuthClient {
	client := AuthClient{
		username:    clientCfg.Username,
		password:    clientCfg.Password,
		serverCfgs:  serverCfgs,
		clientCfg:   clientCfg,
		agent:       agent,
		accessToken: &atomic.Value{},
	}

	return client
}

func (ac *AuthClient) GetAccessToken() string {
	v := ac.accessToken.Load()
	if v == nil {
		return ""
	}
	return v.(string)
}

func (ac *AuthClient) AutoRefresh() {

	// If the username is not set, the automatic refresh Token is not enabled

	if ac.username == "" {
		return
	}

	go func() {
		timer := time.NewTimer(time.Second * time.Duration(ac.tokenTtl-ac.tokenRefreshWindow))

		for {
			select {
			case <-timer.C:
				_, err := ac.Login()
				if err != nil {
					logger.Errorf("login has error %+v", err)
				}
				timer.Reset(time.Second * time.Duration(ac.tokenTtl-ac.tokenRefreshWindow))
			}
		}
	}()
}

func (ac *AuthClient) Login() (bool, error) {
	var throwable error = nil
	for i := 0; i < len(ac.serverCfgs); i++ {
		result, err := ac.login(ac.serverCfgs[i])
		throwable = err
		if result {
			return true, nil
		}
	}
	return false, throwable
}

func (ac *AuthClient) login(server constant.ServerConfig) (bool, error) {
	if ac.username != "" {
		contextPath := server.ContextPath

		if !strings.HasPrefix(contextPath, "/") {
			contextPath = "/" + contextPath
		}

		if strings.HasSuffix(contextPath, "/") {
			contextPath = contextPath[0 : len(contextPath)-1]
		}

		if server.Scheme == "" {
			server.Scheme = "http"
		}

		reqUrl := server.Scheme + "://" + server.IpAddr + ":" + strconv.FormatInt(int64(server.Port), 10) + contextPath + "/v1/auth/users/login"

		header := http.Header{
			"content-type": []string{"application/x-www-form-urlencoded"},
		}
		resp, err := ac.agent.Post(reqUrl, header, ac.clientCfg.TimeoutMs, map[string]string{
			"username": ac.username,
			"password": ac.password,
		})

		if err != nil {
			return false, err
		}

		var bytes []byte
		bytes, err = ioutil.ReadAll(resp.Body)
		defer resp.Body.Close()
		if err != nil {
			return false, err
		}

		if resp.StatusCode != 200 {
			errMsg := string(bytes)
			return false, errors.New(errMsg)
		}

		var result map[string]interface{}

		err = json.Unmarshal(bytes, &result)

		if err != nil {
			return false, err
		}

		if val, ok := result[constant.KEY_ACCESS_TOKEN]; ok {
			ac.accessToken.Store(val)
			ac.tokenTtl = int64(result[constant.KEY_TOKEN_TTL].(float64))
			ac.tokenRefreshWindow = ac.tokenTtl / 10
		}
	}
	return true, nil

}
