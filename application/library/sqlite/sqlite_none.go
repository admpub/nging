// +build !sqlite

package sqlite

import (
	_ "github.com/admpub/nging/application/library/config"
)

func init() {
	//installer.ExecSQL = ExecSQL
}
