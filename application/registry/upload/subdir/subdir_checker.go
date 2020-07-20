package subdir

import (
	"github.com/admpub/log"
	"github.com/admpub/nging/application/registry/upload/checker"
	"github.com/fatih/color"
)

func CheckerRegister(typ string, checkerFn checker.Checker, fieldNames ...string) {
	params := ParseUploadType(typ)
	subdir := params.MustGetSubdir()
	log.Info(color.GreenString(`checker.register:`), typ)
	if len(fieldNames) > 0 {
		GetOrCreate(subdir).SetChecker(checkerFn, fieldNames...)
		return
	}
	GetOrCreate(subdir).SetChecker(checkerFn, params.Field)
}

func CheckerGet(subdir string, defaults ...string) checker.Checker {
	s := Get(subdir)
	if s != nil {
		return s.MustChecker()
	}
	if len(defaults) == 0 {
		return checker.Default
	}
	return CheckerGet(defaults[0])
}
