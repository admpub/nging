// +build sqlite
// +build !windows windows,i386

package sqlite

import (
	_ "github.com/mattn/go-sqlite3" //sqlite driver
)
