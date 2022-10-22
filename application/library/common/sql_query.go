package common

import (
	"database/sql"
	"fmt"
	"reflect"
	"strings"

	"github.com/admpub/nging/v4/application/dbschema"
	"github.com/admpub/null"
	"github.com/webx-top/db"
	"github.com/webx-top/db/lib/factory"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/param"
)

func NewSQLQueryLimit(ctx echo.Context, offset int, limit int, linkIDs ...int) *SQLQuery {
	var linkID int
	if len(linkIDs) > 0 {
		linkID = linkIDs[0]
	}
	return &SQLQuery{
		ctx:    ctx,
		link:   linkID,
		offset: offset,
		limit:  limit,
		cacher: factory.GetCacher(),
	}
}

func NewSQLQuery(ctx echo.Context) *SQLQuery {
	return &SQLQuery{
		ctx:    ctx,
		cacher: factory.GetCacher(),
	}
}

type SQLQuery struct {
	ctx      echo.Context
	link     int
	offset   int
	limit    int
	sorts    []interface{}
	cacheKey string
	cacheTTL int64
	cacher   factory.Cacher
}

func (s *SQLQuery) LinkID(dbLinkID int) *SQLQuery {
	s.link = dbLinkID
	return s
}

func (s *SQLQuery) LinkName(dbLinkName string) *SQLQuery {
	s.link = factory.IndexByName(dbLinkName)
	return s
}

func (s *SQLQuery) Limit(limit int) *SQLQuery {
	s.limit = limit
	return s
}

func (s *SQLQuery) Offset(offset int) *SQLQuery {
	s.offset = offset
	return s
}

func (s *SQLQuery) OrderBy(sorts ...interface{}) *SQLQuery {
	s.sorts = sorts
	return s
}

func (s *SQLQuery) CacheKey(cacheKey string) *SQLQuery {
	s.cacheKey = cacheKey
	return s
}

func (s *SQLQuery) CacheTTL(ttlSeconds int64) *SQLQuery {
	s.cacheTTL = ttlSeconds
	return s
}

func (s *SQLQuery) repair(query string) string { //[link1] SELECT ...
	query = strings.TrimSpace(query)
	if len(query) < 3 || query[0] != '[' {
		return query
	}
	_query := query[1:]
	qs := strings.SplitN(_query, `]`, 2)
	if len(qs) != 2 {
		return query
	}
	linkName := qs[0]
	if len(linkName) > 0 {
		s.link = factory.IndexByName(linkName)
	}
	query = strings.TrimLeft(qs[1], ` `)
	return query
}

type SetContext interface {
	SetContext(echo.Context)
}

func (s *SQLQuery) query(recv interface{}, fn func() error, args ...interface{}) error {
	if s.cacher != nil && len(s.cacheKey) > 0 {
		cacheKey := s.cacheKey + fmt.Sprintf(`.%d.%d`, s.offset, s.limit) + `:` + strings.TrimSuffix(fmt.Sprintf(`%T`, recv), ` {}`)
		if len(args) > 0 {
			format := strings.Repeat(`%+v,`, len(args))
			cacheKey += `:args(` + fmt.Sprintf(format, args...) + `)`
		}
		//cacheKey=com.Md5(cacheKey)
		defer func() {
			if sc, ok := recv.(SetContext); ok {
				sc.SetContext(s.ctx)
			}
		}()
		return s.cacher.Do(`SQLQuery.`+cacheKey, recv, fn, s.cacheTTL)
	}
	return fn()
}

// GetValue 查询单个字段值
func (s *SQLQuery) GetValue(recv interface{}, query string, args ...interface{}) error {
	fn := func() error {
		query = s.repair(query)
		row := factory.NewParam().SetIndex(s.link).DB().QueryRow(query, args...)
		err := row.Scan(recv)
		if err != nil {
			if err == sql.ErrNoRows {
				err = nil
			}
		}
		return err
	}
	return s.query(recv, fn, args...)
}

func (s *SQLQuery) GetString(query string, args ...interface{}) (null.String, error) {
	result := null.String{}
	err := s.GetValue(&result, query, args...)
	return result, err
}

func (s *SQLQuery) GetUint(query string, args ...interface{}) (null.Uint, error) {
	result := null.Uint{}
	err := s.GetValue(&result, query, args...)
	return result, err
}

