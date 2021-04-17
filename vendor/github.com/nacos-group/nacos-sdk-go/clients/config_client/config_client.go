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
	"errors"
	"math"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/kms"
	"github.com/nacos-group/nacos-sdk-go/clients/cache"
	"github.com/nacos-group/nacos-sdk-go/clients/nacos_client"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/common/http_agent"
	"github.com/nacos-group/nacos-sdk-go/common/logger"
	"github.com/nacos-group/nacos-sdk-go/common/nacos_error"
	"github.com/nacos-group/nacos-sdk-go/model"
	"github.com/nacos-group/nacos-sdk-go/util"
	"github.com/nacos-group/nacos-sdk-go/vo"
)

type ConfigClient struct {
	nacos_client.INacosClient
	kmsClient      *kms.Client
	localConfigs   []vo.ConfigParam
	mutex          sync.Mutex
	configProxy    ConfigProxy
	configCacheDir string
}

const (
	perTaskConfigSize = 3000
	executorErrDelay  = 5 * time.Second
)

var (
	currentTaskCount int
	cacheMap         = cache.NewConcurrentMap()
	schedulerMap     = cache.NewConcurrentMap()
)

type cacheData struct {
	isInitializing    bool
	dataId            string
	group             string
	content           string
	tenant            string
	cacheDataListener *cacheDataListener
	md5               string
	appName           string
	taskId            int
	configClient      *ConfigClient
}

type cacheDataListener struct {
	listener vo.Listener
	lastMd5  string
}

func init() {
	schedulerMap.Set("root", true)
	go delayScheduler(time.NewTimer(1*time.Millisecond), 10*time.Millisecond, "root", listenConfigExecutor())
}

func NewConfigClient(nc nacos_client.INacosClient) (ConfigClient, error) {
	config := ConfigClient{}
	config.INacosClient = nc
	clientConfig, err := nc.GetClientConfig()
	if err != nil {
		return config, err
	}
	serverConfig, err := nc.GetServerConfig()
	if err != nil {
		return config, err
	}
	httpAgent, err := nc.GetHttpAgent()
	if err != nil {
		return config, err
	}
	err = logger.InitLogger(logger.Config{
		Level:        clientConfig.LogLevel,
		OutputPath:   clientConfig.LogDir,
		RotationTime: clientConfig.RotateTime,
		MaxAge:       clientConfig.MaxAge,
	})
	if err != nil {
		return config, err
	}
	config.configCacheDir = clientConfig.CacheDir + string(os.PathSeparator) + "config"
	config.configProxy, err = NewConfigProxy(serverConfig, clientConfig, httpAgent)
	if clientConfig.OpenKMS {
		kmsClient, err := kms.NewClientWithAccessKey(clientConfig.RegionId, clientConfig.AccessKey, clientConfig.SecretKey)
		if err != nil {
			return config, err
		}
		config.kmsClient = kmsClient
	}
	return config, err
}

func (client *ConfigClient) sync() (clientConfig constant.ClientConfig,
	serverConfigs []constant.ServerConfig, agent http_agent.IHttpAgent, err error) {
	clientConfig, err = client.GetClientConfig()
	if err != nil {
		logger.Errorf("getClientConfig catch error:%+v", err)
		return
	}
	serverConfigs, err = client.GetServerConfig()
	if err != nil {
		logger.Errorf("getServerConfig catch error:%+v", err)
		return
	}

	agent, err = client.GetHttpAgent()
	if err != nil {
		logger.Errorf("getHttpAgent catch error:%+v", err)
	}
	return
}

func (client *ConfigClient) GetConfig(param vo.ConfigParam) (content string, err error) {
	content, err = client.getConfigInner(param)

	if err != nil {
		return "", err
	}

	return client.decrypt(param.DataId, content)
}

func (client *ConfigClient) decrypt(dataId, content string) (string, error) {
	if strings.HasPrefix(dataId, "cipher-") && client.kmsClient != nil {
		request := kms.CreateDecryptRequest()
		request.Method = "POST"
		request.Scheme = "https"
		request.AcceptFormat = "json"
		request.CiphertextBlob = content
		response, err := client.kmsClient.Decrypt(request)
		if err != nil {
			return "", errors.New("kms decrypt failed")
		}
		content = response.Plaintext
	}

	return content, nil
}

