package utils

import (
	"github.com/admpub/nging/v5/application/handler"
	"github.com/webx-top/echo"
)

func AjaxListSelectpage[T any](ctx echo.Context, list []T, callback func(v T) echo.H) error {
	data := ctx.Data()
	rows := make([]echo.H, 0, len(list))
	var sk, n int
	_, size, _, pg := handler.PagingWithPagination(ctx)
	pg.SetRows(len(list))
	offset := pg.Offset()
	for _, v := range list {
		if n >= size {
			break
		}
		sk++
		if sk-1 < offset {
			continue
		}
		row := callback(v)
		if row == nil {
			continue
		}
		rows = append(rows, row)
		n++
	}
	data.SetData(echo.H{`listData`: rows, `pagination`: pg})
	return ctx.JSON(data)
}

func AjaxListTypeahead[T any](ctx echo.Context, list []T, callback func(v T) string) error {
	data := ctx.Data()
	names := make([]string, 0, len(list))
	for _, v := range list {
		name := callback(v)
		if len(name) == 0 {
			continue
		}
		names = append(names, name)
	}
	data.SetData(echo.H{`listData`: names})
	return ctx.JSON(data)
}
