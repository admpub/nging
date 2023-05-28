package fileperm

import (
	"path/filepath"
	"sort"
	"strings"

	"github.com/admpub/nging/v5/application/library/common"
	"github.com/webx-top/echo"
)

func ParseForm(ctx echo.Context) (rules Rules, err error) {
	readables := ctx.FormValues(`readables[]`)
	writeables := ctx.FormValues(`writeables[]`)
	rlen := len(readables)
	wlen := len(writeables)
	resources := ctx.FormValues(`resources[]`)
	rules = make([]*Rule, 0, len(resources))
	for index, path := range resources {
		path = strings.TrimSpace(path)
		if len(path) == 0 {
			continue
		}
		r := &Rule{
			Path: path,
		}
		r.Path = filepath.ToSlash(path)
		if index < rlen {
			r.Readable = readables[index] == `Y`
		}
		if index < wlen {
			switch writeables[index] {
			case `Y`:
				r.SetWriteable(true)
			case `N`:
				r.SetWriteable(false)
			}
		}
		err = rules.Add(r)
		if err != nil {
			return
		}
	}
	sort.Sort(rules)
	return
}

func (s Rules) SetForm(ctx echo.Context) {
	for _, rule := range s {
		ctx.Request().Form().Add(`readables[]`, common.BoolToFlag(rule.Readable))
		if rule.Writeable != nil {
			ctx.Request().Form().Add(`writeables[]`, common.BoolToFlag(*rule.Writeable))
		} else {
			ctx.Request().Form().Add(`writeables[]`, ``)
		}
		ctx.Request().Form().Add(`resources[]`, rule.Path)
	}
}
