package formbuilder

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/coscms/forms"
	"github.com/coscms/forms/common"
	"github.com/coscms/forms/config"
	"github.com/coscms/forms/fields"
	"gopkg.in/yaml.v3"

	"github.com/webx-top/com"
	"github.com/webx-top/db/lib/factory"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/formfilter"
	echoMw "github.com/webx-top/echo/middleware"
	"github.com/webx-top/echo/middleware/render/driver"
)

var (
	ErrJSONConfigFileNameInvalid = errors.New("*.form.json name invalid")
)

func New(ctx echo.Context, model interface{}, options ...Option) *FormBuilder {
	f := &FormBuilder{
		Forms: forms.NewForms(forms.New()),
		on:    MethodHooks{},
		ctx:   ctx,
	}
	defaultHooks := []MethodHook{
		BindModel(f),
		ValidModel(f),
	}
	f.OnPost(defaultHooks...)
	f.OnPut(defaultHooks...)
	f.SetModel(model)
	f.Theme = common.BOOTSTRAP
	for _, option := range options {
		if option == nil {
			continue
		}
		option(f)
	}
	f.SetLabelFunc(func(txt string) string {
		return ctx.T(txt)
	})
	f.AddBeforeRender(func() {
		nextURL := ctx.Form(echo.DefaultNextURLVarName)
		if len(nextURL) > 0 {
			f.Elements(fields.HiddenField(echo.DefaultNextURLVarName).SetValue(nextURL))
		}
	})
	if csrfToken := ctx.Internal().String(echoMw.DefaultCSRFConfig.ContextKey); len(csrfToken) > 0 {
		f.AddBeforeRender(func() {
			f.Elements(fields.HiddenField(echoMw.DefaultCSRFConfig.ContextKey).SetValue(csrfToken))
		})
	}
	ctx.Set(`forms`, f.Forms)
	return f
}

// FormBuilder HTML表单构建器
type FormBuilder struct {
	*forms.Forms
	on         MethodHooks
	exit       bool
	err        error
	ctx        echo.Context
	configFile string
	dbi        *factory.DBI
	defaults   map[string]string
	filters    []formfilter.Options
}

// Exited 是否需要退出后续处理。此时一般有err值，用于记录错误原因
func (f *FormBuilder) Exited() bool {
	return f.exit
}

// Exit 设置退出标记
func (f *FormBuilder) Exit(exit ...bool) *FormBuilder {
	if len(exit) > 0 && !exit[0] {
		f.exit = false
	} else {
		f.exit = true
	}
	return f
}

// SetError 记录错误
func (f *FormBuilder) SetError(err error) *FormBuilder {
	f.err = err
	return f
}

// HasError 是否有错误
func (f *FormBuilder) HasError() bool {
	return f.err != nil
}

// Error 返回错误值
func (f *FormBuilder) Error() error {
	return f.err
}

// ParseConfigFile 解析配置文件 xxx.form.json
func (f *FormBuilder) ParseConfigFile(jsonformat ...bool) error {
	configFile := f.configFile + `.form`
	var cfg *config.Config
	renderer, ok := f.ctx.Renderer().(driver.Driver)
	if !ok {
		return fmt.Errorf(`FormBuilder: Expected renderer is "driver.Driver", but got "%T"`, f.ctx.Renderer())
	}
	var isJSON bool
	if len(jsonformat) > 0 {
		isJSON = jsonformat[0]
	}
	if isJSON {
		configFile += `.json`
	} else {
		configFile += `.yaml`
	}
	configFile = renderer.TmplPath(f.ctx, configFile)
	if len(configFile) == 0 {
		return ErrJSONConfigFileNameInvalid
	}
	b, err := renderer.RawContent(configFile)
	if err != nil {
		if !os.IsNotExist(err) || renderer.Manager() == nil {
			return fmt.Errorf(`read file %s: %w`, configFile, err)
		}
		cfg = f.ToConfig()
		if isJSON {
			b, err = f.ToJSONBlob(cfg)
			if err != nil {
				return fmt.Errorf(`[form.ToJSONBlob] %s: %w`, configFile, err)
			}
		} else {
			b, err = yaml.Marshal(cfg)
			if err != nil {
				return fmt.Errorf(`[form:yaml.Marshal] %s: %w`, configFile, err)
			}
		}
		err = renderer.Manager().SetTemplate(configFile, b)
		if err != nil {
			return fmt.Errorf(`%s: %w`, configFile, err)
		}
		f.ctx.Logger().Infof(f.ctx.T(`生成表单配置文件“%v”成功。`), configFile)
	} else {
		if isJSON {
			cfg, err = forms.Unmarshal(b, configFile)
			if err != nil {
				return fmt.Errorf(`[forms.Unmarshal] %s: %w`, configFile, err)
			}
		} else {
			cfg, err = common.GetOrSetCachedConfig(configFile, func() (*config.Config, error) {
				cfg := &config.Config{}
				err := yaml.Unmarshal(b, cfg)
				return cfg, err
			})
			if err != nil {
				return fmt.Errorf(`[form:yaml.Unmarshal] %s: %w`, configFile, err)
			}
		}
	}
	if cfg == nil {
		cfg = f.NewConfig()
	}

	defaultValues := f.DefaultValues()
	if defaultValues != nil && len(defaultValues) > 0 {
		cfg.SetDefaultValue(func(fieldName string) string {
			fieldName = com.Title(fieldName)
			val, _ := defaultValues[fieldName]
			return val
		})
	}
	f.Init(cfg)
	return err
}

// DefaultValues 获取model结构体各个字段在数据库中的默认值
func (f *FormBuilder) DefaultValues() map[string]string {
	if f.defaults != nil {
		return f.defaults
	}
	if f.dbi == nil || f.dbi.Fields == nil {
		return nil
	}
	m, ok := f.Model.(factory.Short)
	if !ok {
		return nil
	}
	fields, ok := f.dbi.Fields[m.Short_()]
	if !ok || fields == nil {
		return nil
	}
	f.defaults = map[string]string{}
	for _, info := range fields {
		if len(info.DefaultValue) > 0 {
			f.defaults[info.GoName] = info.DefaultValue
		}
	}
	return f.defaults
}

// DefaultValue 查询某个结构体字段在数据库中对应的默认值
func (f *FormBuilder) DefaultValue(fieldName string) string {
	defaultValues := f.DefaultValues()
	if defaultValues == nil {
		return ``
	}
	fieldName = com.Title(fieldName)
	val, _ := defaultValues[fieldName]
	return val
}

// RecvSubmission 接收客户端的提交
func (f *FormBuilder) RecvSubmission() error {
	ctx := f.ctx
	method := strings.ToUpper(ctx.Method())
	if f.err = f.on.Fire(method); f.err != nil {
		return f.err
	}
	f.err = f.on.Fire(`*`)
	if ctx.Response().Committed() {
		f.exit = true
	}
	return f.err
}

// Generate 生成表单参数
func (f *FormBuilder) Generate() *FormBuilder {
	f.ParseFromConfig()
	return f
}

// Snippet 表单片段
func (f *FormBuilder) Snippet() *FormBuilder {
	f.Config().Template = common.TmplDir(f.Config().Theme) + `/allfields.html`
	return f
}
