package upload

import (
	"strings"

	"github.com/webx-top/echo"
)

var Subdir = echo.NewKVData()

func init() {
	Subdir.Add(`default`, `默认`)
	Subdir.Add(`avatar`, `头像`)
}

func AllowedSubdir(subdir string) bool {
	parts := strings.SplitN(subdir, `/`, 2)
	if len(parts) != 2 {
		return Subdir.Has(subdir)
	}
	item := Subdir.GetItem(parts[0])
	if item == nil || item.H == nil {
		return false
	}

	return item.H.Has(parts[1])
}
