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
	"unsafe"
)

func sqlite3_libversion() string {
	msgPtr, _, _ := dll_sqlite3_libversion.Call()
	if msgPtr == uintptr(0) {
		return "Unknown"
	}
	return BytePtrToString((*byte)(unsafe.Pointer(msgPtr)))
}

func sqlite3_libversion_number() int {
	intPtr, _, _ := dll_sqlite3_libversion_number.Call()
	if intPtr == uintptr(0) {
		return 0
	}
	return int(intPtr)
}

func sqlite3_sourceid() string {
	msgPtr, _, _ := dll_sqlite3_sourceid.Call()
	if msgPtr == uintptr(0) {
		return "Unknown"
	}
	return BytePtrToString((*byte)(unsafe.Pointer(msgPtr)))
}

func sqlite3_errstr(code int) (msg string) {
	msgPtr, _, _ := dll_sqlite3_errstr.Call(uintptr(code))
	if msgPtr == uintptr(0) {
		return "Unknown Error"
	}
	return BytePtrToString((*byte)(unsafe.Pointer(msgPtr)))
}

func sqlite3_errcode(db sqlite3) (code int) {
	retInt, _, _ := dll_sqlite3_errcode.Call(uintptr(db))
	return int(retInt)
}

func sqlite3_extended_errcode(db sqlite3) (code int) {
	retInt, _, _ := dll_sqlite3_extended_errcode.Call(uintptr(db))
	return int(retInt)
}

func sqlite3_errmsg(db sqlite3) (msg string) {
	msgPtr, _, _ := dll_sqlite3_errmsg.Call(uintptr(db))
	if msgPtr == uintptr(0) {
		return "Unknown Error"
	}
	return BytePtrToString((*byte)(unsafe.Pointer(msgPtr)))
}

func sqlite3_threadsafe() int {
	intPtr, _, _ := dll_sqlite3_threadsafe.Call()
	return int(intPtr)
}

func sqlite3_open_v2(filename string, ppDb *sqlite3, flags int, zVfs string) int {
	var fn uintptr
	if len(filename) > 0 {
		var tfn = []byte(filename)
		fn = uintptr(unsafe.Pointer(&tfn[0]))
	} else {
		fn = uintptr(0)
	}
	var vfs uintptr
	if len(zVfs) > 0 {
		var tvfs = []byte(zVfs)
		vfs = uintptr(unsafe.Pointer(&tvfs[0]))
	} else {
		vfs = uintptr(0)
	}
	retInt, _, _ := dll_sqlite3_open_v2.Call(
		fn,
		uintptr(unsafe.Pointer(ppDb)),
		uintptr(flags),
		vfs,
	)
	return int(retInt)
}

func sqlite3_busy_timeout(db sqlite3, busyTimeout int) int {
	retInt, _, _ := dll_sqlite3_busy_timeout.Call(uintptr(db), uintptr(busyTimeout))
	return int(retInt)
}

func sqlite3_close_v2(db sqlite3) int {
	retInt, _, _ := dll_sqlite3_close_v2.Call(uintptr(db))
	return int(retInt)
}

func sqlite3_prepare_v2(db sqlite3, zSql string) (retCode int, stmtHandle sqlite3_stmt, tail string) {
	var sql = []byte(zSql + "\x00")
	var handle uintptr
	var thandle uintptr
	retInt, _, _ := dll_sqlite3_prepare_v2.Call(
		uintptr(db),
		uintptr(unsafe.Pointer(&sql[0])),
		uintptr(SQLITE_TRANSIENT),
		uintptr(unsafe.Pointer(&handle)),
		uintptr(unsafe.Pointer(&thandle)),
	)

	return int(retInt), sqlite3_stmt(handle), BytePtrToString((*byte)(unsafe.Pointer(thandle)))
}

func sqlite3_get_autocommit(db sqlite3) int {
	retInt, _, _ := dll_sqlite3_get_autocommit.Call(uintptr(db))
	return int(retInt)
}

func sqlite3_finalize(stmt sqlite3_stmt) int {
	retInt, _, _ := dll_sqlite3_finalize.Call(uintptr(stmt))
	return int(retInt)
}

func sqlite3_bind_parameter_count(stmt sqlite3_stmt) int {
	retInt, _, _ := dll_sqlite3_bind_parameter_count.Call(uintptr(stmt))
	return int(retInt)
}

func sqlite3_bind_parameter_index(stmt sqlite3_stmt, name string) int {
	var pName = []byte(name)
	retInt, _, _ := dll_sqlite3_bind_parameter_index.Call(
		uintptr(stmt),
		uintptr(unsafe.Pointer(&pName[0])),
	)
	return int(retInt)
}

func sqlite3_reset(stmt sqlite3_stmt) int {
	retInt, _, _ := dll_sqlite3_reset.Call(uintptr(stmt))
	return int(retInt)
}

func sqlite3_bind_null(stmt sqlite3_stmt, ord int) int {
	retInt, _, _ := dll_sqlite3_bind_null.Call(
		uintptr(stmt),
		uintptr(ord),
	)
	return int(retInt)
}