func (client *ConfigClient) getConfigInner(param vo.ConfigParam) (content string, err error) {
	if len(param.DataId) <= 0 {
		err = errors.New("[client.GetConfig] param.dataId can not be empty")
		return "", err
	}
	if len(param.Group) <= 0 {
		err = errors.New("[client.GetConfig] param.group can not be empty")
		return "", err
	}
	clientConfig, _ := client.GetClientConfig()
	cacheKey := util.GetConfigCacheKey(param.DataId, param.Group, clientConfig.NamespaceId)
	content, err = client.configProxy.GetConfigProxy(param, clientConfig.NamespaceId, clientConfig.AccessKey, clientConfig.SecretKey)

	if err != nil {
		logger.Infof("get config from server error:%+v ", err)
		if _, ok := err.(*nacos_error.NacosError); ok {
			nacosErr := err.(*nacos_error.NacosError)
			if nacosErr.ErrorCode() == "404" {
				cache.WriteConfigToFile(cacheKey, client.configCacheDir, "")
				return "", errors.New("config not found")
			}
			if nacosErr.ErrorCode() == "403" {
				return "", errors.New("get config forbidden")
			}
		}
		content, err = cache.ReadConfigFromFile(cacheKey, client.configCacheDir)
		if err != nil {
			logger.Errorf("get config from cache  error:%+v ", err)
			return "", errors.New("read config from both server and cache fail")
		}

	} else {
		cache.WriteConfigToFile(cacheKey, client.configCacheDir, content)
	}
	return content, nil
}

func (client *ConfigClient) PublishConfig(param vo.ConfigParam) (published bool,
	err error) {
	if len(param.DataId) <= 0 {
		err = errors.New("[client.PublishConfig] param.dataId can not be empty")
	}
	if len(param.Group) <= 0 {
		err = errors.New("[client.PublishConfig] param.group can not be empty")
	}
	if len(param.Content) <= 0 {
		err = errors.New("[client.PublishConfig] param.content can not be empty")
	}
	clientConfig, _ := client.GetClientConfig()
	return client.configProxy.PublishConfigProxy(param, clientConfig.NamespaceId, clientConfig.AccessKey, clientConfig.SecretKey)
}

func (client *ConfigClient) DeleteConfig(param vo.ConfigParam) (deleted bool, err error) {
	if len(param.DataId) <= 0 {
		err = errors.New("[client.DeleteConfig] param.dataId can not be empty")
	}
	if len(param.Group) <= 0 {
		err = errors.New("[client.DeleteConfig] param.group can not be empty")
	}

	clientConfig, _ := client.GetClientConfig()
	return client.configProxy.DeleteConfigProxy(param, clientConfig.NamespaceId, clientConfig.AccessKey, clientConfig.SecretKey)
}

//Cancel Listen Config
func (client *ConfigClient) CancelListenConfig(param vo.ConfigParam) (err error) {
	clientConfig, err := client.GetClientConfig()
	if err != nil {
		logger.Errorf("[checkConfigInfo.GetClientConfig] failed,err:%+v", err)
		return
	}
	cacheMap.Remove(util.GetConfigCacheKey(param.DataId, param.Group, clientConfig.NamespaceId))
	logger.Infof("Cancel listen config DataId:%s Group:%s", param.DataId, param.Group)
	remakeId := int(math.Ceil(float64(len(cacheMap.Keys())) / float64(perTaskConfigSize)))
	if remakeId < currentTaskCount {
		remakeCacheDataTaskId(remakeId)
	}
	return err
}

//Remake cache data taskId
func remakeCacheDataTaskId(remakeId int) {
	for i := 0; i < remakeId; i++ {
		count := 0
		for _, key := range cacheMap.Keys() {
			if count == perTaskConfigSize {
				break
			}
			if value, ok := cacheMap.Get(key); ok {
				cData := value.(cacheData)
				cData.taskId = i
				cacheMap.Set(key, cData)
			}
			count++
		}
	}
}

