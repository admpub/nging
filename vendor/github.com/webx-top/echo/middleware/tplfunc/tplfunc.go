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
	"fmt"
	"html/template"
	"net/url"
	"path/filepath"
	"strings"
	"time"

	"github.com/webx-top/captcha"
	"github.com/webx-top/com"
)

func New() (r template.FuncMap) {
	r = template.FuncMap{}
	for name, function := range TplFuncMap {
		r[name] = function
	}
	return
}

var TplFuncMap template.FuncMap = template.FuncMap{
	"Now":             Now,
	"Eq":              Eq,
	"Add":             Add,
	"Sub":             Sub,
	"IsNil":           IsNil,
	"IsEmpty":         IsEmpty,
	"NotEmpty":        NotEmpty,
	"Html":            ToHTML,
	"Js":              ToJS,
	"Css":             ToCSS,
	"ToJS":            ToJS,
	"ToCSS":           ToCSS,
	"ToURL":           ToURL,
	"ToHTML":          ToHTML,
	"ToHTMLAttr":      ToHTMLAttr,
	"ToHTMLAttrs":     ToHTMLAttrs,
	"ToStrSlice":      ToStrSlice,
	"Concat":          Concat,
	"ElapsedMemory":   com.ElapsedMemory, //内存消耗
	"TotalRunTime":    com.TotalRunTime,  //运行时长(从启动服务时算起)
	"CaptchaForm":     CaptchaForm,       //验证码图片
	"FormatByte":      com.FormatByte,    //字节转为适合理解的格式
	"FriendlyTime":    FriendlyTime,
	"FormatPastTime":  com.FormatPastTime, //以前距离现在多长时间
	"DateFormat":      com.DateFormat,
	"DateFormatShort": com.DateFormatShort,
	"Replace":         strings.Replace, //strings.Replace(s, old, new, n)
	"Contains":        strings.Contains,
	"HasPrefix":       strings.HasPrefix,
	"HasSuffix":       strings.HasSuffix,
	"Split":           strings.Split,
	"Join":            strings.Join,
	"Ext":             filepath.Ext,
	"InExt":           InExt,
	"Str":             com.Str,
	"Int":             com.Int,
	"Int32":           com.Int32,
	"Int64":           com.Int64,
	"Float32":         com.Float32,
	"Float64":         com.Float64,
	"InSlice":         com.InSlice,
	"InSlicex":        com.InSliceIface,
	"Substr":          com.Substr,
	"StripTags":       com.StripTags,
	"Default":         Default,
	"JsonEncode":      JsonEncode,
	"UrlEncode":       com.UrlEncode,
	"UrlDecode":       com.UrlDecode,
	"Base64Encode":    com.Base64Encode,
	"Base64Decode":    com.Base64Decode,
	"Set":             Set,
	"Append":          Append,
	"Nl2br":           NlToBr,
	"AddSuffix":       AddSuffix,
	"InStrSlice":      InStrSlice,
	"SearchStrSlice":  SearchStrSlice,
	"URLValues":       URLValues,
	"ToSlice":         ToSlice,
}

func JsonEncode(s interface{}) string {
	r, _ := com.SetJSON(s)
	return r
}

func URLValues(values ...interface{}) url.Values {
	v := url.Values{}
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

//NlToBr Replaces newlines with <br />
func NlToBr(text string) template.HTML {
	return template.HTML(Nl2br(text))
}

//CaptchaForm 验证码表单域
func CaptchaForm(args ...string) template.HTML {
	id := "captcha"
	format := `<img id="%[2]sImage" src="/captcha/%[1]s.png" alt="Captcha image" onclick="this.src=this.src.split('?')[0]+'?reload='+Math.random();" /><input type="hidden" name="captchaId" id="%[2]sId" value="%[1]s" />`
	switch len(args) {
	case 2:
		format = args[1]
		fallthrough
	case 1:
		id = args[0]
	}
	cid := captcha.New()
	return template.HTML(fmt.Sprintf(format, cid, id))
}

//CaptchaVerify 验证码验证
func CaptchaVerify(captchaSolution string, idGet func(string) string) bool {
	//id := r.FormValue("captchaId")
	id := idGet("captchaId")
	if !captcha.VerifyString(id, captchaSolution) {
		return false
	}
	return true
}

//Nl2br 将换行符替换为<br />
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

func Add(left interface{}, right interface{}) interface{} {
	var rleft, rright int64
	var fleft, fright float64
	isInt := true
	switch left.(type) {
	case int:
		rleft = int64(left.(int))
	case int8:
		rleft = int64(left.(int8))
	case int16:
		rleft = int64(left.(int16))
	case int32:
		rleft = int64(left.(int32))
	case int64:
		rleft = left.(int64)
	case float32:
		fleft = float64(left.(float32))
		isInt = false
	case float64:
		fleft = left.(float64)
		isInt = false
	}

	switch right.(type) {
	case int:
		rright = int64(right.(int))
	case int8:
		rright = int64(right.(int8))
	case int16:
		rright = int64(right.(int16))
	case int32:
		rright = int64(right.(int32))
	case int64:
		rright = right.(int64)
	case float32:
		fright = float64(left.(float32))
		isInt = false
	case float64:
		fleft = left.(float64)
		isInt = false
	}

	intSum := rleft + rright

	if isInt {
		return intSum
	}
	return fleft + fright + float64(intSum)
}

func Sub(left interface{}, right interface{}) interface{} {
	var rleft, rright int64
	var fleft, fright float64
	isInt := true
	switch left.(type) {
	case int:
		rleft = int64(left.(int))
	case int8:
		rleft = int64(left.(int8))
	case int16:
		rleft = int64(left.(int16))
	case int32:
		rleft = int64(left.(int32))
	case int64:
		rleft = left.(int64)
	case float32:
		fleft = float64(left.(float32))
		isInt = false
	case float64:
		fleft = left.(float64)
		isInt = false
	}

	switch right.(type) {
	case int:
		rright = int64(right.(int))
	case int8:
		rright = int64(right.(int8))
	case int16:
		rright = int64(right.(int16))
	case int32:
		rright = int64(right.(int32))
	case int64:
		rright = right.(int64)
	case float32:
		fright = float64(left.(float32))
		isInt = false
	case float64:
		fleft = left.(float64)
		isInt = false
	}

	if isInt {
		return rleft - rright
	}
	return fleft + float64(rleft) - (fright + float64(rright))
}

func Now() time.Time {
	return time.Now()
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

func FriendlyTime(t interface{}, args ...string) string {
	if v, y := t.(time.Duration); y {
		return com.FriendlyTime(v, args...)
	}
	if v, y := t.(int64); y {
		return com.FriendlyTime(time.Duration(v), args...)
	}
	if v, y := t.(int); y {
		return com.FriendlyTime(time.Duration(v), args...)
	}
	if v, y := t.(uint); y {
		return com.FriendlyTime(time.Duration(v), args...)
	}
	if v, y := t.(int32); y {
		return com.FriendlyTime(time.Duration(v), args...)
	}
	if v, y := t.(uint32); y {
		return com.FriendlyTime(time.Duration(v), args...)
	}
	return com.FriendlyTime(time.Duration(com.Int64(t)), args...)
}
