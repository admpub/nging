// +build windows
// Copyright (C) 2016 Samuel Melrose <sam@infitialis.com>.
//
// Based on work by Yasuhiro Matsumoto <mattn.jp@gmail.com>
// https://github.com/mattn/go-sqlite3
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package sqlite3

import (
	"database/sql/driver"
	"io"
	"strings"
	"time"
)

// Close the rows.
func (rc *SQLiteRows) Close() error {
	if rc.s.closed {
		return nil
	}
	if rc.done != nil {
		close(rc.done)
	}
	if rc.cls {
		return rc.s.Close()
	}
	rv := sqlite3_reset(rc.s.s)
	if rv != SQLITE_OK {
		return rc.s.c.lastError()
	}
	return nil
}

// Columns return column names.
func (rc *SQLiteRows) Columns() []string {
	if rc.nc != len(rc.cols) {
		rc.cols = make([]string, rc.nc)
		for i := 0; i < rc.nc; i++ {
			rc.cols[i] = sqlite3_column_name(rc.s.s, i)
		}
	}
	return rc.cols
}

// DeclTypes return column types.
func (rc *SQLiteRows) DeclTypes() []string {
	if rc.decltype == nil {
		rc.decltype = make([]string, rc.nc)
		for i := 0; i < rc.nc; i++ {
			rc.decltype[i] = strings.ToLower(sqlite3_column_decltype(rc.s.s, i))
		}
	}
	return rc.decltype
}

// Next move cursor to next.
func (rc *SQLiteRows) Next(dest []driver.Value) error {
	var rowid, changes int64
	rv := sqlite3_step(rc.s.s, &rowid, &changes)
	if rv == SQLITE_DONE {
		return io.EOF
	}
	if rv != SQLITE_ROW {
		rv = sqlite3_reset(rc.s.s)
		if rv != SQLITE_OK {
			return rc.s.c.lastError()
		}
		return nil
	}

	rc.DeclTypes()

	for i := range dest {
		switch sqlite3_column_type(rc.s.s, i) {
		case SQLITE_INTEGER:
			val := sqlite3_column_int64(rc.s.s, i)
			switch rc.decltype[i] {
			case "timestamp", "datetime", "date":
				var t time.Time
				// Assume a millisecond unix timestamp if it's 13 digits -- too
				// large to be a reasonable timestamp in seconds.
				if val > 1e12 || val < -1e12 {
					val *= int64(time.Millisecond) // convert ms to nsec
				} else {
					val *= int64(time.Second) // convert sec to nsec
				}
				t = time.Unix(0, val).UTC()
				if rc.s.c.loc != nil {
					t = t.In(rc.s.c.loc)
				}
				dest[i] = t
			case "boolean":
				dest[i] = val > 0
			default:
				dest[i] = val
			}
		case SQLITE_FLOAT:
			dest[i] = sqlite3_column_double(rc.s.s, i)
		case SQLITE_BLOB:
			dest[i] = sqlite3_column_blob(rc.s.s, i)
		case SQLITE_NULL:
			dest[i] = nil
		case SQLITE_TEXT:
			var err error
			var timeVal time.Time

			s := sqlite3_column_text(rc.s.s, i)

			switch rc.decltype[i] {
			case "timestamp", "datetime", "date":
				var t time.Time
				s = strings.TrimSuffix(s, "Z")
				for _, format := range SQLiteTimestampFormats {
					if timeVal, err = time.ParseInLocation(format, s, time.UTC); err == nil {
						t = timeVal
						break
					}
				}
				if err != nil {
					// The column is a time value, so return the zero time on parse failure.
					t = time.Time{}
				}
				if rc.s.c.loc != nil {
					t = t.In(rc.s.c.loc)
				}
				dest[i] = t
			default:
				dest[i] = []byte(s)
			}
		}
	}
	return nil
}
