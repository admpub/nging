package map2struct

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/admpub/copier"
)

func Scan(sv interface{}, m map[string]interface{}, tagNames ...string) (err error) {
	val := reflect.ValueOf(sv)
	if val.Kind() != reflect.Ptr {
		return errors.New("non-pointer passed to Unmarshal")
	}
	v := val.Elem()

	typ := v.Type()
	if v.Kind() != reflect.Struct {
		return errors.New("paramter error")
	}

	fields := copier.DeepFindFields(typ, v, ``, copier.AllNilFields)
	mcols := map[string]string{}
	for _, f := range fields {
		var mapKey string
		if len(tagNames) > 0 {
			for _, tagName := range tagNames {
				mapKey = f.Tag.Get(tagName)
				if len(mapKey) > 0 {
					mapKey = strings.SplitN(mapKey, `,`, 2)[0]
					break
				}
			}
			if mapKey == "-" {
				continue
			}
		}
		if len(mapKey) == 0 {
			mapKey = f.Name
		}
		mcols[mapKey] = f.Name
	}
	for k, vv := range m {
		if name, ok := mcols[k]; ok {
			kv := v.FieldByName(name)
			convertAssign(kv, vv)
		}
	}
	return
}

// convertAssign copies to dest the value in src, converting it if possible.
// An error is returned if the copy would result in loss of information.
// dest should be a pointer type.
func convertAssign(dest reflect.Value, src interface{}) error {
	// Common cases, without reflect.
	switch s := src.(type) {
	case string:
		switch dest.Kind() {
		case reflect.String:
			if dest.CanSet() {
				dest.SetString(s)
				return nil
			}
		case reflect.Slice:
			if dest.CanSet() {
				if dest.Elem().Kind() == reflect.Uint8 {
					dest.SetBytes([]byte(s))
					return nil
				}
			}
		}
	case []byte:
		switch dest.Kind() {
		case reflect.String:
			if dest.CanSet() {
				dest.SetString(string(s))
				return nil
			}
		case reflect.Slice:
			if dest.CanSet() {
				if dest.Elem().Kind() == reflect.Uint8 {
					dest.SetBytes([]byte(s))
					return nil
				}
			}
		}
	case time.Time:
		switch dest.Kind() {
		case reflect.String:
			if dest.CanSet() {
				dest.SetString(s.Format(time.RFC3339Nano))
				return nil
			}
		case reflect.Slice:
			if dest.CanSet() {
				if dest.Elem().Kind() == reflect.Uint8 {
					dest.SetBytes([]byte(s.Format(time.RFC3339Nano)))
					return nil
				}
			}
		}
	case nil:
		if dest.CanSet() {
			dest.SetPointer(nil)
		}
	}

	sv := reflect.ValueOf(src)

	switch dest.Kind() {
	case reflect.String:
		if dest.CanSet() {
			switch sv.Kind() {
			case reflect.Bool,
				reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
				reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
				reflect.Float32, reflect.Float64:
				dest.SetString(asString(src))
				return nil
			}

		}
	case reflect.Slice:
		if dest.Elem().Kind() == reflect.Uint8 {
			sv = reflect.ValueOf(src)
			if b, ok := asBytes(nil, sv); ok {
				dest.SetBytes(b)
				return nil
			}
		}
	case reflect.Bool:
		if dest.CanSet() {
			switch sv.Kind() {
			case reflect.Bool:
				dest.SetBool(sv.Bool())
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				dest.SetBool(sv.Int() != 0)
			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				dest.SetBool(int(sv.Uint()) != 0)
			case reflect.Float32, reflect.Float64:
				dest.SetBool(int(sv.Float()) != 0)
				return nil
			}
		}
	case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int:
		if dest.CanSet() {
			switch sv.Kind() {
			case reflect.Bool:
				if sv.Bool() {
					dest.SetInt(1)
				} else {
					dest.SetInt(0)
				}
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				dest.SetInt(sv.Int())
			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				dest.SetInt(int64(sv.Uint()))
			case reflect.Float32, reflect.Float64:
				dest.SetInt(int64(sv.Float()))
				return nil
			}
		}
	case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint:
		if dest.CanSet() {
			switch sv.Kind() {
			case reflect.Bool:
				if sv.Bool() {
					dest.SetUint(1)
				} else {
					dest.SetUint(0)
				}
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				dest.SetUint(uint64(sv.Int()))
			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				dest.SetUint(sv.Uint())
			case reflect.Float32, reflect.Float64:
				dest.SetUint(uint64(sv.Float()))
				return nil
			}
		}
	case reflect.Float32, reflect.Float64:
		if dest.CanSet() {
			switch sv.Kind() {
			case reflect.Bool:
				if sv.Bool() {
					dest.SetFloat(1)
				} else {
					dest.SetFloat(0)
				}
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				dest.SetFloat(float64(sv.Int()))
			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				dest.SetFloat(float64(sv.Uint()))
			case reflect.Float32, reflect.Float64:
				dest.SetFloat(sv.Float())
				return nil
			}
		}
	}

	return fmt.Errorf("unsupported Scan, storing driver.Value type %T into type %T", src, dest)
}

func asString(src interface{}) string {
	switch v := src.(type) {
	case string:
		return v
	case []byte:
		return string(v)
	}
	rv := reflect.ValueOf(src)
	switch rv.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.FormatInt(rv.Int(), 10)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return strconv.FormatUint(rv.Uint(), 10)
	case reflect.Float64:
		return strconv.FormatFloat(rv.Float(), 'g', -1, 64)
	case reflect.Float32:
		return strconv.FormatFloat(rv.Float(), 'g', -1, 32)
	case reflect.Bool:
		return strconv.FormatBool(rv.Bool())
	}
	return fmt.Sprintf("%v", src)
}
func asBytes(buf []byte, rv reflect.Value) (b []byte, ok bool) {
	switch rv.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.AppendInt(buf, rv.Int(), 10), true
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return strconv.AppendUint(buf, rv.Uint(), 10), true
	case reflect.Float32:
		return strconv.AppendFloat(buf, rv.Float(), 'g', -1, 32), true
	case reflect.Float64:
		return strconv.AppendFloat(buf, rv.Float(), 'g', -1, 64), true
	case reflect.Bool:
		return strconv.AppendBool(buf, rv.Bool()), true
	case reflect.String:
		s := rv.String()
		return append(buf, s...), true
	}
	return
}
func strconvErr(err error) error {
	if ne, ok := err.(*strconv.NumError); ok {
		return ne.Err
	}
	return err
}
