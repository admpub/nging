package namedstruct

import (
	"testing"

	"github.com/admpub/pp"
	"github.com/stretchr/testify/assert"
)

type TestStruct struct {
	Name string
}

func TestMakeSlice(t *testing.T) {
	s := NewStructs()
	s.Register(`TestStruct`, &TestStruct{})
	sliceV := s.MakeSlice(`TestStruct`)
	assert.NotNil(t, sliceV)
	pp.Println(sliceV)
	v := s.Make(`TestStruct`)
	assert.NotNil(t, sliceV)
	pp.Println(v)

	sliceV = s.MakeSlice(`TestStructNotExists`)
	assert.Nil(t, sliceV)
	pp.Println(sliceV)
	v = s.Make(`TestStructNotExists`)
	assert.Nil(t, sliceV)
	pp.Println(v)
}

func TestConvertToSlice(t *testing.T) {
	sliceV := ConvertToSlice([]*TestStruct{
		{Name: `1`},
	})
	assert.NotNil(t, sliceV)
	pp.Println(sliceV)
	sliceV = ConvertToSlice(&[]*TestStruct{
		{Name: `2`},
	})
	assert.NotNil(t, sliceV)
	pp.Println(sliceV)

	sliceV = ConvertToSlice(123)
	assert.Nil(t, sliceV)
	pp.Println(sliceV)
}
