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
package driver

import (
	"io"

	"github.com/webx-top/echo"
	"github.com/webx-top/echo/logger"
)

type Driver interface {
	//初始化模板引擎
	Init(...bool)

	//获取模板根路径
	TmplDir() string
	SetLogger(logger.Logger)
	Logger() logger.Logger

	//设置模板内容预处理器
	SetContentProcessor(fn func([]byte) []byte)

	//设置模板函数
	SetFuncMap(func() map[string]interface{})

	//渲染模板
	Render(io.Writer, string, interface{}, echo.Context) error

	//获取模板渲染后的结果
	Fetch(string, interface{}, map[string]interface{}) string

	//读取模板原始内容
	RawContent(string) ([]byte, error)

	//模板目录监控事件
	MonitorEvent(func(string))

	//清除模板对象缓存
	ClearCache()

	//关闭并停用模板引擎
	Close()
}
