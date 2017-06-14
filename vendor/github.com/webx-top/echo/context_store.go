package echo

import (
	"html/template"
	"strconv"
	"sync"
)

var (
	mutex         sync.RWMutex
	emptyHTML     = template.HTML(``)
	emptyJS       = template.JS(``)
	emptyCSS      = template.CSS(``)
	emptyHTMLAttr = template.HTMLAttr(``)
)

type store map[string]interface{}

func (s store) Set(key string, value interface{}) store {
	mutex.Lock()
	s[key] = value
	mutex.Unlock()
	return s
}

func (s store) Get(key string) interface{} {
	mutex.RLock()
	defer mutex.RUnlock()
	if v, y := s[key]; y {
		return v
	}
	return nil
}

func (s store) String(key string) string {
	if v, y := s.Get(key).(string); y {
		return v
	}
	return ``
}

func (s store) HTML(key string) template.HTML {
	val := s.Get(key)
	if v, y := val.(template.HTML); y {
		return v
	}
	if v, y := val.(string); y {
		return template.HTML(v)
	}
	return emptyHTML
}

func (s store) HTMLAttr(key string) template.HTMLAttr {
	val := s.Get(key)
	if v, y := val.(template.HTMLAttr); y {
		return v
	}
	if v, y := val.(string); y {
		return template.HTMLAttr(v)
	}
	return emptyHTMLAttr
}

func (s store) JS(key string) template.JS {
	val := s.Get(key)
	if v, y := val.(template.JS); y {
		return v
	}
	if v, y := val.(string); y {
		return template.JS(v)
	}
	return emptyJS
}

func (s store) CSS(key string) template.CSS {
	val := s.Get(key)
	if v, y := val.(template.CSS); y {
		return v
	}
	if v, y := val.(string); y {
		return template.CSS(v)
	}
	return emptyCSS
}

func (s store) Bool(key string) bool {
	if v, y := s.Get(key).(bool); y {
		return v
	}
	return false
}

func (s store) Float64(key string) float64 {
	val := s.Get(key)
	if v, y := val.(float64); y {
		return v
	}
	if v, y := val.(int64); y {
		return float64(v)
	}
	if v, y := val.(uint64); y {
		return float64(v)
	}
	if v, y := val.(float32); y {
		return float64(v)
	}
	if v, y := val.(int32); y {
		return float64(v)
	}
	if v, y := val.(uint32); y {
		return float64(v)
	}
	if v, y := val.(int); y {
		return float64(v)
	}
	if v, y := val.(uint); y {
		return float64(v)
	}
	if v, y := val.(string); y {
		v, _ := strconv.ParseFloat(v, 64)
		return v
	}
	return 0
}

func (s store) Float32(key string) float32 {
	val := s.Get(key)
	if v, y := val.(float32); y {
		return v
	}
	if v, y := val.(int32); y {
		return float32(v)
	}
	if v, y := val.(uint32); y {
		return float32(v)
	}
	if v, y := val.(string); y {
		v, _ := strconv.ParseFloat(v, 32)
		return float32(v)
	}
	return 0
}

func (s store) Int8(key string) int8 {
	val := s.Get(key)
	if v, y := val.(int8); y {
		return v
	}
	if v, y := val.(string); y {
		v, _ := strconv.ParseInt(v, 10, 8)
		return int8(v)
	}
	return 0
}

func (s store) Int16(key string) int16 {
	val := s.Get(key)
	if v, y := val.(int16); y {
		return v
	}
	if v, y := val.(string); y {
		v, _ := strconv.ParseInt(v, 10, 16)
		return int16(v)
	}
	return 0
}

func (s store) Int(key string) int {
	val := s.Get(key)
	if v, y := val.(int); y {
		return v
	}
	if v, y := val.(string); y {
		v, _ := strconv.Atoi(v)
		return v
	}
	return 0
}

func (s store) Int32(key string) int32 {
	val := s.Get(key)
	if v, y := val.(int32); y {
		return v
	}
	if v, y := val.(string); y {
		v, _ := strconv.ParseInt(v, 10, 32)
		return int32(v)
	}
	return 0
}

func (s store) Int64(key string) int64 {
	val := s.Get(key)
	if v, y := val.(int64); y {
		return v
	}
	if v, y := val.(string); y {
		v, _ := strconv.ParseInt(v, 10, 64)
		return v
	}
	return 0
}

func (s store) Decr(key string, n int64) int64 {
	v, _ := s.Get(key).(int64)
	v -= n
	s.Set(key, v)
	return v
}

func (s store) Incr(key string, n int64) int64 {
	v, _ := s.Get(key).(int64)
	v += n
	s.Set(key, v)
	return v
}

func (s store) Uint8(key string) uint8 {
	val := s.Get(key)
	if v, y := val.(uint8); y {
		return v
	}
	if v, y := val.(string); y {
		v, _ := strconv.ParseUint(v, 10, 8)
		return uint8(v)
	}
	return 0
}

func (s store) Uint16(key string) uint16 {
	val := s.Get(key)
	if v, y := val.(uint16); y {
		return v
	}
	if v, y := val.(string); y {
		v, _ := strconv.ParseUint(v, 10, 16)
		return uint16(v)
	}
	return 0
}

func (s store) Uint(key string) uint {
	val := s.Get(key)
	if v, y := val.(uint); y {
		return v
	}
	if v, y := val.(string); y {
		v, _ := strconv.ParseUint(v, 10, 32)
		return uint(v)
	}
	return 0
}

func (s store) Uint32(key string) uint32 {
	val := s.Get(key)
	if v, y := val.(uint32); y {
		return v
	}
	if v, y := val.(string); y {
		v, _ := strconv.ParseUint(v, 10, 32)
		return uint32(v)
	}
	return 0
}

func (s store) Uint64(key string) uint64 {
	val := s.Get(key)
	if v, y := val.(uint64); y {
		return v
	}
	if v, y := val.(string); y {
		v, _ := strconv.ParseUint(v, 10, 64)
		return v
	}
	return 0
}

func (s store) Delete(keys ...string) {
	mutex.Lock()
	for _, key := range keys {
		if _, y := s[key]; y {
			delete(s, key)
		}
	}
	mutex.Unlock()
}
