package respond

import "github.com/webx-top/echo"

// DropzoneData 数据
type DropzoneData struct {
	URL string `json:"url" xml:"url"`
	ID  string `json:"id" xml:"id"`
}

var dropzoneDefaultData = &DropzoneData{}

// Dropzone 响应内容
func Dropzone(ctx echo.Context, err error, data *DropzoneData) error {
	if err != nil {
		return ctx.String(err.Error(), 500)
	}
	if data == nil {
		data = dropzoneDefaultData
	}
	result := echo.H{
		`error`:  nil,
		`result`: data,
	}
	return ctx.JSON(result)
}
