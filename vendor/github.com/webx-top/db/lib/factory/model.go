package factory

import "github.com/webx-top/db"

type Model interface {
	Trans() *Transaction
	Use(trans *Transaction) Model
	SetNamer(func(string) string) Model
	Name_() string
	Short_() string
	Struct_() string
	SetConnID(connID int) Model
	New(structName string, connID ...int) Model
	NewParam() *Param
	SetParam(param *Param) Model
	Param() *Param
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
	BatchValidate(kvset map[string]interface{}) error
	Validate(field string, value interface{}) error
}
