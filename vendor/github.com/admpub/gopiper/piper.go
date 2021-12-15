package gopiper

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/admpub/regexp2"

	"github.com/PuerkitoBio/goquery"
	simplejson "github.com/bitly/go-simplejson"
)

const (
	// begin new version
	PT_RAW          = "raw"
	PT_INT          = "int"
	PT_FLOAT        = "float"
	PT_BOOL         = "bool"
	PT_STRING       = "string"
	PT_INT_ARRAY    = "int-array"
	PT_FLOAT_ARRAY  = "float-array"
	PT_BOOL_ARRAY   = "bool-array"
	PT_STRING_ARRAY = "string-array"
	PT_HTML_ARRAY   = "html-array"
	PT_MAP          = "map"
	PT_ARRAY        = "array"
	PT_JSON_VALUE   = "json"
	PT_JSON_PARSE   = "jsonparse"
	// end new version

	// begin compatible old version
	PT_TEXT       = "text"
	PT_HREF       = "href"
	PT_HTML       = "html"
	PT_ATTR       = `attr\[([\w\W]+)\]`
	PT_ATTR_ARRAY = `attr-array\[([\w\W]+)\]`
	PT_IMG_SRC    = "src"
	PT_IMG_ALT    = "alt"
	PT_TEXT_ARRAY = "text-array"
	PT_HREF_ARRAY = "href-array"
	PT_OUT_HTML   = "outhtml"
	// end compatible old version

	PAGE_JSON = "json"
	PAGE_HTML = "html"
	PAGE_JS   = "js"
	PAGE_XML  = "xml"
	PAGE_TEXT = "text"

	REGEXP_PRE  = "regexp:"
	REGEXP2_PRE = "regexp2:"
)

var (
	attrExp            = regexp.MustCompile(PT_ATTR)
	attrArrayExp       = regexp.MustCompile(PT_ATTR_ARRAY)
	fnExp              = regexp.MustCompile(`([a-z_]+)(\(([\w\W+]+)\))?`)
	jsonNumberIndexExp = regexp.MustCompile(`^\[(\d+)\]$`)
)

// VerifySelector 验证正则表达式
func VerifySelector(selector string) (err error) {
	if strings.HasPrefix(selector, REGEXP_PRE) {
		_, err = regexp.Compile(strings.TrimPrefix(selector, REGEXP_PRE))
	} else if strings.HasPrefix(selector, REGEXP2_PRE) {
		_, err = regexp2.Compile(strings.TrimPrefix(selector, REGEXP2_PRE), 0)
	}
	return
}

type PipeItem struct {
	Name     string     `json:"name,omitempty"` //只有类型为map的时候才会用到
	Selector string     `json:"selector,omitempty"`
	Type     string     `json:"type"`
	Filter   string     `json:"filter,omitempty"`
	SubItem  []PipeItem `json:"subitem,omitempty"`
	fetcher  Fether
	storer   Storer
	pageType string
	doc      *goquery.Document
}

type Fether func(pageURL string) (body []byte, err error)
type Storer func(fileURL, savePath string, fetched bool) (newPath string, err error)

type htmlSelector struct {
	*goquery.Selection
	attr     string
	selector string
}

func (p *PipeItem) SetFetcher(fetcher Fether) {
	p.fetcher = fetcher
}

func (p *PipeItem) SetStorer(storer Storer) {
	p.storer = storer
}

func (p *PipeItem) CopyFrom(from *PipeItem) {
	p.SetFetcher(from.fetcher)
	p.SetStorer(from.storer)
	p.doc = from.doc
}

func (p *PipeItem) Fetcher() Fether {
	return p.fetcher
}

func (p *PipeItem) Storer() Storer {
	return p.storer
}

func (p *PipeItem) PipeBytes(body []byte, pageType string) (interface{}, error) {
	p.pageType = pageType
	switch pageType {
	case PAGE_HTML:
		doc, err := goquery.NewDocumentFromReader(bytes.NewReader(body))
		if err != nil {
			return nil, err
		}
		p.doc = doc
		return p.pipeSelection(doc.Selection)
	case PAGE_JSON:
		return p.pipeJSON(body)
	case PAGE_TEXT:
		return p.pipeText(body)
	}
	return nil, nil
}

