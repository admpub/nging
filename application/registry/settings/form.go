package settings

import (
	"html/template"

	"github.com/admpub/nging/v4/application/dbschema"
	"github.com/admpub/nging/v4/application/library/common"
	"github.com/coscms/forms"
	formsconfig "github.com/coscms/forms/config"
	"github.com/webx-top/echo"
)

func NewForm(group string, short string, label string, opts ...FormSetter) *SettingForm {
	f := &SettingForm{
		Group: group,
		Short: short,
		Label: label,
		items: map[string]*dbschema.NgingConfig{},
	}
	for _, opt := range opts {
		opt(f)
	}
	return f
}

type SettingForm struct {
	Short        string   //简短标签
	Label        string   //标签文本
	Group        string   //组标识
	Tmpl         []string //输入表单模板路径
	HeadTmpl     []string
	FootTmpl     []string
	items        map[string]*dbschema.NgingConfig //配置项
	hookPost     []func(echo.Context) error       //数据提交逻辑处理
	hookGet      []func(echo.Context) error       //数据读取逻辑处理
	renderer     func(echo.Context) template.HTML
	config       *formsconfig.Config
	form         *forms.Forms
	dataDecoders DataDecoders //Decoder: table => form
	dataEncoders DataEncoders //Encoder: form => table
}

func (s *SettingForm) AddTmpl(tmpl ...string) *SettingForm {
	s.Tmpl = append(s.Tmpl, tmpl...)
	return s
}

func (s *SettingForm) AddHeadTmpl(tmpl ...string) *SettingForm {
	s.HeadTmpl = append(s.HeadTmpl, tmpl...)
	return s
}

func (s *SettingForm) AddFootTmpl(tmpl ...string) *SettingForm {
	s.FootTmpl = append(s.FootTmpl, tmpl...)
	return s
}

func (s *SettingForm) AddHookPost(hook func(echo.Context) error) *SettingForm {
	s.hookPost = append(s.hookPost, hook)
	return s
}

func (s *SettingForm) SetFormConfig(formcfg *formsconfig.Config) *SettingForm {
	s.config = formcfg
	if s.form != nil {
		s.form = nil
	}
	return s
}

// SetDataTransfer 数据转换
// dataDecoder: Decoder(table => form)
// dataEncoder: Encoder(form => table)
func (s *SettingForm) SetDataTransfer(name string, dataDecoder DataDecoder, dataEncoder DataEncoder) *SettingForm {
	if dataDecoder != nil {
		if s.dataDecoders == nil {
			s.dataDecoders = DataDecoders{}
		}
		s.dataDecoders[name] = dataDecoder
	}
	if dataEncoder != nil {
		if s.dataEncoders == nil {
			s.dataEncoders = DataEncoders{}
		}
		s.dataEncoders[name] = dataEncoder
	}
	return s
}

func (s *SettingForm) AddConfig(configs ...*dbschema.NgingConfig) *SettingForm {
	if s.items == nil {
		s.items = map[string]*dbschema.NgingConfig{}
	}
	for _, c := range configs {
		if c.Group != s.Group {
			c.Group = s.Group
		}
		s.items[c.Key] = c
	}
	return s
}

func (s *SettingForm) SetRenderer(renderer func(echo.Context) template.HTML) *SettingForm {
	s.renderer = renderer
	return s
}

func (s *SettingForm) Render(ctx echo.Context) template.HTML {
	if s.renderer != nil {
		return s.renderer(ctx)
	}
	if s.config != nil {
		if s.form == nil {
			s.form = forms.NewForms(forms.NewWithConfig(s.config))
			m := echo.GetStore(`NgingConfig`).GetStore(s.Group)
			s.form.SetModel(m)
			s.form.ParseFromConfig(true)
		}
		return s.form.Render()
	}
	var htmlContent string
	var stored echo.Store
	if fn, ok := ctx.GetFunc(`Stored`).(func() echo.Store); ok {
		stored = fn()
	} else {
		stored = ctx.Stored()
	}
	for _, t := range s.Tmpl {
		if len(t) == 0 {
			continue
		}
		b, err := ctx.Fetch(t, stored)
		if err != nil {
			htmlContent += err.Error()
		} else {
			htmlContent += string(b)
		}
	}
	return template.HTML(htmlContent)
}

func (s *SettingForm) AddHookGet(hook func(echo.Context) error) *SettingForm {
	s.hookGet = append(s.hookGet, hook)
	return s
}

func (s *SettingForm) RunHookPost(ctx echo.Context) error {
	if s.hookPost == nil {
		return nil
	}
	errs := common.NewErrors()
	for _, hook := range s.hookPost {
		err := hook(ctx)
		if err != nil {
			errs.Add(err)
		}
	}
	return errs.ToError()
}

func (s *SettingForm) RunHookGet(ctx echo.Context) error {
	if s.hookGet == nil {
		return nil
	}
	errs := common.NewErrors()
	for _, hook := range s.hookGet {
		err := hook(ctx)
		if err != nil {
			errs.Add(err)
		}
	}
	return errs.ToError()
}
