// Copyright (c) 2012-present The upper.io/db authors. All rights reserved.
//
// Permission is hereby granted, free of charge, to any person obtaining
// a copy of this software and associated documentation files (the
// "Software"), to deal in the Software without restriction, including
// without limitation the rights to use, copy, modify, merge, publish,
// distribute, sublicense, and/or sell copies of the Software, and to
// permit persons to whom the Software is furnished to do so, subject to
// the following conditions:
//
// The above copyright notice and this permission notice shall be
// included in all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
// EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
// MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
// NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE
// LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION
// OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION
// WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

package db

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"regexp"
	"strings"
	"time"

	"github.com/webx-top/com"
)

const (
	fmtLogSessID       = `Session ID:     %05d`
	fmtLogTxID         = `Transaction ID: %05d`
	fmtLogQuery        = `Query:          %s`
	fmtLogArgs         = `Arguments:      %v`
	fmtLogRowsAffected = `Rows affected:  %d`
	fmtLogLastInsertID = `Last insert ID: %d`
	fmtLogError        = `Error:          %v`
	fmtLogTimeTaken    = `Time taken:     %0.5fs`
	fmtLogContext      = `Context:        %v`
)

var (
	reInvisibleChars = regexp.MustCompile(`[\s\r\n\t]+`)
)

// QueryStatus represents the status of a query after being executed.
type QueryStatus struct {
	SessID uint64
	TxID   uint64

	RowsAffected *int64
	LastInsertID *int64

	Query string
	Args  []interface{}

	Err error

	Start time.Time
	End   time.Time
	Slow  bool

	Context context.Context
}

func sqlValue(v interface{}) interface{} {
	switch vv := v.(type) {
	case uint, int, uint8, int8, uint16, int16, uint32, int32, uint64, int64, float32, float64:
		return vv
	case string:
		vv = com.AddSlashes(vv)
		return `'` + vv + `'`
	default:
		rv := reflect.ValueOf(vv)
		rv = reflect.Indirect(rv)
		switch rv.Kind() {
		case reflect.Uint, reflect.Int, reflect.Uint8, reflect.Int8, reflect.Uint16, reflect.Int16, reflect.Uint32, reflect.Int32, reflect.Uint64, reflect.Int64, reflect.Float32, reflect.Float64:
			return rv.Interface()
		case reflect.String:
			return `'` + com.AddSlashes(rv.String()) + `'`
		default:
			iv := rv.Interface()
			if s, y := iv.(fmt.Stringer); y {
				return `'` + com.AddSlashes(s.String()) + `'`
			}
			s := fmt.Sprint(iv)
			s = com.AddSlashes(s)
			return `'` + s + `'`
		}
	}
}

// BuildSQL build sql query
func BuildSQL(query string, args ...interface{}) string {
	if len(query) == 0 {
		return query
	}
	query = reInvisibleChars.ReplaceAllString(query, ` `)
	query = strings.TrimSpace(query)
	if len(args) > 0 {
		if len(args) == 1 {
			if v, y := args[0].([]interface{}); y {
				args = v
			}
		}
		newArgs := make([]interface{}, len(args))
		for k, v := range args {
			newArgs[k] = sqlValue(v)
		}
		query = fmt.Sprintf(strings.Replace(query, `?`, `%v`, -1), newArgs...)
	}
	return query
}

// String returns a formatted log message.
func (q *QueryStatus) Lines() []string {
	lines := make([]string, 0, 9)
	if q.SessID > 0 {
		lines = append(lines, fmt.Sprintf(fmtLogSessID, q.SessID))
	}

	if q.TxID > 0 {
		lines = append(lines, fmt.Sprintf(fmtLogTxID, q.TxID))
	}

	if query := q.Query; len(query) > 0 {
		query = BuildSQL(query, q.Args...)
		lines = append(lines, fmt.Sprintf(fmtLogQuery, query))
	}

	if len(q.Args) > 0 {
		b, _ := json.Marshal(q.Args)
		lines = append(lines, fmt.Sprintf(fmtLogArgs, string(b)))
	}

	if q.RowsAffected != nil {
		lines = append(lines, fmt.Sprintf(fmtLogRowsAffected, *q.RowsAffected))
	}
	if q.LastInsertID != nil {
		lines = append(lines, fmt.Sprintf(fmtLogLastInsertID, *q.LastInsertID))
	}

	if q.Err != nil {
		lines = append(lines, fmt.Sprintf(fmtLogError, q.Err))
	}

	lines = append(lines, fmt.Sprintf(fmtLogTimeTaken, float64(q.End.UnixNano()-q.Start.UnixNano())/float64(1e9)))

	if q.Context != nil {
		var ctx interface{}
		switch v := q.Context.(type) {
		case RequestURI:
			if m, ok := v.(Method); ok {
				ctx = `[` + m.Method() + `] ` + v.RequestURI()
			} else {
				ctx = v.RequestURI()
			}
		case StdContext:
			ctx = v.StdContext()
		default:
			ctx = v
		}
		lines = append(lines, fmt.Sprintf(fmtLogContext, ctx))
	}
	return lines
}

// String returns a formatted log message.
func (q *QueryStatus) String() string {
	return q.Stringify("\n")
}

func (q *QueryStatus) Stringify(sep string) string {
	return strings.Join(q.Lines(), sep)
}

// EnvEnableDebug can be used by adapters to determine if the user has enabled
// debugging.
//
// If the user sets the `UPPERIO_DB_DEBUG` environment variable to a
// non-empty value, all generated statements will be printed at runtime to
// the standard logger.
//
// Example:
//
//	UPPERIO_DB_DEBUG=1 go test
//
//	UPPERIO_DB_DEBUG=1 ./go-program
const (
	EnvEnableDebug = `UPPERIO_DB_DEBUG`
)

// Logger represents a logging collector. You can pass a logging collector to
// db.DefaultSettings.SetLogger(myCollector) to make it collect db.QueryStatus messages
// after executing a query.
type Logger interface {
	Log(*QueryStatus)
}

type defaultLogger struct {
}

func (lg *defaultLogger) Log(q *QueryStatus) {
	log.Println("\n\t" + q.Stringify("\n\t") + "\n")
}

var _ = Logger(&defaultLogger{})

func init() {
	if envEnabled(EnvEnableDebug) {
		DefaultSettings.SetLogging(true)
	}
}