func (p *PipeItem) parseRegexp(body string, useRegexp2 bool) (interface{}, error) {
	var (
		preLen int
		sv     []string
		rs     string
	)
	if useRegexp2 {
		preLen = len(REGEXP2_PRE)
	} else {
		preLen = len(REGEXP_PRE)
	}
	s := p.Selector[preLen:]
	if useRegexp2 {
		exp, err := regexp2.Compile(s, regexp2.None)
		if err != nil {
			return nil, err
		}
		mch, err := exp.FindStringMatch(body)
		if err != nil {
			return nil, err
		}
		if mch != nil {
			sv = mch.Slice()
			//fmt.Println(`[regexp2][matched:`+strconv.Itoa(mch.GroupCount())+`]`, mch.String(), com.Dump(sv, false))
		}
	} else {
		exp, err := regexp.Compile(s)
		if err != nil {
			return nil, err
		}
		sv = exp.FindStringSubmatch(body)
	}

	if len(sv) == 1 {
		rs = sv[0]
	} else if len(sv) > 1 {
		rs = sv[1]
		sv = sv[1:]
	}

	switch p.Type {
	case PT_INT, PT_FLOAT, PT_BOOL:
		val, err := parseTextValue(rs, p.Type)
		if err != nil {
			return nil, err
		}
		return callFilter(p, val, p.Filter)
	case PT_INT_ARRAY, PT_FLOAT_ARRAY, PT_BOOL_ARRAY:
		val, err := parseTextValue(sv, p.Type)
		if err != nil {
			return nil, err
		}
		return callFilter(p, val, p.Filter)
	case PT_TEXT, PT_STRING:
		return callFilter(p, rs, p.Filter)
	case PT_TEXT_ARRAY, PT_STRING_ARRAY:
		return callFilter(p, sv, p.Filter)
	case PT_JSON_PARSE:
		if p.SubItem == nil || len(p.SubItem) <= 0 {
			return nil, ErrJsonparseNeedSubItem
		}
		body, err := text2JSONByte(rs)
		if err != nil {
			return nil, errors.New("jsonparse: text is not a json string: " + err.Error())
		}
		parseItem := p.SubItem[0]
		parseItem.CopyFrom(p)
		res, err := parseItem.pipeJSON(body)
		if err != nil {
			return nil, err
		}
		return callFilter(p, res, p.Filter)
	case PT_JSON_VALUE:
		res, err := text2JSON(rs)
		if err != nil {
			return nil, err
		}
		return callFilter(p, res, p.Filter)
	case PT_MAP:
		if p.SubItem == nil || len(p.SubItem) <= 0 {
			return nil, ErrArrayNeedSubItem
		}
		res := make(map[string]interface{})
		for _, subitem := range p.SubItem {
			if len(subitem.Name) == 0 {
				continue
			}
			subitem.CopyFrom(p)
			subitem.Name = replaceName(subitem.Name, res)
			res[subitem.Name], _ = subitem.pipeText([]byte(rs))
		}
		return callFilter(p, res, p.Filter)
	case PT_RAW:
		return callFilter(p, p.Selector, p.Filter)
	}
	return nil, ErrNotSupportPipeType
}

var namePlaceholder = regexp.MustCompile(`#([^#]*)#`)

func replaceName(name string, data map[string]interface{}) string {
	vt := namePlaceholder.FindAllStringSubmatch(name, -1)
	for _, v := range vt {
		if len(v) < 2 {
			continue
		}
		placeholder := v[1]
		var value string
		val, ok := data[placeholder]
		if ok {
			value, ok = val.(string)
			if !ok {
				value = fmt.Sprint(val)
			}
		}
		name = strings.Replace(name, v[0], value, -1)
	}
	return name
}

