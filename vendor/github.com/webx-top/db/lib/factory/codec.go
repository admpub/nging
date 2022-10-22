package factory

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"reflect"
	"sort"

	"github.com/webx-top/com"
	"github.com/webx-top/echo/param"
)

var _ json.Marshaler = (*Codec)(nil)
var _ json.Marshaler = Codec{}
var _ json.Unmarshaler = (*Codec)(nil)
var _ json.Unmarshaler = Codec{}
var _ xml.Marshaler = (*Codec)(nil)
var _ xml.Marshaler = Codec{}
var _ xml.Unmarshaler = (*Codec)(nil)
var _ xml.Unmarshaler = Codec{}

const (
	KeystyleSnakeCase  = `SnakeCase`  // WebxTop => webx_top
	KeystyleCamelCase  = `CamelCase`  // webx_top => webxTop
	KeystylePascalCase = `PascalCase` // webx_top => WebxTop
)

type Mapper interface {
	AsMap(onlyFields ...string) param.Store
	AsRow(onlyFields ...string) param.Store
	FromRow(row map[string]interface{})
	Set(key interface{}, value ...interface{})
}

func NewCodec(data Mapper, keystyle ...string) *Codec {
	c := &Codec{
		Mapper: data,
	}
	if len(keystyle) > 0 {
		c.keystyle = keystyle[0]
	}
	return c
}

type Codec struct {
	Mapper
	columns  []string // SnakeCase
	fields   []string // PascalCase
	keystyle string
}

func (c *Codec) SetModel(data Mapper) {
	c.Mapper = data
}

func CreateNewInstance(instance interface{}) interface{} {
	t := reflect.Indirect(reflect.ValueOf(instance)).Type()
	if t.Kind() != reflect.Struct {
		panic(fmt.Sprintf(`unsupported instance type: %T`, instance))
	}
	return reflect.New(t).Interface()
}

func (c *Codec) Clone(data ...Mapper) *Codec {
	copied := &Codec{
		Mapper:   c.Mapper,
		columns:  make([]string, len(c.columns)),
		fields:   make([]string, len(c.fields)),
		keystyle: c.keystyle,
	}
	copy(copied.columns, c.columns)
	copy(copied.fields, c.fields)
	if len(data) > 0 {
		copied.Mapper = data[0]
	}
	// else if c.Mapper != nil {
	// 	copied.Mapper = CreateNewInstance(c.Mapper).(Mapper)
	// }
	return copied
}

func (c *Codec) SetColumns(columns ...string) {
	c.columns = columns
}

func (c *Codec) SetFields(fields ...string) {
	c.fields = fields
}

func (c *Codec) SetKeystyle(keystyle string) {
	c.keystyle = keystyle
}

func (c *Codec) Columns() []string {
	if len(c.columns) == len(c.fields) {
		return c.columns
	}
	columns := make([]string, len(c.fields))
	for i, v := range c.fields {
		columns[i] = com.PascalCase(v)
	}
	c.columns = columns
	return columns
}

func (c *Codec) Keystyle() string {
	return c.keystyle
}

func (c *Codec) Fields() []string {
	if len(c.columns) == len(c.fields) {
		return c.fields
	}
	fields := make([]string, len(c.columns))
	for i, v := range c.columns {
		fields[i] = com.PascalCase(v)
	}
	c.fields = fields
	return fields
}

func (c Codec) MakeMap() map[string]interface{} {
	var m map[string]interface{}
	switch c.keystyle {
	case KeystyleCamelCase:
		m = map[string]interface{}{}
		am := c.AsMap(c.Fields()...)
		for k, v := range am {
			k = com.LowerCaseFirst(k)
			m[k] = v
		}
	case KeystylePascalCase:
		m = c.AsMap(c.Fields()...)
	case KeystyleSnakeCase:
		fallthrough
	default:
		m = c.AsRow(c.columns...)
	}
	return m
}

func (c Codec) Import(m map[string]interface{}) {
	switch c.keystyle {
	case KeystyleCamelCase:
		for k, v := range m {
			k = com.PascalCase(k)
			c.Set(k, v)
		}
	case KeystylePascalCase:
		c.Set(m)
	case KeystyleSnakeCase:
		fallthrough
	default:
		c.FromRow(m)
	}
}

func (c Codec) MarshalJSON() ([]byte, error) {
	return json.Marshal(c.MakeMap())
}

func (c Codec) UnmarshalJSON(data []byte) error {
	recv := map[string]interface{}{}
	err := json.Unmarshal(data, &recv)
	c.Import(recv)
	return err
}

func (c Codec) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	if c.keystyle == KeystylePascalCase {
		start.Name.Local = `Item`
	} else {
		start.Name.Local = `item`
	}
	if err := e.EncodeToken(start); err != nil {
		return err
	}
	m := c.MakeMap()
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		v := m[k]
		if err := XMLEncode(e, k, v); err != nil {
			return err
		}
	}
	return e.EncodeToken(xml.EndElement{Name: start.Name})
}

func XMLEncode(e *xml.Encoder, key string, value interface{}, attrs ...xml.Attr) error {
	elem := xml.StartElement{
		Name: xml.Name{Space: ``, Local: key},
		Attr: attrs,
	}
	return e.EncodeElement(value, elem)
}

func (c Codec) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	recv := map[string]interface{}{}
	defer c.Import(recv)
	for {
		token, err := d.Token()
		if err != nil || token == nil {
			return err
		}
		switch t := token.(type) {
		case xml.StartElement:
			e := xml.StartElement(t)
			var q string
			err = d.DecodeElement(&q, &e)
			if err != nil {
				return err
			}
			//println(`start`, e.Name.Local, q)
			recv[e.Name.Local] = q
		case xml.EndElement:
			return err
		}
	}
}
