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
	"github.com/nacos-group/nacos-sdk-go/model"
	"github.com/nacos-group/nacos-sdk-go/vo"
)

//go:generate mockgen -destination ../../mock/mock_service_client_interface.go -package mock -source=./service_client_interface.go

type INamingClient interface {

	//RegisterInstance use to register instance
	//Ip  require
	//Port  require
	//Tenant optional
	//Weight  require,it must be lager than 0
	//Enable  require,the instance can be access or not
	//Healthy  require,the instance is health or not
	//Metadata  optional
	//ClusterName  optional,default:DEFAULT
	//ServiceName require
	//GroupName optional,default:DEFAULT_GROUP
	//Ephemeral optional
	RegisterInstance(param vo.RegisterInstanceParam) (bool, error)

	//DeregisterInstance use to deregister instance
	//Ip required
	//Port required
	//Tenant optional
	//Cluster optional,default:DEFAULT
	//ServiceName  require
	//GroupName  optional,default:DEFAULT_GROUP
	//Ephemeral optional
	DeregisterInstance(param vo.DeregisterInstanceParam) (bool, error)

	//GetService use to get service
	//ServiceName require
	//Clusters optional,default:DEFAULT
	//GroupName optional,default:DEFAULT_GROUP
	GetService(param vo.GetServiceParam) (model.Service, error)

	//SelectAllInstance return all instances,include healthy=false,enable=false,weight<=0
	//ServiceName require
	//Clusters optional,default:DEFAULT
	//GroupName optional,default:DEFAULT_GROUP
	SelectAllInstances(param vo.SelectAllInstancesParam) ([]model.Instance, error)

	//SelectInstances only return the instances of healthy=${HealthyOnly},enable=true and weight>0
	//ServiceName require
	//Clusters optional,default:DEFAULT
	//GroupName optional,default:DEFAULT_GROUP
	//HealthyOnly optional
	SelectInstances(param vo.SelectInstancesParam) ([]model.Instance, error)

	//SelectInstances return one instance by WRR strategy for load balance
	//And the instance should be health=true,enable=true and weight>0
	//ServiceName require
	//Clusters optional,default:DEFAULT
	//GroupName optional,default:DEFAULT_GROUP
	SelectOneHealthyInstance(param vo.SelectOneHealthInstanceParam) (*model.Instance, error)

	//Subscribe use to subscribe service change event
	//ServiceName require
	//Clusters optional,default:DEFAULT
	//GroupName optional,default:DEFAULT_GROUP
	//SubscribeCallback require
	Subscribe(param *vo.SubscribeParam) error

	//Unsubscribe use to unsubscribe service change event
	//ServiceName require
	//Clusters optional,default:DEFAULT
	//GroupName optional,default:DEFAULT_GROUP
	//SubscribeCallback require
	Unsubscribe(param *vo.SubscribeParam) error

	//GetAllServicesInfo use to get all service info by page
	GetAllServicesInfo(param vo.GetAllServiceInfoParam) (model.ServiceList, error)
}
