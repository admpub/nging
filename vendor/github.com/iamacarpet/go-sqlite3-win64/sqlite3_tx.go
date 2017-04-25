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
    "golang.org/x/net/context"
)

// Commit transaction.
func (tx *SQLiteTx) Commit() error {
    _, err := tx.c.exec(context.Background(), "COMMIT", nil)
    if err != nil && err.(Error).Code == SQLITE_BUSY {
        // sqlite3 will leave the transaction open in this scenario.
        // However, database/sql considers the transaction complete once we
        // return from Commit() - we must clean up to honour its semantics.
        tx.c.exec(context.Background(), "ROLLBACK", nil)
    }
    return err
}

// Rollback transaction.
func (tx *SQLiteTx) Rollback() error {
    _, err := tx.c.exec(context.Background(), "ROLLBACK", nil)
    return err
}