func (s *SQLQuery) GetUint32(query string, args ...interface{}) (null.Uint32, error) {
	result := null.Uint32{}
	err := s.GetValue(&result, query, args...)
	return result, err
}

func (s *SQLQuery) GetUint64(query string, args ...interface{}) (null.Uint64, error) {
	result := null.Uint64{}
	err := s.GetValue(&result, query, args...)
	return result, err
}

func (s *SQLQuery) GetInt(query string, args ...interface{}) (null.Int, error) {
	result := null.Int{}
	err := s.GetValue(&result, query, args...)
	return result, err
}

func (s *SQLQuery) GetInt32(query string, args ...interface{}) (null.Int32, error) {
	result := null.Int32{}
	err := s.GetValue(&result, query, args...)
	return result, err
}

func (s *SQLQuery) GetInt64(query string, args ...interface{}) (null.Int64, error) {
	result := null.Int64{}
	err := s.GetValue(&result, query, args...)
	return result, err
}

func (s *SQLQuery) GetFloat64(query string, args ...interface{}) (null.Float64, error) {
	result := null.Float64{}
	err := s.GetValue(&result, query, args...)
	return result, err
}

func (s *SQLQuery) GetModel(structName string, args ...interface{}) (factory.Model, error) {
	m := dbschema.DBI.NewModel(structName, s.link)
	m.SetContext(s.ctx)
	if len(args) == 0 {
		return m, nil
	}
	fn := func() error {
		err := m.Get(func(r db.Result) db.Result {
			return r.OrderBy(s.sorts...)
		}, args...)
		if err != nil {
			if err == sql.ErrNoRows {
				err = nil
			}
		}
		return err
	}
	return m, s.query(m, fn, args...)
}

func ModelObjects(m factory.Model) []interface{} {
	rv := reflect.ValueOf(m).MethodByName(`Objects`)
	if rv.Kind() != reflect.Func {
		return nil
	}
	values := rv.Call(nil)
	rows := make([]interface{}, len(values))
	for index, value := range values {
		rows[index] = value.Interface()
	}
	return rows
}

func (s *SQLQuery) GetModels(structName string, args ...interface{}) ([]interface{}, error) {
	if len(args) == 0 {
		return nil, nil
	}
	var results []interface{}
	fn := func() error {
		m := dbschema.DBI.NewModel(structName, s.link)
		m.SetContext(s.ctx)
		cond := db.NewCompounds()
		var k string
		for i, j := 0, len(args)-1; i <= j; i++ {
			if i%2 == 0 {
				k = param.AsString(args[i])
				continue
			}
			cond.AddKV(k, args[i])
		}
		_, err := m.ListByOffset(nil, func(r db.Result) db.Result {
			return r.OrderBy(s.sorts...)
		}, s.offset, s.limit, cond.And())
		results = ModelObjects(m)
		return err
	}

	return results, s.query(&results, fn, args...)
}

func (s *SQLQuery) GetModelsWithPaging(structName string, args ...interface{}) ([]interface{}, error) {
	if len(args) == 0 {
		return nil, nil
	}
	var results []interface{}
	fn := func() error {
		m := dbschema.DBI.NewModel(structName, s.link)
		m.SetContext(s.ctx)
		cond := db.NewCompounds()
		var k string
		for i, j := 0, len(args)-1; i <= j; i++ {
			if i%2 == 0 {
				k = param.AsString(args[i])
				continue
			}
			cond.AddKV(k, args[i])
		}
		err := m.ListPage(cond, s.sorts...)
		if err != nil {
			return err
		}
		results = ModelObjects(m)
		return nil
	}
	return results, s.query(&results, fn, args...)
}

func (s *SQLQuery) MustGetString(query string, args ...interface{}) null.String {
	result, err := s.GetString(query, args...)
	if err != nil {
		panic(err)
	}
	return result
}

func (s *SQLQuery) MustGetUint(query string, args ...interface{}) null.Uint {
	result, err := s.GetUint(query, args...)
	if err != nil {
		panic(err)
	}
	return result
}

func (s *SQLQuery) MustGetUint32(query string, args ...interface{}) null.Uint32 {
	result, err := s.GetUint32(query, args...)
	if err != nil {
		panic(err)
	}
	return result
}

