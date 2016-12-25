/*

   Copyright 2016 Wenhui Shen <www.webx.top>

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.

*/
package mysql

import (
	"database/sql"
	"encoding/gob"
	"time"

	"github.com/admpub/nging/application/library/dbmanager/result"
	"github.com/webx-top/db/lib/factory"
)

func init() {
	gob.Register(&Result{})
}

var _ result.Resulter = &Result{}

type Result struct {
	SQL          string
	RowsAffected int64
	TimeStart    time.Time
	TimeEnd      time.Time
	Started      string
	Elapsed      string
	err          error
	ErrorString  string
}

func (r *Result) GetSQL() string {
	return r.SQL
}

func (r *Result) GetBeginTime() string {
	return r.Started
}

func (r *Result) GetElapsedTime() string {
	return r.Elapsed
}
func (r *Result) GetAffected() int64 {
	return r.RowsAffected
}

func (r *Result) GetError() string {
	return r.ErrorString
}

func (r *Result) elapsed() time.Duration {
	return r.TimeEnd.Sub(r.TimeStart)
}

func (r *Result) start() *Result {
	r.TimeStart = time.Now()
	return r
}

func (r *Result) end() *Result {
	r.TimeEnd = time.Now()
	r.Started = r.TimeStart.Format(`2006-01-02 15:04:05`)
	r.Elapsed = r.elapsed().String()
	if r.err != nil {
		r.ErrorString = r.err.Error()
	}
	return r
}

func (r *Result) Exec(p *factory.Param) *Result {
	r.start()
	defer r.end()
	result, err := p.SetCollection(r.SQL).Exec()
	r.err = err
	if err != nil {
		return r
	}
	r.RowsAffected, r.err = result.RowsAffected()
	return r
}

func (r *Result) Query(p *factory.Param, readRows func(*sql.Rows) error) *Result {
	r.start()
	defer r.end()
	rows, err := p.SetCollection(r.SQL).Query()
	r.err = err
	if err != nil {
		return r
	}
	r.err = readRows(rows)
	return r
}
