//Package db Copyright (c) 2012-present upper.io/db authors. All rights reserved.
//Package db Copyright (c) 2017-present Hank Shen. All rights reserved.
package db

import (
	"fmt"
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

func NewCompounds() *Compounds {
	return &Compounds{}
}

func (c *Compounds) AddKV(key, value interface{}) *Compounds {
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

var EmptyCond = Cond{}

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
