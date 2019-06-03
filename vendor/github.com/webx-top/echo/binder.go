package echo

import (
	"encoding/xml"
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/admpub/log"
	"github.com/webx-top/echo/encoding/json"
	"github.com/webx-top/echo/engine"
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
		decoders map[string]func(interface{}, Context, ...FormDataFilter) error
	}
)

func NewBinder(e *Echo) Binder {
	return &binder{
		Echo: e,
		decoders: map[string]func(interface{}, Context, ...FormDataFilter) error{
			MIMEApplicationJSON: func(i interface{}, ctx Context, filter ...FormDataFilter) error {
				body := ctx.Request().Body()
				if body == nil {
					return NewHTTPError(http.StatusBadRequest, "Request body can't be nil")
				}
				defer body.Close()
				return json.NewDecoder(body).Decode(i)
			},
			MIMEApplicationXML: func(i interface{}, ctx Context, filter ...FormDataFilter) error {
				body := ctx.Request().Body()
				if body == nil {
					return NewHTTPError(http.StatusBadRequest, "Request body can't be nil")
				}
				defer body.Close()
				return xml.NewDecoder(body).Decode(i)
			},
			MIMEApplicationForm: func(i interface{}, ctx Context, filter ...FormDataFilter) error {
				body := ctx.Request().Body()
				if body == nil {
					return NewHTTPError(http.StatusBadRequest, "Request body can't be nil")
				}
				defer body.Close()
				return NamedStructMap(ctx.Echo(), i, ctx.Request().PostForm().All(), ``, filter...)
			},
			MIMEMultipartForm: func(i interface{}, ctx Context, filter ...FormDataFilter) error {
				body := ctx.Request().Body()
				if body == nil {
					return NewHTTPError(http.StatusBadRequest, "Request body can't be nil")
				}
				defer body.Close()
				return NamedStructMap(ctx.Echo(), i, ctx.Request().Form().All(), ``, filter...)
			},
		},
	}
}

func (b *binder) MustBind(i interface{}, c Context, filter ...FormDataFilter) error {
	contentType := c.Request().Header().Get(HeaderContentType)
	contentType = strings.ToLower(strings.TrimSpace(strings.SplitN(contentType, `;`, 2)[0]))
	if decoder, ok := b.decoders[contentType]; ok {
		return decoder(i, c, filter...)
	}
	return ErrUnsupportedMediaType
}

func (b *binder) Bind(i interface{}, c Context, filter ...FormDataFilter) (err error) {
	err = b.MustBind(i, c, filter...)
	if err == ErrUnsupportedMediaType {
		err = nil
	}
	return
}

func (b *binder) SetDecoders(decoders map[string]func(interface{}, Context, ...FormDataFilter) error) {
	b.decoders = decoders
}

