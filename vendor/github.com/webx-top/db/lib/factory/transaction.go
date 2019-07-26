package factory

import (
	"database/sql"
	"log"
	"strings"

	"github.com/webx-top/db"
	"github.com/webx-top/db/lib/sqlbuilder"
)

type Transaction struct {
	sqlbuilder.Tx
	*Cluster
	*Factory
}

func (t *Transaction) Database(param *Param) db.Database {
	if t.Cluster == nil {
		param.cluster = t.Factory.Cluster(param.Index)
	} else {
		param.cluster = t.Cluster
	}
	if t.Tx != nil {
		if _, ok := t.Tx.Driver().(*sql.Tx); ok {
			return t.Tx
		}
		t.Tx = nil
	}
	if param.ReadOnly {
		return param.cluster.Slave()
	}
	return param.cluster.Master()
}

func (t *Transaction) Driver(param *Param) interface{} {
	return t.Database(param).Driver()
}

func (t *Transaction) DB(param *Param) *sql.DB {
	if db, ok := t.Driver(param).(*sql.DB); ok {
		return db
	}
	panic(db.ErrUnsupported)
}

func (t *Transaction) SQLBuilder(param *Param) sqlbuilder.SQLBuilder {
	if db, ok := t.Database(param).(sqlbuilder.SQLBuilder); ok {
		return db
	}
	panic(db.ErrUnsupported)
}

func (t *Transaction) Result(param *Param) db.Result {
	res := t.C(param).Find(param.Args...)
	if len(param.Cols) > 0 {
		res = res.Select(param.Cols...)
	}
	return res
}

func (t *Transaction) C(param *Param) db.Collection {
	return t.Database(param).Collection(param.TableName())
}

// Exec execute SQL
func (t *Transaction) Exec(param *Param) (sql.Result, error) {
	param.ReadOnly = false
	return t.DB(param).ExecContext(param.Context(), param.Collection, param.Args...)
}

// Query query SQL. sqlRows is an *sql.Rows object, so you can use Scan() on it
// err = sqlRows.Scan(&a, &b, ...)
func (t *Transaction) Query(param *Param) (*sql.Rows, error) {
	return t.DB(param).QueryContext(param.Context(), param.Collection, param.Args...)
}

// QueryTo query SQL. mapping fields into a struct
func (t *Transaction) QueryTo(param *Param) (sqlbuilder.Iterator, error) {
	rows, err := t.Query(param)
	if err != nil {
		return nil, err
	}
	iter := sqlbuilder.NewIterator(t.SQLBuilder(param), rows)
	if param.ResultData != nil {
		err = iter.All(param.ResultData)
	}
	return iter, err
}

// QueryRow query SQL
func (t *Transaction) QueryRow(param *Param) *sql.Row {
	return t.DB(param).QueryRowContext(param.Context(), param.Collection, param.Args...)
}

// ================================
// API
// ================================

// Read ==========================

func (t *Transaction) SelectAll(param *Param) error {

	if t.CheckCached(param) {
		data, err := t.Factory.cacher.Get(param.CachedKey())
		if err == nil && data != nil {
			if v, ok := data.(*Param); ok {
				param = v
				param.factory = t.Factory
				return nil
			}
		}
		defer t.Factory.cacher.Put(param.CachedKey(), param, param.MaxAge)
	}

	selector := t.Select(param)
	if param.Size > 0 {
		selector = selector.Limit(param.Size).Offset(param.GetOffset())
	}
	if param.SelectorMiddleware != nil {
		selector = param.SelectorMiddleware(selector)
	}
	return selector.All(param.ResultData)
}

func (t *Transaction) SelectOne(param *Param) error {

	if t.CheckCached(param) {
		data, err := t.Factory.cacher.Get(param.CachedKey())
		if err == nil && data != nil {
			if v, ok := data.(*Param); ok {
				param = v
				param.factory = t.Factory
				return nil
			}
		}
		defer t.Factory.cacher.Put(param.CachedKey(), param, param.MaxAge)
	}

	selector := t.Select(param).Limit(1)
	if param.SelectorMiddleware != nil {
		selector = param.SelectorMiddleware(selector)
	}
	return selector.One(param.ResultData)
}

