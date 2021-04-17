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

package cache

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"

	"github.com/go-errors/errors"
	"github.com/nacos-group/nacos-sdk-go/common/file"
	"github.com/nacos-group/nacos-sdk-go/common/logger"
	"github.com/nacos-group/nacos-sdk-go/model"
	"github.com/nacos-group/nacos-sdk-go/util"
)

func GetFileName(cacheKey string, cacheDir string) string {
	return cacheDir + string(os.PathSeparator) + cacheKey
}

func WriteServicesToFile(service model.Service, cacheDir string) {
	file.MkdirIfNecessary(cacheDir)
	sb, _ := json.Marshal(service)
	domFileName := GetFileName(util.GetServiceCacheKey(service.Name, service.Clusters), cacheDir)

	err := ioutil.WriteFile(domFileName, sb, 0666)
	if err != nil {
		logger.Errorf("failed to write name cache:%s ,value:%s ,err:%+v", domFileName, string(sb), err)
	}

}

func ReadServicesFromFile(cacheDir string) map[string]model.Service {
	files, err := ioutil.ReadDir(cacheDir)
	if err != nil {
		logger.Errorf("read cacheDir:%s failed!err:%+v", cacheDir, err)
		return nil
	}
	serviceMap := map[string]model.Service{}
	for _, f := range files {
		fileName := GetFileName(f.Name(), cacheDir)
		b, err := ioutil.ReadFile(fileName)
		if err != nil {
			logger.Errorf("failed to read name cache file:%s,err:%+v ", fileName, err)
			continue
		}

		s := string(b)
		service := util.JsonToService(s)

		if service == nil {
			continue
		}

		serviceMap[f.Name()] = *service
	}

	logger.Info("finish loading name cache, total: " + strconv.Itoa(len(files)))
	return serviceMap
}

func WriteConfigToFile(cacheKey string, cacheDir string, content string) {
	file.MkdirIfNecessary(cacheDir)
	fileName := GetFileName(cacheKey, cacheDir)
	err := ioutil.WriteFile(fileName, []byte(content), 0666)
	if err != nil {
		logger.Errorf("failed to write config  cache:%s ,value:%s ,err:%+v", fileName, content, err)
	}
}

func ReadConfigFromFile(cacheKey string, cacheDir string) (string, error) {
	fileName := GetFileName(cacheKey, cacheDir)
	b, err := ioutil.ReadFile(fileName)
	if err != nil {
		return "", errors.New(fmt.Sprintf("failed to read config cache file:%s,err:%+v ", fileName, err))
	}
	return string(b), nil
}
