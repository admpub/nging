package pipe

import (
	"path"
	"strings"

	"github.com/admpub/nging/v5/application/registry/upload/driver"
	"github.com/admpub/qrcode"
	uploadClient "github.com/webx-top/client/upload"
	"github.com/webx-top/echo"
)

func init() {
	Register(`deqr`, Deqr)
}

// Deqr 识别二维码
func Deqr(ctx echo.Context, storer driver.Storer, results uploadClient.Results, data map[string]interface{}) error {
	if len(results) == 0 {
		return nil
	}
	reader, err := storer.Get(results[0].SavePath)
	if reader != nil {
		defer reader.Close()
	}
	if err != nil {
		return err
	}
	raw, err := qrcode.Decode(reader, strings.TrimPrefix(path.Ext(results[0].SavePath), `.`))
	if err != nil {
		raw = err.Error()
	}
	data[`raw`] = raw
	return nil
}