func (t *Transaction) SelectList(param *Param) (func() int64, error) {

	if t.CheckCached(param) {
		data, err := t.Factory.cacher.Get(param.CachedKey())
		if err == nil && data != nil {
			if v, ok := data.(*Param); ok {
				param = v
				param.factory = t.Factory
				return func() int64 {
					return param.Total
				}, nil
			}
		}
		defer t.Factory.cacher.Put(param.CachedKey(), param, param.MaxAge)
	}

	selector := t.Select(param).Limit(param.Size).Offset(param.GetOffset())
	if param.SelectorMiddleware != nil {
		selector = param.SelectorMiddleware(selector)
	}
	countFn := func() int64 {
		cnt, err := t.SelectCount(param)
		if err != nil {
			log.Println(err)
		}
		return cnt
	}
	return countFn, selector.All(param.ResultData)
}

func (t *Transaction) SelectCount(param *Param) (int64, error) {
	if param.Total > 0 {
		return param.Total, nil
	}

	if t.CheckCached(param) {
		data, err := t.Factory.cacher.Get(param.CachedKey())
		if err == nil && data != nil {
			if v, ok := data.(*Param); ok {
				param = v
				param.factory = t.Factory
				return param.Total, nil
			}
		}
		defer t.Factory.cacher.Put(param.CachedKey(), param, param.MaxAge)
	}

	counter := struct {
		Count int64 `db:"_t"`
	}{}
	selector := t.SQLBuilder(param).Select(db.Raw("count(1) AS _t")).From(param.TableName()).Where(param.Args...)
	selector = t.joinSelect(param, selector)
	if param.SelectorMiddleware != nil {
		selector = param.SelectorMiddleware(selector)
	}
	selector = selector.Offset(0).Limit(1).OrderBy()
	if err := selector.IteratorContext(param.Context()).One(&counter); err != nil {
		if err == db.ErrNoMoreRows {
			return 0, nil
		}
		return 0, err
	}
	param.Total = counter.Count
	return counter.Count, nil
}

func (t *Transaction) joinSelect(param *Param, selector sqlbuilder.Selector) sqlbuilder.Selector {
	if param.Joins == nil {
		return selector
	}
	for _, join := range param.Joins {
		coll := join.Collection
		if len(join.Alias) > 0 {
			coll += ` ` + join.Alias
		}
		switch strings.ToUpper(join.Type) {
		case "LEFT":
			selector = selector.LeftJoin(coll)
		case "RIGHT":
			selector = selector.RightJoin(coll)
		case "CROSS":
			selector = selector.CrossJoin(coll)
		case "INNER":
			selector = selector.Join(coll)
		default:
			selector = selector.FullJoin(coll)
		}
		if len(join.Condition) > 0 {
			selector = selector.On(join.Condition)
		}
	}
	return selector
}

func (t *Transaction) Select(param *Param) sqlbuilder.Selector {
	selector := t.SQLBuilder(param).Select(param.Cols...).From(param.TableName()).Where(param.Args...)
	return t.joinSelect(param, selector)
}

func (t *Transaction) CheckCached(param *Param) bool {
	if t.Factory.cacher != nil {
		if param.MaxAge > 0 {
			return true
		}
		if param.MaxAge < 0 {
			err := t.Factory.cacher.Del(param.CachedKey())
			if err != nil {
				log.Println(err)
			}
		}
	}

	return false
}

func (t *Transaction) Cached(param *Param, fn func(*Param) error) error {
	if t.CheckCached(param) {
		data, err := t.Factory.cacher.Get(param.CachedKey())
		if err == nil && data != nil {
			if v, ok := data.(*Param); ok {
				param = v
				param.factory = t.Factory
				return nil
			}
		}
		defer t.Factory.cacher.Put(param.CachedKey(), param, param.MaxAge)
	}

	return fn(param)
}

func (t *Transaction) All(param *Param) error {

	if t.CheckCached(param) {
		data, err := t.Factory.cacher.Get(param.CachedKey())
		if err == nil && data != nil {
			if v, ok := data.(*Param); ok {
				param = v
				param.factory = t.Factory
				return nil
			}
		}
		defer t.Factory.cacher.Put(param.CachedKey(), param, param.MaxAge)
	}

	res := t.Result(param)
	if param.Size > 0 {
		res = res.Limit(param.Size).Offset(param.GetOffset())
	}
	if param.Middleware != nil {
		res = param.Middleware(res)
	}
	return res.All(param.ResultData)
}

