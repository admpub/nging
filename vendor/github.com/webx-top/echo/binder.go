package echo

import (
	"database/sql"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/admpub/copier"
	"github.com/admpub/log"

	"github.com/webx-top/echo/engine"
	"github.com/webx-top/echo/param"
	"github.com/webx-top/tagfast"
)

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
		Echo:     e,
		decoders: DefaultBinderDecoders,
	}
}

func (b *binder) MustBind(i interface{}, c Context, filter ...FormDataFilter) error {
	contentType := c.Request().Header().Get(HeaderContentType)
	contentType = strings.ToLower(strings.TrimSpace(strings.SplitN(contentType, `;`, 2)[0]))
	if decoder, ok := b.decoders[contentType]; ok {
		return decoder(i, c, filter...)
	}
	if decoder, ok := b.decoders[`*`]; ok {
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
	for key, values := range data {

		if len(topName) > 0 {
			if !strings.HasPrefix(key, topName) {
				continue
			}
			key = key[len(topName)+1:]
		}

		names := strings.Split(key, `.`)
		var propPath, checkPath string
		if len(names) == 1 && strings.HasSuffix(key, `]`) {
			key = strings.TrimSuffix(key, `[]`)
			names = FormNames(key)
		}
		err := parseFormItem(e, m, tc, vc, names, propPath, checkPath, key, values, filters...)
		if err == nil {
			continue
		}
		if err == ErrBreak {
			err = nil
			break
		}
		return err
	}
	return nil
}

func parseFormItem(e *Echo, m interface{}, typev reflect.Type, value reflect.Value, names []string, propPath string, checkPath string, key string, values []string, filters ...FormDataFilter) error {
	length := len(names)
	vc := value
	tc := typev
	for i, name := range names {
		name = strings.Title(name)
		if i > 0 {
			propPath += `.`
			checkPath += `.`
		}
		propPath += name

		// check
		vk := checkPath + name // example: Name or *.Name or *.*.Password
		for _, filter := range filters {
			vk, values = filter(vk, values)
			if len(vk) == 0 || len(values) == 0 {
				e.Logger().Debugf(`binder: skip %v%v (%v) => %v`, checkPath, name, propPath, values)
				return nil
			}
		}
		checkPath += `*`

		//最后一个元素
		if i == length-1 {
			err := setField(e, tc, vc, key, name, value, typev, values)
			if err == nil {
				continue
			}
			if err == ErrBreak {
				return nil
			}
			return err
		}

		//不是最后一个元素
		switch value.Kind() {
		case reflect.Slice:
			index, err := strconv.Atoi(name)
			if err != nil {
				e.Logger().Warnf(`binder: can not convert index number %T#%v -> %v`, m, propPath, err.Error())
				return nil
			}
			if e.FormSliceMaxIndex > 0 && index > e.FormSliceMaxIndex {
				return fmt.Errorf(`%w, greater than %d`, ErrSliceIndexTooLarge, e.FormSliceMaxIndex)
			}
			if value.IsNil() {
				value.Set(reflect.MakeSlice(value.Type(), 1, 1))
			}
			itemT := value.Type()
			if itemT.Kind() == reflect.Ptr {
				itemT = itemT.Elem()
				value = value.Elem()
			}
			itemT = itemT.Elem()
			if index >= value.Len() {
				for i := value.Len(); i <= index; i++ {
					tempv := reflect.New(itemT)
					value.Set(reflect.Append(value, tempv.Elem()))
				}
			}
			newV := value.Index(index)
			newT := newV.Type()
			switch newT.Kind() {
			case reflect.Struct:
			case reflect.Ptr:
				newT = newT.Elem()
				if newV.IsNil() {
					newV.Set(reflect.New(newT))
				}
				newV = newV.Elem()
			default:
				return errors.New(`binder: unsupported type ` + tc.Kind().String())
			}
			return parseFormItem(e, m, newT, newV, names[i+1:], propPath+`.`, checkPath+`.`, key, values, filters...)
		case reflect.Map:
			if value.IsNil() {
				value.Set(reflect.MakeMap(value.Type()))
			}
			itemT := value.Type()
			if itemT.Kind() == reflect.Ptr {
				itemT = itemT.Elem()
				value = value.Elem()
			}
			itemT = itemT.Elem()
			index := reflect.ValueOf(name)
			newV := value.MapIndex(index)
			if !newV.IsValid() {
				newV = reflect.New(itemT).Elem()
				value.SetMapIndex(index, newV)
			}
			newT := newV.Type()
			switch newT.Kind() {
			case reflect.Struct:
			case reflect.Ptr:
				newT = newT.Elem()
				if newV.IsNil() {
					newV = reflect.New(newT)
					value.SetMapIndex(index, newV)
				}
				newV = newV.Elem()
			default:
				return errors.New(`binder: unsupported type ` + tc.Kind().String())
			}
			return parseFormItem(e, m, newT, newV, names[i+1:], propPath+`.`, checkPath+`.`, key, values, filters...)
		case reflect.Struct:
			f, _ := typev.FieldByName(name)
			if tagfast.Value(tc, f, `form_options`) == `-` {
				return nil
			}
			value = value.FieldByName(name)
			if !value.IsValid() {
				e.Logger().Debugf(`binder: %T#%v value is not valid %v`, m, propPath, value)
				return nil
			}
			if !value.CanSet() {
				e.Logger().Warnf(`binder: can not set %T#%v -> %v`, m, propPath, value.Interface())
				return nil
			}
			if value.Kind() == reflect.Ptr {
				if value.IsNil() {
					value.Set(reflect.New(value.Type().Elem()))
				}
				value = value.Elem()
			}
			switch value.Kind() {
			case reflect.Struct:
			case reflect.Slice, reflect.Map:
				return parseFormItem(e, m, value.Type(), value, names[i+1:], propPath+`.`, checkPath+`.`, key, values, filters...)
			default:
				e.Logger().Warnf(`binder: arg error, value %T#%v kind is %v`, m, propPath, value.Kind())
				return nil
			}
			typev = value.Type()
			f, _ = typev.FieldByName(name)
			if tagfast.Value(tc, f, `form_options`) == `-` {
				return nil
			}
		default:
			e.Logger().Warnf(`binder: arg error, value kind is %v`, value.Kind())
			return nil
		}
	}
	return nil
}

var (
	ErrBreak              = errors.New("[BREAK]")
	ErrContinue           = errors.New("[CONTINUE]")
	ErrExit               = errors.New("[EXIT]")
	ErrReturn             = errors.New("[RETURN]")
	ErrSliceIndexTooLarge = errors.New("The slice index value of the form field is too large")
)

func SafeGetFieldByName(parentT reflect.Type, parentV reflect.Value, name string, value reflect.Value) (v reflect.Value) {
	defer func() {
		if r := recover(); r != nil {
			switch fmt.Sprint(r) {
			case `reflect: indirection through nil pointer to embedded struct`:
				copier.InitNilFields(parentT, parentV, ``, copier.AllNilFields)
				v = value.FieldByName(name)
			default:
				panic(r)
			}
		}
	}()
	v = value.FieldByName(name)
	return
}

func setField(e *Echo, parentT reflect.Type, parentV reflect.Value, k string, name string, value reflect.Value, typev reflect.Type, values []string) error {
	tv := SafeGetFieldByName(parentT, parentV, name, value)
	if !tv.IsValid() {
		return ErrBreak
	}
	if !tv.CanSet() {
		e.Logger().Warnf(`binder: can not set %v to %v`, k, tv)
		return ErrBreak
	}
	f, _ := typev.FieldByName(name)
	if tagfast.Value(parentT, f, `form_options`) == `-` {
		return ErrBreak
	}
	if tv.Kind() == reflect.Ptr {
		tv.Set(reflect.New(tv.Type().Elem()))
		tv = tv.Elem()
	}
	v := values[0]
	var l interface{}
	switch kind := tv.Kind(); kind {
	case reflect.String:
		switch tagfast.Value(parentT, f, `form_filter`) {
		case `html`:
			v = DefaultHTMLFilter(v)
		default:
			delimter := tagfast.Value(parentT, f, `form_delimiter`)
			if len(delimter) > 0 {
				v = strings.Join(values, delimter)
			}
		}
		l = v
		tv.Set(reflect.ValueOf(l))
	case reflect.Bool:
		l = (v != `false` && v != `0` && v != ``)
		tv.Set(reflect.ValueOf(l))
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32:
		dateformat := tagfast.Value(parentT, f, `form_format`)
		if len(dateformat) > 0 {
			t, err := time.ParseInLocation(dateformat, v, time.Local)
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
		dateformat := tagfast.Value(parentT, f, `form_format`)
		if len(dateformat) > 0 {
			t, err := time.ParseInLocation(dateformat, v, time.Local)
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
		dateformat := tagfast.Value(parentT, f, `form_format`)
		var x uint64
		var bitSize int
		switch kind {
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
			t, err := time.ParseInLocation(dateformat, v, time.Local)
			if err != nil {
				e.Logger().Warnf(`binder: arg %v as uint: %v`, v, err)
				x = uint64(0)
			} else {
				x = uint64(t.Unix())
			}
		} else {
			var err error
			x, err = strconv.ParseUint(v, 10, bitSize)
			if err != nil {
				e.Logger().Warnf(`binder: arg %v as uint: %v`, v, err)
			}
		}
		switch kind {
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
		switch rawType := tv.Interface().(type) {
		case FromConversion:
			if err := rawType.FromString(v); err != nil {
				e.Logger().Warnf(`binder: struct %v invoke FromString faild`, rawType)
			}
		case time.Time:
			x, err := time.ParseInLocation(`2006-01-02 15:04:05.000 -0700`, v, time.Local)
			if err != nil {
				x, err = time.ParseInLocation(`2006-01-02 15:04:05`, v, time.Local)
				if err != nil {
					x, err = time.ParseInLocation(`2006-01-02`, v, time.Local)
					if err != nil {
						e.Logger().Warnf(`binder: unsupported time format %v, %v`, v, err)
					}
				}
			}
			l = x
			tv.Set(reflect.ValueOf(l))
		default:
			if scanner, ok := tv.Addr().Interface().(sql.Scanner); ok {
				if err := scanner.Scan(values[0]); err != nil {
					e.Logger().Warnf(`binder: struct %v invoke Scan faild`, rawType)
				}
			}
		}
	case reflect.Ptr:
		e.Logger().Warn(`binder: can not set an ptr of ptr`)
	case reflect.Slice, reflect.Array:
		setSlice(e, name, tv, values)
	default:
		return ErrBreak
	}

	//validation
	valid := tagfast.Value(parentT, f, `valid`)
	if len(valid) == 0 {
		return nil
	}
	result := e.Validator.Validate(name, fmt.Sprintf(`%v`, l), valid)
	return result.Error()
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
				t, e := time.ParseInLocation(layout, v[0], time.Local)
				if e != nil {
					log.Error(e)
					return k, []string{`0`}
				}
				return k, []string{fmt.Sprint(t.Unix())}
			}
			return k, []string{`0`}
		}
	}
	//TimestampToDate 时间戳转日期时间
	TimestampToDate = func(layouts ...string) FormDataFilter {
		layout := `2006-01-02 15:04:05`
		if len(layouts) > 0 && len(layouts[0]) > 0 {
			layout = layouts[0]
		}
		return func(k string, v []string) (string, []string) {
			if len(v) > 0 && len(v[0]) > 0 {
				tsi := strings.SplitN(v[0], `.`, 2)
				var sec, nsec int64
				switch len(tsi) {
				case 2:
					nsec = param.AsInt64(tsi[1])
					fallthrough
				case 1:
					sec = param.AsInt64(tsi[0])
				}
				t := time.Unix(sec, nsec)
				if t.IsZero() {
					return k, []string{``}
				}
				return k, []string{t.Format(layout)}
			}
			return k, v
		}
	}
	//JoinValues 组合数组为字符串
	JoinValues = func(seperators ...string) FormDataFilter {
		sep := `,`
		if len(seperators) > 0 {
			sep = seperators[0]
		}
		return func(k string, v []string) (string, []string) {
			return k, []string{strings.Join(v, sep)}
		}
	}
	//SplitValues 拆分字符串为数组
	SplitValues = func(seperators ...string) FormDataFilter {
		sep := `,`
		if len(seperators) > 0 {
			sep = seperators[0]
		}
		return func(k string, v []string) (string, []string) {
			if len(v) > 0 && len(v[0]) > 0 {
				v = strings.Split(v[0], sep)
			}
			return k, v
		}
	}
	TimestampStringer  = param.TimestampStringer
	DateTimeStringer   = param.DateTimeStringer
	WhitespaceStringer = param.WhitespaceStringer
	Ignored            = param.Ignored
)

func TranslateStringer(t Translator, args ...interface{}) param.Stringer {
	return param.StringerFunc(func(v interface{}) string {
		return t.T(param.AsString(v), args...)
	})
}

//FormatFieldValue 格式化字段值
func FormatFieldValue(formatters map[string]FormDataFilter) FormDataFilter {
	newFormatters := map[string]FormDataFilter{}
	for k, v := range formatters {
		newFormatters[strings.Title(k)] = v
	}
	return func(k string, v []string) (string, []string) {
		tk := strings.Title(k)
		if formatter, ok := newFormatters[tk]; ok {
			return formatter(k, v)
		}
		return k, v
	}
}

//IncludeFieldName 包含字段
func IncludeFieldName(fieldNames ...string) FormDataFilter {
	for k, v := range fieldNames {
		fieldNames[k] = strings.Title(v)
	}
	return func(k string, v []string) (string, []string) {
		tk := strings.Title(k)
		for _, fv := range fieldNames {
			if fv == tk {
				return k, v
			}
		}
		return ``, v
	}
}

//ExcludeFieldName 排除字段
func ExcludeFieldName(fieldNames ...string) FormDataFilter {
	for k, v := range fieldNames {
		fieldNames[k] = strings.Title(v)
	}
	return func(k string, v []string) (string, []string) {
		tk := strings.Title(k)
		for _, fv := range fieldNames {
			if fv == tk {
				return ``, v
			}
		}
		return k, v
	}
}

func SetFormValue(f engine.URLValuer, fName string, index int, value interface{}) {
	if index == 0 {
		f.Set(fName, fmt.Sprint(value))
	} else {
		f.Add(fName, fmt.Sprint(value))
	}
}

//FlatStructToForm 映射struct到form
func FlatStructToForm(ctx Context, m interface{}, topName string, fieldNameFormatter FieldNameFormatter, formatters ...param.StringerMap) {
	StructToForm(ctx, m, ``, fieldNameFormatter, formatters...)
}

//StructToForm 映射struct到form
func StructToForm(ctx Context, m interface{}, topName string, fieldNameFormatter FieldNameFormatter, formatters ...param.StringerMap) {
	var formatter param.StringerMap
	if len(formatters) > 0 {
		formatter = formatters[0]
	}
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
		if formatter != nil {
			result, found, skip := formatter.String(fName, fVal.Interface())
			if skip {
				continue
			}
			if found {
				f.Set(fName, result)
				continue
			}
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
				switch v := fVal.Interface().(type) {
				case ToConversion:
					f.Set(fName, v.ToString())
				default:
					f.Set(fName, fmt.Sprint(v))
				}
			}
		}
	}
}
