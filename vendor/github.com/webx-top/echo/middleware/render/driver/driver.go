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
	"bytes"
	"io"
	"regexp"
	"strconv"

	"github.com/webx-top/echo"
	"github.com/webx-top/echo/logger"
)

type Driver interface {
	//初始化模板引擎
	Init(...bool)

	//获取模板根路径
	TmplDir() string
	Debug() bool
	SetDebug(bool)
	SetLogger(logger.Logger)
	Logger() logger.Logger

	//设置模板内容预处理器
	SetContentProcessor(fn func([]byte) []byte)
	SetManager(Manager)

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

type NopRenderer struct{}

func (n *NopRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return nil
}

func (n *NopRenderer) Debug() bool { return false }

func (n *NopRenderer) SetDebug(_ bool) {}

func (n *NopRenderer) Init(_ ...bool) {}

func (n *NopRenderer) TmplDir() string { return `` }

func (n *NopRenderer) SetLogger(_ logger.Logger) {}

func (n *NopRenderer) Logger() logger.Logger { return nil }

func (n *NopRenderer) SetContentProcessor(fn func([]byte) []byte) {}

func (n *NopRenderer) SetManager(_ Manager) {}

func (n *NopRenderer) SetFuncMap(_ func() map[string]interface{}) {}

func (n *NopRenderer) Fetch(_ string, _ interface{}, _ map[string]interface{}) string { return `` }

func (n *NopRenderer) RawContent(_ string) ([]byte, error) { return nil, nil }

func (n *NopRenderer) MonitorEvent(_ func(string)) {}

func (n *NopRenderer) ClearCache() {}

func (n *NopRenderer) Close() {}

var (
	FE       = []byte(`$1 $2`)
	First    = []byte(`$1`)
	preRegex = regexp.MustCompile(`(?is)<pre( [^>]*)?>.*?<\/pre>`)
	eolRegex = regexp.MustCompile("(?s)(\r?\n){2,}")
)

func ReplacePRE(b []byte) ([]byte, [][]byte) {
	var pres [][]byte
	b = preRegex.ReplaceAllFunc(b, func(r []byte) []byte {
		index := strconv.Itoa(len(pres))
		pres = append(pres, r)
		return []byte(`<!-- <[#pre:` + index + `#]> -->`)
	})
	return b, pres
}

func RemoveMultiCRLF(b []byte) []byte {
	return eolRegex.ReplaceAll(b, First)
}

func RecoveryPRE(b []byte, pres [][]byte) []byte {
	for k, v := range pres {
		b = bytes.Replace(b, []byte(`<!-- <[#pre:`+strconv.Itoa(k)+`#]> -->`), v, 1)
	}
	return b
}
