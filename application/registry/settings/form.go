package settings

import (
	"html/template"

	"github.com/admpub/nging/v3/application/dbschema"
	"github.com/coscms/forms"
	formsconfig "github.com/coscms/forms/config"
	"github.com/webx-top/echo"
)

func NewForm() *SettingForm {
	return &SettingForm{
		items: map[string]*dbschema.NgingConfig{},
	}
}

type SettingForm struct {
	Short       string                           //简短标签
	Label       string                           //标签文本
	Group       string                           //组标识
	Tmpl        []string                         //输入表单模板路径
	items       map[string]*dbschema.NgingConfig //配置项
	hookPost    []func(echo.Context) error       //数据提交逻辑处理
	hookGet     []func(echo.Context) error       //数据读取逻辑处理
	renderer    func(echo.Context) template.HTML
	config      *formsconfig.Config
	form        *forms.Forms
	dataInitors DataInitors
	dataFroms   DataFroms
}

func (s *SettingForm) AddTmpl(tmpl string) *SettingForm {
	s.Tmpl = append(s.Tmpl, tmpl)
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

func (s *SettingForm) SetDataTransfer(name string, dataInitor DataInitor, dataFrom DataFrom) *SettingForm {
	if dataInitor != nil {
		if s.dataInitors == nil {
			s.dataInitors = DataInitors{}
		}
		s.dataInitors[name] = dataInitor
	}
	if dataFrom != nil {
		if s.dataFroms == nil {
			s.dataFroms = DataFroms{}
		}
		s.dataFroms[name] = dataFrom
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
	for _, t := range s.Tmpl {
		if len(t) == 0 {
			continue
		}
		b, err := ctx.Fetch(t, ctx.Stored())
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
	for _, hook := range s.hookPost {
		err := hook(ctx)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *SettingForm) RunHookGet(ctx echo.Context) error {
	if s.hookGet == nil {
		return nil
	}
	for _, hook := range s.hookGet {
		err := hook(ctx)
		if err != nil {
			return err
		}
	}
	return nil
}
