package factory

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	"github.com/webx-top/db"
	"github.com/webx-top/db/lib/sqlbuilder"
)

type Transaction struct {
	tx      sqlbuilder.Tx
	cluster *Cluster
	factory *Factory
}

func (t *Transaction) Database(param *Param) db.Database {
	if t.cluster == nil {
		param.cluster = t.factory.Cluster(param.index)
	} else {
		param.cluster = t.cluster
	}
	if t.tx != nil {
		if _, ok := t.tx.Driver().(*sql.Tx); ok {
			return t.tx
		}
		t.tx = nil
	}
	if param.readOnly {
		return param.cluster.Slave()
	}
	return param.cluster.Master()
}

func (t *Transaction) Driver(param *Param) interface{} {
	return t.Database(param).Driver()
}

func (t *Transaction) DB(param *Param) *sql.DB {
	d := t.Driver(param)
	if db, ok := d.(*sql.DB); ok {
		return db
	}
	panic(fmt.Sprintf(`%v: %T`, db.ErrUnsupported.Error(), d))
}

func (t *Transaction) SQLBuilder(param *Param) sqlbuilder.SQLBuilder {
	d := t.Database(param)
	if db, ok := d.(sqlbuilder.SQLBuilder); ok {
		return db
	}
	panic(fmt.Sprintf(`%v: %T`, db.ErrUnsupported.Error(), d))
}

func (t *Transaction) result(param *Param) db.Result {
	res := t.C(param).Find(param.args...)
	if len(param.cols) > 0 {
		res = res.Select(param.cols...)
	}
	return res
}

func (t *Transaction) Result(param *Param) db.Result {
	res := t.result(param)
	if param.middleware != nil {
		res = param.middleware(res)
	}
	return res
}

func (t *Transaction) C(param *Param) db.Collection {
	return t.Database(param).Collection(param.TableName())
}

// Exec execute SQL
func (t *Transaction) Exec(param *Param) (sql.Result, error) {
	param.readOnly = false
	return t.SQLBuilder(param).ExecContext(param.Context(), param.collection, param.args...)
}

// Query query SQL. sqlRows is an *sql.Rows object, so you can use Scan() on it
// err = sqlRows.Scan(&a, &b, ...)
func (t *Transaction) Query(param *Param) (*sql.Rows, error) {
	return t.SQLBuilder(param).QueryContext(param.Context(), param.collection, param.args...)
}

// QueryTo query SQL. mapping fields into a struct
func (t *Transaction) QueryTo(param *Param) (sqlbuilder.Iterator, error) {
	rows, err := t.Query(param)
	if err != nil {
		return nil, err
	}
	iter := sqlbuilder.NewIterator(t.SQLBuilder(param), rows)
	if param.result != nil {
		err = iter.All(param.result)
	}
	return iter, err
}

// QueryRow query SQL
func (t *Transaction) QueryRow(param *Param) (*sql.Row, error) {
	return t.SQLBuilder(param).QueryRowContext(param.Context(), param.collection, param.args...)
}

// ================================
// API
// ================================

// Read ==========================

func (t *Transaction) SelectAll(param *Param) error {
	selector := t.Select(param)
	if param.size > 0 {
		selector = selector.Limit(param.size).Offset(param.GetOffset())
	}
	if param.middlewareSelector != nil {
		selector = param.middlewareSelector(selector)
	}
	return selector.All(param.result)
}

func (t *Transaction) SelectOne(param *Param) error {
	selector := t.Select(param).Limit(1)
	if param.middlewareSelector != nil {
		selector = param.middlewareSelector(selector)
	}
	return selector.One(param.result)
}

func (t *Transaction) SelectList(param *Param) (func() int64, error) {
	selector := t.Select(param).Limit(param.size).Offset(param.GetOffset())
	if param.middlewareSelector != nil {
		selector = param.middlewareSelector(selector)
	}
	countFn := func() int64 {
		cnt, err := t.SelectCount(param)
		if err != nil {
			log.Println(err)
		}
		return cnt
	}
	return countFn, selector.All(param.result)
}

func (t *Transaction) SelectCount(param *Param) (int64, error) {
	if param.total > 0 {
		return param.total, nil
	}
	counter := struct {
		Count int64 `db:"_t"`
	}{}
	selector := t.SQLBuilder(param).Select(db.Raw("count(1) AS _t")).From(param.TableName()).Where(param.args...)
	selector = t.joinSelect(param, selector)
	if param.middlewareSelector != nil {
		selector = param.middlewareSelector(selector)
	}
	selector = selector.Offset(0).Limit(1).OrderBy()
	if err := selector.IteratorContext(param.Context()).One(&counter); err != nil {
		if err == db.ErrNoMoreRows {
			return 0, nil
		}
		return 0, err
	}
	param.total = counter.Count
	return counter.Count, nil
}