func (s *SQLQuery) MustGetUint64(query string, args ...interface{}) null.Uint64 {
	result, err := s.GetUint64(query, args...)
	if err != nil {
		panic(err)
	}
	return result
}

func (s *SQLQuery) MustGetInt(query string, args ...interface{}) null.Int {
	result, err := s.GetInt(query, args...)
	if err != nil {
		panic(err)
	}
	return result
}

func (s *SQLQuery) MustGetInt32(query string, args ...interface{}) null.Int32 {
	result, err := s.GetInt32(query, args...)
	if err != nil {
		panic(err)
	}
	return result
}

func (s *SQLQuery) MustGetInt64(query string, args ...interface{}) null.Int64 {
	result, err := s.GetInt64(query, args...)
	if err != nil {
		panic(err)
	}
	return result
}

func (s *SQLQuery) MustGetFloat64(query string, args ...interface{}) null.Float64 {
	result, err := s.GetFloat64(query, args...)
	if err != nil {
		panic(err)
	}
	return result
}

func (s *SQLQuery) MustGetModel(structName string, args ...interface{}) factory.Model {
	result, err := s.GetModel(structName, args...)
	if err != nil {
		panic(err)
	}
	return result
}

func (s *SQLQuery) MustGetModels(structName string, args ...interface{}) []interface{} {
	result, err := s.GetModels(structName, args...)
	if err != nil {
		panic(err)
	}
	return result
}

func (s *SQLQuery) MustGetModelsWithPaging(structName string, args ...interface{}) []interface{} {
	result, err := s.GetModelsWithPaging(structName, args...)
	if err != nil {
		panic(err)
	}
	return result
}

// GetRow 查询一行多个字段值
func (s *SQLQuery) GetRow(query string, args ...interface{}) (null.StringMap, error) {
	result := null.StringMap{}
	fn := func() error {
		query = s.repair(query)
		rows, err := factory.NewParam().SetIndex(s.link).DB().Query(query, args...)
		if err != nil {
			return err
		}
		defer rows.Close()
		columns, err := rows.Columns()
		if err != nil {
			return err
		}
		size := len(columns)
		if rows.Next() {
			values := make([]interface{}, size)
			for k := range columns {
				values[k] = &null.String{}
			}
			err = rows.Scan(values...)
			if err != nil {
				if err == sql.ErrNoRows {
					err = nil
				}
				return err
			}
			for k, colName := range columns {
				result[colName] = *values[k].(*null.String)
			}
		}
		return err
	}
	return result, s.query(&result, fn, args...)
}

func (s *SQLQuery) MustGetRow(query string, args ...interface{}) null.StringMap {
	result, err := s.GetRow(query, args...)
	if err != nil {
		panic(err)
	}
	return result
}

// GetRows 查询多行
func (s *SQLQuery) GetRows(query string, args ...interface{}) (null.StringMapSlice, error) {
	result := null.StringMapSlice{}
	fn := func() error {
		query = s.repair(query)
		rows, err := factory.NewParam().SetIndex(s.link).DB().Query(query, args...)
		if err != nil {
			return err
		}
		defer rows.Close()
		columns, err := rows.Columns()
		if err != nil {
			return err
		}
		size := len(columns)
		read := func(rows *sql.Rows) error {
			values := make([]interface{}, size)
			for k := range columns {
				values[k] = &null.String{}
			}
			err := rows.Scan(values...)
			if err != nil {
				return err
			}
			v := null.StringMap{}
			for k, colName := range columns {
				v[colName] = *values[k].(*null.String)
			}
			result = append(result, v)
			return nil
		}
		if s.limit > 0 {
			for i := 0; i < s.limit && rows.Next(); i++ {
				err = read(rows)
				if err != nil {
					if err == sql.ErrNoRows {
						err = nil
					}
					return err
				}
			}
		} else {
			for rows.Next() {
				err = read(rows)
				if err != nil {
					if err == sql.ErrNoRows {
						err = nil
					}
					return err
				}
			}
		}
		return err
	}
	return result, s.query(&result, fn, args...)
}

func (s *SQLQuery) MustGetRows(query string, args ...interface{}) null.StringMapSlice {
	result, err := s.GetRows(query, args...)
	if err != nil {
		panic(err)
	}
	return result
}
