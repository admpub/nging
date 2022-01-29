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

package upload

import (
	"github.com/webx-top/echo"
)

// Client 上次客户端处理接口
type Client interface {
	//初始化
	Init(echo.Context, *Result)
	SetUploadMaxSize(maxSize int64) Client
	SetReadBeforeHook(hooks ...ReadBeforeHook) Client
	AddReadBeforeHook(hooks ...ReadBeforeHook) Client

	//分片上传支持
	SetChunkUpload(cu *ChunkUpload) Client
	SetFieldMapping(fm map[string]string) Client

	UploadMaxSize() int64
	SetError(err error) Client
	GetError() error
	ErrorString() string

	//file表单域name属性值
	Name() string
	SetName(formField string)

	//文件内容
	Body() (ReadCloserWithSize, error)
	Upload(...OptionsSetter) Client
	BatchUpload(...OptionsSetter) Client
	GetUploadResult() *Result
	GetBatchUploadResults() Results

	//构建结果
	BuildResult() Client

	GetRespData() interface{}
	SetRespData(data interface{}) Client

	Response() error
	Reset()
}
