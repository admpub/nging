package collector

import (
	"errors"
	"regexp"
	"strconv"

	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/moxar/arithmetic"
)

var ErrMissingParam = errors.New(`missing param`)

func NewGoQuery(ctx *Context) *GoQuery {
	return &GoQuery{Context: ctx}
}

type GoQuery struct {
	Context *Context
}

func (g *GoQuery) Parse() (err error) {
	for index, page := range g.Context.Pages {
		err = g.ParsePage(index, page)
		if err != nil {
			break
		}
	}
	return
}

func (g *GoQuery) ParsePage(index int, config *PageConfig) (err error) {
	reader := g.Context.Reader(config)
	var doc *goquery.Document
	doc, err = goquery.NewDocumentFromReader(reader)
	config.Rule = strings.TrimSpace(config.Rule)
	for _, rule := range strings.Split(config.Rule, "\n") {
		if len(rule) == 0 {
			continue
		}
		//$('#ID').find('.class').text => title
		//$('#ID').find('.class').text => @next@params->post->name
		//$('#ID').find('.class').text => @next@url
		info := strings.SplitN(rule, `=>`, 2)
		if len(info) != 2 {
			continue
		}
		rule = strings.TrimSpace(info[0])
		name := strings.TrimSpace(info[1])
		if len(name) < 1 {
			continue
		}
		var value string
		value, err = g.parseSelections(doc.Selection, strings.Split(rule, "."))
		if err != nil {
			return
		}
		if name[0] == '@' {
			name = name[1:]
			if pos := strings.Index(name, `@`); pos > -1 {
				var idx int
				idxStr := name[0:pos]
				othVal := name[pos+1:]
				switch idxStr {
				case `next`:
					idx = index + 1
				default:
					if i, e := strconv.Atoi(idxStr); e == nil {
						idx = i
					}
				}
				if idx >= len(g.Context.Pages) {
					continue
				}
				var typ string
				args := []string{}
				for i, v := range strings.Split(othVal, `->`) {
					v = strings.TrimSpace(v)
					if i == 0 {
						typ = v
						continue
					}
					args = append(args, v)
				}
				argsNum := len(args)
				switch typ {
				case `params`:
					if argsNum < 3 {
						args = append(args, value)
					}
					config.SetParam(args...)
				case `url`:
					if argsNum > 0 {
						g.Context.Pages[idx].URL = args[0]
					} else {
						g.Context.Pages[idx].URL = value
					}
				}
			}
			continue
		}
		g.Context.Set(name, value)
	}
	return
}

func (g *GoQuery) parseSelections(s *goquery.Selection, selections []string) (r string, err error) {
	for _, selector := range selections {
		var (
			function  = selector
			selection string
			quote     string
		)
		if pos := strings.Index(selector, "("); pos > 0 {
			function = selector[0:pos]
			selection = selector[pos+1:]
			selection = strings.TrimSuffix(selection, ")")
			if len(selection) > 0 {
				quote = selection[0:1]
				selection = strings.Trim(selection, quote)
			}
		}
		switch function {
		case "$":
			s = g.Document.Selection
			if len(selection) > 0 {
				s = s.Find(selection)
			}
		case "find":
			s = s.Find(selection)
		case "children":
			if len(selection) > 0 {
				s = s.ChildrenFiltered(selection)
			} else {
				s = s.Children()
			}
		case "parent":
			if len(selection) > 0 {
				s = s.ParentFiltered(selection)
			} else {
				s = s.Parent()
			}
		case "parents":
			if len(selection) > 0 {
				s = s.ParentsFiltered(selection)
			} else {
				s = s.Parents()
			}
		case "closest":
			s = s.Closest(selection)
		case "siblings":
			if len(selection) > 0 {
				s = s.SiblingsFiltered(selection)
			} else {
				s = s.Siblings()
			}
		case "next":
			if len(selection) > 0 {
				s = s.NextFiltered(selection)
			} else {
				s = s.Next()
			}
		case "prev":
			if len(selection) > 0 {
				s = s.PrevFiltered(selection)
			} else {
				s = s.Prev()
			}
		case "attr":
			r, _ = s.Attr(selection)
		case "text":
			r = s.Text()
		case "html":
			r, err = s.Html()
		case "outerHTML":
			r, err = goquery.OuterHtml(s)
		case "style":
			r, _ = s.Attr("style")
		case "href":
			r, _ = s.Attr("href")
		case "src":
			r, _ = s.Attr("src")
		case "class":
			r, _ = s.Attr("class")
		case "id":
			r, _ = s.Attr("id")
		case "calc":
			var v interface{}
			v, err = arithmetic.Parse(r)
			if err != nil {
				return
			}
			n, _ := arithmetic.ToFloat(v)
			prec := 2
			if len(selection) > 0 {
				var i64 int64
				i64, err = strconv.ParseInt(selection, 10, 32)
				if err != nil {
					return
				}
				prec = int(i64)
			}
			r = strconv.FormatFloat(n, 'g', prec, 64)
		case "expand":
			var (
				rx *regexp.Regexp
				ry *regexp.Regexp
			)
			ry, err = regexp.Compile(quote + `[\s]*,[\s]*` + quote)
			if err != nil {
				return
			}
			params := ry.Split(selection, 2)
			if len(params) != 2 {
				err = ErrMissingParam
				return
			}
			rx, err = regexp.Compile(params[0])
			if err != nil {
				return
			}
			src := r
			dst := []byte{}
			m := rx.FindStringSubmatchIndex(src)
			s := rx.ExpandString(dst, params[1], src, m)
			r = string(s)
		case "match":
			var rx *regexp.Regexp
			rx, err = regexp.Compile(selection)
			if err != nil {
				return
			}
			rs := rx.FindAllStringSubmatch(r, -1)
			if len(rs) > 0 && len(rs[0]) > 1 {
				r = rs[0][1]
			}
		}
	}
	return
}