func (p *PipeItem) pipeSelection(s *goquery.Selection) (interface{}, error) {
	if p.Type == PT_RAW {
		return callFilter(p, p.Selector, p.Filter)
	}
	var (
		sel = htmlSelector{s, "", p.Selector}
		err error
	)
	if strings.HasPrefix(p.Selector, REGEXP_PRE) {
		body, _ := sel.Html()
		return p.parseRegexp(body, false)
	}
	if strings.HasPrefix(p.Selector, REGEXP2_PRE) {
		body, _ := sel.Html()
		return p.parseRegexp(body, true)
	}
	selector := p.Selector
	if len(selector) > 0 {
		sel, err = p.parseHTMLSelector(s, selector)
		if err != nil {
			return nil, err
		}
		selector = sel.selector
	}

	if sel.Size() == 0 {
		return nil, errors.New("Selector can't Find node: " + selector)
	}

	if attrExp.MatchString(p.Type) { // 例如：attr[href] 或 attr[src] 等
		vt := attrExp.FindStringSubmatch(p.Type)
		res, has := sel.Attr(vt[1])
		if !has {
			return nil, errors.New("Can't Find attribute: " + p.Type + " selector: " + selector)
		}
		return callFilter(p, res, p.Filter)
	}
	if attrArrayExp.MatchString(p.Type) { // 例如：attr-array[href] 或 attr-array[src] 等
		vt := attrArrayExp.FindStringSubmatch(p.Type)
		res := make([]string, 0)
		sel.Each(func(index int, child *goquery.Selection) {
			href, has := child.Attr(vt[1])
			if has {
				res = append(res, href)
			}
		})
		return callFilter(p, res, p.Filter)
	}

	switch p.Type {
	case PT_INT, PT_FLOAT, PT_BOOL, PT_STRING, PT_TEXT, PT_INT_ARRAY, PT_FLOAT_ARRAY, PT_BOOL_ARRAY, PT_STRING_ARRAY:
		val, err := parseHTMLAttr(sel, p.Type)
		if err != nil {
			return nil, err
		}
		return callFilter(p, val, p.Filter)
	case PT_HTML_ARRAY:
		res := make([]string, 0)
		sel.Each(func(index int, child *goquery.Selection) {
			str, _ := child.Html()
			res = append(res, str)
		})
		return callFilter(p, res, p.Filter)
	case PT_HTML:
		var html string
		sel.Each(func(idx int, child *goquery.Selection) {
			str, _ := child.Html()
			html += str
		})
		return callFilter(p, html, p.Filter)
	case PT_OUT_HTML:
		var html string
		sel.Each(func(idx int, child *goquery.Selection) {
			str, _ := goquery.OuterHtml(child)
			html += str
		})
		return callFilter(p, html, p.Filter)
	case PT_HREF, PT_IMG_SRC, PT_IMG_ALT:
		res, has := sel.Attr(p.Type)
		if !has {
			return nil, errors.New("Can't Find attribute: " + p.Type + " selector: " + selector)
		}
		return callFilter(p, res, p.Filter)
	case PT_TEXT_ARRAY:
		res := make([]string, 0)
		sel.Each(func(index int, child *goquery.Selection) {
			res = append(res, child.Text())
		})
		return callFilter(p, res, p.Filter)
	case PT_HREF_ARRAY:
		res := make([]string, 0)
		sel.Each(func(index int, child *goquery.Selection) {
			href, has := child.Attr("href")
			if has {
				res = append(res, href)
			}
		})
		return callFilter(p, res, p.Filter)
	case PT_ARRAY:
		if p.SubItem == nil || len(p.SubItem) <= 0 {
			return nil, ErrArrayNeedSubItem
		}
		arrayItem := p.SubItem[0]
		arrayItem.CopyFrom(p)
		res := make([]interface{}, 0)
		sel.Each(func(index int, child *goquery.Selection) {
			v, _ := arrayItem.pipeSelection(child)
			res = append(res, v)
		})
		return callFilter(p, res, p.Filter)
	case PT_MAP:
		if p.SubItem == nil || len(p.SubItem) <= 0 {
			return nil, ErrArrayNeedSubItem
		}
		res := make(map[string]interface{})
		for _, subitem := range p.SubItem {
			if len(subitem.Name) == 0 {
				continue
			}
			subitem.CopyFrom(p)
			subitem.Name = replaceName(subitem.Name, res)
			res[subitem.Name], _ = subitem.pipeSelection(sel.Selection)
		}

		return callFilter(p, res, p.Filter)
	default:
		return callFilter(p, 0, p.Filter)
	}
}

func (p *PipeItem) CallFilter(src interface{}, filters string) (interface{}, error) {
	return callFilter(p, src, filters)
}

