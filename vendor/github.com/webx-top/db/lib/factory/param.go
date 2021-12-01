//Package factory added by swh@admpub.com
package factory

import (
	"context"
	"database/sql"
	"encoding/gob"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"sync"

	"github.com/webx-top/db"
	"github.com/webx-top/db/lib/sqlbuilder"
	"github.com/webx-top/tagfast"
)

var (
	tagParser = func(tag string) interface{} {
		if len(tag) == 0 {
			return nil
		}
		return strings.Split(tag, `,`)
	}

	ErrExpectingStruct = errors.New(`bean must be an address of struct or struct`)
	_                  = gob.Register
)

func NewJoin(joinType string, collection string, alias string, condition string) *Join {
	return &Join{
		Collection: collection,
		Alias:      alias,
		Condition:  condition,
		Type:       joinType,
	}
}

var paramPool = sync.Pool{
	New: func() interface{} {
		p := NewParam()
		p.inPool = true
		return p
	},
}

func ParamPoolGet() *Param {
	return paramPool.Get().(*Param)
}

func ParamPoolRelease(c *Param) {
	c.Reset()
	paramPool.Put(c)
}

type Join struct {
	Collection string
	Alias      string
	Condition  string
	Type       string
}

type Param struct {
	inPool             bool
	noRelease          bool
	ctx                context.Context
	factory            *Factory
	index              int //数据库对象元素所在的索引位置
	readOnly           bool
	collection         string //集合名或表名称
	alias              string //表别名
	middleware         func(db.Result) db.Result
	middlewareName     string
	middlewareSelector func(sqlbuilder.Selector) sqlbuilder.Selector
	middlewareTx       func(*Transaction) error
	result             interface{}   //查询后保存的结果
	args               []interface{} //Find方法的条件参数
	cols               []interface{} //使用Selector要查询的列
	joins              []*Join
	save               interface{} //增加和更改数据时要保存到数据库中的数据
	offset             int
	page               int   //页码
	size               int   //每页数据量
	total              int64 //数据表中符合条件的数据行数
	maxAge             int64 //缓存有效时间（单位：秒），为0时代表临时关闭缓存，为-1时代表删除缓存
	trans              *Transaction
	cachedKey          string
	cluster            *Cluster
	model              Model
}

func NewParam(args ...interface{}) *Param {
	p := &Param{
		factory: DefaultFactory,
		page:    1,
		offset:  -1,
	}
	p.init(args...)
	return p
}

func (p *Param) Reset() {
	p.noRelease = false
	p.ctx = nil
	//p.factory = nil
	p.index = 0
	p.readOnly = false
	p.collection = ``
	p.alias = ``
	p.middleware = nil
	p.middlewareName = ``
	p.middlewareSelector = nil
	p.middlewareTx = nil
	p.result = nil
	p.args = nil
	p.cols = nil
	p.joins = nil
	p.save = nil
	p.offset = -1
	p.page = 1
	p.size = 0
	p.total = 0
	p.maxAge = 0
	p.trans = nil
	p.cachedKey = ``
	p.cluster = nil
	p.model = nil
}

func (p *Param) Release() {
	if p.inPool && !p.noRelease {
		ParamPoolRelease(p)
	}
}

func (p *Param) init(args ...interface{}) *Param {
	if len(args) > 0 {
		for _, v := range args {
			if factory, ok := v.(*Factory); ok {
				p.factory = factory
				continue
			}
			if param, ok := v.(*Param); ok {
				p.TransFrom(param)
				continue
			}
			if option, ok := v.(ParamOption); ok {
				option(p)
			}
		}
	}
	return p
}

func (p *Param) Set(opts ...ParamOption) *Param {
	for _, o := range opts {
		o(p)
	}
	return p
}

func (p *Param) SetIndex(index int) *Param {
	p.index = index
	return p
}

func (p *Param) SetContext(ctx context.Context) *Param {
	p.ctx = ctx
	return p
}

func (p *Param) Context() context.Context {
	if p.ctx == nil {
		p.ctx = context.Background()
	}
	return p.ctx
}

func (p *Param) SetModel(model Model) *Param {
	p.model = model
	p.trans = model.Trans()
	if len(p.collection) == 0 {
		p.collection = model.Name_()
	}
	return p
}

func (p *Param) Model() Model {
	return p.model.Use(p.trans).SetParam(p)
}

