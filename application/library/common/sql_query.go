package common

import (
	"database/sql"
	"fmt"
	"reflect"
	"strings"

	"github.com/admpub/nging/v4/application/dbschema"
	"github.com/admpub/nging/v4/application/library/namedstruct"
	"github.com/admpub/nging/v4/application/response"
	"github.com/admpub/null"
	"github.com/webx-top/db"
	"github.com/webx-top/db/lib/factory"
	"github.com/webx-top/db/lib/factory/pagination"
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

// [link1] SELECT ...
// [conn=default;cacheKey=testCacheKey;cacheTTL=86400;offset=10;limit=10;orderBy=-updated,id] SELECT ...
func (s *SQLQuery) repair(query string) string {
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
		for _, option := range strings.Split(linkName, `;`) {
			parts := strings.SplitN(option, `=`, 2)
			if len(parts) != 2 {
				s.link = factory.IndexByName(linkName)
				continue
			}
			for index, part := range parts {
				parts[index] = strings.TrimSpace(part)
			}
			if len(parts[1]) == 0 {
				continue
			}
			switch parts[0] {
			case `conn`:
				s.link = factory.IndexByName(parts[1])
			case `cacheKey`:
				s.cacheKey = parts[1]
			case `cacheTTL`:
				s.cacheTTL = param.AsInt64(parts[1])
			case `offset`:
				s.offset = param.AsInt(parts[1])
			case `limit`:
				s.limit = param.AsInt(parts[1])
			case `orderBy`:
				for _, sort := range strings.Split(parts[1], `,`) {
					sort = strings.TrimSpace(sort)
					if len(sort) == 0 {
						continue
					}
					s.sorts = append(s.sorts, sort)
				}
			}
		}
	}
	query = strings.TrimLeft(qs[1], ` `)
	return query
}

type SetContext interface {
	SetContext(echo.Context)
}

func (s *SQLQuery) query(name string, recv interface{}, fn func() error, args ...interface{}) error {
	if s.cacher != nil && len(s.cacheKey) > 0 {
		cacheKey := s.cacheKey + fmt.Sprintf(`.%d.%d`, s.offset, s.limit) + `:` + name
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
	return s.query(fmt.Sprintf(`%T$%s`, recv, query), recv, fn, args...)
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

func parseStructName(name string) (structName string, recv interface{}, err error) {
	structName = name
	parts := strings.SplitN(name, `@`, 2) // modelStructName@responseStructName
	if len(parts) != 2 {
		return
	}
	structName = parts[1]
	responseStructName := parts[1]
	recv = response.Registry.Make(responseStructName)
	if recv == nil {
		err = namedstruct.ErrNotExist
	}
	return
}

func makeCond(args []interface{}) *db.Compounds {
	cond := db.NewCompounds()
	var k string
	for i, j := 0, len(args)-1; i <= j; i++ {
		if i%2 == 0 {
			k = param.AsString(args[i])
			continue
		}
		cond.AddKV(k, args[i])
	}
	return cond
}

func (s *SQLQuery) GetModel(name string, args ...interface{}) (interface{}, error) {
	if len(args) == 0 {
		return nil, nil
	}
	name = s.repair(name)
	structName, recv, err := parseStructName(name)
	if err != nil {
		return nil, err
	}
	m := dbschema.DBI.NewModel(structName, s.link)
	m.SetContext(s.ctx)
	fn := func() error {
		var err error
		if recv == nil {
			err = m.Get(func(r db.Result) db.Result {
				return r.OrderBy(s.sorts...)
			}, args...)
		} else {
			err = m.NewParam().SetArgs(args...).SetMiddleware(func(r db.Result) db.Result {
				return r.OrderBy(s.sorts...)
			}).SetRecv(recv).One()
		}
		if err != nil {
			if err == sql.ErrNoRows {
				err = nil
			}
		}
		return err
	}
	if recv == nil {
		return m, s.query(name, m, fn, args...)
	}
	return recv, s.query(name, recv, fn, args...)
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

func (s *SQLQuery) GetModels(name string, args ...interface{}) ([]interface{}, error) {
	if len(args) == 0 {
		return nil, nil
	}
	name = s.repair(name)
	structName, recv, err := parseStructName(name)
	if err != nil {
		return nil, err
	}
	var results []interface{}
	fn := func() error {
		m := dbschema.DBI.NewModel(structName, s.link)
		m.SetContext(s.ctx)
		cond := makeCond(args)
		_, err := m.ListByOffset(recv, func(r db.Result) db.Result {
			return r.OrderBy(s.sorts...)
		}, s.offset, s.limit, cond.And())
		if err != nil {
			return err
		}
		if recv == nil {
			results = ModelObjects(m)
		} else {
			results = namedstruct.ConvertToSlice(recv)
		}
		return err
	}

	return results, s.query(name, &results, fn, args...)
}

func (s *SQLQuery) GetModelsWithPaging(name string, args ...interface{}) ([]interface{}, error) {
	if len(args) == 0 {
		return nil, nil
	}
	name = s.repair(name)
	structName, recv, err := parseStructName(name)
	if err != nil {
		return nil, err
	}
	var results []interface{}
	fn := func() error {
		m := dbschema.DBI.NewModel(structName, s.link)
		m.SetContext(s.ctx)
		cond := makeCond(args)
		_, err := pagination.NewLister(m, recv, func(r db.Result) db.Result {
			return r.OrderBy(s.sorts...)
		}, cond.And()).Paging(m.Context())
		if err != nil {
			return err
		}
		if recv == nil {
			results = ModelObjects(m)
		} else {
			results = namedstruct.ConvertToSlice(recv)
		}
		return nil
	}
	return results, s.query(name, &results, fn, args...)
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

func (s *SQLQuery) MustGetModel(structName string, args ...interface{}) interface{} {
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
	return result, s.query(query, &result, fn, args...)
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
	return result, s.query(query, &result, fn, args...)
}

func (s *SQLQuery) MustGetRows(query string, args ...interface{}) null.StringMapSlice {
	result, err := s.GetRows(query, args...)
	if err != nil {
		panic(err)
	}
	return result
}