func (p *PipeItem) parseHTMLSelector(s *goquery.Selection, selector string) (htmlSelector, error) {
	var attr string
	if len(selector) == 0 {
		return htmlSelector{s, attr, selector}, nil
	}
	if strings.HasPrefix(selector, `$.`) {
		selector = strings.TrimPrefix(selector, `$`)
		s = p.doc.Selection
	}
	// html: <a class="bn-sharing" data-type="book"></a>
	// selector: a.bn-sharing//attr[data-type]
	if idx := strings.Index(selector, "//"); idx > 0 {
		attr = strings.TrimSpace(selector[idx+2:])
		selector = strings.TrimSpace(selector[:idx])
	}

	// selector: ul > li | eq(2)
	subs := SplitParams(selector, `|`)
	leng := len(subs)
	if leng < 2 {
		return htmlSelector{s.Find(selector), attr, selector}, nil
	}
	subs[0] = strings.TrimSpace(subs[0])
	s = s.Find(subs[0])
	for i := 1; i < leng; i++ {
		subs[i] = strings.TrimSpace(subs[i])
		if !fnExp.MatchString(subs[i]) {
			return htmlSelector{s, attr, selector}, errors.New("error parse html selector: " + subs[i])
		}

		vt := fnExp.FindStringSubmatch(subs[i])
		fn := vt[1]
		var params string
		if len(vt) > 3 {
			params = strings.TrimSpace(vt[3])
		}

		switch fn {
		case "eq":
			pm, _ := strconv.Atoi(params)
			s = s.Eq(pm)
		case "next":
			if len(params) > 0 {
				s = s.NextFiltered(params)
			} else {
				s = s.Next()
			}
		case "prev":
			if len(params) > 0 {
				s = s.PrevFiltered(params)
			} else {
				s = s.Prev()
			}
		case "first":
			s = s.First()
		case "last":
			s = s.Last()
		case "siblings":
			if len(params) > 0 {
				s = s.SiblingsFiltered(params)
			} else {
				s = s.Siblings()
			}
		case "nextall":
			if len(params) > 0 {
				s = s.NextAllFiltered(params)
			} else {
				s = s.NextAll()
			}
		case "children":
			if len(params) > 0 {
				s = s.ChildrenFiltered(params)
			} else {
				s = s.Children()
			}
		case "parent":
			if len(params) > 0 {
				s = s.ParentFiltered(params)
			} else {
				s = s.Parent()
			}
		case "parents":
			if len(params) > 0 {
				s = s.ParentsFiltered(params)
			} else {
				s = s.Parents()
			}
		case "not":
			if len(params) > 0 {
				s = s.Not(params)
			}
		case "filter":
			if len(params) > 0 {
				s = s.Filter(params)
			}
		case "prevall":
			if len(params) > 0 {
				s = s.PrevAllFiltered(params)
			} else {
				s = s.PrevAll()
			}
		case "rm", "remove":
			if len(params) > 0 {
				s.Find(params).Remove()
			}
		case "attr":
			if len(params) > 0 {
				return htmlSelector{s, `attr[` + params + `]`, selector}, nil
			}
		}
	}
	return htmlSelector{s, attr, selector}, nil
}

func parseTextValue(text interface{}, tp string) (interface{}, error) {
	switch tp {
	case PT_INT, PT_INT_ARRAY:
		return text2int(text)
	case PT_FLOAT, PT_FLOAT_ARRAY:
		return text2float(text)
	case PT_BOOL, PT_BOOL_ARRAY:
		return text2bool(text)
	}
	return text, nil
}

func parseHTMLAttr(sel htmlSelector, tp string) (interface{}, error) {
	switch tp {
	case PT_INT, PT_FLOAT, PT_BOOL, PT_TEXT, PT_STRING:
		text, err := getHTMLAttr(sel.Selection, sel.attr, sel.selector)
		if err != nil {
			return nil, err
		}
		return parseTextValue(text, tp)
	case PT_INT_ARRAY, PT_FLOAT_ARRAY, PT_BOOL_ARRAY, PT_STRING_ARRAY:
		text, err := getHTMLAttrArray(sel.Selection, sel.attr, sel.selector)
		if err != nil {
			return nil, err
		}
		return parseTextValue(text, tp)
	}

	return nil, ErrUnknowHTMLAttr
}

func getHTMLAttrArray(sel *goquery.Selection, attr, selector string) ([]string, error) {
	res := make([]string, 0)
	if len(attr) == 0 {
		sel.Each(func(index int, child *goquery.Selection) {
			res = append(res, child.Text())
		})
		return res, nil
	}

	if attrExp.MatchString(attr) {
		vt := attrExp.FindStringSubmatch(attr)
		sel.Each(func(index int, child *goquery.Selection) {
			text, has := child.Attr(vt[1])
			if has {
				res = append(res, text)
			}
		})
		return res, nil
	}
	if attr == "html" {
		sel.Each(func(idx int, s1 *goquery.Selection) {
			str, _ := s1.Html()
			res = append(res, str)
		})
		return res, nil
	}
	if attr == "outhtml" {
		sel.Each(func(idx int, s1 *goquery.Selection) {
			str, _ := goquery.OuterHtml(s1)
			res = append(res, str)
		})
		return res, nil
	}

	return res, nil
}

