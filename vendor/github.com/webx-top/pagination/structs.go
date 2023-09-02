package pagination

import "github.com/webx-top/echo"

type pageData struct {
	Page      int    `json:"page" xml:"page"`
	Rows      int    `json:"rows" xml:"rows"`
	Size      int    `json:"size" xml:"size"`
	Limit     int    `json:"limit" xml:"limit"`
	Pages     int    `json:"pages" xml:"pages"`
	URLLayout string `json:"urlLayout" xml:"urlLayout"`
	Data      echo.H `json:"data,omitempty" xml:"data,omitempty"`
}

func (d *pageData) Apply(p *Pagination) {
	p.page = d.Page
	p.rows = d.Rows
	p.size = d.Size
	if p.size <= 0 && d.Limit > 0 {
		p.size = d.Limit
	}
	p.pages = d.Pages
	p.urlLayout = d.URLLayout
	p.data = d.Data
}

type positionData struct {
	Curr      string `json:"curr" xml:"curr"`
	Prev      string `json:"prev" xml:"prev"`
	Next      string `json:"next" xml:"next"`
	Size      int    `json:"size" xml:"size"`
	Limit     int    `json:"limit" xml:"limit"`
	URLLayout string `json:"urlLayout" xml:"urlLayout"`
	Data      echo.H `json:"data,omitempty" xml:"data,omitempty"`
}

func (d *positionData) Apply(p *Pagination) {
	p.position = d.Curr
	p.prevPosition = d.Prev
	p.nextPosition = d.Next
	p.size = d.Size
	if p.size <= 0 && d.Limit > 0 {
		p.size = d.Limit
	}
	p.urlLayout = d.URLLayout
	p.data = d.Data
}

type Applier interface {
	Apply(*Pagination)
}
