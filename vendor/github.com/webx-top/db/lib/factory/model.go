package factory

import (
	"github.com/webx-top/db"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/param"
)

type Base struct {
	param    *Param
	trans    *Transaction
	namer    func(Model) string
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

func (b *Base) EventOFF(off ...bool) *Base {
	if len(off) == 0 {
		b.eventOff = true
	} else {
		b.eventOff = off[0]
	}
	return b
}

func (b *Base) Eventable() bool {
	return !b.eventOff
}

func (b *Base) EventON(on ...bool) *Base {
	if len(on) == 0 {
		b.eventOff = false
	} else {
		b.eventOff = !on[0]
	}
	return b
}

func (b *Base) SetParam(param *Param) *Base {
	b.param = param
	return b
}

func (b *Base) Param() *Param {
	return b.param
}

func (b *Base) Trans() *Transaction {
	return b.trans
}

func (b *Base) Use(trans *Transaction) {
	b.trans = trans
}

func (b *Base) SetContext(ctx echo.Context) *Base {
	b.context = ctx
	if ctx == nil {
		return b
	}
	if setter, ok := ctx.(ModelBaseSetter); ok {
		setter.SetModelBase(b)
	}
	switch t := ctx.Transaction().(type) {
	case *echo.BaseTransaction:
		if tr, ok := t.Transaction.(*Param); ok {
			b.trans = tr.T()
		}
	case *Param:
		b.trans = t.T()
	}
	return b
}

func (b *Base) Context() echo.Context {
	return b.context
}

func (b *Base) SetConnID(connID int) *Base {
	b.connID = connID
	return b
}

func (b *Base) ConnID() int {
	return b.connID
}

func (b *Base) SetNamer(namer func(Model) string) *Base {
	b.namer = namer
	return b
}

func (b *Base) Namer() func(Model) string {
	return b.namer
}

func (b *Base) FieldInfo(dbi *DBI, tableName, columnName string) FieldInfor {
	info, _ := dbi.Fields.Find(tableName, columnName)
	return info
}

type Model interface {
	Trans() *Transaction
	Use(trans *Transaction) Model
	SetContext(ctx echo.Context) Model
	Context() echo.Context
	SetNamer(func(Model) string) Model
	Namer() func(Model) string
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
	UpdateFields(mw func(db.Result) db.Result, kvset map[string]interface{}, args ...interface{}) error
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
