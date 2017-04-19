package echo

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/webx-top/tagfast"
	"github.com/webx-top/validation"
)

// DefaultHTMLFilter html filter (`form_filter:"html"`)
var DefaultHTMLFilter = func(v string) (r string) {
	return v
}

type (
	// Binder is the interface that wraps the Bind method.
	Binder interface {
		Bind(interface{}, Context, ...FormDataFilter) error
		MustBind(interface{}, Context, ...FormDataFilter) error
	}
	binder struct {
		*Echo
	}
)

func (b *binder) MustBind(i interface{}, c Context, filter ...FormDataFilter) (err error) {
	r := c.Request()
	body := r.Body()
	if body == nil {
		err = NewHTTPError(http.StatusBadRequest, "Request body can't be nil")
		return
	}
	defer body.Close()
	ct := r.Header().Get(HeaderContentType)
	err = ErrUnsupportedMediaType
	if strings.HasPrefix(ct, MIMEApplicationJSON) {
		err = json.NewDecoder(body).Decode(i)
	} else if strings.HasPrefix(ct, MIMEApplicationXML) {
		err = xml.NewDecoder(body).Decode(i)
	} else if strings.HasPrefix(ct, MIMEApplicationForm) {
		err = b.structMap(i, r.PostForm().All(), filter...)
	} else if strings.Contains(ct, MIMEMultipartForm) {
		err = b.structMap(i, r.Form().All(), filter...)
	}
	return
}

func (b *binder) Bind(i interface{}, c Context, filter ...FormDataFilter) (err error) {
	err = b.MustBind(i, c, filter...)
	if err == ErrUnsupportedMediaType {
		err = nil
	}
	return
}

// StructMap function mapping params to controller's properties
func (b *binder) structMap(m interface{}, data map[string][]string, filter ...FormDataFilter) error {
	return NamedStructMap(b.Echo, m, data, ``, filter...)
}

// FormNames user[name][test]
func FormNames(s string) []string {
	var res []string
	hasLeft := false
	hasRight := true
	var val []rune
	for _, r := range s {
		if r == '[' {
			if hasRight {
				res = append(res, string(val))
				val = []rune{}
			}
			hasLeft = true
			hasRight = false
			continue
		}
		if r == ']' {
			if hasLeft {
				hasRight = true
			}
			continue
		}
		val = append(val, r)
	}
	if len(val) > 0 {
		res = append(res, string(val))
	}
	return res
}

