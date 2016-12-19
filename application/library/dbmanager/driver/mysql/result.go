package mysql

import (
	"database/sql"
	"time"

	"github.com/webx-top/db/lib/factory"
)

type Result struct {
	SQL          string
	RowsAffected int64
	timeStart    time.Time
	timeEnd      time.Time
	Started      string
	Elapsed      string
	Error        error
}

func (r *Result) elapsed() time.Duration {
	return r.timeEnd.Sub(r.timeStart)
}

func (r *Result) start() *Result {
	r.timeStart = time.Now()
	return r
}

func (r *Result) end() *Result {
	r.timeEnd = time.Now()
	r.Started = r.timeStart.Format(`2006-01-02 15:04:05`)
	r.Elapsed = r.elapsed().String()
	return r
}

func (r *Result) Exec(p *factory.Param) *Result {
	r.start()
	defer r.end()
	result, err := p.SetCollection(r.SQL).Exec()
	r.Error = err
	if err != nil {
		return r
	}
	r.RowsAffected, r.Error = result.RowsAffected()
	return r
}

func (r *Result) Query(p *factory.Param, readRows func(*sql.Rows) error) *Result {
	r.start()
	defer r.end()
	rows, err := p.SetCollection(r.SQL).Query()
	r.Error = err
	if err != nil {
		return r
	}
	r.Error = readRows(rows)
	return r
}
