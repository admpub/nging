package formbuilder

import (
	"errors"
	"fmt"
	"os"
	"strings"

	ncommon "github.com/admpub/nging/application/library/common"

	"github.com/coscms/forms"
	"github.com/coscms/forms/common"
	"github.com/coscms/forms/config"
	"github.com/coscms/forms/fields"

	"github.com/webx-top/echo"
	"github.com/webx-top/echo/middleware/render/driver"
)

var (
	ErrJSONConfigFileNameInvalid = errors.New("*.form.json name invalid")
)

func New(ctx echo.Context, model interface{}, options ...Option) *FormBuilder {
	f := &FormBuilder{
		Forms: forms.NewForms(forms.New()),
		on: MethodHooks{
			echo.POST: {BindModel, ValidModel},
			echo.PUT:  {BindModel, ValidModel},
		},
		ctx: ctx,
	}
	f.SetModel(model)
	f.Style = common.BOOTSTRAP
	for _, option := range options {
		if option == nil {
			continue
		}
		option(ctx, f)
	}
	f.SetLabelFunc(func(txt string) string {
		return ctx.T(txt)
	})
	f.AddBeforeRender(func() {
		nextURL := ctx.Form(ncommon.DefaultReturnToURLVarName)
		if len(nextURL) > 0 {
			f.Elements(fields.HiddenField(ncommon.DefaultReturnToURLVarName).SetValue(nextURL))
		}
	})
	csrfToken, ok := ctx.Get(`csrf`).(string)
	if ok {
		f.AddBeforeRender(func() {
			f.Elements(fields.HiddenField(`csrf`).SetValue(csrfToken))
		})
	}
	ctx.Set(`forms`, f.Forms)
	return f
}

type FormBuilder struct {
	*forms.Forms
	on         MethodHooks
	exit       bool
	err        error
	ctx        echo.Context
	configFile string
}

func (f *FormBuilder) Exited() bool {
	return f.exit
}

func (f *FormBuilder) Exit(exit ...bool) *FormBuilder {
	if len(exit) > 0 && !exit[0] {
		f.exit = false
	} else {
		f.exit = true
	}
	return f
}

func (f *FormBuilder) SetError(err error) *FormBuilder {
	f.err = err
	return f
}

func (f *FormBuilder) HasError() bool {
	return f.err != nil
}

func (f *FormBuilder) Error() error {
	return f.err
}

func (f *FormBuilder) ParseConfigFile() error {
	jsonFile := f.configFile
	var cfg *config.Config
	renderer := f.ctx.Renderer().(driver.Driver)
	jsonFile += `.form.json`
	jsonFile = renderer.TmplPath(f.ctx, jsonFile)
	if len(jsonFile) == 0 {
		return ErrJSONConfigFileNameInvalid
	}
	b, err := renderer.RawContent(jsonFile)
	if err != nil {
		if !os.IsNotExist(err) /* && !strings.Contains(err.Error(), `file does not exist`)*/ || renderer.Manager() == nil {
			return fmt.Errorf(`read file %s: %w`, jsonFile, err)
		}
		cfg = f.ToConfig()
		var jsonb []byte
		jsonb, err = f.ToJSONBlob(cfg)
		if err != nil {
			return fmt.Errorf(`[form.ToJSONBlob] %s: %w`, jsonFile, err)
		}
		err = renderer.Manager().SetTemplate(jsonFile, jsonb)
		if err != nil {
			return fmt.Errorf(`%s: %w`, jsonFile, err)
		}
		f.ctx.Logger().Infof(f.ctx.T(`生成表单配置文件“%v”成功。`), jsonFile)
	} else {
		cfg, err = forms.Unmarshal(b, jsonFile)
		if err != nil {
			return fmt.Errorf(`[forms.Unmarshal] %s: %w`, jsonFile, err)
		}
	}
	if cfg == nil {
		cfg = f.NewConfig()
	}
	f.Init(cfg)
	return err
}

func (f *FormBuilder) RecvSubmission() error {
	method := strings.ToUpper(f.ctx.Method())
	if f.err = f.on.Fire(method, f.ctx, f); f.err != nil {
		return f.err
	}
	f.err = f.on.Fire(`*`, f.ctx, f)
	if f.ctx.Response().Committed() {
		f.exit = true
	}
	return f.err
}

func (f *FormBuilder) Generate() *FormBuilder {
	f.ParseFromConfig()
	return f
}