func getHTMLAttr(sel *goquery.Selection, attr, selector string) (string, error) {
	if len(attr) == 0 {
		return sel.Text(), nil
	}

	if attrExp.MatchString(attr) {
		vt := attrExp.FindStringSubmatch(attr)
		res, has := sel.Attr(vt[1])
		if !has {
			return "", errors.New("Can't Find attribute: " + attr + " selector: " + selector)
		}
		return res, nil
	}
	if attr == "html" {
		var html string
		sel.Each(func(idx int, s1 *goquery.Selection) {
			str, _ := s1.Html()
			html += str
		})
		return html, nil
	}
	if attr == "outhtml" {
		var html string
		sel.Each(func(idx int, s1 *goquery.Selection) {
			str, _ := goquery.OuterHtml(s1)
			html += str
		})
		return html, nil
	}

	return sel.Text(), nil
}

func parseJSONSelector(js *simplejson.Json, selector string) (*simplejson.Json, error) {
	subs := strings.Split(selector, ".")
	for _, s := range subs {
		if index := strings.Index(s, "["); index >= 0 {
			if index > 0 {
				k := s[:index]
				if k != "this" {
					js = js.Get(k)
				}
			}
			s = s[index:]
			if !jsonNumberIndexExp.MatchString(s) {
				return nil, errors.New("parse json selector error:  " + selector)
			}
			v := jsonNumberIndexExp.FindStringSubmatch(s)
			intV, err := strconv.Atoi(v[1])
			if err != nil {
				return nil, err
			}
			js = js.GetIndex(intV)
		} else {
			if s == "this" {
				continue
			}
			js = js.Get(s)
		}
	}
	return js, nil
}

func (p *PipeItem) pipeJSON(body []byte) (interface{}, error) {
	if p.Type == PT_RAW {
		return callFilter(p, p.Selector, p.Filter)
	}
	js, err := simplejson.NewJson(body)
	if err != nil {
		return nil, err
	}

	if len(p.Selector) > 0 {
		js, err = parseJSONSelector(js, p.Selector)
		if err != nil {
			return nil, err
		}
	}

	switch p.Type {
	case PT_INT:
		return callFilter(p, js.MustInt64(0), p.Filter)
	case PT_FLOAT:
		return callFilter(p, js.MustFloat64(0.0), p.Filter)
	case PT_BOOL:
		return callFilter(p, js.MustBool(false), p.Filter)
	case PT_TEXT, PT_STRING:
		return callFilter(p, js.MustString(""), p.Filter)
	case PT_TEXT_ARRAY, PT_STRING_ARRAY:
		v, err := js.StringArray()
		if err != nil {
			return nil, err
		}
		return callFilter(p, v, p.Filter)
	case PT_JSON_VALUE:
		return callFilter(p, js.Interface(), p.Filter)
	case PT_JSON_PARSE:
		if p.SubItem == nil || len(p.SubItem) <= 0 {
			return nil, ErrJsonparseNeedSubItem
		}
		bodyStr := strings.TrimSpace(js.MustString(""))
		if len(bodyStr) == 0 {
			return nil, nil
		}
		body, err := text2JSONByte(bodyStr)
		if err != nil {
			return nil, errors.New("jsonparse: text is not a json string: " + err.Error())
		}
		parseItem := p.SubItem[0]
		parseItem.CopyFrom(p)
		res, err := parseItem.pipeJSON(body)
		if err != nil {
			return nil, err
		}
		return callFilter(p, res, p.Filter)
	case PT_ARRAY:
		v, err := js.Array()
		if err != nil {
			return nil, err
		}

		if p.SubItem == nil || len(p.SubItem) <= 0 {
			return nil, ErrArrayNeedSubItem
		}
		arrayItem := p.SubItem[0]
		arrayItem.CopyFrom(p)
		res := make([]interface{}, 0)
		for _, r := range v {
			data, _ := json.Marshal(r)
			vl, _ := arrayItem.pipeJSON(data)
			res = append(res, vl)
		}
		return callFilter(p, res, p.Filter)
	case PT_MAP:
		if p.SubItem == nil || len(p.SubItem) <= 0 {
			return nil, ErrArrayNeedSubItem
		}
		data, _ := json.Marshal(js)
		res := make(map[string]interface{})
		for _, subitem := range p.SubItem {
			if len(subitem.Name) == 0 {
				continue
			}
			subitem.CopyFrom(p)
			subitem.Name = replaceName(subitem.Name, res)
			res[subitem.Name], _ = subitem.pipeJSON(data)
		}

		return callFilter(p, res, p.Filter)
	default:
		return callFilter(p, 0, p.Filter)
	}
}

