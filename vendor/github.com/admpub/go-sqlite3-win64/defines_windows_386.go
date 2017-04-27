// Copyright (C) 2016 Samuel Melrose <sam@infitialis.com>.
//
// Based on work by Yasuhiro Matsumoto <mattn.jp@gmail.com>
// https://github.com/mattn/go-sqlite3
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package sqlite3

import "math"

const (
	SQLITE_TRANSIENT = math.MaxUint32 // Can't do -1 for overflow like in C, so use largest unsigned 32bit int.
)
