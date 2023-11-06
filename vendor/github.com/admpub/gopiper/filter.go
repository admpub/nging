package gopiper

import (
	"errors"
	"fmt"
	"html"
	"path"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/admpub/regexp2"
	"github.com/webx-top/com"
)

func init() {
	RegisterFilter("preadd", preadd, "添加前缀", `preadd(prefix)`, ``)
	RegisterFilter("postadd", postadd, "添加后缀", `postadd(suffix)`, ``)
	RegisterFilter("replace", replace, "替换", `replace(find,replace)`, ``)
	RegisterFilter("split", split, "将字符串按指定分隔符分割成数组", `split(-)`, ``)
	RegisterFilter("join", join, "合并数组为字符串", `join(-)`, ``)
	RegisterFilter("trim", trim, "剪掉头尾指定字符", `trim(;)`, ``)
	RegisterFilter("trimleft", trimleft, "从左边剪掉指定字符串", `trimleft(a-)`, ``)
	RegisterFilter("trimright", trimright, "从右边剪掉指定字符串", `trimright(.html)`, ``)
	RegisterFilter("trimspace", trimspace, "剪掉头尾空白", `trimspace`, ``)
	RegisterFilter("substr", substr, "获取子字符串。字符串总是从左向右从0开始编号，参数1和参数2分别用来指定要截取的起止位置编号，截取子字符串时，总是包含起始编号的字符，不包含终止编号的字符", `substr(0,5)`, ``)
	RegisterFilter("intval", intval, "转换为整数", `intval`, ``)
	RegisterFilter("floatval", floatval, "转换为小数", `floatval`, ``)
	RegisterFilter("hrefreplace", hrefreplace, "替换href属性。$2为捕获到的href属性值", `hrefreplace(data-url="$2")`, ``)
	RegisterFilter("regexpreplace", regexpreplace, "正则替换(regexp2引擎)。参数1为正则表达式，参数2为替换成的新内容，参数3为起始位置编号(从0开始)，参数4为替换次数(-1代表相对全部替换,-2代表绝对全部替换)", `regexpreplace(^A$,B,0,-1)`, ``)
	RegisterFilter("wraphtml", wraphtml, "将采集到的数据用HTML标签包围起来", `wraphtml(a)`, ``)
	RegisterFilter("tosbc", tosbc, "将全角的标点符号和英文字母转换为半角", `tosbc`, ``)
	RegisterFilter("unescape", unescape, "解码HTML", `unescape`, ``)
	RegisterFilter("escape", escape, "编码HTML", `escape`, ``)
	RegisterFilter("sprintf", sprintf, "格式化", `sprintf(%s)`, ``)
	RegisterFilter("sprintfmap", sprintfmap, "用map值格式化(前提是采集到的数据必须是map类型)。参数1为模板字符串，其它参数用于指定相应map元素值的键值", `sprintfmap(%v-%v,a,b)`, ``)
	RegisterFilter("unixtime", unixtime, "UNIX时间戳(秒)。如果带参数则代表将获取到的数据按照参数指定的格式转为时间戳；不带参数则获取当前时间戳", `unixtime(DateTime)`, `unixtime、unixtime(Y-m-d H:i:s)、unixtime(DateTime) 或 unixtime(2006-01-02 15:04:05)`)
	RegisterFilter("unixmill", unixmill, "获取当前UNIX时间戳(毫秒)", `unixmill`, ``)
	RegisterFilter("paging", paging, "分页。参数1为起始页码，参数2为终止页码，参数3为步进值(可选)。需要在网址中添加页码占位符“{0}”，一般与sprintf组合起来使用。经过paging处理后的网址会变成网址数组", `paging(1,10,1)`, `sprintf(%s?page={0})|paging(1,10)`)
	RegisterFilter("quote", quote, "用双引号包起来", `quote`, ``)
	RegisterFilter("unquote", unquote, "取消双引号包围", `unquote`, ``)
	RegisterFilter("saveto", saveto, "下载并保存文件到指定位置", `saveto(savePath)`, ``)
	RegisterFilter("fetch", fetch, "抓取网址内容。参数pageType仅支持html、json、text这三个值", `fetch(pageType,selector)`, ``)
	RegisterFilter("basename", basename, "获取文件名", `basename`, ``)
	RegisterFilter("extension", extension, "获取扩展名", `extension`, ``)
}

type FilterFunction func(pipe *PipeItem, src interface{}, params string) (interface{}, error)

func NewFilter(name string, fn FilterFunction, description, usage, example string) *Filter {
	return &Filter{
		Name:        name,
		function:    fn,
		Description: description,
		Usage:       usage,
		Example:     example,
	}
}