// NamedStructMap 自动将map值映射到结构体
func NamedStructMap(e *Echo, m interface{}, data map[string][]string, topName string, filters ...FormDataFilter) error {
	vc := reflect.ValueOf(m)
	tc := reflect.TypeOf(m)

	switch tc.Kind() {
	case reflect.Struct:
	case reflect.Ptr:
		vc = vc.Elem()
		tc = tc.Elem()
	}
	var (
		validator *validation.Validation
		filter    FormDataFilter
	)
	if len(filters) > 0 {
		filter = filters[0]
	}
	if filter == nil {
		filter = DefaultNopFilter
	}
	for k, t := range data {
		k, t = filter(k, t)
		if k == `` || k[0] == '_' {
			continue
		}

		if topName != `` {
			if !strings.HasPrefix(k, topName) {
				continue
			}
			k = k[len(topName)+1:]
		}

		v := t[0]
		names := strings.Split(k, `.`)
		var err error
		length := len(names)
		if length == 1 && strings.HasSuffix(k, `]`) {
			names = FormNames(k)
			length = len(names)
		}
		value := vc
		typev := tc
		for i, name := range names {
			name = strings.Title(name)

			//不是最后一个元素
			if i != length-1 {
				if value.Kind() != reflect.Struct {
					e.Logger().Warnf(`arg error, value kind is %v`, value.Kind())
					break
				}
				value = value.FieldByName(name)
				if !value.IsValid() {
					e.Logger().Warnf(`(%v value is not valid %v)`, name, value)
					break
				}
				if !value.CanSet() {
					e.Logger().Warnf(`can not set %v -> %v`, name, value.Interface())
					break
				}
				if value.Kind() == reflect.Ptr {
					if value.IsNil() {
						value.Set(reflect.New(value.Type().Elem()))
					}
					value = value.Elem()
				}
				typev = value.Type()
				f, _ := typev.FieldByName(name)
				if tagfast.Value(tc, f, `form_options`) == `-` {
					break
				}
				continue
			}

			//最后一个元素
			if value.Kind() != reflect.Struct {
				e.Logger().Warnf(`arg error, value %v kind is %v`, name, value.Kind())
				break
			}
			tv := value.FieldByName(name)
			if !tv.IsValid() {
				break
			}
			if !tv.CanSet() {
				e.Logger().Warnf(`can not set %v to %v`, k, tv)
				break
			}
			f, _ := typev.FieldByName(name)
			if tagfast.Value(tc, f, `form_options`) == `-` {
				break
			}
			if tv.Kind() == reflect.Ptr {
				tv.Set(reflect.New(tv.Type().Elem()))
				tv = tv.Elem()
			}

			var l interface{}
			switch k := tv.Kind(); k {
			case reflect.String:
				switch tagfast.Value(tc, f, `form_filter`) {
				case `html`:
					v = DefaultHTMLFilter(v)
				}
				l = v
				tv.Set(reflect.ValueOf(l))
			case reflect.Bool:
				l = (v != `false` && v != `0` && v != ``)
				tv.Set(reflect.ValueOf(l))
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32:
				dateformat := tagfast.Value(tc, f, `form_format`)
				if dateformat != `` {
					t, err := time.Parse(dateformat, v)
					if err != nil {
						e.Logger().Warnf(`arg %v as int: %v`, v, err)
						l = int(0)
					} else {
						l = int(t.Unix())
					}
				} else {
					x, err := strconv.Atoi(v)
					if err != nil {
						e.Logger().Warnf(`arg %v as int: %v`, v, err)
					}
					l = x
				}
				tv.Set(reflect.ValueOf(l))
			case reflect.Int64:
				dateformat := tagfast.Value(tc, f, `form_format`)
				if dateformat != `` {
					t, err := time.Parse(dateformat, v)
					if err != nil {
						e.Logger().Warnf(`arg %v as int64: %v`, v, err)
						l = int64(0)
					} else {
						l = t.Unix()
					}
				} else {
					x, err := strconv.ParseInt(v, 10, 64)
					if err != nil {
						e.Logger().Warnf(`arg %v as int64: %v`, v, err)
					}
					l = x
				}
				tv.Set(reflect.ValueOf(l))
			case reflect.Float32, reflect.Float64:
				x, err := strconv.ParseFloat(v, 64)
				if err != nil {
					e.Logger().Warnf(`arg %v as float64: %v`, v, err)
				}
				l = x
				tv.Set(reflect.ValueOf(l))
			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				dateformat := tagfast.Value(tc, f, `form_format`)
				var x uint64
				var bitSize int
				switch k {
				case reflect.Uint8:
					bitSize = 8
				case reflect.Uint16:
					bitSize = 16
				case reflect.Uint32:
					bitSize = 32
				default:
					bitSize = 64
				}
				if dateformat != `` {
					t, err := time.Parse(dateformat, v)
					if err != nil {
						e.Logger().Warnf(`arg %v as uint: %v`, v, err)
						x = uint64(0)
					} else {
						x = uint64(t.Unix())
					}
				} else {
					x, err = strconv.ParseUint(v, 10, bitSize)
					if err != nil {
						e.Logger().Warnf(`arg %v as uint: %v`, v, err)
					}
				}
				switch k {
				case reflect.Uint:
					l = uint(x)
				case reflect.Uint8:
					l = uint8(x)
				case reflect.Uint16:
					l = uint16(x)
				case reflect.Uint32:
					l = uint32(x)
				default:
					l = x
				}
				tv.Set(reflect.ValueOf(l))
			case reflect.Struct:
				if tvf, ok := tv.Interface().(FromConversion); ok {
					err := tvf.FromString(v)
					if err != nil {
						e.Logger().Warnf(`struct %v invoke FromString faild`, tvf)
					}
				} else if tv.Type().String() == `time.Time` {
					x, err := time.Parse(`2006-01-02 15:04:05.000 -0700`, v)
					if err != nil {
						x, err = time.Parse(`2006-01-02 15:04:05`, v)
						if err != nil {
							x, err = time.Parse(`2006-01-02`, v)
							if err != nil {
								e.Logger().Warnf(`unsupported time format %v, %v`, v, err)
							}
						}
					}
					l = x
					tv.Set(reflect.ValueOf(l))
				} else {
					e.Logger().Warn(`can not set an struct which is not implement Fromconversion interface`)
				}
			case reflect.Ptr:
				e.Logger().Warn(`can not set an ptr of ptr`)
			case reflect.Slice, reflect.Array:
				tt := tv.Type().Elem()
				tk := tt.Kind()

				if tv.IsNil() {
					tv.Set(reflect.MakeSlice(tv.Type(), len(t), len(t)))
				}

				for i, s := range t {
					var err error
					switch tk {
					case reflect.Int, reflect.Int16, reflect.Int32, reflect.Int8, reflect.Int64:
						var v int64
						v, err = strconv.ParseInt(s, 10, tt.Bits())
						if err == nil {
							tv.Index(i).SetInt(v)
						}
					case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
						var v uint64
						v, err = strconv.ParseUint(s, 10, tt.Bits())
						if err == nil {
							tv.Index(i).SetUint(v)
						}
					case reflect.Float32, reflect.Float64:
						var v float64
						v, err = strconv.ParseFloat(s, tt.Bits())
						if err == nil {
							tv.Index(i).SetFloat(v)
						}
					case reflect.Bool:
						var v bool
						v, err = strconv.ParseBool(s)
						if err == nil {
							tv.Index(i).SetBool(v)
						}
					case reflect.String:
						tv.Index(i).SetString(s)
					case reflect.Complex64, reflect.Complex128:
						// TODO:
						err = fmt.Errorf(`unsupported slice element type %v`, tk.String())
					default:
						err = fmt.Errorf(`unsupported slice element type %v`, tk.String())
					}
					if err != nil {
						e.Logger().Warnf(`slice error: %v, %v`, name, err)
					}
				}
			default:
				break
			}

			//validation
			valid := tagfast.Value(tc, f, `valid`)
			if len(valid) == 0 {
				continue
			}
			if validator == nil {
				validator = validation.New()
			}
			ok, err := validator.ValidSimple(name, fmt.Sprintf(`%v`, l), valid)
			if !ok {
				return validator.Errors[0].WithField()
			}
			if err != nil {
				e.Logger().Warn(err)
			}
		}
	}
	return nil
}

