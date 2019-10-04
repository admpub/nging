package factory

import (
	"github.com/webx-top/db"
	"github.com/webx-top/echo"
)

type Base struct {
	param   *Param
	trans   *Transaction
	namer   func(string) string
	connID  int
	context echo.Context
}

func (this *Base) SetParam(param *Param) {
	this.param = param
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

func (this *Base) SetContext(ctx echo.Context) {
	this.context = ctx
}

func (this *Base) Context() echo.Context {
	return this.context
}

func (this *Base) SetConnID(connID int) {
	this.connID = connID
}

func (this *Base) ConnID() int {
	return this.connID
}

func (this *Base) SetNamer(namer func(string) string) {
	this.namer = namer
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
	Param() *Param
	NewObjects() Ranger
	Get(mw func(db.Result) db.Result, args ...interface{}) error
	List(recv interface{}, mw func(db.Result) db.Result, page, size int, args ...interface{}) (func() int64, error)
	ListByOffset(recv interface{}, mw func(db.Result) db.Result, offset, size int, args ...interface{}) (func() int64, error)
	Add() (interface{}, error)
	Edit(mw func(db.Result) db.Result, args ...interface{}) error
	Upsert(mw func(db.Result) db.Result, args ...interface{}) (interface{}, error)
	Delete(mw func(db.Result) db.Result, args ...interface{}) error
	Count(mw func(db.Result) db.Result, args ...interface{}) (int64, error)
	SetField(mw func(db.Result) db.Result, field string, value interface{}, args ...interface{}) error
	SetFields(mw func(db.Result) db.Result, kvset map[string]interface{}, args ...interface{}) error
	AsMap() map[string]interface{}
	AsRow() map[string]interface{}
	FromRow(row map[string]interface{})
	Set(key interface{}, value ...interface{})
	BatchValidate(kvset map[string]interface{}) error
	Validate(field string, value interface{}) error
}
