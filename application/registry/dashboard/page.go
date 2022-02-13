package dashboard

import (
	"strings"

	"github.com/admpub/nging/v4/application/library/common"
	"github.com/webx-top/echo"
)

func NewPage(key string, atmpls ...map[string][]string) *Page {
	var tmpls map[string][]string
	if len(atmpls) > 0 {
		tmpls = atmpls[0]
	}
	if tmpls == nil {
		tmpls = map[string][]string{}
	}
	return &Page{
		Key:   key,
		Tmpls: tmpls,
		hooks: map[string][]func(echo.Context) error{},
	}
}

type Page struct {
	Key   string
	Tmpls map[string][]string
	hooks map[string][]func(echo.Context) error
}

func (s *Page) AddTmpl(position string, tmpl ...string) *Page {
	if _, ok := s.Tmpls[position]; !ok {
		s.Tmpls[position] = []string{}
	}
	s.Tmpls[position] = append(s.Tmpls[position], tmpl...)
	return s
}

func (s *Page) Tmpl(position string) []string {
	return s.Tmpls[position]
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
