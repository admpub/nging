// +build sqlite
// +build !windll

package sqlite

import (
	_ "github.com/mattn/go-sqlite3" //sqlite driver
)

func init() {
	register()
}
