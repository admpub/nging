package tagfast

import (
	"reflect"
	"sync"
)

var (
	lock   = new(sync.RWMutex)
	caches = make(map[string]map[string]Faster) //{"pkg.struct":{"field":Faster}}
)

func Tag(t reflect.Type, f reflect.StructField, tagName string) (value string, faster Faster) {
	faster = tag(t, f)
	if faster == nil {
		return
	}
	value = faster.Get(tagName)
	return
}

func setCache(structName string, fieldName string, fast Faster) {
	lock.Lock()
	_, ok := caches[structName]
	if !ok {
		caches[structName] = make(map[string]Faster)
	}
	caches[structName][fieldName] = fast
	lock.Unlock()
}

func getCache(structName string, fieldName string) (fast Faster) {
	lock.RLock()
	cc, ok := caches[structName]
	if ok {
		tf, ok := cc[fieldName]
		if ok {
			fast = tf
		}
	}
	lock.RUnlock()
	return fast
}

func tag(t reflect.Type, f reflect.StructField) Faster {
	if len(f.Tag) == 0 {
		return nil
	}
	name := t.PkgPath() + "." + t.Name()
	fast := getCache(name, f.Name)
	if fast == nil {
		fast = New(f.Tag)
		setCache(name, f.Name, fast)
	}
	return fast
}

func Parsed(t reflect.Type, f reflect.StructField, tagName string, fns ...func() interface{}) interface{} {
	faster := tag(t, f)
	if faster == nil {
		return nil
	}
	return faster.Parsed(tagName, fns...)
}

func GetParsed(t reflect.Type, f reflect.StructField, tagName string, fns ...func(string) interface{}) interface{} {
	faster := tag(t, f)
	if faster == nil {
		return nil
	}
	return faster.GetParsed(tagName, fns...)
}

func Value(t reflect.Type, f reflect.StructField, tagName string) (value string) {
	value, _ = Tag(t, f, tagName)
	return
}

func Caches() map[string]map[string]Faster {
	return caches
}
