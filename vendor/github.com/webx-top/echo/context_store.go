package echo

import (
	"encoding/xml"
	"fmt"
	"html/template"
	"strconv"
	"sync"
	"time"

	"github.com/webx-top/echo/param"
)

var (
	mutex      sync.RWMutex
	emptyStore = Store{}
)

type Store map[string]interface{}

func (s Store) Set(key string, value interface{}) Store {
	mutex.Lock()
	s[key] = value
	mutex.Unlock()
	return s
}

func (s Store) Has(key string) bool {
	mutex.RLock()
	defer mutex.RUnlock()
	_, y := s[key]
	return y
}

func (s Store) Get(key string, defaults ...interface{}) interface{} {
	mutex.RLock()
	defer mutex.RUnlock()
	if v, y := s[key]; y {
		if v == nil && len(defaults) > 0 {
			return defaults[0]
		}
		return v
	}
	if len(defaults) > 0 {
		return defaults[0]
	}
	return nil
}

func (s Store) String(key string, defaults ...interface{}) string {
	return param.AsString(s.Get(key, defaults...))
}

func (s Store) Split(key string, sep string, limit ...int) param.StringSlice {
	return param.Split(s.Get(key), sep, limit...)
}

func (s Store) Trim(key string, defaults ...interface{}) param.String {
	return param.Trim(s.Get(key, defaults...))
}

func (s Store) HTML(key string, defaults ...interface{}) template.HTML {
	return param.AsHTML(s.Get(key, defaults...))
}

func (s Store) HTMLAttr(key string, defaults ...interface{}) template.HTMLAttr {
	return param.AsHTMLAttr(s.Get(key, defaults...))
}

func (s Store) JS(key string, defaults ...interface{}) template.JS {
	return param.AsJS(s.Get(key, defaults...))
}

func (s Store) CSS(key string, defaults ...interface{}) template.CSS {
	return param.AsCSS(s.Get(key, defaults...))
}

func (s Store) Bool(key string, defaults ...interface{}) bool {
	return param.AsBool(s.Get(key, defaults...))
}

func (s Store) Float64(key string, defaults ...interface{}) float64 {
	return param.AsFloat64(s.Get(key, defaults...))
}

func (s Store) Float32(key string, defaults ...interface{}) float32 {
	return param.AsFloat32(s.Get(key, defaults...))
}

func (s Store) Int8(key string, defaults ...interface{}) int8 {
	return param.AsInt8(s.Get(key, defaults...))
}

func (s Store) Int16(key string, defaults ...interface{}) int16 {
	return param.AsInt16(s.Get(key, defaults...))
}

func (s Store) Int(key string, defaults ...interface{}) int {
	return param.AsInt(s.Get(key, defaults...))
}

func (s Store) Int32(key string, defaults ...interface{}) int32 {
	return param.AsInt32(s.Get(key, defaults...))
}

func (s Store) Int64(key string, defaults ...interface{}) int64 {
	return param.AsInt64(s.Get(key, defaults...))
}

func (s Store) Decr(key string, n int64, defaults ...interface{}) int64 {
	v := param.Decr(s.Get(key, defaults...), n)
	s.Set(key, v)
	return v
}

func (s Store) Incr(key string, n int64, defaults ...interface{}) int64 {
	v := param.Incr(s.Get(key, defaults...), n)
	s.Set(key, v)
	return v
}

func (s Store) Uint8(key string, defaults ...interface{}) uint8 {
	return param.AsUint8(s.Get(key, defaults...))
}

func (s Store) Uint16(key string, defaults ...interface{}) uint16 {
	return param.AsUint16(s.Get(key, defaults...))
}

func (s Store) Uint(key string, defaults ...interface{}) uint {
	return param.AsUint(s.Get(key, defaults...))
}

func (s Store) Uint32(key string, defaults ...interface{}) uint32 {
	return param.AsUint32(s.Get(key, defaults...))
}

func (s Store) Uint64(key string, defaults ...interface{}) uint64 {
	return param.AsUint64(s.Get(key, defaults...))
}

func (s Store) Timestamp(key string, defaults ...interface{}) time.Time {
	return param.AsTimestamp(s.Get(key, defaults...))
}

func (s Store) DateTime(key string, layouts ...string) time.Time {
	return param.AsDateTime(s.Get(key), layouts...)
}

func (s Store) Children(keys ...interface{}) Store {
	r := s
	for _, key := range keys {
		r = r.Store(fmt.Sprint(key))
	}
	return r
}

func (s Store) Store(key string, defaults ...interface{}) Store {
	return AsStore(s.Get(key, defaults...))
}

func (s Store) Delete(keys ...string) {
	mutex.Lock()
	for _, key := range keys {
		if _, y := s[key]; y {
			delete(s, key)
		}
	}
	mutex.Unlock()
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

// ToData conversion to *RawData
func (s Store) ToData() *RawData {
	var info, zone, data interface{}
	if v, y := s["Data"]; y {
		data = v
	}
	if v, y := s["Zone"]; y {
		zone = v
	}
	if v, y := s["Info"]; y {
		info = v
	}
	var code State
	if v, y := s["Code"]; y {
		switch c := v.(type) {
		case State:
			code = c
		case int:
			code = State(c)
		case string:
			i, _ := strconv.Atoi(c)
			code = State(i)
		default:
			s := fmt.Sprint(c)
			i, _ := strconv.Atoi(s)
			code = State(i)
		}
	}
	return &RawData{
		Code: code,
		Info: info,
		Zone: zone,
		Data: data,
	}
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
		sourceM, sourceOk := value.(H)
		destM, destOk := destValue.(H)
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

func AsStore(val interface{}) Store {
	switch v := val.(type) {
	case Store:
		return v
	case map[string]interface{}:
		return Store(v)
	case map[string]uint64:
		r := Store{}
		for k, a := range v {
			r[k] = interface{}(a)
		}
		return r
	case map[string]int64:
		r := Store{}
		for k, a := range v {
			r[k] = interface{}(a)
		}
		return r
	case map[string]uint:
		r := Store{}
		for k, a := range v {
			r[k] = interface{}(a)
		}
		return r
	case map[string]int:
		r := Store{}
		for k, a := range v {
			r[k] = interface{}(a)
		}
		return r
	case map[string]uint32:
		r := Store{}
		for k, a := range v {
			r[k] = interface{}(a)
		}
		return r
	case map[string]int32:
		r := Store{}
		for k, a := range v {
			r[k] = interface{}(a)
		}
		return r
	case map[string]float32:
		r := Store{}
		for k, a := range v {
			r[k] = interface{}(a)
		}
		return r
	case map[string]float64:
		r := Store{}
		for k, a := range v {
			r[k] = interface{}(a)
		}
		return r
	case map[string]string:
		r := Store{}
		for k, a := range v {
			r[k] = interface{}(a)
		}
		return r
	default:
		return emptyStore
	}
}
