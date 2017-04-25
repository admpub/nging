// +build windows,amd64
// Copyright (C) 2016 Samuel Melrose <sam@infitialis.com>.
//
// Based on work by Yasuhiro Matsumoto <mattn.jp@gmail.com>
// https://github.com/mattn/go-sqlite3
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package sqlite3

import (
    "fmt"
    "time"
    "runtime"
    "database/sql/driver"
    "golang.org/x/net/context"
)

// Close the statement.
func (s *SQLiteStmt) Close() error {
    if s.closed {
        return nil
    }
    s.closed = true
    if s.c == nil || s.c.db == sqlite3(uintptr(0)) {
        return fmt.Errorf("sqlite statement with already closed database connection")
    }
    rv := sqlite3_finalize(s.s)
    if rv != SQLITE_OK {
        return s.c.lastError()
    }
    runtime.SetFinalizer(s, nil)
    return nil
}

// NumInput return a number of parameters.
func (s *SQLiteStmt) NumInput() int {
    return sqlite3_bind_parameter_count(s.s)
}

func (s *SQLiteStmt) bind(args []namedValue) error {
    rv := sqlite3_reset(s.s)
    if rv != SQLITE_ROW && rv != SQLITE_OK && rv != SQLITE_DONE {
        return s.c.lastError()
    }

    for i, v := range args {
        if v.Name != "" {
            args[i].Ordinal = sqlite3_bind_parameter_index(s.s, v.Name)
        }
    }

    for _, arg := range args {
        n := arg.Ordinal
        switch v := arg.Value.(type) {
            case nil:
                rv = sqlite3_bind_null(s.s, n)
            case string:
                rv = sqlite3_bind_text(s.s, n, v)
            case int64:
                rv = sqlite3_bind_int64(s.s, n, v)
            case bool:
                if bool(v) {
                    rv = sqlite3_bind_int(s.s, n, 1)
                } else {
                    rv = sqlite3_bind_int(s.s, n, 0)
                }
            case float64:
                rv = sqlite3_bind_double(s.s, n, v)
            case []byte:
                rv = sqlite3_bind_blob(s.s, n, v)
            case time.Time:
                b := v.Format(SQLiteTimestampFormats[0])
                rv = sqlite3_bind_text(s.s, n, b)
        }
        if rv != SQLITE_OK {
            return s.c.lastError()
        }
    }
    return nil
}

// Query the statement with arguments. Return records.
func (s *SQLiteStmt) Query(args []driver.Value) (driver.Rows, error) {
    list := make([]namedValue, len(args))
    for i, v := range args {
        list[i] = namedValue{
            Ordinal: i + 1,
            Value:   v,
        }
    }
    return s.query(context.Background(), list)
}

func (s *SQLiteStmt) query(ctx context.Context, args []namedValue) (driver.Rows, error) {
    if err := s.bind(args); err != nil {
        return nil, err
    }

    rows := &SQLiteRows{
        s:        s,
        nc:       sqlite3_column_count(s.s),
        cols:     nil,
        decltype: nil,
        cls:      s.cls,
        done:     make(chan struct{}),
    }

    go func() {
        select {
            case <-ctx.Done():
                sqlite3_interrupt(s.c.db)
                rows.Close()
            case <-rows.done:
        }
    }()

    return rows, nil
}

// LastInsertId teturn last inserted ID.
func (r *SQLiteResult) LastInsertId() (int64, error) {
    return r.id, nil
}

// RowsAffected return how many rows affected.
func (r *SQLiteResult) RowsAffected() (int64, error) {
    return r.changes, nil
}

// Exec execute the statement with arguments. Return result object.
func (s *SQLiteStmt) Exec(args []driver.Value) (driver.Result, error) {
    list := make([]namedValue, len(args))
    for i, v := range args {
        list[i] = namedValue{
            Ordinal: i + 1,
            Value:   v,
        }
    }
    return s.exec(context.Background(), list)
}

func (s *SQLiteStmt) exec(ctx context.Context, args []namedValue) (driver.Result, error) {
    if err := s.bind(args); err != nil {
        sqlite3_reset(s.s)
        sqlite3_clear_bindings(s.s)
        return nil, err
    }

    done := make(chan struct{})
    defer close(done)
    go func() {
        select {
            case <-ctx.Done():
                sqlite3_interrupt(s.c.db)
            case <-done:
        }
    }()

    var rowid, changes int64
    rv := sqlite3_step(s.s, &rowid, &changes)
    if rv != SQLITE_ROW && rv != SQLITE_OK && rv != SQLITE_DONE {
        err := s.c.lastError()
        sqlite3_reset(s.s)
        sqlite3_clear_bindings(s.s)
        return nil, err
    }

    return &SQLiteResult{id: rowid, changes: changes}, nil
}