func (client *ConfigClient) ListenConfig(param vo.ConfigParam) (err error) {
	if len(param.DataId) <= 0 {
		err = errors.New("[client.ListenConfig] DataId can not be empty")
		return err
	}
	if len(param.Group) <= 0 {
		err = errors.New("[client.ListenConfig] Group can not be empty")
		return err
	}
	clientConfig, err := client.GetClientConfig()
	if err != nil {
		err = errors.New("[checkConfigInfo.GetClientConfig] failed")
		return err
	}

	key := util.GetConfigCacheKey(param.DataId, param.Group, clientConfig.NamespaceId)
	var cData cacheData
	if v, ok := cacheMap.Get(key); ok {
		cData = v.(cacheData)
		cData.isInitializing = true
	} else {
		var (
			content string
			md5Str  string
		)
		content, fileErr := cache.ReadConfigFromFile(key, client.configCacheDir)
		if fileErr != nil {
			logger.Errorf("[cache.ReadConfigFromFile] error: %+v", err)
		}
		if len(content) > 0 {
			md5Str = util.Md5(content)
		}
		listener := &cacheDataListener{
			listener: param.OnChange,
			lastMd5:  md5Str,
		}

		cData = cacheData{
			isInitializing:    true,
			dataId:            param.DataId,
			group:             param.Group,
			tenant:            clientConfig.NamespaceId,
			content:           content,
			md5:               md5Str,
			cacheDataListener: listener,
			taskId:            len(cacheMap.Keys()) / perTaskConfigSize,
			configClient:      client,
		}
	}
	cacheMap.Set(key, cData)
	return
}

//Delay Scheduler
//initialDelay the time to delay first execution
//delay the delay between the termination of one execution and the commencement of the next
func delayScheduler(t *time.Timer, delay time.Duration, taskId string, execute func() error) {
	for {
		if v, ok := schedulerMap.Get(taskId); ok {
			if !v.(bool) {
				return
			}
		}
		<-t.C
		d := delay
		if err := execute(); err != nil {
			d = executorErrDelay
		}
		t.Reset(d)
	}
}

//Listen for the configuration executor
func listenConfigExecutor() func() error {
	return func() error {
		listenerSize := len(cacheMap.Keys())
		taskCount := int(math.Ceil(float64(listenerSize) / float64(perTaskConfigSize)))

		if taskCount > currentTaskCount {
			for i := currentTaskCount; i < taskCount; i++ {
				schedulerMap.Set(strconv.Itoa(i), true)
				go delayScheduler(time.NewTimer(1*time.Millisecond), 10*time.Millisecond, strconv.Itoa(i), longPulling(i))
			}
			currentTaskCount = taskCount
		} else if taskCount < currentTaskCount {
			for i := taskCount; i < currentTaskCount; i++ {
				if _, ok := schedulerMap.Get(strconv.Itoa(i)); ok {
					schedulerMap.Set(strconv.Itoa(i), false)
				}
			}
			currentTaskCount = taskCount
		}
		return nil
	}
}