// FromConversion a struct implements this interface can be convert from request param to a struct
type FromConversion interface {
	FromString(content string) error
}

// ToConversion a struct implements this interface can be convert from struct to template variable
// Not Implemented
type ToConversion interface {
	ToString() string
}

type (
	//FieldNameFormatter 结构体字段值映射到表单时，结构体字段名称格式化处理
	FieldNameFormatter func(topName, fieldName string) string
	//FormDataFilter 将map映射到结构体时，对名称和值的过滤处理，如果返回的名称为空，则跳过本字段
	FormDataFilter func(key string, values []string) (string, []string)
)

var (
	//DefaultNopFilter 默认过滤器(map->struct)
	DefaultNopFilter FormDataFilter = func(k string, v []string) (string, []string) {
		return k, v
	}
	//DefaultFieldNameFormatter 默认格式化函数(struct->form)
	DefaultFieldNameFormatter FieldNameFormatter = func(topName, fieldName string) string {
		var fName string
		if len(topName) == 0 {
			fName = fieldName
		} else {
			fName = topName + "." + fieldName
		}
		return fName
	}
	//LowerCaseFirstLetter 小写首字母(struct->form)
	LowerCaseFirstLetter FieldNameFormatter = func(topName, fieldName string) string {
		var fName string
		s := []rune(fieldName)
		if len(s) > 0 {
			s[0] = unicode.ToLower(s[0])
			fieldName = string(s)
		}
		if len(topName) == 0 {
			fName = fieldName
		} else {
			fName = topName + "." + fieldName
		}
		return fName
	}
)

//StructToForm 映射struct到form
func StructToForm(ctx Context, m interface{}, topName string, fieldNameFormatter FieldNameFormatter) {
	vc := reflect.ValueOf(m)
	tc := reflect.TypeOf(m)

	switch tc.Kind() {
	case reflect.Struct:
	case reflect.Ptr:
		vc = vc.Elem()
		tc = tc.Elem()
	}
	l := tc.NumField()
	f := ctx.Request().Form()
	if fieldNameFormatter == nil {
		fieldNameFormatter = DefaultFieldNameFormatter
	}

	for i := 0; i < l; i++ {
		fVal := vc.Field(i)
		fTyp := tc.Field(i)

		fName := fieldNameFormatter(topName, fTyp.Name)
		if !fVal.CanInterface() || len(fName) == 0 {
			continue
		}
		switch fTyp.Type.String() {
		case "time.Time":
			if t, y := fVal.Interface().(time.Time); y {
				dateformat := tagfast.Value(tc, fTyp, `form_format`)
				if len(dateformat) > 0 {
					f.Add(fName, t.Format(dateformat))
				} else {
					f.Add(fName, t.Format(`2006-01-02 15:04:05`))
				}
			}
		case "struct":
			StructToForm(ctx, fVal.Interface(), fName, fieldNameFormatter)
		default:
			f.Add(fName, fmt.Sprint(fVal.Interface()))
		}
	}
}