type Filter struct {
	Name        string
	function    FilterFunction
	Description string `json:",omitempty"`
	Usage       string `json:",omitempty"`
	Example     string `json:",omitempty"`
}

var filters = make(map[string]*Filter)

func RegisterFilter(name string, fn FilterFunction, description, usage, example string) {
	_, existing := filters[name]
	if existing {
		panic(fmt.Sprintf("Filter with name '%s' is already registered.", name))
	}
	filters[name] = NewFilter(name, fn, description, usage, example)
}

func ReplaceFilter(name string, fn FilterFunction, description, usage, example string) {
	_, existing := filters[name]
	if !existing {
		panic(fmt.Sprintf("Filter with name '%s' does not exist (therefore cannot be overridden).", name))
	}
	filters[name] = NewFilter(name, fn, description, usage, example)
}

func AllFilter() map[string]*Filter {
	return filters
}

var (
	filterExp      = regexp.MustCompile(`([a-zA-Z0-9\-_]+)(?:\(([\w\W]*?)\))?(\||$)`)
	hrefFilterExp  = regexp.MustCompile(`href(?:\s*)=(?:\s*)(['"])?([^'" ]*)(['"])?`)
	hrefFilterExp2 = regexp2.MustCompile(`href(?:\s*)=(?:\s*)(['"]?)([^'" ]*)\1`, regexp2.IgnoreCase)
)

func applyFilter(pipe *PipeItem, name string, src interface{}, params string) (interface{}, error) {
	filter, existing := filters[name]
	if !existing {
		return nil, fmt.Errorf("Filter with name '%s' not found", name)
	}
	return filter.function(pipe, src, params)
}

func callFilter(pipe *PipeItem, src interface{}, value string) (interface{}, error) {

	if src == nil || len(value) == 0 {
		return src, nil
	}

	vt := filterExp.FindAllStringSubmatch(value, -1)

	for _, v := range vt {
		if len(v) < 3 {
			continue
		}
		name := v[1]
		params := v[2]
		next, err := applyFilter(pipe, name, src, params)
		if err != nil {
			if err == ErrInvalidContent {
				return next, err
			}
			continue
		}
		src = next
	}

	return src, nil
}

// fetch(pageType,selector)
func fetch(pipe *PipeItem, src interface{}, params string) (interface{}, error) {
	if pipe.fetcher == nil {
		return src, ErrFetcherNotRegistered
	}
	var (
		pageType = pipe.pageType
		selector string
	)
	paramList := SplitParams(params, `,`)
	switch len(paramList) {
	case 2:
		selector = paramList[1]
		fallthrough
	case 1:
		pageType = paramList[0]
	}
	return _filterValue(src, func(v string) (interface{}, error) {
		body, err := pipe.fetcher(v)
		if err != nil {
			return nil, err
		}
		if len(selector) == 0 {
			return string(body), nil
		}
		pipe2 := &PipeItem{
			Name:     ``,
			Selector: selector,
			Type:     PT_STRING,
			Filter:   ``,
		}
		return pipe2.PipeBytes(body, pageType)
	})
}

// saveto(savePath)
func saveto(pipe *PipeItem, src interface{}, params string) (interface{}, error) {
	if pipe.storer == nil {
		return src, ErrStorerNotRegistered
	}
	var (
		fetched  bool
		savePath string
	)
	paramList := SplitParams(params, `,`)
	switch len(paramList) {
	case 2:
		fetched, _ = strconv.ParseBool(strings.TrimSpace(paramList[1]))
		fallthrough
	case 1:
		savePath = strings.TrimSpace(paramList[0])
	}
	return _filterValue(src, func(v string) (interface{}, error) {
		return pipe.storer(v, savePath, fetched)
	})
}

// preadd(prefix) => {prefix}{src}
func preadd(pipe *PipeItem, src interface{}, params string) (interface{}, error) {
	return _filterValue(src, func(v string) (interface{}, error) {
		return params + v, nil
	}, func(_ interface{}) (interface{}, error) {
		return params, nil
	})
}

// postadd(suffix) => {src}{suffix}
func postadd(pipe *PipeItem, src interface{}, params string) (interface{}, error) {
	return _filterValue(src, func(v string) (interface{}, error) {
		return v + params, nil
	}, func(_ interface{}) (interface{}, error) {
		return params, nil
	})
}

func _substr(src string, params string) string {
	vt := strings.Split(params, ",")
	switch len(vt) {
	case 1:
		start, _ := strconv.Atoi(vt[0])
		return src[start:]
	case 2:
		start, _ := strconv.Atoi(vt[0])
		end, _ := strconv.Atoi(vt[1])
		return src[start:end]
	}
	return src
}

