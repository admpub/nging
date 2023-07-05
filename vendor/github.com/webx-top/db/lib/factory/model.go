package factory

import (
	"github.com/webx-top/db"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/param"
)

func NewBase(connID int) *Base {
	return &Base{connID: connID}
}

var _ Baser = &Base{}

type Base struct {
	param *Param

	context       echo.Context
	transactioner Transactioner

	namer  func(Model) string
	connID int

	eventOff bool
}

type ModelBaseSetter interface {
	SetModelBase(Baser)
}

type ModelSetter interface {
	SetModel(Model)
}

func (b *Base) EventOFF(off ...bool) Baser {
	if len(off) == 0 {
		b.eventOff = true
	} else {
		b.eventOff = off[0]
	}
	return b
}

func (b Base) Eventable() bool {
	return !b.eventOff
}

func (b *Base) EventON(on ...bool) Baser {
	if len(on) == 0 {
		b.eventOff = false
	} else {
		b.eventOff = !on[0]
	}
	return b
}

func (b *Base) SetParam(param *Param) Baser {
	b.param = param
	return b
}

func (b Base) Param() *Param {
	return b.param
}

func (b Base) T() *Transaction {
	tr := b.Trans()
	if tr == nil {
		return nil
	}
	return tr.T()
}

func (b Base) Trans() Transactioner {
	if b.transactioner == nil && b.context != nil {
		switch t := b.context.Transaction().(type) {
		case echo.UnwrapTransaction:
			if tr, ok := t.Unwrap().(Transactioner); ok {
				return tr
			}
		case Transactioner:
			return t
		}
	}
	return b.transactioner
}

func (b *Base) Use(trans Transactioner) {
	b.transactioner = trans
}

func (b *Base) SetContext(ctx echo.Context) Baser {
	b.context = ctx
	if ctx == nil {
		return b
	}
	if setter, ok := ctx.(ModelBaseSetter); ok {
		setter.SetModelBase(b)
	}
	if b.transactioner != nil {
		switch t := ctx.Transaction().(type) {
		case echo.UnwrapTransaction:
			if tr, ok := t.Unwrap().(Transactioner); ok {
				b.Use(tr)
			}
		case Transactioner:
			b.Use(t)
		}
	}
	return b
}

func (b Base) Context() echo.Context {
	return b.context
}

func (b *Base) SetConnID(connID int) Baser {
	b.connID = connID
	return b
}

func (b Base) ConnID() int {
	return b.connID
}

func (b *Base) SetNamer(namer func(Model) string) Baser {
	b.namer = namer
	return b
}

func (b Base) Namer() func(Model) string {
	return b.namer
}

func (b Base) FieldInfo(dbi *DBI, tableName, columnName string) FieldInfor {
	info, _ := dbi.Fields.Find(tableName, columnName)
	return info
}

func (b Base) New(structName string, connID ...int) Model {
	var m Model
	if len(connID) > 0 {
		m = NewModel(structName, connID[0])
	} else {
		m = NewModel(structName, b.connID)
	}
	return m.SetContext(b.context)
}

type Baser interface {
	EventOFF(off ...bool) Baser
	Eventable() bool
	EventON(on ...bool) Baser
	SetParam(param *Param) Baser
	Param() *Param
	T() *Transaction
	Trans() Transactioner
	Use(trans Transactioner)
	SetContext(ctx echo.Context) Baser
	Context() echo.Context
	SetConnID(connID int) Baser
	ConnID() int
	SetNamer(namer func(Model) string) Baser
	Namer() func(Model) string
	FieldInfo(dbi *DBI, tableName, columnName string) FieldInfor
	New(structName string, connID ...int) Model
}

type Transactioner interface {
	T() *Transaction
}

type Model interface {
	Trans() Transactioner
	Use(trans Transactioner) Model
	SetContext(ctx echo.Context) Model
	Context() echo.Context
	SetNamer(func(Model) string) Model
	Namer() func(Model) string
	CPAFrom(source Model) Model //CopyAttrFrom
	Base_() Baser
	Name_() string
	Short_() string
	Struct_() string
	SetConnID(connID int) Model
	ConnID() int
	New(structName string, connID ...int) Model
	NewParam() *Param
	SetParam(param *Param) Model
	Param(mw func(db.Result) db.Result, args ...interface{}) *Param
	NewObjects() Ranger
	Get(mw func(db.Result) db.Result, args ...interface{}) error
	List(recv interface{}, mw func(db.Result) db.Result, page, size int, args ...interface{}) (func() int64, error)
	ListByOffset(recv interface{}, mw func(db.Result) db.Result, offset, size int, args ...interface{}) (func() int64, error)
	Insert() (interface{}, error)
	Update(mw func(db.Result) db.Result, args ...interface{}) error
	Updatex(mw func(db.Result) db.Result, args ...interface{}) (affected int64, err error)
	UpdateByFields(mw func(db.Result) db.Result, fields []string, args ...interface{}) (err error)
	UpdatexByFields(mw func(db.Result) db.Result, fields []string, args ...interface{}) (affected int64, err error)
	Upsert(mw func(db.Result) db.Result, args ...interface{}) (interface{}, error)
	Delete(mw func(db.Result) db.Result, args ...interface{}) error
	Deletex(mw func(db.Result) db.Result, args ...interface{}) (affected int64, err error)
	Count(mw func(db.Result) db.Result, args ...interface{}) (int64, error)
	Exists(mw func(db.Result) db.Result, args ...interface{}) (bool, error)
	UpdateField(mw func(db.Result) db.Result, field string, value interface{}, args ...interface{}) error
	UpdatexField(mw func(db.Result) db.Result, field string, value interface{}, args ...interface{}) (int64, error)
	UpdateFields(mw func(db.Result) db.Result, kvset map[string]interface{}, args ...interface{}) error
	UpdatexFields(mw func(db.Result) db.Result, kvset map[string]interface{}, args ...interface{}) (int64, error)
	UpdateValues(mw func(db.Result) db.Result, keysValues *db.KeysValues, args ...interface{}) error
	AsMap(onlyFields ...string) param.Store
	AsRow(onlyFields ...string) param.Store
	FromRow(row map[string]interface{})
	Set(key interface{}, value ...interface{})
	BatchValidate(kvset map[string]interface{}) error
	Validate(field string, value interface{}) error
	EventON(on ...bool) Model
	EventOFF(off ...bool) Model
	ListPage(cond *db.Compounds, sorts ...interface{}) error
	ListPageAs(recv interface{}, cond *db.Compounds, sorts ...interface{}) error
}