func (t *Transaction) List(param *Param) (func() int64, error) {

	if t.CheckCached(param) {
		data, err := t.Factory.cacher.Get(param.CachedKey())
		if err == nil && data != nil {
			if v, ok := data.(*Param); ok {
				param = v
				param.factory = t.Factory
				return func() int64 {
					return param.Total
				}, nil
			}
		}
		defer t.Factory.cacher.Put(param.CachedKey(), param, param.MaxAge)
	}

	var res db.Result
	if param.Middleware == nil {
		param.CountFunc = func() int64 {
			if param.Total <= 0 {
				res := t.Result(param)
				count, _ := res.Count()
				param.Total = int64(count)
			}
			return param.Total
		}
		if param.Size >= 0 {
			res = t.Result(param).Limit(param.Size).Offset(param.GetOffset())
		} else {
			res = t.Result(param).Offset(param.GetOffset())
		}
	} else {
		param.CountFunc = func() int64 {
			if param.Total <= 0 {
				res := param.Middleware(t.Result(param)).OrderBy()
				count, _ := res.Count()
				param.Total = int64(count)
			}
			return param.Total
		}
		if param.Size >= 0 {
			res = param.Middleware(t.Result(param).Limit(param.Size).Offset(param.GetOffset()))
		} else {
			res = param.Middleware(t.Result(param).Offset(param.GetOffset()))
		}
	}
	return param.CountFunc, res.All(param.ResultData)
}

func (t *Transaction) One(param *Param) error {

	if t.CheckCached(param) {
		data, err := t.Factory.cacher.Get(param.CachedKey())
		if err == nil && data != nil {
			if v, ok := data.(*Param); ok {
				param = v
				param.factory = t.Factory
				return nil
			}
		}
		defer t.Factory.cacher.Put(param.CachedKey(), param, param.MaxAge)
	}

	res := t.Result(param)
	if param.Middleware != nil {
		res = param.Middleware(res)
	}
	return res.One(param.ResultData)
}

func (t *Transaction) Count(param *Param) (int64, error) {

	if t.CheckCached(param) {
		data, err := t.Factory.cacher.Get(param.CachedKey())
		if err == nil && data != nil {
			if v, ok := data.(*Param); ok {
				param = v
				param.factory = t.Factory
				return param.Total, nil
			}
		}
		defer t.Factory.cacher.Put(param.CachedKey(), param, param.MaxAge)
	}

	var cnt uint64
	var err error

	res := t.Result(param)
	if param.Middleware != nil {
		res = param.Middleware(res)
	}
	cnt, err = res.Count()
	param.Total = int64(cnt)
	return param.Total, err
}

// Write ==========================

func (t *Transaction) Insert(param *Param) (interface{}, error) {
	param.ReadOnly = false
	return t.C(param).Insert(param.SaveData)
}

func (t *Transaction) Update(param *Param) error {
	param.ReadOnly = false
	res := t.Result(param)
	if param.Middleware != nil {
		res = param.Middleware(res)
	}
	return res.Update(param.SaveData)
}

func (t *Transaction) Updatex(param *Param) (affected int64, err error) {
	param.ReadOnly = false
	res, err := t.SQLBuilder(param).Update(param.TableName()).Set(param.SaveData).Where(param.Args...).ExecContext(param.Context())
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

func (t *Transaction) Upsert(param *Param, beforeUpsert ...func()) (interface{}, error) {
	param.ReadOnly = false
	res := t.Result(param)
	if param.Middleware != nil {
		res = param.Middleware(res)
	}
	cnt, err := res.Count()
	if err != nil {
		if err == db.ErrNoMoreRows {
			if len(beforeUpsert) > 1 && beforeUpsert[1] != nil {
				beforeUpsert[1]()
			}
			return t.C(param).Insert(param.SaveData)
		}
		return nil, err
	}
	if cnt < 1 {
		if len(beforeUpsert) > 1 && beforeUpsert[1] != nil {
			beforeUpsert[1]()
		}
		return t.C(param).Insert(param.SaveData)
	}
	if len(beforeUpsert) > 0 && beforeUpsert[0] != nil {
		beforeUpsert[0]()
	}
	return nil, t.Update(param)
}

func (t *Transaction) Delete(param *Param) error {
	param.ReadOnly = false
	res := t.Result(param)
	if param.Middleware != nil {
		res = param.Middleware(res)
	}
	return res.Delete()
}

func (t *Transaction) Deletex(param *Param) (affected int64, err error) {
	param.ReadOnly = false
	res, err := t.SQLBuilder(param).DeleteFrom(param.TableName()).Where(param.Args...).ExecContext(param.Context())
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}
