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

type Selection struct {
	Function   string
	Parameters []string
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
		value, err = g.parseSelections(doc.Selection, rule)
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

func (g *GoQuery) parseSelections(rootSelection *goquery.Selection, rule string) (r string, err error) {
	s := rootSelection
	for _, selector := range parseSelections(rule) {
		switch selector.Function {
		case "$":
			s = rootSelection
			if len(selector.Parameters) > 0 && len(selector.Parameters[0]) > 0 {
				s = s.Find(selector.Parameters[0])
			}
		case "find":
			if len(selector.Parameters) > 0 && len(selector.Parameters[0]) > 0 {
				s = s.Find(selector.Parameters[0])
			}
		case "children":
			if len(selector.Parameters) > 0 && len(selector.Parameters[0]) > 0 {
				s = s.ChildrenFiltered(selector.Parameters[0])
			} else {
				s = s.Children()
			}
		case "parent":
			if len(selector.Parameters) > 0 && len(selector.Parameters[0]) > 0 {
				s = s.ParentFiltered(selector.Parameters[0])
			} else {
				s = s.Parent()
			}
		case "parents":
			if len(selector.Parameters) > 0 && len(selector.Parameters[0]) > 0 {
				s = s.ParentsFiltered(selector.Parameters[0])
			} else {
				s = s.Parents()
			}
		case "closest":
			s = s.Closest(selector.Parameters[0])
		case "siblings":
			if len(selector.Parameters) > 0 && len(selector.Parameters[0]) > 0 {
				s = s.SiblingsFiltered(selector.Parameters[0])
			} else {
				s = s.Siblings()
			}
		case "next":
			if len(selector.Parameters) > 0 && len(selector.Parameters[0]) > 0 {
				s = s.NextFiltered(selector.Parameters[0])
			} else {
				s = s.Next()
			}
		case "prev":
			if len(selector.Parameters) > 0 && len(selector.Parameters[0]) > 0 {
				s = s.PrevFiltered(selector.Parameters[0])
			} else {
				s = s.Prev()
			}
		case "attr":
			if len(selector.Parameters) > 0 && len(selector.Parameters[0]) > 0 {
				r, _ = s.Attr(selector.Parameters[0])
			}
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
			if len(selector.Parameters) > 0 && len(selector.Parameters[0]) > 0 {
				var i64 int64
				i64, err = strconv.ParseInt(selector.Parameters[0], 10, 32)
				if err != nil {
					return
				}
				prec = int(i64)
			}
			r = strconv.FormatFloat(n, 'g', prec, 64)
		case "expand":
			var (
				rx *regexp.Regexp
			)
			if len(selector.Parameters) != 2 {
				err = ErrMissingParam
				return
			}
			rx, err = regexp.Compile(selector.Parameters[0])
			if err != nil {
				return
			}
			src := r
			dst := []byte{}
			m := rx.FindStringSubmatchIndex(src)
			s := rx.ExpandString(dst, selector.Parameters[1], src, m)
			r = string(s)
		case "match":
			if len(selector.Parameters) > 0 && len(selector.Parameters[0]) > 0 {
				var rx *regexp.Regexp
				rx, err = regexp.Compile(selector.Parameters[0])
				if err != nil {
					return
				}
				rs := rx.FindAllStringSubmatch(r, -1)
				if len(rs) > 0 && len(rs[0]) > 1 {
					r = rs[0][1]
				}
			}
		}
	}
	return
}

func parseSelections(rule string) []*Selection {
	selections := []*Selection{}
	//$('#ID').find('.class').text
	var (
		function     string
		parameters   []string
		paramItem    []rune
		paramStarted bool
		quoteStarted bool
		quote        rune
		slashAdded   bool
	)
	for index, v := range rule {
		if index == 0 && v == '$' {
			function = "$"
			continue
		}
		if !paramStarted {
			if v == '(' {
				paramStarted = true
				continue
			}
			if v == '.' {
				continue
			}
			function += string(v)
			continue
		}
		if !quoteStarted {
			if v == ')' {
				if paramItem != nil {
					parameters = append(parameters, string(paramItem))
				}
				selections = append(selections, &Selection{
					Function:   function,
					Parameters: parameters,
				})
				function = ``
				paramItem = nil
				paramStarted = false
				parameters = []string{}
				slashAdded = false
				continue
			}
			if v == '\'' || v == '"' {
				quote = v
				quoteStarted = true
			}
			continue
		}
		if !slashAdded {
			if v == '\\' {
				slashAdded = true
				continue
			}
			if quote == v {
				if paramItem != nil {
					parameters = append(parameters, string(paramItem))
				}
				paramItem = nil
				quoteStarted = false
				slashAdded = false
				continue
			}
		}
		paramItem = append(paramItem, v)
		slashAdded = false
	}
	if len(function) > 0 {
		if paramItem != nil {
			parameters = append(parameters, string(paramItem))
		}
		selections = append(selections, &Selection{
			Function:   function,
			Parameters: parameters,
		})
	}
	return selections
}
