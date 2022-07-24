package standard

import (
	"net/url"
	"sync"

	"github.com/webx-top/echo/engine"
)

var (
	_ engine.URLValuer = &UrlValue{}
	_ engine.URLValuer = &Value{}
)

type UrlValue struct {
	values *url.Values
	initFn func() *url.Values
}

func (u *UrlValue) Add(key string, value string) {
	u.init()
	u.values.Add(key, value)
}

func (u *UrlValue) Del(key string) {
	u.init()
	u.values.Del(key)
}

func (u *UrlValue) Get(key string) string {
	u.init()
	return u.values.Get(key)
}

func (u *UrlValue) Gets(key string) []string {
	u.init()
	if v, ok := (*u.values)[key]; ok {
		return v
	}
	return []string{}
}

func (u *UrlValue) Set(key string, value string) {
	u.init()
	u.values.Set(key, value)
}

func (u *UrlValue) Encode() string {
	u.init()
	return u.values.Encode()
}

func (u *UrlValue) All() map[string][]string {
	u.init()
	return *u.values
}

func (u *UrlValue) Reset(data url.Values) {
	u.values = &data
}

func (u *UrlValue) init() {
	if u.values != nil {
		return
	}
	u.values = u.initFn()
}

func (u *UrlValue) Merge(data url.Values) {
	u.init()
	for key, values := range data {
		for index, value := range values {
			if index == 0 {
				u.values.Set(key, value)
			} else {
				u.values.Add(key, value)
			}
		}
	}
}

func NewValue(r *Request) *Value {
	v := &Value{
		queryArgs: &UrlValue{initFn: func() *url.Values {
			q := r.url.Query()
			return &q
		}},
		request: r,
	}
	v.postArgs = &UrlValue{initFn: func() *url.Values {
		r.MultipartForm()
		return &r.request.PostForm
	}}
	return v
}

type Value struct {
	request   *Request
	queryArgs *UrlValue
	postArgs  *UrlValue
	form      *url.Values
	lock      sync.RWMutex
}

func (v *Value) Add(key string, value string) {
	v.lock.Lock()
	v.init()
	v.form.Add(key, value)
	v.lock.Unlock()
}

func (v *Value) Del(key string) {
	v.lock.Lock()
	v.init()
	v.form.Del(key)
	v.lock.Unlock()
}

func (v *Value) Get(key string) string {
	v.lock.Lock()
	v.init()
	val := v.form.Get(key)
	v.lock.Unlock()
	return val
}

func (v *Value) Gets(key string) []string {
	v.lock.Lock()
	defer v.lock.Unlock()
	v.init()
	form := *v.form
	if v, ok := form[key]; ok {
		return v
	}
	return []string{}
}

func (v *Value) Set(key string, value string) {
	v.lock.Lock()
	v.init()
	v.form.Set(key, value)
	v.lock.Unlock()
}

func (v *Value) Encode() string {
	v.lock.Lock()
	v.init()
	val := v.form.Encode()
	v.lock.Unlock()
	return val
}

func (v *Value) init() {
	if v.form != nil {
		return
	}
	//v.request.request.ParseForm()
	v.request.MultipartForm()
	v.form = &v.request.request.Form
}

func (v *Value) All() map[string][]string {
	v.lock.Lock()
	v.init()
	m := *v.form
	v.lock.Unlock()
	return m
}

func (v *Value) Reset(data url.Values) {
	v.lock.Lock()
	v.form = &data
	v.lock.Unlock()
}

func (v *Value) Merge(data url.Values) {
	v.lock.Lock()
	v.init()
	for key, values := range data {
		for index, value := range values {
			if index == 0 {
				v.form.Set(key, value)
			} else {
				v.form.Add(key, value)
			}
		}
	}
	v.lock.Unlock()
}
