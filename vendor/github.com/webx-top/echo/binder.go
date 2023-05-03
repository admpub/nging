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

	"github.com/admpub/log"

	"github.com/webx-top/com"
	"github.com/webx-top/echo/engine"
	"github.com/webx-top/echo/logger"
	"github.com/webx-top/echo/param"
	"github.com/webx-top/tagfast"
)

type (
	// Binder is the interface that wraps the Bind method.
	Binder interface {
		Bind(interface{}, Context, ...FormDataFilter) error
		BindAndValidate(interface{}, Context, ...FormDataFilter) error
		MustBind(interface{}, Context, ...FormDataFilter) error
		MustBindAndValidate(interface{}, Context, ...FormDataFilter) error
	}
	binder struct {
		decoders map[string]func(interface{}, Context, ...FormDataFilter) error
	}
)

func NewBinder(e *Echo) Binder {
	return &binder{
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

func (b *binder) MustBindAndValidate(i interface{}, c Context, filter ...FormDataFilter) error {
	if err := b.MustBind(i, c, filter...); err != nil {
		return err
	}
	return ValidateStruct(c, i)
}

func (b *binder) Bind(i interface{}, c Context, filter ...FormDataFilter) (err error) {
	err = b.MustBind(i, c, filter...)
	if err == ErrUnsupportedMediaType {
		err = nil
	}
	return
}

func (b *binder) BindAndValidate(i interface{}, c Context, filter ...FormDataFilter) error {
	if err := b.MustBind(i, c, filter...); err != nil {
		if err != ErrUnsupportedMediaType {
			return err
		}
	}
	return ValidateStruct(c, i)
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

type BinderFormTopNamer interface {
	BinderFormTopName() string
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
	keyNormalizer := com.Title
	if bkn, ok := m.(BinderKeyNormalizer); ok {
		keyNormalizer = bkn.BinderKeyNormalizer
	}
	topNameLen := len(topName)
	if topNameLen == 0 {
		if topNamer, ok := m.(BinderFormTopNamer); ok {
			topName = topNamer.BinderFormTopName()
			topNameLen = len(topName)
		}
	}
	for key, values := range data {

		if topNameLen > 0 {
			if topNameLen+1 >= len(key) { //key = topName.field
				continue
			}
			if key[0:topNameLen] != topName {
				continue
			}
			key = key[topNameLen:]
			if key[0] == '.' {
				if len(key) <= 1 {
					continue
				}
				key = key[1:]
			} else if key[0] == '[' {
				if len(key) <= 1 {
					continue
				}
				key = key[1:]
			}
		}

		names := strings.Split(key, `.`)
		var propPath, checkPath string
		if len(names) == 1 && strings.HasSuffix(key, `]`) {
			key = strings.TrimSuffix(key, `[]`)
			if len(key) == 0 {
				continue
			}
			names = FormNames(key)
		}
		err := parseFormItem(keyNormalizer, e, m, tc, vc, names, propPath, checkPath, key, values, filters...)
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

type BinderKeyNormalizer interface {
	BinderKeyNormalizer(string) string
}

func parseFormItem(keyNormalizer func(string) string, e *Echo, m interface{}, typev reflect.Type, value reflect.Value, names []string, propPath string, checkPath string, key string, values []string, filters ...FormDataFilter) error {
	length := len(names)
	vc := value
	tc := typev
	isMap := value.Kind() == reflect.Map
	for i, name := range names {
		if !isMap {
			name = keyNormalizer(name)
		}
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
			var err error
			switch value.Kind() {
			case reflect.Map:
				err = setMap(e.Logger(), tc, vc, key, name, value, typev, values)
			default:
				err = setStructField(e.Logger(), tc, vc, key, name, value, typev, values)
			}
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
			case reflect.Map:
				if newV.IsNil() {
					newV.Set(reflect.MakeMap(newT))
				}
			default:
				return errors.New(`binder: [parseFormItem#slice] unsupported type ` + tc.Kind().String() + `: ` + propPath)
			}
			return parseFormItem(keyNormalizer, e, m, newT, newV, names[i+1:], propPath+`.`, checkPath+`.`, key, values, filters...)
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
			case reflect.Map:
				if newV.IsNil() {
					newV = reflect.MakeMap(newT)
					value.SetMapIndex(index, newV)
				}
			default:
				return errors.New(`binder: [parseFormItem#map] unsupported type ` + tc.Kind().String() + `: ` + propPath)
			}
			return parseFormItem(keyNormalizer, e, m, newT, newV, names[i+1:], propPath+`.`, checkPath+`.`, key, values, filters...)
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
				return parseFormItem(keyNormalizer, e, m, value.Type(), value, names[i+1:], propPath+`.`, checkPath+`.`, key, values, filters...)
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
	ErrSliceIndexTooLarge = errors.New("the slice index value of the form field is too large")
)

func SafeGetFieldByName(value reflect.Value, name string) (v reflect.Value) {
	var destFieldNotSet bool
	if f, ok := value.Type().FieldByName(name); ok {
		// only initialize parent embedded struct pointer in the path
		for idx := range f.Index[:len(f.Index)-1] {
			destField := value.FieldByIndex(f.Index[:idx+1])

			if destField.Kind() != reflect.Ptr {
				continue
			}

			if !destField.IsNil() {
				continue
			}
			if !destField.CanSet() {
				destFieldNotSet = true
				break
			}

			// destField is a nil pointer that can be set
			newValue := reflect.New(destField.Type().Elem())
			destField.Set(newValue)
		}
	}

	if destFieldNotSet {
		return
	}
	v = value.FieldByName(name)
	return
}

func setStructField(logger logger.Logger, parentT reflect.Type, parentV reflect.Value, k string, name string, value reflect.Value, typev reflect.Type, values []string) error {
	tv := SafeGetFieldByName(value, name)
	if !tv.IsValid() {
		return ErrBreak
	}
	if !tv.CanSet() {
		logger.Warnf(`binder: can not set %v to %v`, k, tv)
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
	return setField(logger, parentT, tv, f, name, values)
}

func setMap(logger logger.Logger, parentT reflect.Type, parentV reflect.Value, k string, name string, value reflect.Value, typev reflect.Type, values []string) error {
	if value.IsNil() {
		value.Set(reflect.MakeMap(value.Type()))
	}
	value = reflect.Indirect(value)
	var index reflect.Value
	if typev.Key().Kind() != reflect.String {
		mapKey := param.AsType(typev.Key().Kind().String(), name)
		index = reflect.ValueOf(mapKey)
	} else {
		index = reflect.ValueOf(name)
	}
	if index.Type().Name() != typev.Key().Name() && index.CanConvert(typev.Key()) {
		index = index.Convert(typev.Key())
	}
	oldVal := value.MapIndex(index)
	if !oldVal.IsValid() {
		oldType := value.Type().Elem()
		oldVal = reflect.New(oldType).Elem()
		switch oldVal.Kind() {
		case reflect.String:
			value.SetMapIndex(index, reflect.ValueOf(values[0]))
		case reflect.Interface:
			if len(values) > 1 {
				value.SetMapIndex(index, reflect.ValueOf(values))
			} else {
				value.SetMapIndex(index, reflect.ValueOf(values[0]))
			}
		case reflect.Slice:
			setSlice(logger, name, oldVal, values)
			value.SetMapIndex(index, oldVal)
		default:
			mapVal := param.AsType(oldVal.Kind().String(), values[0])
			converted := !reflect.DeepEqual(mapVal, oldVal.Interface())
			if converted {
				originalType := oldVal.Type()
				oldVal = reflect.ValueOf(mapVal)
				if oldVal.Type().Name() != originalType.Name() && oldVal.CanConvert(originalType) {
					oldVal = oldVal.Convert(originalType)
				}
			}
			value.SetMapIndex(index, oldVal)
			if !converted {
				return errors.New(`binder: [setStructField] unsupported type ` + oldVal.Kind().String() + `: ` + name)
			}
		}
		return nil
	}
	if oldVal.Type().Kind() == reflect.Interface {
		oldVal = reflect.Indirect(reflect.ValueOf(oldVal.Interface()))
	}
	isPtr := oldVal.CanAddr()
	if !isPtr {
		oldVal = reflect.New(oldVal.Type())
	}
	err := setField(logger, parentT, oldVal.Elem(), reflect.StructField{Name: name}, name, values)
	if err == nil {
		if !isPtr {
			oldVal = reflect.Indirect(oldVal)
		}
		value.SetMapIndex(index, oldVal)
	}
	return err
}

func SetReflectValue(source interface{}, dest reflect.Value) bool {
	fv := reflect.ValueOf(source)
	destT := dest.Type()
	if destT.Name() == fv.Type().Name() {
		dest.Set(fv)
		return true
	}
	if fv.CanConvert(destT) {
		fv = fv.Convert(destT)
		dest.Set(fv)
		return true
	}
	return false
}

func setField(logger logger.Logger, parentT reflect.Type, tv reflect.Value, f reflect.StructField, name string, values []string) error {
	v := values[0]
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
		SetReflectValue(v, tv)
	case reflect.Bool:
		ok, _ := strconv.ParseBool(v)
		SetReflectValue(ok, tv)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32:
		var l interface{}
		dateformat := tagfast.Value(parentT, f, `form_format`)
		if len(dateformat) > 0 {
			t, err := time.ParseInLocation(dateformat, v, time.Local)
			if err != nil {
				logger.Warnf(`binder: arg %q as int: %v`, v, err)
				l = int(0)
			} else {
				l = int(t.Unix())
			}
		} else {
			x, err := strconv.Atoi(v)
			if err != nil {
				logger.Warnf(`binder: arg %q as int: %v`, v, err)
			}
			l = x
		}
		SetReflectValue(l, tv)
	case reflect.Int64:
		var l interface{}
		switch tv.Interface().(type) {
		case time.Duration:
			l, _ = time.ParseDuration(v)
		default:
			dateformat := tagfast.Value(parentT, f, `form_format`)
			if len(dateformat) > 0 {
				t, err := time.ParseInLocation(dateformat, v, time.Local)
				if err != nil {
					logger.Warnf(`binder: arg %q as int64: %v`, v, err)
					l = int64(0)
				} else {
					l = t.Unix()
				}
			} else {
				x, err := strconv.ParseInt(v, 10, 64)
				if err != nil {
					logger.Warnf(`binder: arg %q as int64: %v`, v, err)
				}
				l = x
			}
		}
		SetReflectValue(l, tv)
	case reflect.Float32, reflect.Float64:
		x, err := strconv.ParseFloat(v, 64)
		if err != nil {
			logger.Warnf(`binder: arg %q as float64: %v`, v, err)
		}
		SetReflectValue(x, tv)
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
				logger.Warnf(`binder: arg %q as uint: %v`, v, err)
				x = uint64(0)
			} else {
				x = uint64(t.Unix())
			}
		} else {
			var err error
			x, err = strconv.ParseUint(v, 10, bitSize)
			if err != nil {
				logger.Warnf(`binder: arg %q as uint: %v`, v, err)
			}
		}
		var l interface{}
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
		SetReflectValue(l, tv)
	case reflect.Struct:
		switch rawType := tv.Interface().(type) {
		case FromConversion:
			if err := rawType.FromString(v); err != nil {
				logger.Warnf(`binder: struct %v invoke FromString faild`, rawType)
			}
		case time.Time:
			x, err := time.ParseInLocation(`2006-01-02 15:04:05.000 -0700`, v, time.Local)
			if err != nil {
				x, err = time.ParseInLocation(`2006-01-02 15:04:05`, v, time.Local)
				if err != nil {
					x, err = time.ParseInLocation(`2006-01-02`, v, time.Local)
					if err != nil {
						logger.Warnf(`binder: unsupported time format %v, %v`, v, err)
					}
				}
			}
			SetReflectValue(x, tv)
		default:
			if scanner, ok := tv.Addr().Interface().(sql.Scanner); ok {
				if err := scanner.Scan(values[0]); err != nil {
					logger.Warnf(`binder: struct %v invoke Scan faild`, rawType)
				}
			}
		}
	case reflect.Ptr:
		setField(logger, parentT, tv.Elem(), f, name, values)
	case reflect.Slice, reflect.Array:
		setSlice(logger, name, tv, values)
	default:
		return ErrBreak
	}
	return nil
}

func setSlice(logger logger.Logger, fieldName string, tv reflect.Value, t []string) {

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
		case reflect.Interface:
			tv.Index(i).Set(reflect.ValueOf(s))
		case reflect.Complex64, reflect.Complex128:
			// TODO:
			err = fmt.Errorf(`binder: unsupported slice element type %v`, tk.String())
		default:
			err = fmt.Errorf(`binder: unsupported slice element type %v`, tk.String())
		}
		if err != nil {
			logger.Warnf(`binder: slice error: %v, %v`, fieldName, err)
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
	//ArrayFieldNameFormatter 格式化函数(struct->form)
	ArrayFieldNameFormatter FieldNameFormatter = func(topName, fieldName string) string {
		var fName string
		if len(topName) == 0 {
			fName = fieldName
		} else {
			fName = topName + `[` + fieldName + `]`
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

// FormatFieldValue 格式化字段值
func FormatFieldValue(formatters map[string]FormDataFilter, keyNormalizerArg ...func(string) string) FormDataFilter {
	newFormatters := map[string]FormDataFilter{}
	keyNormalizer := strings.Title
	if len(keyNormalizerArg) > 0 && keyNormalizerArg[0] != nil {
		keyNormalizer = keyNormalizerArg[0]
	}
	for k, v := range formatters {
		newFormatters[keyNormalizer(k)] = v
	}
	return func(k string, v []string) (string, []string) {
		tk := keyNormalizer(k)
		if formatter, ok := newFormatters[tk]; ok {
			return formatter(k, v)
		}
		return k, v
	}
}

// IncludeFieldName 包含字段
func IncludeFieldName(fieldNames ...string) FormDataFilter {
	for k, v := range fieldNames {
		fieldNames[k] = com.Title(v)
	}
	return func(k string, v []string) (string, []string) {
		tk := com.Title(k)
		for _, fv := range fieldNames {
			if fv == tk {
				return k, v
			}
		}
		return ``, v
	}
}

// ExcludeFieldName 排除字段
func ExcludeFieldName(fieldNames ...string) FormDataFilter {
	for k, v := range fieldNames {
		fieldNames[k] = com.Title(v)
	}
	return func(k string, v []string) (string, []string) {
		tk := com.Title(k)
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

// FlatStructToForm 映射struct到form
func FlatStructToForm(ctx Context, m interface{}, fieldNameFormatter FieldNameFormatter, formatters ...param.StringerMap) {
	StructToForm(ctx, m, ``, fieldNameFormatter, formatters...)
}

// StructToForm 映射struct到form
func StructToForm(ctx Context, m interface{}, topName string, fieldNameFormatter FieldNameFormatter, formatters ...param.StringerMap) {
	var formatter param.StringerMap
	if len(formatters) > 0 {
		formatter = formatters[0]
	}
	if fieldNameFormatter == nil {
		fieldNameFormatter = DefaultFieldNameFormatter
	}
	vc := reflect.ValueOf(m)
	tc := reflect.TypeOf(m)
	if tc.Kind() == reflect.Ptr {
		tc = tc.Elem()
		if vc.IsNil() {
			return
		}
		vc = vc.Elem()
	}
	switch tc.Kind() {
	case reflect.Struct:
	case reflect.Map:
		for _, srcKey := range vc.MapKeys() {
			srcVal := vc.MapIndex(srcKey)
			if !srcVal.CanInterface() || srcVal.Interface() == nil {
				continue
			}
			key := fieldNameFormatter(topName, srcKey.String())
			switch srcVal.Kind() {
			case reflect.Ptr:
				StructToForm(ctx, srcVal.Interface(), key, fieldNameFormatter, formatters...)
			case reflect.Struct:
				StructToForm(ctx, srcVal.Interface(), key, fieldNameFormatter, formatters...)
			default:
				fieldToForm(ctx, tc, reflect.StructField{Name: srcKey.String()}, srcVal, topName, fieldNameFormatter, formatter)
			}
		}
		return
	default:
		//fieldToForm(ctx, tc, reflect.StructField{}, vc, topName, fieldNameFormatter, formatter)
		return
	}
	for i, l := 0, tc.NumField(); i < l; i++ {
		fVal := vc.Field(i)
		fStruct := tc.Field(i)
		fieldToForm(ctx, tc, fStruct, fVal, topName, fieldNameFormatter, formatter)
	}
}

func fieldToForm(ctx Context, parentTyp reflect.Type, fStruct reflect.StructField, fVal reflect.Value, topName string, fieldNameFormatter FieldNameFormatter, formatter param.StringerMap) {
	f := ctx.Request().Form()
	fName := fieldNameFormatter(topName, fStruct.Name)
	if !fVal.CanInterface() || len(fName) == 0 {
		return
	}
	if formatter != nil {
		result, found, skip := formatter.String(fName, fVal.Interface())
		if skip {
			return
		}
		if found {
			f.Set(fName, result)
			return
		}
	}
	switch fVal.Type().String() {
	case `time.Time`:
		if t, y := fVal.Interface().(time.Time); y {
			if t.IsZero() {
				f.Set(fName, ``)
			} else {
				dateformat := tagfast.Value(parentTyp, fStruct, `form_format`)
				if len(dateformat) > 0 {
					f.Set(fName, t.Format(dateformat))
				} else {
					f.Set(fName, t.Format(`2006-01-02 15:04:05`))
				}
			}
		}
	case `time.Duration`:
		if t, y := fVal.Interface().(time.Duration); y {
			f.Set(fName, t.String())
		}
	case `struct`:
		StructToForm(ctx, fVal.Interface(), fName, fieldNameFormatter)
	default:
		switch fVal.Type().Kind() {
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
		case reflect.Map:
			StructToForm(ctx, fVal.Interface(), fName, fieldNameFormatter, formatter)
		case reflect.Ptr:
			StructToForm(ctx, fVal.Interface(), fName, fieldNameFormatter)
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
