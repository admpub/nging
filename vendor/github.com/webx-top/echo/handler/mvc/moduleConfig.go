/*

   Copyright 2016 Wenhui Shen <www.webx.top>

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.

*/
package mvc

import (
	"encoding/json"

	"github.com/admpub/log"
	"github.com/webx-top/echo/engine"
)

type ModuleConfiger interface {
	Init(conf interface{}) ModuleConfiger
	Set(data string) error
	Get(recvFn func(interface{}))
	Config() interface{}
	IsValid() bool
	String() string
	Template() string                           //获取表单模板名称
	SetTemplate(tmplFile string) ModuleConfiger //获取设置表单模板
}

type ModuleConfig struct {
	config   interface{}
	tmplFile string
}

func (a *ModuleConfig) Init(config interface{}) ModuleConfiger {
	a.config = config
	return a
}

/*
Set usage:
var appConf map[string]string
a.Set(`{"Name":"webx"}`)
*/
func (a *ModuleConfig) Set(data string) error {
	if len(data) == 0 || a.config == nil {
		return nil
	}
	err := json.Unmarshal(engine.Str2bytes(data), a.config)
	return err
}

/*
Get usage:
var appConf *map[string]string
a.Get(func(conf interface{}){
	if v, y := conf.(*map[string]string); y {
		appConf = v
	}
})
*/
func (a *ModuleConfig) Get(recvFn func(interface{})) {
	if recvFn != nil {
		recvFn(a.config)
	}
}

func (a *ModuleConfig) IsValid() bool {
	return a.config != nil
}

func (a *ModuleConfig) Config() interface{} {
	return a.config
}

func (a *ModuleConfig) String() string {
	if a.config == nil {
		return ``
	}
	b, err := json.Marshal(a.config)
	if err != nil {
		log.Error(err)
		return ``
	}
	return engine.Bytes2str(b)
}

func (a *ModuleConfig) Template() string {
	return a.tmplFile
}

func (a *ModuleConfig) SetTemplate(tmplFile string) ModuleConfiger {
	a.tmplFile = tmplFile
	return a
}
