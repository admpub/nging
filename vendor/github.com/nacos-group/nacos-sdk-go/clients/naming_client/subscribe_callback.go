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

package naming_client

import (
	"errors"

	"github.com/nacos-group/nacos-sdk-go/clients/cache"
	"github.com/nacos-group/nacos-sdk-go/common/logger"
	"github.com/nacos-group/nacos-sdk-go/model"
	"github.com/nacos-group/nacos-sdk-go/util"
)

type SubscribeCallback struct {
	callbackFuncsMap cache.ConcurrentMap
}

func NewSubscribeCallback() SubscribeCallback {
	ed := SubscribeCallback{}
	ed.callbackFuncsMap = cache.NewConcurrentMap()
	return ed
}

func (ed *SubscribeCallback) AddCallbackFuncs(serviceName string, clusters string, callbackFunc *func(services []model.SubscribeService, err error)) {
	logger.Info("adding " + serviceName + " with " + clusters + " to listener map")
	key := util.GetServiceCacheKey(serviceName, clusters)
	var funcs []*func(services []model.SubscribeService, err error)
	old, ok := ed.callbackFuncsMap.Get(key)
	if ok {
		funcs = append(funcs, old.([]*func(services []model.SubscribeService, err error))...)
	}
	funcs = append(funcs, callbackFunc)
	ed.callbackFuncsMap.Set(key, funcs)
}

func (ed *SubscribeCallback) RemoveCallbackFuncs(serviceName string, clusters string, callbackFunc *func(services []model.SubscribeService, err error)) {
	logger.Info("removing " + serviceName + " with " + clusters + " to listener map")
	key := util.GetServiceCacheKey(serviceName, clusters)
	funcs, ok := ed.callbackFuncsMap.Get(key)
	if ok && funcs != nil {
		var newFuncs []*func(services []model.SubscribeService, err error)
		for _, funcItem := range funcs.([]*func(services []model.SubscribeService, err error)) {
			if funcItem != callbackFunc {
				newFuncs = append(newFuncs, funcItem)
			}
		}
		ed.callbackFuncsMap.Set(key, newFuncs)
	}

}

func (ed *SubscribeCallback) ServiceChanged(service *model.Service) {
	if service == nil || service.Name == "" {
		return
	}
	key := util.GetServiceCacheKey(service.Name, service.Clusters)
	funcs, ok := ed.callbackFuncsMap.Get(key)
	if ok {
		for _, funcItem := range funcs.([]*func(services []model.SubscribeService, err error)) {
			var subscribeServices []model.SubscribeService
			if len(service.Hosts) == 0 {
				(*funcItem)(subscribeServices, errors.New("[client.Subscribe] subscribe failed,hosts is empty"))
				return
			}
			for _, host := range service.Hosts {
				var subscribeService model.SubscribeService
				subscribeService.Valid = host.Valid
				subscribeService.Port = host.Port
				subscribeService.Ip = host.Ip
				subscribeService.Metadata = host.Metadata
				subscribeService.ServiceName = host.ServiceName
				subscribeService.ClusterName = host.ClusterName
				subscribeService.Weight = host.Weight
				subscribeService.InstanceId = host.InstanceId
				subscribeService.Enable = host.Enable
				subscribeServices = append(subscribeServices, subscribeService)
			}
			(*funcItem)(subscribeServices, nil)
		}
	}
}