func (p *PipeItem) pipeText(body []byte) (interface{}, error) {
	if p.Type == PT_RAW {
		return callFilter(p, p.Selector, p.Filter)
	}
	bodyStr := string(body)
	if strings.HasPrefix(p.Selector, REGEXP_PRE) {
		return p.parseRegexp(bodyStr, false)
	}
	if strings.HasPrefix(p.Selector, REGEXP2_PRE) {
		return p.parseRegexp(bodyStr, true)
	}

	switch p.Type {
	case PT_INT, PT_FLOAT, PT_BOOL:
		val, err := parseTextValue(bodyStr, p.Type)
		if err != nil {
			return nil, err
		}
		return callFilter(p, val, p.Filter)
	case PT_TEXT, PT_STRING:
		return callFilter(p, bodyStr, p.Filter)
	case PT_JSON_PARSE:
		if p.SubItem == nil || len(p.SubItem) <= 0 {
			return nil, ErrJsonparseNeedSubItem
		}
		body, err := text2JSONByte(bodyStr)
		if err != nil {
			return nil, errors.New("jsonparse: text is not a json string: " + err.Error())
		}
		parseItem := p.SubItem[0]
		parseItem.CopyFrom(p)
		res, err := parseItem.pipeJSON(body)
		if err != nil {
			return nil, err
		}
		return callFilter(p, res, p.Filter)
	case PT_JSON_VALUE:
		res, err := text2JSON(string(body))
		if err != nil {
			return nil, err
		}
		return callFilter(p, res, p.Filter)
	case PT_MAP:
		if p.SubItem == nil || len(p.SubItem) <= 0 {
			return nil, ErrArrayNeedSubItem
		}
		res := make(map[string]interface{})
		for _, subitem := range p.SubItem {
			if len(subitem.Name) == 0 {
				continue
			}
			subitem.Name = replaceName(subitem.Name, res)
			res[subitem.Name], _ = subitem.pipeText(body)
		}
		return callFilter(p, res, p.Filter)
	default:
		return callFilter(p, 0, p.Filter)
	}
}

func text2int(text interface{}) (interface{}, error) {
	switch val := text.(type) {
	case string:
		return strconv.ParseInt(val, 10, 64)
	case []string:
		vs := make([]int64, 0)
		for _, v := range val {
			n, err := strconv.ParseInt(v, 10, 64)
			if err != nil {
				return nil, err
			}
			vs = append(vs, n)
		}
		return vs, nil
	}
	return nil, ErrUnsupportText2intType
}

func text2float(text interface{}) (interface{}, error) {
	switch val := text.(type) {
	case string:
		return strconv.ParseFloat(val, 64)
	case []string:
		vs := make([]float64, 0)
		for _, v := range val {
			n, err := strconv.ParseFloat(v, 64)
			if err != nil {
				return nil, err
			}
			vs = append(vs, n)
		}
		return vs, nil
	}
	return nil, ErrUnsupportText2floatType
}

func text2bool(text interface{}) (interface{}, error) {
	switch val := text.(type) {
	case string:
		return strconv.ParseBool(val)
	case []string:
		vs := make([]bool, 0)
		for _, v := range val {
			n, err := strconv.ParseBool(v)
			if err != nil {
				return nil, err
			}
			vs = append(vs, n)
		}
		return vs, nil
	}
	return nil, ErrUnsupportText2boolType
}

func text2JSON(text string) (interface{}, error) {
	res, err := textJSONValue(text)
	if err != nil {
		return untextJSONValue(text)
	}
	return res, nil
}

func text2JSONByte(text string) ([]byte, error) {
	val, err := text2JSON(text)
	if err != nil {
		return nil, err
	}
	return json.Marshal(val)
}

func textJSONValue(text string) (interface{}, error) {
	res := map[string]interface{}{}
	if err := json.Unmarshal([]byte(text), &res); err != nil {
		resarray := make([]interface{}, 0)
		if err = json.Unmarshal([]byte(text), &resarray); err != nil {
			return nil, errors.New("parse json value error, text is not json value: " + err.Error())
		}
		return resarray, nil
	}

	return res, nil
}

func untextJSONValue(text string) (interface{}, error) {
	text, err := strconv.Unquote(`"` + text + `"`)
	if err != nil {
		return nil, err
	}
	return textJSONValue(text)
}
