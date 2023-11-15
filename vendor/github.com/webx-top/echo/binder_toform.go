package echo

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/webx-top/echo/engine"
	"github.com/webx-top/echo/param"
	"github.com/webx-top/tagfast"
)

func SetFormValue(f engine.URLValuer, fName string, index int, value interface{}) {
	if index == 0 {
		f.Set(fName, fmt.Sprint(value))
	} else {
		f.Add(fName, fmt.Sprint(value))
	}
}

func SetFormValues(f engine.URLValuer, fName string, values []string) {
	for index, value := range values {
		if index == 0 {
			f.Set(fName, fmt.Sprint(value))
		} else {
			f.Add(fName, fmt.Sprint(value))
		}
	}
}

// FlatStructToForm 映射struct到form
func FlatStructToForm(ctx Context, m interface{}, fieldNameFormatter FieldNameFormatter, formatters ...param.StringerMap) {
	StructToForm(ctx, m, ``, fieldNameFormatter, formatters...)
}

func (e *Echo) binderValueEncode(name string, typev reflect.Type, tv reflect.Value) ([]string, error) {
	f, _ := typev.FieldByName(name)
	encoder := tagfast.Value(typev, f, `form_encoder`)
	if len(encoder) == 0 {
		return nil, ErrNotImplemented
	}
	parts := strings.SplitN(encoder, `:`, 2)
	encoder = parts[0]
	var params string
	if len(parts) == 2 {
		params = parts[1]
	}
	// ErrNotImplemented
	return e.CallBinderValueEncoder(encoder, name, tv.Interface(), params)
}

// StructToForm 映射struct到form

func StructToForm(ctx Context, m interface{}, topName string, fieldNameFormatter FieldNameFormatter, formatters ...param.StringerMap) {
	var stringers param.StringerMap // 这里的 key 为表单字段 name 属性值
	if len(formatters) > 0 {
		stringers = formatters[0]
	}
	if fieldNameFormatter == nil {
		if g, y := m.(FormNameFormatterGetter); y {
			fieldNameFormatter = g.FormNameFormatter(ctx)
		}
		if fieldNameFormatter == nil {
			fieldNameFormatter = DefaultFieldNameFormatter
		}
	}
	var valueEncoders BinderValueCustomEncoders
	if g, y := m.(ValueEncodersGetter); y {
		valueEncoders = g.ValueEncoders(ctx)
	}
	if valueEncoders == nil {
		valueEncoders = BinderValueCustomEncoders{}
	}
	if stringers == nil {
		if g, y := m.(ValueStringersGetter); y {
			stringers = g.ValueStringers(ctx)
		}
	}
	for k, f := range stringers {
		valueEncoders[k] = FormStringer(f)
	}
	structToForm(ctx, m, topName, fieldNameFormatter, valueEncoders)
}

func structToForm(ctx Context, m interface{}, topName string, fieldNameFormatter FieldNameFormatter, valueEncoders BinderValueCustomEncoders) {
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
			structToForm(ctx, srcVal.Interface(), key, fieldNameFormatter, valueEncoders)
		}
		return
	case reflect.Array, reflect.Slice:
		elems := vc.Len()
		for i := 0; i < elems; i++ {
			srcVal := vc.Index(i)
			if !srcVal.CanInterface() || srcVal.Interface() == nil {
				continue
			}
			iStr := strconv.Itoa(i)
			key := fieldNameFormatter(topName, iStr)
			structToForm(ctx, srcVal.Interface(), key, fieldNameFormatter, valueEncoders)
		}
	case reflect.Bool,
		reflect.Float32, reflect.Float64,
		reflect.Uint, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Int, reflect.Int16, reflect.Int32, reflect.Int64, reflect.String:
		if vc.CanInterface() && vc.Interface() != nil {
			ctx.Request().Form().Set(topName, fmt.Sprint(vc.Interface()))
		}
		return
	default:
		//fieldToForm(ctx, tc, reflect.StructField{}, vc, topName, fieldNameFormatter, formatter)
		//fmt.Printf("~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~>%T\n", m)
		return
	}
	for i, l := 0, tc.NumField(); i < l; i++ {
		fVal := vc.Field(i)
		fStruct := tc.Field(i)
		if fStruct.Anonymous {
			if fVal.CanInterface() {
				structToForm(ctx, fVal.Interface(), topName, fieldNameFormatter, valueEncoders)
			}
			continue
		}
		err := fieldToForm(ctx, tc, fStruct, fVal, topName, fieldNameFormatter, valueEncoders)
		if err != nil {
			fpath := fStruct.Name
			if len(topName) > 0 {
				fpath = topName + `.` + fpath
			}
			ctx.Logger().Warnf(`[StructToForm] %s: %v`, fpath, err)
		}
	}
}

func fieldToForm(ctx Context, parentTyp reflect.Type, fStruct reflect.StructField, fVal reflect.Value, topName string, fieldNameFormatter FieldNameFormatter, valueEncoders BinderValueCustomEncoders) error {
	f := ctx.Request().Form()
	fName := fieldNameFormatter(topName, fStruct.Name)
	if !fVal.CanInterface() || len(fName) == 0 {
		return nil
	}
	if valueEncoders != nil {
		encoder, ok := valueEncoders[fName]
		if ok {
			values := encoder(fVal.Interface())
			if len(values) == 0 {
				return nil
			}
			SetFormValues(f, fName, values)
			return nil
		}
	}
	if parentTyp.Kind() == reflect.Struct {
		values, err := ctx.Echo().binderValueEncode(fStruct.Name, parentTyp, fVal)
		if err != nil {
			if err != ErrNotImplemented {
				return err
			}
		} else {
			SetFormValues(f, fName, values)
			return nil
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
	default:
		switch fVal.Type().Kind() {
		case reflect.Struct:
			structToForm(ctx, fVal.Interface(), fName, fieldNameFormatter, valueEncoders)
		case reflect.Ptr:
			structToForm(ctx, fVal.Interface(), fName, fieldNameFormatter, valueEncoders)
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
			structToForm(ctx, fVal.Interface(), fName, fieldNameFormatter, valueEncoders)
		default:
			switch v := fVal.Interface().(type) {
			case ToConversion:
				f.Set(fName, v.ToString())
			default:
				f.Set(fName, fmt.Sprint(v))
			}
		}
	}
	return nil
}
