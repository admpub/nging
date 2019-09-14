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

type Results []*Result

func (r Results) FileURLs() (rs []string) {
	rs = make([]string, len(r))
	for k, v := range r {
		rs[k] = v.FileURL
	}
	return rs
}

func (r *Results) Add(result *Result) {
	*r = append(*r, result)
}

type Result struct {
	FileID            int64
	FileName          string
	FileURL           string
	FileType          FileType
	FileSize          int64
	SavePath          string
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
	err    error
}

func (a *BaseClient) Init(ctx echo.Context, res *Result) {
	a.Context = ctx
	a.Data = res
}

func (a *BaseClient) Name() string {
	return "filedata"
}

func (a *BaseClient) SetError(err error) Client {
	a.err = err
	return a
}

func (a *BaseClient) GetError() error {
	return a.err
}

func (a *BaseClient) Error() string {
	if a.err != nil {
		return a.err.Error()
	}
	return ``
}

func (a *BaseClient) Body() (file ReadCloserWithSize, err error) {
	file, a.Data.FileName, err = Receive(a.Name(), a.Context)
	if err != nil {
		return
	}
	a.Data.FileSize = file.Size()
	return
}

func (a *BaseClient) Result() (r string) {
	status := "1"
	if a.err != nil {
		status = "0"
	}
	r = `{"Code":` + status + `,"Info":"` + a.err.Error() + `","Data":{"Url":"` + a.Data.FileURL + `","Id":"` + a.Data.FileIdString() + `"}}`
	return
}

func (a *BaseClient) Response() error {
	var result string
	if a.Object != nil {
		result = a.Object.Result()
	} else {
		result = a.Result()
	}
	return a.JSONBlob(engine.Str2bytes(result))
}

type Client interface {
	//初始化
	Init(echo.Context, *Result)
	SetError(err error) Client
	GetError() error
	Error() string

	//file表单域name属性值
	Name() string

	//文件内容
	Body() (ReadCloserWithSize, error)

	//返回结果
	Result() string

	Response() error
}
