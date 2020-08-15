package factory

import (
	"github.com/webx-top/db"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/param"
)

type Base struct {
	param    *Param
	trans    *Transaction
	namer    func(string) string
	connID   int
	context  echo.Context
	eventOff bool
}

type ModelBaseSetter interface {
	SetModelBase(*Base)
}

type ModelSetter interface {
	SetModel(Model)
}

func (this *Base) EventOFF(off ...bool) *Base {
	if len(off) == 0 {
		this.eventOff = true
	} else {
		this.eventOff = off[0]
	}
	return this
}

func (this *Base) Eventable() bool {
	return !this.eventOff
}

func (this *Base) EventON(on ...bool) *Base {
	if len(on) == 0 {
		this.eventOff = false
	} else {
		this.eventOff = !on[0]
	}
	return this
}

func (this *Base) SetParam(param *Param) *Base {
	this.param = param
	return this
}

func (this *Base) Param() *Param {
	return this.param
}

func (this *Base) Trans() *Transaction {
	return this.trans
}

func (this *Base) Use(trans *Transaction) {
	this.trans = trans
}

func (this *Base) SetContext(ctx echo.Context) *Base {
	this.context = ctx
	if setter, ok := ctx.(ModelBaseSetter); ok {
		setter.SetModelBase(this)
	}
	switch t := ctx.Transaction().(type) {
	case *echo.BaseTransaction:
		if tr, ok := t.Transaction.(*Param); ok {
			this.trans = tr.T()
		}
	case *Param:
		this.trans = t.T()
	}
	return this
}

func (this *Base) Context() echo.Context {
	return this.context
}

func (this *Base) SetConnID(connID int) *Base {
	this.connID = connID
	return this
}

func (this *Base) ConnID() int {
	return this.connID
}

func (this *Base) SetNamer(namer func(string) string) *Base {
	this.namer = namer
	return this
}

func (this *Base) Namer() func(string) string {
	return this.namer
}

type Model interface {
	Trans() *Transaction
	Use(trans *Transaction) Model
	SetContext(ctx echo.Context) Model
	Context() echo.Context
	SetNamer(func(string) string) Model
	Namer() func(string) string
	CPAFrom(source Model) Model //CopyAttrFrom
	Name_() string
	Short_() string
	Struct_() string
	SetConnID(connID int) Model
	New(structName string, connID ...int) Model
	NewParam() *Param
	SetParam(param *Param) Model
	Param(mw func(db.Result) db.Result, args ...interface{}) *Param
	NewObjects() Ranger
	Get(mw func(db.Result) db.Result, args ...interface{}) error
	List(recv interface{}, mw func(db.Result) db.Result, page, size int, args ...interface{}) (func() int64, error)
	ListByOffset(recv interface{}, mw func(db.Result) db.Result, offset, size int, args ...interface{}) (func() int64, error)
	Add() (interface{}, error)
	Edit(mw func(db.Result) db.Result, args ...interface{}) error
	Upsert(mw func(db.Result) db.Result, args ...interface{}) (interface{}, error)
	Delete(mw func(db.Result) db.Result, args ...interface{}) error
	Count(mw func(db.Result) db.Result, args ...interface{}) (int64, error)
	Exists(mw func(db.Result) db.Result, args ...interface{}) (bool, error)
	SetField(mw func(db.Result) db.Result, field string, value interface{}, args ...interface{}) error
	SetFields(mw func(db.Result) db.Result, kvset map[string]interface{}, args ...interface{}) error
	AsMap() param.Store
	AsRow() param.Store
	FromRow(row map[string]interface{})
	Set(key interface{}, value ...interface{})
	BatchValidate(kvset map[string]interface{}) error
	Validate(field string, value interface{}) error
	EventON(on ...bool) Model
	EventOFF(off ...bool) Model
}
