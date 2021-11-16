package upload

import (
	"fmt"
	"net/http"

	"github.com/webx-top/echo"
)

// DefaultUploadMaxSize 默认最大上传尺寸: 5MB
var DefaultUploadMaxSize int64 = 5 * 1024 * 1024

func New(object Client, formFields ...string) *BaseClient {
	formField := DefaultFormField
	if len(formFields) > 0 {
		formField = formFields[0]
	}
	return &BaseClient{
		Object:        object,
		FormField:     formField,
		Results:       Results{},
		uploadMaxSize: DefaultUploadMaxSize,
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
	uploadMaxSize int64
	chunkUpload   *ChunkUpload
	fieldMapping  map[string]string
}

func (a *BaseClient) Init(ctx echo.Context, res *Result) {
	a.Context = ctx
	a.Data = res
}

func (a *BaseClient) SetUploadMaxSize(maxSize int64) Client {
	a.uploadMaxSize = maxSize
	return a
}

func (a *BaseClient) UploadMaxSize() int64 {
	return a.uploadMaxSize
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
		`Id`:  a.Data.FileIdString(),
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
