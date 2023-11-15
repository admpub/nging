package echo

import (
	"database/sql"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/webx-top/com"
	"github.com/webx-top/echo/logger"
	"github.com/webx-top/echo/param"
	"github.com/webx-top/tagfast"
)

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
	return namedStructMap(e, m, data, topName, nil, filters)
}

func NamedStructMapWithDecoder(e *Echo, m interface{}, data map[string][]string, topName string, valueDecoders BinderValueCustomDecoders, filters ...FormDataFilter) error {
	return namedStructMap(e, m, data, topName, valueDecoders, filters)
}

func FormToStruct(e *Echo, m interface{}, data map[string][]string, topName string, filters ...FormDataFilter) error {
	return namedStructMap(e, m, data, topName, nil, filters)
}

func FormToStructWithDecoder(e *Echo, m interface{}, data map[string][]string, topName string, valueDecoders BinderValueCustomDecoders, filters ...FormDataFilter) error {
	return namedStructMap(e, m, data, topName, valueDecoders, filters)
}

func FormToMap(e *Echo, m interface{}, data map[string][]string, topName string, filters ...FormDataFilter) error {
	return namedStructMap(e, m, data, topName, nil, filters)
}

func FormToMapWithDecoder(e *Echo, m interface{}, data map[string][]string, topName string, valueDecoders BinderValueCustomDecoders, filters ...FormDataFilter) error {
	return namedStructMap(e, m, data, topName, valueDecoders, filters)
}

func namedStructMap(e *Echo, m interface{}, data map[string][]string, topName string, valueDecoders BinderValueCustomDecoders, filters []FormDataFilter) error {
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
		err := e.parseFormItem(keyNormalizer, m, tc, vc, names, propPath, checkPath, values, valueDecoders, filters)
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

func (e *Echo) valueDecode(valueDecoders BinderValueCustomDecoders, propPath string, value reflect.Value, formValues []string) error {
	if valueDecoders == nil {
		return nil
	}
	decoder, ok := valueDecoders[propPath]
	if !ok {
		return nil
	}
	result, err := decoder(formValues)
	if err != nil {
		return err
	}
	vv := reflect.ValueOf(result)
	if vv.Kind() == value.Kind() && vv.Type().AssignableTo(value.Type()) {
		value.Set(vv)
		return ErrBreak
	}
	return nil
}

func (e *Echo) parseFormItem(keyNormalizer func(string) string, m interface{}, typev reflect.Type, value reflect.Value, names []string, propPath string, checkPath string, values []string, valueDecoders BinderValueCustomDecoders, filters []FormDataFilter) error {
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
		propPath += name // 结构体字段或map的key层级路径

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
				err = e.setMap(e.Logger(), tc, vc, name, value, typev, propPath, values)
			case reflect.Struct:
				err = e.setStructField(e.Logger(), tc, vc, name, value, typev, propPath, values, valueDecoders)
			default:
				e.Logger().Debugf(`binder: The last layer field "%v" does not support type: %v`, propPath, value.Kind())
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
			pos := strings.LastIndex(propPath, `.`)
			if pos > -1 {
				propPath = propPath[0:pos] // 忽略切片数字下标
			}
			return e.parseFormItem(keyNormalizer, m, newT, newV, names[i+1:], propPath+`.`, checkPath+`.`, values, valueDecoders, filters)
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
			return e.parseFormItem(keyNormalizer, m, newT, newV, names[i+1:], propPath+`.`, checkPath+`.`, values, valueDecoders, filters)
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
				return e.parseFormItem(keyNormalizer, m, value.Type(), value, names[i+1:], propPath+`.`, checkPath+`.`, values, valueDecoders, filters)
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

func (e *Echo) setStructField(logger logger.Logger,
	parentT reflect.Type, parentV reflect.Value, name string,
	value reflect.Value, typev reflect.Type,
	propPath string, values []string, valueDecoders BinderValueCustomDecoders) error {
	tv := SafeGetFieldByName(value, name)
	if !tv.IsValid() {
		return ErrBreak
	}
	if !tv.CanSet() {
		logger.Warnf(`binder: can not set %v=%+v to %v`, propPath, values, tv.Kind())
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
	if decErr := e.valueDecode(valueDecoders, propPath, tv, values); decErr != nil {
		if decErr == ErrBreak {
			decErr = nil
		}
		return decErr
	}
	if typev.Kind() == reflect.Struct {
		err := e.binderValueDecode(name, typev, tv, values)
		if err == nil || err != ErrNotImplemented {
			return err
		}
	}
	return setField(logger, parentT, tv, f, name, values)
}

func (e *Echo) binderValueDecode(name string, typev reflect.Type, tv reflect.Value, values []string) error {
	f, _ := typev.FieldByName(name)
	decoder := tagfast.Value(typev, f, `form_decoder`)
	if len(decoder) == 0 {
		return ErrNotImplemented
	}
	parts := strings.SplitN(decoder, `:`, 2)
	decoder = parts[0]
	var params string
	if len(parts) == 2 {
		params = parts[1]
	}
	result, err := e.CallBinderValueDecoder(decoder, name, values, params)
	if err != nil { // ErrNotImplemented
		return err
	}
	vv := reflect.ValueOf(result)
	if vv.Kind() == tv.Kind() && vv.Type().AssignableTo(tv.Type()) {
		tv.Set(vv)
		return nil
	}
	return ErrNotImplemented
}

func (e *Echo) setMap(logger logger.Logger,
	parentT reflect.Type, parentV reflect.Value,
	name string, value reflect.Value, typev reflect.Type,
	propPath string, values []string) error {
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
				return errors.New(`binder: [setMap] unsupported type ` + oldVal.Kind().String() + `: ` + name)
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
			seperator := tagfast.Value(parentT, f, `form_seperator`)
			if len(seperator) == 0 {
				seperator = tagfast.Value(parentT, f, `form_delimiter`)
			}
			if len(seperator) > 0 {
				v = strings.Join(values, seperator)
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
		} else if len(v) > 0 {
			x, err := strconv.Atoi(v)
			if err != nil {
				logger.Warnf(`binder: arg %q as int: %v`, v, err)
			}
			l = x
		} else {
			l = int(0)
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
			} else if len(v) > 0 {
				x, err := strconv.ParseInt(v, 10, 64)
				if err != nil {
					logger.Warnf(`binder: arg %q as int64: %v`, v, err)
				}
				l = x
			} else {
				l = int64(0)
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
		case reflect.Uint:
			bitSize = 0
		case reflect.Uint32:
			bitSize = 32
		default:
			bitSize = 64
		}
		if len(dateformat) > 0 {
			t, err := time.ParseInLocation(dateformat, v, time.Local)
			if err != nil {
				logger.Warnf(`binder: arg %q as time: %v`, v, err)
				x = uint64(0)
			} else {
				x = uint64(t.Unix())
			}
		} else if len(v) > 0 {
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
		seperator := tagfast.Value(parentT, f, `form_seperator`)
		if len(seperator) == 0 {
			seperator = tagfast.Value(parentT, f, `form_delimiter`)
		}
		if len(seperator) > 0 {
			var parts []string
			for _, value := range values {
				value = strings.TrimSpace(value)
				if len(value) == 0 {
					continue
				}
				parts = append(parts, strings.Split(value, seperator)...)
			}
			setSlice(logger, name, tv, parts)
		} else {
			setSlice(logger, name, tv, values)
		}
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