//Long polling listening configuration
func longPulling(taskId int) func() error {
	return func() error {
		var listeningConfigs string
		var client *ConfigClient
		initializationList := make([]cacheData, 0)
		for _, key := range cacheMap.Keys() {
			if value, ok := cacheMap.Get(key); ok {
				cData := value.(cacheData)
				client = cData.configClient
				if cData.taskId == taskId {
					if cData.isInitializing {
						initializationList = append(initializationList, cData)
					}
					if len(cData.tenant) > 0 {
						listeningConfigs += cData.dataId + constant.SPLIT_CONFIG_INNER + cData.group + constant.SPLIT_CONFIG_INNER +
							cData.md5 + constant.SPLIT_CONFIG_INNER + cData.tenant + constant.SPLIT_CONFIG
					} else {
						listeningConfigs += cData.dataId + constant.SPLIT_CONFIG_INNER + cData.group + constant.SPLIT_CONFIG_INNER +
							cData.md5 + constant.SPLIT_CONFIG
					}
				}
			}
		}
		if len(listeningConfigs) > 0 {
			clientConfig, err := client.GetClientConfig()
			if err != nil {
				logger.Errorf("[checkConfigInfo.GetClientConfig] err: %+v", err)
				return err
			}
			// http get
			params := make(map[string]string)
			params[constant.KEY_LISTEN_CONFIGS] = listeningConfigs

			var changed string
			changedTmp, err := client.configProxy.ListenConfig(params, len(initializationList) > 0, clientConfig.NamespaceId, clientConfig.AccessKey, clientConfig.SecretKey)
			if err == nil {
				changed = changedTmp
			} else {
				if _, ok := err.(*nacos_error.NacosError); ok {
					changed = changedTmp
				} else {
					logger.Errorf("[client.ListenConfig] listen config error: %+v", err)
				}
				return err
			}
			for _, v := range initializationList {
				v.isInitializing = false
				cacheMap.Set(util.GetConfigCacheKey(v.dataId, v.group, clientConfig.NamespaceId), v)
			}
			if len(strings.ToLower(strings.Trim(changed, " "))) == 0 {
				logger.Info("[client.ListenConfig] no change")
			} else {
				logger.Info("[client.ListenConfig] config changed:" + changed)
				client.callListener(changed, clientConfig.NamespaceId)
			}
		}
		return nil
	}

}

//Execute the Listener callback func()
func (client *ConfigClient) callListener(changed, tenant string) {
	changedConfigs := strings.Split(changed, "%01")
	for _, config := range changedConfigs {
		attrs := strings.Split(config, "%02")
		if len(attrs) >= 2 {
			if value, ok := cacheMap.Get(util.GetConfigCacheKey(attrs[0], attrs[1], tenant)); ok {
				cData := value.(cacheData)
				content, err := client.getConfigInner(vo.ConfigParam{
					DataId: cData.dataId,
					Group:  cData.group,
				})
				if err != nil {
					logger.Errorf("[client.getConfigInner] DataId:[%s] Group:[%s] Error:[%+v]", cData.dataId, cData.group, err)
					continue
				}
				cData.content = content
				cData.md5 = util.Md5(content)
				if cData.md5 != cData.cacheDataListener.lastMd5 {
					go cData.cacheDataListener.listener(tenant, attrs[1], attrs[0], cData.content)
					cData.cacheDataListener.lastMd5 = cData.md5
					cacheMap.Set(util.GetConfigCacheKey(cData.dataId, cData.group, tenant), cData)
				}
			}
		}
	}
}

func (client *ConfigClient) buildBasePath(serverConfig constant.ServerConfig) (basePath string) {
	basePath = "http://" + serverConfig.IpAddr + ":" +
		strconv.FormatUint(serverConfig.Port, 10) + serverConfig.ContextPath + constant.CONFIG_PATH
	return
}

func (client *ConfigClient) SearchConfig(param vo.SearchConfigParm) (*model.ConfigPage, error) {
	return client.searchConfigInner(param)
}

func (client *ConfigClient) searchConfigInner(param vo.SearchConfigParm) (*model.ConfigPage, error) {
	if param.Search != "accurate" && param.Search != "blur" {
		return nil, errors.New("[client.searchConfigInner] param.search must be accurate or blur")
	}
	if param.PageNo <= 0 {
		param.PageNo = 1
	}
	if param.PageSize <= 0 {
		param.PageSize = 10
	}
	clientConfig, _ := client.GetClientConfig()
	configItems, err := client.configProxy.SearchConfigProxy(param, clientConfig.NamespaceId, clientConfig.AccessKey, clientConfig.SecretKey)
	if err != nil {
		logger.Errorf("search config from server error:%+v ", err)
		if _, ok := err.(*nacos_error.NacosError); ok {
			nacosErr := err.(*nacos_error.NacosError)
			if nacosErr.ErrorCode() == "404" {
				return nil, errors.New("config not found")
			}
			if nacosErr.ErrorCode() == "403" {
				return nil, errors.New("get config forbidden")
			}
		}
		return nil, err
	}
	return configItems, nil
}
