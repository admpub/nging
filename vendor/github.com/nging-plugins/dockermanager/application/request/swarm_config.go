package request

import (
	"io"

	"github.com/docker/docker/api/types/swarm"
	"github.com/webx-top/com"
	"github.com/webx-top/echo"
)

var _ echo.ValueDecodersGetter = (*SwarmConfigEdit)(nil)
var _ echo.ValueEncodersGetter = (*SwarmConfigEdit)(nil)
var _ echo.FormNameFormatterGetter = (*SwarmConfigEdit)(nil)
var _ echo.FiltersGetter = (*SwarmConfigEdit)(nil)
var _ echo.AfterValidate = (*SwarmConfigEdit)(nil)

type SwarmConfigEdit struct {
	Content     string
	ContentFrom string `validate:"oneof=input file"`
	swarm.ConfigSpec
}

func (a *SwarmConfigEdit) Filters(echo.Context) []echo.FormDataFilter {
	return []echo.FormDataFilter{
		echo.ExcludeFieldName(`Data`),
	}
}

func (a *SwarmConfigEdit) BeforeValidate(ctx echo.Context) error {
	if a.ContentFrom == `file` {
		a.Content = ``
		fp, _, err := ctx.Request().FormFile(`file`)
		if err != nil {
			return err
		}
		defer fp.Close()
		a.Data, err = io.ReadAll(fp)
		if err != nil {
			return err
		}
	} else {
		if len(a.ContentFrom) == 0 {
			a.ContentFrom = `input`
		}
	}
	return nil
}

func (a *SwarmConfigEdit) AfterValidate(echo.Context) error {
	if a.ContentFrom == `input` {
		if len(a.Content) > 0 {
			a.Data = com.Str2bytes(a.Content)
		}
	}
	if a.Templating != nil {
		if len(a.Templating.Name) == 0 && len(a.Templating.Options) == 0 {
			a.Templating = nil
		}
	}
	return nil
}

func (a *SwarmConfigEdit) ValueDecoders(echo.Context) echo.BinderValueCustomDecoders {
	return map[string]echo.BinderValueCustomDecoder{
		`Labels`:             mapDecoder,
		`Templating.Options`: mapDecoder,
	}
}

func (a *SwarmConfigEdit) FormNameFormatter(echo.Context) echo.FieldNameFormatter {
	return IgnoreDataField()
}

func (a *SwarmConfigEdit) ValueEncoders(echo.Context) echo.BinderValueCustomEncoders {
	return echo.BinderValueCustomEncoders{
		`labels`: func(v interface{}) []string {
			return []string{com.JoinKVRows(v)}
		},
		`templating[options]`: func(v interface{}) []string {
			return []string{com.JoinKVRows(v)}
		},
	}
}
