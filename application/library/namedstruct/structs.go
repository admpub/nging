package namedstruct

import (
	"fmt"
	"reflect"
)

func NewStructs() *Structs {
	return &Structs{
		types: map[string]reflect.Type{},
	}
}

type Structs struct {
	types map[string]reflect.Type
}

func (s *Structs) Register(name string, v interface{}) {
	_, ok := s.GetType(name)
	if ok {
		panic(fmt.Sprintf(`%v: name=%s, type=%T`, ErrNameConflict, name, v))
	}
	s.types[name] = reflect.ValueOf(v).Type()
}

func (s *Structs) MakeSlice(name string) interface{} {
	typ, ok := s.GetType(name)
	if !ok {
		return nil
	}
	sliceType := reflect.SliceOf(typ)
	return reflect.New(sliceType).Interface()
}

func (s *Structs) GetType(name string) (reflect.Type, bool) {
	typ, ok := s.types[name]
	return typ, ok
}

func (s *Structs) Make(name string) interface{} {
	typ, ok := s.GetType(name)
	if !ok {
		return nil
	}
	if typ.Kind() == reflect.Ptr {
		return reflect.New(typ.Elem()).Interface()
	}
	return reflect.New(typ).Interface()
}

func ConvertToSlice(recv interface{}) []interface{} {
	v := reflect.ValueOf(recv)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if v.Kind() != reflect.Slice {
		return nil
	}
	results := make([]interface{}, v.Len())
	for i, j := 0, v.Len(); i < j; i++ {
		results[i] = v.Index(i).Interface()
	}
	return results
}
