// Package db Copyright (c) 2012-present upper.io/db authors. All rights reserved.
// Package db Copyright (c) 2017-present Hank Shen. All rights reserved.
package db

import (
	"context"
	"fmt"
	"sync"
)

func NewKeysValues() *KeysValues {
	return &KeysValues{}
}

type KeysValues struct {
	keys   []string
	values []interface{}
}

func (k *KeysValues) Keys() []string {
	return k.keys
}

func (k *KeysValues) Values() []interface{} {
	return k.values
}

func (k *KeysValues) Add(key string, value interface{}) *KeysValues {
	k.keys = append(k.keys, key)
	k.values = append(k.values, value)
	return k
}

func (k *KeysValues) Reset() *KeysValues {
	k.keys = k.keys[0:0]
	k.values = k.values[0:0]
	return k
}

func (k *KeysValues) String() string {
	return fmt.Sprintf("keys: %#v\nvalues: %#v", k.keys, k.values)
}

// Slice 依次填充key和value
func (k *KeysValues) Slice() []interface{} {
	var data []interface{}
	vl := len(k.values)
	for i, kk := range k.keys {
		data = append(data, kk)
		if i < vl {
			data = append(data, k.values[i])
		} else {
			data = append(data, nil)
		}
	}
	return data
}

func (k *KeysValues) Map() map[string]interface{} {
	data := map[string]interface{}{}
	vl := len(k.values)
	for i, kk := range k.keys {
		if i < vl {
			data[kk] = k.values[i]
		} else {
			data[kk] = nil
		}
	}
	return data
}

type Compounds []Compound

var compoundsPool = sync.Pool{
	New: func() interface{} {
		return NewCompounds()
	},
}

func CompoundsPoolGet() *Compounds {
	return compoundsPool.Get().(*Compounds)
}

func CompoundsPoolRelease(c *Compounds) {
	c.Reset()
	compoundsPool.Put(c)
}

func NewCompounds() *Compounds {
	return &Compounds{}
}

func (c *Compounds) AddKV(key, value interface{}) *Compounds {
	*c = append(*c, Cond{key: value})
	return c
}

func (c *Compounds) AddKVNotEmpty(key, value interface{}) *Compounds {
	switch v := value.(type) {
	case nil:
		return c
	case string:
		if len(v) == 0 {
			return c
		}
	case []string:
		if len(v) == 0 {
			return c
		}
	case uint8:
		if v == 0 {
			return c
		}
	case []uint8:
		if len(v) == 0 {
			return c
		}
	case int8:
		if v == 0 {
			return c
		}
	case []int8:
		if len(v) == 0 {
			return c
		}
	case uint16:
		if v == 0 {
			return c
		}
	case []uint16:
		if len(v) == 0 {
			return c
		}
	case int16:
		if v == 0 {
			return c
		}
	case []int16:
		if len(v) == 0 {
			return c
		}
	case uint32:
		if v == 0 {
			return c
		}
	case []uint32:
		if len(v) == 0 {
			return c
		}
	case int32:
		if v == 0 {
			return c
		}
	case []int32:
		if len(v) == 0 {
			return c
		}
	case uint:
		if v == 0 {
			return c
		}
	case []uint:
		if len(v) == 0 {
			return c
		}
	case int:
		if v == 0 {
			return c
		}
	case []int:
		if len(v) == 0 {
			return c
		}
	case uint64:
		if v == 0 {
			return c
		}
	case []uint64:
		if len(v) == 0 {
			return c
		}
	case int64:
		if v == 0 {
			return c
		}
	case []int64:
		if len(v) == 0 {
			return c
		}
	case float32:
		if v == 0 {
			return c
		}
	case []float32:
		if len(v) == 0 {
			return c
		}
	case float64:
		if v == 0 {
			return c
		}
	case []float64:
		if len(v) == 0 {
			return c
		}
	}
	*c = append(*c, Cond{key: value})
	return c
}

func (c *Compounds) AddKVGtZero(key, value interface{}) *Compounds {
	switch v := value.(type) {
	case int8:
		if v > 0 {
			return c
		}
	case int16:
		if v > 0 {
			return c
		}
	case int32:
		if v > 0 {
			return c
		}
	case int:
		if v > 0 {
			return c
		}
	case int64:
		if v > 0 {
			return c
		}
	case float32:
		if v > 0 {
			return c
		}
	case float64:
		if v > 0 {
			return c
		}
	default:
		return c.AddKVNotEmpty(key, value)
	}
	*c = append(*c, Cond{key: value})
	return c
}

func (c *Compounds) Set(compounds ...Compound) *Compounds {
	*c = compounds
	return c
}

func (c *Compounds) Add(compounds ...Compound) *Compounds {
	*c = append(*c, compounds...)
	return c
}

func (c *Compounds) From(from *Compounds) *Compounds {
	if from.Size() == 0 {
		return c
	}
	return c.Add(from.V()...)
}

func (c *Compounds) Slice() []Compound {
	return *c
}

func (c *Compounds) V() []Compound {
	return c.Slice()
}

func (c *Compounds) Size() int {
	return len(*c)
}

var _ Compound = NewCompounds()

func (c *Compounds) And(compounds ...Compound) Compound {
	c.Add(compounds...)
	switch c.Size() {
	case 0:
		return EmptyCond
	case 1:
		return (*c)[0]
	default:
		return And(*c...)
	}
}

func (c *Compounds) Or(compounds ...Compound) Compound {
	c.Add(compounds...)
	switch c.Size() {
	case 0:
		return EmptyCond
	case 1:
		return (*c)[0]
	default:
		return Or(*c...)
	}
}

// Sentences return each one of the map records as a compound.
func (c *Compounds) Sentences() []Compound {
	return c.Slice()
}

// Operator returns the default compound operator.
func (c *Compounds) Operator() CompoundOperator {
	return OperatorAnd
}

// Empty returns false if there are no conditions.
func (c *Compounds) Empty() bool {
	return c.Size() == 0
}

func (c *Compounds) Reset() {
	if c.Empty() {
		return
	}
	*c = (*c)[0:0]
}

func (c *Compounds) remove(s int) Compounds {
	return append((*c)[:s], (*c)[s+1:]...)
}

func (c *Compounds) Delete(keys ...interface{}) {
	for _, key := range keys {
		for i, v := range *c {
			r, y := v.(Cond)
			if !y {
				continue
			}
			_, ok := r[key]
			if !ok {
				continue
			}
			delete(r, key)
			if len(r) == 0 {
				*c = c.remove(i)
			}
		}
	}
}

type TableName interface {
	TableName() string
}

type tableNameString string

func (t tableNameString) TableName() string {
	return string(t)
}

func Table(tableName string) TableName {
	return tableNameString(tableName)
}

type StdContext interface {
	StdContext() context.Context
}

type RequestURI interface {
	RequestURI() string
}

type Method interface {
	Method() string
}
