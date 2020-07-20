package subdir

import (
	"strings"

	"github.com/admpub/log"
	"github.com/admpub/nging/application/registry/upload/checker"
	"github.com/admpub/nging/application/registry/upload/table"
	"github.com/fatih/color"
)

func CheckerRegister(typ string, checkerFn checker.Checker, fieldNames ...string) {
	log.Info(color.GreenString(`checker.register:`), typ)
	if len(fieldNames) > 0 {
		Get(typ).SetChecker(checkerFn, fieldNames...)
		return
	}
	tableName, fieldName, _ := table.GetTableInfo(typ)
	info := Get(tableName)
	info.SetChecker(checkerFn, fieldName)
}

func CheckerGet(typ string, defaults ...string) checker.Checker {
	s := Get(typ)
	if s != nil {
		return s.MustChecker()
	}
	if len(defaults) == 0 {
		tmp := strings.SplitN(typ, `.`, 2)
		if len(tmp) == 2 {
			return CheckerGet(tmp[0])
		}
		return checker.Default
	}
	return CheckerGet(defaults[0])
}