func (t *Transaction) joinSelect(param *Param, selector sqlbuilder.Selector) sqlbuilder.Selector {
	if param.joins == nil {
		return selector
	}
	for _, join := range param.joins {
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
	selector := t.SQLBuilder(param).Select(param.cols...).From(param.TableName()).Where(param.args...)
	return t.joinSelect(param, selector)
}

func (t *Transaction) CheckCached(param *Param) bool {
	if t.factory.cacher != nil {
		if param.maxAge > 0 {
			return true
		}
		if param.maxAge < 0 {
			err := t.factory.cacher.Del(param.CachedKey())
			if err != nil {
				log.Println(err)
			}
		}
	}

	return false
}

func (t *Transaction) Cached(param *Param, fn func(*Param) error) error {
	if !t.CheckCached(param) {
		return fn(param)
	}
	return t.factory.cacher.Do(param.CachedKey(), param.result, func() error {
		return fn(param)
	}, param.maxAge)
}

func (t *Transaction) All(param *Param) error {
	res := t.Result(param)
	if param.size > 0 {
		res = res.Limit(param.size).Offset(param.GetOffset())
	}
	return res.All(param.result)
}

func (t *Transaction) List(param *Param) (func() int64, error) {
	var res db.Result
	cnt := func() int64 {
		if param.total <= 0 {
			param.total, _ = t.count(param)
		}
		return param.total
	}
	if param.size >= 0 {
		res = t.Result(param).Limit(param.size).Offset(param.GetOffset())
	} else {
		res = t.Result(param).Offset(param.GetOffset())
	}
	return cnt, res.All(param.result)
}

func (t *Transaction) One(param *Param) error {
	res := t.Result(param)
	return res.One(param.result)
}

func (t *Transaction) count(param *Param) (int64, error) {
	res := t.Result(param).OrderBy()
	cnt, err := res.Count()
	return int64(cnt), err
}

func (t *Transaction) Count(param *Param) (int64, error) {
	var err error
	param.total, err = t.count(param)
	return param.total, err
}

func (t *Transaction) Exists(param *Param) (bool, error) {
	return t.Result(param).OrderBy().Exists()
}

// Write ==========================

func (t *Transaction) Insert(param *Param) (interface{}, error) {
	param.readOnly = false
	return t.C(param).Insert(param.save)
}

func (t *Transaction) Update(param *Param) error {
	param.readOnly = false
	res := t.Result(param)
	return res.Update(param.save)
}

func (t *Transaction) Updatex(param *Param) (affected int64, err error) {
	param.readOnly = false
	res, err := t.SQLBuilder(param).Update(param.TableName()).Set(param.save).Where(param.args...).ExecContext(param.Context())
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

func (t *Transaction) Upsert(param *Param, beforeUpsert ...func() error) (interface{}, error) {
	param.readOnly = false
	res := t.Result(param)
	cnt, err := res.Count()
	if err != nil {
		if err != db.ErrNoMoreRows {
			return nil, err
		}
		if len(beforeUpsert) > 1 && beforeUpsert[1] != nil {
			if err = beforeUpsert[1](); err != nil {
				return nil, err
			}
		}
		return t.C(param).Insert(param.save)
	}
	if cnt < 1 {
		if len(beforeUpsert) > 1 && beforeUpsert[1] != nil {
			if err = beforeUpsert[1](); err != nil {
				return nil, err
			}
		}
		return t.C(param).Insert(param.save)
	}
	if len(beforeUpsert) > 0 && beforeUpsert[0] != nil {
		if err = beforeUpsert[0](); err != nil {
			return nil, err
		}
	}
	return nil, res.Update(param.save)
}

func (t *Transaction) Delete(param *Param) error {
	param.readOnly = false
	res := t.Result(param)
	return res.Delete()
}

func (t *Transaction) Deletex(param *Param) (affected int64, err error) {
	param.readOnly = false
	res, err := t.SQLBuilder(param).DeleteFrom(param.TableName()).Where(param.args...).ExecContext(param.Context())
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

// Stat Stat(param,`max`,`score`)
func (t *Transaction) Stat(param *Param, fn string, field string) (float64, error) {
	counter := struct {
		Stat sql.NullFloat64 `db:"_t"`
	}{}
	selector := t.SQLBuilder(param).Select(db.Raw(fn + "(" + field + ") AS _t")).From(param.TableName()).Where(param.args...)
	selector = t.joinSelect(param, selector)
	if param.middlewareSelector != nil {
		selector = param.middlewareSelector(selector)
	}
	selector = selector.Offset(0).Limit(1).OrderBy()
	if err := selector.IteratorContext(param.Context()).One(&counter); err != nil {
		if err == db.ErrNoMoreRows {
			return 0, nil
		}
		return 0, err
	}
	return counter.Stat.Float64, nil
}