// substr(0,5) => src[0:5]
// substr(5) => src[5:]
func substr(pipe *PipeItem, src interface{}, params string) (interface{}, error) {
	return _filterValue(src, func(v string) (interface{}, error) {
		return _substr(v, params), nil
	})
}

func _replace(src string, params string) string {
	vt := SplitParams(params)
	switch len(vt) {
	case 1:
		return strings.Replace(src, vt[0], "", -1)
	case 2:
		return strings.Replace(src, vt[0], vt[1], -1)
	case 3:
		n, _ := strconv.Atoi(vt[2])
		return strings.Replace(src, vt[0], vt[1], n)
	}
	return src
}

// replace(find,replace) => src=findaaa => replaceaaa
// replace(find) => src=findaaa => aaa
func replace(pipe *PipeItem, src interface{}, params string) (interface{}, error) {
	return _filterValue(src, func(v string) (interface{}, error) {
		return _replace(v, params), nil
	})
}

// trim(;) => src=;a; => a
func trim(pipe *PipeItem, src interface{}, params string) (interface{}, error) {
	if len(params) == 0 {
		return src, ErrTrimNilParams
	}
	return _filterValue(src, func(v string) (interface{}, error) {
		return strings.Trim(v, params), nil
	})
}
func trimleft(pipe *PipeItem, src interface{}, params string) (interface{}, error) {
	if len(params) == 0 {
		return src, ErrTrimNilParams
	}
	return _filterValue(src, func(v string) (interface{}, error) {
		return strings.TrimLeft(v, params), nil
	})
}
func trimright(pipe *PipeItem, src interface{}, params string) (interface{}, error) {
	if len(params) == 0 {
		return src, ErrTrimNilParams
	}
	return _filterValue(src, func(v string) (interface{}, error) {
		return strings.TrimRight(v, params), nil
	})
}

// trimspace => src=" \naaa\n " => "aaa"
func trimspace(pipe *PipeItem, src interface{}, params string) (interface{}, error) {
	return _filterValue(src, func(v string) (interface{}, error) {
		return strings.TrimSpace(v), nil
	})
}

// split(:) => src="a:b" => [a,b]
func split(pipe *PipeItem, src interface{}, params string) (interface{}, error) {
	if len(params) == 0 {
		return src, ErrSplitNilParams
	}
	return _filterValue(src, func(v string) (interface{}, error) {
		str := strings.TrimSpace(v)
		if len(str) == 0 {
			return []string{}, nil
		}
		return strings.Split(str, params), nil
	})
}

// join(:) => src=["a","b"] => a:b
func join(pipe *PipeItem, src interface{}, params string) (interface{}, error) {
	if len(params) == 0 {
		return src, ErrJoinNilParams
	}
	switch vt := src.(type) {
	case []string:
		rs := make([]string, 0, len(vt))
		for _, v := range vt {
			if len(v) > 0 {
				rs = append(rs, v)
			}
		}
		return strings.Join(rs, params), nil
	default:
		return vt, nil
	}
}

// intval => src="123" => 123
func intval(pipe *PipeItem, src interface{}, params string) (interface{}, error) {
	return _filterValue(src, func(v string) (interface{}, error) {
		return strconv.Atoi(v)
	})
}

// basename => src="a/b/c.html" => c.html
func basename(pipe *PipeItem, src interface{}, params string) (interface{}, error) {
	return _filterValue(src, func(v string) (interface{}, error) {
		return path.Base(v), nil
	})
}

// extension => src="a/b/c.html" => .html
func extension(pipe *PipeItem, src interface{}, params string) (interface{}, error) {
	return _filterValue(src, func(v string) (interface{}, error) {
		return path.Ext(v), nil
	})
}

// floatval => src="12.3" => 12.3
func floatval(pipe *PipeItem, src interface{}, params string) (interface{}, error) {
	return _filterValue(src, func(v string) (interface{}, error) {
		return strconv.ParseFloat(v, 64)
	})
}

// hrefreplace(data-url="$2") => src=`href="http://www.admpub.com"` => data-url="http://www.admpub.com"
func hrefreplace(pipe *PipeItem, src interface{}, params string) (interface{}, error) {
	return _filterValue(src, func(v string) (interface{}, error) {
		return hrefFilterExp2.Replace(v, params, 0, -1)
		//return hrefFilterExp.ReplaceAllString(v, params), nil
	})
}

