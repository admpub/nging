package param

import (
	"encoding/xml"
	"fmt"
	"html/template"
	"time"
)

var (
	emptyStore = Store{}
)

func StoreStart() Store {
	return Store{}
}

type Store map[string]interface{}

func (s Store) Set(key string, value interface{}) Store {
	s[key] = value
	return s
}

func (s Store) Has(key string) bool {
	_, y := s[key]
	return y
}

func (s Store) Get(key string, defaults ...interface{}) interface{} {
	value, ok := s[key]
	if (!ok || value == nil) && len(defaults) > 0 {
		if fallback, ok := defaults[0].(func() interface{}); ok {
			return fallback()
		}
		return defaults[0]
	}
	return value
}

func (s Store) String(key string, defaults ...interface{}) string {
	return AsString(s.Get(key, defaults...))
}

func (s Store) Split(key string, sep string, limit ...int) StringSlice {
	return Split(s.Get(key), sep, limit...)
}

func (s Store) Trim(key string, defaults ...interface{}) String {
	return Trim(s.Get(key, defaults...))
}

func (s Store) HTML(key string, defaults ...interface{}) template.HTML {
	return AsHTML(s.Get(key, defaults...))
}

func (s Store) HTMLAttr(key string, defaults ...interface{}) template.HTMLAttr {
	return AsHTMLAttr(s.Get(key, defaults...))
}

func (s Store) JS(key string, defaults ...interface{}) template.JS {
	return AsJS(s.Get(key, defaults...))
}

func (s Store) CSS(key string, defaults ...interface{}) template.CSS {
	return AsCSS(s.Get(key, defaults...))
}

func (s Store) Bool(key string, defaults ...interface{}) bool {
	return AsBool(s.Get(key, defaults...))
}

func (s Store) Float64(key string, defaults ...interface{}) float64 {
	return AsFloat64(s.Get(key, defaults...))
}

func (s Store) Float32(key string, defaults ...interface{}) float32 {
	return AsFloat32(s.Get(key, defaults...))
}

func (s Store) Int8(key string, defaults ...interface{}) int8 {
	return AsInt8(s.Get(key, defaults...))
}

func (s Store) Int16(key string, defaults ...interface{}) int16 {
	return AsInt16(s.Get(key, defaults...))
}

func (s Store) Int(key string, defaults ...interface{}) int {
	return AsInt(s.Get(key, defaults...))
}

func (s Store) Int32(key string, defaults ...interface{}) int32 {
	return AsInt32(s.Get(key, defaults...))
}

func (s Store) Int64(key string, defaults ...interface{}) int64 {
	return AsInt64(s.Get(key, defaults...))
}

func (s Store) Decr(key string, n int64, defaults ...interface{}) int64 {
	v := Decr(s.Get(key, defaults...), n)
	s.Set(key, v)
	return v
}

func (s Store) Incr(key string, n int64, defaults ...interface{}) int64 {
	v := Incr(s.Get(key, defaults...), n)
	s.Set(key, v)
	return v
}

func (s Store) Uint8(key string, defaults ...interface{}) uint8 {
	return AsUint8(s.Get(key, defaults...))
}

func (s Store) Uint16(key string, defaults ...interface{}) uint16 {
	return AsUint16(s.Get(key, defaults...))
}

func (s Store) Uint(key string, defaults ...interface{}) uint {
	return AsUint(s.Get(key, defaults...))
}

func (s Store) Uint32(key string, defaults ...interface{}) uint32 {
	return AsUint32(s.Get(key, defaults...))
}

func (s Store) Uint64(key string, defaults ...interface{}) uint64 {
	return AsUint64(s.Get(key, defaults...))
}

func (s Store) Timestamp(key string, defaults ...interface{}) time.Time {
	return AsTimestamp(s.Get(key, defaults...))
}

func (s Store) DateTime(key string, layouts ...string) time.Time {
	return AsDateTime(s.Get(key), layouts...)
}

func (s Store) Children(keys ...interface{}) Store {
	r := s
	for _, key := range keys {
		r = r.GetStore(fmt.Sprint(key))
	}
	return r
}

func (s Store) GetStore(key string, defaults ...interface{}) Store {
	return AsStore(s.Get(key, defaults...))
}

func (s Store) GetStoreByKeys(keys ...string) Store {
	sz := len(keys)
	if sz == 0 {
		return s
	}
	r := s.GetStore(keys[0])
	if sz == 1 {
		return r
	}
	for _, key := range keys[1:] {
		r = r.GetStore(key)
	}
	return r
}

func (s Store) Delete(keys ...string) Store {
	for _, key := range keys {
		if _, y := s[key]; y {
			delete(s, key)
		}
	}
	return s
}

// MarshalXML allows type Store to be used with xml.Marshal
func (s Store) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	if start.Name.Local == `Store` {
		start.Name.Local = `Map`
	}
	if err := e.EncodeToken(start); err != nil {
		return err
	}
	for key, value := range s {
		elem := xml.StartElement{
			Name: xml.Name{Space: ``, Local: key},
			Attr: []xml.Attr{},
		}
		if err := e.EncodeElement(value, elem); err != nil {
			return err
		}
	}
	return e.EncodeToken(xml.EndElement{Name: start.Name})
}

func (s Store) DeepMerge(source Store) {
	for k, value := range source {
		var (
			destValue interface{}
			ok        bool
		)
		if destValue, ok = s[k]; !ok {
			s[k] = value
			continue
		}
		sourceM, sourceOk := value.(Store)
		destM, destOk := destValue.(Store)
		if sourceOk && sourceOk == destOk {
			destM.DeepMerge(sourceM)
		} else {
			s[k] = value
		}
	}
}

func (s Store) Clone() Store {
	r := make(Store)
	for k, value := range s {
		switch v := value.(type) {
		case Store:
			r[k] = v.Clone()
		case []Store:
			vCopy := make([]Store, len(v))
			for i, row := range v {
				vCopy[i] = row.Clone()
			}
			r[k] = vCopy
		default:
			r[k] = value
		}
	}
	return r
}

func (s Store) Transform(transfers map[string]Transfer) Store {
	rmap := Store{}
	for key, transfer := range transfers {
		value, _ := s[key]
		if transfer == nil {
			rmap[key] = value
			continue
		}
		newKey := transfer.Destination()
		if len(newKey) == 0 {
			newKey = key
		}
		rmap[newKey] = transfer.Transform(value, s)
	}
	return rmap
}
