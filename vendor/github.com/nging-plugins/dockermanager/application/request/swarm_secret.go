package request

import (
	"io"

	"github.com/docker/docker/api/types/swarm"
	"github.com/webx-top/com"
	"github.com/webx-top/echo"
)

var _ echo.ValueDecodersGetter = (*SwarmSecretEdit)(nil)
var _ echo.ValueEncodersGetter = (*SwarmSecretEdit)(nil)
var _ echo.FormNameFormatterGetter = (*SwarmSecretEdit)(nil)
var _ echo.FiltersGetter = (*SwarmSecretEdit)(nil)
var _ echo.AfterValidate = (*SwarmSecretEdit)(nil)

type SwarmSecretEdit struct {
	Content     string
	ContentFrom string `validate:"oneof=input file"`
	swarm.SecretSpec
}

func (a *SwarmSecretEdit) Filters(echo.Context) []echo.FormDataFilter {
	return []echo.FormDataFilter{
		echo.ExcludeFieldName(`Data`),
	}
}

func (a *SwarmSecretEdit) BeforeValidate(ctx echo.Context) error {
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

func (a *SwarmSecretEdit) AfterValidate(echo.Context) error {
	if a.ContentFrom == `input` {
		if len(a.Content) > 0 {
			a.Data = com.Str2bytes(a.Content)
		}
	}
	if a.Driver != nil {
		if len(a.Driver.Name) == 0 && len(a.Driver.Options) == 0 {
			a.Driver = nil
		}
	}
	if a.Templating != nil {
		if len(a.Templating.Name) == 0 && len(a.Templating.Options) == 0 {
			a.Templating = nil
		}
	}
	return nil
}

func (a *SwarmSecretEdit) ValueDecoders(echo.Context) echo.BinderValueCustomDecoders {
	return map[string]echo.BinderValueCustomDecoder{
		`Labels`:             mapDecoder,
		`Driver.Options`:     mapDecoder,
		`Templating.Options`: mapDecoder,
	}
}

func (a *SwarmSecretEdit) FormNameFormatter(_ echo.Context) echo.FieldNameFormatter {
	return IgnoreDataField()
}

func (a *SwarmSecretEdit) ValueEncoders(_ echo.Context) echo.BinderValueCustomEncoders {
	return echo.BinderValueCustomEncoders{
		`labels`: func(v interface{}) []string {
			return []string{com.JoinKVRows(v)}
		},
		`dirver[options]`: func(v interface{}) []string {
			return []string{com.JoinKVRows(v)}
		},
		`templating[options]`: func(v interface{}) []string {
			return []string{com.JoinKVRows(v)}
		},
	}
}

func IgnoreDataField() echo.FieldNameFormatter {
	fft := echo.MakeArrayFieldNameFormatter(com.LowerCaseFirst)
	return func(topName, fieldName string) string {
		if fieldName == `data` {
			return ``
		}
		return fft(topName, fieldName)
	}
}
