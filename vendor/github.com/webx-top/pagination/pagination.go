/*

   Copyright since 2017 Wenhui Shen <www.webx.top>

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.

*/

package pagination

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"html/template"
	"math"
	"strconv"
	"strings"

	"github.com/webx-top/echo"
	"github.com/webx-top/echo/engine"
)

const (
	ModePageNumber = iota + 1
	ModePosition
)

func New(ctx echo.Context) *Pagination {
	return &Pagination{context: ctx, pages: -1, data: echo.H{}, mode: ModePageNumber}
}

type Pagination struct {
	context   echo.Context
	tmpl      string
	urlLayout string
	data      echo.H
	mode      int

	// 按基准位置分页
	position     string
	prevPosition string
	nextPosition string

	// 按页码分页
	page  int
	rows  int //total rows
	limit int
	num   int
	pages int //total pages

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

func (p *Pagination) SetPosition(prev string, next string, curr string) *Pagination {
	p.prevPosition = prev
	p.nextPosition = next
	p.position = curr
	p.mode = ModePosition
	return p
}

func (p *Pagination) Set(key string, data interface{}) *Pagination {
	p.data[key] = data
	return p
}

func (p *Pagination) Sets(args ...interface{}) *Pagination {
	var key string
	for i, j := 0, len(args); i < j; i++ {
		if i%2 == 0 {
			key = fmt.Sprint(args[i])
			continue
		}
		p.data[key] = args[i]
	}
	return p
}

func (p *Pagination) Get(key string) interface{} {
	if v, y := p.data[key]; y {
		return v
	}
	return nil
}

func (p *Pagination) Data() echo.H {
	return p.data
}

func (p *Pagination) Position() string {
	return p.position
}

func (p *Pagination) PrevPosition() string {
	return p.prevPosition
}

func (p *Pagination) NextPosition() string {
	return p.nextPosition
}

func (p *Pagination) HasNext() bool {
	if p.mode == ModePageNumber {
		return p.Page() < p.Pages()
	}
	return len(p.NextPosition()) > 0 && p.NextPosition() != `0` && p.NextPosition() != p.Position()
}

func (p *Pagination) HasPrev() bool {
	if p.mode == ModePageNumber {
		return p.Page() > 1
	}
	return p.PrevPosition() != `0` && p.PrevPosition() != p.Position()

}

func (p *Pagination) SetPage(page int) *Pagination {
	p.page = page
	return p
}

func (p *Pagination) Page() int {
	return p.page
}

func (p *Pagination) PrevPage() int {
	if p.page < 2 {
		return 1
	}
	return p.page - 1
}

func (p *Pagination) NextPage() int {
	n := p.page + 1
	if n <= p.pages {
		return n
	}
	return p.pages
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

func (p *Pagination) URL(curr interface{}) (s string) {
	if p.mode == ModePageNumber {
		s = strings.Replace(p.urlLayout, `{page}`, fmt.Sprint(curr), -1)
		s = strings.Replace(s, `{rows}`, strconv.Itoa(p.rows), -1)
		s = strings.Replace(s, `{size}`, strconv.Itoa(p.limit), -1)
		s = strings.Replace(s, `{limit}`, strconv.Itoa(p.limit), -1)
		s = strings.Replace(s, `{pages}`, strconv.Itoa(p.pages), -1)
	} else {
		s = strings.Replace(p.urlLayout, `{curr}`, fmt.Sprint(curr), -1)
		s = strings.Replace(s, `{prev}`, p.prevPosition, -1)
		s = strings.Replace(s, `{next}`, p.nextPosition, -1)
	}
	return s
}

func (p *Pagination) SetURL(s interface{}, delKeys ...string) *Pagination {
	switch v := s.(type) {
	case string:
		p.urlLayout = v
	case map[string]string:
		p.urlLayout = p.RebuildURL(v, delKeys...)
	default:
		panic(`Unsupported type: ` + fmt.Sprintf(`%T`, s))
	}
	return p
}

func (p *Pagination) RebuildQueryString(delKeys ...string) string {
	query := p.context.Request().URL().Query()
	for _, key := range delKeys {
		query.Del(key)
	}
	return query.Encode()
}

func (p *Pagination) RebuildURL(pageVars map[string]string, delKeys ...string) string {
	var (
		pq string
		jn string
	)
	for name, urlVar := range pageVars {
		delKeys = append(delKeys, urlVar)
		pq += jn + urlVar + `={` + name + `}`
		jn = `&`
	}
	q := p.RebuildQueryString(delKeys...)
	if len(q) > 0 {
		q += `&`
	}
	url := p.context.Request().URL().Path() + `?` + q + pq
	return url
}

func (p *Pagination) List(num ...int) []int {
	if len(num) > 0 {
		p.num = num[0]
	}
	pages := p.Pages()
	var (
		pList []int
		start int
		count int
	)
	remainPages := pages - p.page
	if remainPages < p.num {
		start = pages - p.num + 1
	} else {
		start = p.Page() - (p.num / 2)
	}
	if start < 1 {
		start = 1
	}
	for page := start; page <= pages; page++ {
		count++
		if count > p.num {
			break
		}
		pList = append(pList, page)
	}
	return pList
}

func (p *Pagination) setDefault() *Pagination {
	if p.mode == ModePageNumber {
		if p.page < 1 {
			p.page = 1
		}
		if p.limit < 1 {
			p.limit = 50
		}
		if p.num < 1 {
			p.num = 10
		}
	}
	return p
}

func (p *Pagination) Render(settings ...string) interface{} {
	p.setDefault()
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

// MarshalJSON allows type Pagination to be used with json.Marshal
func (p *Pagination) MarshalJSON() ([]byte, error) {
	b, e := json.Marshal(p.data)
	var s string
	if e != nil {
		s = fmt.Sprintf(`%q`, e.Error())
	} else {
		s = engine.Bytes2str(b)
	}
	if p.mode == ModePageNumber {
		p.setDefault()
		s = fmt.Sprintf(`{"page":%d,"rows":%d,"limit":%d,"pages":%d,"urlLayout":%q,"data":%s}`, p.Page(), p.Rows(), p.Limit(), p.Pages(), p.urlLayout, s)
	} else {
		s = fmt.Sprintf(`{"curr":%q,"prev":%q,"next":%q,"urlLayout":%q,"data":%s}`, p.Position(), p.PrevPosition(), p.NextPosition(), p.urlLayout, s)
	}
	return engine.Str2bytes(s), nil
}

func (p *Pagination) SetOptions(m echo.H) *Pagination {
	if _, y := m[`page`]; y {
		p.mode = ModePageNumber
		p.page = m.Int(`page`)
		p.rows = m.Int(`rows`)
		p.limit = m.Int(`limit`)
		p.pages = m.Int(`pages`)
	} else {
		p.mode = ModePosition
		p.position = m.String(`curr`)
		p.prevPosition = m.String(`prev`)
		p.nextPosition = m.String(`next`)
	}
	p.urlLayout = m.String(`urlLayout`)
	p.data = m.Store(`data`)
	return p
}

func (p *Pagination) Options() echo.H {
	m := echo.H{}
	if p.mode == ModePageNumber {
		m.Set(`page`, p.page)
		m.Set(`rows`, p.rows)
		m.Set(`limit`, p.limit)
		m.Set(`pages`, p.pages)
	} else {
		m.Set(`curr`, p.position)
		m.Set(`prev`, p.prevPosition)
		m.Set(`next`, p.nextPosition)
	}
	m.Set(`urlLayout`, p.urlLayout)
	m.Set(`data`, p.data)
	return m
}

// MarshalXML allows type Pagination to be used with xml.Marshal
func (p *Pagination) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	start.Name.Local = `Pagination`
	if err := e.EncodeToken(start); err != nil {
		return err
	}
	if p.mode == ModePageNumber {
		p.setDefault()
		if err := xmlEncode(e, `page`, p.Page()); err != nil {
			return err
		}
		if err := xmlEncode(e, `rows`, p.Rows()); err != nil {
			return err
		}
		if err := xmlEncode(e, `limit`, p.Limit()); err != nil {
			return err
		}
		if err := xmlEncode(e, `pages`, p.Pages()); err != nil {
			return err
		}
	} else {
		if err := xmlEncode(e, `curr`, p.Position()); err != nil {
			return err
		}
		if err := xmlEncode(e, `prev`, p.PrevPosition()); err != nil {
			return err
		}
		if err := xmlEncode(e, `next`, p.NextPosition()); err != nil {
			return err
		}
	}
	if err := xmlEncode(e, `urlLayout`, p.urlLayout); err != nil {
		return err
	}
	if err := xmlEncode(e, `data`, p.data); err != nil {
		return err
	}
	return e.EncodeToken(xml.EndElement{Name: start.Name})
}

func xmlEncode(e *xml.Encoder, key string, value interface{}) error {
	elem := xml.StartElement{
		Name: xml.Name{Space: ``, Local: key},
		Attr: []xml.Attr{},
	}
	return e.EncodeElement(value, elem)
}
