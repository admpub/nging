package upload

import (
	"fmt"
	"net/http"

	"github.com/webx-top/echo"
)

func New(object Client, formFields ...string) *BaseClient {
	formField := DefaultFormField
	if len(formFields) > 0 {
		formField = formFields[0]
	}
	return &BaseClient{
		Object:    object,
		FormField: formField,
		Results:   Results{},
	}
}

var DefaultFormField = `filedata`

type BaseClient struct {
	Data *Result
	echo.Context
	Object        Client
	FormField     string // 表单文件字段名
	Code          int    // HTTP code
	ContentType   string
	JSONPVarName  string
	RespData      interface{}
	Results       Results
	err           error
	uploadMaxSize int64 // 单位字节 (0 代表未设置，小于 0 代表不限制)
	readBefore    []ReadBeforeHook
	chunkUpload   *ChunkUpload
	fieldMapping  map[string]string
}

func (a *BaseClient) Init(ctx echo.Context, res *Result) {
	a.Context = ctx
	a.Data = res
}

func (a *BaseClient) Reset() {
	a.Data = nil
	a.Context = nil
	a.Object = nil
	a.FormField = ``
	a.Code = 0
	a.ContentType = ``
	a.JSONPVarName = ``
	a.RespData = nil
	a.Results = nil
	a.err = nil
	a.uploadMaxSize = 0
	a.chunkUpload = nil
	a.fieldMapping = nil
}

func (a *BaseClient) SetUploadMaxSize(maxSize int64) Client {
	a.uploadMaxSize = maxSize
	return a
}

func (a *BaseClient) SetReadBeforeHook(hooks ...ReadBeforeHook) Client {
	a.readBefore = hooks
	return a
}

func (a *BaseClient) AddReadBeforeHook(hooks ...ReadBeforeHook) Client {
	a.readBefore = append(a.readBefore, hooks...)
	return a
}

func (a *BaseClient) UploadMaxSize() int64 {
	if a.uploadMaxSize != 0 {
		return a.uploadMaxSize
	}

	return int64(a.Context.Request().MaxSize())
}

func (a *BaseClient) Name() string {
	if len(a.FormField) == 0 {
		return DefaultFormField
	}
	return a.FormField
}

func (a *BaseClient) SetName(formField string) {
	a.FormField = formField
}

func (a *BaseClient) SetError(err error) Client {
	a.err = err
	return a
}

func (a *BaseClient) GetError() error {
	return a.err
}

func (a *BaseClient) ErrorString() string {
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

func (a *BaseClient) BuildResult() Client {
	data := a.Context.Data()
	data.SetData(echo.H{
		`Url`: a.Data.FileURL,
		`Id`:  a.Data.FileIDString(),
	}, 1)
	if a.err != nil {
		data.SetError(a.err)
	}
	a.RespData = data
	return a
}

func (a *BaseClient) GetRespData() interface{} {
	if a.RespData != nil {
		return a.RespData
	}
	if a.Object != nil {
		a.Object.BuildResult()
	} else {
		a.BuildResult()
	}
	return a.RespData
}

func (a *BaseClient) GetUploadResult() *Result {
	return a.Data
}

func (a *BaseClient) GetBatchUploadResults() Results {
	return a.Results
}

func (a *BaseClient) SetRespData(data interface{}) Client {
	a.RespData = data
	return a
}

func (a *BaseClient) SetChunkUpload(cu *ChunkUpload) Client {
	a.chunkUpload = cu
	return a
}

func (a *BaseClient) SetFieldMapping(fm map[string]string) Client {
	a.fieldMapping = fm
	return a
}

func (a *BaseClient) Response() error {
	if a.Code > 0 {
		return a.responseContentType(a.Code)
	}
	return a.responseContentType(http.StatusOK)
}

func (a *BaseClient) responseContentType(code int) error {
	switch a.ContentType {
	case `string`:
		return a.String(fmt.Sprint(a.GetRespData()), code)
	case `xml`:
		return a.XML(a.GetRespData(), code)
	case `redirect`:
		a.Context.Response().Redirect(fmt.Sprint(a.GetRespData()), code)
		return nil
	default:
		if len(a.JSONPVarName) > 0 {
			return a.JSONP(a.JSONPVarName, a.GetRespData(), code)
		}
		return a.JSON(a.GetRespData(), code)
	}
}
