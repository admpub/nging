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
	"fmt"
	"path/filepath"
	"time"

	"github.com/webx-top/echo"
	"github.com/webx-top/echo/engine"
)

type Result struct {
	FileID            int64
	FileName          string
	FileURL           string
	FileType          FileType
	Addon             interface{}
	distFileGenerator func(string) (string, error)
}

func (r *Result) SetDistFileGenerator(generator func(string) (string, error)) *Result {
	r.distFileGenerator = generator
	return r
}

func (r *Result) DistFile() (string, error) {
	if r.distFileGenerator == nil {
		return filepath.Join(time.Now().Format("2006/0102"), r.FileName), nil
	}
	return r.distFileGenerator(r.FileName)
}

func (r *Result) FileIdString() string {
	return fmt.Sprintf(`%d`, r.FileID)
}

func New(object Client) *BaseClient {
	return &BaseClient{Object: object}
}

type BaseClient struct {
	Data *Result
	echo.Context
	Object Client
}

func (a *BaseClient) Init(ctx echo.Context, res *Result) {
	a.Context = ctx
	a.Data = res
}

func (a *BaseClient) Name() string {
	return "filedata"
}

func (a *BaseClient) Body() (file ReadCloserWithSize, err error) {
	file, a.Data.FileName, err = Receive(a.Name(), a.Context)
	if err != nil {
		return
	}
	return
}

func (a *BaseClient) Result(errMsg string) (r string) {
	status := "1"
	if len(errMsg) > 0 {
		status = "0"
	}
	r = `{"Code":` + status + `,"Info":"` + errMsg + `","Data":{"Url":"` + a.Data.FileURL + `","Id":"` + a.Data.FileIdString() + `"}}`
	return
}

func (a *BaseClient) Response(errMsg string) error {
	var result string
	if a.Object != nil {
		result = a.Object.Result(errMsg)
	} else {
		result = a.Result(errMsg)
	}
	return a.JSONBlob(engine.Str2bytes(result))
}

type Client interface {
	//初始化
	Init(echo.Context, *Result)

	//file表单域name属性值
	Name() string

	//文件内容
	Body() (ReadCloserWithSize, error)

	//返回结果
	Result(string) string

	Response(string) error
}
