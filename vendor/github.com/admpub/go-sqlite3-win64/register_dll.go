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
	"fmt"
	"path/filepath"
	"runtime"
	"syscall"
)

var (
	DLLPath       string
	dllRegistered bool = false

	modSQLite3 *syscall.LazyDLL

	dll_sqlite3_db_handle         *syscall.LazyProc
	dll_sqlite3_last_insert_rowid *syscall.LazyProc
	dll_sqlite3_changes           *syscall.LazyProc

	dll_sqlite3_libversion        *syscall.LazyProc
	dll_sqlite3_libversion_number *syscall.LazyProc
	dll_sqlite3_sourceid          *syscall.LazyProc

	dll_sqlite3_errstr           *syscall.LazyProc
	dll_sqlite3_errcode          *syscall.LazyProc
	dll_sqlite3_extended_errcode *syscall.LazyProc
	dll_sqlite3_errmsg           *syscall.LazyProc
	dll_sqlite3_threadsafe       *syscall.LazyProc

	dll_sqlite3_open_v2      *syscall.LazyProc
	dll_sqlite3_busy_timeout *syscall.LazyProc
	dll_sqlite3_close_v2     *syscall.LazyProc

	dll_sqlite3_prepare_v2     *syscall.LazyProc
	dll_sqlite3_get_autocommit *syscall.LazyProc

	dll_sqlite3_finalize             *syscall.LazyProc
	dll_sqlite3_bind_parameter_count *syscall.LazyProc
	dll_sqlite3_bind_parameter_index *syscall.LazyProc
	dll_sqlite3_reset                *syscall.LazyProc
	dll_sqlite3_bind_null            *syscall.LazyProc
	dll_sqlite3_bind_int64           *syscall.LazyProc
	dll_sqlite3_bind_int             *syscall.LazyProc
	dll_sqlite3_bind_text            *syscall.LazyProc
	dll_sqlite3_bind_double          *syscall.LazyProc
	dll_sqlite3_bind_blob            *syscall.LazyProc
	dll_sqlite3_column_count         *syscall.LazyProc
	dll_sqlite3_column_name          *syscall.LazyProc
	dll_sqlite3_interrupt            *syscall.LazyProc
	dll_sqlite3_clear_bindings       *syscall.LazyProc

	dll_sqlite3_step            *syscall.LazyProc
	dll_sqlite3_column_decltype *syscall.LazyProc
	dll_sqlite3_column_type     *syscall.LazyProc
	dll_sqlite3_column_int64    *syscall.LazyProc
	dll_sqlite3_column_double   *syscall.LazyProc
	dll_sqlite3_column_bytes    *syscall.LazyProc
	dll_sqlite3_column_blob     *syscall.LazyProc
	dll_sqlite3_column_text     *syscall.LazyProc
)

func registerDLLFunctions() {
	if err := registerDLL(); err != nil {
		dllRegistered = false
		return
	}

	dll_sqlite3_db_handle = modSQLite3.NewProc("sqlite3_db_handle")
	dll_sqlite3_last_insert_rowid = modSQLite3.NewProc("sqlite3_last_insert_rowid")
	dll_sqlite3_changes = modSQLite3.NewProc("sqlite3_changes")

	dll_sqlite3_libversion = modSQLite3.NewProc("sqlite3_libversion")
	dll_sqlite3_libversion_number = modSQLite3.NewProc("sqlite3_libversion_number")
	dll_sqlite3_sourceid = modSQLite3.NewProc("sqlite3_sourceid")

	dll_sqlite3_errstr = modSQLite3.NewProc("sqlite3_errstr")
	dll_sqlite3_errcode = modSQLite3.NewProc("sqlite3_errcode")
	dll_sqlite3_extended_errcode = modSQLite3.NewProc("sqlite3_extended_errcode")
	dll_sqlite3_errmsg = modSQLite3.NewProc("sqlite3_errmsg")
	dll_sqlite3_threadsafe = modSQLite3.NewProc("sqlite3_threadsafe")

	dll_sqlite3_open_v2 = modSQLite3.NewProc("sqlite3_open_v2")
	dll_sqlite3_busy_timeout = modSQLite3.NewProc("sqlite3_busy_timeout")
	dll_sqlite3_close_v2 = modSQLite3.NewProc("sqlite3_close_v2")

	dll_sqlite3_prepare_v2 = modSQLite3.NewProc("sqlite3_prepare_v2")
	dll_sqlite3_get_autocommit = modSQLite3.NewProc("sqlite3_get_autocommit")

	dll_sqlite3_finalize = modSQLite3.NewProc("sqlite3_finalize")
	dll_sqlite3_bind_parameter_count = modSQLite3.NewProc("sqlite3_bind_parameter_count")
	dll_sqlite3_bind_parameter_index = modSQLite3.NewProc("sqlite3_bind_parameter_index")
	dll_sqlite3_reset = modSQLite3.NewProc("sqlite3_reset")
	dll_sqlite3_bind_null = modSQLite3.NewProc("sqlite3_bind_null")
	dll_sqlite3_bind_int64 = modSQLite3.NewProc("sqlite3_bind_int64")
	dll_sqlite3_bind_int = modSQLite3.NewProc("sqlite3_bind_int")
	dll_sqlite3_bind_text = modSQLite3.NewProc("sqlite3_bind_text")
	dll_sqlite3_bind_double = modSQLite3.NewProc("sqlite3_bind_double")
	dll_sqlite3_bind_blob = modSQLite3.NewProc("sqlite3_bind_blob")
	dll_sqlite3_column_count = modSQLite3.NewProc("sqlite3_column_count")
	dll_sqlite3_column_name = modSQLite3.NewProc("sqlite3_column_name")
	dll_sqlite3_interrupt = modSQLite3.NewProc("sqlite3_interrupt")
	dll_sqlite3_clear_bindings = modSQLite3.NewProc("sqlite3_clear_bindings")

	dll_sqlite3_step = modSQLite3.NewProc("sqlite3_step")
	dll_sqlite3_column_decltype = modSQLite3.NewProc("sqlite3_column_decltype")
	dll_sqlite3_column_type = modSQLite3.NewProc("sqlite3_column_type")
	dll_sqlite3_column_int64 = modSQLite3.NewProc("sqlite3_column_int64")
	dll_sqlite3_column_double = modSQLite3.NewProc("sqlite3_column_double")
	dll_sqlite3_column_bytes = modSQLite3.NewProc("sqlite3_column_bytes")
	dll_sqlite3_column_blob = modSQLite3.NewProc("sqlite3_column_blob")
	dll_sqlite3_column_text = modSQLite3.NewProc("sqlite3_column_text")

	dllRegistered = true
}

func registerDLL() error {
	dllName := "sqlite3_" + runtime.GOARCH + ".dll"
	dllPath := DLLPath
	if len(dllPath) == 0 {
		dllPath = basePath()
	}
	filePath := filepath.Join(dllPath, dllName)
	if exist, _ := exists(filePath); exist {
		modSQLite3 = syscall.NewLazyDLL(filePath)
		return nil
	}

	filePath = filepath.Join(dllPath, "support", dllName)
	if exist, _ := exists(filePath); exist {
		modSQLite3 = syscall.NewLazyDLL(filePath)
		return nil
	}

	return fmt.Errorf("%s not found.", dllName)
}
