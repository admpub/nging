package dashboard

import (
	"strings"

	"github.com/admpub/nging/v4/application/library/common"
	"github.com/webx-top/echo"
)

func NewPage(key string) *Page {
	return &Page{
		Key:   key,
		hooks: map[string][]func(echo.Context) error{},
	}
}

type Page struct {
	Key      string
	BodyTmpl []string
	HeadTmpl []string
	FootTmpl []string
	hooks    map[string][]func(echo.Context) error
}

func (s *Page) AddBodyTmpl(tmpl ...string) *Page {
	s.BodyTmpl = append(s.BodyTmpl, tmpl...)
	return s
}

func (s *Page) AddHeadTmpl(tmpl ...string) *Page {
	s.HeadTmpl = append(s.HeadTmpl, tmpl...)
	return s
}

func (s *Page) AddFootTmpl(tmpl ...string) *Page {
	s.FootTmpl = append(s.FootTmpl, tmpl...)
	return s
}

func (s *Page) On(method string, hook func(echo.Context) error) *Page {
	method = strings.ToUpper(method)
	if _, ok := s.hooks[method]; !ok {
		s.hooks[method] = []func(echo.Context) error{}
	}
	s.hooks[method] = append(s.hooks[method], hook)
	return s
}

func (s *Page) Fire(ctx echo.Context) error {
	method := strings.ToUpper(ctx.Method())
	hooks, ok := s.hooks[method]
	if !ok {
		return nil
	}
	errs := common.NewErrors()
	for _, hook := range hooks {
		err := hook(ctx)
		if err != nil {
			errs.Add(err)
		}
	}
	return errs.ToError()
}
