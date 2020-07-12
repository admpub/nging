package common

import (
	"database/sql"

	"github.com/admpub/null"
	"github.com/webx-top/db/lib/factory"
)

func SQLQueryLimit(limit int, linkIDs ...int) sqlQuery {
	var linkID int
	if len(linkIDs) > 0 {
		linkID = linkIDs[0]
	}
	return sqlQuery{link: linkID, limit: limit}
}

func SQLQuery() sqlQuery {
	return sqlQuery{}
}

type sqlQuery struct {
	link  int
	limit int
}

func (s sqlQuery) Link(dbLink int) sqlQuery {
	r := sqlQuery{
		link:  dbLink,
		limit: s.limit,
	}
	return r
}

func (s sqlQuery) Limit(limit int) sqlQuery {
	r := sqlQuery{
		link:  s.link,
		limit: limit,
	}
	return r
}

// GetValue 查询单个字段值
func (s sqlQuery) GetValue(query string, args ...interface{}) (null.String, error) {
	result := null.String{}
	row := factory.NewParam().SetIndex(s.link).DB().QueryRow(query, args...)
	err := row.Scan(&result)
	return result, err
}

// GetRow 查询一行多个字段值
func (s sqlQuery) GetRow(query string, args ...interface{}) (null.StringMap, error) {
	result := null.StringMap{}
	rows, err := factory.NewParam().SetIndex(s.link).DB().Query(query, args...)
	if err != nil {
		return result, err
	}
	defer rows.Close()
	columns, err := rows.Columns()
	if err != nil {
		return result, err
	}
	size := len(columns)
	if rows.Next() {
		values := make([]interface{}, size)
		for k := range columns {
			values[k] = &null.String{}
		}
		err = rows.Scan(values...)
		if err != nil {
			return result, err
		}
		for k, colName := range columns {
			result[colName] = *values[k].(*null.String)
		}
	}
	return result, err
}

// GetRows 查询多行
func (s sqlQuery) GetRows(query string, args ...interface{}) (null.StringMapSlice, error) {
	result := null.StringMapSlice{}
	rows, err := factory.NewParam().SetIndex(s.link).DB().Query(query, args...)
	if err != nil {
		return result, err
	}
	defer rows.Close()
	columns, err := rows.Columns()
	if err != nil {
		return result, err
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
				return result, err
			}
		}
	} else {
		for rows.Next() {
			err = read(rows)
			if err != nil {
				return result, err
			}
		}
	}
	return result, err
}
