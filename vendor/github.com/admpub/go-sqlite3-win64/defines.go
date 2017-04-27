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
	"time"
)

// SQLiteTimestampFormats is timestamp formats understood by both this module
// and SQLite.  The first format in the slice will be used when saving time
// values into the database. When parsing a string from a timestamp or datetime
// column, the formats are tried in order.
var SQLiteTimestampFormats = []string{
	// By default, store timestamps with whatever timezone they come with.
	// When parsed, they will be returned with the same timezone.
	"2006-01-02 15:04:05.999999999-07:00",
	"2006-01-02T15:04:05.999999999-07:00",
	"2006-01-02 15:04:05.999999999",
	"2006-01-02T15:04:05.999999999",
	"2006-01-02 15:04:05",
	"2006-01-02T15:04:05",
	"2006-01-02 15:04",
	"2006-01-02T15:04",
	"2006-01-02",
}

const (
	SQLITE_OK                  = 0
	SQLITE_OPEN_READONLY       = 0x00000001 /* Ok for sqlite3_open_v2() */
	SQLITE_OPEN_READWRITE      = 0x00000002 /* Ok for sqlite3_open_v2() */
	SQLITE_OPEN_CREATE         = 0x00000004 /* Ok for sqlite3_open_v2() */
	SQLITE_OPEN_DELETEONCLOSE  = 0x00000008 /* VFS only */
	SQLITE_OPEN_EXCLUSIVE      = 0x00000010 /* VFS only */
	SQLITE_OPEN_AUTOPROXY      = 0x00000020 /* VFS only */
	SQLITE_OPEN_URI            = 0x00000040 /* Ok for sqlite3_open_v2() */
	SQLITE_OPEN_MEMORY         = 0x00000080 /* Ok for sqlite3_open_v2() */
	SQLITE_OPEN_MAIN_DB        = 0x00000100 /* VFS only */
	SQLITE_OPEN_TEMP_DB        = 0x00000200 /* VFS only */
	SQLITE_OPEN_TRANSIENT_DB   = 0x00000400 /* VFS only */
	SQLITE_OPEN_MAIN_JOURNAL   = 0x00000800 /* VFS only */
	SQLITE_OPEN_TEMP_JOURNAL   = 0x00001000 /* VFS only */
	SQLITE_OPEN_SUBJOURNAL     = 0x00002000 /* VFS only */
	SQLITE_OPEN_MASTER_JOURNAL = 0x00004000 /* VFS only */
	SQLITE_OPEN_NOMUTEX        = 0x00008000 /* Ok for sqlite3_open_v2() */
	SQLITE_OPEN_FULLMUTEX      = 0x00010000 /* Ok for sqlite3_open_v2() */
	SQLITE_OPEN_SHAREDCACHE    = 0x00020000 /* Ok for sqlite3_open_v2() */
	SQLITE_OPEN_PRIVATECACHE   = 0x00040000 /* Ok for sqlite3_open_v2() */
	SQLITE_OPEN_WAL            = 0x00080000 /* VFS only */
	SQLITE_STATIC              = 0

	SQLITE_BUSY = 5
	SQLITE_DONE = 101
	SQLITE_ROW  = 100

	SQLITE_INTEGER = 1
	SQLITE_FLOAT   = 2
	SQLITE_BLOB    = 4
	SQLITE_NULL    = 5
	SQLITE_TEXT    = 3
)

type sqlite3 uintptr
type sqlite3_stmt uintptr

// SQLiteDriver implement sql.Driver.
type SQLiteDriver struct {
	Extensions  []string
	ConnectHook func(*SQLiteConn) error
}

// SQLiteConn implement sql.Conn.
type SQLiteConn struct {
	db     sqlite3
	loc    *time.Location
	txlock string
}

// SQLiteTx implemen sql.Tx.
type SQLiteTx struct {
	c *SQLiteConn
}

// SQLiteStmt implement sql.Stmt.
type SQLiteStmt struct {
	c      *SQLiteConn
	s      sqlite3_stmt
	t      string
	closed bool
	cls    bool
}

// SQLiteResult implement sql.Result.
type SQLiteResult struct {
	id      int64
	changes int64
}

// SQLiteRows implement sql.Rows.
type SQLiteRows struct {
	s        *SQLiteStmt
	nc       int
	cols     []string
	decltype []string
	cls      bool
	done     chan struct{}
}

type namedValue struct {
	Name    string
	Ordinal int
	Value   driver.Value
}

type bindArg struct {
	n int
	v driver.Value
}
