package null

import (
	"time"

	"github.com/webx-top/echo/param"
)

func (p String) Stringx() param.String {
	return param.String(p.String)
}

func (p String) Interface() interface{} {
	return interface{}(p.String)
}

func (p String) Int() int {
	return p.Stringx().Int()
}

func (p String) Int64() int64 {
	return p.Stringx().Int64()
}

func (p String) Int32() int32 {
	return p.Stringx().Int32()
}

func (p String) Uint() uint {
	return p.Stringx().Uint()
}

func (p String) Uint64() uint64 {
	return p.Stringx().Uint64()
}

func (p String) Uint32() uint32 {
	return p.Stringx().Uint32()
}

func (p String) Float32() float32 {
	return p.Stringx().Float32()
}

func (p String) Float64() float64 {
	return p.Stringx().Float64()
}

func (p String) Bool() bool {
	return p.Stringx().Bool()
}

func (p String) Timestamp() time.Time {
	return p.Stringx().Timestamp()
}

func (p String) DateTime() time.Time {
	return p.Stringx().DateTime()
}