func (p *Param) SelectLink(index int) *Param {
	p.index = index
	return p
}

func (p *Param) SelectLinkName(name string) *Param {
	p.index = p.factory.IndexByName(name)
	return p
}

func (p *Param) CachedKey() string {
	if len(p.cachedKey) == 0 {
		p.cachedKey = fmt.Sprintf(`%v-%v-%v-%v-%v-%v-%v-%v-%v`, p.index, p.collection, p.cols, p.args, p.offset, p.page, p.size, p.joins, p.middlewareName)
	}
	return p.cachedKey
}

func (p *Param) SetCache(maxAge int64, key ...string) *Param {
	p.maxAge = maxAge
	if len(key) > 0 {
		p.cachedKey = key[0]
	}
	return p
}

func (p *Param) SetCachedKey(key string) *Param {
	p.cachedKey = key
	return p
}

func (p *Param) SetJoin(joins ...*Join) *Param {
	p.joins = joins
	return p
}

func (p *Param) SetTx(tx sqlbuilder.Tx) *Param {
	p.trans = &Transaction{
		tx:      tx,
		factory: p.factory,
	}
	return p
}

func (p *Param) SetTrans(trans *Transaction) *Param {
	p.trans = trans
	return p
}

func (p *Param) SetRead() *Param {
	p.readOnly = true
	return p
}

func (p *Param) SetWrite() *Param {
	p.readOnly = false
	return p
}

func (p *Param) AddJoin(joinType string, collection string, alias string, condition string) *Param {
	p.joins = append(p.joins, NewJoin(joinType, collection, alias, condition))
	return p
}

func (p *Param) SetCollection(collection string, alias ...string) *Param {
	p.collection = collection
	if len(alias) > 0 {
		p.alias = alias[0]
	}
	return p
}

func (p *Param) SetAlias(alias string) *Param {
	p.alias = alias
	return p
}

func (p *Param) TableName() string {
	if len(p.alias) > 0 {
		return p.collection + ` ` + p.alias
	}
	return p.collection
}

func (p *Param) TableField(m interface{}, structField *string, tableField ...*string) *Param {
	var tblField *string
	if len(tableField) > 0 {
		tblField = tableField[0]
	} else {
		tblField = structField
	}
	parts := strings.Split(*structField, `.`)
	j := len(parts)
	rv := reflect.Indirect(reflect.ValueOf(m))
	rt := rv.Type()
	if j == 1 {
		sf, ok := rt.FieldByName(parts[0])
		if !ok {
			*tblField = ``
			return p
		}
		tag := tagfast.GetParsed(rt, sf, `bson`, tagParser)
		if tag == nil {
			tag = tagfast.GetParsed(rt, sf, `db`, tagParser)
		}
		field := parts[0]
		if tags, ok := tag.([]string); ok && len(tags) > 0 && len(tags[0]) > 0 {
			field = tags[0]
		}
		*tblField = field
		if len(p.alias) > 0 {
			*tblField = p.alias + `.` + *tblField
		}
		return p
	}
	var prefix string
	for i, v := range parts {
		if i+1 == j { //end
			sf, ok := rt.FieldByName(v)
			if !ok {
				*tblField = ``
				break
			}
			tag := tagfast.GetParsed(rt, sf, `bson`, tagParser)
			if tag == nil {
				tag = tagfast.GetParsed(rt, sf, `db`, tagParser)
			}
			field := v
			if tags, ok := tag.([]string); ok && len(tags) > 0 && len(tags[0]) > 0 {
				field = tags[0]
			} else {
				field = ToSnakeCase(field)
			}
			*tblField = prefix + field
			break
		}
		sf, ok := rt.FieldByName(v)
		if !ok {
			*tblField = ``
			break
		}
		tag := tagfast.GetParsed(rt, sf, `bson`, tagParser)
		if tag == nil {
			tag = tagfast.GetParsed(rt, sf, `db`, tagParser)
		}

		rv = rv.FieldByName(v)
		if !rv.IsValid() {
			*tblField = ``
			break
		}
		if rv.Kind() == reflect.Ptr {
			rt = rv.Type().Elem()
			if rt.Kind() == reflect.Struct {
				fieldPtr := rv
				rv = rv.Elem()
				if !rv.IsValid() || fieldPtr.IsNil() {
					rv = reflect.New(rt).Elem()
				}
			}
		} else {
			rt = rv.Type()
		}

		var table string
		if tags, ok := tag.([]string); ok && len(tags) > 0 && len(tags[0]) > 0 {
			table = tags[0]
		}
		if len(p.joins) > 0 {
			var rawTableName string
			if len(table) < 1 {
				rawTableName = ToSnakeCase(rt.Name())
			} else {
				rawTableName = table
			}
			for _, jo := range p.joins {
				if jo.Collection == rawTableName {
					if len(jo.Alias) > 0 {
						table = jo.Alias
					} else {
						table = v
					}
					break
				}
			}
		}

		if len(table) == 0 {
			if len(p.alias) > 0 {
				table = p.alias
			} else {
				table = v
			}
		}
		prefix += table + `.`
	}
	return p
}

