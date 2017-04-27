// +build sqlite
// +build windows,amd64

package sqlite

import (
	_ "github.com/iamacarpet/go-sqlite3-win64" //sqlite driver
	"github.com/webx-top/db/sqlite"
)

func init() {
	sqlite.FixFilePath = func(file string) string {
		return strings.TrimPrefix(file, `file:///`)
	}
}
