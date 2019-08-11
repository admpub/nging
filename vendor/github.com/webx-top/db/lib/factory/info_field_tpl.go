package factory

import (
	"bytes"
	"html/template"
	"sync"
)

var DefaultHTMLTmpl = &HTMLTmpl{
	Group: `<div class="form-group">
	<label class="col-sm-{{.labelCols}} control-label">{{.label}}{{.labelSuffix}}</label>
	<div class="col-sm-{{.inputCols}}">{{.input}}</div>
  	</div>`,
	Required: ` <span class="text-danger star-required">*</span>`,
	Inputs: map[string]string{
		`radio`: `<div class="radio radio-{{if .theme}}{{.theme}}{{else}}primary{{end}}{{if .inline}} radio-inline{{end}}">
		<input type="radio" value="{{.value}}" name="{{.name}}" id="{{.id}}"{{range $k,$v:= .attrs}} {{$v.K}}="{{$v.V}}"{{end}}> <label for="{{.id}}">{{.label}}</label>
	</div>`,
		`checkbox`: `<div class="checkbox checkbox-{{if .theme}}{{.theme}}{{else}}primary{{end}}{{if .inline}} checkbox-inline{{end}}">
		<input type="checkbox" value="{{.value}}" name="{{.name}}" id="{{.id}}"{{range $k,$v:= .attrs}} {{$v.K}}="{{$v.V}}"{{end}}> <label for="{{.id}}">{{.label}}</label>
	</div>`,
		`text`:     `<input type="{{.type}}" class="form-control" name="{{.name}}" value="{{.value}}"{{range $k,$v:= .attrs}} {{$v.K}}="{{$v.V}}"{{end}} />`,
		`textarea`: `<textarea class="form-control" name="{{.name}}" {{range $k,$v:= .attrs}} {{$v.K}}="{{$v.V}}"{{end}}>{{.value}}</textarea>`,
	},
}

type HTMLAttrs []*HTMLAttr

type HTMLAttr struct {
	K template.HTMLAttr
	V template.HTML
}

func (a *HTMLAttrs) Add(k string, v ...string) {
	var value string
	if len(v) > 0 {
		value = v[0]
	} else {
		value = k
	}
	*a = append(*a, &HTMLAttr{template.HTMLAttr(k), template.HTML(value)})
}

type HTMLTmpl struct {
	Inputs   map[string]string
	Group    string
	groupT   *template.Template
	inputT   sync.Map
	Required string
}

func (f *HTMLTmpl) Clear() {
	f.groupT = nil
	f.inputT = sync.Map{}
}

func (f *HTMLTmpl) ToGroup(data interface{}) string {
	if f.groupT == nil {
		f.groupT = template.New(`group`)
		_, err := f.groupT.Parse(f.Group)
		if err != nil {
			f.groupT = nil
			return err.Error()
		}
	}
	buf := bytes.NewBuffer(nil)
	err := f.groupT.Execute(buf, data)
	if err != nil {
		return err.Error()
	}
	return buf.String()
}

func (f *HTMLTmpl) ToInput(typ string, data interface{}) string {
	v, y := f.inputT.Load(typ)
	var t *template.Template
	if !y {
		t = template.New(typ)
		_, err := t.Parse(f.Inputs[typ])
		if err != nil {
			return err.Error()
		}
	} else {
		t = v.(*template.Template)
	}
	buf := bytes.NewBuffer(nil)
	err := t.Execute(buf, data)
	if err != nil {
		return err.Error()
	}
	return buf.String()
}
