package pagination

import (
	"html/template"
	"math"
	"strconv"
	"strings"

	"github.com/webx-top/echo"
)

func New(ctx echo.Context) *Pagination {
	return &Pagination{context: ctx, pages: -1}
}

type Pagination struct {
	context   echo.Context
	page      int
	rows      int //total rows
	limit     int
	num       int
	tmpl      string
	pages     int //total pages
	urlLayout string
}

func (p *Pagination) SetAll(tmpl string, rows int, pnl ...int) *Pagination {
	switch len(pnl) {
	case 3:
		p.limit = pnl[2]
		fallthrough
	case 2:
		p.num = pnl[1]
		fallthrough
	case 1:
		p.page = pnl[0]
	}
	p.rows = rows
	p.tmpl = tmpl
	p.pages = -1
	return p
}

func (p *Pagination) SetPage(page int) *Pagination {
	p.page = page
	return p
}

func (p *Pagination) Page() int {
	return p.page
}

func (p *Pagination) SetRows(rows int) *Pagination {
	p.pages = -1
	p.rows = rows
	return p
}

func (p *Pagination) Rows() int {
	return p.rows
}

func (p *Pagination) SetLimit(limit int) *Pagination {
	p.pages = -1
	p.limit = limit
	return p
}

func (p *Pagination) Limit() int {
	return p.limit
}

func (p *Pagination) SetNum(num int) *Pagination {
	p.num = num
	return p
}

func (p *Pagination) Num() int {
	return p.num
}

func (p *Pagination) SetTmpl(tmpl string) *Pagination {
	p.tmpl = tmpl
	return p
}

func (p *Pagination) Tmpl() string {
	return p.tmpl
}

func (p *Pagination) Pages() int {
	if p.pages == -1 {
		p.pages = int(math.Ceil(float64(p.rows) / float64(p.limit)))
	}
	return p.pages
}

func (p *Pagination) URL(page int) string {
	s := strings.Replace(p.urlLayout, `{page}`, strconv.Itoa(page), -1)
	s = strings.Replace(s, `{rows}`, strconv.Itoa(p.rows), -1)
	s = strings.Replace(s, `{limit}`, strconv.Itoa(p.limit), -1)
	s = strings.Replace(s, `{pages}`, strconv.Itoa(p.pages), -1)
	return s
}

func (p *Pagination) SetURL(s string) *Pagination {
	p.urlLayout = s
	return p
}

func (p *Pagination) List(num ...int) []int {
	if len(num) > 0 {
		p.num = num[0]
	}
	r := []int{}
	half := p.num / 2
	lefts := []int{}
	for i, j := p.page, p.num; i > 0 && j > half; i-- {
		lefts = append(lefts, i)
		j--
	}
	c := len(lefts)
	for i := c - 1; i >= 0; i-- {
		r = append(r, lefts[i])
	}
	pages := p.Pages()
	for i := p.page + 1; i <= pages && c < p.num; i++ {
		r = append(r, i)
		c++
	}
	return r
}

func (p *Pagination) Render(settings ...string) interface{} {
	if p.page < 1 {
		p.page = 1
	}
	if p.limit < 1 {
		p.limit = 50
	}
	if p.num < 1 {
		p.num = 10
	}
	switch len(settings) {
	case 1:
		p.tmpl = settings[0]
	}
	if len(p.tmpl) == 0 {
		p.tmpl = `pagination`
	}
	b, e := p.context.Fetch(p.tmpl, p)
	if e != nil {
		return e
	}
	return template.HTML(string(b))
}