func sqlite3_bind_text(stmt sqlite3_stmt, ord int, data string) int {
	var b []byte
	if len(data) == 0 {
		b = []byte{0}
	} else {
		b = []byte(data)
	}
	retInt, _, _ := dll_sqlite3_bind_text.Call(
		uintptr(stmt),
		uintptr(ord),
		uintptr(unsafe.Pointer(&b[0])),
		uintptr(len(data)),
		uintptr(SQLITE_TRANSIENT),
	)
	return int(retInt)
}

func sqlite3_bind_int64(stmt sqlite3_stmt, ord int, data int64) int {
	retInt, _, _ := dll_sqlite3_bind_int64.Call(
		uintptr(stmt),
		uintptr(ord),
		uintptr(data),
	)
	return int(retInt)
}

func sqlite3_bind_int(stmt sqlite3_stmt, ord int, data int) int {
	retInt, _, _ := dll_sqlite3_bind_int.Call(
		uintptr(stmt),
		uintptr(ord),
		uintptr(data),
	)
	return int(retInt)
}

func sqlite3_bind_double(stmt sqlite3_stmt, ord int, data float64) int {
	retInt, _, _ := dll_sqlite3_bind_double.Call(
		uintptr(stmt),
		uintptr(ord),
		uintptr(data),
	)
	return int(retInt)
}

func sqlite3_bind_blob(stmt sqlite3_stmt, ord int, data []byte) int {
	var pData uintptr
	if len(data) == 0 {
		pData = 0
	} else {
		pData = uintptr(unsafe.Pointer(&data[0]))
	}
	retInt, _, _ := dll_sqlite3_bind_blob.Call(
		uintptr(stmt),
		uintptr(ord),
		pData,
		uintptr(len(data)),
		uintptr(SQLITE_TRANSIENT),
	)
	return int(retInt)
}

func sqlite3_column_count(stmt sqlite3_stmt) int {
	retInt, _, _ := dll_sqlite3_column_count.Call(uintptr(stmt))
	return int(retInt)
}

func sqlite3_column_name(stmt sqlite3_stmt, index int) string {
	msgPtr, _, _ := dll_sqlite3_column_name.Call(
		uintptr(stmt),
		uintptr(index),
	)
	return BytePtrToString((*byte)(unsafe.Pointer(msgPtr)))
}

func sqlite3_interrupt(db sqlite3) {
	dll_sqlite3_interrupt.Call(uintptr(db))
}

func sqlite3_clear_bindings(stmt sqlite3_stmt) {
	dll_sqlite3_clear_bindings.Call(uintptr(stmt))
}

func sqlite3_step(stmt sqlite3_stmt, rowid *int64, changes *int64) int {
	retInt, _, _ := dll_sqlite3_step.Call(
		uintptr(stmt),
		uintptr(unsafe.Pointer(rowid)),
		uintptr(unsafe.Pointer(changes)),
	)
	return int(retInt)
}

func sqlite3_column_decltype(stmt sqlite3_stmt, index int) string {
	msgPtr, _, _ := dll_sqlite3_column_decltype.Call(
		uintptr(stmt),
		uintptr(index),
	)
	if msgPtr == uintptr(0) {
		return ""
	}
	return BytePtrToString((*byte)(unsafe.Pointer(msgPtr)))
}

func sqlite3_column_type(stmt sqlite3_stmt, index int) int {
	retInt, _, _ := dll_sqlite3_column_type.Call(
		uintptr(stmt),
		uintptr(index),
	)
	return int(retInt)
}

func sqlite3_column_int64(stmt sqlite3_stmt, index int) int64 {
	intPtr, _, _ := dll_sqlite3_column_int64.Call(
		uintptr(stmt),
		uintptr(index),
	)
	return int64(intPtr)
}

func sqlite3_column_double(stmt sqlite3_stmt, index int) float64 {
	intPtr, _, _ := dll_sqlite3_column_double.Call(
		uintptr(stmt),
		uintptr(index),
	)
	return float64(intPtr)
}

func sqlite3_column_bytes(stmt sqlite3_stmt, index int) int {
	intPtr, _, _ := dll_sqlite3_column_bytes.Call(
		uintptr(stmt),
		uintptr(index),
	)
	return int(intPtr)
}

func sqlite3_column_blob(stmt sqlite3_stmt, index int) []byte {
	bytesPtr, _, _ := dll_sqlite3_column_blob.Call(
		uintptr(stmt),
		uintptr(index),
	)

	n := sqlite3_column_bytes(stmt, index)

	slice := make([]byte, n)
	copy(slice[:], (*[1 << 30]byte)(unsafe.Pointer(bytesPtr))[0:n])
	return slice
}

func sqlite3_column_text(stmt sqlite3_stmt, index int) string {
	bytesPtr, _, _ := dll_sqlite3_column_text.Call(
		uintptr(stmt),
		uintptr(index),
	)

	n := sqlite3_column_bytes(stmt, index)

	slice := make([]byte, n)
	copy(slice[:], (*[1 << 30]byte)(unsafe.Pointer(bytesPtr))[0:n])
	return string(slice)
}