// regexpreplace(^1) => src="1233" => "233"
// regexpreplace(^1,2) => src="1233" => "2233"
func regexpreplace(pipe *PipeItem, src interface{}, params string) (interface{}, error) {
	vt := SplitParams(params)
	var (
		expr    string
		repl    string
		startAt int
		count   = -1
	)
	switch len(vt) {
	case 4:
		count, _ = strconv.Atoi(vt[3])
		fallthrough
	case 3:
		startAt, _ = strconv.Atoi(vt[2])
		fallthrough
	case 2:
		repl = vt[1]
		fallthrough
	case 1:
		expr = vt[0]
	}
	re, err := regexp2.Compile(expr, 0)
	if err != nil {
		return src, err
	}
	return _filterValue(src, func(v string) (interface{}, error) {
		if count < -1 {
			find, err := re.MatchString(v)
			for find {
				v, err = re.Replace(v, repl, startAt, -1)
				if err != nil {
					return v, err
				}
				if len(v) == 0 {
					break
				}
				find, err = re.MatchString(v)
			}
			return v, err
		}
		return re.Replace(v, repl, startAt, count)
	})
}

// 将全角的标点符号和英文字母转换为半角
func _tosbc(src string) string {
	var res string
	for _, t := range src {
		if t == 12288 {
			t = 32
		} else if t > 65280 && t < 65375 {
			t = t - 65248
		}
		res += string(t)
	}
	return res
}

// tosbc => src="1～2" => "1~2"
func tosbc(pipe *PipeItem, src interface{}, params string) (interface{}, error) {
	return _filterValue(src, func(v string) (interface{}, error) {
		return _tosbc(v), nil
	})
}

// unescape => src="&lt;" => "<"
func unescape(pipe *PipeItem, src interface{}, params string) (interface{}, error) {
	return _filterValue(src, func(v string) (interface{}, error) {
		return html.UnescapeString(v), nil
	})
}

// escape => src="<" => "&lt;"
func escape(pipe *PipeItem, src interface{}, params string) (interface{}, error) {
	return _filterValue(src, func(v string) (interface{}, error) {
		return html.EscapeString(v), nil
	})
}

// wraphtml(a) => <a>{src}</a>
func wraphtml(pipe *PipeItem, src interface{}, params string) (interface{}, error) {
	if len(params) == 0 {
		return src, errors.New("filter wraphtml nil params")
	}

	return _filterValue(src, func(v string) (interface{}, error) {
		return fmt.Sprintf("<%s>%s</%s>", params, v, params), nil
	})
}

// sprintf_multi_param(%veee%v) src=[1,2] => 1eee2
func sprintf_multi_param(pipe *PipeItem, src interface{}, params string) (interface{}, error) {
	if len(params) == 0 {
		return src, errors.New("filter split nil params ")
	}

	srcValue := reflect.ValueOf(src)
	srcType := srcValue.Type()
	if srcType.Kind() == reflect.Array || srcType.Kind() == reflect.Slice {
		count := strings.Count(params, "%")
		ret := make([]interface{}, 0)
		for i := 0; i < srcValue.Len(); i++ {
			ret = append(ret, srcValue.Index(i).Interface())
		}
		if len(ret) > count {
			return fmt.Sprintf(params, ret[:count]...), nil
		}
		return fmt.Sprintf(params, ret...), nil
	}

	return fmt.Sprintf(params, src), nil
}

// sprintf(%s) src=a => a
func sprintf(pipe *PipeItem, src interface{}, params string) (interface{}, error) {
	if len(params) == 0 {
		return src, errors.New("filter split nil params")
	}
	return _filterValue(src, func(v string) (interface{}, error) {
		return fmt.Sprintf(params, v), nil
	})
}

// sprintfmap(%v-%v,a,b) src={"a":1,"b":2} => "1-2"
func sprintfmap(pipe *PipeItem, src interface{}, params string) (interface{}, error) {
	if len(params) == 0 {
		return src, errors.New("filter split nil params")
	}
	msrc, ok := src.(map[string]interface{})
	if ok == false {
		return src, errors.New("value is not map[string]interface{}")
	}
	vt := SplitParams(params)
	if len(vt) <= 1 {
		return src, errors.New("params length must > 1")
	}
	pArray := []interface{}{}
	for _, x := range vt[1:] {
		if vm, ok := msrc[x]; ok {
			pArray = append(pArray, vm)
		} else {
			pArray = append(pArray, nil)
		}
	}
	return fmt.Sprintf(vt[0], pArray...), nil
}

