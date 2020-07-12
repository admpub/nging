package null

import (
	"time"

	"github.com/webx-top/echo"
	"github.com/webx-top/echo/param"
)

var EmptyString = String{}

type StringMapSlice []StringMap

func StringMapSliceFrom(list []echo.H) StringMapSlice {
	result := make([]StringMap, len(list))
	for k, values := range list {
		mp := StringMap{}
		result[k] = mp.MapFrom(values)
	}
	return result
}

type StringMap map[string]String

func (p StringMap) Interfaces() map[string]interface{} {
	result := map[string]interface{}{}
	for k, v := range p {
		if v.Valid {
			result[k] = interface{}(v.String)
		} else {
			result[k] = nil
		}
	}
	return result
}

func (p StringMap) StringMap() param.StringMap {
	result := param.StringMap{}
	for k, v := range p {
		result[k] = param.String(v.String)
	}
	return result
}

func (p StringMap) String(key string) string {
	return p.Get(key).String
}

func (p StringMap) IsZero(key string) bool {
	return p.Get(key).IsZero()
}

func (p StringMap) Get(key string) String {
	if value, ok := p[key]; ok {
		return value
	}
	return EmptyString
}

func (p StringMap) MapFrom(values echo.H) StringMap {
	for key, value := range values {
		p[key] = NewString(values.String(key), value != nil)
	}
	return p
}

func (p StringMap) Split(key string, sep string, limit ...int) []string {
	return p.Stringx(key).Split(sep, limit...)
}

func (p StringMap) Stringx(key string) param.String {
	return param.String(p.Get(key).String)
}

func (p StringMap) Interface(key string) interface{} {
	return interface{}(p.Get(key).String)
}

func (p StringMap) Int(key string) int {
	return p.Stringx(key).Int()
}

func (p StringMap) Int64(key string) int64 {
	return p.Stringx(key).Int64()
}

func (p StringMap) Int32(key string) int32 {
	return p.Stringx(key).Int32()
}

func (p StringMap) Uint(key string) uint {
	return p.Stringx(key).Uint()
}

func (p StringMap) Uint64(key string) uint64 {
	return p.Stringx(key).Uint64()
}

func (p StringMap) Uint32(key string) uint32 {
	return p.Stringx(key).Uint32()
}

func (p StringMap) Float32(key string) float32 {
	return p.Stringx(key).Float32()
}

func (p StringMap) Float64(key string) float64 {
	return p.Stringx(key).Float64()
}

func (p StringMap) Bool(key string) bool {
	return p.Stringx(key).Bool()
}

func (p StringMap) Timestamp(key string) time.Time {
	return p.Stringx(key).Timestamp()
}

func (p StringMap) DateTime(key string, layouts ...string) time.Time {
	return p.Stringx(key).DateTime(layouts...)
}
