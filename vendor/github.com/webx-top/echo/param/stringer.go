package param

import "time"

type (
	Stringer interface {
		String(v interface{}) string
	}
	Ignorer interface {
		Ignore() bool
	}
	StringerFunc   func(interface{}) string
	StringerMap    map[string]Stringer
	StringerList   []Stringer
	StringerIgnore struct{}
)

func StringerMapStart() StringerMap {
	return StringerMap{}
}

var Ignored = &StringerIgnore{}

func (f StringerFunc) String(v interface{}) string {
	return f(v)
}

func (s *StringerIgnore) Ignore() bool {
	return true
}

func (s *StringerIgnore) String(_ interface{}) string {
	return ``
}

func (s StringerList) Ignore() bool {
	for _, f := range s {
		if ig, ok := f.(Ignorer); ok {
			if ig.Ignore() {
				return true
			}
		}
	}
	return false
}

func (s StringerList) String(v interface{}) string {
	for _, f := range s {
		v = f.String(v)
	}
	return AsString(v)
}

func (s StringerList) Size() int {
	return len(s)
}

func (s StringerMap) Set(key string, value Stringer) StringerMap {
	s[key] = value
	return s
}

func (s StringerMap) SetFunc(key string, value func(interface{}) string) StringerMap {
	s[key] = StringerFunc(value)
	return s
}

func (s StringerMap) Add(key string, value Stringer) StringerMap {
	if st, ok := s[key]; ok {
		switch v := st.(type) {
		case StringerList:
			if sl, ok := value.(StringerList); ok {
				v = append(v, sl...)
			} else {
				v = append(v, value)
			}
			s[key] = v
		default:
			if sl, ok := value.(StringerList); ok {
				v := append(StringerList{v}, sl...)
				s[key] = v
			} else {
				s[key] = StringerList{v, value}
			}
		}
		return s
	}
	s[key] = value
	return s
}

func (s StringerMap) AddFunc(key string, value func(interface{}) string) StringerMap {
	return s.Add(key, StringerFunc(value))
}

func (s StringerMap) Has(key string) bool {
	_, y := s[key]
	return y
}

func (s StringerMap) Get(key string, defaults ...Stringer) Stringer {
	value, ok := s[key]
	if (!ok || value == nil) && len(defaults) > 0 {
		return defaults[0]
	}
	return value
}

func (s StringerMap) String(key string, value interface{}) (result string, found bool, ignore bool) {
	formatter := s.Get(key)
	if formatter == nil {
		return
	}
	if ig, ok := formatter.(Ignorer); ok {
		ignore = ig.Ignore()
		if ignore {
			return
		}
	}
	if sl, ok := formatter.(StringerList); ok {
		if sl.Size() == 0 {
			return
		}
	}
	found = true
	result = formatter.String(value)
	return
}

func (s StringerMap) Delete(keys ...string) StringerMap {
	for _, key := range keys {
		if _, y := s[key]; y {
			delete(s, key)
		}
	}
	return s
}

func TimestampStringer(layouts ...string) Stringer {
	layout := DateTimeNormal
	if len(layouts) > 0 {
		layout = layouts[0]
	}
	return StringerFunc(func(v interface{}) string {
		t := AsTimestamp(v)
		if t.IsZero() {
			return ``
		}
		return t.Format(layout)
	})
}

func WhitespaceStringer() Stringer {
	return StringerFunc(func(_ interface{}) string {
		return ``
	})
}

func DateTimeStringer(layouts ...string) Stringer {
	layout := DateTimeNormal
	if len(layouts) > 0 {
		layout = layouts[0]
	}
	return StringerFunc(func(v interface{}) string {
		switch t := v.(type) {
		case time.Time:
			if t.IsZero() {
				return ``
			}
			return t.Format(layout)
		default:
			return AsString(v)
		}
	})
}
