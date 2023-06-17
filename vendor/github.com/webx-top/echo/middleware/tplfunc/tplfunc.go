/*

   Copyright 2016 Wenhui Shen <www.webx.top>

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

package tplfunc

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"math"
	"net/url"
	"path"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/admpub/decimal"

	"github.com/webx-top/captcha"
	"github.com/webx-top/com"
	"github.com/webx-top/echo/param"
)

func New() (r template.FuncMap) {
	r = template.FuncMap{}
	for name, function := range TplFuncMap {
		r[name] = function
	}
	return
}

var TplFuncMap template.FuncMap = template.FuncMap{
	// ======================
	// time
	// ======================
	"Now":             Now,
	"UnixTime":        UnixTime,
	"ElapsedMemory":   com.ElapsedMemory, //内存消耗
	"TotalRunTime":    com.TotalRunTime,  //运行时长(从启动服务时算起)
	"CaptchaForm":     CaptchaForm,       //验证码图片
	"FormatByte":      com.FormatBytes,   //字节转为适合理解的格式
	"FormatBytes":     com.FormatBytes,   //字节转为适合理解的格式
	"FriendlyTime":    FriendlyTime,
	"FormatPastTime":  com.FormatPastTime, //以前距离现在多长时间
	"DateFormat":      com.DateFormat,
	"DateFormatShort": com.DateFormatShort,
	"Ts2time":         TsToTime, // 时间戳数字转time.Time
	"Ts2date":         TsToDate, // 时间戳数字转日期字符串

	// ======================
	// compare
	// ======================
	"Eq":       Eq,
	"Add":      Add,
	"Sub":      Sub,
	"Div":      Div,
	"Mul":      Mul,
	"IsNil":    IsNil,
	"IsEmpty":  IsEmpty,
	"NotEmpty": NotEmpty,
	"IsNaN":    IsNaN,
	"IsInf":    IsInf,

	// ======================
	// conversion type
	// ======================
	"Html":           ToHTML,
	"Js":             ToJS,
	"Css":            ToCSS,
	"ToJS":           ToJS,
	"ToCSS":          ToCSS,
	"ToURL":          ToURL,
	"ToHTML":         ToHTML,
	"ToHTMLAttr":     ToHTMLAttr,
	"ToHTMLAttrs":    ToHTMLAttrs,
	"ToStrSlice":     ToStrSlice,
	"ToDuration":     ToDuration,
	"Str":            com.Str,
	"Int":            com.Int,
	"Int32":          com.Int32,
	"Int64":          com.Int64,
	"Uint":           com.Uint,
	"Uint32":         com.Uint32,
	"Uint64":         com.Uint64,
	"Float32":        com.Float32,
	"Float64":        com.Float64,
	"Float2int":      com.Float2int,
	"Float2uint":     com.Float2uint,
	"Float2int64":    com.Float2int64,
	"Float2uint64":   com.Float2uint64,
	"ToFloat64":      ToFloat64,
	"ToFixed":        ToFixed,
	"ToDecimal":      ToDecimal,
	"Math":           Math,
	"NumberFormat":   NumberFormat,
	"NumberTrim":     NumberTrim,
	"DurationFormat": DurationFormat,
	"DelimLeft":      DelimLeft,
	"DelimRight":     DelimRight,
	"TemplateTag":    TemplateTag,

	// ======================
	// string
	// ======================
	"Contains":   strings.Contains,
	"HasPrefix":  strings.HasPrefix,
	"HasSuffix":  strings.HasSuffix,
	"Trim":       strings.TrimSpace,
	"TrimLeft":   strings.TrimLeft,
	"TrimRight":  strings.TrimRight,
	"TrimPrefix": strings.TrimPrefix,
	"TrimSuffix": strings.TrimSuffix,

	"ToLower":        strings.ToLower,
	"ToUpper":        strings.ToUpper,
	"Title":          strings.Title,
	"LowerCaseFirst": com.LowerCaseFirst,
	"UpperCaseFirst": com.UpperCaseFirst,
	"CamelCase":      com.CamelCase,
	"PascalCase":     com.PascalCase,
	"SnakeCase":      com.SnakeCase,
	"Reverse":        com.Reverse,
	"Dir":            filepath.Dir,
	"Base":           filepath.Base,
	"Ext":            filepath.Ext,
	"Dirname":        path.Dir,
	"Basename":       path.Base,
	"Extension":      path.Ext,
	"InExt":          InExt,

	"Concat":    Concat,
	"Replace":   strings.Replace, //strings.Replace(s, old, new, n)
	"Split":     strings.Split,
	"Join":      strings.Join,
	"Substr":    com.Substr,
	"StripTags": com.StripTags,
	"Nl2br":     NlToBr, // \n替换为<br>
	"AddSuffix": AddSuffix,

	// ======================
	// encode & decode
	// ======================
	"JSONEncode":       JSONEncode,
	"JSONDecode":       JSONDecode,
	"URLEncode":        com.URLEncode,
	"URLDecode":        URLDecode,
	"RawURLEncode":     com.RawURLEncode,
	"RawURLDecode":     URLDecode,
	"Base64Encode":     com.Base64Encode,
	"Base64Decode":     Base64Decode,
	"UnicodeDecode":    UnicodeDecode,
	"SafeBase64Encode": com.SafeBase64Encode,
	"SafeBase64Decode": SafeBase64Decode,
	"Hash":             Hash,
	"Unquote":          Unquote,
	"Quote":            strconv.Quote,

	// ======================
	// map & slice
	// ======================
	"MakeMap":        MakeMap,
	"InSet":          com.InSet,
	"InSlice":        com.InSlice,
	"InSlicex":       com.InSliceIface,
	"Set":            Set,
	"Append":         Append,
	"InStrSlice":     InStrSlice,
	"SearchStrSlice": SearchStrSlice,
	"URLValues":      URLValues,
	"ToSlice":        ToSlice,
	"StrToSlice":     StrToSlice,
	"GetByIndex":     param.GetByIndex,
	"ToParamString":  func(v string) param.String { return param.String(v) },

	// ======================
	// regexp
	// ======================
	"Regexp":      regexp.MustCompile,
	"RegexpPOSIX": regexp.MustCompilePOSIX,

	// ======================
	// other
	// ======================
	"Ignore":        Ignore,
	"Default":       Default,
	"WithURLParams": com.WithURLParams,
	"FullURL":       com.FullURL,
	"IsFullURL":     com.IsFullURL,
}

var (
	HashSalt          = time.Now().Format(time.RFC3339)
	HashClipPositions = []uint{1, 3, 8, 9}
	NumberFormat      = com.NumberFormat
)

func Hash(text, salt string, positions ...uint) string {
	if len(salt) < 1 {
		salt = HashSalt
	}
	if len(positions) < 1 {
		positions = HashClipPositions
	}
	return com.MakePassword(text, salt, positions...)
}

func Unquote(s string) string {
	r, _ := strconv.Unquote(`"` + s + `"`)
	return r
}

func UnicodeDecode(str string) string {
	buf := bytes.NewBuffer(nil)
	i, j := 0, len(str)
	for i < j {
		x := i + 6
		if x > j {
			buf.WriteString(str[i:])
			break
		}
		if str[i] == '\\' && str[i+1] == 'u' {
			hex := str[i+2 : x]
			r, err := strconv.ParseUint(hex, 16, 64)
			if err == nil {
				buf.WriteRune(rune(r))
			} else {
				buf.WriteString(str[i:x])
			}
			i = x
		} else {
			buf.WriteByte(str[i])
			i++
		}
	}
	return buf.String()
}

func JSONEncode(s interface{}, indents ...string) string {
	r, _ := com.JSONEncode(s, indents...)
	return string(r)
}

func JSONDecode(s string) map[string]interface{} {
	r := map[string]interface{}{}
	e := com.JSONDecode([]byte(s), &r)
	if e != nil {
		log.Println(e)
	}
	return r
}

func URLDecode(s string) string {
	r, e := com.URLDecode(s)
	if e != nil {
		log.Println(e)
	}
	return r
}

func Base64Decode(s string) string {
	r, e := com.Base64Decode(s)
	if e != nil {
		log.Println(e)
	}
	return r
}

func SafeBase64Decode(s string) string {
	r, e := com.SafeBase64Decode(s)
	if e != nil {
		log.Println(e)
	}
	return r
}

func Ignore(_ interface{}) interface{} {
	return nil
}

func URLValues(values ...interface{}) url.Values {
	v := url.Values{}
	return AddURLValues(v, values...)
}

func AddURLValues(v url.Values, values ...interface{}) url.Values {
	var k string
	for i, j := 0, len(values); i < j; i++ {
		if i%2 == 0 {
			k = fmt.Sprint(values[i])
			continue
		}
		v.Add(k, fmt.Sprint(values[i]))
		k = ``
	}
	if len(k) > 0 {
		v.Add(k, ``)
		k = ``
	}
	return v
}

func ToStrSlice(s ...string) []string {
	return s
}

func ToSlice(s ...interface{}) []interface{} {
	return s
}

func StrToSlice(s string, sep string) []interface{} {
	ss := strings.Split(s, sep)
	r := make([]interface{}, len(ss))
	for i, s := range ss {
		r[i] = s
	}
	return r
}

func Concat(s ...string) string {
	return strings.Join(s, ``)
}

func InExt(fileName string, exts ...string) bool {
	ext := filepath.Ext(fileName)
	ext = strings.ToLower(ext)
	for _, _ext := range exts {
		if ext == strings.ToLower(_ext) {
			return true
		}
	}
	return false
}

func Default(defaultV interface{}, v interface{}) interface{} {
	switch val := v.(type) {
	case nil:
		return defaultV
	case string:
		if len(val) == 0 {
			return defaultV
		}
	case uint8, int8, uint, int, uint32, int32, int64, uint64:
		if val == 0 {
			return defaultV
		}
	case float32, float64:
		if val == 0.0 {
			return defaultV
		}
	default:
		if len(com.Str(v)) == 0 {
			return defaultV
		}
	}
	return v
}

func Set(renderArgs map[string]interface{}, key string, value interface{}) string {
	renderArgs[key] = value
	return ``
}

func Append(renderArgs map[string]interface{}, key string, value interface{}) string {
	if renderArgs[key] == nil {
		renderArgs[key] = []interface{}{value}
	} else {
		renderArgs[key] = append(renderArgs[key].([]interface{}), value)
	}
	return ``
}

// NlToBr Replaces newlines with <br />
func NlToBr(text string) template.HTML {
	return template.HTML(Nl2br(text))
}

// CaptchaForm 验证码表单域
func CaptchaForm(args ...interface{}) template.HTML {
	return CaptchaFormWithURLPrefix(``, args...)
}

// CaptchaFormWithURLPrefix 验证码表单域
func CaptchaFormWithURLPrefix(urlPrefix string, args ...interface{}) template.HTML {
	id := "captcha"
	msg := "页面验证码已经失效，必须重新请求当前页面。确定要刷新本页面吗？"
	onErr := "if(this.src.indexOf('?reload=')!=-1 && confirm('%s')) window.location.reload();"
	format := `<img id="%[2]sImage" src="` + urlPrefix + `/captcha/%[1]s.png" alt="Captcha image" onclick="this.src=this.src.split('?')[0]+'?reload='+Math.random();" onerror="%[3]s" style="cursor:pointer" /><input type="hidden" name="captchaId" id="%[2]sId" value="%[1]s" />`
	var (
		customOnErr bool
		cid         string
	)
	switch len(args) {
	case 3:
		switch v := args[2].(type) {
		case template.JS:
			onErr = string(v)
			customOnErr = true
		case string:
			msg = v
		}
		fallthrough
	case 2:
		if args[1] != nil {
			v := fmt.Sprint(args[1])
			format = v
		}
		fallthrough
	case 1:
		switch v := args[0].(type) {
		case template.JS:
			onErr = string(v)
			customOnErr = true
		case template.HTML:
			format = string(v)
		case string:
			id = v
		case param.Store:
			cid = v.String(`captchaId`)
			if v.Has(`onErr`) {
				onErr = v.String(`onErr`)
			}
			if v.Has(`format`) {
				format = v.String(`format`)
			}
			if v.Has(`id`) {
				id = v.String(`id`)
			}
		case map[string]interface{}:
			h := param.Store(v)
			cid = h.String(`captchaId`)
			if h.Has(`onErr`) {
				onErr = h.String(`onErr`)
			}
			if h.Has(`format`) {
				format = h.String(`format`)
			}
			if h.Has(`id`) {
				id = h.String(`id`)
			}
		}
	}
	if len(cid) == 0 {
		cid = captcha.New()
	}
	if !customOnErr {
		onErr = fmt.Sprintf(onErr, msg)
	}
	return template.HTML(fmt.Sprintf(format, cid, id, onErr))
}

// CaptchaVerify 验证码验证
func CaptchaVerify(captchaSolution string, idGet func(string, ...string) string) bool {
	//id := r.FormValue("captchaId")
	id := idGet("captchaId")
	return captcha.VerifyString(id, captchaSolution)
}

// Nl2br 将换行符替换为<br />
func Nl2br(text string) string {
	return com.Nl2br(template.HTMLEscapeString(text))
}

func IsNil(a interface{}) bool {
	switch a.(type) {
	case nil:
		return true
	}
	return false
}

func interface2Int64(value interface{}) (int64, bool) {
	switch v := value.(type) {
	case uint:
		return int64(v), true
	case uint8:
		return int64(v), true
	case uint16:
		return int64(v), true
	case uint32:
		return int64(v), true
	case uint64:
		return int64(v), true
	case int:
		return int64(v), true
	case int8:
		return int64(v), true
	case int16:
		return int64(v), true
	case int32:
		return int64(v), true
	case int64:
		return v, true
	default:
		return 0, false
	}
}

func interface2Float64(value interface{}) (float64, bool) {
	switch v := value.(type) {
	case float32:
		return float64(v), true
	case float64:
		return v, true
	default:
		return 0, false
	}
}

func ToFloat64(value interface{}) float64 {
	if v, ok := interface2Int64(value); ok {
		return float64(v)
	}
	if v, ok := interface2Float64(value); ok {
		return v
	}
	return com.Float64(value)
}

func Add(left interface{}, right interface{}) interface{} {
	var rleft, rright int64
	var fleft, fright float64
	var isInt bool
	rleft, isInt = interface2Int64(left)
	if !isInt {
		fleft, _ = interface2Float64(left)
	}
	rright, isInt = interface2Int64(right)
	if !isInt {
		fright, _ = interface2Float64(right)
	}
	intSum := rleft + rright

	if isInt {
		return intSum
	}
	return fleft + fright + float64(intSum)
}

func Div(left interface{}, right interface{}) interface{} {
	return ToFloat64(left) / ToFloat64(right)
}

func Mul(left interface{}, right interface{}) interface{} {
	return ToFloat64(left) * ToFloat64(right)
}

func Math(op string, args ...interface{}) interface{} {
	length := len(args)
	if length < 1 {
		return float64(0)
	}
	switch op {
	case `mod`: //模
		if length < 2 {
			return float64(0)
		}
		return math.Mod(ToFloat64(args[0]), ToFloat64(args[1]))
	case `abs`:
		return math.Abs(ToFloat64(args[0]))
	case `acos`:
		return math.Acos(ToFloat64(args[0]))
	case `acosh`:
		return math.Acosh(ToFloat64(args[0]))
	case `asin`:
		return math.Asin(ToFloat64(args[0]))
	case `asinh`:
		return math.Asinh(ToFloat64(args[0]))
	case `atan`:
		return math.Atan(ToFloat64(args[0]))
	case `atan2`:
		if length < 2 {
			return float64(0)
		}
		return math.Atan2(ToFloat64(args[0]), ToFloat64(args[1]))
	case `atanh`:
		return math.Atanh(ToFloat64(args[0]))
	case `cbrt`:
		return math.Cbrt(ToFloat64(args[0]))
	case `ceil`:
		return math.Ceil(ToFloat64(args[0]))
	case `copysign`:
		if length < 2 {
			return float64(0)
		}
		return math.Copysign(ToFloat64(args[0]), ToFloat64(args[1]))
	case `cos`:
		return math.Cos(ToFloat64(args[0]))
	case `cosh`:
		return math.Cosh(ToFloat64(args[0]))
	case `dim`:
		if length < 2 {
			return float64(0)
		}
		return math.Dim(ToFloat64(args[0]), ToFloat64(args[1]))
	case `erf`:
		return math.Erf(ToFloat64(args[0]))
	case `erfc`:
		return math.Erfc(ToFloat64(args[0]))
	case `exp`:
		return math.Exp(ToFloat64(args[0]))
	case `exp2`:
		return math.Exp2(ToFloat64(args[0]))
	case `floor`:
		return math.Floor(ToFloat64(args[0]))
	case `max`:
		if length < 2 {
			return float64(0)
		}
		return math.Max(ToFloat64(args[0]), ToFloat64(args[1]))
	case `min`:
		if length < 2 {
			return float64(0)
		}
		return math.Min(ToFloat64(args[0]), ToFloat64(args[1]))
	case `pow`: //幂
		if length < 2 {
			return float64(0)
		}
		return math.Pow(ToFloat64(args[0]), ToFloat64(args[1]))
	case `sqrt`: //平方根
		return math.Sqrt(ToFloat64(args[0]))
	case `sin`:
		return math.Sin(ToFloat64(args[0]))
	case `log`:
		return math.Log(ToFloat64(args[0]))
	case `log2`:
		return math.Log2(ToFloat64(args[0]))
	case `log10`:
		return math.Log10(ToFloat64(args[0]))
	case `tan`:
		return math.Tan(ToFloat64(args[0]))
	case `tanh`:
		return math.Tanh(ToFloat64(args[0]))
	case `add`: //加
		if length < 2 {
			return float64(0)
		}
		return Add(ToFloat64(args[0]), ToFloat64(args[1]))
	case `sub`: //减
		if length < 2 {
			return float64(0)
		}
		return Sub(ToFloat64(args[0]), ToFloat64(args[1]))
	case `mul`: //乘
		if length < 2 {
			return float64(0)
		}
		return Mul(ToFloat64(args[0]), ToFloat64(args[1]))
	case `div`: //除
		if length < 2 {
			return float64(0)
		}
		return Div(ToFloat64(args[0]), ToFloat64(args[1]))
	}
	return nil
}

func IsNaN(v interface{}) bool {
	return math.IsNaN(ToFloat64(v))
}

func IsInf(v interface{}, s interface{}) bool {
	return math.IsInf(ToFloat64(v), com.Int(s))
}

func Sub(left interface{}, right interface{}) interface{} {
	var rleft, rright int64
	var fleft, fright float64
	var isInt bool
	rleft, isInt = interface2Int64(left)
	if !isInt {
		fleft, _ = interface2Float64(left)
	}
	rright, isInt = interface2Int64(right)
	if !isInt {
		fright, _ = interface2Float64(right)
	}
	if isInt {
		return rleft - rright
	}
	return fleft + float64(rleft) - (fright + float64(rright))
}

func ToFixed(value interface{}, precision interface{}) string {
	return fmt.Sprintf("%.*f", com.Int(precision), ToFloat64(value))
}

func Now() time.Time {
	return time.Now()
}

func UnixTime() int64 {
	return time.Now().Unix()
}

func Eq(left interface{}, right interface{}) bool {
	leftIsNil := (left == nil)
	rightIsNil := (right == nil)
	if leftIsNil || rightIsNil {
		if leftIsNil && rightIsNil {
			return true
		}
		return false
	}
	return fmt.Sprintf("%v", left) == fmt.Sprintf("%v", right)
}

func ToHTML(raw string) template.HTML {
	return template.HTML(raw)
}

func ToHTMLAttr(raw string) template.HTMLAttr {
	return template.HTMLAttr(raw)
}

func ToHTMLAttrs(raw map[string]interface{}) (r map[template.HTMLAttr]interface{}) {
	r = make(map[template.HTMLAttr]interface{})
	for k, v := range raw {
		r[ToHTMLAttr(k)] = v
	}
	return
}

func ToJS(raw string) template.JS {
	return template.JS(raw)
}

func ToCSS(raw string) template.CSS {
	return template.CSS(raw)
}

func ToURL(raw string) template.URL {
	return template.URL(raw)
}

func AddSuffix(s string, suffix string, args ...string) string {
	beforeChar := `.`
	if len(args) > 0 {
		beforeChar = args[0]
		if beforeChar == `` {
			return s + suffix
		}
	}
	p := strings.LastIndex(s, beforeChar)
	if p < 0 {
		return s
	}
	return s[0:p] + suffix + s[p:]
}

func IsEmpty(a interface{}) bool {
	switch v := a.(type) {
	case nil:
		return true
	case string:
		return len(v) == 0
	case []interface{}:
		return len(v) < 1
	default:
		switch fmt.Sprintf(`%v`, a) {
		case `<nil>`, ``, `[]`:
			return true
		}
	}
	return false
}

func NotEmpty(a interface{}) bool {
	return !IsEmpty(a)
}

func InStrSlice(values []string, value string) bool {
	for _, v := range values {
		if v == value {
			return true
		}
	}
	return false
}

func SearchStrSlice(values []string, value string) int {
	for i, v := range values {
		if v == value {
			return i
		}
	}
	return -1
}

func DurationFormat(lang interface{}, t interface{}, args ...string) *com.Durafmt {
	duration := ToDuration(t, args...)
	return com.ParseDuration(duration, lang)
}

func ToTime(t interface{}) time.Time {
	switch v := t.(type) {
	case time.Time:
		return v
	case string:
		t, err := time.ParseInLocation(`2006-01-02 15:04:05`, v, time.Local)
		if err != nil {
			panic(err)
		}
		return t
	default:
		return TsToTime(t)
	}
}

func ToDuration(t interface{}, args ...string) time.Duration {
	td := time.Second
	if len(args) > 0 {
		switch args[0] {
		case `ns`:
			td = time.Nanosecond
		case `us`:
			td = time.Microsecond
		case `s`:
			td = time.Second
		case `ms`:
			td = time.Millisecond
		case `h`:
			td = time.Hour
		case `m`:
			td = time.Minute
		}
	}
	switch v := t.(type) {
	case time.Duration:
		return v
	case int64:
		td = time.Duration(v) * td
	case int:
		td = time.Duration(v) * td
	case uint:
		td = time.Duration(v) * td
	case int32:
		td = time.Duration(v) * td
	case uint32:
		td = time.Duration(v) * td
	case uint64:
		td = time.Duration(v) * td
	default:
		td = time.Duration(com.Int64(t)) * td
	}
	return td
}

func FriendlyTime(t interface{}, args ...interface{}) string {
	var td time.Duration
	switch v := t.(type) {
	case time.Duration:
		td = v
	case int64:
		td = time.Duration(v)
	case int:
		td = time.Duration(v)
	case uint:
		td = time.Duration(v)
	case int32:
		td = time.Duration(v)
	case uint32:
		td = time.Duration(v)
	case uint64:
		td = time.Duration(v)
	default:
		td = time.Duration(com.Int64(t))
	}
	return com.FriendlyTime(td, args...)
}

func TsToTime(timestamp interface{}) time.Time {
	return TimestampToTime(timestamp)
}

func TsToDate(format string, timestamp interface{}) string {
	t := TimestampToTime(timestamp)
	if t.IsZero() {
		return ``
	}
	return t.Format(format)
}

func TimestampToTime(timestamp interface{}) time.Time {
	var ts int64
	switch v := timestamp.(type) {
	case int64:
		ts = v
	case uint:
		ts = int64(v)
	case int:
		ts = int64(v)
	case uint32:
		ts = int64(v)
	case int32:
		ts = int64(v)
	case uint64:
		ts = int64(v)
	default:
		i, e := strconv.ParseInt(fmt.Sprint(timestamp), 10, 64)
		if e != nil {
			log.Println(e)
		}
		ts = i
	}
	return time.Unix(ts, 0)
}

func ToDecimal(number interface{}) decimal.Decimal {
	money := ToFloat64(number)
	return decimal.NewFromFloat(money)
}

func NumberTrim(number interface{}, precision int, separator ...string) string {
	money := ToFloat64(number)
	s := decimal.NewFromFloat(money).Truncate(int32(precision)).String()
	return com.NumberTrim(s, precision, separator...)
}

func MakeMap(values ...interface{}) param.Store {
	h := param.Store{}
	if len(values) == 0 {
		return h
	}
	var k string
	for i, j := 0, len(values); i < j; i++ {
		if i%2 == 0 {
			k = fmt.Sprint(values[i])
			continue
		}
		h.Set(k, values[i])
		k = ``
	}
	if len(k) > 0 {
		h.Set(k, nil)
	}
	return h
}

func DelimLeft() string {
	return `{{`
}

func DelimRight() string {
	return `}}`
}

func TemplateTag(name string) string {
	return DelimLeft() + name + DelimRight()
}
