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
    "strings"
    "runtime"
    "database/sql/driver"
    "golang.org/x/net/context"
)

// Exec implements Execer.
func (c *SQLiteConn) Exec(query string, args []driver.Value) (driver.Result, error) {
    list := make([]namedValue, len(args))
    for i, v := range args {
        list[i] = namedValue{
            Ordinal: i + 1,
            Value:   v,
        }
    }
    return c.exec(context.Background(), query, list)
}

func (c *SQLiteConn) exec(ctx context.Context, query string, args []namedValue) (driver.Result, error) {
    start := 0
    for {
        s, err := c.prepare(ctx, query)
        if err != nil {
            return nil, err
        }
        var res driver.Result
        if s.(*SQLiteStmt).s != sqlite3_stmt(uintptr(0)) {
            na := s.NumInput()
            if len(args) < na {
                return nil, fmt.Errorf("Not enough args to execute query. Expected %d, got %d.", na, len(args))
            }
            for i := 0; i < na; i++ {
                args[i].Ordinal -= start
            }
            res, err = s.(*SQLiteStmt).exec(ctx, args[:na])
            if err != nil && err != driver.ErrSkip {
                s.Close()
                return nil, err
            }
            args = args[na:]
            start += na
        }
        tail := s.(*SQLiteStmt).t
        s.Close()
        if tail == "" {
            return res, nil
        }
        query = tail
    }
}

// Query implements Queryer.
func (c *SQLiteConn) Query(query string, args []driver.Value) (driver.Rows, error) {
    list := make([]namedValue, len(args))
    for i, v := range args {
        list[i] = namedValue{
            Ordinal: i + 1,
            Value:   v,
        }
    }
    return c.query(context.Background(), query, list)
}

func (c *SQLiteConn) query(ctx context.Context, query string, args []namedValue) (driver.Rows, error) {
    start := 0
    for {
        s, err := c.prepare(ctx, query)
        if err != nil {
            return nil, err
        }
        s.(*SQLiteStmt).cls = true
        na := s.NumInput()
        if len(args) < na {
            return nil, fmt.Errorf("Not enough args to execute query. Expected %d, got %d.", na, len(args))
        }
        for i := 0; i < na; i++ {
            args[i].Ordinal -= start
        }
        rows, err := s.(*SQLiteStmt).query(ctx, args[:na])
        if err != nil && err != driver.ErrSkip {
            s.Close()
            return rows, err
        }
        args = args[na:]
        start += na
        tail := s.(*SQLiteStmt).t
        if tail == "" {
            return rows, nil
        }
        rows.Close()
        s.Close()
        query = tail
    }
}

// Begin transaction.
func (c *SQLiteConn) Begin() (driver.Tx, error) {
    return c.begin(context.Background())
}

func (c *SQLiteConn) begin(ctx context.Context) (driver.Tx, error) {
    if _, err := c.exec(ctx, c.txlock, nil); err != nil {
        return nil, err
    }
    return &SQLiteTx{c}, nil
}

// Prepare the query string. Return a new statement.
func (c *SQLiteConn) Prepare(query string) (driver.Stmt, error) {
    return c.prepare(context.Background(), query)
}

func (c *SQLiteConn) prepare(ctx context.Context, query string) (driver.Stmt, error) {
    rv, s, tail := sqlite3_prepare_v2(c.db, query)
    if rv != SQLITE_OK {
        return nil, c.lastError()
    }
    t := strings.TrimSpace(tail)
    ss := &SQLiteStmt{c: c, s: s, t: t}
    runtime.SetFinalizer(ss, (*SQLiteStmt).Close)
    return ss, nil
}

// AutoCommit return which currently auto commit or not.
func (c *SQLiteConn) AutoCommit() bool {
    return sqlite3_get_autocommit(c.db) != 0
}

// Close the connection.
func (c *SQLiteConn) Close() error {
    rv := sqlite3_close_v2(c.db)
    if rv != SQLITE_OK {
        return c.lastError()
    }
    c.db = sqlite3(uintptr(0))
    runtime.SetFinalizer(c, nil)
    return nil
}

func (c *SQLiteConn) lastError() Error {
    return Error{
        Code:         ErrNo(sqlite3_errcode(c.db)),
        ExtendedCode: ErrNoExtended(sqlite3_extended_errcode(c.db)),
        err:          sqlite3_errmsg(c.db),
    }
}
