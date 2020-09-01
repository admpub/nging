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
	"net/http"
	"path/filepath"
	"time"

	"github.com/webx-top/echo"
)

type Results []*Result
type Checker func(r *Result) error

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
	Md5               string
	Addon             interface{}
	fileNameGenerator func(string) (string, error)
}

func (r *Result) SetFileNameGenerator(generator func(string) (string, error)) *Result {
	r.fileNameGenerator = generator
	return r
}

func (r *Result) GenFileName() (string, error) {
	if r.fileNameGenerator == nil {
		return filepath.Join(time.Now().Format("2006/0102"), r.FileName), nil
	}
	return r.fileNameGenerator(r.FileName)
}

func (r *Result) FileIdString() string {
	return fmt.Sprintf(`%d`, r.FileID)
}

func New(object Client, formFields ...string) *BaseClient {
	formField := DefaultFormField
	if len(formFields) > 0 {
		formField = formFields[0]
	}
	return &BaseClient{Object: object, FormField: formField}
}

var DefaultFormField = `filedata`

type BaseClient struct {
	Data *Result
	echo.Context
	Object       Client
	FormField    string // 表单文件字段名
	Code         int    // HTTP code
	ContentType  string
	JSONPVarName string
	err          error
	RespData     interface{}
}

func (a *BaseClient) Init(ctx echo.Context, res *Result) {
	a.Context = ctx
	a.Data = res
}

func (a *BaseClient) Name() string {
	if len(a.FormField) == 0 {
		return DefaultFormField
	}
	return a.FormField
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
	a.Data.Md5, err = file.Md5()
	return
}

func (a *BaseClient) BuildResult() {
	data := a.Context.Data()
	data.SetData(echo.H{
		`Url`: a.Data.FileURL,
		`Id`:  a.Data.FileIdString(),
	}, 1)
	if a.err != nil {
		data.SetError(a.err)
	}
	a.RespData = data
	return
}

func (a *BaseClient) Response() error {
	if a.Object != nil {
		a.Object.BuildResult()
	} else {
		a.BuildResult()
	}
	if a.Code > 0 {
		return a.responseContentType(a.Code)
	}
	return a.responseContentType(http.StatusOK)
}

func (a *BaseClient) responseContentType(code int) error {
	switch a.ContentType {
	case `string`:
		return a.String(fmt.Sprint(a.RespData), code)
	case `xml`:
		return a.XML(a.RespData, code)
	case `redirect`:
		return a.Redirect(fmt.Sprint(a.RespData), code)
	default:
		if len(a.JSONPVarName) > 0 {
			return a.JSONP(a.JSONPVarName, a.RespData, code)
		}
		return a.JSON(a.RespData, code)
	}
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

	//构建结果
	BuildResult()

	Response() error
}
