// +build sqlite
// +build !windows windows,386

package sqlite

import (
	_ "github.com/mattn/go-sqlite3" //sqlite driver
)