func (b *binder) AddDecoder(mime string, decoder func(interface{}, Context, ...FormDataFilter) error) {
	b.decoders[mime] = decoder
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
				res = append(res, string(val))
				val = []rune{}
				hasLeft = false
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
	default:
		return errors.New(`binder: unsupported type ` + tc.Kind().String())
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
		if len(k) == 0 || k[0] == '_' {
			continue
		}

		if len(topName) > 0 {
			if !strings.HasPrefix(k, topName) {
				continue
			}
			k = k[len(topName)+1:]
		}

		v := t[0]
		names := strings.Split(k, `.`)
		var (
			err      error
			propPath string
		)
		length := len(names)
		if length == 1 && strings.HasSuffix(k, `]`) {
			k = strings.TrimRight(k, `[]`)
			names = FormNames(k)
			length = len(names)
		}
		value := vc
		typev := tc
		for i, name := range names {
			name = strings.Title(name)
			if i > 0 {
				propPath += `.`
			}
			propPath += name

			//不是最后一个元素
			if i != length-1 {
				if value.Kind() != reflect.Struct {
					e.Logger().Warnf(`binder: arg error, value kind is %v`, value.Kind())
					break
				}
				f, _ := typev.FieldByName(name)
				if tagfast.Value(tc, f, `form_options`) == `-` {
					break
				}
				value = value.FieldByName(name)
				if !value.IsValid() {
					e.Logger().Debugf(`binder: %T#%v value is not valid %v`, m, propPath, value)
					break
				}
				if !value.CanSet() {
					e.Logger().Warnf(`binder: can not set %T#%v -> %v`, m, propPath, value.Interface())
					break
				}
				if value.Kind() == reflect.Ptr {
					if value.IsNil() {
						value.Set(reflect.New(value.Type().Elem()))
					}
					value = value.Elem()
				}
				if value.Kind() != reflect.Struct {
					e.Logger().Warnf(`binder: arg error, value %T#%v kind is %v`, m, propPath, value.Kind())
					break
				}
				typev = value.Type()
				f, _ = typev.FieldByName(name)
				if tagfast.Value(tc, f, `form_options`) == `-` {
					break
				}
				continue
			}

			//最后一个元素
			tv := value.FieldByName(name)
			if !tv.IsValid() {
				break
			}
			if !tv.CanSet() {
				e.Logger().Warnf(`binder: can not set %v to %v`, k, tv)
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
				default:
					delimter := tagfast.Value(tc, f, `form_delimiter`)
					if len(delimter) > 0 {
						v = strings.Join(t, delimter)
					}
				}
				l = v
				tv.Set(reflect.ValueOf(l))
			case reflect.Bool:
				l = (v != `false` && v != `0` && v != ``)
				tv.Set(reflect.ValueOf(l))
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32:
				dateformat := tagfast.Value(tc, f, `form_format`)
				if len(dateformat) > 0 {
					t, err := time.Parse(dateformat, v)
					if err != nil {
						e.Logger().Warnf(`binder: arg %v as int: %v`, v, err)
						l = int(0)
					} else {
						l = int(t.Unix())
					}
				} else {
					x, err := strconv.Atoi(v)
					if err != nil {
						e.Logger().Warnf(`binder: arg %v as int: %v`, v, err)
					}
					l = x
				}
				tv.Set(reflect.ValueOf(l))
			case reflect.Int64:
				dateformat := tagfast.Value(tc, f, `form_format`)
				if len(dateformat) > 0 {
					t, err := time.Parse(dateformat, v)
					if err != nil {
						e.Logger().Warnf(`binder: arg %v as int64: %v`, v, err)
						l = int64(0)
					} else {
						l = t.Unix()
					}
				} else {
					x, err := strconv.ParseInt(v, 10, 64)
					if err != nil {
						e.Logger().Warnf(`binder: arg %v as int64: %v`, v, err)
					}
					l = x
				}
				tv.Set(reflect.ValueOf(l))
			case reflect.Float32, reflect.Float64:
				x, err := strconv.ParseFloat(v, 64)
				if err != nil {
					e.Logger().Warnf(`binder: arg %v as float64: %v`, v, err)
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
				if len(dateformat) > 0 {
					t, err := time.Parse(dateformat, v)
					if err != nil {
						e.Logger().Warnf(`binder: arg %v as uint: %v`, v, err)
						x = uint64(0)
					} else {
						x = uint64(t.Unix())
					}
				} else {
					x, err = strconv.ParseUint(v, 10, bitSize)
					if err != nil {
						e.Logger().Warnf(`binder: arg %v as uint: %v`, v, err)
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
						e.Logger().Warnf(`binder: struct %v invoke FromString faild`, tvf)
					}
				} else if tv.Type().String() == `time.Time` {
					x, err := time.Parse(`2006-01-02 15:04:05.000 -0700`, v)
					if err != nil {
						x, err = time.Parse(`2006-01-02 15:04:05`, v)
						if err != nil {
							x, err = time.Parse(`2006-01-02`, v)
							if err != nil {
								e.Logger().Warnf(`binder: unsupported time format %v, %v`, v, err)
							}
						}
					}
					l = x
					tv.Set(reflect.ValueOf(l))
				} else {
					e.Logger().Warn(`binder: can not set an struct which is not implement Fromconversion interface`)
				}
			case reflect.Ptr:
				e.Logger().Warn(`binder: can not set an ptr of ptr`)
			case reflect.Slice, reflect.Array:
				setSlice(e, name, tv, t)
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

func setSlice(e *Echo, fieldName string, tv reflect.Value, t []string) {

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
			err = fmt.Errorf(`binder: unsupported slice element type %v`, tk.String())
		default:
			err = fmt.Errorf(`binder: unsupported slice element type %v`, tk.String())
		}
		if err != nil {
			e.Logger().Warnf(`binder: slice error: %v, %v`, fieldName, err)
		}
	}
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

	//DateToTimestamp 日期时间转时间戳
	DateToTimestamp = func(layouts ...string) FormDataFilter {
		layout := `2006-01-02`
		if len(layouts) > 0 && len(layouts[0]) > 0 {
			layout = layouts[0]
		}
		return func(k string, v []string) (string, []string) {
			if len(v) > 0 && len(v[0]) > 0 {
				t, e := time.Parse(layout, v[0])
				if e != nil {
					log.Error(e)
					return k, []string{`0`}
				}
				return k, []string{fmt.Sprint(t.Unix())}
			}
			return k, []string{`0`}
		}
	}
)

func SetFormValue(f engine.URLValuer, fName string, index int, value interface{}) {
	if index == 0 {
		f.Set(fName, fmt.Sprint(value))
	} else {
		f.Add(fName, fmt.Sprint(value))
	}
}

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
		case `time.Time`:
			if t, y := fVal.Interface().(time.Time); y {
				dateformat := tagfast.Value(tc, fTyp, `form_format`)
				if len(dateformat) > 0 {
					f.Set(fName, t.Format(dateformat))
				} else {
					f.Set(fName, t.Format(`2006-01-02 15:04:05`))
				}
			}
		case `struct`:
			StructToForm(ctx, fVal.Interface(), fName, fieldNameFormatter)
		default:
			switch fTyp.Type.Kind() {
			case reflect.Slice:
				switch sl := fVal.Interface().(type) {
				case []uint:
					for k, v := range sl {
						SetFormValue(f, fName, k, v)
					}
				case []uint16:
					for k, v := range sl {
						SetFormValue(f, fName, k, v)
					}
				case []uint32:
					for k, v := range sl {
						SetFormValue(f, fName, k, v)
					}
				case []uint64:
					for k, v := range sl {
						SetFormValue(f, fName, k, v)
					}
				case []int:
					for k, v := range sl {
						SetFormValue(f, fName, k, v)
					}
				case []int16:
					for k, v := range sl {
						SetFormValue(f, fName, k, v)
					}
				case []int32:
					for k, v := range sl {
						SetFormValue(f, fName, k, v)
					}
				case []int64:
					for k, v := range sl {
						SetFormValue(f, fName, k, v)
					}
				case []float32:
					for k, v := range sl {
						SetFormValue(f, fName, k, v)
					}
				case []float64:
					for k, v := range sl {
						SetFormValue(f, fName, k, v)
					}
				case []string:
					for k, v := range sl {
						SetFormValue(f, fName, k, v)
					}
				case []interface{}:
					for k, v := range sl {
						SetFormValue(f, fName, k, v)
					}
				default:
					// ignore
				}
			default:
				f.Set(fName, fmt.Sprint(fVal.Interface()))
			}
		}
	}
}