func (p *Param) SetMiddleware(middleware func(db.Result) db.Result, name ...string) *Param {
	p.middleware = middleware
	if len(name) > 0 {
		p.middlewareName = name[0]
	}
	return p
}

func (p *Param) SetMiddlewareSelector(middleware func(sqlbuilder.Selector) sqlbuilder.Selector, name ...string) *Param {
	p.middlewareSelector = middleware
	if len(name) > 0 {
		p.middlewareName = name[0]
	}
	return p
}

// SetMW is SetMiddleware's alias.
func (p *Param) SetMW(middleware func(db.Result) db.Result, name ...string) *Param {
	p.SetMiddleware(middleware, name...)
	return p
}

func (p *Param) SetMiddlewareTx(middleware func(*Transaction) error) *Param {
	p.middlewareTx = middleware
	return p
}

func (p *Param) SetMWTx(middleware func(*Transaction) error) *Param {
	p.SetMiddlewareTx(middleware)
	return p
}

// SetSelMW is SetSelectorMiddleware's alias.
func (p *Param) SetMWSel(middleware func(sqlbuilder.Selector) sqlbuilder.Selector, name ...string) *Param {
	p.SetMiddlewareSelector(middleware, name...)
	return p
}

func (p *Param) SetRecv(result interface{}) *Param {
	p.result = result
	return p
}

func (p *Param) Recv() interface{} {
	return p.result
}

func (p *Param) SetArgs(args ...interface{}) *Param {
	p.args = args
	return p
}

func (p *Param) AddArgs(args ...interface{}) *Param {
	p.args = append(p.args, args...)
	return p
}

func (p *Param) SetCols(args ...interface{}) *Param {
	p.cols = args
	return p
}

func (p *Param) AddCols(args ...interface{}) *Param {
	p.cols = append(p.cols, args...)
	return p
}

func (p *Param) SetSend(save interface{}) *Param {
	p.save = save
	return p
}

func (p *Param) SetPage(n int) *Param {
	if n < 1 {
		p.page = 1
	} else {
		p.page = n
	}
	return p
}

func (p *Param) SetOffset(offset int) *Param {
	p.offset = offset
	return p
}

func (p *Param) SetSize(size int) *Param {
	p.size = size
	return p
}

func (p *Param) SetTotal(total int64) *Param {
	p.total = total
	return p
}

func (p *Param) Total() int64 {
	return p.total
}

func (p *Param) Trans() *Transaction {
	return p.trans
}

func (p *Param) TransTo(param *Param) *Param {
	param.trans = p.trans
	return p
}

func (p *Param) TransFrom(param *Param) *Param {
	p.trans = param.trans
	return p
}

func (p *Param) GetOffset() int {
	if p.offset > -1 {
		return p.offset
	}
	if p.size < 0 {
		return 0
	}
	if p.page < 1 {
		p.page = 1
	}
	return (p.page - 1) * p.size
}

func (p *Param) NewTx(ctx context.Context) (*Transaction, error) {
	return p.factory.NewTx(ctx, p.index)
}

func (p *Param) Tx(ctx context.Context) error {
	return p.factory.Tx(p, ctx)
}

func (p *Param) MustTx(ctx context.Context) *Transaction {
	trans, err := p.NewTx(ctx)
	if err != nil {
		panic(err.Error())
	}
	return trans
}

func (p *Param) Begin(ctx context.Context) (err error) {
	p.trans, err = p.NewTx(ctx)
	return
}

func (p *Param) MustBegin(ctx context.Context) *Param {
	p.trans = p.MustTx(ctx)
	return p
}

