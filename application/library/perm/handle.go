package perm

import (
	"github.com/webx-top/echo"
)

func NewHandle() *Handle {
	return &Handle{}
}

type Handler interface {
	SetGenerator(fn func(ctx echo.Context) (string, error)) Handler
	SetChecker(fn func(ctx echo.Context, parsed interface{}, current string) (interface{}, error)) Handler
	SetParser(fn func(ctx echo.Context, rule string) (interface{}, error)) Handler
	SetItemLister(fn func(ctx echo.Context) ([]interface{}, error)) Handler
	OnRender(fn func(ctx echo.Context) error) Handler
	SetIsValid(fn func(ctx echo.Context) bool) Handler
	SetTmpl(tmpl string, typ ...string) Handler
	Tmpl(typ ...string) string
	Generate(ctx echo.Context) (string, error)
	Check(ctx echo.Context, parsed interface{}, current string) (interface{}, error)
	Parse(ctx echo.Context, rule string) (interface{}, error)
	ListItems(ctx echo.Context) ([]interface{}, error)
	FireRender(ctx echo.Context) error
	IsValid(ctx echo.Context) bool
}

type Handle struct {
	generator  func(ctx echo.Context) (string, error)
	checker    func(ctx echo.Context, parsed interface{}, current string) (interface{}, error)
	parser     func(ctx echo.Context, rule string) (interface{}, error)
	itemLister func(ctx echo.Context) ([]interface{}, error)
	onRender   func(ctx echo.Context) error
	isValid    func(ctx echo.Context) bool
	tmpl       string
	tmplHead   string
	tmplFoot   string
}

func (u *Handle) SetGenerator(fn func(ctx echo.Context) (string, error)) Handler {
	u.generator = fn
	return u
}

func (u *Handle) SetIsValid(fn func(ctx echo.Context) bool) Handler {
	u.isValid = fn
	return u
}

func (u *Handle) SetChecker(fn func(ctx echo.Context, parsed interface{}, current string) (interface{}, error)) Handler {
	u.checker = fn
	return u
}

func (u *Handle) SetParser(fn func(ctx echo.Context, rule string) (interface{}, error)) Handler {
	u.parser = fn
	return u
}

func (u *Handle) SetItemLister(fn func(ctx echo.Context) ([]interface{}, error)) Handler {
	u.itemLister = fn
	return u
}

func (u *Handle) OnRender(fn func(ctx echo.Context) error) Handler {
	u.onRender = fn
	return u
}

func (u *Handle) SetTmpl(tmpl string, typ ...string) Handler {
	if len(typ) == 0 {
		u.tmpl = tmpl
	} else {
		switch typ[0] {
		case `head`:
			u.tmplHead = tmpl
		case `foot`:
			u.tmplFoot = tmpl
		}
	}
	return u
}

func (u *Handle) Tmpl(typ ...string) string {
	if len(typ) == 0 {
		return u.tmpl
	}
	switch typ[0] {
	case `head`:
		return u.tmplHead
	case `foot`:
		return u.tmplFoot
	default:
		return ``
	}
}

func (u *Handle) Generate(ctx echo.Context) (string, error) {
	if u.generator == nil {
		return ``, nil
	}
	return u.generator(ctx)
}

func (u *Handle) Check(ctx echo.Context, parsed interface{}, current string) (interface{}, error) {
	if u.checker == nil {
		return nil, nil
	}
	return u.checker(ctx, parsed, current)
}

func (u *Handle) Parse(ctx echo.Context, rule string) (interface{}, error) {
	if u.parser == nil {
		return nil, nil
	}
	return u.parser(ctx, rule)
}

func (u *Handle) ListItems(ctx echo.Context) ([]interface{}, error) {
	if u.itemLister == nil {
		return []interface{}{}, nil
	}
	return u.itemLister(ctx)
}

func (u *Handle) FireRender(ctx echo.Context) error {
	if u.onRender == nil {
		return nil
	}
	return u.onRender(ctx)
}

func (u *Handle) IsValid(ctx echo.Context) bool {
	if u.isValid == nil {
		return true
	}
	return u.isValid(ctx)
}

func HandleFireRender(ctx echo.Context, config *echo.KVData) (err error) {
	for _, item := range config.Slice() {
		err = item.X.(*Handle).FireRender(ctx)
		if err != nil {
			break
		}
	}
	return
}

func HandleGenerate(ctx echo.Context, config *echo.KVData) (mp map[string]string, err error) {
	mp = map[string]string{}
	var rule string
	for _, item := range config.Slice() {
		rule, err = item.X.(*Handle).Generate(ctx)
		if err != nil {
			break
		}
		if len(rule) > 0 {
			mp[item.K] = rule
		}
	}
	return
}

func HandleCheck(ctx echo.Context, config *echo.KVData, current string, typ string, permission string, parsed interface{}) (interface{}, error) {
	item := config.GetItem(typ)
	if item == nil {
		return nil, nil
	}
	h := item.X.(*Handle)
	if parsed == nil {
		var err error
		parsed, err = h.Parse(ctx, permission)
		if err != nil {
			return nil, err
		}
	}
	return h.Check(ctx, parsed, current)
}

func HandleBatchCheck(ctx echo.Context, config *echo.KVData, current string, rules map[string]string) (mp map[string]interface{}, err error) {
	mp = map[string]interface{}{}
	for typ, permission := range rules {
		item := config.GetItem(typ)
		if item == nil {
			continue
		}
		h := item.X.(*Handle)
		var (
			parsed interface{}
			result interface{}
		)
		parsed, err = h.Parse(ctx, permission)
		if err != nil {
			break
		}
		result, err = h.Check(ctx, parsed, current)
		if err != nil {
			break
		}
		mp[item.K] = result
	}
	return
}