var timeFormatNames = map[string]string{
	`Layout`:      time.Layout,
	`ANSIC`:       time.ANSIC,
	`UnixDate`:    time.UnixDate,
	`RubyDate`:    time.RubyDate,
	`RFC822`:      time.RFC822,
	`RFC822Z`:     time.RFC822Z,
	`RFC850`:      time.RFC850,
	`RFC1123`:     time.RFC1123,
	`RFC1123Z`:    time.RFC1123Z,
	`RFC3339`:     time.RFC3339,
	`RFC3339Nano`: time.RFC3339Nano,
	`Kitchen`:     time.Kitchen,
	// Handy time stamps.
	`Stamp`:      time.Stamp,
	`StampMilli`: time.StampMilli,
	`StampMicro`: time.StampMicro,
	`StampNano`:  time.StampNano,
	`DateTime`:   time.DateTime,
	`DateOnly`:   time.DateOnly,
	`TimeOnly`:   time.TimeOnly,
}

// unixtime 时间戳(总秒数)
func unixtime(pipe *PipeItem, src interface{}, params string) (interface{}, error) {
	if len(params) == 0 {
		return time.Now().Unix(), nil
	}
	layout := params
	if len(layout) == 0 {
		return time.Now().Unix(), nil
	}
	srcString, ok := src.(string)
	if !ok {
		return src, errors.New("value is not string")
	}
	if v, y := timeFormatNames[layout]; y {
		layout = v
	} else {
		layout = com.ConvDateFormat(layout)
	}
	t, err := time.Parse(layout, srcString)
	if err != nil {
		return t, err
	}
	return t.Unix(), nil
}

// unixmill 时间戳(总毫秒数)
func unixmill(pipe *PipeItem, src interface{}, params string) (interface{}, error) {
	return time.Now().UnixNano() / int64(time.Millisecond), nil
}

// paging(startAt,endAt,step)
// paging(1,10) / paging(1,10,2)
func paging(pipe *PipeItem, src interface{}, params string) (interface{}, error) {
	if len(params) == 0 {
		return src, errors.New("filter paging nil params")
	}
	srcType := reflect.TypeOf(src).Kind()
	if srcType != reflect.Slice && srcType != reflect.Array && srcType != reflect.String {
		return src, errors.New("value is not slice,array or string")
	}
	vt := strings.Split(params, ",")
	if len(vt) < 2 {
		return src, errors.New("params length must > 1")
	}

	start, err := strconv.Atoi(vt[0])
	if err != nil {
		return src, errors.New("params type error:need int." + err.Error())
	}
	end, err := strconv.Atoi(vt[1])
	if err != nil {
		return src, errors.New("params type error:need int." + err.Error())
	}

	offset := -1
	if len(vt) == 3 {
		offset, err = strconv.Atoi(vt[2])
		if err != nil {
			return src, errors.New("params type error:need int." + err.Error())
		}
		if offset < 1 {
			return src, errors.New("offset must > 0")
		}
	}

	var result []string
	switch vt := src.(type) {
	case []interface{}:
		for i := start; i <= end; i++ {
			for _, v := range vt {
				if offset > 0 {
					result = append(result, sprintf_replace(v.(string), []string{strconv.Itoa(i * offset), strconv.Itoa((i + 1) * offset)}))
				} else {
					result = append(result, sprintf_replace(v.(string), []string{strconv.Itoa(i)}))
				}
			}

		}
		return result, nil

	case []string:
		for i := start; i <= end; i++ {
			for _, v := range vt {
				if offset > 0 {
					result = append(result, sprintf_replace(v, []string{strconv.Itoa(i * offset), strconv.Itoa((i + 1) * offset)}))
				} else {
					result = append(result, sprintf_replace(v, []string{strconv.Itoa(i)}))
				}
			}

		}
		return result, nil

	case string:
		for i := start; i <= end; i++ {
			if offset > 0 {
				result = append(result, sprintf_replace(vt, []string{strconv.Itoa(i * offset), strconv.Itoa((i + 1) * offset)}))
			} else {
				result = append(result, sprintf_replace(vt, []string{strconv.Itoa(i)}))
			}
		}
		return result, nil

	default:
		return vt, errors.New("do nothing,src type not support!")
	}
}

func sprintf_replace(src string, param []string) string {
	for i := range param {
		src = strings.Replace(src, "{"+strconv.Itoa(i)+"}", param[i], -1)
	}
	return src
}

// quote => src=`a` => `"a"`
func quote(pipe *PipeItem, src interface{}, params string) (interface{}, error) {
	return _filterValue(src, func(v string) (interface{}, error) {
		return strconv.Quote(v), nil
	})
}

// unquote => src=`"a"` => `a`
func unquote(pipe *PipeItem, src interface{}, params string) (interface{}, error) {
	return _filterValue(src, func(v string) (interface{}, error) {
		return strconv.Unquote(`"` + v + `"`)
	})
}