func (p *Param) Rollback(ctx context.Context) error {
	t := p.T()
	if t.tx == nil {
		return nil
	}
	return t.tx.Rollback()
}

func (p *Param) Commit(ctx context.Context) error {
	t := p.T()
	if t.tx == nil {
		return nil
	}
	return t.tx.Commit()
}

func (p *Param) End(ctx context.Context, succeed bool) error {
	if succeed {
		return p.Commit(ctx)
	}
	return p.Rollback(ctx)
}

func (p *Param) T() *Transaction {
	if p.trans != nil {
		return p.trans
	}
	return p.factory.Transaction
}

func (p *Param) Driver() interface{} {
	return p.T().Driver(p)
}

func (p *Param) DB() *sql.DB {
	return p.T().DB(p)
}

func (p *Param) Result() db.Result {
	return p.T().Result(p)
}

func (p *Param) CheckCached() bool {
	return p.T().CheckCached(p)
}

func (p *Param) SQLBuilder() sqlbuilder.SQLBuilder {
	return p.T().SQLBuilder(p)
}

// Read ==========================

// Cached query support cache
func (p *Param) Cached(cachedKey string, fn func(*Param) error, maxAge int64) error {
	p.SetCache(maxAge, cachedKey)
	return p.T().Cached(p, fn)
}

// Query query SQL. sqlRows is an *sql.Rows object, so you can use Scan() on it
// err = sqlRows.Scan(&a, &b, ...)
func (p *Param) Query() (*sql.Rows, error) {
	return p.T().Query(p)
}

// QueryTo query SQL. mapping fields into a struct
func (p *Param) QueryTo() (sqlbuilder.Iterator, error) {
	return p.T().QueryTo(p)
}

// QueryRow query SQL
func (p *Param) QueryRow() (*sql.Row, error) {
	return p.T().QueryRow(p)
}

func (p *Param) SelectAll() error {
	return p.T().SelectAll(p)
}

func (p *Param) SelectOne() error {
	return p.T().SelectOne(p)
}

func (p *Param) SelectCount() (int64, error) {
	return p.T().SelectCount(p)
}

func (p *Param) SelectList() (func() int64, error) {
	return p.T().SelectList(p)
}

func (p *Param) Select() sqlbuilder.Selector {
	return p.T().Select(p)
}

func (p *Param) All() error {
	return p.T().All(p)
}

func (p *Param) List() (func() int64, error) {
	return p.T().List(p)
}

func (p *Param) One() error {
	return p.T().One(p)
}

func (p *Param) Count() (int64, error) {
	return p.T().Count(p)
}

func (p *Param) Exists() (bool, error) {
	return p.T().Exists(p)
}

// Stat Stat(`max`,`score`)
func (p *Param) Stat(fn string, field string) (float64, error) {
	return p.T().Stat(p, fn, field)
}

// Write ==========================

// Exec execute SQL
func (p *Param) Exec() (sql.Result, error) {
	return p.T().Exec(p)
}

func (p *Param) Insert() (interface{}, error) {
	return p.T().Insert(p)
}

// UpdateField 通过设置数据库字段名和字段值来更新数据(无需使用SetSend)
// 支持“key1, value1, key2, value2, ..., keyN, valueN”的顺序赋值
// e.g. UpdateField("user","admin") or UpdateField("user","admin","score",5)
func (p *Param) UpdateField(keyAndValue ...interface{}) error {
	if len(keyAndValue) == 0 {
		return nil
	}
	p.SetSend(keyAndValue)
	return p.T().Update(p)
}

// UpdateByStruct 通过指定Struct中的字段名来更新相应数据(无需使用SetSend)
func (p *Param) UpdateByStruct(bean interface{}, fields ...string) error {
	err := p.UsingStructField(bean)
	if err != nil {
		return err
	}
	return p.T().Update(p)
}

func (p *Param) Update() error {
	return p.T().Update(p)
}

func (p *Param) Updatex() (int64, error) {
	return p.T().Updatex(p)
}

func (p *Param) Upsert(beforeUpsert ...func() error) (interface{}, error) {
	return p.T().Upsert(p, beforeUpsert...)
}

func (p *Param) Delete() error {
	return p.T().Delete(p)
}

func (p *Param) Deletex() (int64, error) {
	return p.T().Deletex(p)
}
