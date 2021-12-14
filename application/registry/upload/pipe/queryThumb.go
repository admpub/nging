package pipe

import (
	"regexp"
	"strings"

	modelFile "github.com/admpub/nging/v4/application/model/file"
	uploadChecker "github.com/admpub/nging/v4/application/registry/upload/checker"
	"github.com/admpub/nging/v4/application/registry/upload/driver"
	uploadClient "github.com/webx-top/client/upload"
	"github.com/webx-top/db"
	"github.com/webx-top/echo"
)

func init() {
	Register(`_queryThumb`, queryThumb) // 以下划线开始表示这个独立的功能
}

// WidthAndHeightRegexp 宽和高
var WidthAndHeightRegexp = regexp.MustCompile(`^[\d]+x[\d]+$`)

// queryThumb 查询缩略图
func queryThumb(ctx echo.Context, _ driver.Storer, _ uploadClient.Results, data map[string]interface{}) error {
	viewURL := ctx.Form(`file`)
	size := ctx.Form(`size`)
	if len(size) == 0 {
		return ctx.E(`尺寸格式不正确`)
	}
	if !WidthAndHeightRegexp.MatchString(size) {
		return ctx.E(`尺寸格式不正确`)
	}
	sizes := strings.SplitN(size, "x", 2)
	width := sizes[0]
	height := sizes[1]
	m := modelFile.NewThumb(ctx)
	viewURL = modelFile.GetViewURLByOriginalURL(viewURL, width, height)
	//panic(viewURL)
	err := m.GetByViewURL(viewURL)
	if err != nil {
		if err == db.ErrNoMoreRows {
			return nil
		}
		return err
	}
	data[`thumb`] = viewURL
	data[`token`] = uploadChecker.Token(`file`, viewURL, `width`, width, `height`, height)
	return nil
}
